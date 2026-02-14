package e2e

import (
	"net/http"
	"testing"
)

func TestAuth_Login_MissingPhone(t *testing.T) {
	cleanupTables(t)

	resp := doRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{}, "")
	assertStatus(t, resp, http.StatusBadRequest)

	var body map[string]string
	parseJSON(t, resp, &body)
	if body["error"] != "phone is required" {
		t.Errorf("expected 'phone is required', got %q", body["error"])
	}
}

func TestAuth_Login_InvalidBody(t *testing.T) {
	resp, err := http.Post(testServer.URL+"/api/v1/auth/login", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestAuth_Login_ValidPhone(t *testing.T) {
	cleanupTables(t)

	resp := doRequest(t, http.MethodPost, "/api/v1/auth/login",
		map[string]string{"phone": "9801234567"}, "")
	assertStatus(t, resp, http.StatusOK)

	var body map[string]string
	parseJSON(t, resp, &body)
	if body["message"] != "OTP sent successfully" {
		t.Errorf("expected 'OTP sent successfully', got %q", body["message"])
	}
}

func TestAuth_VerifyOTP_MissingFields(t *testing.T) {
	cases := []struct {
		name string
		body map[string]string
	}{
		{"missing both", map[string]string{}},
		{"missing otp", map[string]string{"phone": "9801234567"}},
		{"missing phone", map[string]string{"otp": "123456"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, http.MethodPost, "/api/v1/auth/verify-otp", tc.body, "")
			assertStatus(t, resp, http.StatusBadRequest)
		})
	}
}

func TestAuth_Refresh_MissingToken(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/auth/refresh",
		map[string]string{}, "")
	assertStatus(t, resp, http.StatusBadRequest)

	var body map[string]string
	parseJSON(t, resp, &body)
	if body["error"] != "refresh_token is required" {
		t.Errorf("expected 'refresh_token is required', got %q", body["error"])
	}
}

func TestAuth_Refresh_InvalidToken(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/auth/refresh",
		map[string]string{"refresh_token": "invalid-token"}, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestAuth_Logout_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/auth/logout",
		map[string]string{"refresh_token": "xyz"}, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestAuth_LogoutAll_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/api/v1/auth/logout-all", nil, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestAuth_Logout_WithAuth(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9801234567", "Test User")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/auth/logout",
		map[string]string{}, token)
	assertStatus(t, resp, http.StatusBadRequest)

	var body map[string]string
	parseJSON(t, resp, &body)
	if body["error"] != "refresh_token is required" {
		t.Errorf("expected 'refresh_token is required', got %q", body["error"])
	}
}

func TestAuth_LogoutAll_WithAuth(t *testing.T) {
	cleanupTables(t)
	userID := createTestUser(t, "9801234567", "Test User")
	token := generateTestToken(userID, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/auth/logout-all", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body map[string]string
	parseJSON(t, resp, &body)
	if body["message"] != "logged out from all devices" {
		t.Errorf("expected 'logged out from all devices', got %q", body["message"])
	}
}

func TestAuth_ProtectedEndpoint_NoToken(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/members/me", nil, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestAuth_ProtectedEndpoint_InvalidToken(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/members/me", nil, "invalid-token")
	assertStatus(t, resp, http.StatusUnauthorized)
}

func TestAuth_ProtectedEndpoint_MalformedHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/api/v1/members/me", nil)
	req.Header.Set("Authorization", "NotBearer some-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, http.StatusUnauthorized)
}
