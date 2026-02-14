package e2e

import (
	"net/http"
	"testing"
)

func TestOrg_ListGyms_Public(t *testing.T) {
	cleanupTables(t)

	// Create some data
	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	createTestOrg(t, adminID, "Gym Alpha")
	createTestOrg(t, adminID, "Gym Beta")

	resp := doRequest(t, http.MethodGet, "/api/v1/gyms", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data, ok := body["data"].([]interface{})
	if !ok {
		t.Fatal("expected data array in response")
	}
	if len(data) < 2 {
		t.Errorf("expected at least 2 gyms, got %d", len(data))
	}
}

func TestOrg_GetGym_Public(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	orgID := createTestOrg(t, adminID, "Gym Gamma")

	resp := doRequest(t, http.MethodGet, "/api/v1/gyms/"+orgID, nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "Gym Gamma" {
		t.Errorf("expected name 'Gym Gamma', got %v", body["name"])
	}
}

func TestOrg_GetGym_NotFound(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/gyms/00000000-0000-0000-0000-000000000000", nil, "")
	assertStatus(t, resp, http.StatusNotFound)
}

func TestOrg_Create_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs",
		map[string]string{"name": "New Gym"}, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestOrg_Create_NotSuperAdmin(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800000001", "Regular User")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs",
		map[string]string{"name": "New Gym"}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestOrg_Create_AsSuperAdmin(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs",
		map[string]string{"name": "Brand New Gym", "slug": "brand-new-gym"}, token)
	assertStatus(t, resp, http.StatusCreated)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "Brand New Gym" {
		t.Errorf("expected name 'Brand New Gym', got %v", body["name"])
	}
}

func TestOrg_Create_MissingName(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs",
		map[string]string{"slug": "no-name"}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestOrg_Update_AsAdmin(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	orgID := createTestOrg(t, adminID, "Old Name")
	token := generateTestToken(adminID, "admin")

	newName := "Updated Name"
	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
		map[string]interface{}{"name": &newName}, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestOrg_Update_NotMember(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800000001", "Super Admin")
	orgID := createTestOrg(t, adminID, "Some Gym")

	outsiderID := createTestUser(t, "9800000002", "Outsider")
	token := generateTestToken(outsiderID, "member")

	newName := "Hacked Name"
	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
		map[string]interface{}{"name": &newName}, token)
	assertStatus(t, resp, http.StatusForbidden)
}
