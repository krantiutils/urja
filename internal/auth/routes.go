package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts the auth routes on the given chi router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Auth-specific rate limiter: 5 requests per 15 minutes per IP
	authLimiter := middleware.NewRateLimiter(5.0/(15*60), 5)

	r.Route("/auth", func(r chi.Router) {
		r.Use(authLimiter.Limit())
		r.Post("/login", h.Login)
		r.Post("/verify-otp", h.VerifyOTP)
		// TODO: POST /refresh, POST /logout, POST /logout-all
	})
}
