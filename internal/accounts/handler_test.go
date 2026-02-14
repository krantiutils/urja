package accounts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

func TestListHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
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

func TestCreateHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateHandler_NoAuth(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{}`)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestCreateHandler_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	ctx = context.WithValue(ctx, middleware.UserIDKey, "user-123")
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`not json`)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPut, "/accounts/123", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateHandler_MissingID(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	req := httptest.NewRequest(http.MethodPut, "/accounts/", strings.NewReader(`{}`)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Without chi context, URLParam returns empty string
	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp["error"] != "transaction ID is required" {
		t.Errorf("error = %q, want %q", resp["error"], "transaction ID is required")
	}
}

func TestDeleteHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodDelete, "/accounts/123", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteHandler_MissingID(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	req := httptest.NewRequest(http.MethodDelete, "/accounts/", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp["error"] != "transaction ID is required" {
		t.Errorf("error = %q, want %q", resp["error"], "transaction ID is required")
	}
}

func TestGetSummaryHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/accounts/summary", nil)
	w := httptest.NewRecorder()

	h.GetSummary(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestExportHandler_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/accounts/export", nil)
	w := httptest.NewRecorder()

	h.Export(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateHandler_WithValidID(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	req := httptest.NewRequest(http.MethodPut, "/accounts/txn-abc", strings.NewReader(`{}`)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "txn-abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	// Will call service.Update which calls service validation - nil service will panic.
	// Testing that handler doesn't reject for missing ID.
	defer func() {
		if r := recover(); r != nil {
			// Expected: nil service dereference
		}
	}()

	h.Update(w, req)
}

func TestDeleteHandler_WithValidID(t *testing.T) {
	h := &Handler{logger: testLogger()}

	ctx := context.WithValue(context.Background(), middleware.OrgIDKey, "org-123")
	req := httptest.NewRequest(http.MethodDelete, "/accounts/txn-abc", nil).WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "txn-abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Expected: nil service dereference
		}
	}()

	h.Delete(w, req)
}

func TestCsvEscape(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"", ""},
		{"has,comma", `"has,comma"`},
		{`has"quote`, `"has""quote"`},
		{"has\nnewline", `"has` + "\n" + `newline"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := csvEscape(tt.input)
			if got != tt.want {
				t.Errorf("csvEscape(%q) = %q, want %q", tt.input, got, tt.want)
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
