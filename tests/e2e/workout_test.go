package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// --- Template helpers ---

func createTestWorkoutTemplate(t *testing.T, orgID, createdBy, name string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO workout_templates (organization_id, name, category, difficulty, duration_minutes, exercises, created_by)
		 VALUES ($1, $2, 'strength', 'beginner', 30, '[{"name":"Push-ups","sets":3,"reps":10}]'::jsonb, $3)
		 RETURNING id`,
		orgID, name, createdBy,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestWorkoutTemplate(%s): %v", name, err)
	}
	return id
}

// --- Template CRUD tests ---

func TestWorkoutTemplate_List(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	createTestWorkoutTemplate(t, orgID, adminID, "Test Template")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/workout-templates", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) < 1 {
		t.Errorf("expected at least 1 template, got %d", len(data))
	}
}

func TestWorkoutTemplate_List_IncludesPresets(t *testing.T) {
	cleanupTables(t)

	// Re-seed presets since cleanupTables truncates workout_templates
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO workout_templates (organization_id, name, category, difficulty, duration_minutes, exercises, is_preset)
		 VALUES (NULL, 'Preset Template', 'cardio', 'beginner', 20, '[]'::jsonb, true)`)
	if err != nil {
		t.Fatalf("seeding preset: %v", err)
	}

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/workout-templates", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})

	foundPreset := false
	for _, item := range data {
		tpl := item.(map[string]interface{})
		if tpl["is_preset"] == true {
			foundPreset = true
			break
		}
	}
	if !foundPreset {
		t.Error("expected preset templates to be included in listing")
	}
}

func TestWorkoutTemplate_Get(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Get Template")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/workout-templates/"+templateID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	var tpl map[string]interface{}
	parseJSON(t, resp, &tpl)
	if tpl["name"] != "Get Template" {
		t.Errorf("expected name 'Get Template', got %v", tpl["name"])
	}
}

func TestWorkoutTemplate_Get_NotFound(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/workout-templates/00000000-0000-0000-0000-000000000000", nil, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestWorkoutTemplate_Create(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/workout-templates",
		map[string]interface{}{
			"name":             "New Template",
			"category":         "strength",
			"difficulty":       "intermediate",
			"duration_minutes": 45,
			"exercises":        []map[string]interface{}{{"name": "Squats", "sets": 3, "reps": 12}},
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var tpl map[string]interface{}
	parseJSON(t, resp, &tpl)
	if tpl["name"] != "New Template" {
		t.Errorf("expected name 'New Template', got %v", tpl["name"])
	}
	if tpl["is_preset"] != false {
		t.Error("expected is_preset to be false for org template")
	}
}

func TestWorkoutTemplate_Create_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/workout-templates",
		map[string]interface{}{
			"name":      "Should Fail",
			"exercises": []map[string]interface{}{},
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestWorkoutTemplate_Create_ValidationError(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	token := generateTestToken(adminID, "admin")

	// Empty name
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/workout-templates",
		map[string]interface{}{"name": ""}, token)
	assertStatus(t, resp, http.StatusBadRequest)

	// Invalid difficulty
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/workout-templates",
		map[string]interface{}{"name": "Test", "difficulty": "expert"}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestWorkoutTemplate_Update(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Old Name")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/workout-templates/"+templateID,
		map[string]interface{}{
			"name":       "Updated Name",
			"difficulty": "advanced",
		}, token)
	assertStatus(t, resp, http.StatusOK)

	var tpl map[string]interface{}
	parseJSON(t, resp, &tpl)
	if tpl["name"] != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %v", tpl["name"])
	}
}

func TestWorkoutTemplate_Delete(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "To Delete")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/workout-templates/"+templateID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	// Verify it's gone
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/workout-templates/"+templateID, nil, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestWorkoutTemplate_Delete_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "No Delete")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/workout-templates/"+templateID, nil, token)
	assertStatus(t, resp, http.StatusForbidden)
}

// --- Plan assignment tests ---

func TestWorkout_AssignPlan(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Plan Template")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/assign-plan",
		map[string]string{"workout_template_id": templateID}, token)
	assertStatus(t, resp, http.StatusOK)

	var plan map[string]interface{}
	parseJSON(t, resp, &plan)
	if plan["workout_template_id"] != templateID {
		t.Errorf("expected template_id %s, got %v", templateID, plan["workout_template_id"])
	}
}

func TestWorkout_AssignPlan_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Plan Template")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/assign-plan",
		map[string]string{"workout_template_id": templateID}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestWorkout_GetMyPlan(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "My Plan Template")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	// Assign plan directly
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO member_workout_plans (user_id, organization_id, workout_template_id, assigned_by)
		 VALUES ($1, $2, $3, $4)`,
		memberID, orgID, templateID, adminID)
	if err != nil {
		t.Fatalf("assigning plan: %v", err)
	}

	token := generateTestToken(memberID, "member")
	resp := doRequest(t, http.MethodGet,
		"/api/v1/members/me/workout-plan?organization_id="+orgID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	var plan map[string]interface{}
	parseJSON(t, resp, &plan)
	if plan["workout_template_id"] != templateID {
		t.Errorf("expected template_id %s, got %v", templateID, plan["workout_template_id"])
	}
	// Should include embedded template
	tpl := plan["template"].(map[string]interface{})
	if tpl["name"] != "My Plan Template" {
		t.Errorf("expected template name 'My Plan Template', got %v", tpl["name"])
	}
}

func TestWorkout_GetMyPlan_NoPlan(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	token := generateTestToken(memberID, "member")
	resp := doRequest(t, http.MethodGet,
		"/api/v1/members/me/workout-plan?organization_id="+orgID, nil, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestWorkout_GetMyPlan_MissingOrgID(t *testing.T) {
	cleanupTables(t)

	memberID := createTestUser(t, "9800800002", "Member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/workout-plan", nil, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

// --- Workout logging tests ---

func TestWorkout_CreateLog(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	exercises := []map[string]interface{}{
		{"name": "Push-ups", "sets": 3, "reps": 15},
		{"name": "Squats", "sets": 3, "reps": 12},
	}

	resp := doRequest(t, http.MethodPost, "/api/v1/members/me/workout-logs",
		map[string]interface{}{
			"organization_id": orgID,
			"exercises":       exercises,
			"duration_minutes": 30,
			"notes":           "Good session",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var log map[string]interface{}
	parseJSON(t, resp, &log)
	if log["organization_id"] != orgID {
		t.Errorf("expected org_id %s, got %v", orgID, log["organization_id"])
	}
	if log["notes"] != "Good session" {
		t.Errorf("expected notes 'Good session', got %v", log["notes"])
	}
}

func TestWorkout_CreateLog_WithTemplate(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Log Template")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/members/me/workout-logs",
		map[string]interface{}{
			"organization_id":    orgID,
			"workout_template_id": templateID,
			"exercises":          []map[string]interface{}{{"name": "Push-ups", "sets": 3, "reps": 10}},
			"duration_minutes":   25,
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var log map[string]interface{}
	parseJSON(t, resp, &log)
	if log["workout_template_id"] != templateID {
		t.Errorf("expected template_id %s, got %v", templateID, log["workout_template_id"])
	}
}

func TestWorkout_CreateLog_MissingOrgID(t *testing.T) {
	cleanupTables(t)

	memberID := createTestUser(t, "9800800002", "Member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/members/me/workout-logs",
		map[string]interface{}{
			"exercises": []map[string]interface{}{},
		}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestWorkout_ListLogs(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	// Insert logs directly
	for i := 0; i < 3; i++ {
		_, err := testPool.Exec(context.Background(),
			`INSERT INTO workout_logs (user_id, organization_id, exercises, duration_minutes)
			 VALUES ($1, $2, '[]'::jsonb, 30)`,
			memberID, orgID)
		if err != nil {
			t.Fatalf("inserting log: %v", err)
		}
	}

	token := generateTestToken(memberID, "member")
	resp := doRequest(t, http.MethodGet,
		"/api/v1/members/me/workout-logs?organization_id="+orgID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 3 {
		t.Errorf("expected 3 logs, got %d", len(data))
	}
}

func TestWorkout_ListLogs_MissingOrgID(t *testing.T) {
	cleanupTables(t)

	memberID := createTestUser(t, "9800800002", "Member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/members/me/workout-logs", nil, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestWorkout_ListLogs_DateFiltering(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")

	// Insert a log with a specific date
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO workout_logs (user_id, organization_id, exercises, duration_minutes, logged_at)
		 VALUES ($1, $2, '[]'::jsonb, 30, '2026-01-15 10:00:00+00')`,
		memberID, orgID)
	if err != nil {
		t.Fatalf("inserting log: %v", err)
	}

	token := generateTestToken(memberID, "member")

	// Query with date range that includes the log
	resp := doRequest(t, http.MethodGet,
		"/api/v1/members/me/workout-logs?organization_id="+orgID+"&from=2026-01-01&to=2026-01-31",
		nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 log in date range, got %d", len(data))
	}

	// Query with date range that excludes the log
	resp = doRequest(t, http.MethodGet,
		"/api/v1/members/me/workout-logs?organization_id="+orgID+"&from=2026-02-01&to=2026-02-28",
		nil, token)
	assertStatus(t, resp, http.StatusOK)

	parseJSON(t, resp, &body)
	// data can be nil when no results
	rawData := body["data"]
	if rawData != nil {
		data = rawData.([]interface{})
		if len(data) != 0 {
			t.Errorf("expected 0 logs in date range, got %d", len(data))
		}
	}
}

// --- Staff can create templates ---

func TestWorkoutTemplate_Create_StaffAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")

	staffID := createTestUser(t, "9800800002", "Staff")
	createTestOrgMember(t, staffID, orgID, "staff")
	token := generateTestToken(staffID, "staff")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/workout-templates",
		map[string]interface{}{
			"name":      "Staff Template",
			"exercises": json.RawMessage(`[]`),
		}, token)
	assertStatus(t, resp, http.StatusCreated)
}

// --- Staff can assign plans ---

func TestWorkout_AssignPlan_StaffAllowed(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Workout Gym")
	templateID := createTestWorkoutTemplate(t, orgID, adminID, "Assign Template")

	staffID := createTestUser(t, "9800800003", "Staff")
	createTestOrgMember(t, staffID, orgID, "staff")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(staffID, "staff")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/assign-plan",
		map[string]string{"workout_template_id": templateID}, token)
	assertStatus(t, resp, http.StatusOK)
}
