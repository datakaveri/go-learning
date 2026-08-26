---
title: API and transport conventions
description: Public REST, internal gRPC, contracts, errors, streaming, files, and compatibility.
---

# API and transport conventions

**Status: Partially implemented.** Public REST conventions are implemented broadly. Internal gRPC is the target and has shared unary infrastructure, but fleet adoption and streaming interception are incomplete.

![Public REST and internal gRPC boundaries](/img/architecture/rest-grpc-boundaries.svg)

## Choose the boundary

| Need | Transport | Rule |
|---|---|---|
| Browser, CLI, partner, or public platform API | HTTPS REST/JSON through `dx-gateway-go` | Service owns OpenAPI and public compatibility. |
| Internal synchronous service contract | gRPC/protobuf with workload identity | Caller sets a destination-bound audience; callee allowlists callers where needed. |
| Domain fact or projection update | RabbitMQ event | At-least-once, versioned envelope, idempotent consumer, reconciliation. |
| OGC, NGSI-LD, webhook, blob, or SSE representation | Standards-native HTTP | Use `httpx.HandleRaw`; do not wrap an external standard in the CDPG envelope. |

Do not use events for a request that needs an immediate answer. Do not use synchronous calls to broadcast facts to multiple owners.

## Public REST

- Use resource-oriented nouns and an explicit versioned base path.
- Give every operation a stable, unique OpenAPI `operationId` and the matching `httpx.OpID`.
- Treat OpenAPI as source-controlled contract: validate requests, test route/spec drift, and review breaking changes.
- Use platform binding and validation tags. Handlers must not manually decode the same JSON a platform adapter understands.
- Return the platform success envelope for owned JSON APIs and the platform problem/error mapping for failures.
- Use `httpx.Created[T]` for 201, `httpx.Accepted[T]` for 202, and `HandleVoid` for 204.
- Use cursor pagination when data changes frequently; if an existing contract uses offset/limit, enforce bounded limits and deterministic ordering.
- Accept an idempotency key for externally retried side effects such as purchase, payment callback, subscription creation, or high-risk agent execution. Persist the result in the side-effect owner.

## Errors

Classify once in the domain/application layer, wrap with `%w`, and map at the edge:

| Class | HTTP | gRPC | Retry |
|---|---:|---|---|
| Invalid input | 400 | `InvalidArgument` | No |
| Missing/invalid authentication | 401 | `Unauthenticated` | Only after obtaining valid credentials |
| Denied authorization | 403 | `PermissionDenied` | No unless policy changes |
| Not found | 404 | `NotFound` | No |
| Conflict/idempotency mismatch | 409 | `AlreadyExists` or `Aborted` | Conditional |
| Dependency unavailable | 503 | `Unavailable` | Bounded retry if idempotent |
| Deadline | 504 | `DeadlineExceeded` | Conditional |
| Unclassified internal failure | 500 | `Internal` | Do not expose details |

Return stable machine codes and a request/trace reference. Never expose SQL, stack traces, credentials, provider payloads, or policy internals.

## Correlation, timeouts, and identity

Propagate W3C trace context and the platform request ID. Set a client deadline; the server derives work from the request context. Public identity and represented-subject context are established at trusted middleware, not accepted from arbitrary JSON fields.

Non-streaming routes use finite request and upstream timeouts. `httpx.Streaming()` exempts a route from the ordinary timeout/compression path. Streaming handlers must flush, stop on cancellation, bound buffers, and support resumability where the contract promises it.

## Files and webhooks

Stream large files; never buffer an unbounded body. Enforce content length, content type, object ownership, malware/processing state, and authorization before issuing a short-lived capability. Treat a presigned URL as a credential and never log it.

A public webhook needs an exact route and provider-specific authenticity verification. “Public” means no user OIDC token; it does not mean unauthenticated business input. Persist provider event IDs for idempotency.

## Internal gRPC

- Define protobuf contracts under `api/proto` and generate code; do not hand-edit generated files.
- Preserve field numbers, reserve removed fields/names, add fields compatibly, and avoid changing enum meanings.
- Attach workload credentials whose audience is the destination service.
- Use standard unary interceptors for workload verification, trace/correlation, logging, metrics, recovery, and error mapping.
- Declare deadlines at the caller and honor cancellation at every database or network operation.
- Add Buf/protobuf compatibility checks and client/server contract tests.

The platform gRPC server is currently unary-oriented. A streaming internal contract requires a reviewed extension to shared interceptors and shutdown behavior; it must not silently bypass identity or telemetry.

## Contract verification

Every public operation needs handler tests, OpenAPI request/response validation, route/spec drift tests, and gateway route tests. Every internal operation needs protobuf compatibility and workload-identity tests. See the [testing strategy](../operations/testing.md).

