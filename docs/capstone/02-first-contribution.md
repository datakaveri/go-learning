---
title: First Contribution
sidebar_label: First Contribution
description: Select, scope, implement, and verify a low-risk production contribution.
---

# First contribution

## Choose evidence-backed work

Good first contributions:

- a missing test for a supported behavior;
- an explicit security or validation defect with a small fix;
- route/OpenAPI drift;
- an unbounded timeout, body, sort, or concurrency setting;
- a service-local duplicate of an existing platform seam;
- a documentation mismatch verified against source;
- one roadmap item already scoped by maintainers.

Avoid broad gateway, authorization-model, schema, fleet-wide, or cross-service changes until you have traced and tested the affected contracts.

## Understand before editing

1. Reproduce or observe current behavior.
2. Identify owning repository and layer.
3. Trace callers, stores, events, deployment, and client contract.
4. Read the applicable platform and SDK reference.
5. Write the smallest acceptance criteria.
6. Record risk, compatibility, and rollback.

## Implement

- preserve unrelated local work;
- keep the diff focused;
- use current dx-common-go seams;
- add the failing test first when practical;
- include negative, cancellation, and authorization cases;
- update the authoritative documentation only;
- do not mix mechanical cleanup with behavior change.

## Verify

Run repository unit, race, vet, integration, contract, and smoke checks proportional to risk. Record exact results and skipped checks. Inspect the final diff for secrets, generated files, accidental API changes, and unrelated formatting.

## Review handoff

Your PR description should answer:

- What user or operator problem is solved?
- What source evidence established the old behavior?
- Which contracts change or remain unchanged?
- What are the security and operational effects?
- Which commands passed?
- How is the change deployed, observed, and reversed?
- What is deliberately outside scope?

## Post-deploy

Watch rollout health, request error and latency, relevant dependency metrics, event backlog, and domain outcome. A contribution is not finished at merge if the delivery workflow includes deployment verification.

## Graduation checkpoint

Explain the change to a reviewer from the gateway boundary through application, data, events, tests, and deployment. If any arrow relies on assumption, identify how you would prove it.
