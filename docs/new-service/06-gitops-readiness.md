---
title: 6. GitOps and production-readiness review
description: Add deployment values, identity, network, health, telemetry, rollback, and final evidence.
---

# 6. GitOps and production-readiness review

**Status: GitOps framework implemented; service onboarding and operational proof are service-specific.** A merged values file is not proof that a capability has run in production.

![Kubernetes and GitOps deployment](/img/architecture/kubernetes-gitops.svg)

## Step 1 — service values

Add the service to the intended `dx-gitops` ApplicationSet/environment inventory and shared chart values. The current chart supports deployment/service, ExternalSecret, NetworkPolicy, migration job, PodDisruptionBudget, HorizontalPodAutoscaler, ServiceAccount, probes, resources, security context, read-only root with `/tmp`, and topology spread.

```yaml title="Shape verified against current service values; names and resources are illustrative"
nameOverride: dx-example-go
image:
  repository: <approved-registry>/dx-example-go
  tag: <immutable-version-or-digest>
containerPort: 8080
migrations:
  enabled: true
networkPolicy:
  enabled: true
externalSecrets:
  enabled: true
resources:
  requests: {cpu: 100m, memory: 128Mi}
  limits: {memory: 256Mi}
```

Copy no resource number without load evidence. Pin the image and use the same artifact for migration and serving.

## Step 2 — identity and network

- Provision a unique Keycloak workload client secret through External Secrets.
- Grant only required destination audience scopes.
- Configure destination caller and subject-asserter allowlists.
- Permit NetworkPolicy edges only for gateway, required callers, owned data stores, broker, telemetry, and approved external services.
- Deny unnecessary cross-service/database access.

Provision identity before enabling the route. A service must not fall back to an unverified internal call when token issuance fails.

## Step 3 — probes, metrics, and resources

- Liveness checks process health only.
- Readiness checks whether the advertised contract can be served.
- Migration job completes before rollout.
- Metrics scraping and alerts cover requests, authn/authz denials/errors, dependency state, pool/cache, outbox/queue/DLQ/projection lag, workers, and shutdown.
- Requests/limits come from profiling and concurrency/load tests.
- Termination grace exceeds measured drain duration.

## Step 4 — rollout and rollback

Use an expand-compatible schema, deploy, observe, and contract later. Define rollback behavior if schema, event, or protobuf versions move forward. Preserve idempotency records and outbox data across rollback. A failed policy/runtime readiness gate must keep the service unready and route disabled.

Verify:

```text
rendered values -> config-check -> migration job -> workload identity
-> NetworkPolicy -> readiness -> route enable -> smoke/security tests
-> SLO/alerts -> rollback drill
```

## Step 5 — final checklist

- [ ] Service charter, ownership, API, event, datastore, and authorization mapping approved.
- [ ] Required project layout and dependency direction pass static review.
- [ ] Shared platform packages used; no copied infrastructure wrapper without ADR.
- [ ] Typed configuration and all security settings fail at startup when invalid.
- [ ] Migrations are embedded, tested, and run by one actor.
- [ ] State/outbox are atomic; consumers idempotent; retry/DLQ/replay/reconciliation tested.
- [ ] Gateway route, workload audience, callers, subject asserters, and limits verified.
- [ ] Default deny, organisation isolation, revocation, and failure behavior verified.
- [ ] OpenAPI/protobuf/event compatibility and code generation clean.
- [ ] Unit, integration, contract, security, race, shutdown, and smoke suites pass.
- [ ] Container non-root/read-only, scanned, health-tested, and signal-tested.
- [ ] GitOps values, secrets, NetworkPolicy, probes, metrics, resources, and rollback reviewed.
- [ ] Runbook, ownership, backup/restore, DLQ/replay, reconciliation, SLOs, and alerts exist.
- [ ] Planned OPA/carried-decision/agent features are not claimed as operational.

Apply the [service review scorecard](../standards/review-scorecard.md). A passing score makes the service a candidate for release review; production readiness still requires deployment and operational evidence.

## Reference implementation tree

```text
dx-example-go/
├── cmd/server/main.go                 # bootstrap.Spec + wire only
├── internal/config/config.go          # config.Base + validation
├── internal/domain/widget.go          # invariants/value types
├── internal/application/widget.go     # ports/use cases/transactions/authz
├── internal/transport/http/           # typed handlers/routes/DTO mapping
├── internal/transport/grpc/           # generated-interface adapter if needed
├── internal/repository/postgres/      # parameterized queries and row mapping
├── internal/integration/authz/        # current dx-authz-go adapter
├── internal/worker/                   # consumers/dispatcher/reconciliation
├── db/migrations/                     # schema and outbox
├── api/openapi/openapi.yaml           # public contract
├── api/proto/example/v1/              # internal contract if needed
├── configs/config.yaml                # non-secret local defaults
├── tests/                             # integration/contract/security/smoke
├── Dockerfile
├── go.mod
└── README.md
```

The repository may omit adapters it does not need. It may not collapse business logic into transport, repository, or `main` for convenience.
