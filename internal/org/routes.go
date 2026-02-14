package org

import (
	"github.com/go-chi/chi/v5"
)

// RegisterPublicRoutes mounts public gym listing routes.
// GET /api/v1/gyms      → List
// GET /api/v1/gyms/{id} → Get
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}

// RegisterAdminRoutes mounts authenticated org management routes.
// POST /api/v1/orgs            → Create (super_admin)
// PUT  /api/v1/orgs/{orgId}    → Update (org admin)
func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Route("/{orgId}", func(r chi.Router) {
		r.Put("/", h.Update)
	})
}
