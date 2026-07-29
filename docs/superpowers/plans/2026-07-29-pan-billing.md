# PAN Billing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let a gym issue a PAN tax invoice to a customer, and correct one lawfully by cancelling it or raising a credit note.

**Architecture:** A new `internal/invoice` domain following the repository/service/handler/routes layering used by every other module. Invoices are immutable once issued — enforced by a database trigger, not convention — and numbered from a locked counter row rather than a Postgres sequence so the run of numbers can never contain a gap. Two pure helper packages (`pkg/nepalidate`, `pkg/moneywords`) supply the Bikram Sambat date and the amount-in-words line that a Nepali bill must carry.

**Tech Stack:** Go 1.22, chi v5, pgx v5, Postgres, golang-migrate; Next.js 14 App Router, TypeScript, Tailwind; Playwright for web tests, stdlib `testing` for Go.

## Global Constraints

- **Spec:** `docs/superpowers/specs/2026-07-29-pan-billing-design.md`. Read it before Task 1.
- **Do not touch `internal/billing`.** That is the gym's own SaaS subscription to Urja. The new domain is `internal/invoice`.
- **No edit path.** There is no `PUT /invoices/{id}` and no draft state. An issued bill is immutable; corrections are cancel or credit note only.
- **PAN-only.** `vat_rate` and `vat_amount` ship as columns fixed at `0`. No VAT UI, no VAT arithmetic. VAT is a later phase.
- **Migration numbering:** the next free number is `000053`. Files go in `db/migrations/` as `000053_pan_billing.up.sql` and `.down.sql`.
- **Test commands:** unit `go test ./internal/... ./pkg/...`; e2e `go test ./tests/e2e/...` (needs Postgres on `localhost:5433`, db `urja_test`, user `urja`, password `urja_secret`, and Redis on `localhost:6379`); web `cd web && npx playwright test`.
- **Cross-tenant reads return 404, not 403.** A 403 confirms the row exists.
- **i18n:** every new user-facing string goes in both `en` and `ne` in `web/src/lib/i18n.ts`. `web/tests/i18n-parity.spec.ts` fails otherwise.
- **Money columns** are `DECIMAL(12,2)`; Go side uses `float64` to match the existing `dues`/`transactions` code.
- **Commit style:** imperative subject, explain *why* in the body. End with `Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>`.

## File Structure

| File | Responsibility |
|------|----------------|
| `pkg/nepalidate/nepalidate.go` | AD↔BS conversion, fiscal-year derivation. Pure, no I/O. |
| `pkg/nepalidate/data.go` | The days-per-month lookup table. Data only. |
| `pkg/moneywords/moneywords.go` | Amount → English words in lakh/crore. Pure, no I/O. |
| `db/migrations/000053_pan_billing.up.sql` | Tax columns, four tables, immutability triggers. |
| `internal/invoice/models.go` | Structs, sentinel errors, input types. |
| `internal/invoice/repository.go` | All SQL. Counter allocation, issue, list, get, cancel, credit note, print. |
| `internal/invoice/service.go` | Validation, PAN guard, totals arithmetic, logging. |
| `internal/invoice/handler.go` | HTTP decode/encode, error→status mapping. |
| `internal/invoice/routes.go` | Route table + role gate. |
| `web/src/lib/api.ts` | Client methods (modify). |
| `web/src/app/[lang]/dashboard/invoices/page.tsx` | List. |
| `web/src/app/[lang]/dashboard/invoices/new/page.tsx` | Issue form. |
| `web/src/app/[lang]/dashboard/invoices/[id]/page.tsx` | Document view + actions. |
| `web/src/components/invoice/InvoiceDocument.tsx` | The printable bill. Shared by detail and print. |

---

## Task 1: `pkg/nepalidate` — Bikram Sambat dates

**Files:**
- Create: `pkg/nepalidate/data.go`
- Create: `pkg/nepalidate/nepalidate.go`
- Test: `pkg/nepalidate/nepalidate_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces: `nepalidate.FromAD(time.Time) (Date, error)`, `nepalidate.Date{Year, Month, Day int}`, `(Date) String() string` → `"2082-04-14"`, `(Date) FiscalYear() string` → `"2082-83"`.

> **Read this before you start.** The Bikram Sambat calendar is a lookup table, not
> an algorithm: the number of days in each month varies year to year and cannot be
> computed. A wrong entry silently misdates tax documents, which is a compliance
> defect that no other test in this plan will catch. Step 1 exists to catch it.

- [ ] **Step 1: Write the verification tests first**

These are the safety net for the table. Two independent checks: every year's months must sum to a real year length, and known anchor dates must convert exactly.

```go
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
```

- [ ] **Step 2: Write the lookup table**

Create `pkg/nepalidate/data.go`. `epochAD` is the AD date corresponding to 2000-01-01 BS, the standard epoch used by Nepali calendar implementations.

```go
package nepalidate

import "time"

// epochBSYear is the first year present in daysInMonth.
const epochBSYear = 2000

// epochAD is 2000-01-01 BS expressed in AD.
var epochAD = time.Date(1943, time.April, 14, 0, 0, 0, 0, time.UTC)

// daysInMonth[bsYear] = days in each of the 12 BS months, Baisakh first.
//
// This is reference data, not derived. Cross-check any change against a
// published Nepali calendar; TestTable_YearLengthsAreSane and
// TestFromAD_Anchors are the guards.
var daysInMonth = map[int][12]int{
	2000: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2001: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2002: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2003: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2004: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2005: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2006: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2007: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2008: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31},
	2009: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2010: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2011: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2012: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30},
	2013: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2014: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2015: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2016: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30},
	2017: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2018: {31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2019: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2020: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2021: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2022: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2023: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2024: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2025: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2026: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2027: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2028: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2029: {31, 31, 32, 31, 32, 30, 30, 29, 30, 29, 30, 30},
	2030: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2031: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2032: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2033: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2034: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2035: {30, 32, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31},
	2036: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2037: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2038: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2039: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30},
	2040: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2041: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2042: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2043: {31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30},
	2044: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2045: {31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2046: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2047: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2048: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2049: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2050: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2051: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2052: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2053: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2054: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2055: {31, 32, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30},
	2056: {31, 31, 32, 31, 32, 30, 30, 29, 30, 29, 30, 30},
	2057: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2058: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2059: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2060: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2061: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2062: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2063: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2064: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2065: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31},
	2066: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2067: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2068: {31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30},
	2069: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2070: {31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2071: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2072: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2073: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2074: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2075: {31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30},
	2076: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2077: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2078: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2079: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31},
	2080: {31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2081: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2082: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30},
	2083: {31, 31, 32, 31, 31, 30, 30, 30, 29, 30, 29, 31},
	2084: {31, 31, 32, 31, 31, 30, 30, 30, 29, 30, 29, 31},
	2085: {31, 32, 31, 32, 30, 31, 30, 30, 29, 30, 29, 31},
	2086: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2087: {31, 31, 32, 31, 31, 31, 30, 30, 29, 30, 29, 31},
	2088: {30, 31, 32, 32, 30, 31, 30, 30, 29, 30, 29, 31},
	2089: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
	2090: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31},
}
```

- [ ] **Step 3: Run the table tests**

Run: `go test ./pkg/nepalidate/ -run 'TestTable_YearLengthsAreSane' -v`
Expected: PASS. If a year fails, that row is wrong — correct it against a published calendar before continuing. Do not relax the assertion.

- [ ] **Step 4: Implement conversion**

Create `pkg/nepalidate/nepalidate.go`.

```go
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
```

- [ ] **Step 5: Add the fiscal-year test**

Append to `pkg/nepalidate/nepalidate_test.go`:

```go
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
```

- [ ] **Step 6: Run all package tests**

Run: `go test ./pkg/nepalidate/ -v`
Expected: PASS, all four test functions.

- [ ] **Step 7: Commit**

```bash
git add pkg/nepalidate/
git commit -m "$(cat <<'EOF'
feat: add Bikram Sambat date conversion

A Nepali tax invoice carries a BS date and belongs to a fiscal year running
Shrawan to Ashad, neither of which can be derived from a Gregorian date
arithmetically — BS month lengths vary year to year, so the conversion is
table-driven reference data.

The table is guarded by two independent tests: every year must total 365 or
366 days, and known anchor dates must convert exactly. A silent typo here
would misdate tax documents, and nothing downstream would catch it.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: `pkg/moneywords` — amount in words

**Files:**
- Create: `pkg/moneywords/moneywords.go`
- Test: `pkg/moneywords/moneywords_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces: `moneywords.Rupees(amount float64) string` → `"Three thousand five hundred rupees only"`.

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/moneywords/ -v`
Expected: FAIL — `undefined: Rupees`.

- [ ] **Step 3: Implement**

```go
// Package moneywords renders an amount as English words using the Nepali
// numbering system (crore, lakh, thousand) for the "amount in words" line
// every Nepali tax invoice must carry.
package moneywords

import (
	"fmt"
	"math"
	"strings"
)

var ones = []string{
	"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
	"seventeen", "eighteen", "nineteen",
}

var tens = []string{
	"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
}

// under1000 renders 0..999. Returns "" for 0 so callers can omit empty groups.
func under1000(n int) string {
	if n == 0 {
		return ""
	}
	var parts []string
	if h := n / 100; h > 0 {
		parts = append(parts, ones[h], "hundred")
	}
	r := n % 100
	switch {
	case r == 0:
	case r < 20:
		parts = append(parts, ones[r])
	default:
		parts = append(parts, tens[r/10])
		if r%10 > 0 {
			parts = append(parts, ones[r%10])
		}
	}
	return strings.Join(parts, " ")
}

// words renders a whole number using crore/lakh/thousand grouping.
func words(n int64) string {
	if n == 0 {
		return "zero"
	}
	var parts []string

	if crore := n / 10000000; crore > 0 {
		parts = append(parts, words(crore), "crore")
		n %= 10000000
	}
	if lakh := n / 100000; lakh > 0 {
		parts = append(parts, under1000(int(lakh)), "lakh")
		n %= 100000
	}
	if thousand := n / 1000; thousand > 0 {
		parts = append(parts, under1000(int(thousand)), "thousand")
		n %= 1000
	}
	if rest := under1000(int(n)); rest != "" {
		parts = append(parts, rest)
	}

	return strings.Join(parts, " ")
}

// capitalise upper-cases the first letter, leaving the rest untouched.
func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Rupees renders amount as words, e.g. "Three thousand five hundred rupees only".
// Paisa are included only when non-zero. The amount is rounded to two decimals
// first, so 1500.505 and 1500.51 render identically.
func Rupees(amount float64) string {
	if amount < 0 {
		return "Minus " + strings.ToLower(Rupees(-amount))
	}

	total := int64(math.Round(amount * 100))
	rupees := total / 100
	paisa := total % 100

	unit := "rupees"
	if rupees == 1 {
		unit = "rupee"
	}

	out := fmt.Sprintf("%s %s", words(rupees), unit)
	if paisa > 0 {
		out = fmt.Sprintf("%s and %s paisa", out, words(paisa))
	}
	return capitalise(out) + " only"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./pkg/moneywords/ -v`
Expected: PASS, all 11 cases.

- [ ] **Step 5: Commit**

```bash
git add pkg/moneywords/
git commit -m "$(cat <<'EOF'
feat: render amounts as words in Nepali numbering

Every Nepali tax invoice carries an "amount in words" line, and it groups by
crore and lakh rather than million and billion, so the standard English
rendering is wrong here.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: Migration — schema and immutability

**Files:**
- Create: `db/migrations/000053_pan_billing.up.sql`
- Create: `db/migrations/000053_pan_billing.down.sql`
- Modify: `tests/e2e/setup_test.go` (add new tables to `cleanupTables`)
- Test: `tests/e2e/invoice_schema_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces: tables `invoices`, `invoice_items`, `invoice_prints`, `invoice_counters`; columns `organizations.pan_number`, `.is_vat_registered`, `.tax_legal_name`, `.tax_address`.

- [ ] **Step 1: Write the up migration**

Copy the DDL verbatim from the spec's "Data model" section into
`db/migrations/000053_pan_billing.up.sql`, in this order: the `ALTER TABLE
organizations` block, `invoices`, its indexes (including
`idx_invoices_one_per_transaction`), `invoice_items`, `invoice_prints`,
`invoice_counters`, then the triggers below.

Append the immutability guards:

```sql
-- Issued invoices are immutable. This is the guarantee the whole feature rests
-- on, so it lives in the database rather than only in the service layer.
CREATE OR REPLACE FUNCTION invoices_guard_immutable() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION 'invoice %: issued invoices cannot be deleted', OLD.invoice_number;
    END IF;

    -- Only cancellation fields, print_count and synced_at may ever change.
    -- NOTE: this comparison is positional. A new column added to `invoices`
    -- must be added here too, or it silently becomes mutable.
    IF (NEW.organization_id, NEW.fiscal_year, NEW.sequence, NEW.invoice_number,
        NEW.doc_type, NEW.credit_note_for,
        NEW.seller_name, NEW.seller_pan, NEW.seller_address, NEW.seller_vat_registered,
        NEW.customer_user_id, NEW.customer_name, NEW.customer_pan,
        NEW.customer_address, NEW.customer_phone,
        NEW.issued_date, NEW.issued_date_bs,
        NEW.subtotal, NEW.discount, NEW.taxable_amount,
        NEW.vat_rate, NEW.vat_amount, NEW.total, NEW.amount_in_words,
        NEW.payment_method, NEW.transaction_id, NEW.member_package_id,
        NEW.issued_by, NEW.created_at)
       IS DISTINCT FROM
       (OLD.organization_id, OLD.fiscal_year, OLD.sequence, OLD.invoice_number,
        OLD.doc_type, OLD.credit_note_for,
        OLD.seller_name, OLD.seller_pan, OLD.seller_address, OLD.seller_vat_registered,
        OLD.customer_user_id, OLD.customer_name, OLD.customer_pan,
        OLD.customer_address, OLD.customer_phone,
        OLD.issued_date, OLD.issued_date_bs,
        OLD.subtotal, OLD.discount, OLD.taxable_amount,
        OLD.vat_rate, OLD.vat_amount, OLD.total, OLD.amount_in_words,
        OLD.payment_method, OLD.transaction_id, OLD.member_package_id,
        OLD.issued_by, OLD.created_at)
    THEN
        RAISE EXCEPTION 'invoice %: only cancellation may change an issued invoice',
                        OLD.invoice_number;
    END IF;

    IF OLD.status = 'cancelled' AND NEW.status <> 'cancelled' THEN
        RAISE EXCEPTION 'invoice %: cancellation cannot be undone', OLD.invoice_number;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER invoices_immutable
    BEFORE UPDATE OR DELETE ON invoices
    FOR EACH ROW EXECUTE FUNCTION invoices_guard_immutable();

-- Line items are insert-only.
CREATE OR REPLACE FUNCTION invoice_items_guard_immutable() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'invoice line items are immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER invoice_items_immutable
    BEFORE UPDATE OR DELETE ON invoice_items
    FOR EACH ROW EXECUTE FUNCTION invoice_items_guard_immutable();
```

- [ ] **Step 2: Write the down migration**

```sql
DROP TRIGGER IF EXISTS invoice_items_immutable ON invoice_items;
DROP TRIGGER IF EXISTS invoices_immutable ON invoices;
DROP FUNCTION IF EXISTS invoice_items_guard_immutable();
DROP FUNCTION IF EXISTS invoices_guard_immutable();

DROP TABLE IF EXISTS invoice_prints;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoice_counters;
DROP TABLE IF EXISTS invoices;

ALTER TABLE organizations
    DROP COLUMN IF EXISTS tax_address,
    DROP COLUMN IF EXISTS tax_legal_name,
    DROP COLUMN IF EXISTS is_vat_registered,
    DROP COLUMN IF EXISTS pan_number;
```

- [ ] **Step 3: Add the new tables to test cleanup**

In `tests/e2e/setup_test.go`, in `cleanupTables`, add to the front of the
`tables` slice (before `"payments"`, so child rows go first):

```go
		"invoice_prints", "invoice_items", "invoices", "invoice_counters",
```

- [ ] **Step 4: Write the schema tests**

Create `tests/e2e/invoice_schema_test.go`. These assert the database refuses
what the service layer must never be trusted to prevent.

```go
package e2e

import (
	"context"
	"strings"
	"testing"
)

// seedInvoice inserts a minimal issued invoice directly and returns its id.
func seedInvoice(t *testing.T, orgID, issuedBy string, seq int) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			seller_name, seller_pan, customer_name,
			issued_date, issued_date_bs,
			subtotal, taxable_amount, total, amount_in_words, issued_by
		) VALUES ($1, '2082-83', $2, '2082-83/' || lpad($2::text, 6, '0'),
			'Test Gym', '123456789', 'Ram Bahadur',
			CURRENT_DATE, '2082-04-14',
			1000, 1000, 1000, 'One thousand rupees only', $3)
		RETURNING id`,
		orgID, seq, issuedBy,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedInvoice: %v", err)
	}
	return id
}

func TestInvoiceSchema_CannotUpdateAmount(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000101", "Admin")
	orgID := createTestOrg(t, admin, "Immutable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET total = 5000 WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject an amount change, got nil")
	}
	if !strings.Contains(err.Error(), "only cancellation may change") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_CannotDelete(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000102", "Admin")
	orgID := createTestOrg(t, admin, "Undeletable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(), `DELETE FROM invoices WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject a delete, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be deleted") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_CancellationIsAllowedAndOneWay(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000103", "Admin")
	orgID := createTestOrg(t, admin, "Cancellable Gym")
	id := seedInvoice(t, orgID, admin, 1)

	// Cancelling is the one permitted mutation.
	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET status = 'cancelled', cancelled_at = NOW(),
		        cancelled_by = $2, cancellation_reason = 'wrong customer'
		 WHERE id = $1`, id, admin)
	if err != nil {
		t.Fatalf("cancelling should be permitted: %v", err)
	}

	// Un-cancelling is not.
	_, err = testPool.Exec(context.Background(),
		`UPDATE invoices SET status = 'issued' WHERE id = $1`, id)
	if err == nil {
		t.Fatal("expected un-cancelling to be rejected, got nil")
	}
	if !strings.Contains(err.Error(), "cancellation cannot be undone") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvoiceSchema_LineItemsAreInsertOnly(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000104", "Admin")
	orgID := createTestOrg(t, admin, "Lineitem Gym")
	id := seedInvoice(t, orgID, admin, 1)

	_, err := testPool.Exec(context.Background(),
		`INSERT INTO invoice_items (invoice_id, line_no, description, quantity, unit_price, amount)
		 VALUES ($1, 1, 'Monthly Boxing', 1, 1000, 1000)`, id)
	if err != nil {
		t.Fatalf("inserting a line item should work: %v", err)
	}

	_, err = testPool.Exec(context.Background(),
		`UPDATE invoice_items SET amount = 9999 WHERE invoice_id = $1`, id)
	if err == nil {
		t.Fatal("expected line item update to be rejected, got nil")
	}
}

func TestInvoiceSchema_SequenceIsUniquePerOrgAndYear(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000105", "Admin")
	orgA := createTestOrg(t, admin, "Seq Gym A")
	orgB := createTestOrg(t, admin, "Seq Gym B")

	seedInvoice(t, orgA, admin, 1)

	// The same sequence in a different org is fine.
	seedInvoice(t, orgB, admin, 1)

	// The same sequence in the same org and year is not.
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			seller_name, seller_pan, customer_name, issued_date, issued_date_bs,
			subtotal, taxable_amount, total, amount_in_words, issued_by
		) VALUES ($1, '2082-83', 1, '2082-83/000001',
			'Test Gym', '123456789', 'Someone', CURRENT_DATE, '2082-04-14',
			1000, 1000, 1000, 'One thousand rupees only', $2)`,
		orgA, admin)
	if err == nil {
		t.Fatal("expected duplicate sequence to be rejected, got nil")
	}
}
```

- [ ] **Step 5: Run the schema tests**

Run: `go test ./tests/e2e/ -run 'TestInvoiceSchema' -v`
Expected: PASS, all five.

- [ ] **Step 6: Commit**

```bash
git add db/migrations/000053_pan_billing.up.sql db/migrations/000053_pan_billing.down.sql tests/e2e/setup_test.go tests/e2e/invoice_schema_test.go
git commit -m "$(cat <<'EOF'
feat: add invoice schema with database-enforced immutability

An issued tax invoice cannot lawfully be edited or deleted, and its number is
consumed permanently. Enforcing that only in the service layer would mean one
careless UPDATE could quietly falsify the books, so the guarantee lives in
triggers: invoices reject DELETE outright and reject UPDATE to anything but
the cancellation fields, and line items are insert-only.

Numbering uses a locked counter row rather than a Postgres sequence, because
sequences are non-transactional and a rolled-back insert would leave a gap in
a run of numbers that is required to be unbroken.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: Org tax settings

**Files:**
- Modify: `internal/org/handler.go` (add fields to update request)
- Modify: `internal/org/repository.go` (persist and return them)
- Test: `tests/e2e/invoice_settings_test.go`

**Interfaces:**
- Consumes: Task 3's columns.
- Produces: `PUT /api/v1/orgs/{orgId}` accepts and returns `pan_number`, `tax_legal_name`, `tax_address`; `GET` surfaces them.

- [ ] **Step 1: Write the failing test**

```go
package e2e

import (
	"net/http"
	"testing"
)

func TestOrgSettings_SetPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000201", "Admin")
	orgID := createTestOrg(t, admin, "PAN Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "601234567",
		"tax_legal_name": "PAN Gym Pvt Ltd",
		"tax_address":    "Kirtipur, Kathmandu",
	}, token)
	assertStatus(t, resp, http.StatusOK)

	var got struct {
		PANNumber    string `json:"pan_number"`
		TaxLegalName string `json:"tax_legal_name"`
	}
	parseJSON(t, resp, &got)
	if got.PANNumber != "601234567" {
		t.Errorf("pan_number = %q, want %q", got.PANNumber, "601234567")
	}
	if got.TaxLegalName != "PAN Gym Pvt Ltd" {
		t.Errorf("tax_legal_name = %q, want %q", got.TaxLegalName, "PAN Gym Pvt Ltd")
	}
}

func TestOrgSettings_RejectsMalformedPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000202", "Admin")
	orgID := createTestOrg(t, admin, "Bad PAN Gym")
	token := generateTestToken(admin, "member")

	for _, bad := range []string{"12345", "abcdefghi", "6012345678"} {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
			map[string]any{"pan_number": bad}, token)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("PAN %q: status = %d, want 400", bad, resp.StatusCode)
		}
		resp.Body.Close()
	}
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./tests/e2e/ -run 'TestOrgSettings' -v`
Expected: FAIL — `pan_number` comes back empty because the field is not wired.

- [ ] **Step 3: Add the fields to the org update path**

In `internal/org/handler.go`, add to the update request struct (near the
existing `Settings *json.RawMessage`):

```go
	PANNumber    *string `json:"pan_number"`
	TaxLegalName *string `json:"tax_legal_name"`
	TaxAddress   *string `json:"tax_address"`
```

Validate before calling the service — a malformed PAN on a tax document is
worse than a rejected save:

```go
	if req.PANNumber != nil && *req.PANNumber != "" {
		if !regexp.MustCompile(`^[0-9]{9}$`).MatchString(*req.PANNumber) {
			writeJSON(w, http.StatusBadRequest,
				map[string]string{"error": "pan_number must be exactly 9 digits"})
			return
		}
	}
```

Thread the three fields through the service call and into the repository's
`UPDATE organizations SET ... pan_number = COALESCE($n, pan_number)` clause,
following the `COALESCE` pattern already used there. Add all three to the
`SELECT` list and the scan target of the org read so they round-trip.

- [ ] **Step 4: Run to verify it passes**

Run: `go test ./tests/e2e/ -run 'TestOrgSettings' -v`
Expected: PASS, both tests.

- [ ] **Step 5: Commit**

```bash
git add internal/org/ tests/e2e/invoice_settings_test.go
git commit -m "$(cat <<'EOF'
feat: store a gym's PAN number and registered tax identity

A tax invoice cannot be issued without the seller's PAN, and the registered
name and address on it frequently differ from the trading ones, so they are
stored separately rather than reusing the org's display fields.

The 9-digit format is validated at the handler as well as by a CHECK
constraint: a malformed PAN reaching a printed bill is worse than a rejected
save.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: Issue a bill

**Files:**
- Create: `internal/invoice/models.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`
- Modify: `cmd/api/main.go`, `tests/e2e/setup_test.go`
- Test: `internal/invoice/service_test.go`, `tests/e2e/invoice_test.go`

**Interfaces:**
- Consumes: `nepalidate.Today()`, `nepalidate.Date.FiscalYear()`, `moneywords.Rupees()`.
- Produces:
  - `invoice.NewRepository(*pgxpool.Pool) *Repository`
  - `invoice.NewService(*Repository, *slog.Logger) *Service`
  - `invoice.NewHandler(*Service, *slog.Logger) *Handler`
  - `(*Handler) RegisterRoutes(chi.Router)`
  - `Service.Issue(ctx, orgID, issuedBy string, in IssueInput) (*Invoice, error)`
  - `IssueInput{CustomerUserID, CustomerName, CustomerPAN, CustomerAddress, CustomerPhone, PaymentMethod, TransactionID, MemberPackageID string; Discount float64; Items []ItemInput}`
  - `ItemInput{Description, DescriptionNe string; Quantity, UnitPrice float64}`
  - Errors: `ErrPANNotConfigured`, `ErrNoItems`, `ErrInvalidDiscount`, `ErrAlreadyBilled`, `ErrNotFound`

- [ ] **Step 1: Write the service validation tests**

Create `internal/invoice/service_test.go`. Following the repo's convention,
service tests exercise validation with a nil repository — anything reaching SQL
is covered by the e2e suite.

```go
package invoice

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_Issue_RejectsEmptyItems(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		CustomerName: "Ram",
		Items:        nil,
	})
	if !errors.Is(err, ErrNoItems) {
		t.Fatalf("error = %v, want ErrNoItems", err)
	}
}

func TestService_Issue_RejectsMissingCustomerName(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		Items: []ItemInput{{Description: "Boxing", Quantity: 1, UnitPrice: 1000}},
	})
	if err == nil {
		t.Fatal("expected an error for a bill with no customer name, got nil")
	}
}

func TestService_Issue_RejectsBadLineValues(t *testing.T) {
	tests := []struct {
		name string
		item ItemInput
	}{
		{"zero quantity", ItemInput{Description: "Boxing", Quantity: 0, UnitPrice: 1000}},
		{"negative quantity", ItemInput{Description: "Boxing", Quantity: -1, UnitPrice: 1000}},
		{"negative price", ItemInput{Description: "Boxing", Quantity: 1, UnitPrice: -5}},
		{"blank description", ItemInput{Description: "  ", Quantity: 1, UnitPrice: 1000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{logger: testLogger()}
			_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
				CustomerName: "Ram",
				Items:        []ItemInput{tt.item},
			})
			if err == nil {
				t.Fatalf("expected an error for %s, got nil", tt.name)
			}
		})
	}
}

func TestService_Issue_RejectsDiscountAboveSubtotal(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		CustomerName: "Ram",
		Discount:     2000,
		Items:        []ItemInput{{Description: "Boxing", Quantity: 1, UnitPrice: 1000}},
	})
	if !errors.Is(err, ErrInvalidDiscount) {
		t.Fatalf("error = %v, want ErrInvalidDiscount", err)
	}
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./internal/invoice/ -v`
Expected: FAIL — package does not compile, nothing is defined yet.

- [ ] **Step 3: Write `models.go`**

```go
package invoice

import (
	"errors"
	"time"
)

var (
	// ErrPANNotConfigured is returned when a gym tries to bill before setting
	// its PAN. Handled as its own code so the UI can link to settings rather
	// than dead-ending on a message.
	ErrPANNotConfigured = errors.New("pan_not_configured")
	ErrNoItems          = errors.New("an invoice needs at least one line item")
	ErrInvalidDiscount  = errors.New("discount cannot exceed the subtotal")
	ErrAlreadyBilled    = errors.New("this payment already has a bill")
	ErrNotFound         = errors.New("invoice not found")
	ErrAlreadyCancelled = errors.New("invoice is already cancelled")
	ErrReasonRequired   = errors.New("a cancellation reason is required")
	ErrInvoiceCancelled = errors.New("cannot credit a cancelled invoice")
	ErrInvalidParent    = errors.New("a credit note cannot reference another credit note")
	ErrCreditTooLarge   = errors.New("credit exceeds the uncredited balance of the invoice")
)

// Item is one line on a bill.
type Item struct {
	LineNo        int     `json:"line_no"`
	Description   string  `json:"description"`
	DescriptionNe string  `json:"description_ne,omitempty"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	Amount        float64 `json:"amount"`
}

// Invoice is an issued tax document. Every field a bill prints is snapshotted
// here rather than joined at read time, so a bill never changes retroactively
// when the gym or the customer edits their details.
type Invoice struct {
	ID            string `json:"id"`
	OrgID         string `json:"organization_id"`
	FiscalYear    string `json:"fiscal_year"`
	Sequence      int    `json:"sequence"`
	InvoiceNumber string `json:"invoice_number"`
	DocType       string `json:"doc_type"`
	CreditNoteFor string `json:"credit_note_for,omitempty"`

	SellerName          string `json:"seller_name"`
	SellerPAN           string `json:"seller_pan"`
	SellerAddress       string `json:"seller_address,omitempty"`
	SellerVATRegistered bool   `json:"seller_vat_registered"`

	CustomerUserID  string `json:"customer_user_id,omitempty"`
	CustomerName    string `json:"customer_name"`
	CustomerPAN     string `json:"customer_pan,omitempty"`
	CustomerAddress string `json:"customer_address,omitempty"`
	CustomerPhone   string `json:"customer_phone,omitempty"`

	IssuedDate   string `json:"issued_date"`
	IssuedDateBS string `json:"issued_date_bs"`

	Subtotal      float64 `json:"subtotal"`
	Discount      float64 `json:"discount"`
	TaxableAmount float64 `json:"taxable_amount"`
	VATRate       float64 `json:"vat_rate"`
	VATAmount     float64 `json:"vat_amount"`
	Total         float64 `json:"total"`
	AmountInWords string  `json:"amount_in_words"`

	PaymentMethod string `json:"payment_method,omitempty"`

	Status             string     `json:"status"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`

	TransactionID   string `json:"transaction_id,omitempty"`
	MemberPackageID string `json:"member_package_id,omitempty"`

	IssuedBy   string    `json:"issued_by"`
	PrintCount int       `json:"print_count"`
	CreatedAt  time.Time `json:"created_at"`

	Items []Item `json:"items,omitempty"`
}

// ItemInput is a requested line on a new bill.
type ItemInput struct {
	Description   string  `json:"description"`
	DescriptionNe string  `json:"description_ne"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
}

// IssueInput is a request to issue a bill.
type IssueInput struct {
	CustomerUserID  string      `json:"customer_user_id"`
	CustomerName    string      `json:"customer_name"`
	CustomerPAN     string      `json:"customer_pan"`
	CustomerAddress string      `json:"customer_address"`
	CustomerPhone   string      `json:"customer_phone"`
	PaymentMethod   string      `json:"payment_method"`
	TransactionID   string      `json:"transaction_id"`
	MemberPackageID string      `json:"member_package_id"`
	Discount        float64     `json:"discount"`
	Items           []ItemInput `json:"items"`
}

// ListFilter narrows a list query.
type ListFilter struct {
	OrgID      string
	Status     string
	FiscalYear string
	CustomerID string
	From       string
	To         string
	Query      string
	Limit      int
	Offset     int
}
```

- [ ] **Step 4: Write `repository.go` — the counter allocation and insert**

The whole design rests on this function. Every step happens in one transaction
so that a failure anywhere rolls the number back rather than burning it.

```go
package invoice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository owns all invoice SQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new invoice repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// sellerIdentity is the gym's tax identity at the moment of issue.
type sellerIdentity struct {
	Name          string
	PAN           string
	Address       string
	VATRegistered bool
}

// loadSeller reads the org's tax identity, falling back to its display name and
// address when the registered ones are blank.
func (r *Repository) loadSeller(ctx context.Context, tx pgx.Tx, orgID string) (sellerIdentity, error) {
	var s sellerIdentity
	var pan, legalName, taxAddr, name, addr *string
	err := tx.QueryRow(ctx,
		`SELECT pan_number, tax_legal_name, tax_address, name, address, is_vat_registered
		   FROM organizations WHERE id = $1 AND is_active = true`,
		orgID,
	).Scan(&pan, &legalName, &taxAddr, &name, &addr, &s.VATRegistered)
	if errors.Is(err, pgx.ErrNoRows) {
		return s, ErrNotFound
	}
	if err != nil {
		return s, fmt.Errorf("loading seller identity: %w", err)
	}

	if pan == nil || strings.TrimSpace(*pan) == "" {
		return s, ErrPANNotConfigured
	}
	s.PAN = *pan

	s.Name = deref(name)
	if legalName != nil && strings.TrimSpace(*legalName) != "" {
		s.Name = *legalName
	}
	s.Address = deref(addr)
	if taxAddr != nil && strings.TrimSpace(*taxAddr) != "" {
		s.Address = *taxAddr
	}
	return s, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// allocateSequence reserves the next number for this org and fiscal year.
//
// A Postgres sequence would be wrong: sequences are non-transactional, so a
// rolled-back insert burns its number and leaves a gap in a run that is
// required to be unbroken. The no-op DO UPDATE is what takes the row lock —
// ON CONFLICT DO NOTHING returns no row and would race.
func allocateSequence(ctx context.Context, tx pgx.Tx, orgID, fiscalYear string) (int, error) {
	var seq int
	err := tx.QueryRow(ctx,
		`INSERT INTO invoice_counters (organization_id, fiscal_year, next_sequence)
		 VALUES ($1, $2, 1)
		 ON CONFLICT (organization_id, fiscal_year)
		 DO UPDATE SET next_sequence = invoice_counters.next_sequence
		 RETURNING next_sequence`,
		orgID, fiscalYear,
	).Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("allocating invoice number: %w", err)
	}

	if _, err := tx.Exec(ctx,
		`UPDATE invoice_counters SET next_sequence = next_sequence + 1
		  WHERE organization_id = $1 AND fiscal_year = $2`,
		orgID, fiscalYear,
	); err != nil {
		return 0, fmt.Errorf("advancing invoice number: %w", err)
	}

	return seq, nil
}

// issueParams carries everything the service computed for a new document.
type issueParams struct {
	OrgID         string
	FiscalYear    string
	DocType       string
	CreditNoteFor string
	IssuedDate    string
	IssuedDateBS  string
	Subtotal      float64
	Discount      float64
	TaxableAmount float64
	Total         float64
	AmountInWords string
	IssuedBy      string
	In            IssueInput
}

// insertDocument writes an invoice or credit note and its line items inside tx.
func insertDocument(ctx context.Context, tx pgx.Tx, p issueParams, seller sellerIdentity) (string, string, int, error) {
	seq, err := allocateSequence(ctx, tx, p.OrgID, p.FiscalYear)
	if err != nil {
		return "", "", 0, err
	}
	number := fmt.Sprintf("%s/%06d", p.FiscalYear, seq)

	var id string
	err = tx.QueryRow(ctx,
		`INSERT INTO invoices (
			organization_id, fiscal_year, sequence, invoice_number,
			doc_type, credit_note_for,
			seller_name, seller_pan, seller_address, seller_vat_registered,
			customer_user_id, customer_name, customer_pan, customer_address, customer_phone,
			issued_date, issued_date_bs,
			subtotal, discount, taxable_amount, vat_rate, vat_amount, total, amount_in_words,
			payment_method, transaction_id, member_package_id, issued_by
		 ) VALUES (
			$1, $2, $3, $4,
			$5, NULLIF($6, '')::uuid,
			$7, $8, $9, $10,
			NULLIF($11, '')::uuid, $12, NULLIF($13, ''), NULLIF($14, ''), NULLIF($15, ''),
			$16::date, $17,
			$18, $19, $20, 0, 0, $21, $22,
			NULLIF($23, ''), NULLIF($24, '')::uuid, NULLIF($25, '')::uuid, $26
		 ) RETURNING id`,
		p.OrgID, p.FiscalYear, seq, number,
		p.DocType, p.CreditNoteFor,
		seller.Name, seller.PAN, seller.Address, seller.VATRegistered,
		p.In.CustomerUserID, p.In.CustomerName, p.In.CustomerPAN, p.In.CustomerAddress, p.In.CustomerPhone,
		p.IssuedDate, p.IssuedDateBS,
		p.Subtotal, p.Discount, p.TaxableAmount, p.Total, p.AmountInWords,
		p.In.PaymentMethod, p.In.TransactionID, p.In.MemberPackageID, p.IssuedBy,
	).Scan(&id)
	if err != nil {
		// The partial unique index on transaction_id is what stops the same
		// payment being billed twice.
		if strings.Contains(err.Error(), "idx_invoices_one_per_transaction") {
			return "", "", 0, ErrAlreadyBilled
		}
		return "", "", 0, fmt.Errorf("inserting invoice: %w", err)
	}

	for i, it := range p.In.Items {
		if _, err := tx.Exec(ctx,
			`INSERT INTO invoice_items (invoice_id, line_no, description, description_ne, quantity, unit_price, amount)
			 VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7)`,
			id, i+1, it.Description, it.DescriptionNe, it.Quantity, it.UnitPrice,
			round2(it.Quantity*it.UnitPrice),
		); err != nil {
			return "", "", 0, fmt.Errorf("inserting line item %d: %w", i+1, err)
		}
	}

	return id, number, seq, nil
}

// Issue writes a new bill and returns it.
func (r *Repository) Issue(ctx context.Context, p issueParams) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	seller, err := r.loadSeller(ctx, tx, p.OrgID)
	if err != nil {
		return nil, err
	}

	id, _, _, err := insertDocument(ctx, tx, p, seller)
	if err != nil {
		return nil, err
	}

	inv, err := getInTx(ctx, tx, p.OrgID, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing invoice: %w", err)
	}
	return inv, nil
}
```

Add `getInTx`, `Get`, and `List`. `Get` is org-scoped in the `WHERE` clause so a
cross-tenant read returns `ErrNotFound` rather than another org's document:

```go
const invoiceColumns = `
	id, organization_id, fiscal_year, sequence, invoice_number,
	doc_type, COALESCE(credit_note_for::text, ''),
	seller_name, seller_pan, COALESCE(seller_address, ''), seller_vat_registered,
	COALESCE(customer_user_id::text, ''), customer_name,
	COALESCE(customer_pan, ''), COALESCE(customer_address, ''), COALESCE(customer_phone, ''),
	issued_date::text, issued_date_bs,
	subtotal, discount, taxable_amount, vat_rate, vat_amount, total, amount_in_words,
	COALESCE(payment_method, ''),
	status, cancelled_at, COALESCE(cancellation_reason, ''),
	COALESCE(transaction_id::text, ''), COALESCE(member_package_id::text, ''),
	issued_by, print_count, created_at`

func scanInvoice(row pgx.Row) (*Invoice, error) {
	var v Invoice
	err := row.Scan(
		&v.ID, &v.OrgID, &v.FiscalYear, &v.Sequence, &v.InvoiceNumber,
		&v.DocType, &v.CreditNoteFor,
		&v.SellerName, &v.SellerPAN, &v.SellerAddress, &v.SellerVATRegistered,
		&v.CustomerUserID, &v.CustomerName,
		&v.CustomerPAN, &v.CustomerAddress, &v.CustomerPhone,
		&v.IssuedDate, &v.IssuedDateBS,
		&v.Subtotal, &v.Discount, &v.TaxableAmount, &v.VATRate, &v.VATAmount,
		&v.Total, &v.AmountInWords,
		&v.PaymentMethod,
		&v.Status, &v.CancelledAt, &v.CancellationReason,
		&v.TransactionID, &v.MemberPackageID,
		&v.IssuedBy, &v.PrintCount, &v.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning invoice: %w", err)
	}
	return &v, nil
}

func getInTx(ctx context.Context, tx pgx.Tx, orgID, id string) (*Invoice, error) {
	inv, err := scanInvoice(tx.QueryRow(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE id = $1 AND organization_id = $2`,
		id, orgID))
	if err != nil {
		return nil, err
	}
	items, err := loadItems(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

func loadItems(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, invoiceID string) ([]Item, error) {
	rows, err := q.Query(ctx,
		`SELECT line_no, description, COALESCE(description_ne, ''), quantity, unit_price, amount
		   FROM invoice_items WHERE invoice_id = $1 ORDER BY line_no`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("loading line items: %w", err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.LineNo, &it.Description, &it.DescriptionNe,
			&it.Quantity, &it.UnitPrice, &it.Amount); err != nil {
			return nil, fmt.Errorf("scanning line item: %w", err)
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// Get reads one invoice. Scoping by org in the WHERE clause means another
// gym's invoice is indistinguishable from a missing one.
func (r *Repository) Get(ctx context.Context, orgID, id string) (*Invoice, error) {
	inv, err := scanInvoice(r.db.QueryRow(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE id = $1 AND organization_id = $2`,
		id, orgID))
	if err != nil {
		return nil, err
	}
	items, err := loadItems(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

// List returns invoices for an org, newest first.
func (r *Repository) List(ctx context.Context, f ListFilter) ([]Invoice, int, error) {
	conds := []string{"organization_id = $1"}
	args := []any{f.OrgID}

	add := func(clause string, val any) {
		args = append(args, val)
		conds = append(conds, fmt.Sprintf(clause, len(args)))
	}
	if f.Status != "" {
		add("status = $%d", f.Status)
	}
	if f.FiscalYear != "" {
		add("fiscal_year = $%d", f.FiscalYear)
	}
	if f.CustomerID != "" {
		add("customer_user_id = $%d::uuid", f.CustomerID)
	}
	if f.From != "" {
		add("issued_date >= $%d::date", f.From)
	}
	if f.To != "" {
		add("issued_date <= $%d::date", f.To)
	}
	if f.Query != "" {
		args = append(args, "%"+f.Query+"%")
		conds = append(conds, fmt.Sprintf(
			"(invoice_number ILIKE $%d OR customer_name ILIKE $%d)", len(args), len(args)))
	}
	where := strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM invoices WHERE `+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting invoices: %w", err)
	}

	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE `+where+
			fmt.Sprintf(" ORDER BY sequence DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args)),
		args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing invoices: %w", err)
	}
	defer rows.Close()

	out := []Invoice{}
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *inv)
	}
	return out, total, rows.Err()
}
```

- [ ] **Step 5: Write `service.go`**

```go
package invoice

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/urja-gym/urja/pkg/moneywords"
	"github.com/urja-gym/urja/pkg/nepalidate"
)

// Service holds invoice business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new invoice service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func round2(f float64) float64 {
	return math.Round(f*100) / 100
}

var validPaymentMethods = map[string]bool{
	"cash": true, "bank_transfer": true, "khalti": true,
}

// validateIssue checks an issue request and returns the computed money.
func validateIssue(in IssueInput) (subtotal, taxable float64, err error) {
	if strings.TrimSpace(in.CustomerName) == "" {
		return 0, 0, fmt.Errorf("customer name is required")
	}
	if len(in.Items) == 0 {
		return 0, 0, ErrNoItems
	}
	if in.PaymentMethod != "" && !validPaymentMethods[in.PaymentMethod] {
		return 0, 0, fmt.Errorf("invalid payment method: %s", in.PaymentMethod)
	}
	if in.CustomerPAN != "" && !isNineDigits(in.CustomerPAN) {
		return 0, 0, fmt.Errorf("customer PAN must be exactly 9 digits")
	}

	for i, it := range in.Items {
		if strings.TrimSpace(it.Description) == "" {
			return 0, 0, fmt.Errorf("line %d: description is required", i+1)
		}
		if it.Quantity <= 0 {
			return 0, 0, fmt.Errorf("line %d: quantity must be greater than 0", i+1)
		}
		if it.UnitPrice < 0 {
			return 0, 0, fmt.Errorf("line %d: unit price cannot be negative", i+1)
		}
		subtotal += round2(it.Quantity * it.UnitPrice)
	}
	subtotal = round2(subtotal)

	if in.Discount < 0 {
		return 0, 0, fmt.Errorf("discount cannot be negative")
	}
	if in.Discount > subtotal {
		return 0, 0, ErrInvalidDiscount
	}

	return subtotal, round2(subtotal - in.Discount), nil
}

func isNineDigits(s string) bool {
	if len(s) != 9 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Issue validates and writes a new bill.
func (s *Service) Issue(ctx context.Context, orgID, issuedBy string, in IssueInput) (*Invoice, error) {
	subtotal, taxable, err := validateIssue(in)
	if err != nil {
		return nil, err
	}

	bs, err := nepalidate.Today()
	if err != nil {
		return nil, fmt.Errorf("determining Nepali date: %w", err)
	}
	loc, _ := time.LoadLocation("Asia/Kathmandu")

	inv, err := s.repo.Issue(ctx, issueParams{
		OrgID:         orgID,
		FiscalYear:    bs.FiscalYear(),
		DocType:       "invoice",
		IssuedDate:    time.Now().In(loc).Format("2006-01-02"),
		IssuedDateBS:  bs.String(),
		Subtotal:      subtotal,
		Discount:      round2(in.Discount),
		TaxableAmount: taxable,
		// PAN-only: no VAT is added, so the total is the taxable amount.
		Total:         taxable,
		AmountInWords: moneywords.Rupees(taxable),
		IssuedBy:      issuedBy,
		In:            in,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("invoice issued",
		"invoice_number", inv.InvoiceNumber, "org_id", orgID,
		"total", inv.Total, "issued_by", issuedBy)
	return inv, nil
}

// Get reads one invoice, scoped to the org.
func (s *Service) Get(ctx context.Context, orgID, id string) (*Invoice, error) {
	return s.repo.Get(ctx, orgID, id)
}

// List returns invoices for an org.
func (s *Service) List(ctx context.Context, f ListFilter) ([]Invoice, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.Status != "" && f.Status != "issued" && f.Status != "cancelled" {
		return nil, 0, fmt.Errorf("invalid status filter: %s", f.Status)
	}
	return s.repo.List(ctx, f)
}
```

- [ ] **Step 6: Run the service tests**

Run: `go test ./internal/invoice/ -v`
Expected: PASS, all four test functions.

- [ ] **Step 7: Write `handler.go` and `routes.go`**

`routes.go`:

```go
package invoice

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts invoice routes (under /orgs/{orgId}/invoices).
// Billing is an admin and staff function; members never reach these.
//
// There is deliberately no PUT and no DELETE: an issued bill is immutable and
// is corrected by cancelling it or raising a credit note.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("admin", "staff"))

	r.Get("/", h.List)
	r.Post("/", h.Issue)
	r.Get("/next-number", h.NextNumber)
	r.Get("/{invoiceId}", h.Get)
	r.Post("/{invoiceId}/cancel", h.Cancel)
	r.Post("/{invoiceId}/credit-note", h.CreditNote)
	r.Post("/{invoiceId}/print", h.Print)
}
```

`handler.go` — decode, delegate, map errors. The error map is the contract the
UI depends on:

```go
package invoice

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler serves the invoice HTTP endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new invoice handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// writeErr maps a domain error onto the status and code the UI expects.
// Cross-tenant reads deliberately return 404: a 403 would confirm the row
// exists in another gym.
func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "invoice not found", "code": "not_found"})
	case errors.Is(err, ErrPANNotConfigured):
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "set your PAN number in settings before issuing a bill",
			"code":  "pan_not_configured"})
	case errors.Is(err, ErrAlreadyCancelled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "already_cancelled"})
	case errors.Is(err, ErrInvoiceCancelled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "invoice_cancelled"})
	case errors.Is(err, ErrAlreadyBilled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "already_billed"})
	case errors.Is(err, ErrCreditTooLarge):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "credit_exceeds_balance"})
	case errors.Is(err, ErrInvalidParent):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "invalid_parent"})
	case errors.Is(err, ErrReasonRequired):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "reason_required"})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}

// Issue handles POST /api/v1/orgs/{orgId}/invoices
func (h *Handler) Issue(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var in IssueInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	inv, err := h.service.Issue(r.Context(), orgID, userID, in)
	if err != nil {
		h.logger.Error("issuing invoice failed", "error", err, "org_id", orgID)
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, inv)
}

// Get handles GET /api/v1/orgs/{orgId}/invoices/{invoiceId}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	inv, err := h.service.Get(r.Context(), orgID, chi.URLParam(r, "invoiceId"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

// List handles GET /api/v1/orgs/{orgId}/invoices
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	list, total, err := h.service.List(r.Context(), ListFilter{
		OrgID:      orgID,
		Status:     q.Get("status"),
		FiscalYear: q.Get("fiscal_year"),
		CustomerID: q.Get("customer_user_id"),
		From:       q.Get("from"),
		To:         q.Get("to"),
		Query:      q.Get("q"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": list, "total": total})
}
```

Stub the three not-yet-built handlers so the package compiles; Tasks 6–8
replace each body:

```go
// Cancel handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// CreditNote handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/credit-note
func (h *Handler) CreditNote(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// Print handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/print
func (h *Handler) Print(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// NextNumber handles GET /api/v1/orgs/{orgId}/invoices/next-number
func (h *Handler) NextNumber(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}
```

- [ ] **Step 8: Wire the domain into both routers**

In `cmd/api/main.go`, alongside the other domain wiring:

```go
	invoiceRepo := invoice.NewRepository(pool)
	invoiceSvc := invoice.NewService(invoiceRepo, logger)
	invoiceHandler := invoice.NewHandler(invoiceSvc, logger)
```

and inside `r.Route("/orgs/{orgId}", ...)`:

```go
				r.Route("/invoices", func(r chi.Router) {
					invoiceHandler.RegisterRoutes(r)
				})
```

Make the identical two changes in `tests/e2e/setup_test.go` — it mirrors
`main.go`, and a route registered in only one of them produces tests that pass
against a server the product does not actually run.

- [ ] **Step 9: Write the e2e issue tests**

Create `tests/e2e/invoice_test.go`:

```go
package e2e

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// setPAN gives an org a PAN so it can issue bills.
func setPAN(t *testing.T, orgID, pan string) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`UPDATE organizations SET pan_number = $2 WHERE id = $1`, orgID, pan)
	if err != nil {
		t.Fatalf("setPAN: %v", err)
	}
}

func issueBody(customer string, qty, price float64) map[string]any {
	return map[string]any{
		"customer_name":  customer,
		"payment_method": "cash",
		"items": []map[string]any{
			{"description": "Monthly Boxing", "quantity": qty, "unit_price": price},
		},
	}
}

func TestInvoice_IssueRequiresPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000301", "Admin")
	orgID := createTestOrg(t, admin, "No PAN Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram Bahadur", 1, 3000), token)
	assertStatus(t, resp, http.StatusBadRequest)

	var body struct{ Code string }
	parseJSON(t, resp, &body)
	if body.Code != "pan_not_configured" {
		t.Errorf("code = %q, want %q", body.Code, "pan_not_configured")
	}
}

func TestInvoice_IssueComputesTotalsAndNumber(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000302", "Admin")
	orgID := createTestOrg(t, admin, "Billing Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	body := issueBody("Ram Bahadur", 2, 1500)
	body["discount"] = 500

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusCreated)

	var inv struct {
		InvoiceNumber string  `json:"invoice_number"`
		Sequence      int     `json:"sequence"`
		Subtotal      float64 `json:"subtotal"`
		Discount      float64 `json:"discount"`
		Total         float64 `json:"total"`
		VATAmount     float64 `json:"vat_amount"`
		AmountInWords string  `json:"amount_in_words"`
		SellerPAN     string  `json:"seller_pan"`
		Status        string  `json:"status"`
		Items         []struct {
			Amount float64 `json:"amount"`
		} `json:"items"`
	}
	parseJSON(t, resp, &inv)

	if inv.Sequence != 1 {
		t.Errorf("sequence = %d, want 1", inv.Sequence)
	}
	if inv.Subtotal != 3000 || inv.Discount != 500 || inv.Total != 2500 {
		t.Errorf("totals = subtotal %.2f discount %.2f total %.2f, want 3000/500/2500",
			inv.Subtotal, inv.Discount, inv.Total)
	}
	if inv.VATAmount != 0 {
		t.Errorf("vat_amount = %.2f, want 0 on a PAN-only bill", inv.VATAmount)
	}
	if inv.SellerPAN != "601234567" {
		t.Errorf("seller_pan = %q, want the org's PAN snapshotted", inv.SellerPAN)
	}
	if inv.Status != "issued" {
		t.Errorf("status = %q, want issued", inv.Status)
	}
	if len(inv.Items) != 1 || inv.Items[0].Amount != 3000 {
		t.Errorf("items = %+v, want one line of 3000", inv.Items)
	}
	if inv.AmountInWords == "" {
		t.Error("amount_in_words is empty")
	}
}

// The core guarantee: numbers must run 1..N with no gap and no duplicate even
// when several staff issue at the same moment.
func TestInvoice_ConcurrentIssuesAreGapless(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000303", "Admin")
	orgID := createTestOrg(t, admin, "Concurrent Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	const n = 12
	seqs := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
				issueBody(fmt.Sprintf("Customer %d", i), 1, 1000), token)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				return
			}
			var inv struct {
				Sequence int `json:"sequence"`
			}
			parseJSON(t, resp, &inv)
			seqs[i] = inv.Sequence
		}(i)
	}
	wg.Wait()

	seen := map[int]bool{}
	for _, s := range seqs {
		if s == 0 {
			t.Fatal("an issue failed; all 12 should succeed")
		}
		if seen[s] {
			t.Fatalf("sequence %d was allocated twice", s)
		}
		seen[s] = true
	}
	for want := 1; want <= n; want++ {
		if !seen[want] {
			t.Errorf("sequence %d is missing — the run has a gap", want)
		}
	}
}

func TestInvoice_CrossTenantIsInvisible(t *testing.T) {
	cleanupTables(t)
	adminA := createTestUser(t, "9800000304", "Admin A")
	adminB := createTestUser(t, "9800000305", "Admin B")
	orgA := createTestOrg(t, adminA, "Tenant A Gym")
	orgB := createTestOrg(t, adminB, "Tenant B Gym")
	setPAN(t, orgA, "601234567")
	tokenA := generateTestToken(adminA, "member")
	tokenB := generateTestToken(adminB, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgA+"/invoices",
		issueBody("Ram", 1, 1000), tokenA)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	// B asking for A's invoice through B's own org must 404.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgB+"/invoices/"+inv.ID, nil, tokenB)
	assertStatus(t, resp, http.StatusNotFound)

	// And B must not reach it through A's org either.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgA+"/invoices/"+inv.ID, nil, tokenB)
	if resp.StatusCode == http.StatusOK {
		t.Error("org B read org A's invoice through org A's path")
	}
	resp.Body.Close()
}

func TestInvoice_MemberCannotIssue(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000306", "Admin")
	memberID := createTestUser(t, "9800000307", "Member")
	orgID := createTestOrg(t, admin, "Gated Gym")
	createTestOrgMember(t, memberID, orgID, "member")
	setPAN(t, orgID, "601234567")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), generateTestToken(memberID, "member"))
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403 — billing is staff and admin only", resp.StatusCode)
	}
	resp.Body.Close()
}
```

- [ ] **Step 10: Run the e2e tests**

Run: `go test ./tests/e2e/ -run 'TestInvoice_' -v`
Expected: PASS, all five. If `TestInvoice_ConcurrentIssuesAreGapless` fails,
the counter lock is wrong — do not move on.

- [ ] **Step 11: Commit**

```bash
git add internal/invoice/ cmd/api/main.go tests/e2e/setup_test.go tests/e2e/invoice_test.go
git commit -m "$(cat <<'EOF'
feat: issue PAN tax invoices

A gym could record money in three places and produce no document for any of
it. This adds the document, numbered from a locked counter row so a run of
invoice numbers can never contain a gap even when several staff bill at once.

Seller and customer details are snapshotted onto the invoice rather than
joined at read time: correcting the gym's PAN next year must not silently
rewrite bills issued this year.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: Cancel a bill

**Files:**
- Modify: `internal/invoice/repository.go`, `service.go`, `handler.go`
- Test: `tests/e2e/invoice_cancel_test.go`

**Interfaces:**
- Consumes: Task 5's `Service`, `Repository`, error values.
- Produces: `Service.Cancel(ctx, orgID, invoiceID, cancelledBy, reason string) (*Invoice, error)`.

- [ ] **Step 1: Write the failing test**

```go
package e2e

import (
	"net/http"
	"testing"
)

func TestInvoiceCancel_KeepsNumberConsumed(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000401", "Admin")
	orgID := createTestOrg(t, admin, "Cancel Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var first struct {
		ID       string `json:"id"`
		Sequence int    `json:"sequence"`
	}
	parseJSON(t, resp, &first)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+first.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	var cancelled struct {
		Status string `json:"status"`
		Reason string `json:"cancellation_reason"`
	}
	parseJSON(t, resp, &cancelled)
	if cancelled.Status != "cancelled" {
		t.Errorf("status = %q, want cancelled", cancelled.Status)
	}
	if cancelled.Reason != "wrong customer" {
		t.Errorf("reason = %q, want it recorded", cancelled.Reason)
	}

	// The cancelled number is spent: the next bill takes the one after it.
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Shyam", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var second struct {
		Sequence int `json:"sequence"`
	}
	parseJSON(t, resp, &second)
	if second.Sequence != first.Sequence+1 {
		t.Errorf("next sequence = %d, want %d — a cancelled number must not be reused",
			second.Sequence, first.Sequence+1)
	}
}

func TestInvoiceCancel_RequiresReason(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000402", "Admin")
	orgID := createTestOrg(t, admin, "Reason Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel",
		map[string]any{"reason": "   "}, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestInvoiceCancel_IsNotRepeatable(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000403", "Admin")
	orgID := createTestOrg(t, admin, "Twice Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	body := map[string]any{"reason": "duplicate"}
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel", body, token)
	assertStatus(t, resp, http.StatusOK)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel", body, token)
	assertStatus(t, resp, http.StatusConflict)
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./tests/e2e/ -run 'TestInvoiceCancel' -v`
Expected: FAIL — 501 Not Implemented from the stub.

- [ ] **Step 3: Implement the repository method**

```go
// Cancel marks an invoice cancelled. The number stays consumed; that is the
// point of cancelling rather than deleting.
func (r *Repository) Cancel(ctx context.Context, orgID, id, cancelledBy, reason string) (*Invoice, error) {
	tag, err := r.db.Exec(ctx,
		`UPDATE invoices
		    SET status = 'cancelled', cancelled_at = NOW(),
		        cancelled_by = $3, cancellation_reason = $4
		  WHERE id = $1 AND organization_id = $2 AND status = 'issued'`,
		id, orgID, cancelledBy, reason)
	if err != nil {
		return nil, fmt.Errorf("cancelling invoice: %w", err)
	}
	if tag.RowsAffected() == 0 {
		// Either it does not exist in this org, or it is already cancelled.
		existing, getErr := r.Get(ctx, orgID, id)
		if getErr != nil {
			return nil, getErr
		}
		if existing.Status == "cancelled" {
			return nil, ErrAlreadyCancelled
		}
		return nil, ErrNotFound
	}
	return r.Get(ctx, orgID, id)
}
```

- [ ] **Step 4: Implement the service method**

```go
// Cancel withdraws a bill issued in error. It does not touch the ledger — if
// money actually moved, a credit note is the correct instrument.
func (s *Service) Cancel(ctx context.Context, orgID, id, cancelledBy, reason string) (*Invoice, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrReasonRequired
	}
	inv, err := s.repo.Cancel(ctx, orgID, id, cancelledBy, reason)
	if err != nil {
		return nil, err
	}
	s.logger.Info("invoice cancelled",
		"invoice_number", inv.InvoiceNumber, "org_id", orgID,
		"cancelled_by", cancelledBy, "reason", reason)
	return inv, nil
}
```

- [ ] **Step 5: Replace the handler stub**

```go
type cancelRequest struct {
	Reason string `json:"reason"`
}

// Cancel handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req cancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	inv, err := h.service.Cancel(r.Context(), orgID, chi.URLParam(r, "invoiceId"), userID, req.Reason)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}
```

- [ ] **Step 6: Run to verify it passes**

Run: `go test ./tests/e2e/ -run 'TestInvoiceCancel' -v`
Expected: PASS, all three.

- [ ] **Step 7: Commit**

```bash
git add internal/invoice/ tests/e2e/invoice_cancel_test.go
git commit -m "$(cat <<'EOF'
feat: cancel an invoice without freeing its number

Cancelling is the lawful correction for a bill issued in error before it
reaches the customer. The number stays consumed — that is the whole
difference between cancelling and deleting, and deleting is not available.

A reason is required, because a cancelled bill with no explanation is exactly
what an audit asks about.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: Credit note

**Files:**
- Modify: `internal/invoice/repository.go`, `service.go`, `handler.go`
- Test: `tests/e2e/invoice_credit_test.go`

**Interfaces:**
- Consumes: Task 5's `insertDocument`, `issueParams`; Task 6's patterns.
- Produces: `Service.CreditNote(ctx, orgID, parentID, issuedBy string, in CreditInput) (*Invoice, error)`; `CreditInput{Reason string; Items []ItemInput}`.

- [ ] **Step 1: Write the failing test**

```go
package e2e

import (
	"context"
	"net/http"
	"testing"
)

func creditBody(qty, price float64) map[string]any {
	return map[string]any{
		"reason": "member left mid-month",
		"items": []map[string]any{
			{"description": "Refund: Monthly Boxing", "quantity": qty, "unit_price": price},
		},
	}
}

func TestInvoiceCredit_ReversesAndLinks(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000501", "Admin")
	orgID := createTestOrg(t, admin, "Credit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 3000), token)
	assertStatus(t, resp, http.StatusCreated)
	var parent struct {
		ID       string `json:"id"`
		Sequence int    `json:"sequence"`
	}
	parseJSON(t, resp, &parent)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note",
		creditBody(1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)

	var note struct {
		DocType       string  `json:"doc_type"`
		CreditNoteFor string  `json:"credit_note_for"`
		Total         float64 `json:"total"`
		Sequence      int     `json:"sequence"`
	}
	parseJSON(t, resp, &note)

	if note.DocType != "credit_note" {
		t.Errorf("doc_type = %q, want credit_note", note.DocType)
	}
	if note.CreditNoteFor != parent.ID {
		t.Errorf("credit_note_for = %q, want the parent invoice", note.CreditNoteFor)
	}
	if note.Total != 1000 {
		t.Errorf("total = %.2f, want 1000", note.Total)
	}
	// Credit notes share the invoice sequence — one unbroken run per year.
	if note.Sequence != parent.Sequence+1 {
		t.Errorf("sequence = %d, want %d", note.Sequence, parent.Sequence+1)
	}

	// A reversing ledger row must exist for the credited amount.
	var refunds int
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'expense'
		    AND category = 'refund' AND amount = 1000`, orgID).Scan(&refunds)
	if err != nil {
		t.Fatalf("counting refunds: %v", err)
	}
	if refunds != 1 {
		t.Errorf("refund ledger rows = %d, want exactly 1", refunds)
	}
}

func TestInvoiceCredit_CannotExceedBalance(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000502", "Admin")
	orgID := createTestOrg(t, admin, "Overcredit Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	// Credit 800 of 1000 — fine.
	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 800), token)
	assertStatus(t, resp, http.StatusCreated)

	// Another 400 would exceed the remaining 200.
	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 400), token)
	assertStatus(t, resp, http.StatusBadRequest)

	var body struct{ Code string }
	parseJSON(t, resp, &body)
	if body.Code != "credit_exceeds_balance" {
		t.Errorf("code = %q, want credit_exceeds_balance", body.Code)
	}
}

func TestInvoiceCredit_RefusesCancelledParent(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000503", "Admin")
	orgID := createTestOrg(t, admin, "Cancelled Parent Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var parent struct{ ID string }
	parseJSON(t, resp, &parent)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/cancel",
		map[string]any{"reason": "error"}, token)
	assertStatus(t, resp, http.StatusOK)

	resp = doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/invoices/"+parent.ID+"/credit-note", creditBody(1, 500), token)
	assertStatus(t, resp, http.StatusConflict)
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./tests/e2e/ -run 'TestInvoiceCredit' -v`
Expected: FAIL — 501 from the stub.

- [ ] **Step 3: Implement the repository method**

```go
// creditedSoFar totals the credit notes already raised against an invoice.
func creditedSoFar(ctx context.Context, tx pgx.Tx, parentID string) (float64, error) {
	var sum float64
	err := tx.QueryRow(ctx,
		`SELECT COALESCE(SUM(total), 0) FROM invoices
		  WHERE credit_note_for = $1 AND status = 'issued'`, parentID).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("totalling existing credit notes: %w", err)
	}
	return sum, nil
}

// CreditNote raises a credit note against parentID and writes the reversing
// ledger row, both inside one transaction.
func (r *Repository) CreditNote(ctx context.Context, p issueParams, parentID string) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	parent, err := getInTx(ctx, tx, p.OrgID, parentID)
	if err != nil {
		return nil, err
	}
	if parent.Status == "cancelled" {
		return nil, ErrInvoiceCancelled
	}
	if parent.DocType != "invoice" {
		return nil, ErrInvalidParent
	}

	already, err := creditedSoFar(ctx, tx, parentID)
	if err != nil {
		return nil, err
	}
	if round2(already+p.Total) > parent.Total {
		return nil, ErrCreditTooLarge
	}

	seller, err := r.loadSeller(ctx, tx, p.OrgID)
	if err != nil {
		return nil, err
	}

	p.CreditNoteFor = parentID
	id, number, _, err := insertDocument(ctx, tx, p, seller)
	if err != nil {
		return nil, err
	}

	// A refund is a real movement of money, so it reverses in the ledger
	// regardless of who wrote the original income row.
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (organization_id, category, description, transaction_date,
		                           transaction_type, amount, payment_type, reference, entry_by)
		 VALUES ($1, 'refund', $2, CURRENT_DATE, 'expense', $3, $4, $5, $6)`,
		p.OrgID,
		"Credit note "+number+" against "+parent.InvoiceNumber,
		p.Total,
		defaultTo(parent.PaymentMethod, "cash"),
		number,
		p.IssuedBy,
	); err != nil {
		return nil, fmt.Errorf("writing refund ledger row: %w", err)
	}

	note, err := getInTx(ctx, tx, p.OrgID, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing credit note: %w", err)
	}
	return note, nil
}

func defaultTo(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
```

- [ ] **Step 4: Implement the service method**

```go
// CreditInput is a request to credit part or all of an invoice.
type CreditInput struct {
	Reason string      `json:"reason"`
	Items  []ItemInput `json:"items"`
}

// CreditNote reverses part or all of an issued bill. Unlike cancelling, this
// writes a reversing row to the ledger: the money genuinely goes back.
func (s *Service) CreditNote(ctx context.Context, orgID, parentID, issuedBy string, in CreditInput) (*Invoice, error) {
	if strings.TrimSpace(in.Reason) == "" {
		return nil, ErrReasonRequired
	}

	parent, err := s.repo.Get(ctx, orgID, parentID)
	if err != nil {
		return nil, err
	}

	issue := IssueInput{
		CustomerUserID:  parent.CustomerUserID,
		CustomerName:    parent.CustomerName,
		CustomerPAN:     parent.CustomerPAN,
		CustomerAddress: parent.CustomerAddress,
		CustomerPhone:   parent.CustomerPhone,
		PaymentMethod:   parent.PaymentMethod,
		Items:           in.Items,
	}
	subtotal, taxable, err := validateIssue(issue)
	if err != nil {
		return nil, err
	}

	bs, err := nepalidate.Today()
	if err != nil {
		return nil, fmt.Errorf("determining Nepali date: %w", err)
	}
	loc, _ := time.LoadLocation("Asia/Kathmandu")

	note, err := s.repo.CreditNote(ctx, issueParams{
		OrgID:         orgID,
		FiscalYear:    bs.FiscalYear(),
		DocType:       "credit_note",
		IssuedDate:    time.Now().In(loc).Format("2006-01-02"),
		IssuedDateBS:  bs.String(),
		Subtotal:      subtotal,
		TaxableAmount: taxable,
		Total:         taxable,
		AmountInWords: moneywords.Rupees(taxable),
		IssuedBy:      issuedBy,
		In:            issue,
	}, parentID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("credit note raised",
		"credit_note", note.InvoiceNumber, "against", parent.InvoiceNumber,
		"org_id", orgID, "amount", note.Total, "reason", in.Reason)
	return note, nil
}
```

- [ ] **Step 5: Replace the handler stub**

```go
// CreditNote handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/credit-note
func (h *Handler) CreditNote(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var in CreditInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	note, err := h.service.CreditNote(r.Context(), orgID, chi.URLParam(r, "invoiceId"), userID, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}
```

- [ ] **Step 6: Run to verify it passes**

Run: `go test ./tests/e2e/ -run 'TestInvoiceCredit' -v`
Expected: PASS, all three.

- [ ] **Step 7: Commit**

```bash
git add internal/invoice/ tests/e2e/invoice_credit_test.go
git commit -m "$(cat <<'EOF'
feat: raise credit notes against issued invoices

Once a bill reaches the customer it cannot be cancelled — the lawful
correction is a credit note referencing it. Credit notes draw from the same
number sequence as invoices so the year's run stays unbroken, and are capped
at the parent's uncredited remainder.

Unlike cancelling, this writes a reversing ledger row: a refund is a real
movement of money, and leaving the original income on the books would
overstate takings.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: Print audit and next-number preview

**Files:**
- Modify: `internal/invoice/repository.go`, `service.go`, `handler.go`
- Test: `tests/e2e/invoice_print_test.go`

**Interfaces:**
- Consumes: Task 5's `Repository`.
- Produces: `Service.RecordPrint(ctx, orgID, id, printedBy string) (*Invoice, string, error)` returning the invoice and the copy label; `Service.NextNumber(ctx, orgID string) (string, error)`.

- [ ] **Step 1: Write the failing test**

```go
package e2e

import (
	"net/http"
	"testing"
)

func TestInvoicePrint_FirstIsOriginalThenCopy(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000601", "Admin")
	orgID := createTestOrg(t, admin, "Print Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	var out struct {
		CopyLabel string `json:"copy_label"`
		Invoice   struct {
			PrintCount int `json:"print_count"`
		} `json:"invoice"`
	}

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/print", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &out)
	if out.CopyLabel != "original" {
		t.Errorf("first print label = %q, want original", out.CopyLabel)
	}
	if out.Invoice.PrintCount != 1 {
		t.Errorf("print_count = %d, want 1", out.Invoice.PrintCount)
	}

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/print", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &out)
	if out.CopyLabel != "copy" {
		t.Errorf("second print label = %q, want copy", out.CopyLabel)
	}
	if out.Invoice.PrintCount != 2 {
		t.Errorf("print_count = %d, want 2", out.Invoice.PrintCount)
	}
}

func TestInvoiceNextNumber_PreviewsWithoutConsuming(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000602", "Admin")
	orgID := createTestOrg(t, admin, "Preview Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	var preview struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/invoices/next-number", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &preview)

	// Previewing twice must return the same number — it reserves nothing.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/invoices/next-number", nil, token)
	var second struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	parseJSON(t, resp, &second)
	if preview.InvoiceNumber != second.InvoiceNumber {
		t.Errorf("preview consumed a number: %q then %q", preview.InvoiceNumber, second.InvoiceNumber)
	}

	// And the bill actually issued must take that number.
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	parseJSON(t, resp, &inv)
	if inv.InvoiceNumber != preview.InvoiceNumber {
		t.Errorf("issued %q but preview promised %q", inv.InvoiceNumber, preview.InvoiceNumber)
	}
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./tests/e2e/ -run 'TestInvoicePrint|TestInvoiceNextNumber' -v`
Expected: FAIL — 501 from the stubs.

- [ ] **Step 3: Implement the repository methods**

```go
// RecordPrint logs a print and returns the label the document should carry.
func (r *Repository) RecordPrint(ctx context.Context, orgID, id, printedBy string) (*Invoice, string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	inv, err := getInTx(ctx, tx, orgID, id)
	if err != nil {
		return nil, "", err
	}

	label := "original"
	if inv.PrintCount > 0 {
		label = "copy"
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO invoice_prints (invoice_id, printed_by, copy_label) VALUES ($1, $2, $3)`,
		id, printedBy, label); err != nil {
		return nil, "", fmt.Errorf("logging print: %w", err)
	}
	if _, err := tx.Exec(ctx,
		`UPDATE invoices SET print_count = print_count + 1 WHERE id = $1`, id); err != nil {
		return nil, "", fmt.Errorf("incrementing print count: %w", err)
	}

	updated, err := getInTx(ctx, tx, orgID, id)
	if err != nil {
		return nil, "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, "", fmt.Errorf("committing print: %w", err)
	}
	return updated, label, nil
}

// PeekNextSequence reports the number the next bill will take without
// reserving it. A read, never a write — previewing must not consume.
func (r *Repository) PeekNextSequence(ctx context.Context, orgID, fiscalYear string) (int, error) {
	var seq int
	err := r.db.QueryRow(ctx,
		`SELECT next_sequence FROM invoice_counters
		  WHERE organization_id = $1 AND fiscal_year = $2`, orgID, fiscalYear).Scan(&seq)
	if errors.Is(err, pgx.ErrNoRows) {
		return 1, nil
	}
	if err != nil {
		return 0, fmt.Errorf("reading next invoice number: %w", err)
	}
	return seq, nil
}
```

- [ ] **Step 4: Implement the service methods**

```go
// RecordPrint logs that a bill was printed and returns the copy label.
func (s *Service) RecordPrint(ctx context.Context, orgID, id, printedBy string) (*Invoice, string, error) {
	inv, label, err := s.repo.RecordPrint(ctx, orgID, id, printedBy)
	if err != nil {
		return nil, "", err
	}
	s.logger.Info("invoice printed",
		"invoice_number", inv.InvoiceNumber, "org_id", orgID,
		"printed_by", printedBy, "label", label)
	return inv, label, nil
}

// NextNumber previews the next bill's number without consuming it.
func (s *Service) NextNumber(ctx context.Context, orgID string) (string, error) {
	bs, err := nepalidate.Today()
	if err != nil {
		return "", fmt.Errorf("determining Nepali date: %w", err)
	}
	fy := bs.FiscalYear()
	seq, err := s.repo.PeekNextSequence(ctx, orgID, fy)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%06d", fy, seq), nil
}
```

- [ ] **Step 5: Replace the handler stubs**

```go
// Print handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/print
func (h *Handler) Print(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	inv, label, err := h.service.RecordPrint(r.Context(), orgID, chi.URLParam(r, "invoiceId"), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"invoice": inv, "copy_label": label})
}

// NextNumber handles GET /api/v1/orgs/{orgId}/invoices/next-number
func (h *Handler) NextNumber(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	number, err := h.service.NextNumber(r.Context(), orgID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"invoice_number": number})
}
```

- [ ] **Step 6: Run the whole Go suite**

Run: `go test ./internal/... ./pkg/... && go test ./tests/e2e/ -run 'TestInvoice|TestOrgSettings' -v`
Expected: PASS throughout.

- [ ] **Step 7: Commit**

```bash
git add internal/invoice/ tests/e2e/invoice_print_test.go
git commit -m "$(cat <<'EOF'
feat: log invoice prints and preview the next number

Reprints have to be traceable, so every print is logged with who and when,
and the document is marked Original or Copy from the count.

The next-number preview is a read and never a write: reserving a number just
to show it would burn one every time somebody opened the form and abandoned
it, leaving exactly the gaps the counter design exists to prevent.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: Web — types, API client, translations

**Files:**
- Modify: `web/src/types/index.ts`, `web/src/lib/api.ts`, `web/src/lib/i18n.ts`
- Test: `web/tests/i18n-parity.spec.ts` (existing, must keep passing)

**Interfaces:**
- Consumes: the JSON shapes from Tasks 5–8.
- Produces: TS types `Invoice`, `InvoiceItem`, `InvoiceItemInput`, `IssueInvoiceInput`; client methods `listInvoices`, `getInvoice`, `issueInvoice`, `cancelInvoice`, `creditNote`, `printInvoice`, `nextInvoiceNumber`; dictionary key `invoices` in `en` and `ne`.

- [ ] **Step 1: Add the types**

Append to `web/src/types/index.ts`:

```ts
export interface InvoiceItem {
  line_no: number;
  description: string;
  description_ne?: string;
  quantity: number;
  unit_price: number;
  amount: number;
}

export interface Invoice {
  id: string;
  organization_id: string;
  fiscal_year: string;
  sequence: number;
  invoice_number: string;
  doc_type: "invoice" | "credit_note";
  credit_note_for?: string;

  seller_name: string;
  seller_pan: string;
  seller_address?: string;
  seller_vat_registered: boolean;

  customer_user_id?: string;
  customer_name: string;
  customer_pan?: string;
  customer_address?: string;
  customer_phone?: string;

  issued_date: string;
  issued_date_bs: string;

  subtotal: number;
  discount: number;
  taxable_amount: number;
  vat_rate: number;
  vat_amount: number;
  total: number;
  amount_in_words: string;

  payment_method?: string;
  status: "issued" | "cancelled";
  cancelled_at?: string;
  cancellation_reason?: string;

  transaction_id?: string;
  member_package_id?: string;

  issued_by: string;
  print_count: number;
  created_at: string;

  items?: InvoiceItem[];
}

export interface InvoiceItemInput {
  description: string;
  description_ne?: string;
  quantity: number;
  unit_price: number;
}

export interface IssueInvoiceInput {
  customer_user_id?: string;
  customer_name: string;
  customer_pan?: string;
  customer_address?: string;
  customer_phone?: string;
  payment_method?: string;
  transaction_id?: string;
  member_package_id?: string;
  discount?: number;
  items: InvoiceItemInput[];
}
```

Also extend the existing `UpdateOrganizationRequest` (around line 539) and the
`Organization` interface with the tax fields, or Task 10 will not typecheck:

```ts
  pan_number?: string;
  is_vat_registered?: boolean;
  tax_legal_name?: string;
  tax_address?: string;
```

- [ ] **Step 2: Add the client methods**

Add inside the `ApiClient` class in `web/src/lib/api.ts`, next to the guide methods:

```ts
  async listInvoices(
    orgId: string,
    params: {
      status?: string;
      fiscal_year?: string;
      q?: string;
      limit?: number;
      offset?: number;
    } = {}
  ): Promise<{ data: Invoice[]; total: number }> {
    const q = new URLSearchParams();
    if (params.status) q.set("status", params.status);
    if (params.fiscal_year) q.set("fiscal_year", params.fiscal_year);
    if (params.q) q.set("q", params.q);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/invoices${qs ? `?${qs}` : ""}`);
  }

  async getInvoice(orgId: string, id: string): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}`);
  }

  async issueInvoice(orgId: string, data: IssueInvoiceInput): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async cancelInvoice(orgId: string, id: string, reason: string): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/cancel`, {
      method: "POST",
      body: JSON.stringify({ reason }),
    });
  }

  async creditNote(
    orgId: string,
    id: string,
    data: { reason: string; items: InvoiceItemInput[] }
  ): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/credit-note`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async printInvoice(
    orgId: string,
    id: string
  ): Promise<{ invoice: Invoice; copy_label: "original" | "copy" }> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/print`, {
      method: "POST",
    });
  }

  async nextInvoiceNumber(orgId: string): Promise<{ invoice_number: string }> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/next-number`);
  }
```

Add `Invoice`, `InvoiceItemInput`, `IssueInvoiceInput` to the existing type
import at the top of the file.

- [ ] **Step 3: Add the translations**

Add to the `en` dictionary in `web/src/lib/i18n.ts`:

```ts
    invoices: {
      title: "Bills",
      subtitle: "PAN tax invoices issued to customers.",
      newBill: "New bill",
      empty: "No bills yet.",
      number: "Bill no.",
      customer: "Customer",
      date: "Date",
      amount: "Amount",
      status: "Status",
      issued: "Issued",
      cancelled: "Cancelled",
      creditNote: "Credit note",
      print: "Print",
      revise: "Revise",
      cancel: "Cancel bill",
      cancelReason: "Why is this being cancelled?",
      cancelConfirm: "Cancel this bill? The number stays used and cannot be reissued.",
      creditReason: "Why is this being credited?",
      creditAmount: "Amount to credit",
      panMissing: "Set your PAN number before issuing bills.",
      panMissingCta: "Go to settings",
      description: "Description",
      quantity: "Qty",
      unitPrice: "Rate",
      lineTotal: "Amount",
      addLine: "Add line",
      removeLine: "Remove line",
      subtotal: "Subtotal",
      discount: "Discount",
      total: "Total",
      amountInWords: "In words",
      paymentMethod: "Payment method",
      customerName: "Customer name",
      customerPan: "Customer PAN (optional)",
      customerAddress: "Address (optional)",
      walkIn: "Walk-in customer",
      existingMember: "Existing member",
      original: "Original",
      copy: "Copy",
      reviseExplainCancel:
        "This bill has not been printed, so it can be cancelled and reissued.",
      reviseExplainCredit:
        "This bill has already been printed, so it must be corrected with a credit note.",
    },
```

Add the same keys to `ne` with Nepali values:

```ts
    invoices: {
      title: "बिलहरू",
      subtitle: "ग्राहकलाई जारी गरिएका प्यान कर बिजक।",
      newBill: "नयाँ बिल",
      empty: "अहिलेसम्म कुनै बिल छैन।",
      number: "बिल नं.",
      customer: "ग्राहक",
      date: "मिति",
      amount: "रकम",
      status: "स्थिति",
      issued: "जारी",
      cancelled: "रद्द",
      creditNote: "क्रेडिट नोट",
      print: "प्रिन्ट",
      revise: "सच्याउनुहोस्",
      cancel: "बिल रद्द गर्नुहोस्",
      cancelReason: "किन रद्द गर्दै हुनुहुन्छ?",
      cancelConfirm: "यो बिल रद्द गर्ने? नम्बर प्रयोग भइसकेको रहन्छ र पुनः जारी हुँदैन।",
      creditReason: "किन क्रेडिट दिँदै हुनुहुन्छ?",
      creditAmount: "क्रेडिट रकम",
      panMissing: "बिल जारी गर्नुअघि आफ्नो प्यान नम्बर राख्नुहोस्।",
      panMissingCta: "सेटिङमा जानुहोस्",
      description: "विवरण",
      quantity: "परिमाण",
      unitPrice: "दर",
      lineTotal: "रकम",
      addLine: "पङ्क्ति थप्नुहोस्",
      removeLine: "पङ्क्ति हटाउनुहोस्",
      subtotal: "जम्मा",
      discount: "छुट",
      total: "कुल",
      amountInWords: "अक्षरमा",
      paymentMethod: "भुक्तानी माध्यम",
      customerName: "ग्राहकको नाम",
      customerPan: "ग्राहकको प्यान (वैकल्पिक)",
      customerAddress: "ठेगाना (वैकल्पिक)",
      walkIn: "आगन्तुक ग्राहक",
      existingMember: "पुरानो सदस्य",
      original: "सक्कल",
      copy: "प्रतिलिपि",
      reviseExplainCancel:
        "यो बिल प्रिन्ट भएको छैन, त्यसैले रद्द गरी पुनः जारी गर्न सकिन्छ।",
      reviseExplainCredit:
        "यो बिल प्रिन्ट भइसकेको छ, त्यसैले क्रेडिट नोटबाट मात्र सच्याउन मिल्छ।",
    },
```

Also add `invoices: "Bills"` to `en.nav` and `invoices: "बिलहरू"` to `ne.nav`.

- [ ] **Step 4: Verify parity and types**

Run: `cd web && npx tsc --noEmit && npx playwright test tests/i18n-parity.spec.ts`
Expected: zero TypeScript errors; parity spec PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/types/index.ts web/src/lib/api.ts web/src/lib/i18n.ts
git commit -m "$(cat <<'EOF'
feat: add invoice types, API client and translations

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: Web — PAN in settings

**Files:**
- Modify: `web/src/app/[lang]/dashboard/settings/page.tsx`
- Test: `web/tests/invoices.spec.ts` (create)

**Interfaces:**
- Consumes: Task 4's org update endpoint, Task 9's translations.
- Produces: a Tax section that saves `pan_number`, `tax_legal_name`, `tax_address`.

- [ ] **Step 1: Write the failing test**

Create `web/tests/invoices.spec.ts`. These specs never reach a real backend —
`web/tests/` mocks the API with `page.route` and fakes the session with
`injectAuth`, which decodes into `localStorage` without server verification.
Note the ordering: `goto` first (so there is a document to evaluate against),
then `injectAuth`, then the route mocks.

```ts
import { test, expect, type Page, type Route } from "@playwright/test";
import { injectAuth } from "./helpers";

const ORG = {
  id: "org-001",
  name: "Test Gym",
  slug: "test-gym",
  pan_number: "",
  tax_legal_name: "",
  tax_address: "",
};

function mockOrgAPI(page: Page) {
  let org = { ...ORG };
  return page.route("**/api/v1/orgs/org-001", (route: Route) => {
    if (route.request().method() === "PUT") {
      org = { ...org, ...route.request().postDataJSON() };
    }
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(org),
    });
  });
}

test.describe("Tax settings", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockOrgAPI(page);
  });

  test("rejects a PAN that is not nine digits", async ({ page }) => {
    await page.goto("/en/dashboard/settings");

    const pan = page.getByLabel(/PAN/i);
    await expect(pan).toBeVisible();

    await pan.fill("12345");
    await page.getByRole("button", { name: /save/i }).click();
    await expect(page.getByText(/9 digits/i)).toBeVisible();
  });

  test("accepts a nine-digit PAN", async ({ page }) => {
    await page.goto("/en/dashboard/settings");

    await page.getByLabel(/PAN/i).fill("601234567");
    await page.getByRole("button", { name: /save/i }).click();
    await expect(page.getByText(/9 digits/i)).toHaveCount(0);
  });
});
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd web && npx playwright test tests/invoices.spec.ts`
Expected: FAIL — no PAN field on the settings page.

- [ ] **Step 3: Add the Tax section**

In the settings page, add a section alongside the existing ones. Validate on
the client before calling the API so the user sees the problem immediately:

```tsx
const [pan, setPan] = useState(org?.pan_number ?? "");
const [taxLegalName, setTaxLegalName] = useState(org?.tax_legal_name ?? "");
const [taxAddress, setTaxAddress] = useState(org?.tax_address ?? "");
const [panError, setPanError] = useState<string | null>(null);

async function saveTax(e: React.FormEvent) {
  e.preventDefault();
  setPanError(null);
  if (pan && !/^\d{9}$/.test(pan)) {
    setPanError(t.invoices.panDigits);
    return;
  }
  await api.updateOrg(orgId, {
    pan_number: pan,
    tax_legal_name: taxLegalName,
    tax_address: taxAddress,
  });
}
```

Add `panDigits: "PAN must be exactly 9 digits"` / `"प्यान ठ्याक्कै ९ अंकको हुनुपर्छ"`
to both dictionaries. Render the three inputs with a `<label htmlFor>` on each
so the Playwright `getByLabel` selector resolves.

- [ ] **Step 4: Run to verify it passes**

Run: `cd web && npx playwright test tests/invoices.spec.ts`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/app/\[lang\]/dashboard/settings/page.tsx web/src/lib/i18n.ts web/tests/invoices.spec.ts
git commit -m "$(cat <<'EOF'
feat: let a gym enter its PAN and registered tax identity

Without a PAN no bill can be issued, so this is the first thing a gym needs
and the API returns a distinct pan_not_configured code to send them here.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: Web — the printable document

**Files:**
- Create: `web/src/components/invoice/InvoiceDocument.tsx`
- Create: `web/src/components/invoice/invoice-print.css`

**Interfaces:**
- Consumes: Task 9's `Invoice` type and translations.
- Produces: `<InvoiceDocument invoice={inv} locale={locale} t={t} copyLabel="original" />`.

- [ ] **Step 1: Build the component**

This renders the legal document. Every field IRD expects appears here and comes
from the invoice's own snapshot, never from live org or member data.

```tsx
"use client";

import type { Dictionary } from "@/lib/i18n";
import type { Invoice, Locale } from "@/types";

/**
 * The printed bill.
 *
 * Every value is read from the invoice's own snapshot rather than from the
 * org or the member: a bill reprinted next year must show what it showed the
 * day it was issued, even if the gym has since changed its name or PAN.
 */
export function InvoiceDocument({
  invoice,
  locale,
  t,
  copyLabel,
}: {
  invoice: Invoice;
  locale: Locale;
  t: Dictionary;
  copyLabel?: "original" | "copy";
}) {
  const isCredit = invoice.doc_type === "credit_note";
  const money = (n: number) =>
    n.toLocaleString(locale === "ne" ? "ne-NP" : "en-IN", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

  return (
    <article className="invoice-doc bg-white text-black p-8 max-w-[210mm] mx-auto">
      <header className="text-center border-b border-black/20 pb-4">
        <h1 className="text-lg font-bold">{invoice.seller_name}</h1>
        {invoice.seller_address && <p className="text-sm">{invoice.seller_address}</p>}
        <p className="text-sm font-mono">PAN: {invoice.seller_pan}</p>
        <p className="mt-2 text-sm font-semibold uppercase tracking-widest">
          {isCredit ? t.invoices.creditNote : t.invoices.title}
          {copyLabel && (
            <span className="ml-2 font-normal">
              ({copyLabel === "original" ? t.invoices.original : t.invoices.copy})
            </span>
          )}
        </p>
      </header>

      <section className="grid grid-cols-2 gap-4 py-4 text-sm">
        <div>
          <p>
            <span className="text-black/60">{t.invoices.customer}: </span>
            {invoice.customer_name}
          </p>
          {invoice.customer_pan && (
            <p>
              <span className="text-black/60">PAN: </span>
              {invoice.customer_pan}
            </p>
          )}
          {invoice.customer_address && <p>{invoice.customer_address}</p>}
        </div>
        <div className="text-right">
          <p className="font-mono">{invoice.invoice_number}</p>
          <p>{invoice.issued_date_bs} (BS)</p>
          <p className="text-black/60">{invoice.issued_date} (AD)</p>
        </div>
      </section>

      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-y border-black/20">
            <th className="text-left py-1 w-8">#</th>
            <th className="text-left py-1">{t.invoices.description}</th>
            <th className="text-right py-1 w-16">{t.invoices.quantity}</th>
            <th className="text-right py-1 w-24">{t.invoices.unitPrice}</th>
            <th className="text-right py-1 w-28">{t.invoices.lineTotal}</th>
          </tr>
        </thead>
        <tbody>
          {(invoice.items ?? []).map((it) => (
            <tr key={it.line_no} className="border-b border-black/10">
              <td className="py-1">{it.line_no}</td>
              <td className="py-1">
                {locale === "ne" && it.description_ne ? it.description_ne : it.description}
              </td>
              <td className="py-1 text-right tabular-nums">{it.quantity}</td>
              <td className="py-1 text-right tabular-nums">{money(it.unit_price)}</td>
              <td className="py-1 text-right tabular-nums">{money(it.amount)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <section className="mt-4 ml-auto w-64 text-sm">
        <Row label={t.invoices.subtotal} value={money(invoice.subtotal)} />
        {invoice.discount > 0 && (
          <Row label={t.invoices.discount} value={`- ${money(invoice.discount)}`} />
        )}
        {/* VAT is deferred; the row appears only if a bill ever carries one. */}
        {invoice.vat_amount > 0 && (
          <Row label={`VAT ${invoice.vat_rate}%`} value={money(invoice.vat_amount)} />
        )}
        <div className="flex justify-between border-t border-black/30 pt-1 font-semibold">
          <span>{t.invoices.total}</span>
          <span className="tabular-nums">Rs {money(invoice.total)}</span>
        </div>
      </section>

      <p className="mt-3 text-sm">
        <span className="text-black/60">{t.invoices.amountInWords}: </span>
        {invoice.amount_in_words}
      </p>

      {invoice.status === "cancelled" && (
        <p className="mt-4 border border-black px-3 py-2 text-sm font-semibold uppercase">
          {t.invoices.cancelled}
          {invoice.cancellation_reason ? ` — ${invoice.cancellation_reason}` : ""}
        </p>
      )}
    </article>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between">
      <span className="text-black/60">{label}</span>
      <span className="tabular-nums">{value}</span>
    </div>
  );
}
```

- [ ] **Step 2: Add the print stylesheet**

```css
@media print {
  /* Only the document reaches paper. */
  body * { visibility: hidden; }
  .invoice-doc, .invoice-doc * { visibility: visible; }
  .invoice-doc {
    position: absolute;
    inset: 0;
    margin: 0;
    max-width: none;
    padding: 12mm;
  }
  @page { size: A4; margin: 0; }
}
```

Import it from the invoice detail page.

- [ ] **Step 3: Verify it compiles**

Run: `cd web && npx tsc --noEmit`
Expected: zero errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/components/invoice/
git commit -m "$(cat <<'EOF'
feat: render the printable tax invoice

Every value comes from the invoice's own snapshot rather than live org or
member data, so a reprint next year shows what the bill showed the day it was
issued.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: Web — list, issue, and detail screens

**Files:**
- Create: `web/src/app/[lang]/dashboard/invoices/page.tsx`
- Create: `web/src/app/[lang]/dashboard/invoices/new/page.tsx`
- Create: `web/src/app/[lang]/dashboard/invoices/[id]/page.tsx`
- Modify: `web/src/components/layout/Sidebar.tsx`
- Test: `web/tests/invoices.spec.ts` (extend)

**Interfaces:**
- Consumes: Tasks 9 and 11.
- Produces: the three screens and a sidebar entry.

- [ ] **Step 1: Write the failing tests**

Append to `web/tests/invoices.spec.ts`:

```ts
function baseInvoice(overrides: Record<string, unknown> = {}) {
  return {
    id: "inv-001",
    organization_id: "org-001",
    fiscal_year: "2082-83",
    sequence: 1,
    invoice_number: "2082-83/000001",
    doc_type: "invoice",
    seller_name: "Test Gym",
    seller_pan: "601234567",
    seller_vat_registered: false,
    customer_name: "Ram Bahadur",
    issued_date: "2025-07-20",
    issued_date_bs: "2082-04-04",
    subtotal: 3000,
    discount: 0,
    taxable_amount: 3000,
    vat_rate: 0,
    vat_amount: 0,
    total: 3000,
    amount_in_words: "Three thousand rupees only",
    status: "issued",
    issued_by: "test-user-001",
    print_count: 0,
    created_at: "2025-07-20T04:00:00Z",
    items: [
      {
        line_no: 1,
        description: "Monthly Boxing",
        quantity: 1,
        unit_price: 3000,
        amount: 3000,
      },
    ],
    ...overrides,
  };
}

function mockInvoiceAPI(page: Page) {
  let invoice = baseInvoice();
  return page.route("**/api/v1/orgs/*/invoices**", (route: Route) => {
    const url = new URL(route.request().url());
    const method = route.request().method();

    if (url.pathname.endsWith("/next-number")) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ invoice_number: "2082-83/000001" }),
      });
    }

    if (url.pathname.endsWith("/cancel")) {
      const { reason } = route.request().postDataJSON();
      if (invoice.status === "cancelled") {
        return route.fulfill({
          status: 409,
          contentType: "application/json",
          body: JSON.stringify({ error: "already cancelled", code: "already_cancelled" }),
        });
      }
      invoice = baseInvoice({ status: "cancelled", cancellation_reason: reason });
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(invoice),
      });
    }

    if (method === "POST" && url.pathname.endsWith("/invoices")) {
      const body = route.request().postDataJSON();
      invoice = baseInvoice({ customer_name: body.customer_name });
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(invoice),
      });
    }

    if (method === "GET" && url.pathname.endsWith("/invoices")) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: [invoice], total: 1 }),
      });
    }

    // GET one
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(invoice),
    });
  });
}

test.describe("Bills", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockInvoiceAPI(page);
  });

  test("issuing a bill lands on its document", async ({ page }) => {
    await page.goto("/en/dashboard/invoices/new");

    await page.getByLabel(/customer name/i).fill("Ram Bahadur");
    await page.getByLabel(/description/i).first().fill("Monthly Boxing");
    await page.getByLabel(/rate/i).first().fill("3000");
    await page.getByRole("button", { name: /issue/i }).click();

    await expect(page).toHaveURL(/\/dashboard\/invoices\/inv-001/);
    await expect(page.getByText("Ram Bahadur")).toBeVisible();
    await expect(page.getByText("2082-83/000001")).toBeVisible();
  });

  test("a cancelled bill is marked and offers no second cancel", async ({ page }) => {
    await page.goto("/en/dashboard/invoices/inv-001");

    await page.getByRole("button", { name: /cancel bill/i }).click();
    await page.getByLabel(/why is this being cancelled/i).fill("wrong customer");
    await page.getByRole("button", { name: /confirm/i }).click();

    await expect(page.getByText(/cancelled/i).first()).toBeVisible();
    await expect(page.getByRole("button", { name: /cancel bill/i })).toHaveCount(0);
  });

  test("the list shows the bill number and customer", async ({ page }) => {
    await page.goto("/en/dashboard/invoices");
    await expect(page.getByText("2082-83/000001")).toBeVisible();
    await expect(page.getByText("Ram Bahadur")).toBeVisible();
  });
});
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd web && npx playwright test tests/invoices.spec.ts`
Expected: FAIL — the routes 404.

- [ ] **Step 3: Build the list page**

`web/src/app/[lang]/dashboard/invoices/page.tsx` — mirror the structure of
`dashboard/guides/page.tsx`: `useAuth()` for `org_id`, `useCallback` fetch,
loading spinner, error banner, empty state. Columns: number, customer, BS date,
total, status chip. Each row links to `./invoices/{id}`. A "New bill" button
links to `./invoices/new`.

When `listInvoices` throws `ApiRequestError` with status 400, show
`t.invoices.panMissing` and a link to settings instead of a bare error.

- [ ] **Step 4: Build the issue page**

`web/src/app/[lang]/dashboard/invoices/new/page.tsx`:

- Customer: a toggle between `t.invoices.existingMember` (member picker via
  `api.listMembers`) and `t.invoices.walkIn` (free-text name). Selecting a
  member fills `customer_user_id`, name and phone.
- Line items: a repeating row of description / quantity / rate, with add and
  remove. Default one empty row, quantity 1.
- Live totals computed client-side; the server recomputes and is authoritative.
- Show the previewed number from `api.nextInvoiceNumber` as a hint, labelled as
  a preview so nobody treats it as reserved.
- On success, `router.push` to the new bill's detail page.

Every input needs an `htmlFor`-linked `<label>` — the Playwright selectors
depend on it, and it is what a screen reader needs anyway.

- [ ] **Step 5: Build the detail page**

`web/src/app/[lang]/dashboard/invoices/[id]/page.tsx`:

- Render `<InvoiceDocument>`.
- **Print** calls `api.printInvoice`, then `window.print()` with the returned
  `copy_label`.
- **Cancel bill** shows only when `status === "issued"`; opens a reason prompt
  and calls `api.cancelInvoice`.
- **Revise** explains which correction applies and routes to it:

```tsx
// Cancelling is only honest while the bill has not left the counter. Once it
// has been printed the customer may be holding it, so the lawful correction
// is a credit note.
const canCancel = invoice.status === "issued" && invoice.print_count === 0;
```

Show `t.invoices.reviseExplainCancel` or `t.invoices.reviseExplainCredit`
accordingly. There is no edit control anywhere on this page.

- [ ] **Step 6: Add the sidebar entry**

In `web/src/components/layout/Sidebar.tsx`, add to the "Run Operations" group,
after the accounts entry:

```tsx
        { label: t.nav.invoices, href: `${base}/invoices`, icon: <Receipt className="w-4 h-4" /> },
```

Import `Receipt` from `lucide-react`.

- [ ] **Step 7: Run the web tests**

Run: `cd web && npx tsc --noEmit && npx playwright test tests/invoices.spec.ts tests/i18n-parity.spec.ts`
Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add web/src/app/\[lang\]/dashboard/invoices/ web/src/components/layout/Sidebar.tsx web/tests/invoices.spec.ts
git commit -m "$(cat <<'EOF'
feat: add the bills screens

List, issue and view a PAN tax invoice. The detail screen offers Print,
Cancel and Revise, and Revise explains which correction applies rather than
offering a choice: a bill that has been printed may be in the customer's
hands, so it can only be corrected by credit note.

There is no edit control anywhere, by design.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: Web — issue a bill from a package sale

**Files:**
- Modify: `web/src/app/[lang]/dashboard/members/[id]/page.tsx` (or wherever package assignment lands — grep for `assignPackage`)
- Test: `web/tests/invoices.spec.ts` (extend)

**Interfaces:**
- Consumes: Task 12's issue page.
- Produces: an "Issue bill" action that pre-fills the form and links the existing ledger row.

- [ ] **Step 1: Locate the assignment UI**

Run: `grep -rn "assignPackage\|assign_package" web/src --include=*.tsx`
The action belongs next to the result of a successful assignment.

- [ ] **Step 2: Write the failing test**

Append to the `Bills` describe block in `web/tests/invoices.spec.ts`, which
already provides `injectAuth` and `mockInvoiceAPI` in its `beforeEach`:

```ts
  test("a package sale pre-fills the bill form", async ({ page }) => {
    await page.goto(
      "/en/dashboard/invoices/new" +
        "?customer=Ram%20Bahadur&package=Monthly%20Boxing&amount=3000"
    );

    await expect(page.getByLabel(/customer name/i)).toHaveValue("Ram Bahadur");
    await expect(page.getByLabel(/description/i).first()).toHaveValue("Monthly Boxing");
    await expect(page.getByLabel(/rate/i).first()).toHaveValue("3000");
  });
```

- [ ] **Step 3: Read the pre-fill from query params**

In the issue page, seed initial state from `useSearchParams()`:
`customer`, `customer_user_id`, `package` (line description), `amount`
(unit price), `transaction_id`, `member_package_id`.

`transaction_id` is the important one: passing it links the bill to the income
row the package assignment already wrote, so the sale is never counted twice.
The partial unique index rejects a second live bill against the same payment,
which surfaces as `already_billed` — render that as "this payment already has a
bill" with a link to it rather than a raw error.

- [ ] **Step 4: Add the action button**

Where a package assignment succeeds, add a link to
`./invoices/new?customer=…&customer_user_id=…&package=…&amount=…&transaction_id=…&member_package_id=…`.

- [ ] **Step 5: Run the tests**

Run: `cd web && npx playwright test tests/invoices.spec.ts`
Expected: PASS.

- [ ] **Step 6: Full suite**

Run: `go test ./internal/... ./pkg/... && go test ./tests/e2e/... && cd web && npx tsc --noEmit && npx playwright test`
Expected: PASS throughout.

- [ ] **Step 7: Commit**

```bash
git add web/src
git commit -m "$(cat <<'EOF'
feat: issue a bill straight from a package sale

Assigning a package already writes an income row, so the bill links that row
rather than writing its own — a partial unique index makes billing the same
payment twice impossible instead of merely discouraged.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Self-Review Notes

Checked against the spec:

- Org tax settings → Task 4. PAN guard → Tasks 5 (`ErrPANNotConfigured`), 10, 12.
- Gapless numbering → Task 3 (counter table), Task 5 (`allocateSequence`, concurrency test).
- Immutability → Task 3 triggers + schema tests.
- Cancel → Task 6. Credit note → Task 7. Print audit → Task 8.
- Snapshots → Task 5 (`loadSeller`, insert column list), Task 11 (document reads only snapshots).
- Cross-tenant isolation → Task 5 test; every repository read is org-scoped in the `WHERE` clause.
- One bill per payment → Task 3 index, Task 5 `ErrAlreadyBilled`, Task 13 UI handling.
- Ledger effects → Task 7 (credit note writes the refund row; cancel writes nothing).
- BS dates and amount-in-words → Tasks 1 and 2.
- i18n parity → Task 9.

**Deliberately not covered, matching the spec's "Out" column:** VAT arithmetic
and UI, CBMS sync, split payment methods, bulk issue, customer-facing download,
mobile clients.

Corrected during self-review:

- The Playwright tests originally imported a `loginAsAdmin` helper that does not
  exist. `web/tests/helpers.ts` exports `injectAuth`, and the web specs mock the
  API with `page.route` rather than reaching a live backend — Tasks 10, 12 and 13
  were rewritten to match.
- `UpdateOrganizationRequest` and `Organization` need the four tax fields added,
  or Task 10's settings code does not typecheck. Folded into Task 9.
- Verified `middleware.RequireOrgRole` returns 403, which is what
  `TestInvoice_MemberCannotIssue` asserts.

**Known gap carried from the spec:** the fiscal-year rollover test (spec test 5)
cannot run against `nepalidate.Today()`, which reads the real clock. Task 1's
`TestDate_FiscalYear` covers the boundary at the unit level instead. Making the
service clock injectable would be the way to test it end-to-end; that is a
worthwhile follow-up but not a blocker.
