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
