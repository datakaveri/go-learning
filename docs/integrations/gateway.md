---
title: Gateway integration
description: Register and secure a CDPG service route through dx-gateway-go.
---

# Gateway integration

**Status: Implemented routing and authentication; authorization profiles and some resilience controls are partially implemented.** `dx-gateway-go` is the public edge and a Policy Enforcement Point (PEP). It is not the Policy Administration Point (PAP), Policy Decision Point (PDP), relationship engine, identity provider, or owner of service business rules.

![Gateway route and trust flow](/img/architecture/gateway-integration.svg)

## Preconditions

- The service owns a unique public prefix and OpenAPI contract.
- Its internal upstream and health endpoints exist.
- Required user/app authentication is defined.
- The service’s workload audience and allowed callers/subject asserters are provisioned.
- Object/action mapping is registered or the route remains disabled.
- Body, timeout, rate, and streaming behavior are known.

## Verified route fields

The current `config.Route` supports:

| Field | Meaning |
|---|---|
| `path_prefix` | Public path prefix. Matching respects path-segment boundaries. |
| `upstream` | Internal service base URL. |
| `enabled` | Explicit false returns not found without proxying. |
| `strip_prefix` | Remove the public prefix before forwarding. |
| `exact` | Match only the exact path; essential for a public callback beside protected routes. |
| `auth_mode` | `required`, `optional`, or `none`. Always set explicitly. |
| `workload_destination` | Destination service name used to obtain an audience-bound workload credential. |
| `allow_app_auth` | Permit the verified application-auth path when the route supports it. |
| `require_relation` | Non-empty OpenFGA relation name that the gateway PEP checks, for example `viewer`. |
| `resource_type` | Relationship object type. Must match the canonical authorization vocabulary. |
| `resource_id_from` | Verified extraction rule for the resource identifier. |
| `flush` | Boolean switch that disables reverse-proxy response buffering for streaming/SSE. |
| `keep_authorization` | Preserve external bearer only for a reviewed token-exchange flow. Default is to remove it. |

The older boolean authentication switch exists in source only as a compatibility field. New route entries use `auth_mode`.

## Authentication modes

| Mode | No bearer token | Valid bearer token | Invalid bearer token |
|---|---|---|---|
| `required` | 401 | Verified subject propagated | 401 |
| `optional` | Anonymous request | Verified subject propagated | 401 |
| `none` | Accepted at OIDC layer | Not used as user authentication | Not used |

`none` is appropriate for health, public metadata, or a webhook with its own authenticity check. It does not remove schema validation, rate limits, body limits, provider signature verification, idempotency, or service authorization.

## Worked route

```yaml title="Illustrative route using verified current fields"
routes:
  - path_prefix: /example
    upstream: http://dx-example-go:8080
    enabled: true
    strip_prefix: true
    exact: false
    auth_mode: required
    workload_destination: dx-example-go
    require_relation: "viewer"
    resource_type: resource
    resource_id_from: path:1
    flush: false
    keep_authorization: false
```

If the client calls `GET /example/widgets/2c12...`, prefix stripping forwards `/widgets/2c12...`. The extraction expression must be tested against the real public route; this example is not permission to reuse `path:1` for a different shape.

Corresponding destination policy:

```yaml
workload:
  enforcement: required
  service: dx-example-go
  allowed_callers: [dx-gateway-go]
  subject_asserters: [dx-gateway-go]
```

## Routing behavior

Routes are ordered by longest prefix. Prefix matches stop at a segment boundary, so `/cat` does not own `/catalogue-extra`. Place an `exact: true`, `auth_mode: none` callback before/protected alongside a broader route; the matcher still uses longest/exact semantics rather than source order as a security mechanism.

The gateway strips all client-supplied internal workload/subject headers before adding verified context. It validates the external OIDC token, obtains a credential for `workload_destination`, and sends the request upstream. The service verifies that credential before it honors subject context.

## Identity propagation

Propagate separate values for:

- authenticated human subject;
- acting agent/application when delegated;
- delegation/grant context;
- organisation scope;
- calling workload;
- request and W3C trace context.

The current implementation projects subject context after workload authentication. A cryptographically request-bound serialized subject assertion has not been approved; see [Identity](./identity.md). Do not invent headers or treat unsigned client input as trusted.

## Authorization profiles

Gateway relationship checks are suitable when a resource ID can be extracted unambiguously before proxying. The destination still enforces object state, ownership changes, organisation boundaries, typed obligations, and any authorization that depends on loaded business data.

When `require_relation` is non-empty, missing mapping, invalid identity, authorization dependency failure, or denial must stop the request. A route that cannot map its operation safely remains disabled.

## Limits and streams

- Set route/service body limits for JSON and a separate reviewed path for large uploads.
- Set finite upstream dial/header/idle timeouts; current gateway transport timeout hardening is a tracked gap.
- Rate-limit on trusted client address plus verified subject/actor/workload as appropriate. The limiter’s dependency failure posture must be explicit.
- For SSE, configure `flush`, disable buffering/compression on the matching route, honor cancellation, and test duration beyond ordinary timeouts.
- For downloads/uploads, stream with backpressure and avoid logging URLs, ranges, or protected metadata unnecessarily.

## Route tests

At minimum verify:

```go title="Test intent; use gateway's current config/proxy test helpers"
func TestExampleRoutePosture(t *testing.T) {
    // longest-prefix and segment-boundary match
    // prefix is stripped exactly once
    // missing token is 401; invalid optional token is also 401
    // spoofed internal identity headers are removed
    // workload token audience is dx:svc:dx-example-go
    // relationship denial/dependency error does not reach upstream
    // streaming response flushes and cancels cleanly
}
```

Run the gateway’s route-table and GitOps conformance suites after changing local and deployment route inventories. Local and GitOps configuration must agree on public exposure, auth mode, destination, and enabled state.

## Failure and observability

Emit route ID/prefix, auth mode, verified caller kind, destination, authorization outcome/reason, status, latency, upstream error class, rate-limit result, and trace/request IDs. Never log bearer tokens or propagated credentials.

| Failure | Edge behavior |
|---|---|
| Unknown/disabled route | 404 |
| Required token absent/invalid | 401 |
| Optional token invalid | 401, never anonymous fallback |
| Workload credential acquisition fails | 503; do not proxy uncredentialed |
| Relationship denied | 403 |
| Relationship engine unavailable | Fail closed; status follows safe gateway mapping |
| Upstream timeout/unavailable | 502/503/504 with safe problem response |

## GitOps checklist

- Add the service to the canonical exposure inventory.
- Add the route to every intended environment value file.
- Provision Keycloak audience/client scope before enabling the route.
- Add NetworkPolicy from gateway to service and service to required dependencies.
- Run route conformance tests and a negative identity matrix.
- Keep the route disabled until the service deployment, policy projection, and readiness checks are available.
