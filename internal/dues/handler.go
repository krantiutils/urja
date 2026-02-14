package dues

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for dues endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new dues handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// List handles GET /api/v1/orgs/{orgId}/dues
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	dues, total, err := h.service.List(r.Context(), orgID, search, status, limit, offset)
	if err != nil {
		h.logger.Error("failed to list dues", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  dues,
		"total": total,
	})
}

type payRequest struct {
	Amount           float64 `json:"amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentReference string  `json:"payment_reference,omitempty"`
	DueID            string  `json:"due_id"`
}

// Pay handles POST /api/v1/orgs/{orgId}/dues/{memberId}/pay
func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
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

	var req payRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.DueID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "due_id is required"})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "amount must be positive"})
		return
	}

	if req.PaymentMethod == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "payment_method is required"})
		return
	}

	rec, err := h.service.Pay(r.Context(), orgID, req.DueID, userID, req.Amount, req.PaymentMethod, req.PaymentReference)
	if err != nil {
		h.logger.Error("payment failed", "error", err, "org_id", orgID, "due_id", req.DueID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

// BlockAccess handles POST /api/v1/orgs/{orgId}/dues/{memberId}/block-access
func (h *Handler) BlockAccess(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "member ID is required"})
		return
	}

	count, err := h.service.BlockAccess(r.Context(), orgID, memberID)
	if err != nil {
		h.logger.Error("block access failed", "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to block access"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":           "NFC access blocked",
		"cards_deactivated": count,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
