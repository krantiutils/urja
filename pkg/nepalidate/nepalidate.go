// Package nepalidate converts Gregorian dates to Bikram Sambat and derives the
// Nepali fiscal year. Conversion is table-driven: BS month lengths vary by year
// and cannot be computed, so `data.go` carries the reference table.
//
// This package is deliberately Go-only. The API returns both the AD and BS date
// as strings, so the web client never needs a second implementation that could
// drift from this one.
package nepalidate

import (
	"fmt"
	"time"
)

// Date is a Bikram Sambat calendar date.
type Date struct {
	Year  int
	Month int // 1 = Baisakh
	Day   int
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// FiscalYear returns the Nepali fiscal year containing d, as "2082-83".
// The year runs Shrawan 1 (month 4) to the end of Ashad (month 3).
func (d Date) FiscalYear() string {
	startYear := d.Year
	if d.Month < 4 {
		startYear = d.Year - 1
	}
	return fmt.Sprintf("%04d-%02d", startYear, (startYear+1)%100)
}

// FromAD converts a Gregorian date to Bikram Sambat. The clock time and
// location of ad are ignored; only its calendar date is used.
func FromAD(ad time.Time) (Date, error) {
	target := time.Date(ad.Year(), ad.Month(), ad.Day(), 0, 0, 0, 0, time.UTC)
	if target.Before(epochAD) {
		return Date{}, fmt.Errorf("nepalidate: %s is before the supported range", target.Format("2006-01-02"))
	}

	// Days elapsed since 2000-01-01 BS, then walked forward through the table.
	remaining := int(target.Sub(epochAD).Hours() / 24)

	year := epochBSYear
	month := 1
	day := 1

	for remaining > 0 {
		months, ok := daysInMonth[year]
		if !ok {
			return Date{}, fmt.Errorf("nepalidate: BS year %d is outside the supported range", year)
		}
		inMonth := months[month-1]

		if day+remaining <= inMonth {
			day += remaining
			remaining = 0
			break
		}

		remaining -= inMonth - day + 1
		day = 1
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	return Date{Year: year, Month: month, Day: day}, nil
}

// Today returns the current BS date in Kathmandu time, which is the only
// timezone that matters for a Nepali tax document.
func Today() (Date, error) {
	loc, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		return Date{}, fmt.Errorf("nepalidate: loading Kathmandu timezone: %w", err)
	}
	return FromAD(time.Now().In(loc))
}
