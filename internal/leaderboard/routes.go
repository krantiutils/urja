package leaderboard

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped leaderboard routes (under /orgs/{orgId}/leaderboard).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.GetLeaderboard)
}
