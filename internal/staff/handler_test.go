package staff

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type errorResponse struct {
	Error string `json:"error"`
}

func TestCreateHandler_InvalidBody(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodPost, "/staff", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateHandler_MissingFields(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	tests := []struct {
		name string
		body string
	}{
		{"missing all", `{}`},
		{"missing name", `{"phone":"9801234567","staff_role":"trainer"}`},
		{"missing phone", `{"name":"John","staff_role":"trainer"}`},
		{"missing staff_role", `{"phone":"9801234567","name":"John"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/staff", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Create(w, req)

			// Without org context, should be 400
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestCreateHandler_NoOrgContext(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	body := `{"phone":"9801234567","name":"John","staff_role":"trainer"}`
	req := httptest.NewRequest(http.MethodPost, "/staff", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp.Error != "missing organization" {
		t.Errorf("error = %q, want %q", resp.Error, "missing organization")
	}
}

func TestListHandler_NoOrgContext(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodGet, "/staff", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetHandler_NoOrgContext(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodGet, "/staff/some-id", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateHandler_InvalidBody(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodPut, "/staff/some-id", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Update(w, req)

	// Without org context, should be 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteHandler_NoOrgContext(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodDelete, "/staff/some-id", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListAttendanceHandler_NoOrgContext(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodGet, "/staff/attendance", nil)
	w := httptest.NewRecorder()

	h.ListAttendance(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRecordAttendanceHandler_InvalidBody(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	req := httptest.NewRequest(http.MethodPost, "/staff/attendance", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RecordAttendance(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRecordAttendanceHandler_MissingFields(t *testing.T) {
	svc := NewService(nil, testLogger())
	h := NewHandler(svc, testLogger())

	tests := []struct {
		name string
		body string
	}{
		{"missing all", `{}`},
		{"missing action", `{"staff_id":"some-id"}`},
		{"missing staff_id", `{"action":"check_in","method":"manual"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/staff/attendance", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.RecordAttendance(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
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
