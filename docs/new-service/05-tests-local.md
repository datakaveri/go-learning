---
title: 5. Tests, container, and local integration
description: Verify the service and integrate it safely into the orchestration workspace.
---

# 5. Tests, container, and local integration

**Status: Implemented process; individual service conformance varies.** Local success is an integration gate, not production evidence.

## Step 1 — minimum test suite

Add domain, application, handler, OpenAPI, repository, migration, event/outbox, consumer/idempotency, workload identity, subject asserter, OpenFGA model, organisation isolation, default deny, revocation, shutdown, race, and smoke tests. Add OPA/obligation tests only when those planned contracts are approved and implemented; until then test that unsupported context is denied.

Verification:

```bash
gofmt -w .
go build ./...
go vet ./...
go test ./...
go test -race ./...
govulncheck ./...
```

Expected: clean output and no skipped security suite without a recorded reason.

## Step 2 — container image

Use a multi-stage build, pinned Go/base images, a minimal runtime, non-root user, read-only root filesystem, no compiler/shell where unnecessary, and one binary that supports serve/config-check/migrate-only.

```dockerfile title="Illustrative; pin image digests in the service repository"
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/service /service
USER nonroot:nonroot
ENTRYPOINT ["/service"]
```

Common mistake: floating image tags, secrets in build args/layers, root execution, writable root, or a migration image built from different code.

Verification: image scan/SBOM, run as non-root/read-only, liveness/readiness, SIGTERM drain. Expected: no critical known vulnerability exception without owner/date.

## Step 3 — place the repository

Clone the service inside the CDPG orchestration workspace beside `dx-common-go`, because Compose build contexts and local module replacements expect that structure. Update the orchestrator’s repository clone inventory; do not commit the nested service repository into the orchestrator.

## Step 4 — infrastructure and database

Declare only required PostgreSQL/PostGIS, Redis, RabbitMQ, Elasticsearch, S3-compatible, Keycloak, or OpenFGA dependencies. Add the service database creation step and a migration-only command. Use `claude-docs/PORTS.md` to allocate ports; never copy one from this tutorial.

Verification: start from empty volumes, create database, migrate, run readiness; restart each dependency. Expected: correct dependency failure behavior and recovery without manual hidden state.

## Step 5 — Compose and gateway

Add the image/build, internal network, environment, secret references, health check, depends-on only for boot convenience, and no unnecessary published host port. Add the gateway route and workload identity provisioning. Keep `enabled: false` until the service and authorization projection are ready.

## Step 6 — smoke workflow

Test through the gateway, not only the host-bound service:

1. obtain a development user token without recording credentials in docs/logs;
2. create or select a resource/grant through the owning Control Plane service;
3. wait for relationship projection;
4. call the allowed operation;
5. call with a denied user and swapped organisation/resource;
6. revoke and verify denial;
7. inspect audit, traces, outbox, queue, projection, and service metrics;
8. restart service/broker during work and confirm recovery/idempotency.

Expected: allowed path succeeds, every negative path fails at the intended enforcement point, and correlated signals explain why.

## Local topology

![Local development topology](/img/architecture/local-development.svg)

Use the ports page as the live registry. Service-to-service traffic uses Compose DNS/container ports; host ports are developer conveniences, not internal discovery.

## Checkpoint

Attach command output, contract reports, security matrix, smoke evidence, and unresolved gaps to the review. Next: [GitOps and readiness](./06-gitops-readiness.md).

