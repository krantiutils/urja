package absentee

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped absentee routes (under /orgs/{orgId}/absentees).
// Only staff/admin can view absentees and send notifications.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireRole("staff", "admin"))

	r.Get("/", h.List)
	r.Post("/{memberId}/notify", h.Notify)
}
