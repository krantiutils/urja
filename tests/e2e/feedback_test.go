package e2e

import (
	"net/http"
	"testing"
)

func TestFeedback_Create_AsMember(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Feedback Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/feedbacks",
		map[string]string{"message": "Great gym!"}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestFeedback_Create_EmptyMessage(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Feedback Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/feedbacks",
		map[string]string{"message": ""}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestFeedback_List_AdminAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Feedback Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/feedbacks", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestFeedback_List_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Feedback Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/feedbacks", nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}
