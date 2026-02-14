package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"

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

// RequestOTP generates a 6-digit OTP, hashes it, stores it in Redis, and sends via SMS.
// It enforces rate limiting: max OTPHourlyLimit requests per hour per phone.
func (s *Service) RequestOTP(ctx context.Context, phone string) error {
	if !sms.ValidatePhone(phone) {
		return fmt.Errorf("invalid phone number")
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

// VerifyOTP verifies the OTP against the stored hash and issues tokens.
// It enforces max attempts per verification window.
func (s *Service) VerifyOTP(ctx context.Context, phone, otp string) (accessToken, refreshToken string, err error) {
	// Check attempt count
	attempts, err := s.repo.IncrementOTPAttempts(ctx, phone, s.cfg.OTPExpiry)
	if err != nil {
		return "", "", fmt.Errorf("checking OTP attempts: %w", err)
	}
	if attempts > int64(s.cfg.OTPMaxAttempts) {
		return "", "", fmt.Errorf("too many OTP attempts")
	}

	// Get stored hash
	storedHash, err := s.repo.GetOTPHash(ctx, phone)
	if err != nil {
		return "", "", fmt.Errorf("OTP expired or not found")
	}

	// Verify
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(otp)); err != nil {
		return "", "", fmt.Errorf("invalid OTP")
	}

	// OTP verified — clean up
	if err := s.repo.DeleteOTP(ctx, phone); err != nil {
		s.logger.Warn("failed to delete OTP after verification", "error", err)
	}

	// Find or create user
	userID, role, err := s.repo.FindOrCreateUserByPhone(ctx, phone)
	if err != nil {
		return "", "", fmt.Errorf("finding user: %w", err)
	}

	// Issue tokens (PRD: tokens only issued AFTER OTP verification)
	// TODO: Implement JWT token generation
	_ = userID
	_ = role

	return "", "", fmt.Errorf("JWT token generation not yet implemented")
}

// ValidateAccessToken validates a JWT access token and returns the user ID and role.
// This implements the middleware.TokenValidator interface.
func (s *Service) ValidateAccessToken(_ string) (string, string, error) {
	// TODO: Implement JWT validation
	return "", "", fmt.Errorf("not implemented")
}

// generateOTP generates a cryptographically secure 6-digit OTP.
func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
