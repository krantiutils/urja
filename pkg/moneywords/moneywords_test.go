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
