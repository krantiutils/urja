package middleware

import (
	"net/http"
	"sync"
	"time"
)

// visitor tracks rate limit state for a single client.
type visitor struct {
	tokens    float64
	lastSeen  time.Time
	maxTokens float64
	rate      float64 // tokens per second
}

// RateLimiter implements a token bucket rate limiter per client IP.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     float64
	burst    int
	cleanup  time.Duration
}

// NewRateLimiter creates a rate limiter.
// rate is requests per second, burst is the maximum burst size.
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  3 * time.Minute,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:    float64(rl.burst) - 1,
			lastSeen:  time.Now(),
			maxTokens: float64(rl.burst),
			rate:      rl.rate,
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(v.lastSeen).Seconds()
	v.lastSeen = now

	v.tokens += elapsed * v.rate
	if v.tokens > v.maxTokens {
		v.tokens = v.maxTokens
	}

	if v.tokens < 1 {
		return false
	}

	v.tokens--
	return true
}

// Limit returns middleware that rate limits requests by client IP.
func (rl *RateLimiter) Limit() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := realIP(r)
			if !rl.allow(ip) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// realIP extracts the client IP, preferring X-Real-IP and X-Forwarded-For headers.
func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
