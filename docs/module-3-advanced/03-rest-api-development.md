---
title: REST API Development
sidebar_label: REST APIs
description: Resources, response contracts, errors, paging, OpenAPI, and idempotency.
---

# REST API development

## Resource design

Use stable plural resource paths and HTTP semantics:

| Operation | Method | Expected success |
|---|---|---:|
| List | GET /widgets | 200 |
| Get | `GET /widgets/{id}` | 200 |
| Create | POST /widgets | 201 |
| Replace | `PUT /widgets/{id}` | 200 or 204 |
| Partial update | `PATCH /widgets/{id}` | 200 or 204 |
| Delete | `DELETE /widgets/{id}` | 204 |
| Start async work | POST /jobs | 202 |

The public gateway path and strip-prefix behavior are part of the contract.

## Platform responses

platform/http.Handle renders the standard success envelope. Created and Accepted select status. HandleVoid returns 204. HandleRaw is reserved for OGC/GeoJSON, MCP, file streams, SSE, redirects, and 304.

Return classified platform/errors from application code. The HTTP adapter maps Validation to 400, Unauthorized to 401, Forbidden to 403, NotFound to 404, Conflict to 409, TooManyRequests to 429, and mandatory dependency failure to 503.

## Paging and sorting

Use paging.NewRequest and paging.Page[T]. Enforce deterministic order and map client sort names through a column allowlist. Never concatenate an unchecked identifier into SQL. Use cursor or keyset pagination for deep mutable collections.

## OpenAPI

Every route carries an operation ID that exists in the embedded OpenAPI document. Contract tests compare method, path, operation ID, auth metadata, media type, status, and schema. OpenAPI is executable review evidence, not a marketing document.

## Idempotency and concurrency

Webhook, callback, and event-driven handlers must deduplicate. A client-retryable create with external effects should accept an idempotency key. Updates to contested resources use a version, ETag, or transaction lock rather than last-write-wins by accident.

## Exercise

Design a bookmarks API. Declare paths, auth relationship, request/response types, errors, paging, idempotency behavior, and OpenAPI operation IDs before writing code. Then implement one typed handler.

## Check yourself

- When is HandleRaw appropriate?
- Why is gateway prefix behavior a client contract?
- How does a sort allowlist prevent injection?
- What makes a POST safe to retry?
