package jwt

import (
	"testing"
	"time"
)

const (
	testAccessSecret  = "test-access-secret-that-is-long-enough-for-validation"
	testRefreshSecret = "test-refresh-secret-that-is-long-enough-for-validation"
)

func TestIssueAndValidateAccessToken(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	pair, err := issuer.IssueTokenPair("user-123", "9841234567", "member", "org-456")
	if err != nil {
		t.Fatalf("IssueTokenPair() error = %v", err)
	}

	if pair.AccessToken == "" {
		t.Fatal("access token is empty")
	}
	if pair.RefreshToken == "" {
		t.Fatal("refresh token is empty")
	}
	if pair.ExpiresIn != int64(AccessTokenExpiry.Seconds()) {
		t.Errorf("ExpiresIn = %d, want %d", pair.ExpiresIn, int64(AccessTokenExpiry.Seconds()))
	}

	// Validate access token
	claims, err := issuer.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.Subject != "user-123" {
		t.Errorf("Subject = %q, want %q", claims.Subject, "user-123")
	}
	if claims.Phone != "9841234567" {
		t.Errorf("Phone = %q, want %q", claims.Phone, "9841234567")
	}
	if claims.Role != "member" {
		t.Errorf("Role = %q, want %q", claims.Role, "member")
	}
	if claims.OrgID != "org-456" {
		t.Errorf("OrgID = %q, want %q", claims.OrgID, "org-456")
	}
	if claims.ID == "" {
		t.Error("JTI is empty")
	}
	if claims.Issuer != "urja" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "urja")
	}
}

func TestValidateAccessTokenWithRefreshSecret(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	pair, err := issuer.IssueTokenPair("user-123", "9841234567", "member", "")
	if err != nil {
		t.Fatalf("IssueTokenPair() error = %v", err)
	}

	// Trying to validate a refresh token as an access token should fail
	_, err = issuer.ValidateAccessToken(pair.RefreshToken)
	if err == nil {
		t.Error("expected error validating refresh token as access token")
	}
}

func TestValidateRefreshToken(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	pair, err := issuer.IssueTokenPair("user-123", "9841234567", "member", "")
	if err != nil {
		t.Fatalf("IssueTokenPair() error = %v", err)
	}

	claims, err := issuer.ValidateRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if claims.Subject != "user-123" {
		t.Errorf("Subject = %q, want %q", claims.Subject, "user-123")
	}
	if claims.ID == "" {
		t.Error("JTI is empty")
	}

	// Access token should not validate as refresh token
	_, err = issuer.ValidateRefreshToken(pair.AccessToken)
	if err == nil {
		t.Error("expected error validating access token as refresh token")
	}
}

func TestTokenExpiry(t *testing.T) {
	if AccessTokenExpiry != 15*time.Minute {
		t.Errorf("AccessTokenExpiry = %v, want 15m", AccessTokenExpiry)
	}
	if RefreshTokenExpiry != 7*24*time.Hour {
		t.Errorf("RefreshTokenExpiry = %v, want 7d", RefreshTokenExpiry)
	}
}

func TestEmptyOrgID(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	pair, err := issuer.IssueTokenPair("user-123", "9841234567", "member", "")
	if err != nil {
		t.Fatalf("IssueTokenPair() error = %v", err)
	}

	claims, err := issuer.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.OrgID != "" {
		t.Errorf("OrgID = %q, want empty", claims.OrgID)
	}
}

func TestUniqueJTIs(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	pair1, _ := issuer.IssueTokenPair("user-123", "9841234567", "member", "")
	pair2, _ := issuer.IssueTokenPair("user-123", "9841234567", "member", "")

	if pair1.AccessJTI == pair2.AccessJTI {
		t.Error("access JTIs should be unique across token pairs")
	}
	if pair1.RefreshJTI == pair2.RefreshJTI {
		t.Error("refresh JTIs should be unique across token pairs")
	}
	if pair1.AccessJTI == pair1.RefreshJTI {
		t.Error("access and refresh JTI should differ within the same pair")
	}
}

func TestInvalidTokenString(t *testing.T) {
	issuer := NewIssuer(testAccessSecret, testRefreshSecret)

	_, err := issuer.ValidateAccessToken("garbage-token")
	if err == nil {
		t.Error("expected error for garbage access token")
	}

	_, err = issuer.ValidateRefreshToken("garbage-token")
	if err == nil {
		t.Error("expected error for garbage refresh token")
	}
}
