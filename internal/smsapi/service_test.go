package smsapi

import (
	"context"
	"testing"
)

func TestBuyCredits_InvalidQuantity(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 0, 1.5, 150, "cash", "user-1")
	if err == nil {
		t.Fatal("expected error for zero quantity")
	}
	if err.Error() != "quantity must be positive" {
		t.Errorf("error = %q, want %q", err.Error(), "quantity must be positive")
	}
}

func TestBuyCredits_InvalidRate(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 100, 0, 0, "cash", "user-1")
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
	if err.Error() != "rate must be positive" {
		t.Errorf("error = %q, want %q", err.Error(), "rate must be positive")
	}
}

func TestBuyCredits_InvalidAmount(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 100, 1.5, 0, "cash", "user-1")
	if err == nil {
		t.Fatal("expected error for zero amount")
	}
	if err.Error() != "amount must be positive" {
		t.Errorf("error = %q, want %q", err.Error(), "amount must be positive")
	}
}

func TestBuyCredits_InvalidPaymentMethod(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 100, 1.5, 150, "bitcoin", "user-1")
	if err == nil {
		t.Fatal("expected error for invalid payment method")
	}
	if err.Error() != "invalid payment method: bitcoin" {
		t.Errorf("error = %q, want %q", err.Error(), "invalid payment method: bitcoin")
	}
}

func TestSendSMS_EmptyMessage(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.SendSMS(context.Background(), "org-1", "", "user-1", nil)
	if err == nil {
		t.Fatal("expected error for empty message")
	}
	if err.Error() != "message is required" {
		t.Errorf("error = %q, want %q", err.Error(), "message is required")
	}
}

func TestSendSMS_WhitespaceMessage(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.SendSMS(context.Background(), "org-1", "   ", "user-1", nil)
	if err == nil {
		t.Fatal("expected error for whitespace-only message")
	}
}
