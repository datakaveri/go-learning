---
title: Logging
sidebar_label: Logging
description: Structured logging with zap — levels, fields, request correlation, and what not to log.
---

# Logging

## Learning objectives

- Log structured events with zap: typed fields, not formatted strings.
- Choose levels deliberately and know what each means operationally.
- Correlate every log line of a request with a request ID.
- Apply the platform rules: no `fmt.Println`, log-or-return, no secrets in logs.

## Prerequisites

- [Error Handling](../module-1-go-fundamentals/error-handling), [Dependency Injection](dependency-injection)

## Time estimate

**2.5 hours**

## Concepts

### Structured beats formatted

A log line is data for machines first (aggregators, alerts, dashboards) and humans second. Compare:

```go
fmt.Printf("policy %s created by %s in %dms\n", id, user, ms)   // grep-only
logger.Info("policy created",                                    // queryable
	zap.String("policy_id", id),
	zap.String("user_id", user),
	zap.Int64("duration_ms", ms),
)
```

The second emits JSON with typed fields you can filter on (`policy_id="..."`). The platform logger is **zap** (`go.uber.org/zap`) everywhere — chosen for typed fields and near-zero allocation. `fmt.Println` in service code fails review; so does stdlib `log` outside `main`'s bootstrap.

```go
logger, err := zap.NewProduction() // JSON, ISO timestamps, sampling
defer logger.Sync()                // flush buffers on exit — always defer this in main
```

### Levels, operationally defined

| Level | Meaning | Someone should… |
|---|---|---|
| `Debug` | Developer detail, off in production | nobody — it's for you |
| `Info` | Normal, notable events (boot, policy created, consumer started) | see it on a dashboard |
| `Warn` | Degraded but functioning (optional dep down, retry succeeded) | investigate this week |
| `Error` | Request/operation failed | investigate today |
| `Fatal` | Cannot start/continue — logs then `os.Exit(1)` | get paged. **Only in main**, during boot |

The hard/optional dependency split from the DI page maps directly: Postgres unreachable at boot → `Fatal`; RabbitMQ unreachable → `Warn` and a no-op publisher.

### Request correlation

One request produces many lines across many functions. They're joinable only if they share a **request ID**. The pattern: middleware generates (or receives) the ID, puts it in the context, and a request-scoped logger carries it automatically:

```go
func loggerFrom(ctx context.Context, base *zap.Logger) *zap.Logger {
	if id, ok := requestIDFrom(ctx); ok {
		return base.With(zap.String("request_id", id))
	}
	return base
}
```

`logger.With(...)` returns a child logger with fields pre-attached — build one per request, per consumer message, or per worker cycle, and every subsequent line is correlated for free.

### What not to log

- **Secrets and credentials** — tokens, passwords, DSNs, HMAC signatures. Ever.
- **Full request/response bodies** — PII and noise; log IDs and sizes.
- **The same error at every layer** — the log-or-return rule from [Error Handling](../module-1-go-fundamentals/error-handling). One failure, one `Error` line, at the top.

:::info[Platform connection]
`dx-common-go`'s `Logger()` middleware does the whole correlation dance for you: it captures method, path, status, duration, and the request ID injected by the `RequestID()` middleware one step earlier — which is exactly why the standard middleware order puts RequestID first. Services receive `*zap.Logger` via constructor injection (never a global), and GO-SERVICE-STANDARDS bans `fmt.Println` outright. When you read any DX service's logs in `docker compose logs`, every line is JSON with `request_id` — now you know the machinery.
:::

## Exercises

1. Convert a `fmt.Printf`-riddled program to zap with proper levels and typed fields. Run with `NewDevelopment()` (human-readable) and `NewProduction()` (JSON) and compare.
2. Write an HTTP middleware that generates a request ID, stores it in context, and logs one line per request with method/path/status/duration — you'll compare it to the platform's in Module 3.
3. Add logging to your file-hasher mini-project: `Info` summary, `Warn` for unreadable files, one child logger per worker (`zap.Int("worker", i)`).
4. Grep exercise: given 200 lines of JSON logs (generate them), extract all lines for one request with `jq 'select(.request_id=="...")'`.

## Check yourself

- Why typed fields instead of `fmt.Sprintf` into the message?
- Which level for: a failed optional dependency at boot? A 404 on a GET? A panic recovered by middleware?
- How does a request ID get onto every log line without threading a logger parameter everywhere?
- Why is `Fatal` restricted to `main`?

## References

- [zap docs](https://pkg.go.dev/go.uber.org/zap)
- [Google SRE — implementing service level objectives](https://sre.google/workbook/implementing-slos/) (why logs are data)
- Platform: `dx-common-go` middleware (`RequestID`, `Logger`); GO-SERVICE-STANDARDS.md (observability rules)
