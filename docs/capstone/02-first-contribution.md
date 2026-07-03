---
title: First Real Contribution
sidebar_label: First Contribution
description: Choosing a starter task, working it end to end, and landing your first merged PR on the platform.
---

# First Real Contribution

## Learning objectives

Land a real, merged PR on a real service — the final milestone (M5) and the end of onboarding. Equally important: experience the full loop (ticket → design → implement → gate → review → merge) at production stakes, with a safety-appropriate scope.

## Prerequisites

- [Capstone service](capstone-service) built and reviewed.

## Time estimate

**4–8 hours** of work — spread over the review cycle's real-world latency.

## Choosing the right first task

Work with your lead/buddy. A good first contribution is **real but bounded**: it matters enough to review properly and small enough that the platform machinery — not the problem — is the learning.

Good shapes, roughly in order of preference:

1. **A parity gap** — porting a small missing legacy behavior in a partial-parity service (the `SERVICES.md` table and `ROADMAP.md`'s P1/P2 items are the shopping list). Well-specified by definition: the legacy behavior *is* the spec, and your [contract-testing](../module-4-platform/testing-strategy) habit applies directly.
2. **A hardening item** — the platform review's open findings (a missing DLQ, audit middleware absent on a service, a missing handler test). High value, crisp acceptance criteria.
3. **An endpoint addition** to a healthy service — your capstone, at production stakes.
4. **A shared-library improvement** — only if something from the [known gaps list](../module-4-platform/dx-common-go-tour) genuinely blocks other work you're doing; library PRs get more scrutiny by design.

Avoid for a *first* PR: gateway config, anything in the auth path, schema changes to legacy tables, cross-service refactors. Not because you can't — because blast radius should grow with track record.

## The workflow, at full fidelity

This is [Repos & Workflow](../module-4-platform/repo-structure-workflow)'s ticket-to-merge list, now with stakes. Points that distinguish a smooth first PR:

**Before writing code**
- Reproduce/observe the current behavior in your local stack first. For a parity gap, call the legacy endpoint and save the response — that's your spec and your test oracle.
- Write your plan in three sentences (files, approach, tests) and get a nod from your buddy *before* implementing. Cheap insurance against a week in the wrong direction.

**While writing**
- Match the service's existing conventions even where you might choose differently — consistency within a service beats your preference. Real deviations from *standards* are a finding; note them for a separate ticket rather than fixing drive-by.
- Keep the diff minimal: no opportunistic reformatting, no unrelated renames. Reviewers review what changed; noise hides signal.

**The PR itself**
- Description: what + why, the ticket link, how you tested (paste the contract diff / smoke output — evidence, not claims), and anything you're unsure about *called out explicitly*. Flagging your own open questions is a strength signal, not a weakness.
- Confirm the gate ran locally. A first PR that fails CI on `gofmt` sets exactly the wrong tone.

**Review and after**
- Respond to every comment — with a change or a (brief, genuine) case for the current form. Push fixes as new commits during review so re-review is a diff, not a re-read.
- After merge: follow your change through the [pipeline](../module-4-platform/deployment) — image, gitops bump, and your service's logs in dev. You break it, you watch it; you shipped it, you also watch it.

## After the merge

That's M5 — onboarding complete. Three habits carry forward:

- **The standards are your review voice.** Cite sections, kindly, when reviewing others — you'll start receiving review requests soon.
- **The curriculum stays a reference.** The platform-connection boxes map topics to source; the roadmap's checkpoints work in reverse when you need a refresher.
- **Close the loop.** Something in these pages wrong, stale, or missing by the time you finish? The go-learning repo takes PRs too — the next joiner learns from yours. It's also, fittingly, a fine second contribution.

Welcome to the team. 🎉
