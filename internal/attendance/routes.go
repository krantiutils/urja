package attendance

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped attendance routes (under /orgs/{orgId}/attendance).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.ListByOrg)
	r.Post("/check-in", h.ManualCheckIn)
}

// RegisterSelfRoutes mounts member self-service routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Post("/check-in", h.SelfCheckIn)
	r.Get("/attendance", h.ListMine)
	r.Get("/streaks", h.GetMyStreaks)
	r.Get("/calendar", h.GetMyCalendar)
}
