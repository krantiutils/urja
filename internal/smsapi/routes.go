package smsapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped SMS routes (under /orgs/{orgId}/sms).
// All SMS operations require admin role.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireRole("admin"))

	r.Get("/balance", h.GetBalance)
	r.Post("/send", h.Send)
	r.Post("/buy", h.Buy)
	r.Get("/history", h.History)
}
