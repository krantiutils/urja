package e2e

import (
	"net/http"
	"testing"
)

func TestNFC_ListCards(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801100001", "Admin")
	orgID := createTestOrg(t, adminID, "NFC Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/nfc-cards", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestNFC_RegisterCard(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801100001", "Admin")
	orgID := createTestOrg(t, adminID, "NFC Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/nfc-cards",
		map[string]string{"card_hex": "AABBCCDD"}, token)
	assertStatus(t, resp, http.StatusCreated)
}

func TestNFC_AssignCard(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801100001", "Admin")
	orgID := createTestOrg(t, adminID, "NFC Gym")

	memberID := createTestUser(t, "9801100002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	cardID := createTestNFCCard(t, orgID, "11223344")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPut,
		"/api/v1/orgs/"+orgID+"/nfc-cards/"+cardID+"/assign",
		map[string]string{"user_id": memberID}, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestNFC_UnassignCard(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801100001", "Admin")
	orgID := createTestOrg(t, adminID, "NFC Gym")

	memberID := createTestUser(t, "9801100002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	cardID := createTestNFCCard(t, orgID, "55667788")
	token := generateTestToken(adminID, "admin")

	// Assign first
	resp := doRequest(t, http.MethodPut,
		"/api/v1/orgs/"+orgID+"/nfc-cards/"+cardID+"/assign",
		map[string]string{"user_id": memberID}, token)
	assertStatus(t, resp, http.StatusOK)

	// Then unassign
	resp = doRequest(t, http.MethodPut,
		"/api/v1/orgs/"+orgID+"/nfc-cards/"+cardID+"/unassign", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestNFC_ListDevices(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801100001", "Admin")
	orgID := createTestOrg(t, adminID, "NFC Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/nfc-devices", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}
