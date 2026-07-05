---
title: Observability
sidebar_label: Observability
description: Structured logs, Prometheus metrics, health endpoints — and where tracing fits.
---

# Observability

## Learning objectives

- Assemble the three pillars — logs, metrics, traces — and know what question each answers.
- Expose Prometheus metrics from a Go service and know the four metric types.
- Implement liveness and readiness endpoints with dependency checkers.
- Know the platform's current state: logs + metrics + health standardized; distributed tracing a recorded gap.

## Prerequisites

- [Logging](../module-2-intermediate/logging), [Distributed Systems](distributed-systems)

## Time estimate

**3.5 hours**

## Concepts

### Three pillars, three questions

| Pillar | Question | DX tooling |
|---|---|---|
| **Logs** | *What happened* in this request? | zap, request-ID-correlated (Module 2) |
| **Metrics** | *How is the service doing* in aggregate? | Prometheus `/metrics` |
| **Traces** | *Where did the time go* across services? | concept here; not yet wired (see below) |

Logs are events, metrics are numbers over time, traces are causality across processes. You need all three because each is terrible at the others' job — you can't alert on grep, and you can't debug one request from an average.

### Prometheus metrics

The model: your service exposes current values at `/metrics` in text form; Prometheus scrapes periodically; dashboards and alerts query the result. Four types, in order of how often you'll use them:

- **Counter** — only goes up: `http_requests_total`. Rates come from `rate()` at query time.
- **Histogram** — observations in buckets: `http_request_duration_seconds`. This is how you get p95/p99 latency — and percentiles, not averages, are what you watch (a great average hides a terrible tail).
- **Gauge** — goes up and down: in-flight requests, queue depth, pool size.
- **Summary** — client-side quantiles; rarely what you want over histograms.

```go
var requestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "route", "status"},
)

// in middleware:
requestDuration.WithLabelValues(m, route, status).Observe(elapsed.Seconds())
```

**Label discipline** is the failure mode to learn early: every distinct label combination is a separate time series. `route` (the pattern, `/iudx/v2/notes/{id}`) is a fine label; raw `path` (every actual ID) is a cardinality explosion that takes down your monitoring. Never label by user ID, request ID, or anything unbounded.

The RED trio per endpoint — **R**ate, **E**rrors, **D**uration — is the standard dashboard starting point, and the request-metrics middleware gives you all three from one histogram.

### Health endpoints

Two endpoints with different consumers and different rules ([Distributed Systems](distributed-systems) explained why they must differ):

```go
hh := health.NewHandler()
hh.Register("postgres", health.NewPgxPoolChecker(pool))

r.Get("/healthz/live", hh.Live)   // always 200 while the process runs
r.Get("/healthz/ready", hh.Ready) // 200 iff registered checkers pass
```

- **Liveness** — "don't restart me needlessly." Checks nothing external; if the process answers, it's alive. Putting a DB check here means a DB outage gets your pods restart-looped for no benefit.
- **Readiness** — "route traffic to me?" Checks hard dependencies. Postgres down → not ready → the load balancer stops sending requests → clients hit healthy replicas instead of timeouts.

Optional dependencies (RabbitMQ, per the boot contract) generally *don't* fail readiness — the service degrades but still serves.

### Tracing — the concept, and the honest status

A distributed trace follows one request across services: gateway → authz check → upstream → database, each step a timed **span** sharing a **trace ID** propagated via headers (W3C `traceparent`). It answers "this request took 2s — *where*?" — which logs can only answer within one service. The modern standard is **OpenTelemetry**.

**Platform status, honestly:** DX services standardize logs, metrics, and health today; **OTel tracing is a recorded gap** on the shared-library roadmap (no `dx-common-go` tracing wiring yet). What exists instead is the poor-man's version: the request ID, logged at every hop. When tracing lands, it will arrive as common middleware — another payoff of every service sharing one stack. Learn the concepts now; don't be surprised by the absence in the code.

:::info[Platform connection]
`dx-common-go` provides the whole current stack: `metrics` (Prometheus handler + `RequestMetrics` middleware), `health` (handler + `PgxPoolChecker`), and the zap/request-ID machinery from Module 2. GO-SERVICE-STANDARDS requires all four endpoints on every service: `/healthz/live`, `/healthz/ready`, `/metrics`, `/docs`. Try it now: `curl localhost:8000/metrics` against your local stack and read what the gateway exports.
:::

## Exercises

1. Add `/metrics` to `dx-scratch-go` with the RED histogram above (route-pattern labels!). Generate load; eyeball p95 via the raw bucket counts.
2. Add liveness + readiness with a real Postgres checker. Stop the database container; verify readiness flips to 503 while liveness stays 200 — then explain what Kubernetes would do with each signal.
3. Break monitoring on purpose (in scratch code): label the histogram by raw path, hit 1,000 distinct URLs, and look at the size of `/metrics`. Cardinality, experienced once, is remembered forever.
4. Read your request ID across hops: curl the local gateway, find the same request's line in gateway logs and upstream logs (`docker compose logs | grep <id>`). Write down what a trace would have added.

## Check yourself

- Which pillar answers: "error rate spiked at 14:02"? "this specific upload failed"? "checkout is slow somewhere between four services"?
- Why percentiles instead of averages for latency?
- Why must liveness *not* check the database while readiness must?
- What's unbounded-label cardinality and how did you avoid it in exercise 1?

## References

- [Prometheus docs](https://prometheus.io/docs/introduction/overview/) · [client_golang](https://pkg.go.dev/github.com/prometheus/client_golang)
- [Google SRE Book — Monitoring Distributed Systems](https://sre.google/sre-book/monitoring-distributed-systems/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) — for when the gap closes
- Platform: `dx-common-go/{metrics,health}`; GO-SERVICE-STANDARDS.md (observability)
