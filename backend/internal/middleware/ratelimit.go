package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RateLimitConfig defines parameters for a rate limiter.
type RateLimitConfig struct {
	// Prefix is the Redis key prefix for this limiter.
	Prefix string
	// Max is the maximum number of requests allowed in the window.
	Max int64
	// Window is the time window for the rate limit.
	Window time.Duration
	// KeyFunc extracts the rate limit key from the request (e.g., IP, phone).
	// If nil, the client IP is used.
	KeyFunc func(r *http.Request) string
}

// RateLimiter creates a rate limiting middleware backed by Redis.
// Uses a sliding window counter approach.
func RateLimiter(rdb *redis.Client, cfg RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := clientIP(r)
			if cfg.KeyFunc != nil {
				key = cfg.KeyFunc(r)
			}

			redisKey := fmt.Sprintf("ratelimit:%s:%s", cfg.Prefix, key)

			allowed, err := checkRateLimit(r.Context(), rdb, redisKey, cfg.Max, cfg.Window)
			if err != nil {
				log.Error().Err(err).Str("key", redisKey).Msg("rate limit check failed")
				// Fail open — don't block requests if Redis is down
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(cfg.Window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"error":"rate_limit_exceeded","message":"Too many requests, please try again later"}`)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// checkRateLimit implements a simple sliding window counter using Redis INCR + EXPIRE.
func checkRateLimit(ctx context.Context, rdb *redis.Client, key string, max int64, window time.Duration) (bool, error) {
	pipe := rdb.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := incrCmd.Val()
	return count <= max, nil
}

// clientIP extracts the client IP from a request, respecting X-Forwarded-For
// and X-Real-IP headers (set by reverse proxies).
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain (client IP)
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
