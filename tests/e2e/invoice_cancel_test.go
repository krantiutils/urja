package e2e

import (
	"context"
	"net/http"
	"testing"
)

func TestInvoiceCancel_KeepsNumberConsumed(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000401", "Admin")
	orgID := createTestOrg(t, admin, "Cancel Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var first struct {
		ID       string `json:"id"`
		Sequence int    `json:"sequence"`
	}
	parseJSON(t, resp, &first)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+first.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	var cancelled struct {
		Status string `json:"status"`
		Reason string `json:"cancellation_reason"`
	}
	parseJSON(t, resp, &cancelled)
	if cancelled.Status != "cancelled" {
		t.Errorf("status = %q, want cancelled", cancelled.Status)
	}
	if cancelled.Reason != "wrong customer" {
		t.Errorf("reason = %q, want it recorded", cancelled.Reason)
	}

	// The cancelled number is spent: the next bill takes the one after it.
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Shyam", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var second struct {
		Sequence int `json:"sequence"`
	}
	parseJSON(t, resp, &second)
	if second.Sequence != first.Sequence+1 {
		t.Errorf("next sequence = %d, want %d — a cancelled number must not be reused",
			second.Sequence, first.Sequence+1)
	}
}

func TestInvoiceCancel_RequiresReason(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000402", "Admin")
	orgID := createTestOrg(t, admin, "Reason Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel",
		map[string]any{"reason": "   "}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestInvoiceCancel_IsNotRepeatable(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000403", "Admin")
	orgID := createTestOrg(t, admin, "Twice Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	body := map[string]any{"reason": "duplicate"}
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel", body, token)
	assertStatus(t, resp, http.StatusOK)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel", body, token)
	assertStatus(t, resp, http.StatusConflict)
}

// Cancelling through another org's path must 404, not 403 — a 403 would
// confirm the invoice exists in a gym the caller has no business seeing —
// and it must not actually cancel the invoice.
func TestInvoiceCancel_CrossTenantDoesNotCancel(t *testing.T) {
	cleanupTables(t)
	adminA := createTestUser(t, "9800000404", "Admin A")
	adminB := createTestUser(t, "9800000405", "Admin B")
	orgA := createTestOrg(t, adminA, "Tenant A Cancel")
	orgB := createTestOrg(t, adminB, "Tenant B Cancel")
	setPAN(t, orgB, "601234567")
	tokenA := generateTestToken(adminA, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgB+"/invoices",
		issueBody("Ram", 1, 1000), generateTestToken(adminB, "member"))
	assertStatus(t, resp, http.StatusCreated)
	var invB struct{ ID string }
	parseJSON(t, resp, &invB)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgA+"/invoices/"+invB.ID+"/cancel",
		map[string]any{"reason": "not yours"}, tokenA)
	assertStatus(t, resp, http.StatusNotFound)

	var status string
	err := testPool.QueryRow(context.Background(),
		`SELECT status FROM invoices WHERE id = $1`, invB.ID).Scan(&status)
	if err != nil {
		t.Fatalf("reading invoice status: %v", err)
	}
	if status != "issued" {
		t.Errorf("status = %q, want issued — org A must not have cancelled org B's invoice", status)
	}
}

// Regression guard on the path Cancel's credit-note fix does NOT touch. This
// used to assert that cancelling a plain invoice wrote nothing at all to the
// ledger — true only because, before task 8b, a from-scratch bill recorded no
// income in the first place. Task 8b's own tests (invoice_ledger_test.go)
// cover that gap; what's left worth guarding here, alongside the rest of this
// file's cancel-path tests, is that the net effect still nets to zero rather
// than leaving the reversed income sitting on the books uncancelled.
func TestInvoiceCancel_PlainInvoiceNetsToZeroOnLedger(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000406", "Admin")
	orgID := createTestOrg(t, admin, "No Ledger On Cancel Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	var net float64
	err := testPool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE -amount END), 0)
		   FROM transactions WHERE organization_id = $1`, orgID).Scan(&net)
	if err != nil {
		t.Fatalf("summing net ledger effect: %v", err)
	}
	if net != 0 {
		t.Errorf("net ledger effect after cancelling a plain invoice = %.2f, want 0", net)
	}
}
