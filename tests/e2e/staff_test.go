package e2e

import (
	"net/http"
	"testing"
)

func TestStaff_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/staff", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field in response")
	}
}

func TestStaff_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff",
		map[string]string{
			"phone":      "9800300099",
			"name":       "New Staff",
			"staff_role": "trainer",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestStaff_Create_MissingFields(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff",
		map[string]string{"phone": "9800300099"}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestStaff_Create_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")

	memberID := createTestUser(t, "9800300002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff",
		map[string]string{
			"phone":      "9800300099",
			"name":       "Should Fail",
			"staff_role": "trainer",
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestStaff_GetUpdateDelete(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	// Create staff
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff",
		map[string]string{
			"phone":      "9800300088",
			"name":       "Staff Person",
			"staff_role": "receptionist",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var created map[string]interface{}
	parseJSON(t, resp, &created)
	staffID := created["id"].(string)

	// Get
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/staff/"+staffID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	// Update
	newName := "Updated Staff"
	resp = doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/staff/"+staffID,
		map[string]interface{}{"name": &newName}, token)
	assertStatus(t, resp, http.StatusOK)

	// Delete
	resp = doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/staff/"+staffID, nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestStaff_RecordAttendance(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	// Create staff via API to get the organization_members.id
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff",
		map[string]string{
			"phone":      "9800300077",
			"name":       "Staff Bob",
			"staff_role": "trainer",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var staffBody map[string]interface{}
	parseJSON(t, resp, &staffBody)
	staffMemberID := staffBody["id"].(string)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/staff/attendance",
		map[string]string{
			"staff_id": staffMemberID,
			"action":   "check_in",
			"method":   "manual",
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestStaff_ListAttendance(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800300001", "Admin")
	orgID := createTestOrg(t, adminID, "Staff Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/staff/attendance", nil, token)
	assertStatus(t, resp, http.StatusOK)
}
