package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Repository handles data persistence for authentication.
type Repository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

// NewRepository creates a new auth repository.
func NewRepository(db *pgxpool.Pool, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb}
}

// StoreOTPHash stores a bcrypt-hashed OTP in Redis with expiration.
func (r *Repository) StoreOTPHash(ctx context.Context, phone, hash string, expiry time.Duration) error {
	key := fmt.Sprintf("otp:%s", phone)
	return r.rdb.Set(ctx, key, hash, expiry).Err()
}

// GetOTPHash retrieves the stored OTP hash for a phone number.
func (r *Repository) GetOTPHash(ctx context.Context, phone string) (string, error) {
	key := fmt.Sprintf("otp:%s", phone)
	return r.rdb.Get(ctx, key).Result()
}

// DeleteOTP removes the OTP entry after successful verification.
func (r *Repository) DeleteOTP(ctx context.Context, phone string) error {
	key := fmt.Sprintf("otp:%s", phone)
	return r.rdb.Del(ctx, key).Err()
}

// IncrementOTPAttempts tracks OTP verification attempts.
func (r *Repository) IncrementOTPAttempts(ctx context.Context, phone string, window time.Duration) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s", phone)
	pipe := r.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// GetOTPRequestCount returns how many OTP requests a phone has made in the time window.
func (r *Repository) GetOTPRequestCount(ctx context.Context, phone string) (int64, error) {
	key := fmt.Sprintf("otp_rate:%s", phone)
	return r.rdb.Get(ctx, key).Int64()
}

// IncrementOTPRequestCount increments the OTP request counter with hourly window.
func (r *Repository) IncrementOTPRequestCount(ctx context.Context, phone string) error {
	key := fmt.Sprintf("otp_rate:%s", phone)
	pipe := r.rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

// StoreRefreshToken stores a refresh token in Redis.
func (r *Repository) StoreRefreshToken(ctx context.Context, userID, tokenID string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh:%s:%s", userID, tokenID)
	return r.rdb.Set(ctx, key, "1", expiry).Err()
}

// ValidateRefreshToken checks if a refresh token exists and is valid.
func (r *Repository) ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error) {
	key := fmt.Sprintf("refresh:%s:%s", userID, tokenID)
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// RevokeRefreshToken deletes a specific refresh token.
func (r *Repository) RevokeRefreshToken(ctx context.Context, userID, tokenID string) error {
	key := fmt.Sprintf("refresh:%s:%s", userID, tokenID)
	return r.rdb.Del(ctx, key).Err()
}

// RevokeAllRefreshTokens deletes all refresh tokens for a user.
func (r *Repository) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("refresh:%s:*", userID)
	iter := r.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := r.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// GetUserByID retrieves a user's phone and role by their ID.
func (r *Repository) GetUserByID(ctx context.Context, userID string) (phone, role string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT u.phone, COALESCE(
			(SELECT om.role FROM organization_members om WHERE om.user_id = u.id LIMIT 1),
			'member'
		) FROM users u WHERE u.id = $1`,
		userID,
	).Scan(&phone, &role)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}
	return phone, role, nil
}

// FindOrCreateUserByPhone looks up a user by phone number, creating one if not found.
// Returns the user ID, role, and whether the user was newly created.
func (r *Repository) FindOrCreateUserByPhone(ctx context.Context, phone string) (string, string, bool, error) {
	var userID, role string

	err := r.db.QueryRow(ctx,
		`SELECT id, COALESCE(
			(SELECT role FROM organization_members WHERE user_id = users.id LIMIT 1),
			'member'
		) FROM users WHERE phone = $1`,
		phone,
	).Scan(&userID, &role)

	if err != nil {
		// User doesn't exist; create one.
		err = r.db.QueryRow(ctx,
			`INSERT INTO users (phone) VALUES ($1) RETURNING id`,
			phone,
		).Scan(&userID)
		if err != nil {
			return "", "", false, fmt.Errorf("creating user: %w", err)
		}
		return userID, "member", true, nil
	}

	return userID, role, false, nil
}

// GetUserOnboardingStatus retrieves onboarding_completed for a user.
func (r *Repository) GetUserOnboardingStatus(ctx context.Context, userID string) (bool, error) {
	var completed bool
	err := r.db.QueryRow(ctx,
		`SELECT onboarding_completed FROM users WHERE id = $1`,
		userID,
	).Scan(&completed)
	if err != nil {
		return false, fmt.Errorf("getting onboarding status: %w", err)
	}
	return completed, nil
}
