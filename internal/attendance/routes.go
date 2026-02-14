package attendance

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped attendance routes (under /orgs/{orgId}/attendance).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.ListByOrg)
	r.Post("/check-in", h.CheckIn)
}

// RegisterSelfRoutes mounts member self-service attendance routes (under /members/me/attendance).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/", h.ListMine)
}
