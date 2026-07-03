---
title: "Capstone — Build dx-bookmarks-go"
sidebar_label: Capstone Service
description: Build a complete DX-style service with dx-common-go, graded against the standards checklist.
---

# Capstone — Build `dx-bookmarks-go`

## Learning objectives

This is the integration test for *you*: one service, every platform pattern, no scaffolding provided. Completing it to rubric means you can be handed a real service ticket with confidence — yours and the team's.

## Prerequisites

All of Modules 0–4. Budget **12–16 hours** across 1–2 weeks. Work as if it were real: feature branch, incremental commits, the gate before every push.

## The brief

Build **`dx-bookmarks-go`**: users bookmark datasets and annotate them. Deliberately small domain — every ounce of difficulty should go into *platform correctness*, not business complexity.

**Functional requirements**

| Endpoint | Behavior |
|---|---|
| `POST /iudx/v2/bookmarks` | Create (datasetId + optional note). Authenticated. 409 on duplicate per user+dataset |
| `GET /iudx/v2/bookmarks` | Caller's bookmarks: paginated, sortable (whitelist), filter by datasetId |
| `GET /iudx/v2/bookmarks/{id}` | Fetch one — owner only (403 otherwise) |
| `PATCH /iudx/v2/bookmarks/{id}` | Update the note (pointer-as-optional semantics) |
| `DELETE /iudx/v2/bookmarks/{id}` | Soft delete; invisible to all reads thereafter |

**Platform requirements** — the actual test:

1. **Repo shape**: the canonical layout; `dx-common-go` via `replace`; README with API table, env vars, events.
2. **Boot contract**: `LoadService[T]` + `Validate()`; zap; Postgres hard (Fatal), RabbitMQ optional (Warn + no-op); idempotent schema ensure; `httpserver.Start()`.
3. **HTTP**: `StandardStack`; embedded OpenAPI spec with request validation + `/docs`; resolver middleware (HMAC + JWT) with owner checks from the context user; audit middleware on mutating routes.
4. **Contract**: envelopes via the response writers with your URN prefix (`urn:dx:bookmarks:`); taxonomy errors; `request.Builder`-style pagination.
5. **Persistence**: `BaseDAO[Bookmark]` where it fits, raw parameterized SQL where it doesn't; allowlisted sort keys; explicit soft-delete filtering; unique index behind the 409.
6. **Events**: `bookmark.created` / `bookmark.deleted` via transactional outbox → `ReliablePublisher` to a `bookmarks` topic exchange (durable + DLX + DLQ + TTL); supervised, context-cancelled dispatcher.
7. **Plus a consumer**: a small `cmd/indexer` (or second service) consuming those events into a `bookmark_counts` table — idempotent, bounded retry, poison messages to the DLQ.
8. **Observability**: `/healthz/live`, `/healthz/ready` (pool checker), `/metrics` with RED metrics.
9. **Tests**: table-driven handler tests per endpoint (denial cases included); both spec tests; one env-guarded integration test (repo + soft delete); `smoke.sh`.
10. **Ops**: multi-stage Dockerfile; compose entry with healthcheck; `.golangci.yaml`; the five-command gate green; drafted (not applied) gateway route + dx-gitops files.

## Suggested path

Build in the order the curriculum taught, verifying each layer before the next:

1. **Skeleton** (2h): repo shape, config, boot with Postgres only, health endpoints. *Verify: boots, ready flips with DB down.*
2. **Domain + persistence** (3h): schema, repository, integration test. *Verify: CRUD via psql-checked rows, soft delete filtered.*
3. **API** (3h): handlers, router, spec, validation, envelopes. *Verify: spec tests green; curl every endpoint + every error path.*
4. **Auth + audit** (2h): resolver, owner enforcement, audit middleware. *Verify: 401/403 paths; audit events visible in RMQ UI.*
5. **Events** (3h): outbox, dispatcher, topology, indexer consumer. *Verify: create → count increments exactly once even when you deliver the event twice; poison message lands in DLQ.*
6. **Hardening** (2h): metrics, Dockerfile, compose, smoke.sh, self-review against the standards, gate.

## Grading

Self-review first, then a reviewer (buddy/lead) grades **section by section against GO-SERVICE-STANDARDS.md** — the same way a real service would be reviewed. Expected outcome for a first attempt: mostly compliant with a handful of findings; **fixing the findings is part of the capstone**, because responding to review *is* the job.

Classic findings, from experience: missing soft-delete filter in one query; unsupervised dispatcher goroutine; consumer not idempotent under redelivery; secrets in the baked YAML; `fmt.Println` debugging survivors; sort key reaching SQL unvalidated; log-and-return double handling.

## Stretch goals (optional)

- Register it for real behind your local gateway (route config + policy) and call it with a `make dev-token` JWT through `:8000`.
- Advisory-lock the dispatcher and prove single-instance behavior with two replicas.
- A one-shot `cmd/purge` CronJob binary hard-deleting soft-deleted rows older than 30 days — idempotent, exit codes, YAML included.

When it's green and reviewed: [First Contribution](first-contribution).
