package subscription

import (
	"testing"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{"zero limit defaults to 20", 0, 0, 20, 0},
		{"negative limit defaults to 20", -5, 0, 20, 0},
		{"over 100 limit defaults to 20", 150, 0, 20, 0},
		{"valid limit preserved", 50, 10, 50, 10},
		{"negative offset defaults to 0", 20, -1, 20, 0},
		{"boundary limit 100", 100, 0, 100, 0},
		{"boundary limit 1", 1, 0, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset := normalizePagination(tt.inputLimit, tt.inputOffset)
			if limit != tt.expectedLimit {
				t.Errorf("limit = %d, want %d", limit, tt.expectedLimit)
			}
			if offset != tt.expectedOffset {
				t.Errorf("offset = %d, want %d", offset, tt.expectedOffset)
			}
		})
	}
}

func TestValidateAssignRequest(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name    string
		req     AssignRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing package_id",
			req:     AssignRequest{StartDate: "2026-01-01", AmountPaid: 1000},
			wantErr: true,
			errMsg:  "package_id is required",
		},
		{
			name:    "missing start_date",
			req:     AssignRequest{PackageID: "pkg-1", AmountPaid: 1000},
			wantErr: true,
			errMsg:  "start_date is required",
		},
		{
			name:    "invalid start_date format",
			req:     AssignRequest{PackageID: "pkg-1", StartDate: "01/01/2026", AmountPaid: 1000},
			wantErr: true,
			errMsg:  "YYYY-MM-DD format",
		},
		{
			name:    "invalid payment method",
			req:     AssignRequest{PackageID: "pkg-1", StartDate: "2026-01-01", PaymentMethod: "bitcoin"},
			wantErr: true,
			errMsg:  "invalid payment_method",
		},
		{
			name:    "negative amount_paid",
			req:     AssignRequest{PackageID: "pkg-1", StartDate: "2026-01-01", AmountPaid: -100},
			wantErr: true,
			errMsg:  "amount_paid cannot be negative",
		},
		{
			name:    "negative discount",
			req:     AssignRequest{PackageID: "pkg-1", StartDate: "2026-01-01", Discount: -50},
			wantErr: true,
			errMsg:  "discount cannot be negative",
		},
		{
			name: "valid request with all fields",
			req: AssignRequest{
				PackageID:     "pkg-1",
				StartDate:     "2026-02-01",
				PaymentMethod: "cash",
				AmountPaid:    3000,
				Discount:      500,
			},
			wantErr: false,
		},
		{
			name: "valid request minimal fields",
			req: AssignRequest{
				PackageID:  "pkg-1",
				StartDate:  "2026-02-01",
				AmountPaid: 3000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.validateAssignRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !containsStr(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateRenewRequest(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name    string
		req     RenewRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing member_package_id",
			req:     RenewRequest{AmountPaid: 1000},
			wantErr: true,
			errMsg:  "member_package_id is required",
		},
		{
			name:    "invalid payment method",
			req:     RenewRequest{MemberPackageID: "mp-1", PaymentMethod: "paypal"},
			wantErr: true,
			errMsg:  "invalid payment_method",
		},
		{
			name:    "negative amount_paid",
			req:     RenewRequest{MemberPackageID: "mp-1", AmountPaid: -100},
			wantErr: true,
			errMsg:  "amount_paid cannot be negative",
		},
		{
			name:    "negative discount",
			req:     RenewRequest{MemberPackageID: "mp-1", Discount: -50},
			wantErr: true,
			errMsg:  "discount cannot be negative",
		},
		{
			name: "valid request",
			req: RenewRequest{
				MemberPackageID: "mp-1",
				PaymentMethod:   "khalti",
				AmountPaid:      3000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.validateRenewRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !containsStr(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateExtendRequest(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name    string
		req     ExtendRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing member_package_id",
			req:     ExtendRequest{ExtraDays: 30, AmountPaid: 1000},
			wantErr: true,
			errMsg:  "member_package_id is required",
		},
		{
			name:    "zero extra_days",
			req:     ExtendRequest{MemberPackageID: "mp-1", ExtraDays: 0},
			wantErr: true,
			errMsg:  "extra_days must be positive",
		},
		{
			name:    "negative extra_days",
			req:     ExtendRequest{MemberPackageID: "mp-1", ExtraDays: -10},
			wantErr: true,
			errMsg:  "extra_days must be positive",
		},
		{
			name:    "invalid payment method",
			req:     ExtendRequest{MemberPackageID: "mp-1", ExtraDays: 30, PaymentMethod: "crypto"},
			wantErr: true,
			errMsg:  "invalid payment_method",
		},
		{
			name:    "negative amount_paid",
			req:     ExtendRequest{MemberPackageID: "mp-1", ExtraDays: 30, AmountPaid: -100},
			wantErr: true,
			errMsg:  "amount_paid cannot be negative",
		},
		{
			name:    "negative discount",
			req:     ExtendRequest{MemberPackageID: "mp-1", ExtraDays: 30, Discount: -50},
			wantErr: true,
			errMsg:  "discount cannot be negative",
		},
		{
			name: "valid request",
			req: ExtendRequest{
				MemberPackageID: "mp-1",
				ExtraDays:       30,
				PaymentMethod:   "bank_transfer",
				AmountPaid:      1500,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.validateExtendRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !containsStr(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidPaymentMethods(t *testing.T) {
	valid := []string{"khalti", "cash", "manual", "bank_transfer"}
	invalid := []string{"paypal", "bitcoin", "stripe", "credit_card", ""}

	for _, m := range valid {
		if !validPaymentMethods[m] {
			t.Errorf("expected %q to be a valid payment method", m)
		}
	}
	for _, m := range invalid {
		if validPaymentMethods[m] {
			t.Errorf("expected %q to be an invalid payment method", m)
		}
	}
}

func TestNilIfEmpty(t *testing.T) {
	result := nilIfEmpty("")
	if result != nil {
		t.Errorf("nilIfEmpty(\"\") = %v, want nil", result)
	}

	result = nilIfEmpty("hello")
	if result == nil {
		t.Fatal("nilIfEmpty(\"hello\") = nil, want non-nil")
	}
	if *result != "hello" {
		t.Errorf("nilIfEmpty(\"hello\") = %q, want %q", *result, "hello")
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
