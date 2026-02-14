package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireRole(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name       string
		allowed    []string
		role       string
		hasRole    bool
		wantStatus int
	}{
		{"allowed role", []string{"admin", "staff"}, "admin", true, http.StatusOK},
		{"disallowed role", []string{"admin"}, "member", true, http.StatusForbidden},
		{"no role in context", []string{"admin"}, "", false, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.hasRole {
				ctx := context.WithValue(req.Context(), UserRoleKey, tt.role)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			RequireRole(tt.allowed...)(okHandler).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestRequireOrgRole(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name       string
		allowed    []string
		role       string
		hasRole    bool
		wantStatus int
	}{
		{"allowed staff", []string{"staff", "admin"}, "staff", true, http.StatusOK},
		{"allowed admin", []string{"staff", "admin"}, "admin", true, http.StatusOK},
		{"disallowed member", []string{"staff", "admin"}, "member", true, http.StatusForbidden},
		{"no org role in context", []string{"staff"}, "", false, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.hasRole {
				ctx := context.WithValue(req.Context(), OrgRoleKey, tt.role)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			RequireOrgRole(tt.allowed...)(okHandler).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
