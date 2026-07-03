---
title: Testing Strategy
sidebar_label: Testing Strategy
description: The platform's test pyramid — unit, integration against shared infra, contract checks, and the e2e smoke test.
---

# Testing Strategy

## Learning objectives

- Place any test you write on the platform's four-level pyramid and know what each level may assume.
- Run integration tests against the real local infrastructure.
- Understand contract testing's role during the Java→Go migration.
- Read `make dev-demo` as the executable specification it is.

## Prerequisites

- [Testing](../module-2-intermediate/testing) (mechanics), [Repos & Workflow](repo-structure-workflow)

## Time estimate

**3 hours**

## Concepts

### The pyramid, platform edition

| Level | Speed | Uses | Verifies | When |
|---|---|---|---|---|
| **Unit** | ms | fakes, httptest | logic, handlers, envelopes | every `go test`, every PR |
| **Integration** | seconds | **real** Postgres/RMQ/ES/Redis (local stack or CI containers) | SQL, topology, mapping | per PR / pre-merge |
| **Contract** | seconds | Go service vs legacy service | response parity during migration | scripted, when porting |
| **E2E smoke** | ~30 s | the whole stack | the critical paths, cross-service | `make dev-demo` before promotion |

The proportions matter as much as the levels: **most tests are unit**, a meaningful minority integration, a handful e2e. Inverting the pyramid (leaning on slow e2e for everything) makes suites minutes-long and failures unattributable.

### Unit — what you already know, with the platform norms

From [Module 2](../module-2-intermediate/testing): table-driven, stdlib-only (no testify), fakes behind consumer-defined interfaces, full-router handler tests, and the two spec tests (embedded-spec-loads, spec-covers-all-routes). CI runs them with `-race`. Per the standards, **every endpoint has handler tests** — a reviewer will look.

### Integration — real infra, honest tests

Repository fakes prove nothing about SQL. Integration tests run the repository against **actual Postgres** — and the platform's choice is to use the *shared local stack* (the same one `make dev-up` boots) rather than per-test containers:

```go
func TestPolicyRepo_Integration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — integration tests skipped")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	// … exercise real INSERT/SELECT/soft-delete against real schema
}
```

The env-guard + `t.Skip` pattern keeps `go test ./...` green for anyone without the stack while letting CI and stack-runners get full coverage. Discipline that makes shared-infra testing livable: tests create their own fixtures, clean up after themselves (or use throwaway schemas), and never assume table emptiness — someone else's test data is always a possibility. Same pattern for RabbitMQ: declare a test queue, publish, consume, assert, delete.

### Contract — the migration's safety net

While a Go service replaces a Java one, the question is never "is the Go code nice?" but "**do clients notice?**" Contract checks answer it: same request → both stacks → compare responses (shape, fields, semantics — modulo the documented envelope deltas from `CLIENT-CONTRACT-CHANGES.md`). Today this is scripted/manual during porting work rather than a permanent suite. If you port a legacy endpoint, a contract comparison belongs in your PR description.

### E2E — make dev-demo as executable spec

Read the smoke test's source (it's in the orchestrator repo) and you'll find it *is* the architecture deep-dive's diagrams, asserted: health across the fleet; JWT issuance; gateway enforcement (no token → 401, token without rights → 403, with → 200); proxying to upstreams; and the full **PAP → RMQ → PDP propagation** — create policy, allowed within ~3s; delete, denied within ~3s. That's why "dev-demo green" gates image promotion: it's the minimum claim "the platform works" reduced to a script.

E2E's role is *critical paths only* — it tells you **that** something broke, rarely **what**; the lower levels exist to answer *what*.

### What the platform doesn't do (yet)

Honest inventory, matching the review: no testcontainers convention (shared stack instead), no property-based or fuzz testing tradition, no load-testing harness in the standard kit, contract checks not yet automated. As with the library gaps: absence of these is not an invitation to freelance — it's context, and potentially your future contribution.

:::info[Platform connection]
`claude-docs/TESTING.md` is the authoritative version of this page — read it now; it's short. Then read `make dev-demo`'s implementation and map each check to the architecture diagrams. The capstone requires you to write all of: unit tests, one integration test against the real stack, and a dev-demo-style smoke script for your own service.
:::

## Exercises

1. Write an integration test for `dx-scratch-go`'s repository against the local stack's Postgres — env-guarded, self-cleaning, covering the soft-delete filter. Run it twice in a row to prove the cleanup.
2. Do the same for your RabbitMQ consumer: test queue, publish → consume → assert side effect → clean up.
3. Contract exercise: pick any endpoint that exists on both tracks (see `SERVICES.md`), call both with the same request, and diff the responses. Classify each difference: documented delta, bug, or acceptable?
4. Write `smoke.sh` for your scratch service: boots it (compose), waits for ready, runs five curl checks including one auth-denial and one full write-then-read path, exits nonzero on any failure. Under 50 lines.

## Check yourself

- What may a unit test assume that an integration test must not, and vice versa?
- Why env-guard + skip rather than failing when infra is absent?
- What question does contract testing answer that nothing else does?
- What exactly does dev-demo prove, and what does it deliberately not tell you?

## References

- [Martin Fowler: Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)
- Platform: `claude-docs/TESTING.md` (required); the `dev-demo` target's source; `CLIENT-CONTRACT-CHANGES.md`
