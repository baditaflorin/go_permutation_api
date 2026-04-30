package permutation

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStreamJSON(t *testing.T) {
	var buf bytes.Buffer
	elements := []string{"a", "b", "c"}

	if err := StreamJSON(&buf, elements); err != nil {
		t.Fatalf("StreamJSON error: %v", err)
	}

	var got [][]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error: %v\nbody: %s", err, buf.String())
	}

	if len(got) != 6 {
		t.Errorf("expected 6 permutations, got %d", len(got))
	}
}

func BenchmarkStreamJSON5(b *testing.B) {
	elements := []string{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		StreamJSON(&buf, elements)
	}
}

func BenchmarkStreamJSON8(b *testing.B) {
	elements := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		StreamJSON(&buf, elements)
	}
}
