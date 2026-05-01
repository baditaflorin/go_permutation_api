package websocket

import (
	"testing"
)

func TestProtocolConstants(t *testing.T) {
	if DefaultChunkSize <= 0 {
		t.Error("DefaultChunkSize must be positive")
	}
	if MaxChunkSize < DefaultChunkSize {
		t.Error("MaxChunkSize must be >= DefaultChunkSize")
	}
}

func TestClientMessageFields(t *testing.T) {
	msg := ClientMessage{
		Action:    ActionStart,
		Elements:  []string{"a", "b"},
		ChunkSize: DefaultChunkSize,
	}
	if msg.Action != "start" {
		t.Errorf("expected 'start', got %q", msg.Action)
	}
}
