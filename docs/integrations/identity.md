---
title: Authentication and workload identity
description: User, application, workload, delegated-agent, organisation, and authorization identities.
---

# Authentication and workload identity

**Status: Partially implemented.** OIDC/JWT validation and destination-bound workload credentials exist. Agent token exchange exists in the local Agentic Plane path. A request-bound represented-subject artifact remains an architecture gap.

![User and workload identity flow](/img/architecture/user-workload-identity.svg)

## Keep the concepts separate

| Concept | Question | Authority |
|---|---|---|
| User authentication | Who is the human? | Keycloak-issued OIDC token, validated at the gateway |
| Application authentication | Which client application is calling? | Keycloak/application credential flow |
| Workload authentication | Which CDPG service made this internal call, for which destination? | Destination-bound service token, verified by callee |
| Represented subject | Which human’s authority is carried through a trusted workload? | Verified subject context accepted only from an allowlisted asserter |
| Agent actor | Which agent acts for the subject? | Token-exchange actor/delegation context plus agent registry state |
| Organisation context | Which organisation boundary applies? | Verified identity/membership and operation input; never a header alone |
| Authorization | May these identities perform this action on this object now? | Gateway/application PEP using `dx-authz-go`, OpenFGA, planned OPA, and business state |

Keycloak authenticates and issues tokens. It does not own resource relationships, grant conditions, row filters, or service business decisions.

## User flow

1. A browser uses OIDC Authorization Code with Proof Key for Code Exchange (PKCE).
2. Keycloak authenticates the user and returns short-lived tokens to the approved client.
3. The client sends the access token to `dx-gateway-go`.
4. The gateway validates signature against JWKS, allowed algorithm, issuer, audience, expiry/not-before/issued-at with bounded clock skew, and required subject/client claims.
5. The gateway establishes verified subject context and evaluates route posture.
6. It obtains a workload token for the destination and forwards only sanitized internal identity context.
7. The destination validates workload identity first, then accepts represented subject only if the workload is a configured subject asserter.

Key rotation is handled by JWKS refresh with a bounded cache. Unknown key IDs trigger a controlled refresh; they do not disable signature verification. Cache expiry must not extend token validity.

## Application and workload flow

An application may use an approved client-credentials path for operations designed for applications. Internal services use client credentials through the workload issuer, with an audience derived from the destination: `dx:svc:<service>`.

```go
token, err := issuer.TokenFor(ctx, "dx-example-go")
if err != nil {
    return fmt.Errorf("credential for dx-example-go: %w", err)
}
```

Prefer the supplied HTTP transport or gRPC client integration so credentials, tracing, and refresh are applied consistently. The callee validates issuer, audience, time bounds, caller allowlist, and—before subject propagation—the subject-asserter allowlist.

## Delegated agent flow

The Agent Runtime exchanges the user token under RFC 8693 semantics to obtain a delegated token. The human remains `sub`; the agent appears as the actor. Effective authorization requires both:

1. the human subject may perform the resource relation/action; and
2. the agent has an active, scope/resource-bounded delegation and is not killed or revoked.

High-risk tools additionally require human-in-the-loop approval. An approval does not widen the delegation or replace authorization at execution time.

## Claims and audit context

Validate only claims required by the contract. At minimum, verify issuer, audience, signature algorithm, expiry, not-before/issued-at policy, subject or client identity, and actor/delegation claims when present. An organisation claim is context to validate, not proof of membership by itself.

Audit request ID, trace ID, subject ID, actor ID/kind, workload ID, organisation ID, delegation/grant ID, authorization decision reference when available, action/resource, outcome, and reason. Never log raw tokens, client secrets, private keys, approval codes, or presigned URLs.

## Current architecture gap

Services currently rely on destination-bound workload authentication plus a subject-asserter allowlist before consuming projected subject context. The exact serialized subject assertion and its cryptographic binding to method, path, body digest, destination, and request lifetime are not approved. The security and platform teams must define this in an ADR and contract tests.

Safe interim behavior:

- strip all external internal-identity headers at the gateway;
- require workload verification at authenticated service boundaries;
- keep `subject_asserters` empty unless the caller is explicitly approved;
- authorize using the verified subject and actor, never arbitrary headers/body fields;
- deny if delegated context is required but incomplete or inconsistent.

## Failure tests

Test missing, malformed, expired, not-yet-valid, wrong issuer, wrong audience, unknown key, disallowed algorithm, disallowed caller, disallowed subject asserter, mismatched organisation, absent delegation, revoked agent, and clock-skew boundaries. Optional authentication additionally tests “absent is anonymous, invalid is rejected.”

