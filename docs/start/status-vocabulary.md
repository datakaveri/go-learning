---
title: Status vocabulary
description: The required meanings of implementation-status labels in the CDPG Go developer portal.
---

# Status vocabulary

**Status: Implemented.** Every architecture, integration, and tutorial page uses this vocabulary. A status describes the evidence available on 10 August 2026; it is not a release or support guarantee.

| Label | Meaning |
|---|---|
| **Implemented** | The behavior exists in source, has relevant automated tests, and is integrated in its current repository context. It may still lack production rollout evidence. |
| **Partially implemented** | A useful subset exists, but a required path, contract, test, or deployment integration is missing. |
| **In development** | Active implementation is visible, but the capability is incomplete or not ready to depend on. |
| **Planned** | Target architecture or approved direction without a complete implementation. Examples on planned pages are explicitly labelled pseudocode. |
| **Deferred** | Intentionally outside the current delivery scope. No service should depend on it. |
| **Legacy only** | Vocabulary reserved for documentation governance. It is not used to describe the Go architecture in this portal. |
| **Superseded** | Replaced guidance. The page should point to the current canonical explanation and must not be used for new work. |

## How to read a page

The status at the top applies to the page's main subject. A subsection can carry a more specific status. For example, relationship checks can be **Implemented** while contextual OPA policy is **Planned** on the same authorization page.

Three distinct questions must not be collapsed into one label:

1. Is the code present?
2. Is it integrated through local orchestration or GitOps?
3. Is it operationally proven in a production environment?

The platform is before its first production release. “Implemented” therefore means source-backed and integrated to the stated extent, not “production-proven.”

## Safe behavior when status is uncertain

When evidence is missing or contradictory:

- record an architecture gap instead of filling it with a plausible design;
- identify the required owner or Architecture Decision Record (ADR);
- reject an unsafe configuration at startup;
- default deny security decisions;
- keep externally reachable routes disabled when authorization cannot be enforced correctly.

