package dues

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts dues routes (under /orgs/{orgId}/dues).
// All dues routes require admin or staff role.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("admin", "staff"))

	r.Get("/", h.List)
	r.Post("/", h.Create)
	// The path carries a due id, not a member id — it was named
	// {memberId} while the handler ignored it entirely.
	r.Post("/{dueId}/pay", h.Pay)

	// Block access is admin-only
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("admin"))
		r.Post("/{memberId}/block-access", h.BlockAccess)
	})
}
