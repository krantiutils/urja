package dues

import (
	"context"
	"testing"
)

func TestService_List_DefaultPagination(t *testing.T) {
	// Service.List should clamp limit/offset to safe defaults.
	// We can't test actual DB queries without a real pool, but we can test
	// that validation works by checking invalid statuses.
	s := &Service{logger: testLogger()}

	_, _, err := s.List(context.Background(), "org-1", "", "invalid_status", 10, 0)
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
	if err.Error() != "invalid status filter: invalid_status" {
		t.Errorf("error = %q, want %q", err.Error(), "invalid status filter: invalid_status")
	}
}

func TestService_Pay_Validation(t *testing.T) {
	s := &Service{logger: testLogger()}

	tests := []struct {
		name    string
		amount  float64
		method  string
		wantErr string
	}{
		{
			name:    "zero amount",
			amount:  0,
			method:  "cash",
			wantErr: "amount must be positive",
		},
		{
			name:    "negative amount",
			amount:  -100,
			method:  "cash",
			wantErr: "amount must be positive",
		},
		{
			name:    "invalid method",
			amount:  1000,
			method:  "credit_card",
			wantErr: "invalid payment method: credit_card (must be cash, bank_transfer, or khalti)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Pay(context.Background(), "org-1", "due-1", "user-1", tt.amount, tt.method, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestService_List_InvalidStatus(t *testing.T) {
	s := &Service{logger: testLogger()}

	invalidStatuses := []string{"invalid", "cancelled", "deleted", "active"}
	for _, status := range invalidStatuses {
		t.Run(status, func(t *testing.T) {
			_, _, err := s.List(context.Background(), "org-1", "", status, 10, 0)
			if err == nil {
				t.Fatalf("expected error for invalid status %q, got nil", status)
			}
			expected := "invalid status filter: " + status
			if err.Error() != expected {
				t.Errorf("error = %q, want %q", err.Error(), expected)
			}
		})
	}
}
