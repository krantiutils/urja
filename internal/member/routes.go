package member

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterSelfRoutes mounts member self-service routes (under /members/me).
// All routes require authentication (handled by parent router group).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/", h.GetMe)
	r.Put("/", h.UpdateMe)
	r.Put("/privacy", h.UpdatePrivacy)
	r.Get("/leaderboard", h.GetMyLeaderboard)
}

// RegisterOrgRoutes mounts org-scoped member management routes (under /orgs/{orgId}/members).
// OrgScope middleware is applied by the parent router group.
// Note: PUT/DELETE for individual members are registered via RegisterMemberRoutes
// inside the /{memberId} subrouter to avoid chi route shadowing.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	// Any org member can list members
	r.Get("/", h.ListByOrg)

	// Staff and admin only
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))
		r.Post("/", h.CreateOrgMember)
	})
}

// RegisterMemberRoutes mounts routes for a specific member (under /orgs/{orgId}/members/{memberId}).
// Must be called inside the /{memberId} subrouter.
func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))
		r.Put("/", h.UpdateOrgMember)
		r.Delete("/", h.DeleteOrgMember)
	})
}
