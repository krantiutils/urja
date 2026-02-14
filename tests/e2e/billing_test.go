package e2e

import (
	"net/http"
	"testing"
)

func TestBilling_ListPlans(t *testing.T) {
	// Plans are seeded via migration 000014
	resp := doRequest(t, http.MethodGet, "/api/v1/billing/plans", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestBilling_Subscribe_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/billing/subscribe",
		map[string]string{
			"organization_id": "some-id",
			"plan_id":         "some-plan",
			"billing_period":  "monthly",
			"khalti_pidx":     "some-pidx",
		}, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestBilling_Subscribe_MissingFields(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9801500001", "User")
	token := generateTestToken(userID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/billing/subscribe",
		map[string]string{}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestBilling_Subscribe_MissingPlanID(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9801500001", "User")
	token := generateTestToken(userID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/billing/subscribe",
		map[string]string{
			"organization_id": "some-org-id",
			"billing_period":  "monthly",
			"khalti_pidx":     "test-pidx",
		}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}
