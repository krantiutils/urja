package e2e

import (
	"net/http"
	"testing"
)

func TestActivityLog_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801200001", "Admin")
	orgID := createTestOrg(t, adminID, "Log Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/activity-logs", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	if _, ok := body["data"]; !ok {
		t.Error("expected 'data' field")
	}
}

func TestActivityLog_ListWithDateFilter(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801200001", "Admin")
	orgID := createTestOrg(t, adminID, "Log Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/activity-logs?from=2026-01-01&to=2026-12-31", nil, token)
	assertStatus(t, resp, http.StatusOK)
}

func TestActivityLog_ListInvalidDateFormat(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9801200001", "Admin")
	orgID := createTestOrg(t, adminID, "Log Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/activity-logs?from=not-a-date", nil, token)
	assertStatus(t, resp, http.StatusBadRequest)
}
