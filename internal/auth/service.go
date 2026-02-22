package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urja-gym/urja/internal/config"
	"github.com/urja-gym/urja/pkg/sms"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic.
type Service struct {
	repo      *Repository
	smsClient *sms.Client
	cfg       config.AuthConfig
	logger    *slog.Logger
}

// NewService creates a new auth service.
func NewService(repo *Repository, smsClient *sms.Client, cfg config.AuthConfig, logger *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		smsClient: smsClient,
		cfg:       cfg,
		logger:    logger,
	}
}

// AccessTokenClaims holds the claims embedded in a JWT access token.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Phone string `json:"phone"`
	Role  string `json:"role"`
	OrgID string `json:"org_id,omitempty"`
}

// RefreshTokenClaims holds the claims embedded in a JWT refresh token.
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	TokenID string `json:"jti"`
}

// RequestOTP generates a 6-digit OTP, hashes it, stores it in Redis, and sends via SMS.
// It enforces rate limiting: max OTPHourlyLimit requests per hour per phone.
func (s *Service) RequestOTP(ctx context.Context, phone string) error {
	phone = sms.NormalizePhone(phone)
	if !sms.ValidatePhone(phone) {
		return fmt.Errorf("invalid phone number")
	}

	// Dev bypass: skip SMS for test phones (9800*)
	if s.cfg.DevOTPBypass && strings.HasPrefix(phone, "9800") {
		hash, err := bcrypt.GenerateFromPassword([]byte("123456"), s.cfg.BcryptCost)
		if err != nil {
			return fmt.Errorf("hashing dev OTP: %w", err)
		}
		if err := s.repo.StoreOTPHash(ctx, phone, string(hash), s.cfg.OTPExpiry); err != nil {
			return fmt.Errorf("storing dev OTP: %w", err)
		}
		s.logger.Info("dev OTP bypass: stored test OTP", "phone", phone[:4]+"******")
		return nil
	}

	// Check hourly rate limit
	count, err := s.repo.GetOTPRequestCount(ctx, phone)
	if err == nil && count >= int64(s.cfg.OTPHourlyLimit) {
		return fmt.Errorf("OTP request limit exceeded, try again later")
	}

	// Generate 6-digit OTP
	otp, err := generateOTP()
	if err != nil {
		return fmt.Errorf("generating OTP: %w", err)
	}

	// Hash OTP before storage (PRD: OTP stored as bcrypt hash, NEVER returned in response)
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), s.cfg.BcryptCost)
	if err != nil {
		return fmt.Errorf("hashing OTP: %w", err)
	}

	if err := s.repo.StoreOTPHash(ctx, phone, string(hash), s.cfg.OTPExpiry); err != nil {
		return fmt.Errorf("storing OTP: %w", err)
	}

	if err := s.repo.IncrementOTPRequestCount(ctx, phone); err != nil {
		s.logger.Warn("failed to increment OTP rate counter", "error", err, "phone", phone[:4]+"******")
	}

	// Send SMS
	if err := s.smsClient.SendOTP(ctx, phone, otp); err != nil {
		return fmt.Errorf("sending OTP: %w", err)
	}

	return nil
}

// VerifyOTPResult contains the result of OTP verification.
type VerifyOTPResult struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	IsNewUser           bool   `json:"is_new_user"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
}

// VerifyOTP verifies the OTP against the stored hash and issues tokens.
// It enforces max attempts per verification window.
func (s *Service) VerifyOTP(ctx context.Context, phone, otp string) (*VerifyOTPResult, error) {
	phone = sms.NormalizePhone(phone)

	// Check attempt count
	attempts, err := s.repo.IncrementOTPAttempts(ctx, phone, s.cfg.OTPExpiry)
	if err != nil {
		return nil, fmt.Errorf("checking OTP attempts: %w", err)
	}
	if attempts > int64(s.cfg.OTPMaxAttempts) {
		return nil, fmt.Errorf("too many OTP attempts")
	}

	// Get stored hash
	storedHash, err := s.repo.GetOTPHash(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("OTP expired or not found")
	}

	// Verify
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(otp)); err != nil {
		return nil, fmt.Errorf("invalid OTP")
	}

	// OTP verified — clean up
	if err := s.repo.DeleteOTP(ctx, phone); err != nil {
		s.logger.Warn("failed to delete OTP after verification", "error", err)
	}

	// Find or create user
	userID, role, isNew, err := s.repo.FindOrCreateUserByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	// Check onboarding status
	var onboardingCompleted bool
	if !isNew {
		onboardingCompleted, err = s.repo.GetUserOnboardingStatus(ctx, userID)
		if err != nil {
			s.logger.Warn("failed to get onboarding status, defaulting to false", "error", err)
		}
	}

	// Issue tokens (PRD: tokens only issued AFTER OTP verification)
	accessToken, err := s.generateAccessToken(userID, phone, role)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, tokenID, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Store refresh token in Redis for revocation tracking
	if err := s.repo.StoreRefreshToken(ctx, userID, tokenID, s.cfg.RefreshTokenExpiry); err != nil {
		return nil, fmt.Errorf("storing refresh token: %w", err)
	}

	s.logger.Info("user authenticated", "user_id", userID, "phone", phone[:4]+"******", "is_new", isNew)
	return &VerifyOTPResult{
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		IsNewUser:           isNew,
		OnboardingCompleted: onboardingCompleted,
	}, nil
}

// RefreshAccessToken validates a refresh token and issues a new access token pair.
func (s *Service) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (newAccessToken, newRefreshToken string, err error) {
	// Parse the refresh token
	token, err := jwt.ParseWithClaims(refreshTokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid refresh token claims")
	}

	userID := claims.Subject
	tokenID := claims.TokenID

	// Check if this refresh token is still valid in Redis (not revoked)
	valid, err := s.repo.ValidateRefreshToken(ctx, userID, tokenID)
	if err != nil {
		return "", "", fmt.Errorf("validating refresh token: %w", err)
	}
	if !valid {
		return "", "", fmt.Errorf("refresh token has been revoked")
	}

	// Revoke old refresh token (single-use rotation)
	if err := s.repo.RevokeRefreshToken(ctx, userID, tokenID); err != nil {
		s.logger.Warn("failed to revoke old refresh token", "error", err, "user_id", userID)
	}

	// Look up user to get current role and phone
	phone, role, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("looking up user: %w", err)
	}

	// Issue new token pair
	newAccessToken, err = s.generateAccessToken(userID, phone, role)
	if err != nil {
		return "", "", fmt.Errorf("generating access token: %w", err)
	}

	newRefreshToken, newTokenID, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	if err := s.repo.StoreRefreshToken(ctx, userID, newTokenID, s.cfg.RefreshTokenExpiry); err != nil {
		return "", "", fmt.Errorf("storing refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout revokes a specific refresh token.
func (s *Service) Logout(ctx context.Context, userID, refreshTokenStr string) error {
	token, err := jwt.ParseWithClaims(refreshTokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid refresh token claims")
	}

	// Ensure the token belongs to this user
	if claims.Subject != userID {
		return fmt.Errorf("token does not belong to this user")
	}

	return s.repo.RevokeRefreshToken(ctx, userID, claims.TokenID)
}

// LogoutAll revokes all refresh tokens for the user.
func (s *Service) LogoutAll(ctx context.Context, userID string) error {
	return s.repo.RevokeAllRefreshTokens(ctx, userID)
}

// ValidateAccessToken validates a JWT access token and returns the user ID and role.
// This implements the middleware.TokenValidator interface.
func (s *Service) ValidateAccessToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("invalid access token: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid token claims")
	}

	return claims.Subject, claims.Role, nil
}

func (s *Service) generateAccessToken(userID, phone, role string) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenExpiry)),
			Issuer:    "urja",
		},
		Phone: phone,
		Role:  role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *Service) generateRefreshToken(userID string) (tokenStr string, tokenID string, err error) {
	tokenID, err = generateTokenID()
	if err != nil {
		return "", "", fmt.Errorf("generating token ID: %w", err)
	}

	now := time.Now()
	claims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.RefreshTokenExpiry)),
			Issuer:    "urja",
		},
		TokenID: tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}
	return tokenStr, tokenID, nil
}

// generateOTP generates a cryptographically secure 6-digit OTP.
func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// generateTokenID generates a cryptographically secure random token ID.
func generateTokenID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
