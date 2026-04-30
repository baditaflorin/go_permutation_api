package permutation

import (
	"sort"
)

// Generator handles permutation generation
type Generator struct {
	elements []string
	first    bool
}

// New creates a new permutation generator
func New(elements []string) *Generator {
	// Make a copy to avoid modifying the original slice
	copied := make([]string, len(elements))
	copy(copied, elements)

	// Sort for lexicographical order
	sort.Strings(copied)

	return &Generator{
		elements: copied,
		first:    true,
	}
}

// Next generates the next permutation
// Returns the current permutation and true if more permutations exist
// Returns nil and false when all permutations have been generated
func (g *Generator) Next() ([]string, bool) {
	if g.first {
		g.first = false
		result := make([]string, len(g.elements))
		copy(result, g.elements)
		return result, true
	}

	if !g.nextPermutation() {
		return nil, false
	}

	result := make([]string, len(g.elements))
	copy(result, g.elements)
	return result, true
}

// nextPermutation generates the next lexicographical permutation
// Returns true if successful, false if no more permutations exist
func (g *Generator) nextPermutation() bool {
	arr := g.elements
	n := len(arr)

	if n <= 1 {
		return false
	}

	// Find the rightmost element that is smaller than its right neighbor
	i := n - 2
	for i >= 0 && arr[i] >= arr[i+1] {
		i--
	}

	// If no such element exists, we're at the last permutation
	if i < 0 {
		return false
	}

	// Find the rightmost element that is larger than arr[i]
	j := n - 1
	for arr[j] <= arr[i] {
		j--
	}

	// Swap them
	arr[i], arr[j] = arr[j], arr[i]

	// Reverse the suffix starting at i+1
	reverse(arr, i+1)

	return true
}

// reverse reverses the array from start index to the end
func reverse(arr []string, start int) {
	end := len(arr) - 1
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

// GenerateAll generates all permutations at once
// Warning: This can be memory-intensive for large input sets
func GenerateAll(elements []string) [][]string {
	gen := New(elements)
	var results [][]string

	for {
		perm, ok := gen.Next()
		if !ok {
			break
		}
		results = append(results, perm)
	}

	return results
}

// Count calculates the number of permutations without generating them
// Returns n! for n unique elements
func Count(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 || n == 1 {
		return 1
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}
