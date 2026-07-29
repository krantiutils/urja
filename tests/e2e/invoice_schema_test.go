package e2e

import (
	"context"
	"strings"
	"testing"
)

// seedInvoice inserts a minimal issued invoice directly and returns its id.
func seedInvoice(t *testing.T, orgID, issuedBy string, seq int) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			seller_name, seller_pan, customer_name,
			issued_date, issued_date_bs,
			subtotal, taxable_amount, total, amount_in_words, issued_by
		) VALUES ($1, '2082-83', $2::int, '2082-83/' || lpad($2::text, 6, '0'),
			'Test Gym', '123456789', 'Ram Bahadur',
			CURRENT_DATE, '2082-04-14',
			1000, 1000, 1000, 'One thousand rupees only', $3)
		RETURNING id`,
		orgID, seq, issuedBy,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedInvoice: %v", err)
	}
	return id
}

func TestInvoiceSchema_CannotUpdateAmount(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000101", "Admin")
	orgID := createTestOrg(t, admin, "Immutable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET total = 5000 WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject an amount change, got nil")
	}
	if !strings.Contains(err.Error(), "only cancellation may change") {
		t.Errorf("unexpected error: %v", err)
	}
}

// credit_note_for_number is a snapshot like every other field on this table:
// once written, it must be as immutable as the amounts it sits next to.
func TestInvoiceSchema_CannotUpdateCreditNoteForNumber(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000107", "Admin")
	orgID := createTestOrg(t, admin, "Reference Immutable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET credit_note_for_number = '2082-83/999999' WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject a credit_note_for_number change, got nil")
	}
	if !strings.Contains(err.Error(), "only cancellation may change") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_CannotChangeID(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000106", "Admin")
	orgID := createTestOrg(t, admin, "PK Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET id = gen_random_uuid() WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject a primary key change, got nil")
	}
	if !strings.Contains(err.Error(), "the primary key cannot be changed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_CannotDelete(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000102", "Admin")
	orgID := createTestOrg(t, admin, "Undeletable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(), `DELETE FROM invoices WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject a delete, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be deleted") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_CancellationIsAllowedAndOneWay(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000103", "Admin")
	orgID := createTestOrg(t, admin, "Cancellable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	// Cancelling is the one permitted mutation.
	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET status = 'cancelled', cancelled_at = NOW(),
		        cancelled_by = $2, cancellation_reason = 'wrong customer'
		 WHERE id = $1`, id, admin)
	if err != nil {
		t.Fatalf("cancelling should be permitted: %v", err)
	}

	// Un-cancelling is not.
	_, err = testPool.Exec(context.Background(),
		`UPDATE invoices SET status = 'issued' WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected un-cancelling to be rejected, got nil")
	}
	if !strings.Contains(err.Error(), "cancellation cannot be undone") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_LineItemsAreInsertOnly(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000104", "Admin")
	orgID := createTestOrg(t, admin, "Lineitem Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`INSERT INTO invoice_items (invoice_id, line_no, description, quantity, unit_price, amount)
		 VALUES ($1, 1, 'Monthly Boxing', 1, 1000, 1000)`, id)
	if err != nil {
		t.Fatalf("inserting a line item should work: %v", err)
	}

	_, err = testPool.Exec(context.Background(),
		`UPDATE invoice_items SET amount = 9999 WHERE invoice_id = $1`, id)
	if err == nil {
		t.Fatal("expected line item update to be rejected, got nil")
	}
}

func TestInvoiceSchema_SequenceIsUniquePerOrgAndYear(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000105", "Admin")
	orgA := createTestOrg(t, admin, "Seq Gym A")
	orgB := createTestOrg(t, admin, "Seq Gym B")

	// seedInvoice(orgA, ..., 1) takes sequence=1, invoice_number='2082-83/000001'.
	seedInvoice(t, orgA, admin, 1)

	// The same sequence in a different org is fine.
	seedInvoice(t, orgB, admin, 1)

	// Same sequence, same org and year, but a DISTINCT invoice_number: only
	// the (organization_id, fiscal_year, sequence) constraint can reject this,
	// so this specifically pins sequence uniqueness rather than riding along
	// on the invoice_number constraint.
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			seller_name, seller_pan, customer_name, issued_date, issued_date_bs,
			subtotal, taxable_amount, total, amount_in_words, issued_by
		) VALUES ($1, '2082-83', 1, '2082-83/000999',
			'Test Gym', '123456789', 'Someone', CURRENT_DATE, '2082-04-14',
			1000, 1000, 1000, 'One thousand rupees only', $2)`,
		orgA, admin)
	if err == nil {
		t.Fatal("expected duplicate sequence to be rejected, got nil")
	}

	// Mirror case: same invoice_number, DISTINCT sequence. Only the
	// (organization_id, invoice_number) constraint can reject this, so this
	// pins that constraint independently of the one above.
	_, err = testPool.Exec(context.Background(),
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			seller_name, seller_pan, customer_name, issued_date, issued_date_bs,
			subtotal, taxable_amount, total, amount_in_words, issued_by
		) VALUES ($1, '2082-83', 2, '2082-83/000001',
			'Test Gym', '123456789', 'Someone', CURRENT_DATE, '2082-04-14',
			1000, 1000, 1000, 'One thousand rupees only', $2)`,
		orgA, admin)
	if err == nil {
		t.Fatal("expected duplicate invoice_number to be rejected, got nil")
	}
}
