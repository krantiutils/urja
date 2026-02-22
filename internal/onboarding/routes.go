package onboarding

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts onboarding routes under /members/me.
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Post("/onboarding", h.CompleteOnboarding)
	r.Post("/join-gym", h.JoinGym)
	r.Post("/leave-gym", h.LeaveGym)
}
