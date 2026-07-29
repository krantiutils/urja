package invoice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository owns all invoice SQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new invoice repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// sellerIdentity is the gym's tax identity at the moment of issue.
type sellerIdentity struct {
	Name          string
	PAN           string
	Address       string
	VATRegistered bool
}

// loadSeller reads the org's tax identity, falling back to its display name and
// address when the registered ones are blank.
func (r *Repository) loadSeller(ctx context.Context, tx pgx.Tx, orgID string) (sellerIdentity, error) {
	var s sellerIdentity
	var pan, legalName, taxAddr, name, addr *string
	err := tx.QueryRow(ctx,
		`SELECT pan_number, tax_legal_name, tax_address, name, address, is_vat_registered
		   FROM organizations WHERE id = $1 AND is_active = true`,
		orgID,
	).Scan(&pan, &legalName, &taxAddr, &name, &addr, &s.VATRegistered)
	if errors.Is(err, pgx.ErrNoRows) {
		return s, ErrNotFound
	}
	if err != nil {
		return s, fmt.Errorf("loading seller identity: %w", err)
	}

	if pan == nil || strings.TrimSpace(*pan) == "" {
		return s, ErrPANNotConfigured
	}
	s.PAN = *pan

	s.Name = deref(name)
	if legalName != nil && strings.TrimSpace(*legalName) != "" {
		s.Name = *legalName
	}
	s.Address = deref(addr)
	if taxAddr != nil && strings.TrimSpace(*taxAddr) != "" {
		s.Address = *taxAddr
	}
	return s, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// allocateSequence reserves the next number for this org and fiscal year.
//
// A Postgres sequence would be wrong: sequences are non-transactional, so a
// rolled-back insert burns its number and leaves a gap in a run that is
// required to be unbroken. The no-op DO UPDATE is what takes the row lock —
// ON CONFLICT DO NOTHING returns no row and would race.
func allocateSequence(ctx context.Context, tx pgx.Tx, orgID, fiscalYear string) (int, error) {
	var seq int
	err := tx.QueryRow(ctx,
		`INSERT INTO invoice_counters (organization_id, fiscal_year, next_sequence)
		 VALUES ($1, $2, 1)
		 ON CONFLICT (organization_id, fiscal_year)
		 DO UPDATE SET next_sequence = invoice_counters.next_sequence
		 RETURNING next_sequence`,
		orgID, fiscalYear,
	).Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("allocating invoice number: %w", err)
	}

	if _, err := tx.Exec(ctx,
		`UPDATE invoice_counters SET next_sequence = next_sequence + 1
		  WHERE organization_id = $1 AND fiscal_year = $2`,
		orgID, fiscalYear,
	); err != nil {
		return 0, fmt.Errorf("advancing invoice number: %w", err)
	}

	return seq, nil
}

// verifyOrgReferences checks that every client-supplied foreign key on the
// request actually belongs to orgID, and that a supplied transaction_id is
// still eligible to be linked, before anything is written.
//
// FK constraints alone only prove the row exists somewhere — not that it
// belongs to the org making the request. That gap matters more here than on
// most tables: idx_invoices_one_per_transaction is a *global* unique index,
// not scoped per org, so billing against another org's transaction_id would
// silently consume that org's one-bill slot for a payment it never actually
// billed. Run inside the same transaction as the insert (see Issue) so there
// is no time-of-check/time-of-use gap between this check and the write.
func verifyOrgReferences(ctx context.Context, tx pgx.Tx, orgID string, in IssueInput) error {
	if in.TransactionID != "" {
		var ok bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM transactions WHERE id = $1 AND organization_id = $2)`,
			in.TransactionID, orgID,
		).Scan(&ok); err != nil {
			return fmt.Errorf("checking transaction ownership: %w", err)
		}
		if !ok {
			return ErrTransactionNotInOrg
		}

		// owns_transaction ties a transaction's fate to its owning invoice
		// permanently, not just while that invoice is issued: cancelling the
		// invoice reverses the money but does not release the row, because
		// idx_invoices_one_per_transaction only stops a second *issued*
		// invoice from linking it — a cancelled owner doesn't count there,
		// which would otherwise let a new bill quote the same transaction_id,
		// come back with owns_transaction = false, and leave the ledger
		// under-reporting the sale it's supposedly billing. A package-flow
		// transaction (never owned by any invoice) is unaffected and stays
		// relinkable across cancel-and-reissue, as it must.
		var owned bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM invoices WHERE transaction_id = $1 AND owns_transaction = true)`,
			in.TransactionID,
		).Scan(&owned); err != nil {
			return fmt.Errorf("checking transaction ownership by another invoice: %w", err)
		}
		if owned {
			return ErrTransactionAlreadyOwned
		}
	}

	if in.MemberPackageID != "" {
		var ok bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM member_packages WHERE id = $1 AND organization_id = $2)`,
			in.MemberPackageID, orgID,
		).Scan(&ok); err != nil {
			return fmt.Errorf("checking member package ownership: %w", err)
		}
		if !ok {
			return ErrMemberPackageNotInOrg
		}
	}

	if in.CustomerUserID != "" {
		var ok bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM organization_members
			   WHERE user_id = $1 AND organization_id = $2 AND status = 'active')`,
			in.CustomerUserID, orgID,
		).Scan(&ok); err != nil {
			return fmt.Errorf("checking customer membership: %w", err)
		}
		if !ok {
			return ErrCustomerNotInOrg
		}
	}

	return nil
}

// issueParams carries everything the service computed for a new document.
type issueParams struct {
	OrgID         string
	FiscalYear    string
	DocType       string
	CreditNoteFor string
	IssuedDate    string
	IssuedDateBS  string
	Subtotal      float64
	Discount      float64
	TaxableAmount float64
	Total         float64
	AmountInWords string
	IssuedBy      string
	In            IssueInput
}

// insertDocument writes an invoice or credit note and its line items inside tx.
//
// seq and number are allocated by the caller, not here: the income row a
// from-scratch bill writes (see Issue) needs the invoice number as its
// ledger reference before the invoice row exists, and insertDocument writing
// the invoice row is what used to allocate the number — a cycle. Both
// callers allocate first and pass the result in, inside the same
// transaction, so a rollback still returns the number.
func insertDocument(ctx context.Context, tx pgx.Tx, p issueParams, seller sellerIdentity, seq int, number string, ownsTransaction bool) (string, error) {
	var id string
	err := tx.QueryRow(ctx,
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			doc_type, credit_note_for,
			seller_name, seller_pan, seller_address, seller_vat_registered,
			customer_user_id, customer_name, customer_pan, customer_address, customer_phone,
			issued_date, issued_date_bs,
			subtotal, discount, taxable_amount, vat_rate, vat_amount, total, amount_in_words,
			payment_method, transaction_id, member_package_id, issued_by, owns_transaction
		 ) VALUES (
			$1, $2, $3, $4,
			$5, NULLIF($6, '')::uuid,
			$7, $8, $9, $10,
			NULLIF($11, '')::uuid, $12, NULLIF($13, ''), NULLIF($14, ''), NULLIF($15, ''),
			$16::date, $17,
			$18, $19, $20, 0, 0, $21, $22,
			NULLIF($23, ''), NULLIF($24, '')::uuid, NULLIF($25, '')::uuid, $26, $27
		 ) RETURNING id`,
		p.OrgID, p.FiscalYear, seq, number,
		p.DocType, p.CreditNoteFor,
		seller.Name, seller.PAN, seller.Address, seller.VATRegistered,
		p.In.CustomerUserID, p.In.CustomerName, p.In.CustomerPAN, p.In.CustomerAddress, p.In.CustomerPhone,
		p.IssuedDate, p.IssuedDateBS,
		p.Subtotal, p.Discount, p.TaxableAmount, p.Total, p.AmountInWords,
		p.In.PaymentMethod, p.In.TransactionID, p.In.MemberPackageID, p.IssuedBy, ownsTransaction,
	).Scan(&id)
	if err != nil {
		// The partial unique index on transaction_id is what stops the same
		// payment being billed twice. Match on the structured constraint
		// name, not err.Error() — PgError's Error() format ("Severity:
		// Message (SQLSTATE Code)") doesn't guarantee the constraint name
		// appears in the message text; a Postgres wording change or an
		// index rename would silently degrade this 409 into a generic 400.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "idx_invoices_one_per_transaction" {
			return "", ErrAlreadyBilled
		}
		return "", fmt.Errorf("inserting invoice: %w", err)
	}

	for i, it := range p.In.Items {
		if _, err := tx.Exec(ctx,
			`INSERT INTO invoice_items (invoice_id, line_no, description, description_ne, quantity, unit_price, amount)
			 VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7)`,
			id, i+1, it.Description, it.DescriptionNe, it.Quantity, it.UnitPrice,
			round2(it.Quantity*it.UnitPrice),
		); err != nil {
			return "", fmt.Errorf("inserting line item %d: %w", i+1, err)
		}
	}

	return id, nil
}

// Issue writes a new bill and returns it.
func (r *Repository) Issue(ctx context.Context, p issueParams) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	seller, err := r.loadSeller(ctx, tx, p.OrgID)
	if err != nil {
		return nil, err
	}

	if err := verifyOrgReferences(ctx, tx, p.OrgID, p.In); err != nil {
		return nil, err
	}

	seq, err := allocateSequence(ctx, tx, p.OrgID, p.FiscalYear)
	if err != nil {
		return nil, err
	}
	number := fmt.Sprintf("%s/%06d", p.FiscalYear, seq)

	// A bill raised from scratch is the only record that this money came in, so
	// it creates its own income row. A bill pre-filled from a package sale links
	// the row that flow already wrote — creating a second one would double-count
	// the same payment.
	ownsTransaction := false
	if p.In.TransactionID == "" && p.Total > 0 {
		var transactionID string
		err := tx.QueryRow(ctx,
			`INSERT INTO transactions
			        (organization_id, category, description, transaction_date,
			         transaction_type, amount, payment_type, reference, entry_by)
			 VALUES ($1, 'Sales', $2, CURRENT_DATE, 'income', $3, $4, $5, $6)
			 RETURNING id`,
			p.OrgID,
			"Invoice "+number+" — "+p.In.CustomerName,
			p.Total,
			defaultTo(p.In.PaymentMethod, "cash"),
			number,
			p.IssuedBy,
		).Scan(&transactionID)
		if err != nil {
			return nil, fmt.Errorf("recording invoice income: %w", err)
		}
		p.In.TransactionID = transactionID
		ownsTransaction = true
	}

	id, err := insertDocument(ctx, tx, p, seller, seq, number, ownsTransaction)
	if err != nil {
		return nil, err
	}

	inv, err := getInTx(ctx, tx, p.OrgID, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing invoice: %w", err)
	}
	return inv, nil
}

const invoiceColumns = `
	id, organization_id, fiscal_year, sequence, invoice_number,
	doc_type, COALESCE(credit_note_for::text, ''),
	seller_name, seller_pan, COALESCE(seller_address, ''), seller_vat_registered,
	COALESCE(customer_user_id::text, ''), customer_name,
	COALESCE(customer_pan, ''), COALESCE(customer_address, ''), COALESCE(customer_phone, ''),
	issued_date::text, issued_date_bs,
	subtotal, discount, taxable_amount, vat_rate, vat_amount, total, amount_in_words,
	COALESCE(payment_method, ''),
	status, cancelled_at, COALESCE(cancellation_reason, ''),
	COALESCE(transaction_id::text, ''), COALESCE(member_package_id::text, ''),
	issued_by, print_count, created_at`

func scanInvoice(row pgx.Row) (*Invoice, error) {
	var v Invoice
	err := row.Scan(
		&v.ID, &v.OrgID, &v.FiscalYear, &v.Sequence, &v.InvoiceNumber,
		&v.DocType, &v.CreditNoteFor,
		&v.SellerName, &v.SellerPAN, &v.SellerAddress, &v.SellerVATRegistered,
		&v.CustomerUserID, &v.CustomerName,
		&v.CustomerPAN, &v.CustomerAddress, &v.CustomerPhone,
		&v.IssuedDate, &v.IssuedDateBS,
		&v.Subtotal, &v.Discount, &v.TaxableAmount, &v.VATRate, &v.VATAmount,
		&v.Total, &v.AmountInWords,
		&v.PaymentMethod,
		&v.Status, &v.CancelledAt, &v.CancellationReason,
		&v.TransactionID, &v.MemberPackageID,
		&v.IssuedBy, &v.PrintCount, &v.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning invoice: %w", err)
	}
	return &v, nil
}

func getInTx(ctx context.Context, tx pgx.Tx, orgID, id string) (*Invoice, error) {
	inv, err := scanInvoice(tx.QueryRow(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE id = $1 AND organization_id = $2`,
		id, orgID))
	if err != nil {
		return nil, err
	}
	items, err := loadItems(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

func loadItems(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, invoiceID string) ([]Item, error) {
	rows, err := q.Query(ctx,
		`SELECT line_no, description, COALESCE(description_ne, ''), quantity, unit_price, amount
		   FROM invoice_items WHERE invoice_id = $1 ORDER BY line_no`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("loading line items: %w", err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.LineNo, &it.Description, &it.DescriptionNe,
			&it.Quantity, &it.UnitPrice, &it.Amount); err != nil {
			return nil, fmt.Errorf("scanning line item: %w", err)
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// Get reads one invoice. Scoping by org in the WHERE clause means another
// gym's invoice is indistinguishable from a missing one.
func (r *Repository) Get(ctx context.Context, orgID, id string) (*Invoice, error) {
	inv, err := scanInvoice(r.db.QueryRow(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE id = $1 AND organization_id = $2`,
		id, orgID))
	if err != nil {
		return nil, err
	}
	items, err := loadItems(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

// Cancel marks an invoice cancelled. The number stays consumed; that is the
// point of cancelling rather than deleting.
//
// The ledger must always equal the sum of non-cancelled documents. A bill
// that created its own income row (owns_transaction) must reverse exactly
// that row here, in the same transaction, or the ledger keeps money for a
// document that no longer counts. A bill that only linked a package sale's
// row leaves it alone — that income belongs to the package flow, not to this
// document, so cancelling the bill must not touch it. A credit note DID
// write a refund row at issue, so cancelling one must reverse it too:
// otherwise the refund survives the cancellation, and creditedSoFar (which
// only counts status='issued' credit notes) silently frees up capacity to
// credit the parent all over again.
func (r *Repository) Cancel(ctx context.Context, orgID, id, cancelledBy, reason string) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var docType, invoiceNumber, paymentMethod string
	var total float64
	var ownsTransaction bool
	err = tx.QueryRow(ctx,
		`UPDATE invoices
		    SET status = 'cancelled', cancelled_at = NOW(),
		        cancelled_by = $3, cancellation_reason = $4
		  WHERE id = $1 AND organization_id = $2 AND status = 'issued'
		  RETURNING doc_type, invoice_number, total, COALESCE(payment_method, ''), owns_transaction`,
		id, orgID, cancelledBy, reason,
	).Scan(&docType, &invoiceNumber, &total, &paymentMethod, &ownsTransaction)
	if errors.Is(err, pgx.ErrNoRows) {
		// Either it does not exist in this org, or it is already cancelled.
		existing, getErr := getInTx(ctx, tx, orgID, id)
		if getErr != nil {
			return nil, getErr
		}
		if existing.Status == "cancelled" {
			return nil, ErrAlreadyCancelled
		}
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("cancelling invoice: %w", err)
	}

	if docType == "credit_note" {
		if _, err := tx.Exec(ctx,
			`INSERT INTO transactions (organization_id, category, description, transaction_date,
			                           transaction_type, amount, payment_type, reference, entry_by)
			 VALUES ($1, 'refund_reversal', $2, CURRENT_DATE, 'income', $3, $4, $5, $6)`,
			orgID,
			"Cancelled credit note "+invoiceNumber,
			total,
			defaultTo(paymentMethod, "cash"),
			invoiceNumber,
			cancelledBy,
		); err != nil {
			return nil, fmt.Errorf("writing refund reversal ledger row: %w", err)
		}
	} else if docType == "invoice" && ownsTransaction {
		// This bill created the income row it links (see Issue), so
		// cancelling it must reverse exactly that money, not a package sale's
		// income it never owned.
		if _, err := tx.Exec(ctx,
			`INSERT INTO transactions (organization_id, category, description, transaction_date,
			                           transaction_type, amount, payment_type, reference, entry_by)
			 VALUES ($1, 'Sales reversal', $2, CURRENT_DATE, 'expense', $3, $4, $5, $6)`,
			orgID,
			"Cancelled invoice "+invoiceNumber,
			total,
			defaultTo(paymentMethod, "cash"),
			invoiceNumber,
			cancelledBy,
		); err != nil {
			return nil, fmt.Errorf("writing invoice reversal ledger row: %w", err)
		}
	}

	inv, err := getInTx(ctx, tx, orgID, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing cancellation: %w", err)
	}
	return inv, nil
}

// creditedSoFar totals the credit notes already raised against an invoice.
func creditedSoFar(ctx context.Context, tx pgx.Tx, parentID string) (float64, error) {
	var sum float64
	err := tx.QueryRow(ctx,
		`SELECT COALESCE(SUM(total), 0) FROM invoices
		  WHERE credit_note_for = $1 AND status = 'issued'`, parentID).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("totalling existing credit notes: %w", err)
	}
	return sum, nil
}

// lockParent takes a row lock on the parent invoice before any capacity
// check runs. Without it, two concurrent credit notes both read the same
// creditedSoFar under READ COMMITTED, both pass the cap check, and both
// insert — over-crediting the parent. The lock forces the second transaction
// to block until the first commits, so it re-reads a current sum. Scoped to
// orgID so a cross-tenant parent still 404s rather than blocking on a row
// the caller has no right to see.
func lockParent(ctx context.Context, tx pgx.Tx, orgID, parentID string) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT true FROM invoices WHERE id = $1 AND organization_id = $2 FOR UPDATE`,
		parentID, orgID,
	).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("locking parent invoice: %w", err)
	}
	return nil
}

// CreditNote raises a credit note against parentID and writes the reversing
// ledger row, both inside one transaction.
func (r *Repository) CreditNote(ctx context.Context, p issueParams, parentID string) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock first, before reading anything: the cap check below is only
	// correct if no other transaction can concurrently insert a credit note
	// against this same parent between our read of creditedSoFar and our
	// insert.
	if err := lockParent(ctx, tx, p.OrgID, parentID); err != nil {
		return nil, err
	}

	parent, err := getInTx(ctx, tx, p.OrgID, parentID)
	if err != nil {
		return nil, err
	}
	if parent.Status == "cancelled" {
		return nil, ErrInvoiceCancelled
	}
	if parent.DocType != "invoice" {
		return nil, ErrInvalidParent
	}

	already, err := creditedSoFar(ctx, tx, parentID)
	if err != nil {
		return nil, err
	}
	if round2(already+p.Total) > parent.Total {
		return nil, ErrCreditTooLarge
	}

	seller, err := r.loadSeller(ctx, tx, p.OrgID)
	if err != nil {
		return nil, err
	}

	seq, err := allocateSequence(ctx, tx, p.OrgID, p.FiscalYear)
	if err != nil {
		return nil, err
	}
	number := fmt.Sprintf("%s/%06d", p.FiscalYear, seq)

	p.CreditNoteFor = parentID
	// A credit note never owns a transaction: it writes its own refund row
	// below regardless of who wrote the original income, so it has nothing
	// of its own to reverse on cancel.
	id, err := insertDocument(ctx, tx, p, seller, seq, number, false)
	if err != nil {
		return nil, err
	}

	// A refund is a real movement of money, so it reverses in the ledger
	// regardless of who wrote the original income row.
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (organization_id, category, description, transaction_date,
		                           transaction_type, amount, payment_type, reference, entry_by)
		 VALUES ($1, 'refund', $2, CURRENT_DATE, 'expense', $3, $4, $5, $6)`,
		p.OrgID,
		"Credit note "+number+" against "+parent.InvoiceNumber,
		p.Total,
		defaultTo(parent.PaymentMethod, "cash"),
		number,
		p.IssuedBy,
	); err != nil {
		return nil, fmt.Errorf("writing refund ledger row: %w", err)
	}

	note, err := getInTx(ctx, tx, p.OrgID, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing credit note: %w", err)
	}
	return note, nil
}

func defaultTo(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// RecordPrint logs a print and returns the label the document should carry.
//
// The label is derived from the UPDATE's own RETURNING value, not from a
// prior SELECT. A plain read-then-compare-then-write here would let two
// concurrent prints both observe print_count = 0 and both log themselves as
// "original" — the count would still land on 2 (Postgres serializes the
// UPDATE itself), but the audit trail would then claim two originals of one
// document. Deriving the label from the row lock the UPDATE already takes
// closes that race: the second transaction blocks until the first commits,
// then computes its label from the true post-increment count.
func (r *Repository) RecordPrint(ctx context.Context, orgID, id, printedBy string) (*Invoice, string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var newCount int
	err = tx.QueryRow(ctx,
		`UPDATE invoices SET print_count = print_count + 1
		  WHERE id = $1 AND organization_id = $2
		  RETURNING print_count`,
		id, orgID).Scan(&newCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", fmt.Errorf("incrementing print count: %w", err)
	}

	label := "original"
	if newCount > 1 {
		label = "copy"
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO invoice_prints (invoice_id, printed_by, copy_label) VALUES ($1, $2, $3)`,
		id, printedBy, label); err != nil {
		return nil, "", fmt.Errorf("logging print: %w", err)
	}

	updated, err := getInTx(ctx, tx, orgID, id)
	if err != nil {
		return nil, "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, "", fmt.Errorf("committing print: %w", err)
	}
	return updated, label, nil
}

// PeekNextSequence reports the number the next bill will take without
// reserving it. A read, never a write — previewing must not consume.
func (r *Repository) PeekNextSequence(ctx context.Context, orgID, fiscalYear string) (int, error) {
	var seq int
	err := r.db.QueryRow(ctx,
		`SELECT next_sequence FROM invoice_counters
		  WHERE organization_id = $1 AND fiscal_year = $2`, orgID, fiscalYear).Scan(&seq)
	if errors.Is(err, pgx.ErrNoRows) {
		return 1, nil
	}
	if err != nil {
		return 0, fmt.Errorf("reading next invoice number: %w", err)
	}
	return seq, nil
}

// List returns invoices for an org, newest first.
func (r *Repository) List(ctx context.Context, f ListFilter) ([]Invoice, int, error) {
	conds := []string{"organization_id = $1"}
	args := []any{f.OrgID}

	add := func(clause string, val any) {
		args = append(args, val)
		conds = append(conds, fmt.Sprintf(clause, len(args)))
	}
	if f.Status != "" {
		add("status = $%d", f.Status)
	}
	if f.FiscalYear != "" {
		add("fiscal_year = $%d", f.FiscalYear)
	}
	if f.CustomerID != "" {
		add("customer_user_id = $%d::uuid", f.CustomerID)
	}
	if f.From != "" {
		add("issued_date >= $%d::date", f.From)
	}
	if f.To != "" {
		add("issued_date <= $%d::date", f.To)
	}
	if f.Query != "" {
		args = append(args, "%"+f.Query+"%")
		conds = append(conds, fmt.Sprintf(
			"(invoice_number ILIKE $%d OR customer_name ILIKE $%d)", len(args), len(args)))
	}
	where := strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM invoices WHERE `+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting invoices: %w", err)
	}

	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE `+where+
			fmt.Sprintf(" ORDER BY sequence DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args)),
		args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing invoices: %w", err)
	}
	defer rows.Close()

	out := []Invoice{}
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *inv)
	}
	return out, total, rows.Err()
}
