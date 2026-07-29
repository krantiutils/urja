package org

import (
	"github.com/go-chi/chi/v5"
)

// RegisterPublicRoutes mounts public gym listing routes.
// GET /api/v1/gyms      → List
// GET /api/v1/gyms/{id} → Get
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}

// The authenticated org routes (POST /api/v1/orgs, GET/PUT
// /api/v1/orgs/{orgId}) are registered inline in cmd/api/main.go — and
// mirrored in tests/e2e/setup_test.go — rather than through a helper here.
// A prior RegisterAdminRoutes existed for that purpose but was never called
// by either router, and it did not apply OrgScope to Put: a second,
// hand-maintained copy of the route tree that looked authoritative and
// shipped Update unprotected the moment anyone actually wired it in. Deleted
// rather than fixed in place; making main.go/setup_test.go call a single
// shared registration function instead of each inlining their own is a
// larger refactor for later, not a tail-end fix.
