package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadersMiddleware(t *testing.T) {
	handler := Headers(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
	}
	for header, val := range want {
		if got := w.Header().Get(header); got != val {
			t.Errorf("header %s: got %q, want %q", header, got, val)
		}
	}
}

func TestQuoteIdentifier(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"users", `"users"`, false},
		{"my_table", `"my_table"`, false},
		{"'; DROP TABLE users; --", "", true},
		{"1invalid", "", true},
		{"valid_col", `"valid_col"`, false},
	}
	for _, tt := range tests {
		got, err := QuoteIdentifier(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("QuoteIdentifier(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("QuoteIdentifier(%q): got %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		input string
		check string // must NOT appear in output
	}{
		{"host=localhost password=s3cr3t dbname=mydb", "s3cr3t"},
		{"postgres://user:p@ssw0rd@localhost/db", "p@ssw0rd"},
	}
	for _, tt := range tests {
		out := RedactDSN(tt.input)
		if contains(out, tt.check) {
			t.Errorf("RedactDSN did not redact %q in %q", tt.check, out)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
