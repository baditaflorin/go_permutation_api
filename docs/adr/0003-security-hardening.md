# ADR-0003: Security Hardening for Public Release

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** Black Hat (Caution & Risk)

## Context

The Black Hat demands we catalogue every risk before opening this to the internet.

**Identified risks:**

| Risk | Severity | Current State |
|------|----------|---------------|
| No request size limit | HIGH | An attacker POSTs a 100 MB JSON body |
| Panic in handler exposes stack trace | HIGH | `recovery.go` prints stack to stdout — visible in logs |
| `debug/pprof` endpoints exposed by default | HIGH | Any internet user can profile the server |
| docker-compose default password `postgres` | HIGH | Credentials in version control |
| `DB_PASSWORD` logged if startup fails | MEDIUM | `log.Fatalf` may print connection string |
| No `Strict-Transport-Security` header | MEDIUM | HTTPS deployments lack HSTS |
| No `Content-Security-Policy` header on GUI | MEDIUM | GUI HTML has no CSP |
| SQL table/column injected without quoting | MEDIUM | `database.go` uses `fmt.Sprintf` in query |
| `X-Forwarded-For` trusted without validation | LOW | Rate limiter can be bypassed by header spoofing |
| `ENABLE_CORS` defaults to `*` | LOW | Overly permissive for production |

## Decision

Implement a `internal/security` package that applies:

1. **Request body size limit** — `http.MaxBytesReader(w, r.Body, 1<<20)` (1 MB)
2. **Safe pprof** — guarded by `ENABLE_PPROF=true` env var, default false
3. **Security headers middleware** — HSTS, CSP, X-Frame-Options, X-Content-Type-Options
4. **SQL identifier quoting** — use `pq.QuoteIdentifier` for table/column names
5. **Secret redaction in logs** — connection strings stripped of passwords
6. **Trusted proxy validation** — `TRUSTED_PROXIES` env var controls IP trust
7. **docker-compose** — move secrets to `.env` file, never committed
8. **SECURITY.md** — responsible disclosure policy

## Consequences

**Positive:**
- Eliminates highest-severity attack vectors before public launch
- Demonstrates security-consciousness to enterprise adopters
- `SECURITY.md` establishes responsible disclosure process

**Negative:**
- 1 MB body limit may affect legitimate large POST bodies (mitigatable via config)
- Additional headers add ~200 bytes per response

**Implementation:** `internal/security/middleware.go`, `internal/security/sql.go`, `SECURITY.md`
