package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig describes which browser origins may call the API.
type CORSConfig struct {
	// AllowedOrigins is the exact-match list, for the dashboard and local dev.
	AllowedOrigins []string
	// TenantBaseDomain enables every gym's own subdomain, e.g. "nepalgym.xyz"
	// allows https://ibckirtipur.nepalgym.xyz. Empty disables wildcard matching.
	//
	// Tenant sites cannot be enumerated in a static list: each gym gets its own
	// origin the moment it is created, and the lead form on that site posts to
	// this API. Without this, publishing a new gym silently breaks its only
	// conversion path.
	TenantBaseDomain string
}

// allowedLabelChars mirrors the slug rule enforced when an organization is
// created: lowercase alphanumerics and hyphens, nothing else.
func isValidSubdomainLabel(label string) bool {
	if label == "" || len(label) > 63 {
		return false
	}
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}
	for _, c := range label {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
			return false
		}
	}
	return true
}

// allows reports whether an Origin header may be reflected.
//
// Matching is deliberately strict rather than a suffix test: "evil-nepalgym.xyz"
// and "nepalgym.xyz.attacker.com" both end well for a naive strings.HasSuffix,
// and the response carries Allow-Credentials, so a loose match would hand an
// attacker's page authenticated access.
func (c CORSConfig) allows(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range c.AllowedOrigins {
		if o == origin {
			return true
		}
	}
	if c.TenantBaseDomain == "" {
		return false
	}

	// Only https tenant sites; a plaintext origin is never a published gym.
	host, ok := strings.CutPrefix(origin, "https://")
	if !ok {
		return false
	}
	// A port or path means this is not a bare tenant host.
	if strings.ContainsAny(host, ":/") {
		return false
	}

	label, ok := strings.CutSuffix(host, "."+c.TenantBaseDomain)
	if !ok {
		return false
	}
	// Exactly one label deep: a.b.nepalgym.xyz is not a tenant.
	if strings.Contains(label, ".") {
		return false
	}
	return isValidSubdomainLabel(label)
}

// CORS reflects allowed origins and answers preflight requests.
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// The response body is origin-independent but these headers are not,
			// so any shared cache must key on Origin.
			w.Header().Add("Vary", "Origin")

			if cfg.allows(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				// Only meaningful alongside Allow-Origin; setting it for a
				// rejected origin just muddies what the policy actually is.
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
