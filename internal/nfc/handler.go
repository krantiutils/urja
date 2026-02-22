package nfc

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for NFC endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new NFC handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type registerCardRequest struct {
	CardHex string `json:"card_hex"`
}

type assignCardRequest struct {
	UserID string `json:"user_id"`
}

type registerDeviceRequest struct {
	Name             string `json:"name"`
	DeviceIdentifier string `json:"device_identifier"`
	DeviceSecret     string `json:"device_secret"`
}

type deviceCheckInRequest struct {
	CardUID string `json:"card_uid"`
}

type deviceHeartbeatRequest struct {
	UptimeSeconds int64 `json:"uptime_seconds"`
}

// ListCards handles GET /api/v1/orgs/{orgId}/nfc-cards
func (h *Handler) ListCards(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	cards, total, err := h.service.ListCards(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list nfc cards", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  cards,
		"total": total,
	})
}

// RegisterCard handles POST /api/v1/orgs/{orgId}/nfc-cards
func (h *Handler) RegisterCard(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req registerCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	card, err := h.service.RegisterCard(r.Context(), orgID, req.CardHex)
	if err != nil {
		h.logger.Error("failed to register nfc card", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, card)
}

// AssignCard handles PUT /api/v1/orgs/{orgId}/nfc-cards/{id}/assign
func (h *Handler) AssignCard(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing card ID"})
		return
	}

	var req assignCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	card, err := h.service.AssignCard(r.Context(), cardID, orgID, req.UserID)
	if err != nil {
		h.logger.Error("failed to assign nfc card", "error", err, "org_id", orgID, "card_id", cardID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, card)
}

// UnassignCard handles PUT /api/v1/orgs/{orgId}/nfc-cards/{id}/unassign
func (h *Handler) UnassignCard(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing card ID"})
		return
	}

	card, err := h.service.UnassignCard(r.Context(), cardID, orgID)
	if err != nil {
		h.logger.Error("failed to unassign nfc card", "error", err, "org_id", orgID, "card_id", cardID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, card)
}

// ListDevices handles GET /api/v1/orgs/{orgId}/nfc-devices
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	devices, err := h.service.ListDevices(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list nfc devices", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": devices})
}

// RegisterDevice handles POST /api/v1/orgs/{orgId}/nfc-devices
// Admin provides the device_secret from the ESP32's serial output.
func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req registerDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	device, err := h.service.RegisterDevice(r.Context(), orgID, req.Name, req.DeviceIdentifier, req.DeviceSecret)
	if err != nil {
		h.logger.Error("failed to register nfc device", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, device)
}

// DeviceCheckIn handles POST /api/v1/devices/check-in
// No JWT auth — authenticates via X-Device-Key header (the ESP32's secret).
func (h *Handler) DeviceCheckIn(w http.ResponseWriter, r *http.Request) {
	deviceSecret := strings.TrimSpace(r.Header.Get("X-Device-Key"))
	if deviceSecret == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing X-Device-Key header"})
		return
	}

	var req deviceCheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.service.DeviceCheckIn(r.Context(), deviceSecret, req.CardUID)
	if err != nil {
		h.logger.Error("device check-in failed", "error", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// DeviceHeartbeat handles POST /api/v1/devices/heartbeat
// No JWT auth — authenticates via X-Device-Key header.
func (h *Handler) DeviceHeartbeat(w http.ResponseWriter, r *http.Request) {
	deviceSecret := strings.TrimSpace(r.Header.Get("X-Device-Key"))
	if deviceSecret == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing X-Device-Key header"})
		return
	}

	var req deviceHeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.DeviceHeartbeat(r.Context(), deviceSecret, req.UptimeSeconds); err != nil {
		h.logger.Error("device heartbeat failed", "error", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
