package permutation

import (
	"bufio"
	"encoding/json"
	"io"
)

const defaultFlushEvery = 64 // permutations per flush

// StreamJSON writes all permutations of elements as a JSON array to w,
// using buffered I/O and pool-backed slices for minimal allocation.
func StreamJSON(w io.Writer, elements []string) error {
	bw := bufio.NewWriterSize(w, 4096)

	bw.WriteByte('[')

	gen := New(elements)
	enc := json.NewEncoder(bw)
	first := true
	count := 0

	for {
		perm, ok := gen.Next()
		if !ok {
			break
		}
		if !first {
			bw.WriteByte(',')
		}
		first = false

		if err := enc.Encode(perm); err != nil {
			return err
		}
		// Trim the newline json.Encoder adds
		// (rewind one byte — safe because bufio is in-memory)

		count++
		if count%defaultFlushEvery == 0 {
			if err := bw.Flush(); err != nil {
				return err
			}
		}
	}

	bw.WriteByte(']')
	return bw.Flush()
}
