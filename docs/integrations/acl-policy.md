---
title: ACL and policy administration
description: Create, review, project, verify, and revoke grants with dx-acl-go.
---

# ACL and policy administration

**Status: Partially implemented.** `dx-acl-go` implements policy and access-request APIs and publishes policy/delegation events. Transactional-outbox adoption and canonical relationship vocabulary are incomplete.

`dx-acl-go` is the Policy Administration Point (PAP). It validates who may author or revoke a grant, owns policy/access-request/delegation records, and emits changes. It does not authenticate callers, evaluate the OpenFGA graph, execute OPA, or enforce a downstream request.

![ACL policy lifecycle](/img/architecture/acl-policy-lifecycle.svg)

## Policy lifecycle

1. A provider, owner, or authorised organisation administrator submits a grant.
2. ACL loads the catalogue item through its required internal catalogue dependency and verifies ownership/organisation authority.
3. It validates subject, item type, policy type, expiry, constraints, and request linkage.
4. Policy state and its outbound event are committed. Current outbox adoption is mixed; a broker failure must remain visible and reconcilable.
5. `policy.*` reaches `dx-authz-go`, whose consumer updates OpenFGA tuples idempotently.
6. Enforcement may lag administration. Clients must not assume the relationship graph changed because policy creation returned 201.
7. Revocation changes authoritative ACL state and emits a removal event; caches/projections invalidate and enforcement becomes deny after projection.

## Current endpoints

Service base path: `/iudx/acl/apd/v2`. Through the current gateway prefix, prepend `/acl`.

| Method/path | Purpose |
|---|---|
| `POST /policy` | Batch create grants |
| `DELETE /policy?id=<uuid>` | Soft-delete/revoke a policy |
| `GET /policy/consumer` | Grants usable by caller |
| `GET /policy/provider` | Grants issued by caller |
| `GET /policy/organisation` | Organisation administrator view |
| `GET /policy/platform` | Platform administrator view |
| `POST /verify` | Current policy verification hook |
| `POST /access_request` | Consumer creates pending access request |
| `PUT /access_request` | Provider/organisation administrator grants or rejects |
| `PATCH /access_request/{id}` | Consumer withdraws pending request |
| `POST /access_request/has_access` | Check current access-request/policy state |
| `GET /access_request/consumer`, `/provider`, `/organisation`, `/platform` | Role-scoped listings |

The OpenAPI file in `dx-acl-go/openapi/openapi.yaml` is the payload authority.

## Create a grant

```http
POST /acl/iudx/acl/apd/v2/policy
Authorization: Bearer <user access token>
Content-Type: application/json

{
  "request": [
    {
      "policyType": "INDIVIDUAL",
      "userId": "<consumer-id>",
      "itemId": "<catalogue-item-uuid>",
      "itemType": "DATABANK",
      "expiryTime": "2026-12-31T23:59:59Z",
      "constraints": {"purpose": "research"},
      "providerComment": "Approved for the stated project"
    }
  ]
}
```

The owner and organisation are derived/validated from trusted identity and catalogue state. Do not let a request body self-assert ownership. `itemType` and constraint vocabulary must match the canonical registry; because normalization is currently incomplete, an unrecognized mapping must be rejected rather than guessed.

## Access request and entitlement

Create:

```json
{
  "itemId": "<catalogue-item-uuid>",
  "requestType": "DOWNLOAD",
  "additionalInfo": {"purpose": "research"},
  "constraints": {"maxDownloads": 1}
}
```

Decide:

```json
{
  "requestId": "<request-uuid>",
  "status": "granted",
  "expiryAt": "2026-12-31T23:59:59Z",
  "constraints": {"maxDownloads": 1}
}
```

Approval and policy creation must be atomic. Marketplace-created entitlement follows the same ownership principle: marketplace owns purchase/payment state, ACL owns the grant record, and authorization owns the projected relationship. Stable idempotency keys and reconciliation prevent duplicate payment callbacks or messages from multiplying entitlements.

## Wait for projection and verify

There is no current revision/watermark endpoint that proves a particular ACL event is visible in OpenFGA. This is an architecture gap. For administrative automation:

1. Record the returned policy ID, request/trace ID, and event correlation ID when available.
2. Poll the intended `dx-authz-go` relationship check with bounded exponential backoff.
3. Treat `allowed:false` as not yet effective or denied; stop at a short deadline and surface “projection pending.”
4. Do not bypass enforcement while waiting.
5. Alert/reconcile if outbox age, queue delay, or projection lag exceeds the objective.

Example current relationship check:

```http
POST /v1/check
Content-Type: application/json

{
  "subject_type": "user",
  "subject_id": "<consumer-id>",
  "resource_type": "resource",
  "resource_id": "<catalogue-item-uuid>",
  "relation": "viewer"
}
```

## Revoke and confirm

Call `DELETE /acl/iudx/acl/apd/v2/policy?id=<policy-uuid>`, then wait for the corresponding relationship check to return `allowed:false`. Invalidate only caches whose keys/revisions correspond to that relationship. If the broker/projection is unavailable, authoritative policy is revoked but downstream enforcement may be stale; fail closed for high-risk operations and run reconciliation.

## Events, failures, and signals

Current publishers emit `policy.*` and `delegation.*`; organisation membership comes from `dx-user-go`, while a complete group-membership producer is a known gap. Consumers must be idempotent and tolerate reordering where the contract permits it.

Audit policy ID/version, actor/subject/org, resource, operation, before/after status, expiry/constraints classification, event ID, projection outcome, and authorization decision correlation. Do not copy protected constraint values into every log.

| Failure | Safe behavior |
|---|---|
| Catalogue/ownership lookup unavailable | Reject/503; do not create an unverified grant |
| Invalid subject/type/relation/expiry | 400 and no event |
| Unauthorized administrator | 403 and audit denial |
| Database commit fails | No policy and no event |
| Broker unavailable after outbox commit | Policy remains durable; dispatcher retries; readiness/lag signal degrades |
| Projection message invalid | Quarantine/DLQ; reconcile; no inferred allow |
| Revocation projection delayed | High-risk enforcement denies or validates against a safe authoritative path; alert |
