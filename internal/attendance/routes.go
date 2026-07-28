package attendance

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped attendance routes (under
// /orgs/{orgId}/attendance).
//
// Staff only. OrgScope proves membership but not authority, so without this a
// plain member could check any member in — the handler takes a member id and
// already calls the caller "staffUserID" — and could read the whole gym's
// attendance history. Faked attendance is not cosmetic here: it feeds streaks,
// the leaderboard and the absentee SMS job. Members see their own attendance
// through the /members/me routes below.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Get("/", h.ListByOrg)
		r.Get("/weekly", h.WeeklySummary)
		r.Post("/check-in", h.ManualCheckIn)
	})
}

// RegisterSelfRoutes mounts member self-service routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Post("/attendance/check-in", h.SelfCheckIn)
	r.Get("/attendance", h.ListMine)
	r.Get("/streaks", h.GetMyStreaks)
	r.Get("/attendance/calendar", h.GetMyCalendar)
}
