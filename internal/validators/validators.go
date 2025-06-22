package validators

import (
	"regexp"
	"strings"
	"unicode"
)

// Checks if a string is a valid email format
func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

// Checks if string Length is within min and max bounds (inclusive)
func ValidateLength(s string, min int, max int) bool {
	length := len(s)

	return length >= min && length <= max
}

// checks if a string contains at least one digit
func HasAtLeastOneNumber(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}

	return false
}

// check if a string contains at least one letter
func HasAtLeastOneCharacter(s string) bool {
	for _, char := range s {
		if unicode.IsLetter((char)) {
			return true
		}
	}

	return false
}

func SanitizeInput(s string) string {
	// Remove null bytes and control characters
	sanitized := strings.ReplaceAll(s, "\x00", "")
	sanitized = strings.ReplaceAll(s, " ", "")

	// Remove or replace other problematic characters
	sanitized = strings.ReplaceAll(sanitized, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "&", "&amp;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#39;")

	return strings.TrimSpace(sanitized)
}
