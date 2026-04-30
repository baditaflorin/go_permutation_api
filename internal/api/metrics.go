package api

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// Metrics holds application metrics
type Metrics struct {
	RequestCount      uint64    `json:"request_count"`
	ErrorCount        uint64    `json:"error_count"`
	TotalPermutations uint64    `json:"total_permutations"`
	Uptime            string    `json:"uptime"`
	MemoryUsage       uint64    `json:"memory_usage_bytes"`
	GoRoutines        int       `json:"goroutines"`
	StartTime         time.Time `json:"-"`
}

var globalMetrics = &Metrics{
	StartTime: time.Now(),
}

// IncrementRequests increments the request counter
func IncrementRequests() {
	atomic.AddUint64(&globalMetrics.RequestCount, 1)
}

// IncrementErrors increments the error counter
func IncrementErrors() {
	atomic.AddUint64(&globalMetrics.ErrorCount, 1)
}

// AddPermutations adds to the permutations counter
func AddPermutations(count uint64) {
	atomic.AddUint64(&globalMetrics.TotalPermutations, count)
}

// GetMetrics returns current metrics
func GetMetrics() *Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Metrics{
		RequestCount:      atomic.LoadUint64(&globalMetrics.RequestCount),
		ErrorCount:        atomic.LoadUint64(&globalMetrics.ErrorCount),
		TotalPermutations: atomic.LoadUint64(&globalMetrics.TotalPermutations),
		Uptime:            time.Since(globalMetrics.StartTime).String(),
		MemoryUsage:       m.Alloc,
		GoRoutines:        runtime.NumGoroutine(),
		StartTime:         globalMetrics.StartTime,
	}
}

// HandleMetrics returns metrics as JSON
func (h *Handler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetMetrics())
}

// metricsMiddleware tracks requests
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		IncrementRequests()
		wrapper := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)
		if wrapper.statusCode >= 400 {
			IncrementErrors()
		}
	})
}

// metricsResponseWriter wraps ResponseWriter to capture status code
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}
