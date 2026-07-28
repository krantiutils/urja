package activitylog

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped activity log routes (under /orgs/{orgId}).
//
// Staff only: the log records who did what inside the gym — members added,
// packages assigned, access blocked — which is the gym's internal operations
// and not a member's business. Nothing writes to this table yet, so gating it
// now costs nothing and stops the leak arriving with the first writer.
//
// Grouped rather than r.Use: this mux is shared with sibling handlers that
// register routes on it first, and chi panics if middleware is added to a mux
// that already has routes.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Get("/activity-logs", h.List)
	})
}
