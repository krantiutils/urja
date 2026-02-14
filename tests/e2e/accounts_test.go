package e2e

import (
	"net/http"
	"testing"
)

func TestAccounts_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestAccounts_List_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")

	memberID := createTestUser(t, "9800600002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestAccounts_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/accounts",
		map[string]interface{}{
			"category":         "membership",
			"description":      "Monthly fee",
			"transaction_date": "2026-02-14",
			"transaction_type": "income",
			"amount":           1500.0,
			"payment_type":     "cash",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestAccounts_Create_StaffAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")

	staffID := createTestUser(t, "9800600003", "Staff")
	createTestOrgMember(t, staffID, orgID, "staff")
	token := generateTestToken(staffID, "staff")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/accounts",
		map[string]interface{}{
			"category":         "other",
			"transaction_date": "2026-02-14",
			"transaction_type": "expense",
			"amount":           200.0,
			"payment_type":     "cash",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestAccounts_GetSummary(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts/summary", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestAccounts_Export(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Acct Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts/export", nil, token)
	assertStatus(t, resp, http.StatusOK)

	ct := resp.Header.Get("Content-Type")
	if ct != "text/csv" {
		t.Errorf("expected Content-Type text/csv, got %s", ct)
	}
}
