package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	// UserIDKey is the context key for the authenticated user's ID.
	UserIDKey contextKey = "user_id"
	// UserRoleKey is the context key for the authenticated user's role.
	UserRoleKey contextKey = "user_role"
	// OrgIDKey is the context key for the organization ID from the route.
	OrgIDKey contextKey = "org_id"
	// OrgRoleKey is the context key for the user's role within the current organization.
	OrgRoleKey contextKey = "org_role"
)

// TokenValidator is the interface that the auth service must implement
// for the middleware to validate JWT tokens.
type TokenValidator interface {
	ValidateAccessToken(tokenString string) (userID string, role string, err error)
}

// Auth returns middleware that validates the Authorization header
// and injects user identity into the request context.
func Auth(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			userID, role, err := validator.ValidateAccessToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the user ID from the request context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

// UserRoleFromContext extracts the user role from the request context.
func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

// OrgIDFromContext extracts the organization ID from the request context.
func OrgIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(OrgIDKey).(string)
	return id, ok
}

// OrgRoleFromContext extracts the user's org-specific role from the request context.
func OrgRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(OrgRoleKey).(string)
	return role, ok
}
