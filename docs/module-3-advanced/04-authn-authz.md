---
title: Authentication and Authorization
sidebar_label: AuthN & AuthZ
description: Authentication, workload identity, OpenFGA relationships, authorization boundaries, and delegation.
---

# Authentication and authorization

**Platform status: Partially implemented.** Relationship checks and workload identity exist; composite OPA decisions and carried Data Plane decisions are planned.

Authentication proves identity. Authorization decides whether that identity can perform a relation on an object.

## Token validation

Validate signature, issuer, audience, time claims, and intended token type. Cache JWKS keys with bounded refresh. Do not decode claims without verification. An optional-auth route accepts no token, but must reject an invalid supplied token.

## Downstream identity

The gateway removes client-supplied internal identity headers, validates the user, and calls each destination with an audience-bound workload credential. The service verifies workload identity first and accepts represented subject context only from a configured subject asserter. Workload and human subject are separate identities.

## Relationship authorization

OpenFGA tuples represent subject–relation–object facts. `dx-acl-go` and `dx-user-go` publish policy and membership events; `dx-authz-go` projects them. Application code still checks its operation-specific organisation, ownership, state invariant, and supported obligations. OPA contextual policy remains planned.

## Delegation

Delegated authority names the original human, acting principal, grant, kind, scope, resource boundary, and expiry. Agent actions additionally require active agent state and kill-switch check. Delegation intersects existing authority; it does not expand it.

## Exercise

Threat-model `GET` and `DELETE /widgets/{id}`. Define public/optional/protected status, role, OpenFGA object/relation, tenant check, forged-header test, and audit attribution.

## Check yourself

- Why is a valid token insufficient for resource access?
- Why must a service verify workload identity before accepting subject context?
- What should optional authentication do with an expired token?
- Whose ID is used for authorization and whose for audit under delegation?

Canonical reading: [Authentication and workload identity](../integrations/identity.md), [Authorization service, OpenFGA, and OPA](../integrations/authorization.mdx), and [Authorization workflows](../integrations/authorization-workflows.md).
