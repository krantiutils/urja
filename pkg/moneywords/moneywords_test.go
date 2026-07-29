package moneywords

import (
	"math"
	"testing"
)

func TestRupees(t *testing.T) {
	tests := []struct {
		amount float64
		want   string
	}{
		{0, "Zero rupees only"},
		{1, "One rupee only"},
		{25, "Twenty five rupees only"},
		{100, "One hundred rupees only"},
		{3500, "Three thousand five hundred rupees only"},
		{100000, "One lakh rupees only"},
		{125000, "One lakh twenty five thousand rupees only"},
		{10000000, "One crore rupees only"},
		{10150000, "One crore one lakh fifty thousand rupees only"},
		{1500.50, "One thousand five hundred rupees and fifty paisa only"},
		{0.25, "Zero rupees and twenty five paisa only"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := Rupees(tt.amount); got != tt.want {
				t.Errorf("Rupees(%.2f) = %q, want %q", tt.amount, got, tt.want)
			}
		})
	}
}

func TestRupeesFixes(t *testing.T) {
	tests := []struct {
		amount float64
		want   string
		desc   string
	}{
		// Finding 1: "Minus zero rupees only" — sign decided after rounding
		{-0.001, "Zero rupees only", "tiny negative rounds to zero"},

		// Finding 2: arab/kharab support
		{1000000000, "One arab rupees only", "1 arab = 1e9"},
		{100000000000, "One kharab rupees only", "1 kharab = 1e11"},

		// Finding 3: out-of-range guard — extreme values return sentinel
		{math.Inf(1), "Invalid amount", "positive infinity"},
		{math.Inf(-1), "Invalid amount", "negative infinity"},
		{math.NaN(), "Invalid amount", "NaN"},

		// Finding 4: epsilon rounding for precision — values stored slightly below
		// their decimal representation now round correctly with epsilon.
		{1.005, "One rupee and one paisa only", "1.005 with epsilon rounds to 101 paisa"},
		{0.995, "One rupee only", "0.995 with epsilon rounds to 100 paisa"},
		{0.005, "Zero rupees and one paisa only", "0.005 with epsilon rounds to 1 paisa"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := Rupees(tt.amount); got != tt.want {
				t.Errorf("Rupees(%.10g) = %q, want %q (%s)", tt.amount, got, tt.want, tt.desc)
			}
		})
	}
}

func TestRupeesRound2Fixes(t *testing.T) {
	tests := []struct {
		amount float64
		want   string
		desc   string
	}{
		// Fix Round 2: Regression fix — kharab tier out-of-range panic
		// Value that would cause under1000(3000) -> ones[30] panic before the fix.
		{3e14, "Amount out of range", "3e14 triggers kharab overflow without tighter guard"},

		// Largest representable value: 999 kharab is 9.99e13 rupees < 1e14 boundary.
		{9.99e13, "Nine hundred ninety nine kharab rupees only", "max value just under 1e14 boundary"},

		// Edge case: exactly at the boundary should also be rejected.
		{1e14, "Amount out of range", "1e14 is exactly at the rejection boundary"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := Rupees(tt.amount); got != tt.want {
				t.Errorf("Rupees(%.2e) = %q, want %q (%s)", tt.amount, got, tt.want, tt.desc)
			}
		})
	}
}

func TestUnder1000Defensive(t *testing.T) {
	// Test that under1000() cannot panic even when called with n > 999.
	// It should return an empty string for out-of-contract input.
	// This is an internal function test ensuring the defensive bounds work.
	tests := []struct {
		n    int
		desc string
	}{
		{1000, "n = 1000, just above contract"},
		{3000, "n = 3000, the original panic case from kharab tier"},
		{9223, "n = 9223, would index ones[92] without defensive check"},
		{-1, "n = -1, negative input"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("under1000(%d) panicked: %v (%s)", tt.n, r, tt.desc)
				}
			}()
			_ = under1000(tt.n)
			// If we reach here, no panic occurred — test passes.
		})
	}
}
