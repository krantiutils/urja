package e2e

import (
	"net/http"
	"testing"
)

func TestSubscription_PackageSummary(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")
	createTestPackage(t, orgID, "Monthly", 30, 1000)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/packages/summary", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestSubscription_ListExpiring(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/packages/expiring?days=7", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestSubscription_ListExpired(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/packages/expired", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestSubscription_AssignPackage(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")
	pkgID := createTestPackage(t, orgID, "Monthly", 30, 1000)

	memberID := createTestUser(t, "9801400002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{
			"package_id":     pkgID,
			"start_date":     "2026-02-14",
			"payment_method": "cash",
			"amount_paid":    1000.0,
			"discount":       0.0,
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestSubscription_AssignPackage_MissingFields(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")

	memberID := createTestUser(t, "9801400002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestSubscription_ListPayments(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")

	memberID := createTestUser(t, "9801400002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/payments", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestSubscription_ListSubscriptions(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801400001", "Admin")
	orgID := createTestOrg(t, adminID, "Sub Gym")

	memberID := createTestUser(t, "9801400002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/subscriptions", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}
