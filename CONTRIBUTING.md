# Contributing to Go Permutation API

Thank you for your interest in contributing! This guide covers everything you need to get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Architecture Decisions](#architecture-decisions)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

Be respectful. We follow the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## How to Contribute

### Bug Reports
1. Check [existing issues](https://github.com/baditaflorin/go_permutation_api/issues) first
2. Include: Go version, OS, steps to reproduce, expected vs actual behaviour
3. For security issues, see [SECURITY.md](SECURITY.md) — **do not open public issues**

### Feature Requests
1. Open an issue with the `enhancement` label
2. Describe the use case and expected behaviour
3. Reference relevant ADRs if applicable (see `docs/adr/`)

### Code Contributions
1. Fork the repository
2. Create a branch: `git checkout -b feat/my-feature` or `fix/issue-123`
3. Make your changes (see standards below)
4. Add/update tests
5. Open a pull request against `main`

## Development Setup

### Prerequisites

```bash
# Go 1.23+
go version

# Optional: golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Optional: Docker
docker --version
```

### Quick Start

```bash
git clone https://github.com/baditaflorin/go_permutation_api.git
cd go_permutation_api

# Copy environment template
cp .env.example .env

# Download dependencies
make deps

# Run tests
make test

# Build
make build

# Run locally
make run-api   # API server on :8080
make run-gui   # GUI on :3000
```

### Running with Docker

```bash
make docker-up    # starts API + GUI + PostgreSQL
make docker-logs  # tail logs
make docker-down  # stop
```

## Architecture Decisions

Before adding major features, check `docs/adr/` for accepted decisions:

| ADR | Topic |
|-----|-------|
| [ADR-0001](docs/adr/0001-standardize-api-response-format.md) | Response envelope format |
| [ADR-0002](docs/adr/0002-user-experience-error-messages.md) | Error message standards |
| [ADR-0003](docs/adr/0003-security-hardening.md) | Security requirements |
| [ADR-0004](docs/adr/0004-performance-optimization.md) | Performance: pooling & streaming |
| [ADR-0005](docs/adr/0005-websocket-streaming.md) | WebSocket real-time streaming |
| [ADR-0006](docs/adr/0006-observability-opentelemetry.md) | Observability with slog |

To propose a new ADR, open an issue with the template:

```markdown
# ADR-NNNN: <Title>
**Thinking Hat:** <White|Red|Black|Yellow|Green|Blue>
## Context
## Decision
## Consequences
```

## Coding Standards

### Package Structure

```
internal/     Private application code (not importable externally)
  api/        HTTP handlers, middleware, server
  config/     Configuration loading and validation
  datasource/ CSV, TSV, database readers
  gui/        Web configuration UI
  observability/ Logging, metrics, tracing
  permutation/ Core algorithm
  response/   Envelope, error types, request IDs
  security/   Headers, body limits, SQL safety
  websocket/  WebSocket streaming handler
pkg/          Reusable public packages
  validator/  Input validation helpers
```

### Code Style

- **No comments on obvious code** — names should be self-explanatory
- **Comment the WHY not the WHAT** — only non-obvious invariants/workarounds
- **No magic numbers** — use named constants
- **Errors wrap with context**: `fmt.Errorf("loading elements: %w", err)`
- **Table-driven tests** for all non-trivial logic
- **`go fmt`** before committing — enforced by CI

### Security Requirements (ADR-0003)

Every PR that touches HTTP handlers must:
- Use `security.BodyLimit` middleware (no unguarded `r.Body` reads)
- Return errors via `response.Error()` (never `http.Error` with plain text)
- Use `security.QuoteIdentifier()` for any SQL identifiers
- Never log credentials or connection strings without `security.RedactDSN()`

### Error Handling

Use the structured error types from `internal/response/errors.go`:

```go
// ✅ Good
response.Error(w, r, http.StatusBadRequest, response.ErrTooManyElements(len(elements), max))

// ❌ Bad
http.Error(w, "too many elements", http.StatusBadRequest)
```

## Testing Requirements

All new code must include tests. PRs failing tests will not be merged.

### Test Checklist

- [ ] Unit tests for all exported functions
- [ ] Table-driven tests for input validation
- [ ] HTTP handler tests using `httptest.NewRecorder()`
- [ ] Benchmarks for performance-critical paths (permutation, streaming)
- [ ] Security tests for any new middleware

### Running Tests

```bash
make test           # all tests
make coverage       # with HTML coverage report
make bench          # benchmarks only

# Single package
go test ./internal/permutation/... -v -bench=.
```

### Coverage Expectations

| Package | Minimum Coverage |
|---------|-----------------|
| `internal/permutation` | 90% |
| `internal/api` | 80% |
| `internal/security` | 90% |
| `pkg/validator` | 95% |

## Pull Request Process

1. **Title format**: `type(scope): description`
   - Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
   - Example: `feat(api): add pagination to permutation endpoint`

2. **PR description** must include:
   - What changed and why
   - How to test
   - Reference to issue (`Fixes #123`)
   - ADR reference if applicable

3. **Review checklist** (filled by reviewer):
   - [ ] Tests pass in CI
   - [ ] Code follows style guidelines
   - [ ] Security requirements met (ADR-0003)
   - [ ] Documentation updated
   - [ ] CHANGELOG.md entry added

4. **Merge policy**: Squash merge to `main`, PR author writes final commit message.

## Release Process

Releases are tagged by maintainers:

```bash
git tag -a v2.1.0 -m "Release v2.1.0"
git push origin v2.1.0
```

The CI pipeline builds release binaries and Docker images automatically.

---

Thank you for contributing! 🎉
