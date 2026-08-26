---
title: Control Plane architecture
description: Identity, governance, catalogue, marketplace, registry, credits, audit, notification, subscription, ownership, events, and failure behavior.
---

# Control Plane architecture

**Status: Partially implemented.** Most service repositories and local integrations exist. Shared-platform adoption, event completeness, composite authorization, and production evidence remain incomplete.

The Control Plane owns governance and commercial/administrative state. It never becomes a shared business database: each service owns its records and exposes an API or event contract.

## Service boundaries

| Service | Owns | Synchronous dependencies | Events/failure notes | Status |
|---|---|---|---|---|
| `dx-user-go` | Profiles, organisations/membership, app credentials/delegation state, identity-provider coordination | Keycloak | Produces `org.member.*`; current publisher/outbox adoption partial; group producer gap | Partially implemented |
| `dx-acl-go` | Grants, access requests, delegation policy administration | Catalogue ownership lookup | `policy.*`, `delegation.*`; projection lag/reconciliation required | Partially implemented |
| `dx-authz-go` | Decision API and OpenFGA projection | OpenFGA | Consumes policy/membership/delegation; planned OPA composite path | Partially implemented |
| `dx-catalogue-go` | Dataset metadata, discovery/search representation, ownership facts | Elasticsearch; internal APIs | Ownership/outbox projection; gRPC surface exists | Partially implemented |
| `dx-marketplace-go` | Listings, orders/purchase/payment state | Payment provider, ACL/entitlement workflow | Provider-event idempotency implemented; entitlement reconciliation required | Partially implemented |
| `dx-registry-go` | Resource/ACL server registrations | PostgreSQL | Administrative CRUD and audit | Implemented |
| `dx-credits-go` | Credit balances and billing entries | PostgreSQL/events | Owner of credit mutation/idempotency | Implemented |
| `dx-subscription-go` | Subscription definition/delivery state | Data/event sources | Side-effect idempotency implemented; continued authorization remains partial | Partially implemented |
| `dx-audit-go` | Audit projection and query/export | RabbitMQ/PostgreSQL | Idempotent consumer and read API; producer coverage varies | Implemented |
| `dx-notification-go` | Notification templates/dispatch state | RabbitMQ/provider | Reconnecting consumer; provider effect must be idempotent/reconcilable | Implemented |
| `dx-community-layer-go` | Discussion/challenge state | Owned database/search as configured | Independent bounded context | Partially implemented |

## Identity and membership synchronization

1. `dx-user-go` authenticates administrative intent through trusted subject/workload context.
2. It updates authoritative organisation membership and coordinates identity-provider roles/groups only where that is its contract.
3. State plus event should commit atomically; current publication path still needs consistent outbox adoption.
4. `dx-authz-go` consumes `org.member.*` and updates the OpenFGA projection idempotently.
5. Authorization changes only after projection. Missing producer/event/mapping denies; reconciliation detects drift.

Keycloak is the authentication authority, not the definitive owner of every resource relationship. Organisation isolation uses verified membership plus service-owned resource organisation, never a client field alone.

## Dataset onboarding and discovery

1. Provider submits metadata through gateway/catalogue; gateway and catalogue enforce identity and ownership.
2. Catalogue validates identifiers, provider/organisation, schema, resource type, and supported search fields.
3. Catalogue writes its owned document/index and a durable ownership relationship event.
4. Authorization projection makes the provider/organisation the resource owner.
5. Registry/data/file services may be referenced through contracts; they do not share catalogue storage.
6. Public/optional discovery returns only the view allowed by authentication/authorization. Invalid optional credentials are rejected.

Failures: invalid metadata is rejected; search unavailable fails or uses an explicitly correct degraded path; ownership publication backlog is observable and the resource is not considered safely accessible until projection succeeds.

## Marketplace purchase and entitlement

1. Marketplace creates/retrieves an idempotent order.
2. Client/provider completes payment; exact webhook route verifies provider authenticity and stable provider event ID.
3. Marketplace commits payment/order result once.
4. Entitlement orchestration creates an ACL grant; ACL remains the grant owner.
5. `policy.*` projects to OpenFGA.
6. Consumer access stays pending until authorization verifies the relationship.
7. Expiry/refund/revocation removes entitlement and triggers audit/notification/reconciliation.

Unknown provider outcome is reconciled before retry. Duplicate callback cannot duplicate order, credit, or entitlement.

## Registry, credits, audit, and notifications

Registry records discoverable service endpoints/contracts but does not make them reachable or trusted; gateway/GitOps/identity still govern exposure. Credits owns monetary/usage balance transitions and must use locks/idempotency. Audit consumes immutable security/business facts into an append-oriented projection. Notifications consume facts and own dispatch attempts; they must not become the authoritative business state.

## Configuration and governance

Each service validates its dependencies/security posture at startup. Central governance defines service names, ports, public prefixes, workload audiences, authorization vocabulary, event names/versions, data classification, and deployment policies. Services own domain configuration. Runtime flags cannot bypass required authentication or authorization.

## Common failure posture

- Required owner/PDP/database unavailable: fail request/readiness; never fabricate an allow.
- Broker unavailable: authoritative write succeeds only with durable outbox; backlog alerts and later dispatch.
- Projection lag: report pending/deny for access-sensitive flows and reconcile.
- Duplicate event/webhook: return/reuse idempotent result.
- Notification/audit consumer failure: retry/DLQ/replay; preserve authoritative source fact.
- Cross-organisation identifier mismatch: deny and audit.

See [ACL](../integrations/acl-policy.md), [authorization](../integrations/authorization.mdx), [messaging](../operations/messaging-workers.md), and [observability](../operations/observability.md).

