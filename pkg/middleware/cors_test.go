package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testCORS() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "https://nepalgym.xyz"},
		TenantBaseDomain: "nepalgym.xyz",
	}
}

func TestCORS_AllowsExactOrigins(t *testing.T) {
	cfg := testCORS()
	for _, origin := range []string{"http://localhost:3000", "https://nepalgym.xyz"} {
		if !cfg.allows(origin) {
			t.Errorf("expected %q to be allowed", origin)
		}
	}
}

func TestCORS_AllowsTenantSubdomains(t *testing.T) {
	cfg := testCORS()
	for _, origin := range []string{
		"https://ibckirtipur.nepalgym.xyz",
		"https://pimbahal-gym.nepalgym.xyz",
		"https://gym123.nepalgym.xyz",
	} {
		if !cfg.allows(origin) {
			t.Errorf("expected tenant origin %q to be allowed", origin)
		}
	}
}

// The reflected origin is paired with Allow-Credentials, so a near-miss that
// slips through hands an attacker's page authenticated access.
func TestCORS_RejectsLookalikeOrigins(t *testing.T) {
	cfg := testCORS()
	for _, origin := range []string{
		"https://evil-nepalgym.xyz",           // no separating dot
		"https://nepalgym.xyz.attacker.com",   // base domain as a prefix
		"https://a.b.nepalgym.xyz",            // two labels deep
		"http://ibckirtipur.nepalgym.xyz",     // plaintext
		"https://ibckirtipur.nepalgym.xyz:81", // port
		"https://IBC.nepalgym.xyz",            // uppercase is not a valid slug
		"https://-gym.nepalgym.xyz",           // leading hyphen
		"https://gym-.nepalgym.xyz",           // trailing hyphen
		"https://.nepalgym.xyz",               // empty label
		"https://nepalgym.xyz.evil.com",
		"null",
		"",
	} {
		if cfg.allows(origin) {
			t.Errorf("expected %q to be rejected", origin)
		}
	}
}

func TestCORS_WildcardOffWithoutBaseDomain(t *testing.T) {
	cfg := CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}}
	if cfg.allows("https://ibckirtipur.nepalgym.xyz") {
		t.Error("tenant origins must not be allowed when no base domain is set")
	}
}

func TestCORS_ReflectsAllowedOriginAndVaries(t *testing.T) {
	h := CORS(testCORS())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sites/x", nil)
	req.Header.Set("Origin", "https://ibckirtipur.nepalgym.xyz")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://ibckirtipur.nepalgym.xyz" {
		t.Errorf("Allow-Origin = %q", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Errorf("Vary = %q, want Origin", got)
	}
}

func TestCORS_DisallowedOriginGetsNoAllowHeaders(t *testing.T) {
	h := CORS(testCORS())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sites/x", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin should be empty, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Errorf("Allow-Credentials should not be set for a rejected origin, got %q", got)
	}
	// The request itself still runs; the browser is what enforces the policy.
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestCORS_PreflightShortCircuits(t *testing.T) {
	called := false
	h := CORS(testCORS())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/sites/x/leads", nil)
	req.Header.Set("Origin", "https://ibckirtipur.nepalgym.xyz")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("preflight must not reach the handler")
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("preflight must advertise allowed methods")
	}
}
