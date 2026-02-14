package e2e

import (
	"net/http"
	"testing"
)

func TestQRCode_GeneratePNG(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801300001", "Admin")
	orgID := createTestOrg(t, adminID, "QR Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/qr-code", nil, token)
	assertStatus(t, resp, http.StatusOK)

	ct := resp.Header.Get("Content-Type")
	if ct != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", ct)
	}
}

func TestQRCode_GenerateJSON(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801300001", "Admin")
	orgID := createTestOrg(t, adminID, "QR Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/qr-code?format=json", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["org_id"] != orgID {
		t.Errorf("expected org_id %s, got %v", orgID, body["org_id"])
	}
	if _, ok := body["checkin_url"]; !ok {
		t.Error("expected 'checkin_url' field")
	}
}

func TestQRCode_MemberAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801300001", "Admin")
	orgID := createTestOrg(t, adminID, "QR Gym")

	memberID := createTestUser(t, "9801300002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/qr-code?format=json", nil, token)
	assertStatus(t, resp, http.StatusOK)
}
