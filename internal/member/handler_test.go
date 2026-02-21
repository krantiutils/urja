package member

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urja-gym/urja/pkg/middleware"
)

// setAuthContext injects user ID into request context (simulates auth middleware).
func setAuthContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	return r.WithContext(ctx)
}

// setOrgContext injects org ID into request context (simulates OrgScope middleware).
func setOrgContext(r *http.Request, orgID string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.OrgIDKey, orgID)
	return r.WithContext(ctx)
}

func TestGetMe_Unauthorized(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members/me", nil)
	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["error"] != "unauthorized" {
		t.Fatalf("expected error 'unauthorized', got %q", resp["error"])
	}
}

func TestUpdateMe_Unauthorized(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	body := `{"name": "Ram"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/me", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.UpdateMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestUpdateMe_InvalidBody(t *testing.T) {
	h := NewHandler(&Service{}, nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/me", bytes.NewBufferString("not json"))
	req = setAuthContext(req, "user-123")
	rec := httptest.NewRecorder()

	h.UpdateMe(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdatePrivacy_Unauthorized(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	body := `{"show_email": true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/me/privacy", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.UpdatePrivacy(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestUpdatePrivacy_InvalidBody(t *testing.T) {
	h := NewHandler(&Service{}, nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/me/privacy", bytes.NewBufferString("{bad"))
	req = setAuthContext(req, "user-123")
	rec := httptest.NewRecorder()

	h.UpdatePrivacy(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestListByOrg_MissingOrg(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/xxx/members", nil)
	rec := httptest.NewRecorder()

	h.ListByOrg(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateOrgMember_MissingOrg(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	body := `{"phone": "9841234567", "name": "Ram"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/xxx/members", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.CreateOrgMember(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateOrgMember_InvalidBody(t *testing.T) {
	h := NewHandler(&Service{}, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/xxx/members", bytes.NewBufferString("not json"))
	req = setOrgContext(req, "org-123")
	rec := httptest.NewRecorder()

	h.CreateOrgMember(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteOrgMember_MissingOrg(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orgs/xxx/members/yyy", nil)
	rec := httptest.NewRecorder()

	h.DeleteOrgMember(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusOK, map[string]string{"key": "value"})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type 'application/json', got %q", contentType)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if resp["key"] != "value" {
		t.Fatalf("expected key=value, got key=%q", resp["key"])
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "plain error is validation",
			err:    errNew("name cannot be empty"),
			expect: true,
		},
		{
			name:   "wrapped error is not validation",
			err:    errWrap("updating profile", errNew("db error")),
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidationError(tt.err)
			if got != tt.expect {
				t.Fatalf("isValidationError(%v) = %v, want %v", tt.err, got, tt.expect)
			}
		})
	}
}

type plainError string

func (e plainError) Error() string { return string(e) }

func errNew(msg string) error {
	return plainError(msg)
}

type wrappedError struct {
	msg   string
	cause error
}

func (e *wrappedError) Error() string { return e.msg + ": " + e.cause.Error() }
func (e *wrappedError) Unwrap() error { return e.cause }

func errWrap(msg string, cause error) error {
	return &wrappedError{msg: msg, cause: cause}
}
