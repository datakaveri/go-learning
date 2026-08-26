---
title: Capstone Service
sidebar_label: Capstone Service
description: Build and verify a complete Data Exchange-style bookmarks service.
---

# Capstone service

**Status: Implemented learning assignment.** Follow the canonical [new-service tutorial](../new-service/quick-start.md) and [review scorecard](../standards/review-scorecard.md); this service is not part of the CDPG fleet.

Build dx-bookmarks-go: users bookmark catalogue items, list their bookmarks, annotate them, and remove them.

## Functional contract

- POST /v1/bookmarks creates a bookmark and is idempotent for user + item.
- GET /v1/bookmarks returns a deterministically paged list with allowed sorting.
- `GET /v1/bookmarks/{id}` returns one owned bookmark.
- `PATCH /v1/bookmarks/{id}` changes a bounded note.
- `DELETE /v1/bookmarks/{id}` removes it.
- bookmark.created and bookmark.deleted events are published through an outbox.

All business routes require identity and destination-bound workload verification. Ownership and organisation isolation are enforced even when an ID is guessed. Catalogue existence is checked through an authenticated internal client port with a bounded timeout. Current relationship checks use `dx-authz-go`; planned OPA behavior must not be invented.

## Architecture acceptance

- platform/bootstrap owns lifecycle.
- Typed config embeds platform/config.Base.
- Domain and application import no router, SQL driver, AMQP, or vendor client.
- platform/http declares routes and typed handlers.
- platform/errors and platform/paging define API behavior.
- platform/database/sql owns transactions and PostgreSQL adapter behavior.
- platform/events owns typed topics and outbox dispatch.
- platform/observability/health reports dependencies.

## Data

Create a service-owned PostgreSQL schema with embedded migrations. Run one controlled migration actor before serving replicas and keep business code independent of the operational boot mode. Add a unique user/item constraint, stable paging index, timestamps, and bounded note length.

## Verification

Required evidence:

~~~bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
~~~

Also provide:

- unit tests for invariants and ownership;
- repository/migration integration tests;
- route/OpenAPI contract tests;
- duplicate create and event redelivery tests;
- forged subject, wrong owner, invalid sort, oversized body, and cancellation tests;
- wrong workload audience/caller/subject asserter, organisation swap, authorization outage, and revocation tests;
- a smoke script through the gateway;
- a container SIGTERM drain test;
- rendered GitOps manifests.

## Operational deliverables

Document config, secrets, dependency policy, metrics, logs, traces, audit, dashboards, alerts, backup, rollout, rollback, and event replay. Use an immutable image digest.

## Review

Present:

1. architecture diagram;
2. one request trace;
3. one transaction/outbox trace;
4. threat model;
5. test evidence;
6. deployment and recovery plan;
7. tradeoffs and intentionally deferred work.

The capstone is complete when another engineer can operate and change it—not merely when curl returns 200.
