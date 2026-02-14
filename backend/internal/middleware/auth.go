package middleware

import (
	"context"
	"net/http"
	"strings"

	ujwt "github.com/urja-gym/backend/pkg/jwt"
)

type authContextKey string

const userClaimsKey authContextKey = "user_claims"

// Auth validates the access token in the Authorization header and injects
// claims into the request context.
func Auth(jwtIssuer *ujwt.Issuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, `{"error":"unauthorized","message":"Missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"unauthorized","message":"Invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtIssuer.ValidateAccessToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims extracts the user claims from the context.
func GetUserClaims(ctx context.Context) *ujwt.AccessClaims {
	if claims, ok := ctx.Value(userClaimsKey).(*ujwt.AccessClaims); ok {
		return claims
	}
	return nil
}
