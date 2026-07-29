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
// Defensive lower bound guards against negative indices from out-of-range inputs.
func under1000(n int) string {
	if n <= 0 {
		return ""
	}
	var parts []string
	if h := n / 100; h > 0 {
		parts = append(parts, ones[h], "hundred")
	}
	r := n % 100
	switch {
	case r == 0:
	case r < 0 || r >= 20:
		// Defensive: r >= 20 handles normal path; r < 0 guards against panic
		if r >= 20 {
			parts = append(parts, tens[r/10])
			if r%10 > 0 {
				parts = append(parts, ones[r%10])
			}
		}
	default:
		parts = append(parts, ones[r])
	}
	return strings.Join(parts, " ")
}

// words renders a whole number using kharab/arab/crore/lakh/thousand grouping.
// Stops recursing on the top group to handle kharab and arab.
func words(n int64) string {
	if n == 0 {
		return "zero"
	}
	var parts []string

	// kharab = 1e11, arab = 1e9, crore = 1e7
	if kharab := n / 100000000000; kharab > 0 {
		parts = append(parts, under1000(int(kharab)), "kharab")
		n %= 100000000000
	}
	if arab := n / 1000000000; arab > 0 {
		parts = append(parts, under1000(int(arab)), "arab")
		n %= 1000000000
	}
	if crore := n / 10000000; crore > 0 {
		parts = append(parts, under1000(int(crore)), "crore")
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
// Guards against NaN, Inf, and out-of-range values that could cause panic from
// malformed HTTP input (unbounded JSON floats arriving as unit_price, quantity, etc).
func Rupees(amount float64) string {
	// Guard against NaN and Inf — malformed JSON floats from HTTP clients can produce these.
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return "Invalid amount"
	}

	// Round to paisa (nearest 0.01) before deciding the sign, to avoid "Minus zero".
	// Add small epsilon (1e-9) to handle float64 precision (e.g., 1.005 is stored < 1.005).
	total := int64(math.Round(amount*100 + 1e-9))

	// Guard against out-of-range: if the rounded paisa value doesn't fit in int64,
	// return a sentinel rather than panicking in under1000().
	if total < math.MinInt64/100 || total > math.MaxInt64/100 {
		return "Amount out of range"
	}

	// Now decide the sign from the rounded total.
	negative := total < 0
	if negative {
		total = -total
	}

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
	out = capitalise(out) + " only"

	if negative {
		out = "Minus " + strings.ToLower(out)
	}
	return out
}
