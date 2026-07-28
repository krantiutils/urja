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

// Until dues could be created, the whole feature was unreachable: the table
// could be listed, paid and used to block access, but nothing in the product
// ever wrote to it.
func TestDues_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	memberID := createTestUser(t, "9800900002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues",
		map[string]interface{}{
			"user_id": memberID, "amount": 2500,
			"due_date": "2026-08-15", "description": "Monthly fee",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	list := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues", nil, token)
	assertStatus(t, list, http.StatusOK)
	var body map[string]interface{}
	parseJSON(t, list, &body)
	if data := body["data"].([]interface{}); len(data) != 1 {
		t.Fatalf("expected the due to be listed, got %d", len(data))
	}
}

// Staff of one gym must not be able to raise a due against a stranger.
func TestDues_Create_RejectsNonMember(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	outsider := createTestUser(t, "9800900003", "Outsider")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues",
		map[string]interface{}{"user_id": outsider, "amount": 100, "due_date": "2026-08-15"}, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestDues_Create_RejectsBadAmount(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	memberID := createTestUser(t, "9800900002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues",
		map[string]interface{}{"user_id": memberID, "amount": 0, "due_date": "2026-08-15"}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

// A package assigned with part of the price paid should leave the remainder
// recorded as money owed, rather than losing it.
func TestSubscription_PartialPayment_RaisesDue(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800900001", "Admin")
	orgID := createTestOrg(t, adminID, "Dues Gym")
	memberID := createTestUser(t, "9800900002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	pkgID := createTestPackage(t, orgID, "Monthly", 30, 3000)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{
			"package_id": pkgID, "start_date": "2026-08-01",
			"payment_method": "cash", "amount_paid": 1000,
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	list := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues", nil, token)
	var body map[string]interface{}
	parseJSON(t, list, &body)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected one due for the unpaid balance, got %d", len(data))
	}
	if amount := data[0].(map[string]interface{})["amount"].(float64); amount != 2000 {
		t.Errorf("expected a balance of 2000, got %v", amount)
	}
}
