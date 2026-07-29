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
// request actually belongs to orgID before anything is written.
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
func insertDocument(ctx context.Context, tx pgx.Tx, p issueParams, seller sellerIdentity) (string, string, int, error) {
	seq, err := allocateSequence(ctx, tx, p.OrgID, p.FiscalYear)
	if err != nil {
		return "", "", 0, err
	}
	number := fmt.Sprintf("%s/%06d", p.FiscalYear, seq)

	var id string
	err = tx.QueryRow(ctx,
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			doc_type, credit_note_for,
			seller_name, seller_pan, seller_address, seller_vat_registered,
			customer_user_id, customer_name, customer_pan, customer_address, customer_phone,
			issued_date, issued_date_bs,
			subtotal, discount, taxable_amount, vat_rate, vat_amount, total, amount_in_words,
			payment_method, transaction_id, member_package_id, issued_by
		 ) VALUES (
			$1, $2, $3, $4,
			$5, NULLIF($6, '')::uuid,
			$7, $8, $9, $10,
			NULLIF($11, '')::uuid, $12, NULLIF($13, ''), NULLIF($14, ''), NULLIF($15, ''),
			$16::date, $17,
			$18, $19, $20, 0, 0, $21, $22,
			NULLIF($23, ''), NULLIF($24, '')::uuid, NULLIF($25, '')::uuid, $26
		 ) RETURNING id`,
		p.OrgID, p.FiscalYear, seq, number,
		p.DocType, p.CreditNoteFor,
		seller.Name, seller.PAN, seller.Address, seller.VATRegistered,
		p.In.CustomerUserID, p.In.CustomerName, p.In.CustomerPAN, p.In.CustomerAddress, p.In.CustomerPhone,
		p.IssuedDate, p.IssuedDateBS,
		p.Subtotal, p.Discount, p.TaxableAmount, p.Total, p.AmountInWords,
		p.In.PaymentMethod, p.In.TransactionID, p.In.MemberPackageID, p.IssuedBy,
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
			return "", "", 0, ErrAlreadyBilled
		}
		return "", "", 0, fmt.Errorf("inserting invoice: %w", err)
	}

	for i, it := range p.In.Items {
		if _, err := tx.Exec(ctx,
			`INSERT INTO invoice_items (invoice_id, line_no, description, description_ne, quantity, unit_price, amount)
			 VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7)`,
			id, i+1, it.Description, it.DescriptionNe, it.Quantity, it.UnitPrice,
			round2(it.Quantity*it.UnitPrice),
		); err != nil {
			return "", "", 0, fmt.Errorf("inserting line item %d: %w", i+1, err)
		}
	}

	return id, number, seq, nil
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

	id, _, _, err := insertDocument(ctx, tx, p, seller)
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
