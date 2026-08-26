---
title: New service quick start
description: Build and test the teaching dx-example-go service in five minutes.
---

# New service quick start

**Status: Implemented teaching workflow.** `dx-example-go` is an educational example in this repository, not a platform service or reserved name. Its unit-level path compiles and runs; gateway, identity, PostgreSQL, RabbitMQ, and GitOps steps require the integration tutorial and approved configuration.

![New-service delivery lifecycle](/img/architecture/new-service-lifecycle.svg)

## Five minutes

```bash
cd examples/dx-example-go
go test ./...
go test -race ./...
go vet ./...
```

Expected: all packages build, unit/handler tests pass, and the race detector reports no race. The example’s local `replace` points to the adjacent `dx-common-go` clone in the CDPG orchestration workspace; a released service must pin a reviewed module version and remove the filesystem replacement.

## What you are proving

- domain construction rejects invalid state;
- application logic depends on consumer-owned repository and authorization interfaces;
- a denial stops before persistence;
- a typed platform HTTP handler binds a request and maps a result;
- implementation details remain behind adapters.

You are not yet proving Keycloak, workload identity, OpenFGA, PostgreSQL, RabbitMQ, the gateway, or deployment. Those gates arrive in the progressive tutorial.

## Tutorial path

1. [Choose the boundary and scaffold](./01-boundary-scaffold.md).
2. [Bootstrap, configuration, HTTP, and gRPC](./02-bootstrap-transports.md).
3. [Persistence, cache, events, and workers](./03-persistence-events.md).
4. [Identity, authorization, ACL, and gateway](./04-security-gateway.md).
5. [Tests, container, and local integration](./05-tests-local.md).
6. [GitOps, operations, and production-readiness review](./06-gitops-readiness.md).

Each stage has a stop condition. Do not add an externally reachable route until the negative authorization, workload, organisation-isolation, and failure-path tests pass.

