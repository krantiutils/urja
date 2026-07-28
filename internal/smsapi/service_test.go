package smsapi

import (
	"context"
	"strings"
	"testing"
)

func TestBuyCredits_InvalidQuantity(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 0, "cash", "user-1")
	if err == nil {
		t.Fatal("expected error for zero quantity")
	}
	if err.Error() != "quantity must be positive" {
		t.Errorf("error = %q, want %q", err.Error(), "quantity must be positive")
	}
}

// The rate is no longer a parameter — the server prices the purchase — so what
// is worth asserting is that an absurd quantity is refused rather than silently
// billed.
func TestBuyCredits_QuantityCeiling(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", maxCreditsPerPurchase+1, "cash", "user-1")
	if err == nil {
		t.Fatal("expected error for an excessive quantity")
	}
	if !strings.Contains(err.Error(), "or fewer") {
		t.Errorf("error = %q, want a quantity ceiling message", err.Error())
	}
}

func TestBuyCredits_InvalidPaymentMethod(t *testing.T) {
	s := NewService(nil, nil, testLogger())

	_, err := s.BuyCredits(context.Background(), "org-1", 100, "bitcoin", "user-1")
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
