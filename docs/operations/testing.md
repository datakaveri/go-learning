---
title: Testing strategy
description: Required test layers for CDPG Go services, identity, authorization, data obligations, events, workers, and deployment.
---

# Testing strategy

**Status: Implemented standard; repository coverage is partial.** Tests prove contracts and failure behavior, not only line coverage.

## Test pyramid by boundary

| Layer | What to prove | Preferred technique |
|---|---|---|
| Domain | Invariants, transitions, time/ID edge cases | Table tests, property/fuzz tests, no mocks |
| Application | Orchestration, authz before effects, transaction/event intent | Small handwritten fakes at consumer interfaces |
| HTTP/gRPC | Binding, validation, identity gate, error/status/envelope, cancellation | Direct handler + router/local gRPC server |
| Repository/migration | Real SQL, constraints, locks, isolation, empty/upgrade schema | Disposable PostgreSQL/PostGIS |
| Integration adapters | Timeouts, credentials, provider contract, failure mapping | Fake server plus provider sandbox/contract fixtures |
| Events/workers | Envelope, outbox atomicity, duplicate/retry/DLQ/replay/reconnect/lease | Real broker/database plus controlled faults |
| Gateway/GitOps | Route posture, prefix/exact/strip, audience, network/exposure parity | Config conformance and end-to-end negative matrix |
| End to end | One owned workflow and recovery | Local orchestration, then deployment smoke |

## Authorization tests

- Authentication: absent/malformed/expired/not-yet-valid/wrong issuer/audience/algorithm/key/clock skew.
- Workload: wrong destination, disallowed caller, missing credential, subject assertion from non-asserter.
- OpenFGA: model compilation; user/group/org/owner/editor/viewer; tuple add/remove; deny by default; projection duplicate/reorder/reconciliation.
- Composite: **Planned**—OpenFGA and OPA allow/deny/error matrix, stable reasons, revisions, expiry, cache key and obligations.
- OPA: **Planned**—policy unit/data tests, schema validation, bundle verification/readiness/rollback. Do not create pretend passing tests before the runtime contract exists.
- Organisation isolation: swap resource, org, actor, and membership IDs across every mutation/read.
- Revocation: cached allow, in-flight action, subscription delivery, agent kill switch.
- Agent: human relation ∩ agent delegation, expiry/scope/resource, HITL single-use/race, semantic-firewall injection/tool-output handling.

## Data Plane obligation tests

**Status: Planned.** Fuzz artifact tamper/replay/binding and translators. Verify row/field/spatial/temporal/result/quota intersections never widen a query, all values are parameterized/allowlisted, pagination remains bounded, unsupported obligations deny, and no per-query synchronous PDP call exists.

## API and compatibility tests

- OpenAPI schema lint, generated artifact clean, request/response validation, operation ID ↔ route drift.
- Protobuf generation and breaking-change check; deleted fields reserved.
- Event fixture compatibility for current and supported older schema versions.
- Idempotency same-key/same-payload returns original result; same-key/different-payload conflicts.
- Webhook provider authenticity and duplicate event ID.
- SSE flush/cancel/resume and file streaming/range/body limits.

## Resilience tests

Inject database, broker, Redis, Keycloak/JWKS, OpenFGA, object store, and search outages. Verify timeout budgets, bounded retries/jitter, breaker behavior, fail-closed security, correct degraded mode, reconnect, readiness, backlog recovery, no duplicate effect, and correlated signals.

## Minimum suite for every service

- [ ] Domain table tests and fuzz target for parser/identifier/input boundary.
- [ ] Application allow/deny/error and no-effect-on-deny tests.
- [ ] HTTP handler/router and OpenAPI drift tests.
- [ ] Workload audience/caller/subject-asserter tests.
- [ ] Object authorization, default deny, organisation isolation, revocation tests.
- [ ] Repository and empty/upgrade migration tests.
- [ ] Transaction rollback and outbox atomicity tests.
- [ ] Consumer duplicate/retry/DLQ/reconnect/reconciliation tests when applicable.
- [ ] Worker crash/lease/cancel/shutdown tests when applicable.
- [ ] `go test -race ./...` and goroutine-leak test for lifecycle owners.
- [ ] Container health/SIGTERM and gateway smoke test.

## Commands and evidence

```bash
gofmt -w .
go build ./...
go vet ./...
go test ./...
go test -race ./...
govulncheck ./...
```

CI also runs contract generation/drift, dependency/license/secret/security scans, migration/integration suites, image scan/SBOM, and GitOps rendering/conformance. Archive machine-readable results with the immutable image revision.

