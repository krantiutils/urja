package activitylog

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped activity log routes (under /orgs/{orgId}).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/activity-logs", h.List)
}
