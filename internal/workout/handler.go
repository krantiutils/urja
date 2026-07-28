package workout

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for workout endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new workout handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// --- Template handlers (org-scoped) ---

// ListTemplates handles GET /api/v1/orgs/{orgId}/workout-templates
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	templates, err := h.service.ListTemplates(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list templates", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": templates})
}

// GetTemplate handles GET /api/v1/orgs/{orgId}/workout-templates/{id}
func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	templateID := chi.URLParam(r, "id")
	if templateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing template ID"})
		return
	}

	t, err := h.service.GetTemplate(r.Context(), orgID, templateID)
	if err != nil {
		h.logger.Error("failed to get template", "error", err, "template_id", templateID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// CreateTemplate handles POST /api/v1/orgs/{orgId}/workout-templates
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateTemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	t, err := h.service.CreateTemplate(r.Context(), orgID, &input, userID)
	if err != nil {
		h.logger.Error("failed to create template", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

// UpdateTemplate handles PUT /api/v1/orgs/{orgId}/workout-templates/{id}
func (h *Handler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	templateID := chi.URLParam(r, "id")
	if templateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing template ID"})
		return
	}

	var input UpdateTemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	t, err := h.service.UpdateTemplate(r.Context(), orgID, templateID, &input)
	if err != nil {
		h.logger.Error("failed to update template", "error", err, "template_id", templateID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// DeleteTemplate handles DELETE /api/v1/orgs/{orgId}/workout-templates/{id}
func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	templateID := chi.URLParam(r, "id")
	if templateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing template ID"})
		return
	}

	if err := h.service.DeleteTemplate(r.Context(), orgID, templateID); err != nil {
		h.logger.Error("failed to delete template", "error", err, "template_id", templateID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "template deleted"})
}

// --- Plan assignment handler (org-scoped, under /members/{memberId}) ---

type assignPlanRequest struct {
	WorkoutTemplateID string `json:"workout_template_id"`
}

// AssignPlan handles POST /api/v1/orgs/{orgId}/members/{memberId}/assign-plan
func (h *Handler) AssignPlan(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing member ID"})
		return
	}

	var req assignPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	p, err := h.service.AssignPlan(r.Context(), memberID, orgID, req.WorkoutTemplateID, userID)
	if err != nil {
		h.logger.Error("failed to assign plan", "error", err, "member_id", memberID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// --- Member self-service handlers ---

// GetMyPlan handles GET /api/v1/members/me/workout-plan
func (h *Handler) GetMyPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID, err := h.service.ResolveOrg(r.Context(), userID, r.URL.Query().Get("organization_id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	plan, err := h.service.GetPlan(r.Context(), userID, orgID)
	if err != nil {
		h.logger.Error("failed to get workout plan", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no workout plan assigned"})
		return
	}

	writeJSON(w, http.StatusOK, plan)
}

// CreateMyLog handles POST /api/v1/members/me/workout-logs
func (h *Handler) CreateMyLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateLogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	orgID, err := h.service.ResolveOrg(r.Context(), userID, input.OrgID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	input.OrgID = orgID

	l, err := h.service.CreateLog(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to log workout", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, l)
}

// ListMyLogs handles GET /api/v1/members/me/workout-logs
func (h *Handler) ListMyLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID, err := h.service.ResolveOrg(r.Context(), userID, r.URL.Query().Get("organization_id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	var from, to *time.Time
	if f := r.URL.Query().Get("from"); f != "" {
		parsed, err := time.Parse("2006-01-02", f)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid from date, use YYYY-MM-DD"})
			return
		}
		from = &parsed
	}
	if t := r.URL.Query().Get("to"); t != "" {
		parsed, err := time.Parse("2006-01-02", t)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid to date, use YYYY-MM-DD"})
			return
		}
		// Include the full day
		endOfDay := parsed.Add(24*time.Hour - time.Nanosecond)
		to = &endOfDay
	}

	logs, err := h.service.ListLogs(r.Context(), userID, orgID, from, to, limit, offset)
	if err != nil {
		h.logger.Error("failed to list workout logs", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": logs})
}

// --- Browse, self-assign, and recommend handlers ---

// BrowseTemplates handles GET /api/v1/members/me/workout-templates
func (h *Handler) BrowseTemplates(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("organization_id")

	goal := r.URL.Query().Get("goal")
	difficulty := r.URL.Query().Get("difficulty")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	templates, err := h.service.BrowseTemplates(r.Context(), orgID, goal, difficulty, limit, offset)
	if err != nil {
		h.logger.Error("failed to browse templates", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": templates})
}

type selfAssignPlanRequest struct {
	OrgID             string `json:"organization_id"`
	WorkoutTemplateID string `json:"workout_template_id"`
}

// SelfAssignPlan handles POST /api/v1/members/me/workout-plan
func (h *Handler) SelfAssignPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req selfAssignPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	p, err := h.service.SelfAssignPlan(r.Context(), userID, req.OrgID, req.WorkoutTemplateID)
	if err != nil {
		h.logger.Error("failed to self-assign plan", "error", err, "user_id", userID, "org_id", req.OrgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

type recommendPlanRequest struct {
	OrgID       string `json:"organization_id"`
	Goal        string `json:"goal"`
	Experience  string `json:"experience"`
	DaysPerWeek string `json:"days_per_week"`
}

// RecommendPlan handles POST /api/v1/members/me/workout-plan/recommend
func (h *Handler) RecommendPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req recommendPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	p, err := h.service.RecommendPlan(r.Context(), userID, req.OrgID, req.Goal, req.Experience, req.DaysPerWeek)
	if err != nil {
		h.logger.Error("failed to recommend plan", "error", err, "user_id", userID, "org_id", req.OrgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
