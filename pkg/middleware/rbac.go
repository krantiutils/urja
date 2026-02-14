package middleware

import (
	"net/http"
)

// RequireRole returns middleware that enforces role-based access control.
// Only requests from users with one of the allowed roles will pass through.
func RequireRole(allowed ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := UserRoleFromContext(r.Context())
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				http.Error(w, `{"error":"forbidden: insufficient role"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
