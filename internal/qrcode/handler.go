package qrcode

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	goqrcode "github.com/skip2/go-qrcode"
	"github.com/urja-gym/urja/pkg/middleware"
)

// OrgSlugLookup is an interface for looking up an org's slug by its ID.
type OrgSlugLookup interface {
	GetSlugByID(ctx context.Context, orgID string) (string, error)
}

// Handler handles QR code generation endpoints.
type Handler struct {
	baseURL string
	orgSlug OrgSlugLookup
	logger  *slog.Logger
}

// NewHandler creates a new QR code handler.
// baseURL is the application base URL used to construct check-in URLs.
func NewHandler(baseURL string, orgSlug OrgSlugLookup, logger *slog.Logger) *Handler {
	return &Handler{baseURL: baseURL, orgSlug: orgSlug, logger: logger}
}

// GenerateQR handles GET /api/v1/orgs/{orgId}/qr-code
// Returns a PNG QR code image that members scan to check in.
func (h *Handler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing organization"})
		return
	}

	// Look up the org slug for a human-readable URL
	slug, err := h.orgSlug.GetSlugByID(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to look up org slug", "error", err, "org_id", orgID)
		// Fall back to org ID if slug lookup fails
		slug = orgID
	}

	// The QR code encodes a check-in URL using the org slug
	checkInURL := fmt.Sprintf("%s/checkin?gym=%s", h.baseURL, slug)

	format := r.URL.Query().Get("format")
	size := 256

	if format == "json" {
		// Return the URL as JSON instead of an image
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"org_id":      orgID,
			"slug":        slug,
			"checkin_url": checkInURL,
		})
		return
	}

	// Generate QR code PNG
	png, err := goqrcode.Encode(checkInURL, goqrcode.Medium, size)
	if err != nil {
		h.logger.Error("failed to generate QR code", "error", err, "org_id", orgID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate QR code"})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(png); err != nil {
		h.logger.Error("failed to write QR code response", "error", err)
	}
}
