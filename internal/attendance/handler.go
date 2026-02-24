package attendance

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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

type selfCheckInRequest struct {
	Method  string `json:"method"`   // qr, nfc
	QRToken string `json:"qr_token"` // for qr: base64url_payload.hex_signature
	CardUID string `json:"card_uid"` // for nfc: hex card identifier
	OrgID   string `json:"org_id"`   // for nfc: organization ID
}

type manualCheckInRequest struct {
	MemberID string `json:"member_id"`
}

// SelfCheckIn handles POST /api/v1/members/me/check-in
func (h *Handler) SelfCheckIn(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req selfCheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.service.SelfCheckIn(r.Context(), userID, req.Method, req.QRToken, req.CardUID, req.OrgID)
	if err != nil {
		h.logger.Error("self check-in failed", "error", err, "user_id", userID, "method", req.Method)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ManualCheckIn handles POST /api/v1/orgs/{orgId}/attendance/check-in
func (h *Handler) ManualCheckIn(w http.ResponseWriter, r *http.Request) {
	staffUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req manualCheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.service.ManualCheckIn(r.Context(), staffUserID, req.MemberID, orgID)
	if err != nil {
		h.logger.Error("manual check-in failed",
			"error", err, "staff_id", staffUserID,
			"member_id", req.MemberID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, result)
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

// WeeklySummary handles GET /api/v1/orgs/{orgId}/attendance/weekly
func (h *Handler) WeeklySummary(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	summaries, err := h.service.ListWeeklySummary(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get weekly summary", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": summaries})
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

// GetMyStreaks handles GET /api/v1/members/me/streaks
func (h *Handler) GetMyStreaks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	streaks, err := h.service.GetStreaks(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get streaks", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": streaks})
}

// GetMyCalendar handles GET /api/v1/members/me/attendance/calendar
func (h *Handler) GetMyCalendar(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	month := r.URL.Query().Get("month")
	if month == "" {
		month = time.Now().In(nptLocation).Format("2006-01")
	}

	result, err := h.service.GetMonthlyCalendar(r.Context(), userID, month)
	if err != nil {
		h.logger.Error("failed to get calendar", "error", err, "user_id", userID, "month", month)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
