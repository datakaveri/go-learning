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

Persistence is organized by **backend**, each backend a directory grouping of small single-job packages. Postgres is seven sub-packages (see [Database Patterns](../module-3-advanced/database-patterns) and [SQLC vs the DSL](../module-3-advanced/persistence-sqlc-vs-dsl)); Elasticsearch — restructured 2026-07-05 from the old flat `database/elastic` into the nested `database/elasticsearch/{client,query,repository,mapping,indexing}` (see [Elasticsearch](../module-3-advanced/elasticsearch)) — is five; Redis stays flat. `database/postgres` and `database/elasticsearch` are pure directory groupings with no files of their own, and they **never import each other** (no cross-store dependency).

| Package | What it gives you | Taught in |
|---|---|---|
| `database/postgres/client` | Pool factory (`NewPool`), DSN/pool-size config, pgx tracer composition (`MultiTracer`, `SlowQueryTracer`) | [Database Patterns](../module-3-advanced/database-patterns) |
| `database/postgres/transaction` | `WithTransaction`/`InTransaction` (context-propagated ambient transaction), `InRetryableTransaction`, `WithAdvisoryLock` | [Database Patterns](../module-3-advanced/database-patterns), [Transactions](../module-3-advanced/transactions) |
| `database/postgres/dao` | `BaseDAO[T]`: CRUD, upsert, soft-delete, optimistic locking (`UpdateVersioned`), bulk (`InsertMany`/`CopyFrom`), `FindPage` → `Page[T]`, and `Finder[T]` — the fluent DSL (`Where`/`OrderBy`/`Join`/`Select`/`GroupBy`/`Having`) | same |
| `database/postgres/query` | the DSL's building blocks — typed `Condition` operators, `Join`, safe `ORDER BY`, the `SQLBuilder` that renders them | same |
| `database/postgres/repository` | `Base[R]` — embeddable generic repository wrapper around `BaseDAO[T]`, options-based `New` (`WithTable`/`WithID`/`WithDAOOption`), `NewWithSQL` for the SQLC escape hatch | same |
| `database/postgres/migrate` | embedded, versioned schema migrations (`Run`/`Status`) for a service's own net-new tables | [Schema Migrations](../module-3-advanced/schema-migrations) |
| `database/postgres/sqlcx` | transaction-propagation-aware DBTX provider for SQLC-generated queries | same |
| `database/redis`, `cache` | go-redis v9 client wrapper, caching helpers (`cache.GetOrLoad[T]` — singleflight-protected cache-aside) | — |
| `database/elasticsearch/{client,query,repository,mapping,indexing}` | the reusable ES module: `esclient.New`, the query DSL, typed `Repo[T]`, mappings-as-code + blue/green lifecycle, bulk/`Sync`/`Worker` (catalogue, dataplane-rs) | [Elasticsearch](../module-3-advanced/elasticsearch) |
| `dxtest/containers`, `dxtest/fixtures` | testcontainers `Postgres`/`Redis` with `DX_TEST_PG_DSN` fallback (`WithMigrations`/`WithSetupSQL`) + the shared fixture schema — imported only by test binaries | [Testing Strategy](testing-strategy) |

#### Async & observability

| Package | What it gives you | Taught in |
|---|---|---|
| `messaging/rabbitmq` | connection manager with reconnect + backoff, `ReliablePublisher` (confirms), consumer scaffold (declare, bind, deliver loop) | [RabbitMQ](../module-3-advanced/event-driven-rabbitmq) |
| `auditing` | async audit-event publisher + the audit middleware (mutating endpoints → audit records → queue) | [Security Model](security-model) |
| `health` | liveness/readiness handler, `Checker` interface, `PgxPoolChecker` | [Observability](../module-3-advanced/observability) |
| `metrics` | Prometheus handler + `RequestMetrics` middleware | same |

Plus the long tail — email notifier, S3/appid clients — inventoried in `GO-PLATFORM-REVIEW.md` (~28 packages total).

### Known gaps — don't be surprised, don't fill them ad hoc

Some gaps the platform review once flagged are now closed — bulk-ops and optimistic-locking DAO methods (`InsertMany`/`CopyFrom`/`UpdateVersioned`) have been in `dao/` for a while, OTel tracing wiring now exists (`observability.Init` + `client.WithTracers`, see [Observability](../module-3-advanced/observability)), and a generalized outbox-plus-scheduler pattern replaced acl's original hand-built dispatcher (`messaging/outbox.Dispatcher` + `scheduler.Runner`, see [Workers & Cron](../module-3-advanced/workers-cron)). Still genuinely open: a shared RabbitMQ consumer framework (DLQ/retry boilerplate is still repeated per service) and HTTP retry/backoff helpers for outbound calls. Two implications: (1) don't hunt for the closed ones and conclude you're blind — check the package tables above first; (2) if your work needs something from the still-open list, that's a *library* conversation with the team — exactly the admission rule above — not a private helper in your service.

:::info[Platform connection]
The fastest way to internalize the library is to watch it being consumed: split your screen with `dx-acl-go/cmd/server/main.go` on one side and the packages it imports on the other, and follow each import to its definition. Nearly every line of that main is a row from the tables above.
:::

## Exercises

1. **Inventory check**: open every package in the boot/contract/identity tables and, for each, write one line: main entry point + one design choice you noticed. (Yes, all of them — this exercise *is* the page.)
2. Read `dao/base.go`, `dao/count.go`, and `dao/finder.go` end to end (the package is split into topic files now — CRUD/count/batch/finder each get their own) and answer: how does `FindPage` build its two queries? Where does the soft-delete filter apply? What does `InsertMany` vs `CopyFrom` each require, and when would you pick one over the other?
3. Read `messaging/rabbitmq`'s reconnect loop: what triggers it, what's the backoff, what happens to a consumer mid-reconnect?
4. Migrate `dx-scratch-go` from your handwritten helpers to the real library (envelope writers, StandardStack, health, config loader). Count the lines you delete.
5. Judgment reps: for each, library or service? — a URL slugifier used by one service; retry-with-backoff for outbound HTTP; marketplace price calculation; a `FindByIDs` DAO batch method. Justify each in a sentence.

## Check yourself

- What's the admission rule for the library, in one sentence?
- Which package would you reach for to: bound a query's page size? verify a gateway signature? survive a broker restart? add `/metrics`?
- Name the gaps that are still genuinely open, the ones that have since closed, and the correct process when you hit one that's still open.

## References

- Platform: `dx-common-go` source (the reading *is* the reference); `claude-docs/GO-PLATFORM-REVIEW.md` (full inventory + findings)
