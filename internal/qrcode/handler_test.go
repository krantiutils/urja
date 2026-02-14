package qrcode

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestGenerateQR_PNG(t *testing.T) {
	h := NewHandler("https://urja.app", testLogger())

	req := httptest.NewRequest(http.MethodGet, "/qr-code", nil)
	// Inject org ID into context (simulates OrgScope middleware)
	ctx := context.WithValue(req.Context(), middleware.OrgIDKey, "test-org-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.GenerateQR(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("Content-Type = %q, want %q", ct, "image/png")
	}
	// PNG files start with the PNG magic bytes
	body := rr.Body.Bytes()
	if len(body) < 8 {
		t.Fatal("response body too short to be a valid PNG")
	}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngMagic {
		if body[i] != b {
			t.Errorf("PNG magic byte[%d] = 0x%02X, want 0x%02X", i, body[i], b)
		}
	}
}

func TestGenerateQR_JSON(t *testing.T) {
	h := NewHandler("https://urja.app", testLogger())

	req := httptest.NewRequest(http.MethodGet, "/qr-code?format=json", nil)
	ctx := context.WithValue(req.Context(), middleware.OrgIDKey, "test-org-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.GenerateQR(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["org_id"] != "test-org-id" {
		t.Errorf("org_id = %q, want %q", resp["org_id"], "test-org-id")
	}
	if resp["checkin_url"] != "https://urja.app/checkin/test-org-id" {
		t.Errorf("checkin_url = %q, want %q", resp["checkin_url"], "https://urja.app/checkin/test-org-id")
	}
}

func TestGenerateQR_MissingOrgID(t *testing.T) {
	h := NewHandler("https://urja.app", testLogger())

	req := httptest.NewRequest(http.MethodGet, "/qr-code", nil)
	// No org ID in context
	rr := httptest.NewRecorder()
	h.GenerateQR(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRegisterOrgRoutes(t *testing.T) {
	h := NewHandler("https://urja.app", testLogger())
	r := chi.NewRouter()
	h.RegisterOrgRoutes(r)

	// Walk the routes and verify /qr-code GET exists
	found := false
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if method == "GET" && route == "/qr-code" {
			found = true
		}
		return nil
	})
	if !found {
		t.Error("expected GET /qr-code route to be registered")
	}
}
