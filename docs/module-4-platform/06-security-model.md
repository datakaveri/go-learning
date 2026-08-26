---
title: Security Model
sidebar_label: Security Model
description: Threat-model the gateway, services, data stores, events, callbacks, and agents.
---

# Security model

**Platform status: Partially implemented.** Read the canonical [identity](../integrations/identity.md), [authorization](../integrations/authorization.mdx), [Data Plane authorization](../integrations/data-plane-authorization.mdx), and [Agentic Plane](../architecture/agentic-plane.md) pages before this exercise.

Each arrow crosses a boundary with authentication, authorization, integrity, replay, availability, and observability requirements.

## Gateway

The gateway validates OIDC claims, applies route policy, strips untrusted internal headers, checks configured relationships, obtains a destination-bound workload credential, constrains subject propagation, rate limits, and proxies to allowlisted upstreams. Protected dependency failure fails closed.

## Service

The service verifies the calling workload before subject context and applies operation-specific relationship, organisation, ownership, state, obligation, and input rules. Direct local ports do not make an endpoint trusted; internal callbacks and admin routes need explicit authentication.

## Stores and events

Database roles receive only owned schema privileges. SQL identifiers are allowlisted and values parameterized. Object keys and URLs resist traversal and SSRF. Event producers confirm publication; consumers deduplicate, bound retries, and inspect versions.

## Agent plane

Agent authority includes separate agent identity, original human, scoped expiring grant, active lifecycle, relationship decisions, semantic tool policy, token custody, and human approval for configured high-impact action. Kill-switch behavior is an operational control that must be drilled.

## Threat-model exercise

For multipart file completion, list:

- assets and actors;
- entry points and trust boundaries;
- spoofing, tampering, repudiation, disclosure, denial, and privilege threats;
- controls and tests;
- residual risk and monitoring.

Then repeat for an agent MCP tool that deletes a resource.

## Check yourself

- Why is hiding an endpoint from the gateway insufficient?
- What prevents an arbitrary caller from asserting another subject?
- How do event consumers defend against replay?
- Why must delegated tokens stay outside model-visible data?
