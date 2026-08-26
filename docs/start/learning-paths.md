---
title: Learning paths
description: Role-based routes through the CDPG Go developer portal.
---

# Choose a learning path

**Status: Implemented documentation.** These paths organize the current Go curriculum and developer guide; they do not indicate that every target capability is operational.

## Five-minute orientation

If you have five minutes:

1. Read [Platform orientation](../architecture/platform-orientation.mdx).
2. Scan the [current and target state](../architecture/current-target.md).
3. Learn the [service dependency rule](../standards/service-architecture.md): transports depend on application ports; domain code depends on neither transport nor infrastructure.
4. Read the [gateway](../integrations/gateway.md) and [authorization](../integrations/authorization.mdx) responsibility boundaries.
5. Run the [new-service quick start](../new-service/quick-start.md).

## Beginner: learn Go, then CDPG

Start with [Module 0](../module-0-setup/01-environment.md), complete Modules 1–3, then follow the new-service tutorial. The checkpoints are designed to make concurrency, context, errors, and tests habitual before infrastructure is introduced.

## Intermediate: experienced Go developer joining CDPG

Read the platform orientation, service architecture, shared platform guide, gateway, identity, authorization, persistence, and messaging pages. Then build `dx-example-go` through the tutorial and use the [service review scorecard](../standards/review-scorecard.md).

## Advanced paths

| Audience | Read in this order | Outcome |
|---|---|---|
| Service owner | Architecture → service standard → new-service tutorial → operations → scorecard | Own a bounded context and its data, API, events, deployment, and on-call signals. |
| Platform/framework developer | Shared modules → standards → API transports → persistence → messaging → testing | Change common infrastructure without importing business semantics or breaking consumers. |
| Data Plane developer | Platform orientation → authorization → Data Plane authorization → persistence → observability | Translate typed obligations through safe query builders and fail closed on unsupported decisions. |
| Authorization engineer | Identity → ACL/PAP → composite authorization → workflows → testing | Maintain clean boundaries between authentication, policy administration, relationship evaluation, contextual policy, and enforcement. |
| SRE/platform operator | Current/target → local/deployment → observability → resilience/testing | Configure dependencies, identity, routes, probes, policy readiness, rollout, and recovery. |
| Technical reviewer | Architecture → standards → security integrations → scorecard | Review evidence rather than relying on package names or unsupported readiness claims. |

## Reference versus tutorial

The developer guide explains CDPG contracts and gives verified examples. The curriculum teaches Go progressively. The [platform architecture site](https://datakaveri.github.io/cdpg-docs/) remains the broader architecture authority; the [shared SDK site](https://datakaveri.github.io/dx-common-go-docs/) is the generated API reference.

