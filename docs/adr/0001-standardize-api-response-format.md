# ADR-0001: Standardize API Response Format

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** White Hat (Facts & Information)

## Context

The White Hat demands we look at facts. The current API returns raw JSON arrays with no envelope, no error structure, and no metadata. Real-world API consumers need:

- A predictable response schema they can code against without conditional parsing
- Consistent error objects (not bare text strings)
- Request correlation IDs for debugging distributed systems
- Pagination metadata embedded in responses, not just headers
- Content negotiation via `Accept` header

Observed facts from the codebase:
- `handlers.go` writes `[` then streams arrays then `]` — no envelope
- Errors are returned as `http.Error(w, "message", code)` — plain text body
- No `X-Request-ID` header is generated
- `pagination.go` sets headers but the body has no `meta` field

## Decision

Adopt a versioned JSON envelope for **all** API responses:

```json
{
  "version": "1",
  "request_id": "uuid-v4",
  "data": [...],
  "meta": { "page": 1, "per_page": 100, "total": 6, "total_pages": 1 },
  "error": null
}
```

Error responses follow the same envelope:

```json
{
  "version": "1",
  "request_id": "uuid-v4",
  "data": null,
  "meta": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "too many elements: maximum is 12, got 15",
    "details": []
  }
}
```

The `data` field streams permutations. Legacy behaviour is preserved via `Accept: application/x-ndjson` for raw streaming.

## Consequences

**Positive:**
- All clients can parse a single schema
- Errors are machine-readable with stable `code` strings
- Request tracing becomes trivial via `request_id`

**Negative:**
- Breaking change for existing GET/POST callers expecting bare arrays
- Slightly larger payload (mitigated by gzip compression already present)

**Implementation:** `internal/response/envelope.go`
