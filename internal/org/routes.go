package org

import (
	"github.com/go-chi/chi/v5"
)

// RegisterPublicRoutes mounts public gym listing routes.
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
