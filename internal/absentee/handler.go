package absentee

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for absentee endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new absentee handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type notifyRequest struct {
	OrgName string `json:"org_name"`
}

// List handles GET /api/v1/orgs/{orgId}/absentees
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	minDays, _ := strconv.Atoi(r.URL.Query().Get("min_days"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	absentees, err := h.service.List(r.Context(), orgID, minDays, limit, offset)
	if err != nil {
		h.logger.Error("failed to list absentees", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": absentees})
}

// Notify handles POST /api/v1/orgs/{orgId}/absentees/{memberId}/notify
func (h *Handler) Notify(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing member ID"})
		return
	}

	var req notifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.OrgName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "org_name is required"})
		return
	}

	if err := h.service.Notify(r.Context(), orgID, memberID, req.OrgName); err != nil {
		h.logger.Error("failed to notify absentee", "error", err, "member_id", memberID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "notification sent"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
