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

// The QR carries a signed, five-minute check-in token. That rotation is what
// makes scanning it stand in for being at the gym — but it protects nothing
// against somebody who can mint a fresh one whenever they like. A member with
// this endpoint can check in from home indefinitely, and attendance drives
// streaks, the leaderboard and absentee SMS.
//
// This previously asserted the opposite. Nothing in the web app requests the
// code; it is for a screen at the desk, which is staff.
func TestQRCode_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801300001", "Admin")
	orgID := createTestOrg(t, adminID, "QR Gym")

	memberID := createTestUser(t, "9801300002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/qr-code?format=json", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}
