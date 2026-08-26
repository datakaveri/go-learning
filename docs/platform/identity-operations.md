---
title: Identity, health, observability, resilience, and test support
description: Shared security identity and operational module behavior.
---

# Identity and operational foundations

**Status: Partially implemented.** Workload and subject context, health, platform middleware, metrics/tracing foundations, resilience helpers, and test packages exist; unified adoption and several target authorization types remain incomplete.

## Subject and actor

`platform/security/identity.Subject` represents a verified human subject and optional delegation. `Subject.ID` remains the human whose rights are being exercised. `Subject.Actor()` returns the actual actor for audit/rate attribution. When an agent or application acts for a person, effective authority is the intersection of the subject’s rights and delegation scopes.

Handlers obtain identity through typed request embedding or `identity.Require(ctx)`. Business JSON must not create a trusted `Subject`.

## Workload verification

`platform/security/workload` verifies an internal credential’s issuer, signature, time bounds, and destination audience (`dx:svc:<service>`). `AllowedCallers` optionally narrows the call graph. `SubjectAsserters` names the few workloads allowed to forward represented subject context; an empty list means none.

Enforcement is explicitly `required` or `off`; it has no permissive mode and no empty default. An authenticated surface uses `required`.

`platform/security/workload/issuer` obtains client-credentials tokens per destination and keeps a short, bounded cache. Use `TokenFor(ctx, destination)`, the provided authorization helper/transport, or the gRPC client integration. Do not mint a token once and reuse it for unrelated audiences.

See [Authentication and workload identity](../integrations/identity.md) for the trust flow and current binding gap.

## Health and readiness

`platform/observability/health.Registry` separates:

- liveness: process event loop is alive; do not make it depend on every remote service;
- readiness: service can meet its advertised contract now;
- dependency checks: registered automatically for declared dependencies and explicitly with `App.Probe` for service-specific readiness.

A degraded optional dependency affects readiness only if the documented degraded contract cannot be served correctly.

## Logging, metrics, and tracing

The current platform logger is `zap`. HTTP middleware establishes request ID and trace context, structured request logs, and standard telemetry. OpenTelemetry and Prometheus foundations exist outside the narrow `platform/observability/health` package.

Every service adds domain signals without leaking protected data. Correlate `request_id`, `trace_id`, workload, subject, actor, organisation, resource reference, authorization decision ID when available, event ID, and result. High-cardinality IDs belong in traces/logs, not metric labels.

## Resilience

Use explicit timeouts, bounded retries with jitter, circuit breaking only where its state and fallback are understood, idempotency at the owner of an effect, and backpressure rather than unbounded queues. Authentication, authorization, and decision verification fail closed. A cache or broker outage may degrade only if authoritative correctness remains intact.

## Test utilities

`dx-common-go` currently exposes focused top-level test support, including `dxtest` and package-local fakes/harnesses. There is no generic `platform/test` package to import. Prefer:

- in-memory adapters shipped with a platform capability for unit tests;
- workload verifier/issuer harnesses for trust-boundary tests;
- real dependency containers for SQL, AMQP, OpenFGA, Redis, Elasticsearch, or S3 semantics;
- `httptest`, in-process routers, and gRPC local listeners for transport tests.

Do not let a broad mock framework conceal whether the real middleware, transaction, or message acknowledgement path runs.

## Failure behavior

- Missing verified subject on a protected operation: deny before use case.
- Workload token absent, invalid, wrong audience, expired, or disallowed caller: deny before subject assertion.
- Unapproved subject asserter: ignore/reject asserted context and deny any operation requiring it.
- Readiness dependency fails: remove readiness and emit a low-cardinality reason.
- Telemetry backend fails: continue core work where safe; do not block a request indefinitely or log credentials to compensate.
- Shutdown deadline expires: record incomplete drains and exit; never claim graceful completion.

## Common misuse

- Treating workload identity as the represented user.
- Authorizing `Subject.Actor()` instead of the human subject plus delegation intersection.
- Putting a token, presigned URL, approval code, or private attribute into a log field.
- Making liveness depend on PostgreSQL or RabbitMQ and causing restart storms.
- Catching every dependency error with a retry loop, including invalid input or denials.
- Importing a package name that exists only in target diagrams.
