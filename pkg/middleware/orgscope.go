package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// OrgMembershipChecker verifies that a user belongs to an organization
// and returns their role within it.
type OrgMembershipChecker interface {
	CheckOrgMembership(ctx context.Context, userID, orgID string) (isMember bool, role string, err error)
}

// OrgScope returns middleware that extracts the organization ID from the URL
// and verifies the authenticated user is a member of that organization.
// It also sets the user's org-specific role in the request context.
// This prevents cross-org data access (PRD vulnerability #8).
func OrgScope(checker OrgMembershipChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgID := chi.URLParam(r, "orgId")
			if orgID == "" {
				http.Error(w, `{"error":"missing organization ID"}`, http.StatusBadRequest)
				return
			}

			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			isMember, role, err := checker.CheckOrgMembership(r.Context(), userID, orgID)
			if err != nil {
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
			if !isMember {
				http.Error(w, `{"error":"forbidden: not a member of this organization"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
			ctx = context.WithValue(ctx, OrgRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
