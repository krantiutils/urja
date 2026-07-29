package moneywords

import "testing"

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
