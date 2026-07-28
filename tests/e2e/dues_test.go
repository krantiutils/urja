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

	// The due goes in the path — it is a due id, not a member id.
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues/"+dueID+"/pay",
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

// Recording a payment is the whole point of the Due Payments screen, and it
// could never have worked: the route declared the due in the path while the
// handler read it only from the body, so every attempt from the dashboard
// returned "due_id is required".
func TestDues_Pay_UsesPathParam(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800950001", "Admin")
	orgID := createTestOrg(t, adminID, "Pay Gym")
	memberID := createTestUser(t, "9800950002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(adminID, "admin")

	created := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues",
		map[string]interface{}{"user_id": memberID, "amount": 2000, "due_date": "2026-08-15"}, token)
	assertStatus(t, created, http.StatusCreated)
	var due map[string]interface{}
	parseJSON(t, created, &due)

	// Exactly what the web app sends: the due in the path, nothing but the
	// payment in the body.
	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/dues/"+due["id"].(string)+"/pay",
		map[string]interface{}{"amount": 2000, "payment_method": "cash"}, token)
	assertStatus(t, resp, http.StatusOK)

	unpaid := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues?status=unpaid", nil, token)
	var body map[string]interface{}
	parseJSON(t, unpaid, &body)
	if data, _ := body["data"].([]interface{}); len(data) != 0 {
		t.Errorf("expected the due to be settled, %d still unpaid", len(data))
	}
}

// Money taken at the desk has to reach the books. Package payments were
// written to `payments` but never to `transactions`, which is what the accounts
// screen sums — so a gym's reported income showed dues collections only and
// missed most of its revenue.
func TestSubscription_PackagePayment_ReachesTheBooks(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800960001", "Admin")
	orgID := createTestOrg(t, adminID, "Books Gym")
	memberID := createTestUser(t, "9800960002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	pkgID := createTestPackage(t, orgID, "Monthly", 30, 3000)
	token := generateTestToken(adminID, "admin")

	assign := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{
			"package_id": pkgID, "start_date": "2026-08-01",
			"payment_method": "cash", "amount_paid": 1000,
		}, token)
	assertStatus(t, assign, http.StatusCreated)

	summary := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts/summary", nil, token)
	assertStatus(t, summary, http.StatusOK)
	var books map[string]interface{}
	parseJSON(t, summary, &books)

	if income, _ := books["total_income"].(float64); income != 1000 {
		t.Errorf("expected the 1000 taken at signup to be income, books show %v", income)
	}
}

// The full sequence: part paid at signup, the rest settled later. Both halves
// belong in the books, and together they should equal the package price.
func TestSubscription_PartPaidThenSettled_BooksBalance(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800970001", "Admin")
	orgID := createTestOrg(t, adminID, "Balance Gym")
	memberID := createTestUser(t, "9800970002", "Ram")
	createTestOrgMember(t, memberID, orgID, "member")
	pkgID := createTestPackage(t, orgID, "Monthly", 30, 3000)
	token := generateTestToken(adminID, "admin")

	doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{
			"package_id": pkgID, "start_date": "2026-08-01",
			"payment_method": "cash", "amount_paid": 1000,
		}, token)

	list := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/dues", nil, token)
	var duesBody map[string]interface{}
	parseJSON(t, list, &duesBody)
	data := duesBody["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected a due for the balance, got %d", len(data))
	}
	dueID := data[0].(map[string]interface{})["id"].(string)

	pay := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/dues/"+dueID+"/pay",
		map[string]interface{}{"amount": 2000, "payment_method": "cash"}, token)
	assertStatus(t, pay, http.StatusOK)

	summary := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/accounts/summary", nil, token)
	var books map[string]interface{}
	parseJSON(t, summary, &books)

	if income, _ := books["total_income"].(float64); income != 3000 {
		t.Errorf("1000 at signup plus 2000 settled should be 3000 income, books show %v", income)
	}
}
