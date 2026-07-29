package nepalidate

import (
	"testing"
	"time"
)

// Every BS year is 365 or 366 days. A typo in the table almost always breaks this.
func TestTable_YearLengthsAreSane(t *testing.T) {
	for year, months := range daysInMonth {
		if len(months) != 12 {
			t.Fatalf("BS %d: got %d months, want 12", year, len(months))
		}
		total := 0
		for _, d := range months {
			if d < 28 || d > 33 {
				t.Errorf("BS %d: implausible month length %d", year, d)
			}
			total += d
		}
		if total != 365 && total != 366 {
			t.Errorf("BS %d: year totals %d days, want 365 or 366", year, total)
		}
	}
}

// Anchors cross-checked against the published Nepali calendar. If the table
// drifts, these fail. Do not "fix" an anchor to match the code — fix the table.
func TestFromAD_Anchors(t *testing.T) {
	tests := []struct {
		ad   string
		want string
	}{
		{"2024-04-13", "2081-01-01"}, // Nepali New Year 2081
		{"2025-04-14", "2082-01-01"}, // Nepali New Year 2082
		{"2025-07-17", "2082-04-01"}, // Shrawan 1, 2082 — fiscal year start
	}
	for _, tt := range tests {
		t.Run(tt.ad, func(t *testing.T) {
			ad, err := time.Parse("2006-01-02", tt.ad)
			if err != nil {
				t.Fatalf("bad test input: %v", err)
			}
			got, err := FromAD(ad)
			if err != nil {
				t.Fatalf("FromAD(%s): %v", tt.ad, err)
			}
			if got.String() != tt.want {
				t.Errorf("FromAD(%s) = %s, want %s", tt.ad, got.String(), tt.want)
			}
		})
	}
}

// Consecutive AD days must produce consecutive BS days with no gap or repeat.
func TestFromAD_IsMonotonicAcrossAYear(t *testing.T) {
	d, _ := time.Parse("2006-01-02", "2025-01-01")
	prev, err := FromAD(d)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	for i := 0; i < 365; i++ {
		d = d.AddDate(0, 0, 1)
		cur, err := FromAD(d)
		if err != nil {
			t.Fatalf("FromAD(%s): %v", d.Format("2006-01-02"), err)
		}
		sameMonthNextDay := cur.Year == prev.Year && cur.Month == prev.Month && cur.Day == prev.Day+1
		rolledToNextMonth := cur.Day == 1 && (cur.Month == prev.Month+1 || (cur.Month == 1 && prev.Month == 12))
		if !sameMonthNextDay && !rolledToNextMonth {
			t.Fatalf("non-monotonic: %s -> %s at AD %s", prev, cur, d.Format("2006-01-02"))
		}
		prev = cur
	}
}

func TestDate_FiscalYear(t *testing.T) {
	tests := []struct {
		date Date
		want string
	}{
		{Date{2082, 3, 31}, "2081-82"}, // Ashad — last month of the previous FY
		{Date{2082, 4, 1}, "2082-83"},  // Shrawan 1 — first day of the new FY
		{Date{2082, 12, 30}, "2082-83"},
		{Date{2083, 1, 1}, "2082-83"}, // Baisakh is still the same FY
		{Date{2083, 4, 1}, "2083-84"},
	}
	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			if got := tt.date.FiscalYear(); got != tt.want {
				t.Errorf("FiscalYear(%s) = %s, want %s", tt.date, got, tt.want)
			}
		})
	}
}
