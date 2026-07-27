package guide

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for training guide endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new guide handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type createGuideRequest struct {
	Title         string `json:"title"`
	TitleNe       string `json:"title_ne"`
	Content       string `json:"content"`
	ContentNe     string `json:"content_ne"`
	Category      string `json:"category"`
	CoverImageURL string `json:"cover_image_url"`
}

type updateGuideRequest struct {
	Title         string `json:"title"`
	TitleNe       string `json:"title_ne"`
	Content       string `json:"content"`
	ContentNe     string `json:"content_ne"`
	Category      string `json:"category"`
	CoverImageURL string `json:"cover_image_url"`
}

type publishGuideRequest struct {
	IsPublished *bool `json:"is_published"`
}

// ListPublished handles GET /api/v1/training-guides (public, no auth)
func (h *Handler) ListPublished(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	guides, total, err := h.service.ListPublished(r.Context(), category, search, limit, offset)
	if err != nil {
		h.logger.Error("failed to list published training guides", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": guides, "total": total})
}

// GetPublished handles GET /api/v1/training-guides/{id} (public, no auth)
func (h *Handler) GetPublished(w http.ResponseWriter, r *http.Request) {
	guideID := chi.URLParam(r, "id")
	if guideID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing guide ID"})
		return
	}

	g, err := h.service.GetPublishedByID(r.Context(), guideID)
	if err != nil {
		h.logger.Error("failed to get published training guide", "error", err, "guide_id", guideID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "training guide not found"})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// ListAll handles GET /api/v1/orgs/{orgId}/training-guides (admin)
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	guides, total, err := h.service.ListAll(r.Context(), orgID, category, search, limit, offset)
	if err != nil {
		h.logger.Error("failed to list training guides", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": guides, "total": total})
}

// GetAdmin handles GET /api/v1/orgs/{orgId}/training-guides/{id} (admin, includes unpublished)
func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	guideID := chi.URLParam(r, "id")
	if guideID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing guide ID"})
		return
	}

	g, err := h.service.GetByID(r.Context(), orgID, guideID)
	if err != nil {
		h.logger.Error("failed to get training guide", "error", err, "guide_id", guideID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "training guide not found"})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// Create handles POST /api/v1/orgs/{orgId}/training-guides
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Authentication is checked before organization scope: an anonymous caller
	// gets 401, not a 400 that would imply their request was merely malformed.
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

	var req createGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	g, err := h.service.Create(r.Context(), orgID, req.Title, req.TitleNe, req.Content, req.ContentNe,
		req.Category, req.CoverImageURL, userID)
	if err != nil {
		h.logger.Error("failed to create training guide", "error", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, g)
}

// Update handles PUT /api/v1/orgs/{orgId}/training-guides/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	guideID := chi.URLParam(r, "id")
	if guideID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing guide ID"})
		return
	}

	var req updateGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	g, err := h.service.Update(r.Context(), orgID, guideID, req.Title, req.TitleNe, req.Content, req.ContentNe,
		req.Category, req.CoverImageURL)
	if err != nil {
		h.logger.Error("failed to update training guide", "error", err, "guide_id", guideID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// Publish handles PATCH /api/v1/orgs/{orgId}/training-guides/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	guideID := chi.URLParam(r, "id")
	if guideID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing guide ID"})
		return
	}

	var req publishGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.IsPublished == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "is_published is required"})
		return
	}

	g, err := h.service.SetPublished(r.Context(), orgID, guideID, *req.IsPublished)
	if err != nil {
		h.logger.Error("failed to update training guide publish status", "error", err, "guide_id", guideID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "training guide not found"})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// Delete handles DELETE /api/v1/orgs/{orgId}/training-guides/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	guideID := chi.URLParam(r, "id")
	if guideID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing guide ID"})
		return
	}

	if err := h.service.Delete(r.Context(), orgID, guideID); err != nil {
		h.logger.Error("failed to delete training guide", "error", err, "guide_id", guideID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "training guide deleted"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
