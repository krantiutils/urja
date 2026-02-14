package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/urja-gym/backend/internal/model"
	"github.com/urja-gym/backend/internal/repository"
	ujwt "github.com/urja-gym/backend/pkg/jwt"
	"github.com/urja-gym/backend/pkg/phone"
)

const (
	otpLength     = 6
	otpExpiry     = 5 * time.Minute
	otpMaxAttempts = 3
	otpCooldown   = 60 * time.Second
	otpHourlyMax  = 5

	bcryptCost = 12

	// Redis key prefixes
	keyOTPHash     = "otp:"
	keyOTPAttempts = "otp_attempts:"
	keyOTPCooldown = "otp_cooldown:"
	keyOTPHourly   = "otp_hourly:"
	keyRefresh     = "refresh:"
	keyUserSessions = "user_sessions:"
)

var (
	ErrOTPCooldown      = errors.New("please wait before requesting another OTP")
	ErrOTPRateLimit     = errors.New("too many OTP requests, try again later")
	ErrOTPExpired       = errors.New("OTP expired or not found")
	ErrOTPMaxAttempts   = errors.New("too many failed attempts, request a new OTP")
	ErrOTPInvalid       = errors.New("invalid OTP code")
	ErrInvalidPhone     = errors.New("invalid phone number")
	ErrUserDeactivated  = errors.New("account is deactivated")
	ErrInvalidRefresh   = errors.New("invalid or expired refresh token")
	ErrSessionNotFound  = errors.New("session not found or already invalidated")
)

// AuthService handles authentication business logic.
type AuthService struct {
	rdb      *redis.Client
	userRepo *repository.UserRepository
	sms      SMSSender
	jwt      *ujwt.Issuer
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	rdb *redis.Client,
	userRepo *repository.UserRepository,
	sms SMSSender,
	jwtIssuer *ujwt.Issuer,
) *AuthService {
	return &AuthService{
		rdb:      rdb,
		userRepo: userRepo,
		sms:      sms,
		jwt:      jwtIssuer,
	}
}

// SendOTP generates an OTP, stores its bcrypt hash in Redis, and sends it via SMS.
func (s *AuthService) SendOTP(ctx context.Context, rawPhone string) error {
	normalized, err := phone.Normalize(rawPhone)
	if err != nil {
		return ErrInvalidPhone
	}

	// Check cooldown (60s between OTP requests)
	cooldownKey := keyOTPCooldown + normalized
	exists, err := s.rdb.Exists(ctx, cooldownKey).Result()
	if err != nil {
		return fmt.Errorf("checking OTP cooldown: %w", err)
	}
	if exists > 0 {
		return ErrOTPCooldown
	}

	// Check hourly rate limit (5/hour per phone)
	hourlyKey := keyOTPHourly + normalized
	count, err := s.rdb.Get(ctx, hourlyKey).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("checking hourly rate: %w", err)
	}
	if count >= int64(otpHourlyMax) {
		return ErrOTPRateLimit
	}

	// Generate cryptographically secure 6-digit OTP
	otp, err := generateOTP(otpLength)
	if err != nil {
		return fmt.Errorf("generating OTP: %w", err)
	}

	// Hash OTP with bcrypt before storing
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcryptCost)
	if err != nil {
		return fmt.Errorf("hashing OTP: %w", err)
	}

	// Store OTP hash with 5min TTL
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, keyOTPHash+normalized, string(hash), otpExpiry)

	// Reset attempt counter
	pipe.Del(ctx, keyOTPAttempts+normalized)

	// Set cooldown (60s)
	pipe.Set(ctx, cooldownKey, "1", otpCooldown)

	// Increment hourly counter
	pipe.Incr(ctx, hourlyKey)
	// Set TTL on hourly key if it's new (first request in the window)
	pipe.Expire(ctx, hourlyKey, time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("storing OTP in redis: %w", err)
	}

	// Send OTP via SMS
	if err := s.sms.SendOTP(ctx, normalized, otp); err != nil {
		log.Error().Err(err).Str("phone", normalized).Msg("failed to send OTP SMS")
		// Clean up Redis state on SMS failure so user can retry immediately
		cleanPipe := s.rdb.Pipeline()
		cleanPipe.Del(ctx, keyOTPHash+normalized)
		cleanPipe.Del(ctx, cooldownKey)
		if _, cleanErr := cleanPipe.Exec(ctx); cleanErr != nil {
			log.Error().Err(cleanErr).Msg("failed to clean up OTP state after SMS failure")
		}
		return fmt.Errorf("sending OTP: %w", err)
	}

	return nil
}

// VerifyOTP checks the OTP against the stored hash and issues JWT tokens on success.
func (s *AuthService) VerifyOTP(ctx context.Context, rawPhone, otpCode string) (*model.AuthResponse, error) {
	normalized, err := phone.Normalize(rawPhone)
	if err != nil {
		return nil, ErrInvalidPhone
	}

	// Check attempt counter
	attemptsKey := keyOTPAttempts + normalized
	attempts, err := s.rdb.Get(ctx, attemptsKey).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("checking attempts: %w", err)
	}
	if attempts >= int64(otpMaxAttempts) {
		// Clean up OTP to force re-request
		s.rdb.Del(ctx, keyOTPHash+normalized, attemptsKey)
		return nil, ErrOTPMaxAttempts
	}

	// Fetch stored hash
	hashKey := keyOTPHash + normalized
	storedHash, err := s.rdb.Get(ctx, hashKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrOTPExpired
	}
	if err != nil {
		return nil, fmt.Errorf("fetching OTP hash: %w", err)
	}

	// Increment attempt counter before verification (prevents timing attacks on the counter)
	pipe := s.rdb.Pipeline()
	pipe.Incr(ctx, attemptsKey)
	pipe.Expire(ctx, attemptsKey, otpExpiry)
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("incrementing attempts: %w", err)
	}

	// Verify OTP against bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(otpCode)); err != nil {
		return nil, ErrOTPInvalid
	}

	// OTP verified — clean up OTP state
	s.rdb.Del(ctx, hashKey, attemptsKey, keyOTPCooldown+normalized)

	// Find or create user
	user, err := s.userRepo.FindByPhone(ctx, normalized)
	if errors.Is(err, repository.ErrUserNotFound) {
		user, err = s.userRepo.Create(ctx, normalized)
		if err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("looking up user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserDeactivated
	}

	// Issue JWT pair
	tokens, err := s.jwt.IssueTokenPair(user.ID, user.Phone, user.Role, "")
	if err != nil {
		return nil, fmt.Errorf("issuing tokens: %w", err)
	}

	// Store refresh token session in Redis
	if err := s.storeSession(ctx, user.ID, tokens.RefreshJTI); err != nil {
		return nil, fmt.Errorf("storing session: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user,
	}, nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*model.AuthResponse, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidRefresh
	}

	// Check if session exists in Redis
	sessionKey := keyRefresh + claims.ID
	userID, err := s.rdb.Get(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("checking session: %w", err)
	}

	if userID != claims.Subject {
		return nil, ErrInvalidRefresh
	}

	// Look up user for fresh claims
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("looking up user for refresh: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserDeactivated
	}

	// Invalidate old session
	s.rdb.Del(ctx, sessionKey)
	s.rdb.SRem(ctx, keyUserSessions+userID, claims.ID)

	// Issue new token pair
	tokens, err := s.jwt.IssueTokenPair(user.ID, user.Phone, user.Role, "")
	if err != nil {
		return nil, fmt.Errorf("issuing refreshed tokens: %w", err)
	}

	// Store new session
	if err := s.storeSession(ctx, user.ID, tokens.RefreshJTI); err != nil {
		return nil, fmt.Errorf("storing refreshed session: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user,
	}, nil
}

// Logout invalidates a single refresh token session.
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	claims, err := s.jwt.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		// Token already invalid/expired — treat as success
		return nil
	}

	sessionKey := keyRefresh + claims.ID
	s.rdb.Del(ctx, sessionKey)
	s.rdb.SRem(ctx, keyUserSessions+claims.Subject, claims.ID)

	return nil
}

// LogoutAll invalidates all sessions for a user.
func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	sessionsKey := keyUserSessions + userID
	sessionIDs, err := s.rdb.SMembers(ctx, sessionsKey).Result()
	if err != nil {
		return fmt.Errorf("fetching user sessions: %w", err)
	}

	if len(sessionIDs) > 0 {
		// Build keys to delete
		keys := make([]string, 0, len(sessionIDs)+1)
		for _, sid := range sessionIDs {
			keys = append(keys, keyRefresh+sid)
		}
		keys = append(keys, sessionsKey)

		if err := s.rdb.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("deleting sessions: %w", err)
		}
	}

	return nil
}

// storeSession records a refresh token session in Redis.
func (s *AuthService) storeSession(ctx context.Context, userID, jti string) error {
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, keyRefresh+jti, userID, ujwt.RefreshTokenExpiry)
	pipe.SAdd(ctx, keyUserSessions+userID, jti)
	pipe.Expire(ctx, keyUserSessions+userID, ujwt.RefreshTokenExpiry)
	_, err := pipe.Exec(ctx)
	return err
}

// generateOTP generates a cryptographically secure numeric OTP of the given length.
func generateOTP(length int) (string, error) {
	max := new(big.Int)
	max.SetInt64(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10))
	}

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Zero-pad to the requested length
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n), nil
}
