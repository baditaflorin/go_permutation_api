package permutation

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "sorted input",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "unsorted input",
			input:    []string{"c", "a", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "single element",
			input:    []string{"x"},
			expected: []string{"x"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.input)
			if !reflect.DeepEqual(gen.elements, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, gen.elements)
			}
			// Verify original slice is not modified
			if len(tt.input) > 0 && &tt.input[0] == &gen.elements[0] {
				t.Error("generator should work on a copy of the input")
			}
		})
	}
}

func TestGeneratorNext(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected [][]string
	}{
		{
			name:  "three elements",
			input: []string{"a", "b", "c"},
			expected: [][]string{
				{"a", "b", "c"},
				{"a", "c", "b"},
				{"b", "a", "c"},
				{"b", "c", "a"},
				{"c", "a", "b"},
				{"c", "b", "a"},
			},
		},
		{
			name:  "two elements",
			input: []string{"x", "y"},
			expected: [][]string{
				{"x", "y"},
				{"y", "x"},
			},
		},
		{
			name:     "single element",
			input:    []string{"z"},
			expected: [][]string{{"z"}},
		},
		{
			name:     "empty",
			input:    []string{},
			expected: [][]string{{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.input)
			var results [][]string

			for {
				perm, ok := gen.Next()
				if !ok {
					break
				}
				results = append(results, perm)
			}

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, results)
			}
		})
	}
}

func TestGenerateAll(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		expectedCount int
	}{
		{
			name:          "three elements",
			input:         []string{"a", "b", "c"},
			expectedCount: 6,
		},
		{
			name:          "four elements",
			input:         []string{"1", "2", "3", "4"},
			expectedCount: 24,
		},
		{
			name:          "single element",
			input:         []string{"x"},
			expectedCount: 1,
		},
		{
			name:          "empty",
			input:         []string{},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := GenerateAll(tt.input)
			if len(results) != tt.expectedCount {
				t.Errorf("expected %d permutations, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
		{-1, 0},
	}

	for _, tt := range tests {
		result := Count(tt.n)
		if result != tt.expected {
			t.Errorf("Count(%d) = %d, expected %d", tt.n, result, tt.expected)
		}
	}
}

func TestCountWithinLimit(t *testing.T) {
	if got, exceeds := CountWithinLimit(5, 200); exceeds || got != 120 {
		t.Fatalf("5! within limit: got=%d exceeds=%t", got, exceeds)
	}
	if _, exceeds := CountWithinLimit(12, 100000); !exceeds {
		t.Fatal("12! must exceed the streamed response limit")
	}
}

func TestGeneratorDoesNotModifyInput(t *testing.T) {
	input := []string{"c", "b", "a"}
	original := make([]string, len(input))
	copy(original, input)

	gen := New(input)
	gen.Next()
	gen.Next()

	if !reflect.DeepEqual(input, original) {
		t.Error("generator modified the original input slice")
	}
}

func BenchmarkGenerateAll3Elements(b *testing.B) {
	elements := []string{"a", "b", "c"}
	for i := 0; i < b.N; i++ {
		GenerateAll(elements)
	}
}

func BenchmarkGenerateAll5Elements(b *testing.B) {
	elements := []string{"a", "b", "c", "d", "e"}
	for i := 0; i < b.N; i++ {
		GenerateAll(elements)
	}
}

func BenchmarkNext10Elements(b *testing.B) {
	elements := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen := New(elements)
		for {
			_, ok := gen.Next()
			if !ok {
				break
			}
		}
	}
}
