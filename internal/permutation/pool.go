package permutation

import "sync"

// slicePool reuses []string backing arrays to reduce allocations during
// high-throughput streaming (ADR-0004).
var slicePool = sync.Pool{
	New: func() interface{} {
		s := make([]string, 0, 12)
		return &s
	},
}

// getSlice returns a pooled slice with at least the given capacity.
func getSlice(n int) []string {
	sp := slicePool.Get().(*[]string)
	s := (*sp)[:0]
	if cap(s) < n {
		s = make([]string, n)
	} else {
		s = s[:n]
	}
	return s
}

// putSlice returns a slice back to the pool.
func putSlice(s []string) {
	sp := &s
	slicePool.Put(sp)
}
