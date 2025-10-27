package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/wonjinsin/go-boilerplate/pkg/constants"
)

var emailRegex = regexp.MustCompile(constants.EmailPattern)

// IsValidEmail checks if the given string is a valid email format
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	if len(email) < constants.MinEmailLength || len(email) > constants.MaxEmailLength {
		return false
	}
	return emailRegex.MatchString(email)
}

// NormalizeEmail normalizes email address to lowercase and trims spaces
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeName normalizes name by trimming spaces and title casing
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Simple title case for names
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// IsEmptyOrWhitespace checks if string is empty or contains only whitespace
func IsEmptyOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
