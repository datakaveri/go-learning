---
title: Context Propagation
sidebar_label: Context
description: context.Context — cancellation, deadlines, request-scoped values, and the conventions every DX function follows.
---

# Context Propagation

## Learning objectives

- Explain what `context.Context` carries: cancellation, deadline, and request-scoped values.
- Follow the conventions: `ctx` is the first parameter, flows down every call, and is never stored in a struct.
- Create derived contexts with `WithCancel`, `WithTimeout`, and use `ctx.Done()` in select loops.
- Use context values sparingly and with typed keys.

## Prerequisites

- [Concurrency](concurrency)

## Time estimate

**3 hours**

## Concepts

### One request, one context

Every request that enters a DX service — HTTP call, RabbitMQ delivery, cron tick — gets a `context.Context` that accompanies **every** function call made on its behalf:

```go
func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // created by net/http, cancelled if the client disconnects
	policy, err := h.svc.Create(ctx, req)
	...
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Policy, error) {
	return s.store.Insert(ctx, toDomain(req)) // keep passing it down
}
```

When the client disconnects or a timeout fires, the context is **cancelled**, and everything holding it — database queries, HTTP calls, your loops — can stop promptly instead of wasting work.

### The conventions (memorize these)

1. `ctx context.Context` is the **first parameter**, named `ctx`.
2. **Pass it down; never store it** in a struct. A stored context outlives its request and cancels at the wrong time. (Rare exceptions exist deep in libraries; your code doesn't need them.)
3. Never pass `nil` — use `context.Background()` at the top level (in `main`, tests) and `context.TODO()` as a temporary marker.
4. Cancellation is advisory: functions must **check** it. Long loops include a `ctx.Done()` case.

### Deriving contexts

```go
// Timeout: auto-cancels after 5s — and ALWAYS defer cancel() to release resources
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

row := pool.QueryRow(ctx, query, id) // pgx aborts the query if ctx expires
```

```go
// Cancel: how main stops background workers on shutdown
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
...
cancel() // worker's <-ctx.Done() fires; it returns
```

`ctx.Err()` tells you why: `context.Canceled` or `context.DeadlineExceeded`. Wrap it like any error.

### The shutdown pattern

This is how every DX background loop terminates — the concrete meaning of "context-cancelled on shutdown" from the standards:

```go
func (d *Dispatcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := d.drainOutbox(ctx); err != nil {
				d.log.Warn("outbox drain failed", zap.Error(err))
			}
		case <-ctx.Done():
			return ctx.Err() // graceful exit
		}
	}
}
```

### Context values — rarely, and with typed keys

`context.WithValue` attaches request-scoped data. It's for **cross-cutting metadata that rides along with a request** — request IDs, the authenticated user — not for passing parameters. Use unexported typed keys so packages can't collide:

```go
type ctxKey struct{}

func WithUser(ctx context.Context, u *DxUser) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

func UserFrom(ctx context.Context) (*DxUser, bool) {
	u, ok := ctx.Value(ctxKey{}).(*DxUser) // comma-ok, as always
	return u, ok
}
```

If a function can't work without the value, make it an explicit parameter instead.

:::info[Platform connection]
Both uses are live in `dx-common-go`: the **RequestID middleware** stores the request ID in the context so the zap logger can stamp every line, and the **auth resolver middleware** stores the authenticated `DxUser` so handlers can retrieve the caller with an accessor exactly like `UserFrom` above. Every repository method in every service takes `ctx` first and hands it to pgx — which is what makes statement timeouts and client-disconnect cleanup actually work.
:::

## Exercises

1. Write `slowOp(ctx, d time.Duration)` that respects cancellation via select. Call it with a shorter `WithTimeout` and confirm you get `context.DeadlineExceeded` (test with `errors.Is`).
2. Retrofit your Module-2 file hasher: workers take `ctx`, Ctrl-C (`signal.NotifyContext`) cancels everything, and the program exits cleanly with a partial summary.
3. Implement `WithUser`/`UserFrom` with a typed key; demonstrate that another package using a `string` key cannot collide with yours.
4. Find the bug: a struct that stores `ctx` from its constructor and uses it in a method called minutes later. Explain what goes wrong and fix the API.

## Check yourself

- Why is storing a context in a struct wrong?
- What's the difference between `context.Background()` and `context.TODO()`?
- Why must you `defer cancel()` even for a timeout that "will fire anyway"?
- What belongs in a context value, and what never does?

## References

- [Go Blog: Go Concurrency Patterns: Context](https://go.dev/blog/context) — required reading
- [context package docs](https://pkg.go.dev/context)
- [Google Go Style Guide — Contexts](https://google.github.io/styleguide/go/decisions#contexts)
