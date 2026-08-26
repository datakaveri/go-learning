---
title: Project Structure
sidebar_label: Project Structure
description: Five-layer service layout, dependency direction, and ownership.
---

# Project structure

## The five layers

| Layer | Contains | May depend on |
|---|---|---|
| L0 kernel | Domain values, platform errors, identity, paging | Standard library and L0 |
| L1 ports | Repository and outbound-client interfaces | L0 |
| L2 application | Use cases and workflow orchestration | L0 and L1 |
| L3 adapters | HTTP, SQL, events, search, object storage | Inward layers and platform adapters |
| L4 composition | cmd/server | Every layer for wiring only |

Dependencies point inward. Domain and application code do not import router, pgx, Redis, AMQP, Elasticsearch, S3, or Kubernetes types.

## Standard tree

~~~text
cmd/server/main.go
internal/
  domain/
  service/
  api/
  repository/postgres/
  client/
  config/
db/
  migrations/
  embed.go
configs/config.yaml
openapi/
Dockerfile
~~~

Put an interface next to the application component that consumes it. Put its concrete adapter under repository or client. Keep generated code in a named package and never edit it manually.

## Ownership questions

Before creating a package, ask:

- Which domain owns the behavior and state?
- Is this service-specific or shared runtime policy?
- Which layer should know the vendor?
- Who constructs and closes the resource?
- What is the smallest API another package needs?

Cross-cutting transport, configuration, errors, paging, SQL, cache, events, identity, health, and lifecycle belong in dx-common-go/platform. Domain behavior remains in the service.

## Internal and public

Use internal for service implementation. Public APIs are HTTP/OpenAPI, events, and deliberate Go modules—not exported service packages. This preserves freedom to refactor inside a service.

## Exercise

Take one service and draw every import edge for a request. Flag an outward import from application code or a package called utils. Propose the narrow port and composition change that removes it.

## Check yourself

- Why is composition the only layer allowed to know all implementations?
- Where does an outbound catalogue-client interface live?
- When should code move to dx-common-go?
- What boundary prevents cross-service Go imports?
