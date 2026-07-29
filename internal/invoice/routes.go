package invoice

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts invoice routes (under /orgs/{orgId}/invoices).
// Billing is an admin and staff function; members never reach these.
//
// There is deliberately no PUT and no DELETE: an issued bill is immutable and
// is corrected by cancelling it or raising a credit note.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("admin", "staff"))

	r.Get("/", h.List)
	r.Post("/", h.Issue)
	r.Get("/next-number", h.NextNumber)
	r.Get("/{invoiceId}", h.Get)
	r.Post("/{invoiceId}/cancel", h.Cancel)
	r.Post("/{invoiceId}/credit-note", h.CreditNote)
	r.Post("/{invoiceId}/print", h.Print)
}
