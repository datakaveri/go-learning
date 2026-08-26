---
title: dx-common-go Tour
sidebar_label: dx-common-go Tour
description: Current platform packages, foundation packages, adapter boundaries, and selection rules.
---

# dx-common-go tour

Use the [SDK reference](https://datakaveri.github.io/dx-common-go-docs/) beside this lesson.

**Status: Partially implemented and actively adopted.** Use the canonical [module inventory and responsibility matrix](../platform/index.md). It is source-verified and includes bootstrap/config, HTTP/gRPC, errors/paging, SQL, cache, events/AMQP, idempotency, leases, executor, identity/workload, and health.

Search, S3-compatible storage, user JWT/JWKS authentication, authorization clients, metrics, tracing, audit, OpenAPI, resilience, and test support currently remain focused top-level packages. Do not invent `platform/search`, `platform/storage`, `platform/authz`, or `platform/test`.

## Adapter escape hatches

- platform/database/sql/pgx.Pool for existing adapters needing pgxpool;
- dxsql.SQL, SQLOne, or sqlc for queries outside Repo;
- platform/http.HandleRaw for protocol-owned responses;
- platform/cache/redis and platform/events/amqp for concrete infrastructure.

An escape hatch stays at L3/L4 and is easy to grep.

## Source exercise

1. Read platform/bootstrap.Spec and App.
2. Read platform/http.Handler and Route.
3. Read platform/database/sql.Manager.Do.
4. Read `platform/events.Topic`, `Subscribe`, `Outbox`, and `Dispatcher`.
5. Find each call site in one service.
6. Record any vendor type that escapes into application code.

## Check yourself

- Which package owns safe error rendering?
- Why are adapter subpackages named?
- When should a foundation package be used directly?
- What evidence justifies expanding a platform interface?
