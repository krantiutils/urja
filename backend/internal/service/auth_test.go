package service

import (
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 100; i++ {
		otp, err := generateOTP(6)
		if err != nil {
			t.Fatalf("generateOTP() error = %v", err)
		}

		if len(otp) != 6 {
			t.Errorf("OTP length = %d, want 6: %q", len(otp), otp)
		}

		for _, r := range otp {
			if r < '0' || r > '9' {
				t.Errorf("OTP contains non-digit: %q", otp)
				break
			}
		}

		seen[otp] = true
	}

	// With 100 random 6-digit OTPs, we should see at least a few distinct values.
	// This is a sanity check that the RNG isn't broken, not a statistical test.
	if len(seen) < 10 {
		t.Errorf("expected at least 10 distinct OTPs in 100 generations, got %d", len(seen))
	}
}

func TestGenerateOTPPadding(t *testing.T) {
	// Generate many OTPs and check they're all 6 digits (zero-padded)
	for i := 0; i < 1000; i++ {
		otp, err := generateOTP(6)
		if err != nil {
			t.Fatalf("generateOTP() error = %v", err)
		}
		if len(otp) != 6 {
			t.Fatalf("OTP length = %d, want 6: %q", len(otp), otp)
		}
	}
}
