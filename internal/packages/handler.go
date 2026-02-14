package packages

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for package endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new packages handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// ListPublic handles GET /api/v1/packages — public listing of active packages.
func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	pkgs, err := h.service.ListPublic(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list public packages", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": pkgs})
}

// ListByOrg handles GET /api/v1/orgs/{orgId}/packages — admin lists org packages.
func (h *Handler) ListByOrg(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	pkgs, err := h.service.ListByOrg(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list org packages", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": pkgs})
}

// Create handles POST /api/v1/orgs/{orgId}/packages — admin creates a package.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var input CreatePackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	pkg, err := h.service.Create(r.Context(), orgID, input)
	if err != nil {
		h.logger.Error("failed to create package", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, pkg)
}

// Update handles PUT /api/v1/orgs/{orgId}/packages/{id} — admin updates a package.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	pkgID := chi.URLParam(r, "id")
	if pkgID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing package ID"})
		return
	}

	var input UpdatePackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	pkg, err := h.service.Update(r.Context(), pkgID, orgID, input)
	if err != nil {
		h.logger.Error("failed to update package", "error", err, "id", pkgID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, pkg)
}

// Delete handles DELETE /api/v1/orgs/{orgId}/packages/{id} — admin soft-deletes a package.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	pkgID := chi.URLParam(r, "id")
	if pkgID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing package ID"})
		return
	}

	if err := h.service.Delete(r.Context(), pkgID, orgID); err != nil {
		h.logger.Error("failed to delete package", "error", err, "id", pkgID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "package deleted"})
}

// ListMyPackages handles GET /api/v1/members/me/packages — member's subscriptions.
func (h *Handler) ListMyPackages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	subs, err := h.service.ListMemberPackages(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list member packages", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": subs})
}

// Subscribe handles POST /api/v1/members/me/packages/{id}/subscribe — initiate Khalti payment.
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	packageID := chi.URLParam(r, "id")
	if packageID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing package ID"})
		return
	}

	result, err := h.service.Subscribe(r.Context(), userID, packageID)
	if err != nil {
		h.logger.Error("failed to subscribe", "error", err, "user_id", userID, "package_id", packageID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

type verifyPaymentRequest struct {
	Pidx string `json:"pidx"`
}

// VerifyPayment handles POST /api/v1/members/me/packages/verify-payment — verify Khalti and activate.
func (h *Handler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req verifyPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Pidx == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "pidx is required"})
		return
	}

	result, err := h.service.VerifyPayment(r.Context(), userID, req.Pidx)
	if err != nil {
		h.logger.Error("payment verification failed", "error", err, "user_id", userID, "pidx", req.Pidx)
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
