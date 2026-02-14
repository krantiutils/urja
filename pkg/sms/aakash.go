package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var nepalPhoneRegex = regexp.MustCompile(`^9[6-9]\d{8}$`)

// Client sends SMS messages via the Aakash SMS API.
type Client struct {
	token      string
	apiURL     string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a new Aakash SMS client.
func NewClient(token, apiURL string, logger *slog.Logger) *Client {
	return &Client{
		token:  token,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// aakashResponse represents the response from the Aakash SMS API.
type aakashResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Send sends an arbitrary SMS message to the given phone number.
// Phone must be a valid Nepal number (10 digits, starts with 9[6-9]).
func (c *Client) Send(ctx context.Context, phone, message string) error {
	if !ValidatePhone(phone) {
		return fmt.Errorf("invalid Nepal phone number: %s", phone)
	}
	return c.send(ctx, phone, message)
}

// SendOTP sends an OTP message to the given phone number.
// Phone must be a valid Nepal number (10 digits, starts with 9[6-9]).
func (c *Client) SendOTP(ctx context.Context, phone, otp string) error {
	if !ValidatePhone(phone) {
		return fmt.Errorf("invalid Nepal phone number: %s", phone)
	}

	message := fmt.Sprintf("Your Urja verification code is: %s. Valid for 5 minutes. Do not share this code.", otp)
	return c.send(ctx, phone, message)
}

func (c *Client) send(ctx context.Context, phone, message string) error {
	form := url.Values{
		"auth_token": {c.token},
		"to":         {phone},
		"text":       {message},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("creating SMS request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending SMS: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading SMS response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("SMS API error",
			"status", resp.StatusCode,
			"body", string(body),
			"phone", phone,
		)
		return fmt.Errorf("SMS API returned status %d", resp.StatusCode)
	}

	var apiResp aakashResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("parsing SMS response: %w", err)
	}

	if apiResp.Error != "" {
		return fmt.Errorf("SMS API error: %s", apiResp.Error)
	}

	c.logger.Info("SMS sent successfully",
		"phone", phone[:4]+"******",
	)

	return nil
}

// NormalizePhone strips the +977 prefix from Nepal phone numbers and trims whitespace.
// Handles formats like "+9779801234567", "9779801234567", "+977 9801234567", " 9801234567 ".
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+977")
	phone = strings.TrimPrefix(phone, "977")
	phone = strings.TrimSpace(phone)
	return phone
}

// ValidatePhone checks if a phone number is a valid Nepal mobile number.
// Expects a normalized 10-digit number starting with 9[6-9].
func ValidatePhone(phone string) bool {
	return nepalPhoneRegex.MatchString(phone)
}
