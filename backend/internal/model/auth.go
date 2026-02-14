package model

// LoginRequest is the payload for POST /api/v1/auth/login.
type LoginRequest struct {
	Phone string `json:"phone"`
}

// VerifyOTPRequest is the payload for POST /api/v1/auth/verify-otp.
type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

// RefreshRequest is the payload for POST /api/v1/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest is the payload for POST /api/v1/auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse is the response for successful authentication.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// MessageResponse is a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}
