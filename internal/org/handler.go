package org

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for organization endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new org handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// List handles GET /api/v1/gyms
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	orgs, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list organizations", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": orgs})
}

// Get handles GET /api/v1/gyms/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing gym ID"})
		return
	}

	org, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "gym not found"})
		return
	}

	writeJSON(w, http.StatusOK, org)
}

// createOrgRequest is the JSON body for POST /api/v1/orgs.
type createOrgRequest struct {
	Name          string          `json:"name"`
	NameNe        string          `json:"name_ne"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description"`
	DescriptionNe string          `json:"description_ne"`
	LogoURL       string          `json:"logo_url"`
	Address       string          `json:"address"`
	AddressNe     string          `json:"address_ne"`
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	Latitude      *float64        `json:"latitude"`
	Longitude     *float64        `json:"longitude"`
	Settings      json.RawMessage `json:"settings"`
}

// Create handles POST /api/v1/orgs
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	in := &CreateOrgInput{
		Name:          req.Name,
		NameNe:        req.NameNe,
		Slug:          req.Slug,
		Description:   req.Description,
		DescriptionNe: req.DescriptionNe,
		LogoURL:       req.LogoURL,
		Address:       req.Address,
		AddressNe:     req.AddressNe,
		Phone:         req.Phone,
		Email:         req.Email,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Settings:      req.Settings,
	}

	org, err := h.service.Create(r.Context(), userID, in)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to create organization", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create organization"})
		return
	}

	writeJSON(w, http.StatusCreated, org)
}

// updateOrgRequest is the JSON body for PUT /api/v1/orgs/:orgId.
type updateOrgRequest struct {
	Name          *string          `json:"name"`
	NameNe        *string          `json:"name_ne"`
	Description   *string          `json:"description"`
	DescriptionNe *string          `json:"description_ne"`
	LogoURL       *string          `json:"logo_url"`
	Address       *string          `json:"address"`
	AddressNe     *string          `json:"address_ne"`
	Phone         *string          `json:"phone"`
	Email         *string          `json:"email"`
	Latitude      *float64         `json:"latitude"`
	Longitude     *float64         `json:"longitude"`
	Settings      *json.RawMessage `json:"settings"`
}

// Update handles PUT /api/v1/orgs/{orgId}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization ID"})
		return
	}

	var req updateOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	in := &UpdateOrgInput{
		Name:          req.Name,
		NameNe:        req.NameNe,
		Description:   req.Description,
		DescriptionNe: req.DescriptionNe,
		LogoURL:       req.LogoURL,
		Address:       req.Address,
		AddressNe:     req.AddressNe,
		Phone:         req.Phone,
		Email:         req.Email,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Settings:      req.Settings,
	}

	org, err := h.service.Update(r.Context(), userID, orgID, in)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update organization", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update organization"})
		return
	}

	writeJSON(w, http.StatusOK, org)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
