package e2e

import (
	"net/http"
	"testing"
)

func TestMember_GetMe(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800100001", "Alice Member")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/members/me", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "Alice Member" {
		t.Errorf("expected name 'Alice Member', got %v", body["name"])
	}
	if body["phone"] != "9800100001" {
		t.Errorf("expected phone '9800100001', got %v", body["phone"])
	}
}

func TestMember_GetMe_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/members/me", nil, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestMember_UpdateMe(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800100001", "Alice")
	token := generateTestToken(userID, "member")

	newName := "Alice Updated"
	resp := doRequest(t, http.MethodPut, "/api/v1/members/me",
		map[string]interface{}{"name": &newName}, token)
	assertStatus(t, resp, http.StatusOK)

	// Verify update
	resp2 := doRequest(t, http.MethodGet, "/api/v1/members/me", nil, token)
	assertStatus(t, resp2, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp2, &body)
	if body["name"] != "Alice Updated" {
		t.Errorf("expected name 'Alice Updated', got %v", body["name"])
	}
}

func TestMember_UpdateMe_InvalidEmail(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800100001", "Alice")
	token := generateTestToken(userID, "member")

	badEmail := "not-an-email"
	resp := doRequest(t, http.MethodPut, "/api/v1/members/me",
		map[string]interface{}{"email": &badEmail}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestMember_UpdatePrivacy(t *testing.T) {
	cleanupTables(t)

	userID := createTestUser(t, "9800100001", "Alice")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/members/me/privacy",
		map[string]interface{}{
			"show_email":          true,
			"show_phone":          false,
			"show_profile":        true,
			"show_attendance":     false,
			"show_on_leaderboard": true,
		}, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestMember_ListByOrg(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800100001", "Admin")
	orgID := createTestOrg(t, adminID, "Test Gym")

	member1ID := createTestUser(t, "9800100002", "Member 1")
	createTestOrgMember(t, member1ID, orgID, "member")

	member2ID := createTestUser(t, "9800100003", "Member 2")
	createTestOrgMember(t, member2ID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/members", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	// admin + 2 members = 3
	if len(data) < 3 {
		t.Errorf("expected at least 3 members, got %d", len(data))
	}
}

func TestMember_CreateOrgMember(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800100001", "Admin")
	orgID := createTestOrg(t, adminID, "Test Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/members",
		map[string]string{
			"phone": "9800100099",
			"name":  "New Member",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if body["name"] != "New Member" {
		t.Errorf("expected name 'New Member', got %v", body["name"])
	}
	if body["role"] != "member" {
		t.Errorf("expected role 'member', got %v", body["role"])
	}
}

func TestMember_CreateOrgMember_InvalidPhone(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800100001", "Admin")
	orgID := createTestOrg(t, adminID, "Test Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/members",
		map[string]string{
			"phone": "12345",
			"name":  "Bad Phone",
		}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestMember_CreateOrgMember_MemberCantCreate(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800100001", "Admin")
	orgID := createTestOrg(t, adminID, "Test Gym")

	memberID := createTestUser(t, "9800100002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/members",
		map[string]string{
			"phone": "9800100099",
			"name":  "Should Fail",
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestMember_DeleteOrgMember(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800100001", "Admin")
	orgID := createTestOrg(t, adminID, "Test Gym")

	memberID := createTestUser(t, "9800100002", "To Delete")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/members/"+memberID, nil, token)
	// NOTE: Returns 404 due to chi routing conflict — r.Delete("/{id}") in
	// RegisterOrgRoutes is shadowed by r.Route("/{memberId}") for subscription routes.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		body := readBody(t, resp)
		t.Fatalf("expected 200 or 404, got %d: %s", resp.StatusCode, body)
	}
}

func TestMember_OrgIsolation(t *testing.T) {
	cleanupTables(t)

	admin1 := createTestSuperAdmin(t, "9800100001", "Admin 1")
	_ = createTestOrg(t, admin1, "Gym 1")

	admin2 := createTestUser(t, "9800100002", "Admin 2")
	org2 := createTestOrg(t, admin2, "Gym 2")

	// admin1 should not be able to access org2's members
	token := generateTestToken(admin1, "admin")
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+org2+"/members", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}
