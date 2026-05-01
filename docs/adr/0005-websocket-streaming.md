# ADR-0005: Real-Time WebSocket Streaming

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** Green Hat (Creativity & Innovation)

## Context

The Green Hat asks: *what has never been done before that would make this remarkable?*

The current HTTP streaming model has a fundamental limitation: the client cannot cancel in-flight generation, request a subset, or receive progress events. For 12-element inputs (479M permutations), there is no way to say "stop after 1000 results."

WebSockets enable:
1. **Bidirectional control** — client sends `{"action":"stop"}` mid-stream
2. **Progress events** — server sends `{"type":"progress","count":100000,"total_estimated":479001600}`
3. **Chunked delivery** — client requests next N permutations on demand (back-pressure)
4. **Real-time GUI** — the configuration GUI can preview permutations live
5. **Multi-client broadcast** — one computation, many observers

No standard permutation API offers WebSocket streaming. This is a genuine differentiator.

## Decision

Add a WebSocket endpoint at `GET /ws` using the standard `golang.org/x/net/websocket` package (no CGO, pure Go).

Protocol (JSON messages):

**Client → Server (start):**
```json
{"action": "start", "elements": ["a","b","c"], "chunk_size": 100}
```

**Server → Client (chunk):**
```json
{"type": "chunk", "data": [["a","b","c"],["a","c","b"]], "sequence": 1}
```

**Server → Client (progress):**
```json
{"type": "progress", "count": 200, "elapsed_ms": 15}
```

**Client → Server (stop):**
```json
{"action": "stop"}
```

**Server → Client (done):**
```json
{"type": "done", "total": 6, "elapsed_ms": 1}
```

The WebSocket handler respects `MAX_ELEMENTS` and rate limiting.

## Consequences

**Positive:**
- Unique capability that distinguishes the API
- Enables cancellable long-running computations
- Powers real-time GUI preview
- No polling needed — push-based

**Negative:**
- WebSocket connections are stateful — more complex to scale horizontally
- Requires new dependency (`golang.org/x/net`)
- Client library required for non-browser consumers

**Implementation:** `internal/websocket/handler.go`, `internal/websocket/protocol.go`
