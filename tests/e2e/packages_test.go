package e2e

import (
	"net/http"
	"testing"
)

func TestPackages_ListPublic(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")
	createTestPackage(t, orgID, "Basic Plan", 30, 2000)
	createTestPackage(t, orgID, "Premium Plan", 90, 5000)

	resp := doRequest(t, http.MethodGet, "/api/v1/packages", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) < 2 {
		t.Errorf("expected at least 2 packages, got %d", len(data))
	}
}

func TestPackages_ListByOrg_AdminOnly(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")
	createTestPackage(t, orgID, "Plan A", 30, 1000)

	token := generateTestToken(adminID, "admin")
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/packages", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestPackages_ListByOrg_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")

	memberID := createTestUser(t, "9800200002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(memberID, "member")
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/packages", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestPackages_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/packages",
		map[string]interface{}{
			"name":          "Monthly Basic",
			"duration_days": 30,
			"price":         "1500",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "Monthly Basic" {
		t.Errorf("expected 'Monthly Basic', got %v", body["name"])
	}
}

func TestPackages_Update(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")
	pkgID := createTestPackage(t, orgID, "Old Name", 30, 1000)
	token := generateTestToken(adminID, "admin")

	newName := "New Name"
	newPrice := "1500"
	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/packages/"+pkgID,
		map[string]interface{}{
			"name":  &newName,
			"price": &newPrice,
		}, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "New Name" {
		t.Errorf("expected 'New Name', got %v", body["name"])
	}
}

func TestPackages_Delete(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")
	pkgID := createTestPackage(t, orgID, "To Delete", 30, 1000)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/packages/"+pkgID, nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestPackages_ListMyPackages(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800200001", "Admin")
	orgID := createTestOrg(t, adminID, "Pkg Gym")

	memberID := createTestUser(t, "9800200002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	pkgID := createTestPackage(t, orgID, "Monthly", 30, 1000)
	assignTestPackage(t, memberID, pkgID, orgID)

	token := generateTestToken(memberID, "member")
	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/packages", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) < 1 {
		t.Errorf("expected at least 1 package, got %d", len(data))
	}
}
