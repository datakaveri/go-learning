---
title: Dependency Injection Patterns
sidebar_label: Dependency Injection
description: Constructor injection, wiring in main, and why the platform uses no DI framework.
---

# Dependency Injection Patterns

## Learning objectives

- Practice constructor injection: dependencies in, concrete type out.
- Wire an application in `main` — the composition root — following the platform's boot contract.
- Distinguish hard dependencies (fail fast) from optional ones (degrade gracefully).
- Explain why DX uses no DI framework or global state.

## Prerequisites

- [Interfaces](../module-1-go-fundamentals/interfaces) (consumer-defined interfaces), [Context](context)

## Time estimate

**3 hours**

## Concepts

### Constructor injection

Go's DI "framework" is a function:

```go
type PolicyService struct {
	store    PolicyStore   // interface defined HERE, by the consumer
	notifier Notifier
	log      *zap.Logger
}

func NewPolicyService(store PolicyStore, notifier Notifier, log *zap.Logger) *PolicyService {
	return &PolicyService{store: store, notifier: notifier, log: log}
}
```

Everything the type needs arrives through `New`. No globals, no service locator, no annotations. Consequences worth internalizing:

- **Testability is free** — hand in a fake `PolicyStore` and the service never knows.
- **The dependency graph is explicit** — it's the parameter lists.
- **Construction can't half-succeed** — if `New` returns, the object is usable.

### main is the composition root

All wiring happens once, at the top, in `cmd/server/main.go`. The platform's boot contract makes the order itself a standard:

```go
func main() {
	cfg, err := config.Load()               // 1. config first
	if err != nil { log.Fatal(err) }

	logger := newLogger(cfg)                // 2. logger second
	defer logger.Sync()

	pool, err := postgres.NewPool(ctx, cfg.Postgres) // 3. HARD deps: Fatal on failure
	if err != nil { logger.Fatal("postgres", zap.Error(err)) }

	publisher, err := rabbitmq.NewPublisher(cfg.RMQ) // 4. OPTIONAL deps: Warn + no-op
	if err != nil {
		logger.Warn("rabbitmq unavailable, events disabled", zap.Error(err))
		publisher = events.NewNoopPublisher()
	}

	repo := postgres.NewPolicyRepo(pool)    // 5. repositories
	svc := service.NewPolicyService(repo, publisher, logger) // 6. services
	handler := api.NewHandler(svc, logger)  // 7. handlers
	router := api.NewRouter(handler, ...)   // 8. router

	httpserver.New(cfg.Server, router, logger).Start() // 9. blocks until SIGTERM
}
```

The hard/optional distinction is a deliberate design decision, not improvisation: **a service that can't reach Postgres is useless — fail fast**; a service that can't reach RabbitMQ can still serve reads — **degrade, warn, and keep the readiness probe honest**. The no-op publisher satisfying the same interface is the null-object pattern doing real work.

### Why no DI framework

Java engineers ask this in week one. Reasons the platform (and most of the Go ecosystem) wires by hand:

- A service has perhaps 10–20 components; the wiring is a page of obvious code you can read top to bottom.
- Compile-time safety: a missing dependency is a build error, not a reflection-time surprise at startup.
- Nothing to learn, nothing to debug through. Stack traces contain your functions.

(Google's `wire` generates such wiring for very large graphs; DX services aren't large enough to want it.)

### Anti-patterns to unlearn

- **Package-level singletons** (`var DB *pgxpool.Pool` + `init()`): hidden coupling, impossible to test in parallel, initialization-order puzzles.
- **Service locator** (`container.Get("policyService")`): stringly-typed, runtime failures.
- **Config structs passed everywhere**: pass each component only the narrow config *it* needs (`cfg.Postgres`, not `cfg`).

:::info[Platform connection]
Open `dx-acl-go/cmd/server/main.go` and you'll find the boot contract verbatim — config → logger → pool (Fatal) → RabbitMQ (Warn + no-op fallback) → schema ensure → repos → services → handlers → router → `httpserver.Start()`. Every other service follows the same shape, which is why you can open a service you've never seen and know where everything is wired. GO-SERVICE-STANDARDS.md codifies the order; deviation is a review finding.
:::

## Exercises

1. Take the pluggable-storage mini-project from [Interfaces](../module-1-go-fundamentals/interfaces) and restructure it to strict constructor injection with all wiring in `main`.
2. Add a `Notifier` dependency with a real implementation (writes to a file) and a no-op fallback selected when a `-notify` flag is absent — the optional-dependency pattern.
3. Write a unit test for a service using a hand-written fake store (no mocking library — a struct with function fields is enough).
4. Deliberately create an import cycle between `service` and `repository` packages, read the compiler error, and resolve it with a consumer-defined interface — this is *why* the interface lives on the consumer side.

## Check yourself

- What distinguishes a hard dependency from an optional one at boot, per the platform?
- Why does the no-op publisher exist instead of `if publisher != nil` checks everywhere?
- Three reasons Go projects skip DI containers?
- Where is the only place `New*` functions get called with real dependencies?

## References

- [Google Go Style Guide — Global state](https://google.github.io/styleguide/go/best-practices#global-state)
- [Uber Go Style Guide — Avoid init()](https://github.com/uber-go/guide/blob/master/style.md#avoid-init)
- Platform: `dx-acl-go/cmd/server/main.go` — the canonical composition root; `claude-docs/GO-SERVICE-STANDARDS.md` (boot contract)
