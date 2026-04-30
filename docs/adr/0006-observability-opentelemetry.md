# ADR-0006: Structured Observability with OpenTelemetry

**Status:** Accepted  
**Date:** 2024-01-15  
**Thinking Hat:** Blue Hat (Process & Control)

## Context

The Blue Hat manages the process — it asks: *how do we know the system is working correctly, now and in six months?*

Current observability state:
- Metrics are a custom in-memory counter (not exportable)
- Logging uses a custom `pkg/logger` (not structured JSON)
- No distributed tracing — impossible to correlate a GUI config save with an API request
- No alerting hooks
- `/metrics` endpoint returns JSON but not Prometheus-compatible scrape format

The industry standard for observability is the OpenTelemetry (OTel) specification, which provides:
- **Traces** — distributed request tracing (Jaeger, Zipkin, Tempo compatible)
- **Metrics** — Prometheus-compatible exposition format
- **Logs** — structured JSON with trace correlation

Without OTel, operators cannot:
- Set up SLOs/SLAs
- Debug slow requests in production
- Correlate frontend GUI events with backend API calls
- Integrate with PagerDuty, Grafana, Datadog, etc.

## Decision

Integrate `go.opentelemetry.io/otel` with three exporters (all opt-in via env vars):

| Signal | Exporter | Env Var |
|--------|----------|---------|
| Traces | OTLP gRPC | `OTEL_EXPORTER_OTLP_ENDPOINT` |
| Metrics | Prometheus | `ENABLE_PROMETHEUS=true` → `/metrics` |
| Logs | slog (JSON) | always-on, replaces custom logger |

Replace `pkg/logger` with Go's stdlib `log/slog` (available since Go 1.21) for zero-dependency structured logging.

Add a `internal/observability` package with:
- `Setup(cfg)` — initialises OTel SDK
- `Shutdown(ctx)` — flushes all exporters on graceful shutdown
- `SpanFromContext` — helper for adding spans to handlers
- Prometheus `/metrics` handler when `ENABLE_PROMETHEUS=true`

## Consequences

**Positive:**
- Drop-in compatible with Grafana, Datadog, New Relic, Honeycomb
- Zero-cost when disabled (OTel SDK no-ops)
- `log/slog` replaces custom logger — fewer lines of code
- Enables SLO dashboards and alerting

**Negative:**
- New dependencies: `go.opentelemetry.io/otel` tree (~10 packages)
- Slightly larger binary (+~3 MB)
- Teams must run a collector (e.g., otel-collector) to use traces

**Implementation:** `internal/observability/setup.go`, `internal/observability/prometheus.go`
