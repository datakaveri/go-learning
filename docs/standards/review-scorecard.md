---
title: Service review scorecard
description: Evidence-based architecture, security, quality, and operations review for CDPG Go services.
---

# Service review scorecard

**Status: Implemented documentation; CI automation is partially implemented.** A service is not ready because it has a Dockerfile. Review every row with a source, test, configuration, or operational link.

Score each item `2` (verified), `1` (partial with owner/date), or `0` (missing). Any release blocker below overrides the total.

| Area | Evidence | Score |
|---|---|---:|
| Bounded context | One responsibility; no duplicated owner; service/data/event inventory updated | /2 |
| Layering | Domain has no transport/driver imports; use cases own orchestration; composition root only wires | /2 |
| Contracts | OpenAPI and any protobuf generated/validated; route and compatibility tests pass | /2 |
| Configuration | Typed, validated, secret-free defaults; rendered `config-check` passes | /2 |
| Persistence | Service-owned store; migrations from empty and previous; one migration actor; recovery documented | /2 |
| Transactions/events | State and outbox commit atomically; stable event/version; idempotent consumer and reconciliation | /2 |
| Authentication | External OIDC and internal workload audience/caller tests; optional-auth invalid-token test | /2 |
| Authorization | Object/action mapping, default deny, organisation isolation, revocation, failure/cache behavior tested | /2 |
| Gateway | Explicit auth mode, destination, path/strip/exact behavior, limits/timeouts/streaming, route tests | /2 |
| Observability | Structured logs, metrics, traces, decision/audit correlation, dependency and backlog signals | /2 |
| Resilience | Bounded retries, timeouts, idempotency, leases, cancellation, graceful shutdown | /2 |
| Testing | Unit, application, handler, repository, migration, contract, security, race, smoke suites | /2 |
| Container | Minimal non-root image, read-only root, health check, signal behavior, scan/SBOM | /2 |
| GitOps | Values, secrets, NetworkPolicy, resources, probes, metrics, rollout/rollback, image pin | /2 |
| Operations | SLO/signals, alerts, runbook, ownership, backup/restore, DLQ/replay/reconciliation | /2 |

## Release blockers

- A reachable route cannot enforce its required authorization.
- Workload verification is missing or configured `off` for an authenticated internal surface.
- A subject assertion is accepted from a caller not on the allowlist.
- A service performs unsafe datastore filtering or accepts executable query text from policy.
- State is committed before an event whose delivery is required for correctness, without an outbox/reconciliation design.
- A consumer is non-idempotent, unbounded, cannot reconnect, or has no dead-letter/replay path.
- Required secrets are in source, examples, images, or logs.
- Schema changes run from every application replica rather than one controlled migration actor.
- Planned OPA or carried-decision behavior is represented as operational.
- Required build, vet, test, race, contract, security, or migration checks fail.

## Decision

| Result | Meaning |
|---|---|
| 27–30 and no blockers | Candidate for release review; still requires deployment and operational evidence. |
| 21–26 and no blockers | Integrable with an owned remediation plan; not release-ready. |
| 0–20 or any blocker | Stop and remediate before exposure. |

Attach this scorecard to the service ADR/README and update it when a contract, datastore, authorization profile, or deployment dependency changes.
