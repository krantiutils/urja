package subscription

import (
	"github.com/go-chi/chi/v5"
)

// RegisterPackageRoutes mounts org-scoped package routes (under /orgs/{orgId}/packages).
func (h *Handler) RegisterPackageRoutes(r chi.Router) {
	r.Get("/summary", h.PackageSummary)
	r.Get("/expiring", h.ListExpiring)
	r.Get("/expired", h.ListExpired)
}

// RegisterMemberRoutes mounts member-scoped routes (under /orgs/{orgId}/members/{memberId}).
func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Route("/packages", func(r chi.Router) {
		r.Post("/assign", h.AssignPackage)
		r.Post("/renew", h.RenewPackage)
		r.Post("/extend", h.ExtendPackage)
	})
	r.Get("/payments", h.ListPayments)
	r.Get("/subscriptions", h.ListSubscriptions)
}
