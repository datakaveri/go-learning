---
title: Standards Checklist
sidebar_label: Standards Checklist
description: Review checklist for Data Exchange Go service changes.
---

# Standards checklist

**Status: Superseded as a release checklist by the [service review scorecard](../standards/review-scorecard.md).** This lesson remains a compact exercise; use the scorecard and linked standards in an actual review.

## Architecture

- Business logic is in application/domain packages.
- Ports are narrow and defined by consumers.
- Vendor packages appear only in adapters or composition.
- Cross-cutting concerns use dx-common-go/platform.
- No cross-service database reads or Go implementation imports.

## Process and configuration

- bootstrap.Run declares process dependencies.
- Typed config embeds platform/config.Base and validates at boot.
- Required and degraded dependencies match tested behavior.
- Workers observe context; closers have one owner.
- Shutdown drains HTTP before workers and infrastructure.

## HTTP and contracts

- Routes are declarative with operation IDs.
- Public, optional, protected, role, and stream flags are explicit.
- OpenAPI and route tables are tested for drift.
- Platform errors and paging are used.
- HandleRaw has a protocol-specific justification.
- Gateway path and prefix behavior are tested.

## Security

- Tokens are validated at the gateway or trusted verifier.
- Client-supplied internal identity headers are removed.
- Destination-bound workload identity is verified before an allowlisted caller may assert subject context.
- Subject, actor/delegation, workload, and organisation are kept distinct.
- Planned OPA or carried-decision behavior is never presented as operational.
- Relationship, role, tenant, and ownership rules are tested.
- SQL values are parameterized and sort fields allowlisted.
- Secrets and sensitive payloads are absent from Git and logs.

## Data and events

- One service owns the schema and writes.
- Schema change is backward compatible with rolling deploys.
- Transactions receive the callback context.
- State plus event uses an outbox.
- Event payloads are versioned; consumers are idempotent.
- Retry, drop, dead-letter, and backlog behavior is observable.

## Tests and operations

- Unit, race, adapter, contract, and security tests match risk.
- Health distinguishes liveness, readiness, and degradation.
- Logs, metrics, traces, and audit use safe stable fields.
- Resources, concurrency, request size, and timeouts are bounded.
- Deployment, recovery, and rollback notes exist.
- Documentation updates the authoritative site.

## Exercise

Review a current PR or recent commit against every item. For each finding, cite the exact file/line, risk, and smallest fix. Separate blocking correctness/security findings from follow-up improvements.
