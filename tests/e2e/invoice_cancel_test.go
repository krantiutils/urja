package e2e

import (
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
