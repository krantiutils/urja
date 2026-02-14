package e2e

import (
	"net/http"
	"testing"
)

func TestNotice_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")

	createTestNotice(t, orgID, adminID, "Test Notice", "Notice content")

	token := generateTestToken(adminID, "admin")
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/notices", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) < 1 {
		t.Errorf("expected at least 1 notice, got %d", len(data))
	}
}

func TestNotice_List_MemberAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")

	memberID := createTestUser(t, "9800700002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/notices", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestNotice_Get(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")
	noticeID := createTestNotice(t, orgID, adminID, "Gym Closed", "Closed for holiday")

	token := generateTestToken(adminID, "admin")
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/notices/"+noticeID, nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestNotice_Get_NotFound(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/notices/00000000-0000-0000-0000-000000000000", nil, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestNotice_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/notices",
		map[string]string{
			"title":   "New Notice",
			"content": "Content here",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestNotice_Create_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")

	memberID := createTestUser(t, "9800700002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/notices",
		map[string]string{
			"title":   "Should Fail",
			"content": "Nope",
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestNotice_Update(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")
	noticeID := createTestNotice(t, orgID, adminID, "Old Title", "Old content")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/notices/"+noticeID,
		map[string]string{
			"title":   "Updated Title",
			"content": "Updated content",
		}, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestNotice_Delete(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800700001", "Admin")
	orgID := createTestOrg(t, adminID, "Notice Gym")
	noticeID := createTestNotice(t, orgID, adminID, "To Delete", "Will be deleted")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/notices/"+noticeID, nil, token)
	assertStatus(t, resp, http.StatusOK)
}
