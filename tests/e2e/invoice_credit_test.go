package e2e

import (
	"context"
	"net/http"
	"sync"
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

// A credit note that cannot say which bill it reverses is not a usable tax
// document — credit_note_for is a UUID the customer never sees, so the
// printed reference has to be the parent's own invoice_number, snapshotted at
// credit time. A plain invoice has no parent, so it must carry none.
func TestInvoiceCredit_CarriesParentInvoiceNumber(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000508", "Admin")
	orgID := createTestOrg(t, admin, "Reference Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var parent struct {
		ID                  string `json:"id"`
		InvoiceNumber       string `json:"invoice_number"`
		CreditNoteForNumber string `json:"credit_note_for_number"`
	}
	parseJSON(t, resp, &parent)
	if parent.CreditNoteForNumber != "" {
		t.Errorf("plain invoice credit_note_for_number = %q, want empty", parent.CreditNoteForNumber)
	}

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 500), token)
	assertStatus(t, resp, http.StatusCreated)

	var note struct {
		ID                  string `json:"id"`
		CreditNoteFor       string `json:"credit_note_for"`
		CreditNoteForNumber string `json:"credit_note_for_number"`
	}
	parseJSON(t, resp, &note)
	if note.CreditNoteFor != parent.ID {
		t.Errorf("credit_note_for = %q, want the parent id %q", note.CreditNoteFor, parent.ID)
	}
	if note.CreditNoteForNumber != parent.InvoiceNumber {
		t.Errorf("credit_note_for_number = %q, want the parent's invoice_number %q",
			note.CreditNoteForNumber, parent.InvoiceNumber)
	}

	// Reading the credit note back later (e.g. reprinting it) must still
	// carry the reference — this isn't only present on the create response.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/invoices/"+note.ID, nil, token)
	assertStatus(t, resp, http.StatusOK)
	var reread struct {
		CreditNoteForNumber string `json:"credit_note_for_number"`
	}
	parseJSON(t, resp, &reread)
	if reread.CreditNoteForNumber != parent.InvoiceNumber {
		t.Errorf("re-read credit_note_for_number = %q, want %q", reread.CreditNoteForNumber, parent.InvoiceNumber)
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

// The core guarantee behind the cap: two credit notes racing against the
// same parent must not both pass the check. Parent is 1000, so two
// simultaneous credits of 600 cannot both fit — modelled on
// TestInvoice_ConcurrentIssuesAreGapless, which proves the analogous
// guarantee for the number sequence.
func TestInvoiceCredit_ConcurrentCreditsDoNotOvercredit(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000506", "Admin")
	orgID := createTestOrg(t, admin, "Concurrent Credit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	const n = 8
	statuses := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp := doRequest(t, http.MethodPost,
				"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note",
				creditBody(1, 600), token)
			defer resp.Body.Close()
			statuses[i] = resp.StatusCode
		}(i)
	}
	wg.Wait()

	successes := 0
	for _, code := range statuses {
		if code == http.StatusCreated {
			successes++
		}
	}
	if successes != 1 {
		t.Errorf("successful credit notes = %d, want exactly 1 (1000 cap, 600 each)", successes)
	}

	var total float64
	err := testPool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total), 0) FROM invoices
		  WHERE credit_note_for = $1 AND status = 'issued'`, parent.ID).Scan(&total)
	if err != nil {
		t.Fatalf("summing credit notes: %v", err)
	}
	if total > 1000 {
		t.Errorf("total credited = %.2f, want <= 1000", total)
	}
}

// Cancelling a credit note must reverse its refund — otherwise the refund
// stays on the books while creditedSoFar (which only counts status='issued'
// notes) forgets it ever happened, silently freeing capacity to credit the
// same parent again.
func TestInvoiceCredit_CancellingReversesRefundAndFreesCapacity(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000507", "Admin")
	orgID := createTestOrg(t, admin, "Cancel Credit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var note struct{ ID string }
	parseJSON(t, resp, &note)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+note.ID+"/cancel",
		map[string]any{"reason": "issued by mistake"}, token)
	assertStatus(t, resp, http.StatusOK)

	// Exactly one reversal row, for the credited amount.
	var reversals int
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'income'
		    AND category = 'refund_reversal' AND amount = 1000`, orgID).Scan(&reversals)
	if err != nil {
		t.Fatalf("counting reversals: %v", err)
	}
	if reversals != 1 {
		t.Errorf("refund_reversal rows = %d, want exactly 1", reversals)
	}

	// The net ledger effect of issuing a credit note and then cancelling it
	// is zero: the refund and its reversal cancel out. Scoped to the
	// refund/refund_reversal categories, not the whole org ledger — the
	// parent invoice's own income row (task 8b) is still on the books at
	// this point, since only the credit note was cancelled, not the parent.
	var net float64
	err = testPool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE -amount END), 0)
		   FROM transactions WHERE organization_id = $1
		     AND category IN ('refund', 'refund_reversal')`, orgID).Scan(&net)
	if err != nil {
		t.Fatalf("summing net ledger effect: %v", err)
	}
	if net != 0 {
		t.Errorf("net ledger effect of issue-then-cancel = %.2f, want 0", net)
	}

	// The cancelled credit note no longer counts against the parent's
	// capacity, so the full amount can be credited again.
	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)

	// Two 'refund' expense rows now exist (the cancelled one and the live
	// one) but only one reversal — so exactly one refund's worth remains net
	// on the books, not two.
	var refunds, reversalsAfter int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'expense' AND category = 'refund'`,
		orgID).Scan(&refunds); err != nil {
		t.Fatalf("counting refunds: %v", err)
	}
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'income' AND category = 'refund_reversal'`,
		orgID).Scan(&reversalsAfter); err != nil {
		t.Fatalf("counting reversals: %v", err)
	}
	if refunds != 2 {
		t.Errorf("refund expense rows = %d, want 2 (one cancelled, one live)", refunds)
	}
	if reversalsAfter != 1 {
		t.Errorf("refund_reversal rows = %d, want 1 (only the cancelled note was reversed)", reversalsAfter)
	}

	var netRefund float64
	err = testPool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(CASE WHEN category = 'refund' THEN amount
		                          WHEN category = 'refund_reversal' THEN -amount
		                          ELSE 0 END), 0)
		   FROM transactions WHERE organization_id = $1`, orgID).Scan(&netRefund)
	if err != nil {
		t.Fatalf("summing net refund: %v", err)
	}
	if netRefund != 1000 {
		t.Errorf("net refund on the books = %.2f, want 1000 — exactly one live credit note's worth, not two", netRefund)
	}
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
