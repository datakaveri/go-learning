---
title: Architecture source register
description: Canonical decisions, implementation evidence, precedence, and confirmation gaps for this portal.
---

# Architecture source register

**Status: Implemented documentation control. Last verified: 10 August 2026.** Use source plus approved decisions to update a claim. Do not change a status from a slide, issue title, or repository name alone.

## Precedence

1. The current request and reviewed `cdpg-docs` target architecture define the documentation outcome.
2. Approved Architecture Decision Records (ADRs) define responsibility and contract decisions.
3. Current source/tests establish implemented behavior.
4. GitOps and orchestration establish deployment/local integration.
5. Roadmaps/handover establish status and open work, but do not turn planned behavior into an API.

When sources conflict, record the gap. Do not silently select the most convenient version.

## Canonical design sources

| Subject | Source |
|---|---|
| Full platform architecture | [CDPG architecture site](https://datakaveri.github.io/cdpg-docs/) |
| Shared SDK architecture | [PLATFORM-ARCHITECTURE.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/PLATFORM-ARCHITECTURE.md) |
| Authorization target | [AUTHORIZATION-PLATFORM-TARGET-DESIGN.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/AUTHORIZATION-PLATFORM-TARGET-DESIGN.md) |
| Go service rules | [GO-SERVICE-STANDARDS.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/GO-SERVICE-STANDARDS.md) and [ENGINEERING-PLAYBOOK.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/ENGINEERING-PLAYBOOK.md) |
| API, persistence, messaging | [API-STANDARDS.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/API-STANDARDS.md), [DATABASE.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/DATABASE.md), [MESSAGING-AND-WORKERS.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/MESSAGING-AND-WORKERS.md) |
| Security | [SECURITY-REVIEW.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/SECURITY-REVIEW.md) |
| Agentic Plane | [AGENT-PLANE.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/AGENT-PLANE.md) |
| Deployment/testing/ports | [GITOPS.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/GITOPS.md), [TESTING.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/TESTING.md), [PORTS.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/PORTS.md) |
| Active implementation status | [PLATFORM-IMPLEMENTATION-HANDOFF.md](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-work/PLATFORM-IMPLEMENTATION-HANDOFF.md) and roadmap |

## Key ADRs

- [ADR-06 — workload identity and subject delegation](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/adr/ADR-06-workload-identity-and-subject-delegation.md)
- [ADR-10/A1 — obligations carried forward](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/adr/ADR-10-A1-obligations-carried-forward.md)
- [ADR-14 — internal transport is gRPC](https://github.com/datakaveri/cdpg-claude/blob/dev/claude-docs/adr/ADR-14-internal-transport-is-grpc.md)

## Implementation evidence

The module guide was checked against [`dx-common-go`](https://github.com/datakaveri/dx-common-go). Integration claims were checked against `dx-gateway-go`, `dx-authz-go`, `dx-acl-go`, `dx-user-go`, `dx-catalogue-go`, `dx-registry-go`, `dx-audit-go`, `dx-notification-go`, `dx-files-connect-api-go`, `dx-dataplane-rs-go`, `dx-dataplane-ogc-go`, `dx-agent-registry-go`, `dx-mcp-gateway-go`, `dx-agent-runtime-go`, orchestration, and `dx-gitops` on their current development branches.

## Decisions still needed

- Reconcile the OPA target with the conflicting evaluator ADR direction; approve policy/bundle/runtime governance.
- Approve composite authorization and carried-decision wire contracts, revisions, expiry, attestation, revocation, and typed-obligation schema.
- Approve request-bound represented-subject serialization/binding.
- Normalize ACL item/access values, gateway mappings, OpenFGA types/relations, and agent relation validation.
- Complete internal gRPC adoption and streaming interception decisions.
- Define projection watermark/revision and reconciliation contracts.

Until these are resolved, the safe interim behavior on the relevant page applies, normally default deny.
