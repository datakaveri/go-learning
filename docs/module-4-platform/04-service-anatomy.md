---
title: Anatomy of a Service — dx-acl-go
sidebar_label: Service Anatomy
description: A guided line-by-line walkthrough of the reference service, from main.go to the outbox dispatcher.
---

# Anatomy of a Service — dx-acl-go

## Learning objectives

- Walk `dx-acl-go` top to bottom and account for every structural decision.
- See the boot contract, layered architecture, and outbox implemented for real.
- Acquire the reading order you'll use on *any* unfamiliar DX service.

## Prerequisites

- [dx-common-go Tour](dx-common-go-tour). Have `dx-acl-go` open; this page is a reading companion, not a substitute.

## Time estimate

**4 hours** of guided source reading.

## Concepts

### Why acl-go is the reference

It's a complete specimen: HTTP API + Postgres persistence + transactional outbox + RabbitMQ publishing + an external client + a background worker — every platform pattern in one medium-sized codebase. It's also architecturally central (the PAP from the [deep dive](architecture-deep-dive)). Learn to read this one and the other fourteen are variations.

### The reading order (use it on every unfamiliar service)

1. **README** — API table, env vars, events. The service's self-declaration.
2. **`cmd/server/main.go`** — the wiring reveals every component and dependency.
3. **`internal/api/router.go`** — the API surface and its middleware.
4. **One vertical slice** — a single endpoint, handler → service → repository → SQL.
5. **The background pieces** — workers, consumers, dispatchers.

### main.go — the boot contract, annotated

Read it now; here is what to notice at each step:

```
cfg, err := config.Load()            // dxconfig.LoadService[Config] — nothing hand-rolled
logger := …                          // zap, then defer Sync()
pool, err := …; pool.Ping(ctx)       // HARD dep → logger.Fatal on failure
schema ensure                        // embedded idempotent SQL — additive only,
                                     // legacy tables untouched (the governing principle)
publisher, err := rabbitmq…          // OPTIONAL dep → Warn + no-op fallback, service still boots
repo := postgres.NewPolicyRepo(pool) // repositories take the pool
dispatcher := …; go dispatcher.Run(ctx) // outbox drainer: supervised, ctx-cancelled
svc := service.NewPolicyService(repo, catalogueClient, dispatcher, …)
h := api.NewHandler(svc, logger)
router := api.NewRouter(…)
httpserver.New(cfg.Server, router, logger).Start() // blocks; SIGTERM → drain → exit
```

Every arrow of [Dependency Injection](../module-2-intermediate/dependency-injection) and both halves of the hard/optional split, in maybe a hundred lines. Note what's *absent*: no globals, no init() magic, no framework.

### router.go — the API surface

```go
r.Use(StandardStack...)                          // the ordered seven
r.Use(dxopenapi.ValidationMiddleware(spec, cfg)) // spec-enforced requests
r.Route("/iudx/acl/apd/v2", func(r chi.Router) { // legacy path preserved — clients own contracts
	r.Use(authresolver.Middleware(...))          // HMAC → JWT → anonymous
	r.Use(auditing.Middleware(...))              // mutating calls → audit events
	r.Post("/policies", h.CreatePolicy)
	...
})
r.Get("/healthz/live", …); r.Get("/healthz/ready", …); r.Handle("/metrics", …)
```

Ops endpoints sit *outside* the auth group — Kubernetes probes don't carry tokens. The base path is the documented legacy exception ([REST](../module-3-advanced/rest-api-development)).

### A vertical slice: create policy

Follow `POST /policies` down the layers and watch each one keep its lane:

- **Handler** (`internal/api/handler.go`): decode → call service → translate result to envelope/error URN. No business logic — it doesn't know what a policy *means*.
- **Service** (`internal/service/policy_service.go`): the use case — validates business rules, calls the catalogue client to verify the item exists (an external gRPC dependency, injected as an interface), then asks the repository to persist **policies plus outbox rows**.
- **Repository** (`internal/repository/postgres/policy_repo.go`): `InsertPoliciesWithOutbox` — one transaction, `defer tx.Rollback(ctx)`, parameterized inserts, outbox rows with action + JSON payload. Column names match the legacy Java schema (foreign keys into `user_table` — the governing principle at the row level).
- **Domain** (`internal/domain/policy.go`): the entities everything above shares.

### The background half: outbox dispatcher

`internal/service/outbox_dispatcher.go` is your [Workers & Cron](../module-3-advanced/workers-cron) page running in production: a ticker loop, `select` on `ctx.Done()`, each cycle reading unsent outbox rows, publishing `policy.*` to the `authz` exchange, marking sent. Failures are logged and retried next cycle — the dispatcher never crashes the service, and rows survive process restarts by construction ([Transactions](../module-3-advanced/transactions)).

### The tests worth imitating

Note the pattern of `spec_test`-style checks (embedded spec parses; every route is in the spec) and router-level handler tests with the middleware mounted — the [Testing](../module-2-intermediate/testing) norms, applied. When you write your capstone, copy these shapes.

:::info[Platform connection]
After this walkthrough, prove the transfer: open `dx-registry-go` (a service you haven't studied) and run the same five-step reading order. Time yourself — a standards-shaped service should yield its structure in under twenty minutes. That uniformity dividend is why the standards exist, and why deviations (files-connect's boot, catalogue's env prefix) are tracked as debt rather than tolerated as variety.
:::

## Exercises

1. Do the full reading order on `dx-acl-go`, producing a one-page map: every component in `main.go`, every route group, the create-policy slice, the dispatcher lifecycle.
2. Answer from source (cite file + function): Where does a nonexistent item turn into a client-visible error, and which URN? What exactly is in an outbox row's payload? How does the dispatcher mark rows sent, and what happens if publish succeeds but the mark fails? (Connect that last answer to consumer idempotency.)
3. Run the transfer test on `dx-registry-go` above. Note the two most significant differences you find and *why* they exist (hint: what does registry not need?).
4. Sketch (don't build) how you'd add a `PATCH /policies/{id}` endpoint: which files change, what event publishes, what the OpenAPI diff is, which tests you'd add. Check your sketch against the standards page that follows.

## Check yourself

- Recite the five-step reading order.
- Which layer knows SQL? Which knows URNs? Which knows what a policy means?
- Why does the dispatcher's publish-then-mark ordering require idempotent consumers rather than the reverse ordering?
- What makes ops endpoints exempt from auth middleware?

## References

- Platform: `dx-acl-go` source (primary text); `dx-notification-go/internal/consumer/consumer.go` (the consumer-shaped service for contrast)
