# ADR-0002: Human-Centred Error Messages and API UX

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** Red Hat (Emotions & Feelings)

## Context

The Red Hat asks: *how does a developer feel when they hit an error?* Currently they receive:

```
too many elements: maximum is 12, got 15
```

That is technically accurate but cold. It does not tell the developer:
- What they *should* do next
- Where the limit comes from (config? hard code?)
- Whether they can request a higher limit
- What a valid request looks like

Developers frustration-quit APIs when errors feel like dead ends. Good developer experience (DX) is a retention strategy.

## Decision

Every error response will include:
1. A human-readable `message` (warm, not robotic)
2. A `suggestion` field with an actionable next step
3. A `docs_url` linking to the relevant API docs section
4. A stable machine-readable `code` string

Example:

```json
{
  "error": {
    "code": "TOO_MANY_ELEMENTS",
    "message": "You provided 15 elements, but the maximum allowed is 12.",
    "suggestion": "Reduce your input to 12 or fewer elements, or increase MAX_ELEMENTS in your server configuration.",
    "docs_url": "https://github.com/baditaflorin/go_permutation_api/blob/main/docs/API.md#error-handling"
  }
}
```

Error codes catalogue (stable, versioned):
| Code | HTTP | Meaning |
|------|------|---------|
| `TOO_MANY_ELEMENTS` | 400 | Input exceeds `MAX_ELEMENTS` |
| `EMPTY_ELEMENT` | 400 | An element is blank or whitespace |
| `MISSING_PARAMETER` | 400 | Required query param absent |
| `INVALID_JSON` | 400 | POST body is not valid JSON |
| `METHOD_NOT_ALLOWED` | 405 | Unsupported HTTP method |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `REQUEST_TIMEOUT` | 504 | Processing took too long |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

## Consequences

**Positive:**
- Dramatically better DX; developers spend less time guessing
- Stable error codes let clients handle errors programmatically
- `docs_url` drives traffic to documentation

**Negative:**
- More work per error path
- Must maintain error code catalogue

**Implementation:** `internal/response/errors.go`
