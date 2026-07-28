package boxing

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterSelfRoutes mounts member self-service routes (under /members/me).
// The organization is supplied as a query parameter and verified against real
// membership in the handler — these routes are not under /orgs/{orgId}, so
// OrgScope does not run.
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/boxing", h.GetMyProfile)
	r.Put("/boxing", h.UpdateMyProfile)

	r.Get("/bouts", h.ListMyBouts)
	r.Post("/bouts", h.CreateMyBout)
	r.Delete("/bouts/{boutId}", h.DeleteMyBout)
}

// RegisterMemberRoutes mounts staff-facing routes (under
// /orgs/{orgId}/members/{memberId}).
//
// Grouped rather than using r.Use at the top level: this mux is shared with
// sibling handlers that register routes on it first, and chi panics if
// middleware is added to a mux that already has routes.
func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Get("/boxing", h.GetMemberProfile)
		r.Put("/boxing", h.UpdateMemberProfile)
		// Sparring clearance is a safety gate, so it is a staff-only route and
		// the service additionally refuses self-clearance.
		r.Put("/boxing/sparring-clearance", h.SetSparringClearance)

		r.Post("/bouts", h.CreateMemberBout)
		r.Delete("/bouts/{boutId}", h.DeleteMemberBout)
	})
}
