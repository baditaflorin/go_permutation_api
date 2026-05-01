package websocket

// ClientMessage is sent from browser → server.
type ClientMessage struct {
	Action    string   `json:"action"`               // "start" | "stop"
	Elements  []string `json:"elements,omitempty"`
	ChunkSize int      `json:"chunk_size,omitempty"` // permutations per chunk, default 100
}

// ServerMessage is sent from server → browser.
type ServerMessage struct {
	Type      string     `json:"type"`                // "chunk" | "progress" | "done" | "error"
	Data      [][]string `json:"data,omitempty"`
	Sequence  int        `json:"sequence,omitempty"`
	Count     int        `json:"count,omitempty"`
	Total     int        `json:"total,omitempty"`
	ElapsedMS int64      `json:"elapsed_ms,omitempty"`
	Error     string     `json:"error,omitempty"`
}

const (
	ActionStart = "start"
	ActionStop  = "stop"

	TypeChunk    = "chunk"
	TypeProgress = "progress"
	TypeDone     = "done"
	TypeError    = "error"

	DefaultChunkSize = 100
	MaxChunkSize     = 1000
)
