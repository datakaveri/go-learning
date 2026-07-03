---
title: Security Model
sidebar_label: Security Model
description: The trust chain in depth — secrets, auditing, gateway-origin enforcement, input hardening, and the platform's threat posture.
---

# Security Model

## Learning objectives

- Reason about the platform's trust boundaries and what each control defends against.
- Handle secrets correctly across dev and production.
- Know which endpoints must emit audit events and how the audit pipeline works.
- Apply the input-hardening set: validation, size limits, parameterization, rate limiting.

## Prerequisites

- [AuthN & AuthZ](../module-3-advanced/authn-authz), [Service Anatomy](service-anatomy)

## Time estimate

**3 hours**

## Concepts

### Trust boundaries, enumerated

Security thinking starts with "who is trusted with what, and where is that checked":

| Boundary | Control | Defends against |
|---|---|---|
| Internet → gateway | JWT validation (JWKS, iss/aud/exp) | forged/expired/borrowed identity |
| Gateway → decision | PDP check per request | valid identity, insufficient rights |
| Gateway → upstream | HMAC-signed X-Subject-*, short validity window | header injection, tampering, replay |
| Direct → upstream | full JWT validation (resolver fallback) | bypassing the gateway with garbage |
| Admin routes | `RequireGatewayOrigin()` — HMAC or nothing | direct-access paths to privileged operations |
| Service → DB | parameterized values, allowlisted identifiers | injection |
| Everything | audit events | "who did what" going unanswerable |

Two properties of the design worth restating from Module 3: an **invalid HMAC hard-fails** (no fallback to a weaker path — a forgery attempt must not get a second chance), and OpenFGA is **deny-by-default** (missing tuple = 403; the eventual-consistency window fails safe).

### Secrets

The rules, and their whys:

- **No secrets in code, in baked config, or in Git — ever.** The baked `config.yaml` carries empty defaults; real values arrive by environment: Compose locally, External Secrets Operator (backed by a secret manager) in Kubernetes.
- **Design for rotation.** The HMAC shared secret and DB credentials must be changeable without code changes — which env-injection gives you for free, and hardcoding forecloses forever.
- **Never log a secret** — including inside DSNs (`postgres://user:PASSWORD@…` in a boot log is the classic leak) and Authorization headers. Redact at the logging call site.
- Local dev uses known default credentials (admin/admin everywhere) for convenience; the platform review flags exactly this as a **critical must-change for any real deployment**. Know which mode you're in.

### The audit pipeline

Mutating operations must be attributable, so the platform makes auditing infrastructure, not discipline: the **audit middleware** (`dx-common-go/auditing`) captures who/what/when for state-changing endpoints and publishes events **asynchronously** (a buffered worker — audit latency must not tax request latency) to the auditing exchange, where **`dx-audit-go`** consumes and persists them, exposing a read API + CSV export.

Design tension to understand: async means a crash can lose in-flight audit events — accepted for now, and the review's findings (the audit consumer lacking a DLQ, audit middleware missing from several services) are open hardening items on the roadmap's platform-hardening track. The standard for *your* code is unambiguous: **every mutating endpoint sits behind the audit middleware.** It's on the capstone rubric.

### Input hardening — defense in depth

Layers, outermost first, each catching what the previous missed:

1. **Rate limiting** — at the gateway (per-client token bucket), throttling abuse before it spends compute.
2. **Body size limits** — `MaxUploadSize` middleware; a 2 GB JSON body is a DoS, not a request.
3. **OpenAPI validation** — structural rejection before handler code ([REST](../module-3-advanced/rest-api-development)).
4. **Business validation** — the rules only the service layer can know.
5. **Parameterization + allowlists** — the last line, in the repository ([Database Patterns](../module-3-advanced/database-patterns)).

The general principle: validate at the boundary, then pass *typed, validated* values inward — inner layers shouldn't re-guess whether a string is really a UUID.

### Thinking in threats

The platform review ranks findings by exploitability (critical → high → medium) — a habit to copy. For any change you ship, ask three questions: *What new input paths did I open, and what validates them? What secrets or privileged operations does this touch? If this were misused, what would the audit trail show?* Written answers to those three questions are a perfectly good PR-description security section.

:::info[Platform connection]
Read two short sources now: `claude-docs/AUTH.md` (the canonical-string HMAC format, the JWKS caching details, the realm-issuer mapping trick) and the security findings table in `GO-PLATFORM-REVIEW.md` — the latter is the rare document that shows you *real* gaps ranked by severity, which is how security actually looks inside a living platform: not perfect, but known, ranked, and tracked.
:::

## Exercises

1. Extend your HMAC exercise from Module 3: add the `Issued-At` freshness check with a 60-second window, then demonstrate a replay attack succeeding at 30s and failing at 90s.
2. Grep a service for secret handling: where does the HMAC secret enter, what would it take to rotate it, and confirm nothing logs it (check the DSN log lines especially).
3. Wire the audit middleware into `dx-scratch-go`'s mutating routes, publish to your local RabbitMQ, and consume the events with your Module 3 consumer. You now have the platform's audit pipeline in miniature.
4. Threat-model your scratch service in writing: its trust boundaries table (like the one above), the top three abuse cases, and which existing control catches each. One page, max.

## Check yourself

- Why must invalid HMAC never fall back to JWT?
- Where do secrets live in dev vs production, and what does "designed for rotation" mean concretely?
- Why is auditing middleware rather than a convention, and why async?
- Recite the five input-hardening layers outermost-in.

## References

- [OWASP API Security Top 10](https://owasp.org/API-Security/) — read against the boundary table
- Platform: `claude-docs/AUTH.md`; `GO-PLATFORM-REVIEW.md` (security findings); GO-SERVICE-STANDARDS.md (security section)
