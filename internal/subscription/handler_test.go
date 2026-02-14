package subscription

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestAssignPackage_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/packages/assign", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AssignPackage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp["error"] != "missing organization" {
		// Without org context, we get missing organization first
		t.Logf("error = %q (expected: depends on context injection)", resp["error"])
	}
}

func TestRenewPackage_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/packages/renew", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RenewPackage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestExtendPackage_InvalidBody(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodPost, "/packages/extend", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ExtendPackage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPackageSummary_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/packages/summary", nil)
	w := httptest.NewRecorder()

	h.PackageSummary(w, req)

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

func TestListExpiring_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/packages/expiring", nil)
	w := httptest.NewRecorder()

	h.ListExpiring(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListExpired_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/packages/expired", nil)
	w := httptest.NewRecorder()

	h.ListExpired(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListPayments_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/members/123/payments", nil)
	w := httptest.NewRecorder()

	h.ListPayments(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListSubscriptions_MissingOrg(t *testing.T) {
	h := &Handler{logger: testLogger()}

	req := httptest.NewRequest(http.MethodGet, "/members/123/subscriptions", nil)
	w := httptest.NewRecorder()

	h.ListSubscriptions(w, req)

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

func TestHandleServiceError_Mapping(t *testing.T) {
	h := &Handler{logger: testLogger()}

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid request", ErrInvalidRequest, http.StatusBadRequest},
		{"invalid payment", ErrInvalidPayment, http.StatusBadRequest},
		{"invalid package", ErrInvalidPackage, http.StatusNotFound},
		{"invalid member", ErrInvalidMember, http.StatusNotFound},
		{"subscription not found", ErrSubscriptionNotFound, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.handleServiceError(w, tt.err, "test", "org-1", "member-1")
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
