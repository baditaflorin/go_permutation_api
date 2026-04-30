package websocket

import (
	"encoding/json"
	"log/slog"
	"time"

	"golang.org/x/net/websocket"

	"github.com/baditaflorin/go_permutation_api/internal/config"
	"github.com/baditaflorin/go_permutation_api/internal/permutation"
	"github.com/baditaflorin/go_permutation_api/pkg/validator"
)

// Handler handles WebSocket connections for real-time permutation streaming.
type Handler struct {
	cfg *config.Config
}

// New creates a new WebSocket handler.
func New(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// ServeWS is the websocket.Handler entry point.
func (h *Handler) ServeWS(conn *websocket.Conn) {
	defer conn.Close()

	var msg ClientMessage
	if err := websocket.JSON.Receive(conn, &msg); err != nil {
		send(conn, ServerMessage{Type: TypeError, Error: "invalid message: " + err.Error()})
		return
	}

	if msg.Action != ActionStart {
		send(conn, ServerMessage{Type: TypeError, Error: "first message must be {\"action\":\"start\"}"})
		return
	}

	elements := validator.SanitizeElements(msg.Elements)
	if err := validator.ValidateElements(elements, h.cfg.App.MaxElements); err != nil {
		send(conn, ServerMessage{Type: TypeError, Error: err.Error()})
		return
	}

	chunkSize := msg.ChunkSize
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	// Channel to receive stop signal
	stop := make(chan struct{}, 1)
	go func() {
		var stopMsg ClientMessage
		websocket.JSON.Receive(conn, &stopMsg)
		if stopMsg.Action == ActionStop {
			stop <- struct{}{}
		}
	}()

	start := time.Now()
	gen := permutation.New(elements)
	chunk := make([][]string, 0, chunkSize)
	total := 0
	seq := 0

	for {
		select {
		case <-stop:
			send(conn, ServerMessage{
				Type:      TypeDone,
				Total:     total,
				ElapsedMS: time.Since(start).Milliseconds(),
			})
			return
		default:
		}

		perm, ok := gen.Next()
		if !ok {
			break
		}

		cp := make([]string, len(perm))
		copy(cp, perm)
		chunk = append(chunk, cp)
		total++

		if len(chunk) >= chunkSize {
			seq++
			if err := send(conn, ServerMessage{Type: TypeChunk, Data: chunk, Sequence: seq}); err != nil {
				slog.Error("websocket send error", "err", err)
				return
			}
			// progress every 10 chunks
			if seq%10 == 0 {
				send(conn, ServerMessage{
					Type:      TypeProgress,
					Count:     total,
					ElapsedMS: time.Since(start).Milliseconds(),
				})
			}
			chunk = chunk[:0]
		}
	}

	// Flush remaining
	if len(chunk) > 0 {
		seq++
		send(conn, ServerMessage{Type: TypeChunk, Data: chunk, Sequence: seq})
	}

	send(conn, ServerMessage{
		Type:      TypeDone,
		Total:     total,
		ElapsedMS: time.Since(start).Milliseconds(),
	})
}

func send(conn *websocket.Conn, msg ServerMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return websocket.Message.Send(conn, string(data))
}
