---
title: dx-common-go — The Shared Library Tour
sidebar_label: dx-common-go Tour
description: Package-by-package tour of the shared library every DX service is built from.
---

# dx-common-go — The Shared Library Tour

## Learning objectives

- Know what exists in `dx-common-go` so you never rebuild it in a service.
- For each package: what it does, its key entry points, and which curriculum page taught the underlying concept.
- Know the criteria for adding to the library vs keeping code in a service.

## Prerequisites

- [Repos & Workflow](repo-structure-workflow). Have `dx-common-go` open in your editor while reading.

## Time estimate

**4 hours** — a guided source-reading session. Don't skim; open every file named here.

## Concepts

### Why a shared library

Fifteen services × (config loading + middleware + envelopes + health + DB plumbing + RMQ reconnection) would be fifteen slightly-different copies of everything, each with its own bugs. `dx-common-go` is where a pattern goes once it's needed twice. It's also why the standards are *enforceable*: "use `StandardStack`" is checkable in review; "write middleware in the correct order yourself" is not.

The rule of admission: **shared capability, not shared business logic.** HTTP plumbing, persistence helpers, messaging reliability — yes. Anything with domain vocabulary (policies, datasets) — no; that belongs to its service.

### The tour

Packages grouped as you'll encounter them, with the concept-page each builds on:

#### Boot & HTTP layer

| Package | What it gives you | Taught in |
|---|---|---|
| `config` | `LoadService[T]` — defaults → baked YAML → env, into your struct | [Configuration](../module-2-intermediate/configuration) |
| `httpserver` | `New(cfg, router, logger).Start()` — timeouts, SIGTERM, graceful drain | [HTTP & Middleware](../module-3-advanced/http-servers-middleware) |
| `middleware` | `StandardStack()` (the ordered seven), plus `RateLimit`, `MaxUploadSize`, upload validation | same |
| `openapi` | spec loading (embedded), request-validation middleware, Swagger UI at `/docs` | [REST](../module-3-advanced/rest-api-development) |

#### Contract layer

| Package | What it gives you | Taught in |
|---|---|---|
| `response` | `DxResponse[T]`, `DxPagedResponse[T]`, writers (`WriteSuccess`, `WriteCreated`, `WritePaginated`), `ServiceWriter` with per-service URN prefix | [REST](../module-3-advanced/rest-api-development) |
| `errors` | the taxonomy (`NewValidation`, `NewNotFound`, `NewConflict`, …) → HTTP status + URN mapping, `WriteError` | [Error Handling](../module-1-go-fundamentals/error-handling) |
| `request` / `pagination` / `validation` | query-param parsing with bounds (`request.Builder`), limit/offset model, input validators | [REST](../module-3-advanced/rest-api-development) |

#### Identity & trust

| Package | What it gives you | Taught in |
|---|---|---|
| `auth` | `DxUser` (id, email, roles, org, delegation), context accessors | [Context](../module-2-intermediate/context) |
| `auth/jwt` | RS256 validation, JWKS fetch + cache + refresh, claim checks | [AuthN & AuthZ](../module-3-advanced/authn-authz) |
| `auth/resolver` | the middleware chain: HMAC-verify → JWT fallback → anonymous; `RequireGatewayOrigin` | same |
| `auth/authorization` | role/org access-control helpers | same |
| `trust` / `crypto` | mTLS verification, hashing/encryption utilities | [Security Model](security-model) |

#### Data layer

| Package | What it gives you | Taught in |
|---|---|---|
| `database/postgres` | pool factory, `InTransaction` (context-propagated tx), `InRetryableTransaction`, `WithAdvisoryLock` | [Database Patterns](../module-3-advanced/database-patterns), [Transactions](../module-3-advanced/transactions) |
| `database/postgres/repository` | **embed this** — `Base[R]`: tx-aware generic CRUD + fluent `Query`, the service-repo standard | same |
| `database/postgres/dao` | `BaseDAO[T]`: CRUD, upsert, soft-delete+`Restore`, audit columns, `UpdateVersioned`, bulk, `FindPage` → `Page[T]`, fluent `Finder` | same |
| `database/postgres/query` | builder + spec constructors (`query.Eq/In/And/...`) — typed operators, safe ORDER BY | same |
| `database/postgres/migrate` | golang-migrate runner: embedded migrations, `schema_mode` gate, dirty-state reporting | [Schema Migrations](../module-3-advanced/schema-migrations) |
| `database/redis`, `cache` | go-redis v9 client wrapper, caching helpers | — |
| `database/elastic` | go-elasticsearch v8 client + query DSL builder (catalogue, dataplane-rs) | — |

#### Async & observability

| Package | What it gives you | Taught in |
|---|---|---|
| `messaging/rabbitmq` | connection manager with reconnect + backoff, `ReliablePublisher` (confirms), consumer scaffold (declare, bind, deliver loop) | [RabbitMQ](../module-3-advanced/event-driven-rabbitmq) |
| `auditing` | async audit-event publisher + the audit middleware (mutating endpoints → audit records → queue) | [Security Model](security-model) |
| `health` | liveness/readiness handler, `Checker` interface, `PgxPoolChecker` | [Observability](../module-3-advanced/observability) |
| `metrics` | Prometheus handler + `RequestMetrics` middleware | same |

Plus the long tail — email notifier, S3/appid clients — inventoried in `GO-PLATFORM-REVIEW.md` (~28 packages total).

### Known gaps — don't be surprised, don't fill them ad hoc

The platform review records what the library *doesn't* have yet: OTel tracing wiring, a generalized outbox helper (acl's is hand-built), a consumer framework (DLQ/retry boilerplate is still repeated), HTTP retry/backoff helpers, bulk-ops and optimistic-locking DAO methods. Two implications: (1) don't hunt for these and conclude you're blind; (2) if your work needs one, that's a *library* conversation with the team — exactly the admission rule above — not a private helper in your service.

:::info[Platform connection]
The fastest way to internalize the library is to watch it being consumed: split your screen with `dx-acl-go/cmd/server/main.go` on one side and the packages it imports on the other, and follow each import to its definition. Nearly every line of that main is a row from the tables above.
:::

## Exercises

1. **Inventory check**: open every package in the boot/contract/identity tables and, for each, write one line: main entry point + one design choice you noticed. (Yes, all of them — this exercise *is* the page.)
2. Read `dao/base.go` end to end and answer: how does `FindPage` build its two queries? Where does the soft-delete filter apply? What would bulk-insert require?
3. Read `messaging/rabbitmq`'s reconnect loop: what triggers it, what's the backoff, what happens to a consumer mid-reconnect?
4. Migrate `dx-scratch-go` from your handwritten helpers to the real library (envelope writers, StandardStack, health, config loader). Count the lines you delete.
5. Judgment reps: for each, library or service? — a URL slugifier used by one service; retry-with-backoff for outbound HTTP; marketplace price calculation; a `FindByIDs` DAO batch method. Justify each in a sentence.

## Check yourself

- What's the admission rule for the library, in one sentence?
- Which package would you reach for to: bound a query's page size? verify a gateway signature? survive a broker restart? add `/metrics`?
- Name three recorded gaps and the correct process when you hit one.

## References

- Platform: `dx-common-go` source (the reading *is* the reference); `claude-docs/GO-PLATFORM-REVIEW.md` (full inventory + findings)
