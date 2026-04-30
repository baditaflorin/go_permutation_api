package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baditaflorin/go_permutation_api/internal/api"
	"github.com/baditaflorin/go_permutation_api/internal/config"
)


func TestIntegrationAPIServer(t *testing.T) {
	// Create test configuration
	cfg := config.Default()
	cfg.App.MaxElements = 5

	// Create server
	server := api.NewServer(cfg)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		checkResponse  func(*testing.T, []byte)
	}{
		{
			name:           "health check",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if resp["status"] != "healthy" {
					t.Errorf("expected status 'healthy', got %s", resp["status"])
				}
			},
		},
		{
			name:           "version endpoint",
			method:         "GET",
			path:           "/version",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if resp["version"] == "" {
					t.Error("version should not be empty")
				}
			},
		},
		{
			name:           "GET permutations",
			method:         "GET",
			path:           "/?elements=a,b,c",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var perms [][]string
				if err := json.Unmarshal(body, &perms); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(perms) != 6 {
					t.Errorf("expected 6 permutations, got %d", len(perms))
				}
			},
		},
		{
			name:           "POST permutations",
			method:         "POST",
			path:           "/",
			body:           []string{"x", "y"},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var perms [][]string
				if err := json.Unmarshal(body, &perms); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(perms) != 2 {
					t.Errorf("expected 2 permutations, got %d", len(perms))
				}
			},
		},
		{
			name:           "too many elements",
			method:         "GET",
			path:           "/?elements=1,2,3,4,5,6,7,8,9,10",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("too many elements")) {
					t.Error("expected 'too many elements' error")
				}
			},
		},
		{
			name:           "missing elements parameter",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("required")) {
					t.Error("expected 'required' error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}
