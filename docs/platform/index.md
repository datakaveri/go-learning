---
title: Shared Go platform modules
description: Responsibility map for dx-common-go/platform and adjacent foundation packages.
---

# Shared Go platform modules

**Status: Partially implemented and actively adopted.** Package names on this page were verified in `dx-common-go` source. A capability that does not exist under `platform/` is not presented as if it does.

`dx-common-go/platform` is the infrastructure kernel consumed by CDPG services. It standardizes boot, transports, persistence mechanics, delivery, identity propagation, health, and lifecycle. It must not contain catalogue rules, grant semantics, billing logic, query policy, or another service’s data model.

![Shared Go platform module map](/img/architecture/shared-platform-modules.svg)

## Module inventory

| Module | Owns | Does not own | Status |
|---|---|---|---|
| `platform/bootstrap` | Boot order, dependency construction, serving, supervised workers, drain and close | Business object graph or use cases | Implemented |
| `platform/config` | Typed file/env binding, platform defaults, validation hooks | Secrets or remote configuration authority | Implemented |
| `platform/http` | Typed handlers, route tables, envelopes, validation, error mapping, middleware, health/docs mounting | Route business semantics or authorization policy | Implemented |
| `platform/grpc` and `platform/grpc/server` | Identity-aware clients and unary server/interceptors | Public API ownership or streaming policy | Partially implemented |
| `platform/errors` | Stable error classes and mapping primitives | User-facing business copy | Implemented |
| `platform/paging` | Bounded page requests/results | Query planning | Implemented |
| `platform/database/sql` | Pool-neutral SQL interfaces, transaction propagation, query/repository helpers | Schema/domain design | Implemented |
| `platform/database/sql/pgx` | Explicit driver escape hatch | General repository API | Implemented |
| `platform/cache` and `platform/cache/redis` | Cache contract, memory/Redis adapters, cache-aside helpers | Durable truth or business invalidation facts | Implemented |
| `platform/events` and `platform/events/amqp` | Typed event envelope, bus, outbox, dispatcher, retry/DLQ/replay adapter | Business event vocabulary or consumer side effects | Implemented; fleet adoption partial |
| `platform/idempotency` | Durable idempotency record mechanics | Deciding which operation needs an idempotency key | Implemented |
| `platform/lease` | Durable ownership, renewal, interruption, loss signaling | Work definition or retry policy | Implemented |
| `platform/executor` | Lifecycle-owned per-request background execution | Unlimited fire-and-forget work | Implemented |
| `platform/security/identity` | Verified subject/actor/delegation value and context | Credential verification or authorization | Implemented |
| `platform/security/workload` | Destination audience verification, caller/subject-asserter control | User login or business permission | Implemented |
| `platform/security/workload/issuer` | Client credentials and short destination-token cache | Subject authorization | Implemented |
| `platform/observability/health` | Liveness/readiness registry and dependency checks | Full telemetry backend | Implemented |

Search, Elasticsearch, S3 storage, JWT/JWKS user authentication, metrics, tracing, audit helpers, notification/email, resilience, and test utilities currently exist as focused top-level foundation packages rather than as complete `platform/*` modules. Use their present source-backed APIs where a current service already proves the pattern; do not invent `platform/search`, `platform/storage`, `platform/authz`, or `platform/test` imports.

## Responsibility matrix

| Capability | Shared platform owns | Service owns |
|---|---|---|
| Bootstrap | Correct lifecycle and shutdown order | Dependency declaration and wiring |
| Configuration | Binding, precedence, shared defaults | Service fields, validation, safe values |
| HTTP/gRPC | Transport mechanics and cross-cutting middleware | Contracts, operation semantics, route posture |
| Errors | Classification and rendering | Domain meaning and contextual wrapping |
| SQL | Pool/transaction/query mechanics | Schema, queries, invariants, isolation choice |
| Cache | Storage abstraction and cache-aside mechanics | Keys, TTL, invalidation events, safe fallback |
| Events | Envelope, delivery adapter, outbox mechanics | Event names, versions, payloads, side effects |
| Workers | Supervision, leases, cancellation primitives | Work selection, retry classification, reconciliation |
| Identity | Verified context representation and boundary verification | Required actor/subject semantics for an operation |
| Authorization | Reusable clients/types when implemented | Resource/action mapping and enforcement |
| Observability | Common signal mechanics | Business dimensions, SLOs, alerts, runbooks |
| Testing | Shared fakes/harnesses where available | Domain, contract, security, recovery, and workflow tests |

## Consumption rule

Pin a reviewed `dx-common-go` version. During local workspace development, a `go.work` file or temporary `replace` may point to the adjacent clone. Do not publish a service release whose module file points to a developer’s filesystem.

Upgrade platform code upstream first: add or fix the general capability in `dx-common-go`, test it there, release it, then adopt it in services. Copying shared infrastructure into one repository creates a second behavior and a second security patch path.

Continue with [bootstrap and configuration](./bootstrap-config.md), [HTTP and gRPC](./http-grpc.md), [persistence and cache](./persistence-cache.md), [events and workers](./events-workers.md), and [identity and operations](./identity-operations.md).

