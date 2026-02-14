package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// OrgMembershipChecker verifies that a user belongs to an organization.
type OrgMembershipChecker interface {
	IsOrgMember(ctx context.Context, userID, orgID string) (bool, error)
}

// OrgScope returns middleware that extracts the organization ID from the URL
// and verifies the authenticated user is a member of that organization.
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

			isMember, err := checker.IsOrgMember(r.Context(), userID, orgID)
			if err != nil {
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
			if !isMember {
				http.Error(w, `{"error":"forbidden: not a member of this organization"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
