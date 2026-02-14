package accounts

import (
	"context"
	"testing"
)

func TestService_Create_Validation(t *testing.T) {
	s := &Service{logger: testLogger()}

	tests := []struct {
		name     string
		category string
		txnDate  string
		txnType  string
		amount   float64
		payType  string
		wantErr  string
	}{
		{
			name:    "missing category",
			txnDate: "2026-01-01", txnType: "income", amount: 100, payType: "cash",
			wantErr: "category is required",
		},
		{
			name:     "missing date",
			category: "Sub",
			txnType:  "income", amount: 100, payType: "cash",
			wantErr: "transaction_date is required",
		},
		{
			name:     "invalid date format",
			category: "Sub", txnDate: "01-01-2026",
			txnType: "income", amount: 100, payType: "cash",
			wantErr: "invalid transaction_date format (expected YYYY-MM-DD)",
		},
		{
			name:     "invalid transaction type",
			category: "Sub", txnDate: "2026-01-01",
			txnType: "refund", amount: 100, payType: "cash",
			wantErr: "invalid transaction_type: refund (must be income or expense)",
		},
		{
			name:     "zero amount",
			category: "Sub", txnDate: "2026-01-01",
			txnType: "income", amount: 0, payType: "cash",
			wantErr: "amount must be positive",
		},
		{
			name:     "negative amount",
			category: "Sub", txnDate: "2026-01-01",
			txnType: "income", amount: -50, payType: "cash",
			wantErr: "amount must be positive",
		},
		{
			name:     "invalid payment type",
			category: "Sub", txnDate: "2026-01-01",
			txnType: "income", amount: 100, payType: "credit_card",
			wantErr: "invalid payment_type: credit_card (must be cash, bank_transfer, or khalti)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Create(context.Background(), "org-1",
				tt.category, "", tt.txnDate, tt.txnType,
				tt.amount, tt.payType, "", "user-1",
			)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestService_Update_Validation(t *testing.T) {
	s := &Service{logger: testLogger()}

	tests := []struct {
		name     string
		category string
		txnDate  string
		txnType  string
		amount   float64
		payType  string
		wantErr  string
	}{
		{
			name:    "missing category",
			txnDate: "2026-01-01", txnType: "income", amount: 100, payType: "cash",
			wantErr: "category is required",
		},
		{
			name:     "invalid date",
			category: "Sub", txnDate: "bad-date",
			txnType: "income", amount: 100, payType: "cash",
			wantErr: "invalid transaction_date format (expected YYYY-MM-DD)",
		},
		{
			name:     "invalid type",
			category: "Sub", txnDate: "2026-01-01",
			txnType: "bonus", amount: 100, payType: "cash",
			wantErr: "invalid transaction_type: bonus (must be income or expense)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Update(context.Background(), "org-1", "txn-1",
				tt.category, "", tt.txnDate, tt.txnType,
				tt.amount, tt.payType, "",
			)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestService_List_InvalidFilters(t *testing.T) {
	s := &Service{logger: testLogger()}

	tests := []struct {
		name    string
		from    string
		to      string
		txnType string
		payType string
		wantErr string
	}{
		{
			name: "invalid from date", from: "not-a-date",
			wantErr: "invalid from date format (expected YYYY-MM-DD)",
		},
		{
			name: "invalid to date", to: "not-a-date",
			wantErr: "invalid to date format (expected YYYY-MM-DD)",
		},
		{
			name: "invalid txn type", txnType: "refund",
			wantErr: "invalid transaction type: refund",
		},
		{
			name: "invalid pay type", payType: "bitcoin",
			wantErr: "invalid payment type: bitcoin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := s.List(context.Background(), "org-1",
				tt.from, tt.to, "", tt.txnType, tt.payType, 10, 0,
			)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestService_GetSummary_InvalidDates(t *testing.T) {
	s := &Service{logger: testLogger()}

	_, err := s.GetSummary(context.Background(), "org-1", "bad-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from date, got nil")
	}

	_, err = s.GetSummary(context.Background(), "org-1", "", "bad-date")
	if err == nil {
		t.Fatal("expected error for invalid to date, got nil")
	}
}

func TestService_Export_InvalidFilters(t *testing.T) {
	s := &Service{logger: testLogger()}

	_, err := s.Export(context.Background(), "org-1", "bad", "", "", "income", "")
	if err == nil {
		t.Fatal("expected error for invalid from date")
	}

	_, err = s.Export(context.Background(), "org-1", "", "", "", "refund", "")
	if err == nil {
		t.Fatal("expected error for invalid txn type")
	}

	_, err = s.Export(context.Background(), "org-1", "", "", "", "", "bitcoin")
	if err == nil {
		t.Fatal("expected error for invalid pay type")
	}
}
