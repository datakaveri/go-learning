---
title: Service Anatomy
sidebar_label: Service Anatomy
description: Read one Go service from composition through routes, use cases, persistence, events, and operations.
---

# Service anatomy

Use dx-acl-go as the worked example because it combines protected HTTP, PostgreSQL transactions, outbox publication, authorization projection, scheduled work, and external clients.

## Composition

cmd/server/main.go calls bootstrap.Run with:

- typed config options;
- embedded migrations and required PostgreSQL;
- Wire for repository, service, handler, event bus, dispatcher, clients, workers, and router.

Bootstrap owns server and process lifecycle. Wire owns the service-specific graph.

## HTTP adapter

internal/api declares route tables for policies, access requests, and internal delegations. Typed handlers receive httpx.Actor, bind path/query/body fields, call application methods, and return values or errors.

## Application

internal/service owns policy and access-request invariants. It depends on narrow repository and client ports, not pgx or AMQP. Authorization semantics are expressed as roles, ownership, and relationships rather than raw token claims.

## Persistence and events

internal/repository/postgres uses platform/database/sql. A state change and policy event write share the ambient transaction. The outbox dispatcher publishes to platform/events/amqp; dx-authz-go later projects the relationship into OpenFGA.

## Operational surface

Health and metrics are outside the business base. PostgreSQL is a required readiness dependency. Recoverable consumers use supervised background behavior. Shutdown drains HTTP before workers and clients close.

## Reading exercise

Trace access-request approval:

1. route and operation ID;
2. request binder and actor;
3. service invariant;
4. transaction and locks;
5. policy and event rows;
6. outbox dispatcher and topic;
7. authorization projection;
8. audit/notification side effects;
9. metrics and failure logs.

Write file and symbol references. If a step is inferred rather than source proven, label it as an inference and verify it.

## Check yourself

- Where does business validation live?
- How does the repository join the application transaction?
- What makes policy projection eventually consistent?
- Which worker failure should stop the service?
