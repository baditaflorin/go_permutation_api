package validator

import (
	"testing"
)

func TestValidateElements(t *testing.T) {
	tests := []struct {
		name        string
		elements    []string
		maxElements int
		wantErr     bool
	}{
		{
			name:        "valid elements",
			elements:    []string{"a", "b", "c"},
			maxElements: 5,
			wantErr:     false,
		},
		{
			name:        "too many elements",
			elements:    []string{"a", "b", "c", "d"},
			maxElements: 3,
			wantErr:     true,
		},
		{
			name:        "empty element",
			elements:    []string{"a", "", "c"},
			maxElements: 5,
			wantErr:     true,
		},
		{
			name:        "whitespace only element",
			elements:    []string{"a", "   ", "c"},
			maxElements: 5,
			wantErr:     true,
		},
		{
			name:        "empty slice",
			elements:    []string{},
			maxElements: 5,
			wantErr:     false,
		},
		{
			name:        "single element",
			elements:    []string{"x"},
			maxElements: 1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateElements(tt.elements, tt.maxElements)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateElements() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{
			name:    "valid port",
			port:    "8080",
			wantErr: false,
		},
		{
			name:    "empty port",
			port:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "trim whitespace",
			input:    []string{" a ", " b", "c "},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "no whitespace",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeElements(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
