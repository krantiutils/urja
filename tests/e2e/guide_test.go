package e2e

import (
	"net/http"
	"testing"
)

func TestGuide_ListPublished(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	createTestGuide(t, adminID, "Published Guide", "Some content", "strength", true)
	createTestGuide(t, adminID, "Draft Guide", "Hidden content", "cardio", false)

	// Public endpoint - no auth required
	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 published guide, got %d", len(data))
	}
	total := int(body["total"].(float64))
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
}

func TestGuide_ListPublished_CategoryFilter(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	createTestGuide(t, adminID, "Strength Guide", "Content", "strength", true)
	createTestGuide(t, adminID, "Cardio Guide", "Content", "cardio", true)

	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides?category=strength", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 strength guide, got %d", len(data))
	}
}

func TestGuide_ListPublished_Search(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	createTestGuide(t, adminID, "Yoga Basics", "Flexibility content", "flexibility", true)
	createTestGuide(t, adminID, "Running Tips", "Cardio content", "cardio", true)

	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides?search=yoga", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 guide matching 'yoga', got %d", len(data))
	}
}

func TestGuide_GetPublished(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	guideID := createTestGuide(t, adminID, "Published Guide", "Content", "strength", true)

	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides/"+guideID, nil, "")
	assertStatus(t, resp, http.StatusOK)
}

func TestGuide_GetPublished_DraftNotFound(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	guideID := createTestGuide(t, adminID, "Draft Guide", "Hidden content", "strength", false)

	// Draft should not be accessible via public endpoint
	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides/"+guideID, nil, "")
	assertStatus(t, resp, http.StatusNotFound)
}

func TestGuide_GetPublished_NotFound(t *testing.T) {
	cleanupTables(t)

	resp := doRequest(t, http.MethodGet, "/api/v1/training-guides/00000000-0000-0000-0000-000000000000", nil, "")
	assertStatus(t, resp, http.StatusNotFound)
}

func TestGuide_AdminCreate(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/training-guides",
		map[string]string{
			"title":    "New Guide",
			"title_ne": "नयाँ गाइड",
			"content":  "Guide content here",
			"category": "strength",
		}, token)
	assertStatus(t, resp, http.StatusCreated)

	var guide map[string]interface{}
	parseJSON(t, resp, &guide)
	if guide["title"] != "New Guide" {
		t.Errorf("expected title 'New Guide', got %v", guide["title"])
	}
	if guide["is_published"] != false {
		t.Errorf("expected new guide to be unpublished")
	}
}

func TestGuide_AdminCreate_InvalidCategory(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/training-guides",
		map[string]string{
			"title":    "Bad Guide",
			"content":  "Content",
			"category": "invalid_category",
		}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestGuide_AdminCreate_MemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")

	memberID := createTestUser(t, "9800800002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/training-guides",
		map[string]string{
			"title":   "Should Fail",
			"content": "Nope",
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

func TestGuide_AdminUpdate(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	guideID := createTestGuide(t, adminID, "Old Title", "Old content", "strength", false)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/training-guides/"+guideID,
		map[string]string{
			"title":   "Updated Title",
			"content": "Updated content",
		}, token)
	assertStatus(t, resp, http.StatusOK)

	var guide map[string]interface{}
	parseJSON(t, resp, &guide)
	if guide["title"] != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %v", guide["title"])
	}
}

func TestGuide_AdminPublish(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	guideID := createTestGuide(t, adminID, "Guide", "Content", "strength", false)
	token := generateTestToken(adminID, "admin")

	// Publish
	resp := doRequest(t, http.MethodPatch, "/api/v1/orgs/"+orgID+"/training-guides/"+guideID+"/publish",
		map[string]interface{}{"is_published": true}, token)
	assertStatus(t, resp, http.StatusOK)

	var guide map[string]interface{}
	parseJSON(t, resp, &guide)
	if guide["is_published"] != true {
		t.Errorf("expected guide to be published")
	}

	// Unpublish
	resp = doRequest(t, http.MethodPatch, "/api/v1/orgs/"+orgID+"/training-guides/"+guideID+"/publish",
		map[string]interface{}{"is_published": false}, token)
	assertStatus(t, resp, http.StatusOK)

	parseJSON(t, resp, &guide)
	if guide["is_published"] != false {
		t.Errorf("expected guide to be unpublished")
	}
}

func TestGuide_AdminDelete(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	guideID := createTestGuide(t, adminID, "To Delete", "Will be deleted", "strength", false)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgID+"/training-guides/"+guideID, nil, token)
	assertStatus(t, resp, http.StatusOK)

	// Verify deleted
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/training-guides/"+guideID, nil, token)
	assertStatus(t, resp, http.StatusNotFound)
}

func TestGuide_AdminListAll(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800800001", "Admin")
	orgID := createTestOrg(t, adminID, "Guide Gym")
	createTestGuide(t, adminID, "Published", "Content", "strength", true)
	createTestGuide(t, adminID, "Draft", "Content", "cardio", false)
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/training-guides", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]interface{}
	parseJSON(t, resp, &body)
	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 guides (published + draft), got %d", len(data))
	}
}
