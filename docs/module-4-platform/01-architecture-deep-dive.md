---
title: Architecture Deep Dive
sidebar_label: Architecture Deep Dive
description: Trace the complete Go platform across edge, governance, application, data, and operational planes.
---

# Architecture deep dive

**Platform status: Partially implemented.** Start with the canonical [platform orientation](../architecture/platform-orientation.mdx), then read the [Control Plane](../architecture/control-plane.md), [Data Plane](../architecture/data-plane.md), and [Agentic Plane](../architecture/agentic-plane.md) pages.

The gateway is a public Policy Enforcement Point. It validates external identity, applies route posture, requests configured relationship decisions, obtains a destination-bound workload credential, constrains subject propagation, and routes. `dx-acl-go` administers policy; `dx-authz-go` decides with OpenFGA today and is the planned composite boundary for OPA; application services enforce business state and obligations.

## Synchronous and asynchronous boundaries

A protected read is synchronous through gateway → authorization → service → store. A policy write commits in dx-acl-go, publishes through its outbox, crosses RabbitMQ, and is projected by dx-authz-go into OpenFGA. Audit and notification are also asynchronous.

The client-visible state must name these boundaries. “Write accepted” and “projection available” can occur at different times.

## Data ownership

Each service owns its schema and write path. PostgreSQL/PostGIS, Elasticsearch, Redis, object storage, OpenFGA, Keycloak, and RabbitMQ are infrastructure systems, not domains. No application composes business state with a cross-service SQL join.

## Current and target state

Do not infer readiness from a repository’s existence. Review the [status matrix and architecture gaps](../architecture/current-target.md). Agentic Plane local integration exists without GitOps registration. Federation is deferred. Data Plane carried decisions and contextual OPA policy are planned.

## Architecture exercise

Choose one:

- catalogue item create;
- file multipart completion;
- marketplace purchase;
- OGC process job;
- agent tool requiring approval.

Trace external route, identity, authorization object, owning service, stores, transaction, event, consumers, telemetry, and failure states. Cite source files for every arrow.

## Check yourself

- Which service is authoritative for policy state?
- Which state crosses RabbitMQ before it is visible?
- Why is OpenFGA not the policy administration database?
- How do optional profiles preserve the core boundary?
