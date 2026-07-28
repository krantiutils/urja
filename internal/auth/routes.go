package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts the auth routes on the given chi router.
// validator is used for the authenticated logout routes.
func (h *Handler) RegisterRoutes(r chi.Router, validator middleware.TokenValidator) {
	// Auth-specific rate limiter: 5 requests per 15 minutes per IP
	authLimiter := middleware.NewRateLimiter(5.0/(15*60), 5)

	r.Route("/auth", func(r chi.Router) {
		// Public endpoints (rate limited)
		r.Group(func(r chi.Router) {
			r.Use(authLimiter.Limit())
			r.Post("/login", h.Login)
			r.Post("/verify-otp", h.VerifyOTP)
			// Shares the 5-per-15-minutes limiter with OTP, so password
			// guessing is throttled on the same budget as OTP requests.
			r.Post("/password-login", h.PasswordLogin)
		})

		// Refresh endpoint (no auth required, token is in the body)
		r.Post("/refresh", h.Refresh)

		// Authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(validator))
			r.Get("/password", h.PasswordStatus)
			r.Post("/password", h.SetPassword)
			r.Post("/logout", h.Logout)
			r.Post("/logout-all", h.LogoutAll)
		})
	})
}
