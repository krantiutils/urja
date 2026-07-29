package e2e

import (
	"net/http"
	"testing"
)

func TestOrgSettings_SetPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000201", "Admin")
	orgID := createTestOrg(t, admin, "PAN Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "601234567",
		"tax_legal_name": "PAN Gym Pvt Ltd",
		"tax_address":    "Kirtipur, Kathmandu",
	}, token)
	assertStatus(t, resp, http.StatusOK)

	var got struct {
		PANNumber    string `json:"pan_number"`
		TaxLegalName string `json:"tax_legal_name"`
	}
	parseJSON(t, resp, &got)
	if got.PANNumber != "601234567" {
		t.Errorf("pan_number = %q, want %q", got.PANNumber, "601234567")
	}
	if got.TaxLegalName != "PAN Gym Pvt Ltd" {
		t.Errorf("tax_legal_name = %q, want %q", got.TaxLegalName, "PAN Gym Pvt Ltd")
	}
}

func TestOrgSettings_RejectsMalformedPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000202", "Admin")
	orgID := createTestOrg(t, admin, "Bad PAN Gym")
	token := generateTestToken(admin, "member")

	for _, bad := range []string{"12345", "abcdefghi", "6012345678"} {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
			map[string]any{"pan_number": bad}, token)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("PAN %q: status = %d, want 400", bad, resp.StatusCode)
		}
		resp.Body.Close()
	}
}
