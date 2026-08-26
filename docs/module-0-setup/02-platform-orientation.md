---
title: Platform Orientation
sidebar_label: Platform Orientation
description: Build a mental model of the Go platform planes, request path, ownership, trust, and status.
---

# Platform orientation

## The system you are learning

**Platform status: Partially implemented.** Read the canonical [Go platform orientation](../architecture/platform-orientation.mdx) and [current/target matrix](../architecture/current-target.md); those pages own component and status detail.

The public gateway validates user identity, applies route posture, performs configured relationship checks, obtains a destination-bound workload credential, constrains represented-subject propagation, and proxies. `dx-authz-go` evaluates current OpenFGA relationships; `dx-acl-go` owns policy/access-request administration; `dx-user-go` owns user and organisation state. Contextual OPA policy is planned.

The Control Plane owns governance and entitlements. The Data Plane owns files, NGSI-LD, OGC, search/spatial/temporal execution, and safe filtering. The Agentic Plane owns agent registration, sessions, MCP tool execution, human approval, and kill/revocation controls.

## Data boundaries

Each service owns its database/schema and write path. A controlled migration actor applies embedded versioned schema changes before serving replicas require them. Application behavior is selected through validated configuration, not by a named-environment branch.

RabbitMQ carries policy, membership, audit, notification, subscription, job, and agent activity. OpenFGA stores relationship tuples; Keycloak stores identity; Elasticsearch powers catalogue and NGSI-LD queries; PostGIS powers OGC APIs; object storage contains files and coverages.

## Shared Go SDK

dx-common-go/platform provides the reusable process and adapter contracts. A service should contain business logic, domain adapters, and composition—not another implementation of configuration, response problems, paging, transactions, cache, events, or health.

## Delivery status

Agentic Plane repositories and local Compose paths exist, but GitOps registration is missing. Carried Data Plane decision enforcement and OPA are planned. Federation is deferred. Do not infer operational readiness from source presence.

## Source exercise

1. Open dx-gateway-go/configs and find the route for /iudx/v2/cat.
2. Open dx-catalogue-go/internal/api/router.go and find the matching route table.
3. Find bootstrap.Run in its cmd/server/main.go.
4. Identify its search, authorization, and event dependencies.
5. Draw the read path and one write-side event path.

## Checkpoint

- Where does an external bearer token terminate?
- Which service owns policy state, and which service answers decisions?
- Why must a service not query another service's database?
- When should code use dx-common-go/platform instead of a service-local helper?
- Which capabilities are implemented, partial, planned, or deferred?
