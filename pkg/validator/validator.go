package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateElements validates the input elements for permutation generation
func ValidateElements(elements []string, maxElements int) error {
	if len(elements) > maxElements {
		return fmt.Errorf("too many elements: maximum is %d, got %d", maxElements, len(elements))
	}

	for i, elem := range elements {
		if strings.TrimSpace(elem) == "" {
			return fmt.Errorf("element at index %d is empty or contains only whitespace", i)
		}
		if strings.IndexFunc(elem, unicode.IsControl) >= 0 {
			return fmt.Errorf("element at index %d contains a control character", i)
		}
	}

	return nil
}

// ValidatePort validates a port number
func ValidatePort(port string) error {
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	// Simple validation - can be enhanced with actual port number checking
	return nil
}

// SanitizeElements removes whitespace from elements
func SanitizeElements(elements []string) []string {
	sanitized := make([]string, 0, len(elements))
	seen := make(map[string]struct{}, len(elements))
	for _, elem := range elements {
		elem = strings.TrimSpace(elem)
		if _, exists := seen[elem]; exists {
			continue
		}
		seen[elem] = struct{}{}
		sanitized = append(sanitized, elem)
	}
	return sanitized
}
