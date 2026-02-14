package notice

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for notice endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new notice handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type createNoticeRequest struct {
	Title     string `json:"title"`
	TitleNe   string `json:"title_ne"`
	Content   string `json:"content"`
	ContentNe string `json:"content_ne"`
}

type updateNoticeRequest struct {
	Title     string `json:"title"`
	TitleNe   string `json:"title_ne"`
	Content   string `json:"content"`
	ContentNe string `json:"content_ne"`
	IsActive  *bool  `json:"is_active"`
}

// List handles GET /api/v1/orgs/{orgId}/notices
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	search := r.URL.Query().Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	notices, err := h.service.List(r.Context(), orgID, search, limit, offset)
	if err != nil {
		h.logger.Error("failed to list notices", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": notices})
}

// Get handles GET /api/v1/orgs/{orgId}/notices/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	noticeID := chi.URLParam(r, "id")
	if noticeID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing notice ID"})
		return
	}

	n, err := h.service.GetByID(r.Context(), orgID, noticeID)
	if err != nil {
		h.logger.Error("failed to get notice", "error", err, "notice_id", noticeID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "notice not found"})
		return
	}

	writeJSON(w, http.StatusOK, n)
}

// Create handles POST /api/v1/orgs/{orgId}/notices
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req createNoticeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	n, err := h.service.Create(r.Context(), orgID, req.Title, req.TitleNe, req.Content, req.ContentNe, userID)
	if err != nil {
		h.logger.Error("failed to create notice", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, n)
}

// Update handles PUT /api/v1/orgs/{orgId}/notices/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	noticeID := chi.URLParam(r, "id")
	if noticeID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing notice ID"})
		return
	}

	var req updateNoticeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	n, err := h.service.Update(r.Context(), orgID, noticeID, req.Title, req.TitleNe, req.Content, req.ContentNe, isActive)
	if err != nil {
		h.logger.Error("failed to update notice", "error", err, "notice_id", noticeID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, n)
}

// Delete handles DELETE /api/v1/orgs/{orgId}/notices/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	noticeID := chi.URLParam(r, "id")
	if noticeID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing notice ID"})
		return
	}

	if err := h.service.Delete(r.Context(), orgID, noticeID); err != nil {
		h.logger.Error("failed to delete notice", "error", err, "notice_id", noticeID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "notice deleted"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
