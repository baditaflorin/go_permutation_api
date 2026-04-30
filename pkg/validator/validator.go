package validator

import (
	"fmt"
	"strings"
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
	sanitized := make([]string, len(elements))
	for i, elem := range elements {
		sanitized[i] = strings.TrimSpace(elem)
	}
	return sanitized
}
