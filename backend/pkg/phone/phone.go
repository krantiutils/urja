package phone

import (
	"fmt"
	"strings"
	"unicode"
)

// Normalize takes a phone number string and returns the normalized 10-digit
// Nepal phone number. It strips whitespace, dashes, and the +977 country code.
func Normalize(raw string) (string, error) {
	// Strip all non-digit characters
	var digits strings.Builder
	for _, r := range raw {
		if unicode.IsDigit(r) {
			digits.WriteRune(r)
		}
	}

	number := digits.String()

	// Handle +977 prefix (which becomes 977 after stripping non-digits)
	if strings.HasPrefix(number, "977") && len(number) == 13 {
		number = number[3:]
	}

	if err := Validate(number); err != nil {
		return "", err
	}

	return number, nil
}

// Validate checks that a phone number is a valid Nepal mobile number.
// Expects a 10-digit string starting with 9, 6, 7, or 8.
func Validate(number string) error {
	if len(number) != 10 {
		return fmt.Errorf("phone number must be 10 digits, got %d", len(number))
	}

	for _, r := range number {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("phone number must contain only digits")
		}
	}

	first := number[0]
	if first != '9' && first != '6' && first != '7' && first != '8' {
		return fmt.Errorf("phone number must start with 9, 6, 7, or 8")
	}

	return nil
}
