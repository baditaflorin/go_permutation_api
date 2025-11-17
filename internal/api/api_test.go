package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baditaflorin/go_permutation_api/internal/config"
)

func TestHandleGet(t *testing.T) {
	cfg := config.Default()
	handler := NewHandler(cfg)

	req := httptest.NewRequest("GET", "/permutations?elements=a,b,c", nil)
	w := httptest.NewRecorder()

	handler.HandlePermutations(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result [][]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 6 {
		t.Errorf("expected 6 permutations, got %d", len(result))
	}
}

func TestHandlePost(t *testing.T) {
	cfg := config.Default()
	handler := NewHandler(cfg)

	body := []string{"a", "b", "c"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/permutations", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlePermutations(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result [][]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 6 {
		t.Errorf("expected 6 permutations, got %d", len(result))
	}
}

func TestHandleGetTooManyElements(t *testing.T) {
	cfg := config.Default()
	cfg.App.MaxElements = 3
	handler := NewHandler(cfg)

	req := httptest.NewRequest("GET", "/permutations?elements=a,b,c,d,e", nil)
	w := httptest.NewRecorder()

	handler.HandlePermutations(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleHealth(t *testing.T) {
	cfg := config.Default()
	handler := NewHandler(cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HandleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %s", result["status"])
	}
}

func TestHandleInvalidMethod(t *testing.T) {
	cfg := config.Default()
	handler := NewHandler(cfg)

	req := httptest.NewRequest("DELETE", "/permutations", nil)
	w := httptest.NewRecorder()

	handler.HandlePermutations(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
