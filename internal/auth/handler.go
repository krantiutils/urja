package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/urja-gym/urja/pkg/middleware"
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
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	TokenType           string `json:"token_type"`
	IsNewUser           bool   `json:"is_new_user"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
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

	result, err := h.service.VerifyOTP(r.Context(), req.Phone, req.OTP)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:         result.AccessToken,
		RefreshToken:        result.RefreshToken,
		TokenType:           "Bearer",
		IsNewUser:           result.IsNewUser,
		OnboardingCompleted: result.OnboardingCompleted,
	})
}

// Refresh handles POST /api/v1/auth/refresh — issues new token pair from refresh token.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "refresh_token is required"})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn("token refresh failed", "error", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

// Logout handles POST /api/v1/auth/logout — revokes the provided refresh token.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "authentication required"})
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "refresh_token is required"})
		return
	}

	if err := h.service.Logout(r.Context(), userID, req.RefreshToken); err != nil {
		h.logger.Warn("logout failed", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{Message: "logged out successfully"})
}

// LogoutAll handles POST /api/v1/auth/logout-all — revokes all refresh tokens for the user.
func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "authentication required"})
		return
	}

	if err := h.service.LogoutAll(r.Context(), userID); err != nil {
		h.logger.Error("logout-all failed", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to logout all devices"})
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{Message: "logged out from all devices"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
