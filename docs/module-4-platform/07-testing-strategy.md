---
title: Platform Testing Strategy
sidebar_label: Testing Strategy
description: Fleet test layers, route and event contracts, security, full-stack demos, and evidence.
---

# Platform testing strategy

**Status: Superseded as the canonical checklist by [Testing strategy](../operations/testing.md).** Use this lesson to practice selecting test layers.

## Verification ladder

| Layer | Proves |
|---|---|
| Unit | Domain rules, error classification, handler translation |
| Adapter integration | Real SQL, Redis, AMQP, search, and object behavior |
| Contract | OpenAPI routes, gateway prefixes, event envelopes, client shapes |
| Security | Identity, authorization, injection, isolation, redaction |
| End to end | Cross-service identity, policy projection, routing, jobs |

No layer replaces another. A repository fake cannot prove SQL, and a full-stack demo cannot cover all domain edge cases.

## Platform tools

- Call platform/http typed handlers directly for application behavior.
- Use cache.NewMemory and events.NewMemory for portable semantics.
- Use dxtest/containers and service migrations for adapter tests.
- Verify route table and OpenAPI operation IDs.
- Test platform/database/sql transactions, not-found mapping, paging, and locks.
- Exercise AMQP reconnect, duplicate delivery, retry, drop, and dead letter.
- Test destination workload audience, caller and subject-asserter allowlists.
- Treat OPA and carried Data Plane tests as planned until their contracts exist; test safe denial in the interim.

## Fleet smoke tests

~~~bash
make dev-demo
make dev-demo-agent
~~~

The agent demo is required for agent-plane changes. OGC changes need PostGIS and worker coverage. Data-plane changes need real index naming and mapping verification.

## Race and leak tests

~~~bash
go test -race ./...
~~~

Run goroutine-leak checks for clients, consumers, dispatchers, refreshers, and shutdown. Cancel tests at each lifecycle phase.

## Evidence-based review

Record exact command, source ref, environment, result, and skipped tests. A green command that silently skipped unavailable containers is not integration evidence.

## Exercise

Write a test plan for policy create → outbox → RabbitMQ → dx-authz-go → OpenFGA → gateway allow. Include duplicate event, delayed projection, broker outage, forged identity, and revocation.

## Check yourself

- What proves a route is protected?
- How do you know an integration test really ran?
- Why is the race detector not a unit-test substitute?
- Which failure belongs only in full-stack testing?
