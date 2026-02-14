package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urja-gym/urja/internal/config"
	"github.com/urja-gym/urja/pkg/sms"
)

const testSecret = "test-secret-key-must-be-at-least-32-chars-long"

func testAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		JWTSecret:          testSecret,
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		OTPExpiry:          5 * time.Minute,
		OTPMaxAttempts:     3,
		OTPCooldown:        60 * time.Second,
		OTPHourlyLimit:     5,
		BcryptCost:         4, // Low cost for fast tests
	}
}

func TestGenerateOTP(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		otp, err := generateOTP()
		if err != nil {
			t.Fatalf("generateOTP() error: %v", err)
		}
		if len(otp) != 6 {
			t.Errorf("OTP length = %d, want 6", len(otp))
		}
		for _, c := range otp {
			if c < '0' || c > '9' {
				t.Errorf("OTP contains non-digit character: %c", c)
			}
		}
		// Verify range: 100000-999999
		if otp < "100000" || otp > "999999" {
			t.Errorf("OTP %s out of range [100000, 999999]", otp)
		}
		seen[otp] = true
	}
	// With 100 samples from a 900000-range space, we should have significant uniqueness
	if len(seen) < 90 {
		t.Errorf("poor OTP randomness: only %d unique values in 100 generations", len(seen))
	}
}

func TestGenerateTokenID(t *testing.T) {
	id1, err := generateTokenID()
	if err != nil {
		t.Fatalf("generateTokenID() error: %v", err)
	}
	if len(id1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("token ID length = %d, want 32", len(id1))
	}

	id2, err := generateTokenID()
	if err != nil {
		t.Fatalf("generateTokenID() second call error: %v", err)
	}
	if id1 == id2 {
		t.Error("two token IDs should not be identical")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	cfg := testAuthConfig()
	s := &Service{cfg: cfg}

	tokenStr, err := s.generateAccessToken("user-123", "9801234567", "admin")
	if err != nil {
		t.Fatalf("generateAccessToken() error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("generated token is empty")
	}

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("parsing generated token: %v", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		t.Fatal("token claims invalid")
	}

	if claims.Subject != "user-123" {
		t.Errorf("Subject = %q, want %q", claims.Subject, "user-123")
	}
	if claims.Phone != "9801234567" {
		t.Errorf("Phone = %q, want %q", claims.Phone, "9801234567")
	}
	if claims.Role != "admin" {
		t.Errorf("Role = %q, want %q", claims.Role, "admin")
	}
	if claims.Issuer != "urja" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "urja")
	}

	// Verify expiry is approximately cfg.AccessTokenExpiry from now
	expiry := claims.ExpiresAt.Time
	expectedExpiry := time.Now().Add(cfg.AccessTokenExpiry)
	if expiry.Before(expectedExpiry.Add(-time.Second)) || expiry.After(expectedExpiry.Add(time.Second)) {
		t.Errorf("expiry %v not within 1s of expected %v", expiry, expectedExpiry)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	cfg := testAuthConfig()
	s := &Service{cfg: cfg}

	tokenStr, tokenID, err := s.generateRefreshToken("user-456")
	if err != nil {
		t.Fatalf("generateRefreshToken() error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("generated refresh token is empty")
	}
	if tokenID == "" {
		t.Fatal("generated token ID is empty")
	}

	// Parse and verify
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("parsing generated refresh token: %v", err)
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		t.Fatal("refresh token claims invalid")
	}

	if claims.Subject != "user-456" {
		t.Errorf("Subject = %q, want %q", claims.Subject, "user-456")
	}
	if claims.TokenID != tokenID {
		t.Errorf("TokenID = %q, want %q", claims.TokenID, tokenID)
	}
	if claims.Issuer != "urja" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "urja")
	}

	// Verify expiry
	expiry := claims.ExpiresAt.Time
	expectedExpiry := time.Now().Add(cfg.RefreshTokenExpiry)
	if expiry.Before(expectedExpiry.Add(-time.Second)) || expiry.After(expectedExpiry.Add(time.Second)) {
		t.Errorf("expiry %v not within 1s of expected %v", expiry, expectedExpiry)
	}
}

func TestValidateAccessToken(t *testing.T) {
	cfg := testAuthConfig()
	s := &Service{cfg: cfg}

	// Generate a valid token
	tokenStr, err := s.generateAccessToken("user-789", "9801234567", "member")
	if err != nil {
		t.Fatalf("setup: generateAccessToken error: %v", err)
	}

	// Valid token should parse correctly
	userID, role, err := s.ValidateAccessToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error: %v", err)
	}
	if userID != "user-789" {
		t.Errorf("userID = %q, want %q", userID, "user-789")
	}
	if role != "member" {
		t.Errorf("role = %q, want %q", role, "member")
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	cfg := testAuthConfig()
	s := &Service{cfg: cfg}

	_, _, err := s.ValidateAccessToken("this.is.garbage")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	cfg := testAuthConfig()
	s := &Service{cfg: cfg}

	// Generate token with a different secret
	otherService := &Service{
		cfg: config.AuthConfig{
			JWTSecret:         "another-secret-key-that-is-also-very-long",
			AccessTokenExpiry: 15 * time.Minute,
		},
	}
	tokenStr, err := otherService.generateAccessToken("user-123", "9801234567", "member")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	_, _, err = s.ValidateAccessToken(tokenStr)
	if err == nil {
		t.Error("expected error for token signed with wrong secret, got nil")
	}
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	cfg := testAuthConfig()
	cfg.AccessTokenExpiry = -1 * time.Hour // Already expired
	s := &Service{cfg: cfg}

	tokenStr, err := s.generateAccessToken("user-123", "9801234567", "member")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	// Reset config for validation (validation uses current time)
	s.cfg.AccessTokenExpiry = 15 * time.Minute
	_, _, err = s.ValidateAccessToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestNormalizePhoneIntegration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"9801234567", "9801234567", true},
		{"+9779801234567", "9801234567", true},
		{"9779801234567", "9801234567", true},
		{" 9801234567 ", "9801234567", true},
		{"+977 9801234567", "9801234567", true}, // space between prefix and number
		{"1234567890", "1234567890", false},       // doesn't start with 9[6-9]
		{"980123456", "980123456", false},          // too short
	}
	for _, tt := range tests {
		normalized := sms.NormalizePhone(tt.input)
		if normalized != tt.expected {
			t.Errorf("NormalizePhone(%q) = %q, want %q", tt.input, normalized, tt.expected)
		}
		if sms.ValidatePhone(normalized) != tt.valid {
			t.Errorf("ValidatePhone(%q) = %v, want %v", normalized, !tt.valid, tt.valid)
		}
	}
}
