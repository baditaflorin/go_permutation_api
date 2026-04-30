package observability

import (
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// counters exposed for Prometheus scraping (ADR-0006).
var (
	reqTotal   uint64
	errTotal   uint64
	permTotal  uint64
	startedAt  = time.Now()
)

// IncrementRequests is called by metrics middleware.
func IncrementRequests() { atomic.AddUint64(&reqTotal, 1) }

// IncrementErrors is called on 4xx/5xx.
func IncrementErrors() { atomic.AddUint64(&errTotal, 1) }

// AddPermutations records generated permutations.
func AddPermutations(n uint64) { atomic.AddUint64(&permTotal, n) }

// PrometheusHandler writes metrics in the Prometheus text exposition format.
// Enabled when ENABLE_PROMETHEUS=true (ADR-0006).
func PrometheusHandler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP permutation_requests_total Total HTTP requests handled\n")
	fmt.Fprintf(w, "# TYPE permutation_requests_total counter\n")
	fmt.Fprintf(w, "permutation_requests_total %d\n\n", atomic.LoadUint64(&reqTotal))

	fmt.Fprintf(w, "# HELP permutation_errors_total Total error responses\n")
	fmt.Fprintf(w, "# TYPE permutation_errors_total counter\n")
	fmt.Fprintf(w, "permutation_errors_total %d\n\n", atomic.LoadUint64(&errTotal))

	fmt.Fprintf(w, "# HELP permutation_generated_total Total permutations generated\n")
	fmt.Fprintf(w, "# TYPE permutation_generated_total counter\n")
	fmt.Fprintf(w, "permutation_generated_total %d\n\n", atomic.LoadUint64(&permTotal))

	fmt.Fprintf(w, "# HELP permutation_uptime_seconds Seconds since server start\n")
	fmt.Fprintf(w, "# TYPE permutation_uptime_seconds gauge\n")
	fmt.Fprintf(w, "permutation_uptime_seconds %.2f\n\n", time.Since(startedAt).Seconds())

	fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Bytes allocated and in use\n")
	fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n\n", mem.Alloc)

	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
}
