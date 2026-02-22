package workout

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped workout template routes (under /orgs/{orgId}/workout-templates).
// All org members can list/get templates. Only staff/admin can create, update, delete.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.ListTemplates)
	r.Get("/{id}", h.GetTemplate)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("staff", "admin"))
		r.Post("/", h.CreateTemplate)
		r.Put("/{id}", h.UpdateTemplate)
		r.Delete("/{id}", h.DeleteTemplate)
	})
}

// RegisterMemberRoutes mounts plan assignment route (under /orgs/{orgId}/members/{memberId}).
// Only staff/admin can assign plans.
func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("staff", "admin"))
		r.Post("/assign-plan", h.AssignPlan)
	})
}

// RegisterSelfRoutes mounts member self-service workout routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/workout-plan", h.GetMyPlan)
	r.Get("/workout-logs", h.ListMyLogs)
	r.Post("/workout-logs", h.CreateMyLog)
	r.Get("/workout-templates", h.BrowseTemplates)
	r.Post("/workout-plan", h.SelfAssignPlan)
	r.Post("/workout-plan/recommend", h.RecommendPlan)
}
