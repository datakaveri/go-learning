# CDPG Go developer portal update plan

Status: executed against source and approved-design evidence available on 2026-08-10.

This disposition records the documentation decisions made before changing the production portal. It is intentionally outside `docs/` so it does not appear as learner-facing content.

## Evidence baseline

- Canonical platform architecture: `cdpg-docs` and the current `claude-docs` architecture, security, service-standard, persistence, messaging, testing, deployment, ADR, and implementation-handover material.
- Implementation evidence: `dx-common-go`; gateway, ACL, authorization, user, catalogue, registry, audit, notification, file, Data Plane, and Agentic Plane repositories; orchestration configuration; and `dx-gitops`.
- Port authority: `claude-docs/PORTS.md` only.
- Status date: 2026-08-10. The platform is before its first production release, so this portal must not equate a merged implementation with a production-proven capability.

## Disposition

| Area | Action | Reason |
|---|---|---|
| Home, learning paths, and roadmap | Rewrite | Make the developer guide primary and preserve the progressive Go curriculum as a separate path. |
| Modules 0–3 | Update links and conflicting platform claims; retain Go teaching content | The language material remains useful, but platform references must point to one canonical explanation. |
| Module 4 platform lessons | Rewrite as guided learning pages over the canonical developer guide | Existing pages duplicate architecture and include obsolete or unsupported trust statements. |
| Module 5 capstone | Rewrite | Align the capstone with the verified service layout, shared platform APIs, security posture, test gates, and GitOps workflow. |
| Architecture orientation | Create | Explain Control, Data, and Agentic Planes; public REST/internal gRPC; service ownership; events; trust boundaries; and status. |
| Service architecture and Go standards | Create | Establish dependency direction, project layout, API rules, security rules, Go style, and a review scorecard. |
| New-service tutorial | Create as a progressive series | Take a teaching service from bounded context through bootstrap, persistence, events, gateway, authorization, tests, container, Compose, and GitOps. |
| Shared platform module guide | Create | Document only packages present in `dx-common-go/platform`; identify top-level foundation packages and missing platform modules explicitly. |
| Gateway and identity integration | Create | Document verified route fields, auth modes, workload credentials, subject assertion controls, and failure behavior. |
| ACL, OpenFGA, authorization, and OPA | Create | Establish PAP/PDP/PEP boundaries, current OpenFGA API/model, planned composite decision and OPA responsibilities, lifecycle, gaps, and safe interim behavior. |
| Data Plane authorization | Create | Document planned carried decisions and typed obligation translation without claiming an implemented wire format. |
| Persistence, messaging, workers, observability, testing, and deployment | Create | Provide practical platform-specific operating and verification guidance. |
| Navigation and categories | Rewrite | Add audience-based paths and canonical developer-guide categories before the curriculum. |
| Existing Mermaid diagrams | Retire where they encode platform architecture; retain only where they teach generic Go | Replace platform diagrams with accessible SVGs in the established CDPG visual language. |
| New architecture diagrams | Create | Provide focused diagrams for planes, layering, service lifecycle, modules, transports, gateway, identity, policy, authorization, Data Plane enforcement, events, workers, local topology, GitOps, and current/target state. |
| Unverified API examples | Retire or relabel | Use current source-backed APIs; label target interfaces and unapproved formats as pseudocode. |

## Architecture gaps to expose

1. **OPA governance:** the requested target assigns contextual policy to OPA, while an existing ADR selects a Go-native constraint evaluator. OPA remains **Planned** in this portal. Runtime placement, bundle authority, signing, distribution, revision, rollback, and emergency-disable decisions require a reconciled ADR.
2. **Composite authorization:** `dx-authz-go` currently exposes a simple OpenFGA-backed relationship check. Decision IDs, typed obligations, revisions, expiry, cache semantics, delegation composition, and OPA orchestration are **Planned**.
3. **Carried Data Plane decisions:** the architecture forbids synchronous per-query PDP calls, but the decision artifact, integrity protection, request binding, revocation signal, and obligation schema are not implemented or fully approved. Data Plane exposure remains default deny until those contracts are complete.
4. **Subject assertion binding:** services verify destination-bound workload identity and restrict subject asserters, but the projected subject context is not yet a cryptographically request-bound artifact. The serialized target contract needs an ADR.
5. **Internal gRPC adoption:** public REST and internal gRPC is the target. Catalogue has a current gRPC surface; several services still use HTTP internally and the shared server currently covers unary interception only.
6. **Relationship vocabulary:** ACL item/access values, gateway mapping, and the OpenFGA authorization model are not fully normalized. Services must not silently invent mappings.
7. **Agentic Plane delivery:** local Compose integration exists, but GitOps registration, production streaming behavior, and complete red-team/conformance evidence do not.

## Validation gates

- TypeScript typecheck and Docusaurus production build.
- Broken internal link, anchor, navigation, and image checks through the production build plus targeted source scans.
- Go example formatting, build, vet, tests, and race tests where practical.
- XML parsing and raster rendering of every new SVG at desktop width.
- Browser check for navigation, responsive diagram layout, horizontal overflow, clipped labels, missing images, and console errors.
- Terminology/status scans, canonical-port comparison, forbidden historical-platform/trust scans, OPA-status scan, synchronous-Data-Plane-PDP scan, and `git diff --check`.
- Confirmation that no `CLAUDE.md` changed.
