---
title: Local development and GitOps deployment
description: Workspace placement, infrastructure, ports, Compose, gateway, containers, Kubernetes, verification, and rollback.
---

# Local development and GitOps deployment

**Status: Local orchestration and shared GitOps chart implemented; service coverage varies.** Agentic Plane services are locally integrated but absent from GitOps. `dx-dataplane-rs-go` remains internally exposed while carried authorization is incomplete.

## Workspace

Service repositories are cloned inside the CDPG orchestration checkout beside `dx-common-go`; the orchestrator ignores them in Git while Compose uses them as build contexts. Use the repository inventory/clone script. A local `go.work` or relative `replace` supports coordinated development, but release modules pin reviewed versions.

## Ports

`claude-docs/PORTS.md` is the only port authority. Confirm both host and container ports before editing Compose, gateway, OpenAPI servers, health checks, or GitOps values. Service-to-service calls use network DNS and container ports, not host mappings.

## Local flow

1. Clone/update required service and `dx-common-go` repositories.
2. Start only required PostgreSQL/PostGIS, Redis, RabbitMQ, Elasticsearch, MinIO/S3, Keycloak, and OpenFGA dependencies.
3. Create the service-owned database/schema.
4. Validate config with `DX_BOOT_MODE=config-check`.
5. Apply migrations with `DX_BOOT_MODE=migrate-only`.
6. Provision workload client/audience/caller/subject-asserter configuration.
7. Start service and verify liveness/readiness/metrics.
8. Add the gateway route disabled; run route/auth negative tests.
9. Enable locally and run allow/deny/revoke/recovery smoke workflows.

Never place a real secret in Compose YAML or documentation. Development identities belong in the local bootstrap mechanism and must not be reused outside it.

## Compose registration

Declare build/image, internal networks, container port, healthcheck, environment and secret references, dependency URLs, database, and migration command. Publish a host port only when direct debugging needs it. Configure restart behavior without hiding a permanent config error.

## Container requirements

- immutable, minimal, non-root image;
- read-only root filesystem and explicit writable temp volume;
- no embedded credentials/config secrets;
- one binary for serving/config-check/migrations;
- architecture/version metadata and SBOM;
- liveness/readiness and measured graceful shutdown;
- bounded CPU/memory behavior and body/connection limits.

## GitOps topology

The Argo CD ApplicationSets render environment values through a shared Helm chart. The chart can create Deployment, Service, migration job, ExternalSecret, NetworkPolicy, ServiceAccount, probes, resources, PDB/HPA, security context, and topology spread.

Onboarding sequence:

1. merge service image and contracts;
2. add values to intended environments with immutable image reference;
3. configure External Secrets and workload identity/audiences;
4. restrict NetworkPolicy;
5. render and run config-check;
6. run PreSync migration job from the same image;
7. deploy unexposed and verify readiness/metrics;
8. add/enable gateway route after conformance and security smoke tests;
9. observe SLOs/backlogs/projections and exercise rollback.

## Deployment verification

- desired/actual revision and image digest match;
- migration job completed once and schema readiness is correct;
- pods run non-root/read-only and have requests/limits;
- workload token has correct destination audience; wrong audience/caller fails;
- NetworkPolicy allows only intended flows;
- gateway route matches source/GitOps exposure inventory;
- allowed/denied/revoked workflows and audit correlation pass;
- outbox/consumer/projection catch up after dependency restart;
- SIGTERM drains within grace period;
- dashboards/alerts/runbooks and backup restore evidence exist.

## Rollback and recovery

Rollback the immutable image/config while preserving compatible data. Expand-compatible schema changes support mixed-version rolling updates; destructive contraction follows after the rollback window. Do not delete outbox/idempotency records during rollback. If policy projection or identity is unhealthy, disable the route/fail readiness rather than bypassing enforcement. Reconcile projections after recovery and document user/security impact.

See [new-service GitOps tutorial](../new-service/06-gitops-readiness.md) for the full checklist.
