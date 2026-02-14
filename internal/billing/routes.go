package billing

import (
	"github.com/go-chi/chi/v5"
)

// RegisterPublicRoutes mounts public billing routes.
// GET /api/v1/billing/plans → ListPlans
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/plans", h.ListPlans)
}

// RegisterAuthRoutes mounts authenticated billing routes.
// POST /api/v1/billing/subscribe → Subscribe
func (h *Handler) RegisterAuthRoutes(r chi.Router) {
	r.Post("/subscribe", h.Subscribe)
}
