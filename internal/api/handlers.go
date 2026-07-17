package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/baditaflorin/go_permutation_api/internal/config"
	"github.com/baditaflorin/go_permutation_api/internal/permutation"
	"github.com/baditaflorin/go_permutation_api/pkg/validator"
)

// Handler handles HTTP requests
type Handler struct {
	config *config.Config
}

// maxStreamedPermutations bounds a single HTTP response. Even valid input
// under MaxElements can otherwise request 12! JSON arrays and exhaust a
// worker/network while producing evidence no consumer can use.
const maxStreamedPermutations = 100000

// NewHandler creates a new handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
	}
}

// HandlePermutations handles permutation generation requests
func (h *Handler) HandlePermutations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGet handles GET requests with query parameters
func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	elements, err := h.parseElementsFromQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writePermutations(w, elements)
}

// handlePost handles POST requests with JSON body
func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	elements, err := h.parseElementsFromBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateElements(elements, h.config.App.MaxElements); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validatePermutationVolume(elements); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writePermutations(w, elements)
}

// parseElementsFromQuery extracts elements from query parameters
func (h *Handler) parseElementsFromQuery(r *http.Request) ([]string, error) {
	elementsStr := r.URL.Query().Get("elements")
	if elementsStr == "" {
		return nil, fmt.Errorf("'elements' query parameter is required")
	}

	elements := strings.Split(elementsStr, ",")
	elements = validator.SanitizeElements(elements)

	if err := validator.ValidateElements(elements, h.config.App.MaxElements); err != nil {
		return nil, err
	}
	if err := validatePermutationVolume(elements); err != nil {
		return nil, err
	}

	return elements, nil
}

func validatePermutationVolume(elements []string) error {
	if _, exceeds := permutation.CountWithinLimit(len(elements), maxStreamedPermutations); exceeds {
		return fmt.Errorf("too many permutations: maximum streamed response is %d", maxStreamedPermutations)
	}
	return nil
}

// parseElementsFromBody extracts elements from JSON body
func (h *Handler) parseElementsFromBody(r *http.Request) ([]string, error) {
	var elements []string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&elements); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	elements = validator.SanitizeElements(elements)
	return elements, nil
}

// writePermutations streams permutations as JSON to the response
func (h *Handler) writePermutations(w http.ResponseWriter, elements []string) {
	w.Header().Set("Content-Type", "application/json")

	// Write opening bracket
	w.Write([]byte("["))

	gen := permutation.New(elements)
	first := true

	for {
		perm, ok := gen.Next()
		if !ok {
			break
		}

		if !first {
			w.Write([]byte(","))
		}
		first = false

		// Encode and write permutation
		data, err := json.Marshal(perm)
		if err != nil {
			// Can't change status code at this point, log the error
			break
		}
		w.Write(data)
	}

	// Write closing bracket
	w.Write([]byte("]"))
}

// HandleHealth handles health check requests
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
