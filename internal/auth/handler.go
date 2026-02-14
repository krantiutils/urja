package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Handler handles HTTP requests for authentication endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type loginRequest struct {
	Phone string `json:"phone"`
}

type loginResponse struct {
	Message string `json:"message"`
}

type verifyRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Login handles POST /api/v1/auth/login — requests an OTP.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.Phone == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "phone is required"})
		return
	}

	if err := h.service.RequestOTP(r.Context(), req.Phone); err != nil {
		h.logger.Error("OTP request failed", "error", err)
		writeJSON(w, http.StatusTooManyRequests, errorResponse{Error: err.Error()})
		return
	}

	// SECURITY: Never return the OTP in the response (PRD vulnerability #1)
	writeJSON(w, http.StatusOK, loginResponse{Message: "OTP sent successfully"})
}

// VerifyOTP handles POST /api/v1/auth/verify-otp — verifies OTP and returns tokens.
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.Phone == "" || req.OTP == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "phone and otp are required"})
		return
	}

	accessToken, refreshToken, err := h.service.VerifyOTP(r.Context(), req.Phone, req.OTP)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// At this point headers are already sent; log the error.
		slog.Error("failed to encode JSON response", "error", err)
	}
}
