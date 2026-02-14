package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/urja-gym/backend/internal/middleware"
	"github.com/urja-gym/backend/internal/model"
	"github.com/urja-gym/backend/internal/service"
)

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Login handles POST /api/v1/auth/login — sends an OTP to the given phone number.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	if req.Phone == "" {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Phone number is required",
		})
		return
	}

	err := h.auth.SendOTP(r.Context(), req.Phone)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPhone):
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
				Error:   "invalid_phone",
				Message: "Invalid Nepal phone number. Must be 10 digits starting with 9, 6, 7, or 8.",
			})
		case errors.Is(err, service.ErrOTPCooldown):
			writeJSON(w, http.StatusTooManyRequests, model.ErrorResponse{
				Error:   "otp_cooldown",
				Message: "Please wait 60 seconds before requesting another OTP.",
			})
		case errors.Is(err, service.ErrOTPRateLimit):
			writeJSON(w, http.StatusTooManyRequests, model.ErrorResponse{
				Error:   "otp_rate_limit",
				Message: "Too many OTP requests. Maximum 5 per hour.",
			})
		default:
			log.Error().Err(err).Str("phone", req.Phone).Msg("failed to send OTP")
			writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to send OTP. Please try again.",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, model.MessageResponse{
		Message: "OTP sent successfully",
	})
}

// VerifyOTP handles POST /api/v1/auth/verify-otp — verifies OTP and issues tokens.
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req model.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	if req.Phone == "" || req.OTP == "" {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Phone number and OTP are required",
		})
		return
	}

	if len(req.OTP) != 6 {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "OTP must be 6 digits",
		})
		return
	}

	resp, err := h.auth.VerifyOTP(r.Context(), req.Phone, req.OTP)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPhone):
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
				Error:   "invalid_phone",
				Message: "Invalid Nepal phone number.",
			})
		case errors.Is(err, service.ErrOTPExpired):
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
				Error:   "otp_expired",
				Message: "OTP has expired. Please request a new one.",
			})
		case errors.Is(err, service.ErrOTPMaxAttempts):
			writeJSON(w, http.StatusTooManyRequests, model.ErrorResponse{
				Error:   "otp_max_attempts",
				Message: "Too many failed attempts. Please request a new OTP.",
			})
		case errors.Is(err, service.ErrOTPInvalid):
			writeJSON(w, http.StatusUnauthorized, model.ErrorResponse{
				Error:   "otp_invalid",
				Message: "Invalid OTP code.",
			})
		case errors.Is(err, service.ErrUserDeactivated):
			writeJSON(w, http.StatusForbidden, model.ErrorResponse{
				Error:   "account_deactivated",
				Message: "Your account has been deactivated.",
			})
		default:
			log.Error().Err(err).Str("phone", req.Phone).Msg("OTP verification failed")
			writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
				Error:   "internal_error",
				Message: "Verification failed. Please try again.",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Refresh handles POST /api/v1/auth/refresh — refreshes an access token.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
		return
	}

	resp, err := h.auth.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefresh), errors.Is(err, service.ErrSessionNotFound):
			writeJSON(w, http.StatusUnauthorized, model.ErrorResponse{
				Error:   "invalid_refresh_token",
				Message: "Invalid or expired refresh token. Please log in again.",
			})
		case errors.Is(err, service.ErrUserDeactivated):
			writeJSON(w, http.StatusForbidden, model.ErrorResponse{
				Error:   "account_deactivated",
				Message: "Your account has been deactivated.",
			})
		default:
			log.Error().Err(err).Msg("token refresh failed")
			writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
				Error:   "internal_error",
				Message: "Token refresh failed. Please try again.",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Logout handles POST /api/v1/auth/logout — invalidates a single session.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req model.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
		return
	}

	if err := h.auth.Logout(r.Context(), req.RefreshToken); err != nil {
		log.Error().Err(err).Msg("logout failed")
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Logout failed. Please try again.",
		})
		return
	}

	writeJSON(w, http.StatusOK, model.MessageResponse{
		Message: "Logged out successfully",
	})
}

// LogoutAll handles POST /api/v1/auth/logout-all — invalidates all sessions.
// Requires a valid access token (protected by auth middleware).
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, model.ErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		return
	}

	if err := h.auth.LogoutAll(r.Context(), claims.Subject); err != nil {
		log.Error().Err(err).Str("user_id", claims.Subject).Msg("logout-all failed")
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to logout all devices. Please try again.",
		})
		return
	}

	writeJSON(w, http.StatusOK, model.MessageResponse{
		Message: "Logged out from all devices successfully",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error().Err(err).Msg("failed to encode JSON response")
	}
}
