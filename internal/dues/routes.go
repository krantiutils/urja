package dues

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts dues routes (under /orgs/{orgId}/dues).
// All dues routes require admin or staff role.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.RequireRole("admin", "staff"))

	r.Get("/", h.List)
	r.Post("/{memberId}/pay", h.Pay)

	// Block access is admin-only
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("admin"))
		r.Post("/{memberId}/block-access", h.BlockAccess)
	})
}
