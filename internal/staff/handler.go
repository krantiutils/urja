package staff

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for staff endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new staff handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// List handles GET /api/v1/orgs/{orgId}/staff
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	search := r.URL.Query().Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	members, total, err := h.service.ListByOrg(r.Context(), orgID, search, limit, offset)
	if err != nil {
		h.logger.Error("failed to list staff", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  members,
		"total": total,
	})
}

type createStaffRequest struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	StaffRole string `json:"staff_role"`
}

// Create handles POST /api/v1/orgs/{orgId}/staff
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req createStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Phone == "" || req.Name == "" || req.StaffRole == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "phone, name, and staff_role are required"})
		return
	}

	member, err := h.service.Create(r.Context(), orgID, req.Phone, req.Name, req.Email, req.AvatarURL, req.StaffRole)
	if err != nil {
		h.logger.Error("failed to create staff", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, member)
}

// Get handles GET /api/v1/orgs/{orgId}/staff/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing staff member ID"})
		return
	}

	member, err := h.service.GetByID(r.Context(), orgID, id)
	if err != nil {
		h.logger.Error("failed to get staff member", "error", err, "org_id", orgID, "staff_id", id)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "staff member not found"})
		return
	}

	writeJSON(w, http.StatusOK, member)
}

type updateStaffRequest struct {
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	StaffRole *string `json:"staff_role,omitempty"`
}

// Update handles PUT /api/v1/orgs/{orgId}/staff/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing staff member ID"})
		return
	}

	var req updateStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	member, err := h.service.Update(r.Context(), orgID, id, req.Name, req.Email, req.AvatarURL, req.StaffRole)
	if err != nil {
		h.logger.Error("failed to update staff member", "error", err, "org_id", orgID, "staff_id", id)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, member)
}

// Delete handles DELETE /api/v1/orgs/{orgId}/staff/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing staff member ID"})
		return
	}

	if err := h.service.Delete(r.Context(), orgID, id); err != nil {
		h.logger.Error("failed to delete staff member", "error", err, "org_id", orgID, "staff_id", id)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "staff member removed"})
}

// ListAttendance handles GET /api/v1/orgs/{orgId}/staff/attendance
func (h *Handler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	var from, to *time.Time
	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	records, total, err := h.service.ListAttendance(r.Context(), orgID, from, to, limit, offset)
	if err != nil {
		h.logger.Error("failed to list staff attendance", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  records,
		"total": total,
	})
}

type recordAttendanceRequest struct {
	StaffID string `json:"staff_id"`
	Action  string `json:"action"` // check_in or check_out
	Method  string `json:"method"` // qr, nfc, manual (required for check_in)
}

// RecordAttendance handles POST /api/v1/orgs/{orgId}/staff/attendance
func (h *Handler) RecordAttendance(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req recordAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.StaffID == "" || req.Action == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "staff_id and action are required"})
		return
	}

	rec, err := h.service.RecordAttendance(r.Context(), orgID, req.StaffID, req.Action, req.Method)
	if err != nil {
		h.logger.Error("failed to record staff attendance", "error", err, "org_id", orgID, "staff_id", req.StaffID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, rec)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
