package absentee

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/absentees", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

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

func TestNotifyHandler_NoOrgContext(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/absentees/member-1/notify",
		strings.NewReader(`{"org_name":"Gym"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Notify(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestNotifyHandler_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/absentees/member-1/notify",
		strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Notify(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var data map[string]string
	if err := json.NewDecoder(w.Body).Decode(&data); err != nil {
		t.Fatalf("decoding: %v", err)
	}
	if data["message"] != "ok" {
		t.Errorf("message = %q, want %q", data["message"], "ok")
	}
}
