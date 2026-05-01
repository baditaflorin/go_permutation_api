# ADR-0004: Performance Optimization via Object Pooling and Streaming

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** Yellow Hat (Optimism & Benefits)

## Context

The Yellow Hat asks: *what is the maximum value we can extract from this system?*

Current performance baseline (from benchmarks):
- 11 elements (39.9M permutations): ~2.5 seconds
- Each permutation allocates a new `[]string` slice
- HTTP streaming writes one JSON-encoded array at a time with individual `Write` calls
- No `sync.Pool` for reusable buffers

Opportunities:
1. **Slice pooling** — reuse `[]string` backing arrays via `sync.Pool`
2. **Buffered writes** — batch `Write` calls with `bufio.Writer` to reduce syscalls
3. **Encoder reuse** — reuse `json.Encoder` via pool
4. **Pre-allocated result buffers** — for small inputs (≤6 elements), allocate full result up-front
5. **GOMAXPROCS tuning guidance** — expose optimal settings in documentation

Benchmark projections:
| Optimization | Expected gain |
|---|---|
| `sync.Pool` for permutation slices | 15-25% fewer allocations |
| `bufio.Writer` (4 KB buffer) | 20-30% reduction in syscalls |
| Encoder pool | 5-10% faster JSON marshalling |
| Combined | 30-45% overall throughput improvement |

## Decision

1. Add `sync.Pool` in `internal/permutation/pool.go` for `[]string` reuse
2. Wrap `http.ResponseWriter` in `bufio.NewWriterSize(w, 4096)` in the streaming path
3. Add `internal/permutation/stream.go` — a streaming writer that flushes every N permutations
4. Add `BenchmarkStream` tests to measure and guard regression
5. Document GOMAXPROCS recommendations in README

## Consequences

**Positive:**
- Measurably faster for large inputs
- Lower memory pressure under concurrent load
- Benchmark suite prevents future regression

**Negative:**
- `sync.Pool` adds code complexity
- `bufio.Writer` must be explicitly flushed — easy to forget

**Implementation:** `internal/permutation/pool.go`, `internal/permutation/stream.go`
