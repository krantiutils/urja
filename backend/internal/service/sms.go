package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const aakashSMSEndpoint = "https://sms.aakashsms.com/sms/v3/send"

// SMSSender defines the interface for sending SMS messages.
type SMSSender interface {
	SendOTP(ctx context.Context, phone, otp string) error
}

// AakashSMS sends SMS via the Aakash SMS API.
type AakashSMS struct {
	token  string
	client *http.Client
}

// NewAakashSMS creates a new Aakash SMS sender.
func NewAakashSMS(token string) *AakashSMS {
	return &AakashSMS{
		token: token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendOTP sends an OTP code to the given phone number via Aakash SMS.
func (s *AakashSMS) SendOTP(ctx context.Context, phone, otp string) error {
	message := fmt.Sprintf("Your Urja verification code is: %s. Do not share this code. Expires in 5 minutes.", otp)

	form := url.Values{}
	form.Set("auth_token", s.token)
	form.Set("to", phone)
	form.Set("text", message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aakashSMSEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("creating SMS request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending SMS: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Str("phone", phone).
			Msg("aakash SMS API error")
		return fmt.Errorf("SMS API returned status %d", resp.StatusCode)
	}

	log.Info().Str("phone", phone).Msg("OTP SMS sent successfully")
	return nil
}

// DevSMS is a development SMS sender that logs the OTP instead of sending it.
type DevSMS struct{}

// NewDevSMS creates a development SMS sender.
func NewDevSMS() *DevSMS {
	return &DevSMS{}
}

// SendOTP logs the OTP for development purposes.
func (s *DevSMS) SendOTP(_ context.Context, phone, otp string) error {
	log.Warn().
		Str("phone", phone).
		Str("otp", otp).
		Msg("DEV MODE: OTP not sent via SMS, logged here instead")
	return nil
}
