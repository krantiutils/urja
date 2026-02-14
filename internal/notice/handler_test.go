package notice

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateHandler_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/notices", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	// No org context -> bad request before body parsing
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/notices", strings.NewReader(`{"title":"test","content":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp["error"] != "missing organization" {
		t.Errorf("error = %q, want %q", resp["error"], "missing organization")
	}
}

func TestListHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/notices", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/notices/some-id", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodDelete, "/notices/some-id", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusCreated, map[string]string{"key": "value"})

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var data map[string]string
	if err := json.NewDecoder(w.Body).Decode(&data); err != nil {
		t.Fatalf("decoding: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("key = %q, want %q", data["key"], "value")
	}
}
