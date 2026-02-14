package member

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts member self-service routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/", h.GetMe)
	r.Put("/", h.UpdateMe)
	// TODO: GET /attendance, GET /packages, GET /streaks, etc.
}

// RegisterOrgRoutes mounts org-scoped member management routes (under /orgs/{orgId}/members).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.ListByOrg)
	// TODO: GET /{id}, POST /, PUT /{id}, DELETE /{id}
}
