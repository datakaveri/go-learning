---
title: Current and target developer platform
description: Source-backed implementation status, safe interim behavior, and architecture gaps.
---

# Current and target developer platform

**Status: Implemented documentation of a partially implemented platform.** This snapshot is dated 10 August 2026.

![Current and target developer platform](/img/architecture/current-target.svg)

## Capability matrix

| Capability | Current evidence | Target | Status and safe interim behavior |
|---|---|---|---|
| Shared bootstrap/config/HTTP/SQL/cache/events/identity/health | Present in `dx-common-go/platform`; services are at different adoption levels | One declarative service composition model | Partially implemented; new services use current platform packages and do not copy older local wrappers. |
| Public REST | Gateway and service routes implemented | Operation-owned route manifests and conformance | Implemented at transport level; authorization profiles remain partial. |
| Internal gRPC | Shared unary server/interceptors and catalogue surface exist | Authenticated gRPC for internal service contracts | Partially implemented; do not add a new internal HTTP dependency without an explicit exception. |
| Workload identity | Destination-bound credentials, audience verification, caller and subject-asserter allowlists exist | Request-bound represented-subject contract | Partially implemented; verify workload first and default deny unapproved subject assertions. |
| Relationship authorization | OpenFGA projection and `POST /v1/check` exist | Composite decision with revisions, identifiers, expiry, obligations, delegation, and context | Partially implemented; services retain object/state enforcement and deny ambiguous mappings. |
| Contextual OPA policy | No runtime integration in source | OPA evaluates contextual policy and typed obligations behind `dx-authz-go` | Planned; no service calls OPA directly. Governance ADR required. |
| Data Plane carried decision | No approved artifact or verification implementation | Integrity-protected, request-bound, expiring decision with typed obligations | Planned; keep routes disabled or use explicit in-service controls without per-query PDP calls. |
| Transactional outbox | Platform outbox/dispatcher exists; adoption is mixed | All integration events emitted atomically | Partially implemented; a direct publish after commit is not an acceptable new pattern. |
| Consumers/workers | Supervision, retry/DLQ/replay and lease primitives exist; adoption is mixed | Bounded, cancellable, observable, reconcilable workloads | Partially implemented. |
| Agent Registry, MCP Gateway, Runtime | Repositories and local Compose paths exist | GitOps-deployed, operationally validated Agentic Plane | Partially implemented; do not claim deployment readiness. |
| Federation | No current core dependency | Optional future extension | Deferred. |

## Architecture gaps

### OPA decision authority

The requested target assigns contextual and attribute-based policy to Open Policy Agent (OPA), but the ADR set contains a conflicting evaluator direction. Until a reconciled ADR is approved, OPA stays **Planned**. The ADR must name the policy repository, authoring authority, reviewer, bundle builder, signer, distribution mechanism, runtime placement, readiness gate, revision contract, rollback, and emergency disable procedure.

### Composite decision contract

`dx-authz-go` currently performs a relationship check. A composite result such as `allow`, stable reason codes, decision ID, relationship/policy revision, expiry, and typed obligations is **Planned**. No client may infer an allow from transport success or invent an obligation schema.

### Carried Data Plane decision

The target is approved at the architectural level, but its serialization, integrity protection, replay boundary, request hash/binding, expiry, revision validation, revocation mechanism, and unsupported-obligation response are unresolved. The owning authorization, Data Plane, and security teams need an ADR and contract tests. The interim posture is default deny for routes that cannot enforce policy safely.

### Relationship vocabulary

ACL policy item/access values, gateway `resource_type`/relation mapping, and the current OpenFGA model are not fully normalized. The authorization team owns a canonical type/relation registry and compatibility tests. Until it exists, a missing mapping is a denial and an observable configuration error.

### Internal transport adoption

Public REST and internal gRPC is the target. Catalogue demonstrates current gRPC bootstrap, while several service-to-service calls still use HTTP. The platform team owns streaming interceptors and contract tooling; service teams own conversion of their internal contracts.

### Production evidence

Source implementation, local Compose integration, GitOps configuration, and production operation are separate gates. The platform does not yet have production-user evidence. Release owners must not promote status based only on a merged repository.

## Status update rule

Change a status only with links to implementation, tests, deployment configuration, and—where relevant—an approved ADR. Update the `last_verified` field used by the guide’s source register rather than editing status prose in isolation.

