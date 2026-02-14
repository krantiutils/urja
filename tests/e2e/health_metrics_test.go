package e2e

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestHealthMetrics_LogBMI(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9800100100", "BMI User")
	token := generateTestToken(userID, "user")

	t.Run("success", func(t *testing.T) {
		body := map[string]interface{}{
			"height_cm": 175.0,
			"weight_kg": 70.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		resp = doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, token)
		parseJSON(t, resp, &result)

		if result["metric_type"] != "bmi" {
			t.Fatalf("expected metric_type 'bmi', got %v", result["metric_type"])
		}

		value, ok := result["value"].(map[string]interface{})
		if !ok {
			t.Fatal("expected value to be a JSON object")
		}
		bmi, ok := value["bmi"].(float64)
		if !ok || bmi < 22 || bmi > 23 {
			t.Fatalf("expected BMI ~22.86, got %v", bmi)
		}
		if value["category"] != "normal" {
			t.Fatalf("expected category 'normal', got %v", value["category"])
		}
	})

	t.Run("underweight_bmi", func(t *testing.T) {
		body := map[string]interface{}{
			"height_cm": 180.0,
			"weight_kg": 55.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		value := result["value"].(map[string]interface{})
		if value["category"] != "underweight" {
			t.Fatalf("expected category 'underweight', got %v", value["category"])
		}
	})

	t.Run("obese_bmi", func(t *testing.T) {
		body := map[string]interface{}{
			"height_cm": 170.0,
			"weight_kg": 100.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		value := result["value"].(map[string]interface{})
		if value["category"] != "obese" {
			t.Fatalf("expected category 'obese', got %v", value["category"])
		}
	})

	t.Run("invalid_height", func(t *testing.T) {
		body := map[string]interface{}{
			"height_cm": -5.0,
			"weight_kg": 70.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("unauthorized", func(t *testing.T) {
		body := map[string]interface{}{
			"height_cm": 175.0,
			"weight_kg": 70.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi", body, "")
		assertStatus(t, resp, http.StatusUnauthorized)
	})
}

func TestHealthMetrics_LogMeasurements(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9800200200", "Measurement User")
	token := generateTestToken(userID, "user")

	t.Run("success_with_whr", func(t *testing.T) {
		body := map[string]interface{}{
			"chest_cm":  100.0,
			"waist_cm":  80.0,
			"hips_cm":   95.0,
			"arms_cm":   35.0,
			"thighs_cm": 55.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/measurements", body, token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		parseJSON(t, resp, &result)

		if result["metric_type"] != "measurements" {
			t.Fatalf("expected metric_type 'measurements', got %v", result["metric_type"])
		}

		value := result["value"].(map[string]interface{})
		whr, ok := value["whr"].(float64)
		if !ok || whr < 0.84 || whr > 0.85 {
			t.Fatalf("expected WHR ~0.84, got %v", whr)
		}
	})

	t.Run("partial_measurements_no_whr", func(t *testing.T) {
		body := map[string]interface{}{
			"chest_cm": 100.0,
			"arms_cm":  35.0,
		}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/measurements", body, token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		value := result["value"].(map[string]interface{})
		if _, ok := value["whr"]; ok {
			t.Fatal("expected no WHR when waist/hips not provided")
		}
	})

	t.Run("empty_measurements", func(t *testing.T) {
		body := map[string]interface{}{}
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/measurements", body, token)
		assertStatus(t, resp, http.StatusBadRequest)
	})
}

func TestHealthMetrics_GetMetrics(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9800300300", "Metrics User")
	token := generateTestToken(userID, "user")

	// Log some metrics
	doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi",
		map[string]interface{}{"height_cm": 175.0, "weight_kg": 70.0}, token)
	doRequest(t, http.MethodPost, "/api/v1/members/me/health/bmi",
		map[string]interface{}{"height_cm": 175.0, "weight_kg": 72.0}, token)
	doRequest(t, http.MethodPost, "/api/v1/members/me/health/measurements",
		map[string]interface{}{"waist_cm": 80.0, "hips_cm": 95.0}, token)

	t.Run("get_all_metrics", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		if len(data) != 3 {
			t.Fatalf("expected 3 metrics, got %d", len(data))
		}
	})

	t.Run("filter_by_type", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health?type=bmi", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		if len(data) != 2 {
			t.Fatalf("expected 2 BMI metrics, got %d", len(data))
		}
	})

	t.Run("filter_by_type_measurements", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health?type=measurements", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		if len(data) != 1 {
			t.Fatalf("expected 1 measurements metric, got %d", len(data))
		}
	})

	t.Run("limit_results", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health?limit=1", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		if len(data) != 1 {
			t.Fatalf("expected 1 metric with limit=1, got %d", len(data))
		}
	})

	t.Run("invalid_type", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health?type=invalid", nil, token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("empty_results_for_other_user", func(t *testing.T) {
		otherID := createTestUser(t, "9800300301", "Other User")
		otherToken := generateTestToken(otherID, "user")
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health", nil, otherToken)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		if len(data) != 0 {
			t.Fatalf("expected 0 metrics for other user, got %d", len(data))
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health", nil, "")
		assertStatus(t, resp, http.StatusUnauthorized)
	})
}

func TestProgressPhotos_Upload(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9800400400", "Photo User")
	token := generateTestToken(userID, "user")

	t.Run("success", func(t *testing.T) {
		resp := doMultipartUpload(t, "/api/v1/members/me/health/photos",
			"photo", "test.jpg", []byte("fake-jpeg-content"),
			map[string]string{"category": "front", "notes": "Week 1 progress"},
			token)
		assertStatus(t, resp, http.StatusCreated)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		if result["category"] != "front" {
			t.Fatalf("expected category 'front', got %v", result["category"])
		}
		if result["notes"] != "Week 1 progress" {
			t.Fatalf("expected notes 'Week 1 progress', got %v", result["notes"])
		}
		if result["photo_url"] == nil || result["photo_url"] == "" {
			t.Fatal("expected non-empty photo_url")
		}
	})

	t.Run("no_category", func(t *testing.T) {
		resp := doMultipartUpload(t, "/api/v1/members/me/health/photos",
			"photo", "test2.png", []byte("fake-png-content"),
			map[string]string{},
			token)
		assertStatus(t, resp, http.StatusCreated)
	})

	t.Run("invalid_category", func(t *testing.T) {
		resp := doMultipartUpload(t, "/api/v1/members/me/health/photos",
			"photo", "test3.jpg", []byte("fake-content"),
			map[string]string{"category": "invalid"},
			token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("invalid_file_type", func(t *testing.T) {
		resp := doMultipartUpload(t, "/api/v1/members/me/health/photos",
			"photo", "doc.pdf", []byte("fake-content"),
			map[string]string{},
			token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("missing_file", func(t *testing.T) {
		resp := doRequest(t, http.MethodPost, "/api/v1/members/me/health/photos", nil, token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := doMultipartUpload(t, "/api/v1/members/me/health/photos",
			"photo", "test.jpg", []byte("content"),
			map[string]string{},
			"")
		assertStatus(t, resp, http.StatusUnauthorized)
	})
}

func TestProgressPhotos_List(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9800500500", "Photo List User")
	token := generateTestToken(userID, "user")

	// Upload some photos
	doMultipartUpload(t, "/api/v1/members/me/health/photos",
		"photo", "front1.jpg", []byte("f1"),
		map[string]string{"category": "front"}, token)
	doMultipartUpload(t, "/api/v1/members/me/health/photos",
		"photo", "side1.jpg", []byte("s1"),
		map[string]string{"category": "side"}, token)
	doMultipartUpload(t, "/api/v1/members/me/health/photos",
		"photo", "back1.jpg", []byte("b1"),
		map[string]string{"category": "back"}, token)

	t.Run("list_all", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health/photos", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		total := result["total"].(float64)
		if len(data) != 3 || int(total) != 3 {
			t.Fatalf("expected 3 photos, got %d (total: %v)", len(data), total)
		}
	})

	t.Run("filter_by_category", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health/photos?category=front", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		total := result["total"].(float64)
		if len(data) != 1 || int(total) != 1 {
			t.Fatalf("expected 1 front photo, got %d (total: %v)", len(data), total)
		}
	})

	t.Run("invalid_category", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health/photos?category=invalid", nil, token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("pagination", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/api/v1/members/me/health/photos?limit=2&offset=0", nil, token)
		assertStatus(t, resp, http.StatusOK)

		var result map[string]interface{}
		parseJSON(t, resp, &result)
		data := result["data"].([]interface{})
		total := result["total"].(float64)
		if len(data) != 2 {
			t.Fatalf("expected 2 photos with limit=2, got %d", len(data))
		}
		if int(total) != 3 {
			t.Fatalf("expected total=3, got %v", total)
		}
	})
}

// doMultipartUpload sends a multipart form upload request to the test server.
func doMultipartUpload(t *testing.T, path, fieldName, fileName string, fileContent []byte, fields map[string]string, token string) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("creating form file: %v", err)
	}
	if _, err := part.Write(fileContent); err != nil {
		t.Fatalf("writing file content: %v", err)
	}

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("writing field %s: %v", k, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("closing multipart writer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, testServer.URL+path, &buf)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}
	return resp
}

// TestHealthCheck already exists in health_test.go, but we add the
// unused import guard here for completeness. This file is health_metrics_test.go.
var _ = io.Discard
