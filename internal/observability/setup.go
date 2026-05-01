// Package observability configures structured logging and metrics (ADR-0006).
// It replaces the custom pkg/logger with stdlib log/slog for zero extra deps.
package observability

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// Setup initialises structured logging based on environment.
// LOG_LEVEL: debug | info | warn | error  (default: info)
// LOG_FORMAT: json | text                 (default: json in prod, text locally)
func Setup() {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
	slog.Info("logging initialised", "level", level.String(), "format", effectiveFormat(format))
}

// RequestLogger returns an slog-based HTTP logging middleware.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"remote", r.RemoteAddr,
			"request_id", r.Header.Get("X-Request-ID"),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func effectiveFormat(f string) string {
	if f == "text" {
		return "text"
	}
	return "json"
}
