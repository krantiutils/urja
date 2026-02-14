package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// AccessClaims holds the claims embedded in an access token.
type AccessClaims struct {
	Phone string `json:"phone"`
	Role  string `json:"role"`
	OrgID string `json:"org_id,omitempty"`
	jwt.RegisteredClaims
}

// RefreshClaims holds the claims embedded in a refresh token.
type RefreshClaims struct {
	jwt.RegisteredClaims
}

// TokenPair holds an access token and a refresh token along with their JTIs.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessJTI    string `json:"-"`
	RefreshJTI   string `json:"-"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Issuer generates and validates JWT tokens.
type Issuer struct {
	accessSecret  []byte
	refreshSecret []byte
}

// NewIssuer creates a new JWT issuer with the given secrets.
func NewIssuer(accessSecret, refreshSecret string) *Issuer {
	return &Issuer{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

// IssueTokenPair generates a new access + refresh token pair.
func (i *Issuer) IssueTokenPair(userID, phone, role, orgID string) (*TokenPair, error) {
	now := time.Now()
	accessJTI := uuid.New().String()
	refreshJTI := uuid.New().String()

	accessClaims := AccessClaims{
		Phone: phone,
		Role:  role,
		OrgID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			Issuer:    "urja",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(i.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshClaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        refreshJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExpiry)),
			Issuer:    "urja",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(i.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		AccessJTI:    accessJTI,
		RefreshJTI:   refreshJTI,
		ExpiresIn:    int64(AccessTokenExpiry.Seconds()),
	}, nil
}

// ValidateAccessToken parses and validates an access token string.
func (i *Issuer) ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing access token: %w", err)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}

	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token string.
func (i *Issuer) ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.refreshSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing refresh token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}
