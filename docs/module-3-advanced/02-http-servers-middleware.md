---
title: HTTP Servers & Middleware
sidebar_label: HTTP & Middleware
description: net/http, the chi router, writing middleware, and the platform's StandardStack — order included.
---

# HTTP Servers & Middleware

## Learning objectives

- Understand `net/http`'s core contracts: `Handler`, `HandlerFunc`, `ServeMux` — and what chi adds.
- Write middleware and reason about wrapping order.
- Reproduce the platform's `StandardStack` order and justify each position.
- Run a server with timeouts and graceful shutdown.

## Prerequisites

- [Project Structure](project-structure), [Context](../module-2-intermediate/context), [Logging](../module-2-intermediate/logging)

## Time estimate

**4 hours**

## Concepts

### net/http in one interface

```go
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

Everything — handlers, routers, middleware, whole applications — is a `Handler`. `http.HandlerFunc` adapts a plain function to the interface. This uniformity is why middleware composes at all.

### chi — the platform's router

The stdlib mux is fine but spartan; DX uses **chi** everywhere: stdlib-compatible (`chi.Router` *is* an `http.Handler`), URL parameters, and route groups with scoped middleware:

```go
r := chi.NewRouter()
r.Use(middleware.RequestID()) // router-wide middleware

r.Route("/iudx/v2/policies", func(r chi.Router) {
	r.Use(authResolver)                  // group-scoped middleware
	r.Get("/", h.ListPolicies)
	r.Post("/", h.CreatePolicy)
	r.Get("/{policyID}", h.GetPolicy)    // chi.URLParam(r, "policyID")
})

r.Get("/healthz/live", hh.Live)          // ops routes: outside auth
```

### Middleware — handlers wrapping handlers

```go
func Logger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r) // call the inner handler
			log.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("took", time.Since(start)),
			)
		})
	}
}
```

`r.Use(A, B, C)` produces `A(B(C(handler)))` — requests pass A → B → C → handler, responses unwind back out. **Order is therefore semantics**, which is why the platform standardizes it:

### The StandardStack — order and reasons

```go
r.Use(
	middleware.RequestID(),   // 1. first, so every later line can be correlated
	chimw.RealIP,             // 2. resolve client IP from proxy headers before logging it
	middleware.Logger(log),   // 3. one line per request, request_id attached
	middleware.CORS(cfg),     // 4. answer preflights before doing work
	middleware.Compression(), // 5. wrap the writer before the handler writes
	chimw.Recoverer,          // 6. convert panics below into 500s (and log them)
	chimw.Timeout(15*time.Second), // 7. bound every handler via r.Context()
)
```

Two positions people get wrong: **RequestID before Logger** (or the log line has no ID), and **Recoverer below Logger** (so a panic is still logged as a completed 500 request). After the stack, services add OpenAPI request validation, then auth and audit middleware on protected route groups — you'll wire those in the next two pages.

### Running the server properly

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           r,
	ReadHeaderTimeout: 5 * time.Second,  // slowloris defense
	ReadTimeout:       10 * time.Second,
	WriteTimeout:      20 * time.Second,
}
```

Never `http.ListenAndServe(":8080", r)` bare — it has **no timeouts**. And shut down gracefully: on SIGTERM, stop accepting, drain in-flight requests, then exit:

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

go func() { errCh <- srv.ListenAndServe() }()

<-ctx.Done()
shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_ = srv.Shutdown(shutCtx) // drain, then return
```

:::info[Platform connection]
Both halves are packaged in `dx-common-go`: `middleware.StandardStack()` applies the exact stack above (GO-SERVICE-STANDARDS makes the order normative), and `httpserver.New(cfg, router, logger).Start()` runs the timeout-configured server, blocks until SIGINT/SIGTERM, and drains — every service's `main` ends with that one call. Kubernetes sends SIGTERM on pod rotation; graceful drain is why deploys don't drop requests.
:::

## Exercises

1. Build a chi router for your `dx-scratch-go` notes service: CRUD routes in a `Route` group, health endpoints outside it, URL params via `chi.URLParam`.
2. Write three middleware — request ID, logger, recoverer — from scratch (no chi builtins), stack them in the standard order, and prove order matters by swapping RequestID/Logger and reading the logs.
3. Make a handler panic. Verify Recoverer turns it into a 500 with a logged stack while the process lives on.
4. Wire full graceful shutdown; verify with `curl` that an in-flight slow request completes while Ctrl-C is pending, then the process exits.

## Check yourself

- Why must RequestID precede Logger, and Recoverer sit below it?
- What is `r.Use(A, B)` in function-composition terms?
- What's wrong with bare `http.ListenAndServe`?
- What does graceful shutdown buy during a Kubernetes deploy?

## References

- [net/http docs](https://pkg.go.dev/net/http) · [chi](https://pkg.go.dev/github.com/go-chi/chi/v5)
- [Go Blog: Writing Web Applications](https://go.dev/doc/articles/wiki/)
- Platform: `dx-common-go/middleware` (StandardStack), `dx-common-go/httpserver`, `dx-acl-go/internal/api/router.go`
