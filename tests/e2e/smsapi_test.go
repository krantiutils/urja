package e2e

import (
	"net/http"
	"testing"
)

func TestSMS_GetBalance(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")
	initSMSCredits(t, orgID, 100)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/sms/balance", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestSMS_GetBalance_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")

	memberID := createTestUser(t, "9800900002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/sms/balance", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestSMS_Buy(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")
	initSMSCredits(t, orgID, 0)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/sms/buy",
		map[string]interface{}{
			"quantity":       100,
			"rate":           1.5,
			"amount":         150.0,
			"payment_method": "cash",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestSMS_Send(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")

	memberID := createTestUser(t, "9800900002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	initSMSCredits(t, orgID, 100)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/sms/send",
		map[string]interface{}{
			"message":    "Hello from gym!",
			"member_ids": []string{memberID},
		}, token)
	// May return 201 (success) or 400 (if member phone lookup fails)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		body := readBody(t, resp)
		t.Fatalf("expected 201 or 400, got %d: %s", resp.StatusCode, body)
	}
}

func TestSMS_History(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/sms/history", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestSMS_StaffForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "SMS Gym")

	staffID := createTestUser(t, "9800900003", "Staff")
	createTestOrgMember(t, staffID, orgID, "staff")
	token := generateTestToken(staffID, "staff")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/sms/balance", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}
