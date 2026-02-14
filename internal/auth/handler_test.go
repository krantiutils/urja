package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler_MissingPhone(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp.Error != "phone is required" {
		t.Errorf("error = %q, want %q", resp.Error, "phone is required")
	}
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestVerifyOTPHandler_MissingFields(t *testing.T) {
	h := &Handler{logger: testLogger()}

	tests := []struct {
		name string
		body string
	}{
		{"missing both", `{}`},
		{"missing otp", `{"phone":"9801234567"}`},
		{"missing phone", `{"otp":"123456"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/verify-otp", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.VerifyOTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestRefreshHandler_MissingToken(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp.Error != "refresh_token is required" {
		t.Errorf("error = %q, want %q", resp.Error, "refresh_token is required")
	}
}

func TestLogoutHandler_NoAuth(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(`{"refresh_token":"xyz"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogoutAllHandler_NoAuth(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/auth/logout-all", nil)
	w := httptest.NewRecorder()

	h.LogoutAll(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
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
