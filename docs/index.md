---
title: CDPG Go developer portal
slug: /
description: Canonical architecture, standards, tutorials, and Go curriculum for CDPG service developers.
---

# CDPG Go developer portal

Build and review CDPG Go services from an empty repository through gateway exposure, authorization, event integration, local testing, and GitOps deployment. The portal combines an evidence-backed platform guide with a progressive Go curriculum.

**Platform status: Partially implemented.** The platform is before its first production release. Every architecture and integration page distinguishes implemented source behavior from planned target capability.

## Choose your route

| Starting point | Recommended route |
|---|---|
| New to programming | Complete the Curriculum modules and checkpoints, then the new-service tutorial |
| Experienced engineer, new to Go | Module 0, Module 1 checkpoints, Modules 2–3 gaps, then Developer Guide |
| Experienced Go engineer | [Five-minute orientation](start/learning-paths.md), service standard, integrations, then tutorial |
| Service/platform/Data Plane/security/SRE reviewer | Select the role path in [Learning paths](start/learning-paths.md) |

## Learning outcomes

By the end you can:

- write idiomatic, race-free Go with explicit error and context flow;
- design a layered service whose business logic has no transport or driver imports;
- use dx-common-go/platform for bootstrap, HTTP, SQL, cache, events, identity, errors, paging, and health;
- secure a route through the gateway with OIDC, destination-bound workload identity, constrained subject propagation, and relationship authorization;
- administer grants with `dx-acl-go`, consume current OpenFGA checks, and prepare safely for planned OPA policy;
- enforce typed carried decisions in Data Plane services without accepting executable policy queries or calling the PDP per query;
- test units, adapters, contracts, security boundaries, and full-stack flows;
- deploy and operate a service through the platform's GitOps model;
- submit a focused, evidence-backed contribution.

## Start here

1. [CDPG platform orientation](architecture/platform-orientation.mdx)
2. [Current and target state](architecture/current-target.md)
3. [Go service architecture standard](standards/service-architecture.md)
4. [Shared Go platform modules](platform/index.md)
5. [New service quick start](new-service/quick-start.md)

## How the curriculum works

1. **Learn** the concept and its failure modes.
2. **Read** one current source example.
3. **Practice** in a small exercise.
4. **Check** your understanding without copying the sample.
5. **Verify** with tests, race checks, or a running stack.

## Authoritative references

The curriculum teaches. It does not duplicate reference documentation:

- [Data Exchange platform docs](https://datakaveri.github.io/cdpg-docs/) — architecture, fleet, security, deployment, operations.
- [dx-common-go docs](https://datakaveri.github.io/dx-common-go-docs/) — public SDK APIs, compatibility, examples.
- [Go documentation](https://go.dev/doc/) — language and toolchain.

New Go learners start with [Environment](module-0-setup/01-environment.md), then use the [roadmap](roadmap.md) to track milestones. Service developers start with the [learning paths](start/learning-paths.md).
