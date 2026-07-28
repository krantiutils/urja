package e2e

import (
	"net/http"
	"testing"
)

// A sweep rather than one-offs: every org-scoped endpoint below is either
// staff business or member business, and the difference has been got wrong
// often enough in this codebase to be worth asserting as a set.
//
// OrgScope proves membership, not authority. Anything that manages the gym —
// its people, its money, its records, its operations — must refuse a plain
// member even though that member legitimately belongs to the organization.
func TestAuthz_MemberIsRefusedOrgManagement(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Sweep Gym")
	memberID := createTestUser(t, "9800700002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	base := "/api/v1/orgs/" + orgID

	staffOnly := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"list members", http.MethodGet, base + "/members", nil},
		{"list staff", http.MethodGet, base + "/staff", nil},
		{"list dues", http.MethodGet, base + "/dues", nil},
		{"raise a due", http.MethodPost, base + "/dues",
			map[string]interface{}{"user_id": memberID, "amount": 100, "due_date": "2026-08-01"}},
		{"gym attendance", http.MethodGet, base + "/attendance", nil},
		{"manual check-in", http.MethodPost, base + "/attendance/check-in",
			map[string]string{"member_id": memberID}},
		{"accounts", http.MethodGet, base + "/accounts", nil},
		{"activity log", http.MethodGet, base + "/activity-logs", nil},
		{"check-in QR", http.MethodGet, base + "/qr-code", nil},
		{"NFC cards", http.MethodGet, base + "/nfc-cards", nil},
		{"site settings", http.MethodGet, base + "/site/settings", nil},
		{"site leads", http.MethodGet, base + "/site/leads", nil},
		{"SMS balance", http.MethodGet, base + "/sms/balance", nil},
		{"manage packages", http.MethodGet, base + "/packages", nil},
	}

	for _, tc := range staffOnly {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, tc.method, tc.path, tc.body, token)
			if resp.StatusCode != http.StatusForbidden {
				t.Errorf("a plain member reached %s %s: got %d, want 403",
					tc.method, tc.path, resp.StatusCode)
			}
		})
	}
}

// The other half of the same rule: a member must still be able to run their own
// membership, or the gating above has gone too far.
func TestAuthz_MemberKeepsTheirOwnRoutes(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Sweep Gym")
	memberID := createTestUser(t, "9800700002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	own := []struct {
		name string
		path string
	}{
		{"own profile", "/api/v1/members/me"},
		{"own attendance", "/api/v1/members/me/attendance"},
		{"own streaks", "/api/v1/members/me/streaks"},
		{"own packages", "/api/v1/members/me/packages"},
		{"the leaderboard", "/api/v1/orgs/" + orgID + "/leaderboard"},
	}

	for _, tc := range own {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, http.MethodGet, tc.path, nil, token)
			if resp.StatusCode == http.StatusForbidden {
				t.Errorf("a member was refused their own %s", tc.name)
			}
		})
	}
}
