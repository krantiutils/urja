package attendance

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for attendance endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new attendance handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type checkInRequest struct {
	Method string `json:"method"` // qr, nfc, manual
}

// CheckIn handles POST /api/v1/orgs/{orgId}/attendance/check-in
func (h *Handler) CheckIn(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req checkInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	rec, err := h.service.CheckIn(r.Context(), userID, orgID, req.Method)
	if err != nil {
		h.logger.Error("check-in failed", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, rec)
}

// ListByOrg handles GET /api/v1/orgs/{orgId}/attendance
func (h *Handler) ListByOrg(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	records, err := h.service.ListByOrg(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list attendance", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": records})
}

// ListMine handles GET /api/v1/members/me/attendance
func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	records, err := h.service.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list attendance", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": records})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
