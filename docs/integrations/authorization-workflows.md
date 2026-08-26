---
title: Authorization workflows
description: Decision and enforcement workflows across gateway, services, data, marketplace, files, subscriptions, and agents.
---

# Authorization workflows

**Status: Mixed.** Each workflow names its current status and target gaps. “Decision point” is where allow/deny is computed; “enforcement point” is where execution is stopped or obligations are applied.

## Common rule

Authentication establishes subject, actor, workload, and organisation context. Authorization maps those verified identities to an action and object. A PEP requests/evaluates the decision, enforces it, and audits the outcome. Transport success is not authorization success.

## Gateway-level relationship check

**Status: Implemented for configured routes; mapping coverage is partial.**

- Actors: external subject/application; gateway workload.
- Preconditions: valid route, explicit auth mode, extractable resource ID, canonical relation mapping.
- Inputs: verified identity, method/path, mapped resource/relation.
- Decision: `dx-authz-go` → OpenFGA.
- Enforcement: gateway rejects before proxying.
- Data: OpenFGA projection only; gateway does not read service DB.
- Cache: no revision-aware composite cache contract; never cache beyond identity/grant validity.
- Failure: malformed/missing mapping or PDP failure denies.
- Signals: request/trace, route, subject/actor/workload, object/relation, result/reason, latency.

## Service-level object authorization

**Status: Implemented as distributed service responsibility; consistency is partial.**

1. Workload and subject context are verified.
2. Application loads the minimum service-owned state needed to identify object, owner, organisation, and transition.
3. It checks required OpenFGA relation through the current authorization API or an approved in-service adapter.
4. It enforces state-specific rules in the use case.
5. It commits only after allow.

The service is the enforcement point. An unavailable authorization dependency fails closed. Do not trust a role, organisation header, or gateway check for an object that changed after the edge decision.

## Organisation-scoped mutation

**Status: Partially implemented.**

- Verify caller membership/administrative relationship to the target organisation, not merely a claimed organisation ID.
- Verify the resource belongs to that organisation.
- Verify the requested transition is allowed for the role/relation.
- Scope every query/update by both resource ID and organisation where possible.
- Audit cross-organisation denials and test identifier swapping.

Membership changes flow through `org.member.*`; group producer coverage is incomplete. Projection lag never creates an allow.

## OPA contextual evaluation and obligations

**Status: Planned.** `dx-authz-go` will pass bounded context to OPA, combine its result with OpenFGA, and return typed obligations. Services do not call OPA or compile policy themselves. Unsupported or untranslated obligations deny.

## Revocation and cache invalidation

**Status: Partially implemented.** ACL/agent registry changes authoritative state and publishes events; authorization updates tuples; enforcement/caches observe changes. Revision-bearing decisions and universal invalidation are planned. For high-risk operations, keep TTL short or validate active state at execution. Monitor projection lag and reconcile.

## Agent dual-stage tool authorization

**Status: Partially implemented locally.**

1. Runtime verifies session lease and agent active state.
2. Human subject’s resource relation must allow the action.
3. Agent’s `delegated_*` relationship/scope must allow the same resource/action.
4. MCP Gateway validates tool/schema/resource mapping and semantic-firewall rules.
5. Planned OPA evaluates risk/context/tool guardrails.
6. High-risk action pauses for human approval; the approval is single-use and bounded.
7. Immediately before invocation, kill switch/revocation and authorization are rechecked.
8. Tool call receives the narrowly scoped credential/capability.
9. Output is marked untrusted, bounded, and audited before returning to the planning loop.

Any failed stage denies. Approval does not override a later revocation. Local integration exists; GitOps and full operational evidence do not.

## Marketplace entitlement

**Status: Partially implemented.**

- Marketplace owns product/purchase/payment state and idempotency.
- Payment webhook route is exact/public at the user-auth layer but verifies provider authenticity and event ID.
- A completed purchase requests/causes an ACL entitlement; ACL owns the grant.
- `policy.*` projects the relationship to OpenFGA.
- Access remains pending until projection verifies allow.
- Duplicate webhook/event does not duplicate purchase or grant.
- Refund/expiry/revocation removes authority and triggers reconciliation/audit/notification.

## File download

**Status: Partially implemented; carried decision planned.**

The gateway/file application PEP authenticates and checks object access. File service verifies workload/subject, confirms metadata/processing state, owner/organisation, and authorization. It returns/uses a short-lived object capability bound to object/action/subject where supported. Never log or reuse presigned URLs. Data Plane decision attestation and typed constraints are planned; deny if a required restriction cannot be enforced.

## Subscription creation and delivery

**Status: Partially implemented.**

Creation authorizes subject/resource/query, stores an idempotent subscription and authorization reference, and emits audit. Continued delivery revalidates on expiry/revision/revocation or a bounded schedule; a one-time allow is not indefinite authority. Unsupported filter/limit obligations deny. On revocation, stop delivery, record cursor/state, and audit.

## Audit correlation

**Status: Implemented event service; producer completeness partial.**

Every security-sensitive workflow propagates request ID, trace ID, event ID, subject, actor, workload, organisation, action/resource, decision reference/revision when available, enforcement result, and safe reason. Audit production uses an outbox when tied to a business transaction; consumer is idempotent. A temporary audit outage must not silently erase required audit facts.

