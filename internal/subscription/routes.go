package subscription

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterPackageRoutes mounts org-scoped package routes (under /orgs/{orgId}/packages).
// All package lifecycle operations require staff or admin role.
//
// This router is shared with packages.Handler.RegisterOrgRoutes (see
// cmd/api/main.go), which registers its own routes directly on the same mux — chi
// forbids calling Use() on a mux that already has routes registered on it, so the
// middleware must be scoped via Group() rather than applied at the top level here.
func (h *Handler) RegisterPackageRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Get("/summary", h.PackageSummary)
		r.Get("/expiring", h.ListExpiring)
		r.Get("/expired", h.ListExpired)
	})
}

// RegisterMemberRoutes mounts member-scoped routes (under /orgs/{orgId}/members/{memberId}).
// All member package/payment operations require staff or admin role — a plain member
// must never be able to assign themselves a package or set their own amount_paid.
//
// This router is shared with member.Handler.RegisterMemberRoutes and
// workout.Handler.RegisterMemberRoutes (see cmd/api/main.go); the middleware is
// scoped via Group() for the same reason as RegisterPackageRoutes above.
func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Route("/packages", func(r chi.Router) {
			r.Post("/assign", h.AssignPackage)
			r.Post("/renew", h.RenewPackage)
			r.Post("/extend", h.ExtendPackage)
		})
		r.Get("/payments", h.ListPayments)
		r.Get("/subscriptions", h.ListSubscriptions)
	})
}
