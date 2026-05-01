# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 2.x     | ✅ Yes     |
| 1.x     | ❌ No      |

## Reporting a Vulnerability

**Please do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

Send an email to **security@example.com** (replace with your actual address) with:

1. **Subject:** `[SECURITY] go_permutation_api - <brief description>`
2. **Description:** What the vulnerability is and how it can be exploited
3. **Steps to reproduce:** Minimal reproduction case
4. **Impact:** What an attacker could achieve
5. **Suggested fix** (optional): Your recommended remediation

### What to Expect

| Milestone | Timeline |
|-----------|----------|
| Acknowledgement | Within 48 hours |
| Initial assessment | Within 5 business days |
| Fix or mitigation | Within 30 days for critical issues |
| Public disclosure | Coordinated after fix is released |

We follow [responsible disclosure](https://en.wikipedia.org/wiki/Coordinated_vulnerability_disclosure) principles.

### Scope

**In scope:**
- SQL injection in database queries
- Command injection via configuration values
- Authentication bypass in GUI endpoints
- Information disclosure (credentials, stack traces)
- Denial of Service via large payloads or computation
- SSRF via database connection strings

**Out of scope:**
- Vulnerabilities in dependencies (report upstream)
- Issues requiring physical access
- Social engineering

## Security Design Principles

This project follows these security principles:

### Input Validation (ADR-0003)
- All SQL identifiers (table/column names) validated against `[a-zA-Z_][a-zA-Z0-9_]*`
- Request bodies limited to 1 MB via `http.MaxBytesReader`
- Input elements validated for length and content before processing

### Security Headers
Every HTTP response includes:
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Strict-Transport-Security: max-age=63072000 (HTTPS only)
```

### Credential Handling
- No credentials stored in code or committed to version control
- Passwords redacted in log output
- Database connection strings never logged in plain text
- All secrets loaded from environment variables only

### Rate Limiting
- Token bucket rate limiter protects all endpoints
- Configurable via `RATE_LIMIT_REQUESTS` and `RATE_LIMIT_WINDOW` env vars

### Dependency Management
- Minimal external dependencies
- `go mod verify` enforced in CI
- Dependabot alerts enabled

## Security Configuration Checklist

Before deploying to production:

- [ ] Set a strong `DB_PASSWORD` (not `postgres`)
- [ ] Set `DB_SSL_MODE=require` (not `disable`)
- [ ] Set `TRUSTED_PROXIES` to your load balancer CIDR
- [ ] Set `ENABLE_PPROF=false` (default) — profiling endpoints disabled
- [ ] Bind `SERVER_HOST` to a private interface, not `0.0.0.0`
- [ ] Configure CORS: set `CORS_ALLOWED_ORIGINS` to your domain
- [ ] Enable TLS termination at your reverse proxy (nginx, Caddy)
- [ ] Review `MAX_ELEMENTS` — high values allow slow DoS via large computation
- [ ] Restrict `/metrics` access at the network/proxy level

## Known Security Limitations

1. **GUI authentication**: The configuration GUI (`--gui`) has no built-in authentication. It should only be run on trusted networks or behind an authenticated reverse proxy.

2. **No API authentication**: The permutation API has no built-in API key or OAuth mechanism. For public deployment, add authentication at your API gateway.

3. **WebSocket rate limiting**: WebSocket connections share the HTTP rate limiter but long-lived connections persist beyond the window.

## Hall of Fame

We appreciate responsible disclosure. Researchers who report valid issues will be credited here (with their permission).

*Be the first to be listed here by finding and responsibly disclosing a vulnerability!*

---

*This policy is based on [GitHub's security advisory guidelines](https://docs.github.com/en/code-security/security-advisories).*
