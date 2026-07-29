package e2e

import (
	"context"
	"net/http"
	"testing"
)

func creditBody(qty, price float64) map[string]any {
	return map[string]any{
		"reason": "member left mid-month",
		"items": []map[string]any{
			{"description": "Refund: Monthly Boxing", "quantity": qty, "unit_price": price},
		},
	}
}

func TestInvoiceCredit_ReversesAndLinks(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000501", "Admin")
	orgID := createTestOrg(t, admin, "Credit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 3000), token)
	assertStatus(t, resp, http.StatusCreated)
	var parent struct {
		ID       string `json:"id"`
		Sequence int    `json:"sequence"`
	}
	parseJSON(t, resp, &parent)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note",
		creditBody(1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)

	var note struct {
		DocType       string  `json:"doc_type"`
		CreditNoteFor string  `json:"credit_note_for"`
		Total         float64 `json:"total"`
		Sequence      int     `json:"sequence"`
	}
	parseJSON(t, resp, &note)

	if note.DocType != "credit_note" {
		t.Errorf("doc_type = %q, want credit_note", note.DocType)
	}
	if note.CreditNoteFor != parent.ID {
		t.Errorf("credit_note_for = %q, want the parent invoice", note.CreditNoteFor)
	}
	if note.Total != 1000 {
		t.Errorf("total = %.2f, want 1000", note.Total)
	}
	// Credit notes share the invoice sequence — one unbroken run per year.
	if note.Sequence != parent.Sequence+1 {
		t.Errorf("sequence = %d, want %d", note.Sequence, parent.Sequence+1)
	}

	// A reversing ledger row must exist for the credited amount.
	var refunds int
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'expense'
		    AND category = 'refund' AND amount = 1000`, orgID).Scan(&refunds)
	if err != nil {
		t.Fatalf("counting refunds: %v", err)
	}
	if refunds != 1 {
		t.Errorf("refund ledger rows = %d, want exactly 1", refunds)
	}
}

func TestInvoiceCredit_CannotExceedBalance(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000502", "Admin")
	orgID := createTestOrg(t, admin, "Overcredit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	// Credit 800 of 1000 — fine.
	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 800), token)
	assertStatus(t, resp, http.StatusCreated)

	// Another 400 would exceed the remaining 200.
	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 400), token)
	assertStatus(t, resp, http.StatusBadRequest)

	var body struct{ Code string }
	parseJSON(t, resp, &body)
	if body.Code != "credit_exceeds_balance" {
		t.Errorf("code = %q, want credit_exceeds_balance", body.Code)
	}
}

func TestInvoiceCredit_RefusesCancelledParent(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000503", "Admin")
	orgID := createTestOrg(t, admin, "Cancelled Parent Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/cancel",
		map[string]any{"reason": "error"}, token)
	assertStatus(t, resp, http.StatusOK)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 500), token)
	assertStatus(t, resp, http.StatusConflict)
}

// A credit note against another org's invoice must 404 like every other
// cross-tenant read, and — because this endpoint also writes a document and a
// ledger row — it must not create either as a side effect.
func TestInvoiceCredit_CrossTenantDoesNotCreditOrWriteLedger(t *testing.T) {
	cleanupTables(t)
	adminA := createTestUser(t, "9800000504", "Admin A")
	adminB := createTestUser(t, "9800000505", "Admin B")
	orgA := createTestOrg(t, adminA, "Tenant A Credit")
	orgB := createTestOrg(t, adminB, "Tenant B Credit")
	setPAN(t, orgA, "601234567")
	setPAN(t, orgB, "601234568")
	tokenA := generateTestToken(adminA, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgB+"/invoices",
		issueBody("Ram", 1, 1000), generateTestToken(adminB, "member"))
	assertStatus(t, resp, http.StatusCreated)
	var invB struct{ ID string }
	parseJSON(t, resp, &invB)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgA+"/invoices/"+invB.ID+"/credit-note", creditBody(1, 500), tokenA)
	assertStatus(t, resp, http.StatusNotFound)

	var docs int
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM invoices WHERE credit_note_for = $1`, invB.ID).Scan(&docs)
	if err != nil {
		t.Fatalf("counting credit notes: %v", err)
	}
	if docs != 0 {
		t.Errorf("credit notes against org B's invoice = %d, want 0 — org A must not have credited it", docs)
	}

	var refunds int
	err = testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE transaction_type = 'expense' AND category = 'refund'`).Scan(&refunds)
	if err != nil {
		t.Fatalf("counting refunds: %v", err)
	}
	if refunds != 0 {
		t.Errorf("refund ledger rows = %d, want 0 — no credit note was ever created", refunds)
	}
}
