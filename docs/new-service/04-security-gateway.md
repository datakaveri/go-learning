---
title: 4. Identity, authorization, ACL, and gateway
description: Secure a service boundary and expose it through the CDPG gateway.
---

# 4. Identity, authorization, ACL, and gateway

**Status: Workload identity and OpenFGA checks implemented; composite OPA decisions and carried Data Plane decisions planned.** This stage must be complete before public exposure.

## Step 1 — model object/action authorization

For every operation define:

| Operation | Resource ID source | Required relation/action | State/organisation rule | Planned contextual obligations |
|---|---|---|---|---|
| Create widget | collection/configured scope | `editor` on collection | subject organisation may create | purpose/risk if approved |
| Get widget | path after normalized lookup | `viewer` on widget | row organisation equals verified scope | field mask/result audit |
| Delete widget | loaded widget | `owner` or explicit admin relation | state permits delete | high-risk reason/approval if applicable |

The relationship vocabulary must exist in the canonical model. If it does not, stop and request the authorization ADR/model change.

## Step 2 — call the current authorization service

Use the current typed `dx-common-go/auth/fga` client only for `POST /v1/check`; its additional policy methods do not match the current PAP/PDP boundary and must not be used. Attach a workload issuer configured for destination `dx-authz-go`.

```go
decision, err := authz.Check(ctx, fga.CheckRequest{
    SubjectType:  fga.SubjectTypeUser,
    SubjectID:    subject.ID,
    ResourceType: "resource",
    ResourceID:   widgetID,
    Relation:     "viewer",
})
if err != nil {
    return platformerrors.ServiceUnavailable("authorization unavailable").WithCause(err)
}
if !decision.Allowed {
    return platformerrors.Forbidden("access denied")
}
```

For an agent, enforce the human relation and agent delegated relation plus active/kill-switch state. Do not rely on the current handler accepting an `agent` subject—the server/client validation mismatch is a known gap; use the approved dual-stage path and contract tests.

## Step 3 — enforce business state

The application service remains a PEP. It loads the authoritative object after coarse edge mapping, checks organisation and state, calls authorization, and applies any obligations. A gateway allow does not authorize an object that was swapped or changed.

Planned OPA and obligation examples remain pseudocode. Until implemented, do not claim purpose, field-mask, row-filter, spatial, temporal, quota, or risk policy has been evaluated.

## Step 4 — administer grants through ACL

Service code does not write OpenFGA tuples. Resource owners/admin workflows use `dx-acl-go`; ACL validates ownership and publishes `policy.*`; `dx-authz-go` projects tuples. Treat a new grant as pending until a bounded check observes it. Revoke through ACL and verify denial after projection.

If the service creates a resource, publish the approved ownership fact/outbox event so the authorization projection gains the owner tuple. Do not write another service’s policy database.

## Step 5 — register the gateway route

Add explicit `path_prefix`, upstream, enabled, strip/exact behavior, `auth_mode`, `workload_destination`, relation mapping when safe, body/timeout/rate/stream settings. Update both local orchestration and GitOps inventories.

Destination service configuration requires workload verification, expected audience derived from its service name, caller allowlist, and `dx-gateway-go` in `subject_asserters` only when it forwards subjects.

## Step 6 — audit

Emit an audit fact for grant-relevant changes, denials that need investigation, protected reads, administrative actions, and high-risk agent/tool operations. Correlate request/trace, subject/actor/workload, org, resource/action, decision reference, result, and event. Never log credentials, protected payloads, or presigned URLs.

## Security verification matrix

```text
no credential                 -> 401
invalid credential            -> 401
wrong workload audience       -> 401/403 before handler
disallowed caller             -> 403 before handler
disallowed subject asserter   -> 403 before subject use
valid identity, denied tuple  -> 403
authorization unavailable     -> fail closed
other-organisation object ID  -> 403/404 per disclosure policy
revoked grant/agent           -> deny
unsupported planned obligation-> deny
```

Common mistakes: role-only authorization, trusting organisation/body fields, direct tuple writes, direct OPA calls, caching allow without revision/expiry, or enabling the route before negative tests.

Verification: gateway route, workload, subject assertion, OpenFGA model, default-deny, organisation isolation, revocation, and audit-correlation tests. Expected: no denied or ambiguous request reaches a state-changing adapter. Next: [tests and local integration](./05-tests-local.md).

