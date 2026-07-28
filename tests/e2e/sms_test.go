package e2e

import (
	"net/http"
	"testing"
)


// Audit finding #5, still live: the purchase took the rate and the total from
// the request body and stored both, so anybody who could reach the endpoint
// could award themselves credits for whatever price they cared to name.
func TestSMS_BuyCredits_ServerOwnsThePrice(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800400001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")
	token := generateTestToken(adminID, "admin")

	// The old shape: a thousand credits for one rupee.
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/sms/buy",
		map[string]interface{}{
			"quantity": 1000, "rate": 0.001, "amount": 1,
			"payment_method": "cash",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var body map[string]interface{}
	parseJSON(t, resp, &body)

	// The named price is ignored; the server's rate applies.
	if amount, _ := body["amount"].(float64); amount <= 1 {
		t.Errorf("client-supplied price was honoured: recorded amount %v for 1000 credits", amount)
	}
}
