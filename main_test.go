package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test the permutations function
func TestPermutations(t *testing.T) {
	tests := []struct {
		input    []string
		expected [][]string
	}{
		{
			input:    []string{"a", "b", "c"},
			expected: [][]string{{"a", "b", "c"}, {"a", "c", "b"}, {"b", "a", "c"}, {"b", "c", "a"}, {"c", "a", "b"}, {"c", "b", "a"}},
		},
		{
			input:    []string{"1"},
			expected: [][]string{{"1"}},
		},
		{
			input:    []string{},
			expected: [][]string{},
		},
	}

	for _, test := range tests {
		result := permutations(test.input)
		if !equal(result, test.expected) {
			t.Errorf("For input %v, expected %v, but got %v", test.input, test.expected, result)
		}
	}
}

// Test the HTTP GET handler
func TestHandleGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/permutations?elements=a,b,c", nil)
	w := httptest.NewRecorder()

	handleGet(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}

	var result [][]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	expected := [][]string{{"a", "b", "c"}, {"a", "c", "b"}, {"b", "a", "c"}, {"b", "c", "a"}, {"c", "a", "b"}, {"c", "b", "a"}}
	if !equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Test the HTTP POST handler
func TestHandlePost(t *testing.T) {
	body := strings.NewReader(`["a", "b", "c"]`)
	req := httptest.NewRequest(http.MethodPost, "/permutations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlePost(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}

	var result [][]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	expected := [][]string{{"a", "b", "c"}, {"a", "c", "b"}, {"b", "a", "c"}, {"b", "c", "a"}, {"c", "a", "b"}, {"c", "b", "a"}}
	if !equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Utility function to compare two 2D slices
func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
