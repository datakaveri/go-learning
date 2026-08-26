---
title: Agentic Plane architecture
description: Registry, runtime, MCP gateway, delegated authority, HITL, semantic firewall, kill switch, audit, and failure behavior.
---

# Agentic Plane architecture

**Status: Partially implemented in local integration.** The three repositories and Compose workflow exist. They are absent from GitOps, long-lived streaming rollout is not fully validated, and contextual OPA policy is planned.

The Agentic Plane lets a registered software agent act for a human subject under narrow, expiring, revocable authority. It is not an unrestricted automation channel and does not create a second authorization system.

![Agent lifecycle and tool execution](/img/architecture/agent-lifecycle.svg)

## Components and ownership

| Component | Owns | Must not own |
|---|---|---|
| `dx-agent-registry-go` | Agent records/templates, lifecycle, agent workload client coordination, credentials metadata, active/revoked/killed state, delegation proxy | Runtime sessions, tool execution, resource policy |
| `dx-agent-runtime-go` | Sessions, SSE chat/events, planning loop, durable session lease, token exchange, cancellation, result stream | Tool catalogue/enforcement, grant authoring, business datastore |
| `dx-mcp-gateway-go` | Tool discovery/invocation boundary, schema/resource mapping, semantic firewall, risk classification, HITL state, scoped credential attachment, output trust envelope | Agent planning, resource relationships, service business policy |
| `dx-acl-go` | Agent delegation grant administration and events | Agent runtime state |
| `dx-authz-go` + OpenFGA | Human relationship and delegated-agent relationship checks | HITL workflow or tool execution |
| OPA | Planned contextual agent/tool guardrails and typed obligations | Tool side effects or relationship storage |

## Registration and template lifecycle

1. Authenticated owner selects an approved template or submits allowed metadata.
2. Registry validates owner/organisation, template version, capabilities, allowed tools/scopes/resources, model/runtime configuration, and policy requirements.
3. Registry creates the agent and its distinct workload identity/client material through the approved identity path.
4. Delegation is administered through ACL/registry integration and projected as agent/resource relationships.
5. Agent moves through explicit lifecycle states. Disable/revoke/kill invalidates future sessions/tool calls and produces audit/events.

Templates are versioned and immutable once referenced. Secret values are never returned in listings or logged. Registration is idempotent where clients retry.

## Session lifecycle

1. User authenticates and requests a session for an owned/allowed active agent.
2. Runtime verifies agent state and delegation, creates durable session state, and acquires a lease so one replica executes it.
3. Runtime exchanges the user token for a delegated token: human remains subject, agent is actor.
4. SSE returns ordered, bounded progress/results and stops on cancellation, lease loss, expiry, kill, or revocation.
5. Session checkpoints support recovery; only one holder mutates execution state.

Durable lease and lifecycle-owned executor mechanisms are implemented. A complete rolling-restart/mid-tool end-to-end proof remains incomplete. Streaming beyond ordinary proxy timeout and resume behavior require deployment verification.

## Tool discovery and invocation

1. Runtime asks MCP Gateway for tools visible to the agent/template/delegation.
2. Model proposes a typed tool call; model output is untrusted.
3. MCP Gateway validates tool name/version and JSON schema, normalizes resource/action, bounds input/output, and applies semantic-firewall rules.
4. Human subject relationship and agent delegated relationship both must allow the action.
5. Planned OPA evaluates risk/purpose/context/tool guardrails and obligations.
6. High-risk call creates a single-use, expiring approval request and pauses.
7. Human approves/rejects. Execution rechecks authorization, agent state, kill switch, delegation, approval, and idempotency immediately before side effect.
8. Gateway attaches the narrow destination credential/capability and invokes the owning Control/Data Plane service.
9. Tool output is marked untrusted with provenance/size bounds before it re-enters the planning loop.

No model text, tool description, retrieved content, or tool output is trusted instruction. A semantic firewall validates structure and allowed intent; it does not replace authorization.

## Human in the loop

Approval state uses atomic single-winner transitions so approve/reject races do not both win. Approval binds actor/subject, agent/session, exact tool/action/resource/input digest, risk, expiry, and one execution attempt. It cannot widen delegation. If the side effect’s result is unknown after approval, reconcile the side-effect owner before retrying; a stale approved record is not permission for blind re-execution.

## Kill switch and revocation

Kill/revoke takes effect at registration/session/tool boundaries and via propagated state/events/caches. Runtime stops planning and cancels work; MCP Gateway denies new invocations; credentials/sessions expire or are invalidated; subscriptions/long work stop at safe checkpoints. Every high-risk call rechecks immediately before execution. Projection/cache lag is monitored; ambiguity denies.

## Threats and controls

| Threat | Control | Remaining evidence gap |
|---|---|---|
| Prompt/tool-output injection | Untrusted envelope, schema/bounds, semantic firewall, provenance | Comprehensive red-team corpus/canary suite incomplete |
| Agent exceeds user | Human relation ∩ delegated agent relation | Relationship vocabulary/server validation consistency |
| Approval replay/race | Expiring bound single-use state, atomic transition, owner idempotency | Full fault-injection evidence |
| Stolen credential | Destination audience, short lifetime, workload isolation, secret manager | Request-bound subject assertion contract |
| Killed agent continues | State recheck, lease/cancel, revocation projection | End-to-end latency objective and GitOps proof |
| Duplicate side effect | Stable idempotency at owning service, reconciliation | Coverage varies by tool/service |
| Runaway loop/cost | Step/time/token/tool/result bounds, cancellation, rate/quota | Planned contextual guardrails/quotas |

## Failure behavior and signals

Invalid model output/schema/tool/resource mapping, unavailable required authorization, planned-policy not ready, expired/revoked delegation, missing approval, kill state, lease loss, or unsupported obligation denies/stops. Tool dependency failure returns a bounded typed failure; it does not cause unrestricted replanning or duplicate side effects.

Audit registration/template changes, session/lease, token exchange reference, tool discovery, proposal, every authorization stage, approval state, execution/idempotency, output provenance, kill/revocation, cancellation, and result. Never log prompts/tool payloads containing protected data without an approved redaction/retention policy, and never log tokens/approval codes.

## Deployment status

Agent repositories are present in local orchestration, but current `dx-gitops` service values/ApplicationSets do not register them. Therefore “implemented locally” does not mean deployable. Required work includes values, External Secrets, workload clients/audiences, NetworkPolicies, long-lived SSE ingress behavior, resource/limit profiles, alerts/runbooks, adversarial/security suites, and rollout/kill/recovery drills.
