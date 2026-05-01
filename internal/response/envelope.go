package response

import (
	"encoding/json"
	"net/http"

	"github.com/baditaflorin/go_permutation_api/internal/response/reqid"
)

// Envelope wraps all API responses in a consistent structure.
type Envelope struct {
	Version   string      `json:"version"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
	Meta      *Meta       `json:"meta,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
}

// Meta holds pagination and result metadata.
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// JSON writes a successful envelope response.
func JSON(w http.ResponseWriter, r *http.Request, status int, data interface{}, meta *Meta) {
	env := Envelope{
		Version:   "1",
		RequestID: reqid.FromContext(r.Context()),
		Data:      data,
		Meta:      meta,
	}
	write(w, status, env)
}

// Error writes an error envelope response.
func Error(w http.ResponseWriter, r *http.Request, status int, err *APIError) {
	env := Envelope{
		Version:   "1",
		RequestID: reqid.FromContext(r.Context()),
		Error:     err,
	}
	write(w, status, env)
}

func write(w http.ResponseWriter, status int, env Envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(env)
}
