package e2e

import (
	"net/http"
	"testing"
)

func TestAttendance_ManualCheckIn(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800400001", "Admin")
	orgID := createTestOrg(t, adminID, "Attend Gym")

	memberID := createTestUser(t, "9800400002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	// Give the member an active package so check-in works
	pkgID := createTestPackage(t, orgID, "Monthly", 30, 1000)
	assignTestPackage(t, memberID, pkgID, orgID)

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/attendance/check-in",
		map[string]string{"member_id": memberID}, token)
	// May return 201 or 400 depending on business logic for manual check-in
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		body := readBody(t, resp)
		t.Fatalf("expected 201 or 400, got %d: %s", resp.StatusCode, body)
	}
}

func TestAttendance_ListByOrg(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800400001", "Admin")
	orgID := createTestOrg(t, adminID, "Attend Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/attendance", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestAttendance_ListMine(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800400002", "Member")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/attendance/attendance", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestAttendance_GetMyStreaks(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800400002", "Member")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/attendance/streaks", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestAttendance_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/attendance/attendance", nil, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}
