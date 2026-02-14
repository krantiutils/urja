package e2e

import (
	"net/http"
	"testing"
)

func TestDues_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestDues_List_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")

	memberID := createTestUser(t, "9800500002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestDues_Pay(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")

	memberID := createTestUser(t, "9800500002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	dueID := createTestDue(t, orgID, memberID, 500)

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues/"+memberID+"/pay",
		map[string]interface{}{
			"due_id":         dueID,
			"amount":         500.0,
			"payment_method": "cash",
		}, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestDues_Pay_MissingFields(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues/someid/pay",
		map[string]interface{}{}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestDues_BlockAccess_AdminOnly(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")

	memberID := createTestUser(t, "9800500002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")
	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/dues/"+memberID+"/block-access", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestDues_BlockAccess_StaffForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")

	staffID := createTestUser(t, "9800500003", "Staff")
	createTestOrgMember(t, staffID, orgID, "staff")
	token := generateTestToken(staffID, "staff")

	memberID := createTestUser(t, "9800500002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/dues/"+memberID+"/block-access", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}
