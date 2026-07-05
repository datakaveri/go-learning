---
title: Testing Strategy
sidebar_label: Testing Strategy
description: The DX testing strategy — unit, integration, repository, API, testcontainers via dxtest, migration testing, mocking, Elasticsearch testing, and the e2e smoke test. Real infra, deterministic seams, green everywhere.
---

# Testing Strategy

## Learning objectives

- Place any test on the platform's pyramid and know what each level may assume.
- Write repository and migration tests against **real Postgres** with `dxtest/containers` — testcontainers with a DSN fallback that **skips** (never fails) without Docker.
- Test HTTP/API behavior through the full router, and choose the right mocking strategy (fakes over frameworks; seams over globals).
- Test Elasticsearch code with a transport mock and against a real cluster.
- Read `make dev-demo` as the executable specification it is.

## Prerequisites

- [Testing](../module-2-intermediate/testing) (mechanics), [Database Patterns](../module-3-advanced/database-patterns), [Repos & Workflow](repo-structure-workflow)

## Time estimate

**4 hours**

## Concepts

### The pyramid, DX edition

| Level | Speed | Uses | Verifies | When |
|---|---|---|---|---|
| **Unit** | ms | fakes, httptest, injectable seams | logic, handlers, envelopes | every `go test`, every PR |
| **Integration / repository** | seconds | **real** Postgres/Redis/ES (testcontainers or DSN) | SQL, migrations, mappings, topology | per PR / pre-merge |
| **API / contract** | seconds | full router; Go vs legacy | endpoint behavior; response parity during migration | per PR; scripted when porting |
| **E2E smoke** | ~30 s | the whole stack | the critical paths, cross-service | `make dev-demo` before promotion |

The proportions matter as much as the levels: **most tests are unit**, a meaningful minority integration, a handful e2e. Inverting the pyramid (leaning on slow e2e for everything) makes suites minutes-long and failures unattributable.

The load-bearing principle across every level: **a test must run green for a developer with nothing but the repo.** Tests needing infra are *guarded* — they detect the infra (a container, a DSN env var) and **skip** when it's absent rather than fail. `go test ./...` is always green; CI and stack-runners get the full coverage. This is why the platform can require real-infra tests without punishing anyone who runs the suite on a plane.

### Unit — deterministic, seam-driven

From [Module 2](../module-2-intermediate/testing): table-driven, **stdlib-only (no testify)**, `cmp.Diff` for structs, full-router handler tests, and the two spec tests (embedded-spec-loads, spec-covers-all-routes). CI runs with `-race`.

The addition worth internalizing here is **injectable seams**: framework code takes its clock, its transport, its sleep/jitter as fields so a unit test drives them deterministically — no real sleeps, no real servers, no wall-clock flakiness. `resilience` and `scheduler` are the reference: their tests advance a fake clock instead of waiting. When you write library-grade code, expose the seam; when you test it, drive the seam.

### Mocking strategy — fakes over frameworks, seams over globals

The platform has a definite house style, and reviewers hold to it:

- **Hand-written fakes, not a mock framework.** Because dependencies are small consumer-defined interfaces, a double is a few lines — no `gomock`, no codegen. A `fakeStore` that records what it was asked to insert tests business logic with zero infrastructure.
- **Fake at the *port*, test the *adapter* for real.** Fake the repository *interface* when testing a service's logic. Do **not** fake the database to "test" a repository — the SQL is the thing under test, so the repository is an integration test against real Postgres (below). Mocking the DB proves only that your mock returns what you told it to.
- **No globals to mock.** Dependencies arrive by constructor injection ([DI](../module-2-intermediate/dependency-injection)), so a test wires a fake in directly — there's no package-level singleton to monkey-patch.

```go
type fakeStore struct{ insertErr error; inserted []*domain.Policy }
func (f *fakeStore) Insert(_ context.Context, p *domain.Policy) error {
	if f.insertErr != nil { return f.insertErr }
	f.inserted = append(f.inserted, p); return nil
}
// NewPolicyService(&fakeStore{}, ...) — logic under test, no infra.
```

### Integration & repository tests — real Postgres via `dxtest/containers`

Repository fakes prove nothing about SQL. Repository tests run the real repository against **real Postgres**, and the platform way is `dxtest/containers` (in `dx-common-go`) — a testcontainers Postgres, started **once per test binary** (shared via `sync.Once`, cheaper than one container per test) and torn down by testcontainers' reaper when the binary exits:

```go
func TestAccessRequestRepo(t *testing.T) {
	// Provisions the schema through the SAME migration runner main() uses —
	// so the test exercises the exact migration path production takes.
	h := containers.Postgres(t, containers.WithMigrations(svcdb.Migrations, "migrations"))
	repo := NewAccessRequestRepo(h.Pool)

	// exercise real INSERT / SELECT / soft-delete filter / paging against the real schema
	ok, err := repo.PendingExists(ctx, itemID, consumerID)
	// ... assert
}
```

Two escape hatches on `containers.Postgres`, for two purposes:

- **`WithMigrations(fsys, dir)`** — runs your service's versioned migrations before handing back the handle. Use this for repository/service tests: they get the production schema by the production path.
- **`WithSetupSQL(fsys, dir)`** — runs idempotent `.sql` (e.g. `dx-common-go/dxtest/fixtures/schema.sql`) for *library* tests that need a fixture schema, not a service's migrations. Because the container is shared across the binary, setup SQL must be idempotent (`CREATE TABLE IF NOT EXISTS`, no seed `INSERT`s).

**The DSN fallback** is first-class, not a hack: set `DX_TEST_PG_DSN` (or `DX_TEST_REDIS_ADDR`) and `containers.Postgres`/`.Redis` bind to that external instance instead of starting a container — the Docker-less CI path, and the way to point tests at the shared local stack. When **neither** a DSN nor a reachable Docker daemon is present, the helper calls `t.Skip` — never `t.Fail`.

Discipline that makes real-infra testing livable regardless of backing: tests create their own fixtures, clean up after themselves (or use throwaway data/schemas), and **never assume table emptiness** — a shared instance always has someone else's rows. Run your test twice in a row; if the second run fails, your cleanup is wrong.

### Migration testing — the schema is code too

Because `WithMigrations` runs your real migration files, a passing repository test already proves the migrations *apply cleanly from zero*. Make that explicit and add the reverse:

- **From-zero `up`** — provisioning a fresh container with `WithMigrations` is exactly this. A broken `up` fails the test, not the deploy.
- **`down` → `up` round-trip** — apply all, roll the newest down, re-apply. Catches a `down` that doesn't actually reverse its `up` (the least-tested SQL in any codebase — [Schema Migrations](../module-3-advanced/schema-migrations)).
- **Dirty-state is loud** — a migration that fails partway surfaces as `*migrate.DirtyStateError`; a test that expects a bad migration to fail should assert that type, not just "an error."

This is the CI job the platform standardized toward: from-zero up plus a newest down/up cycle on a scratch Postgres, so a broken migration fails the PR.

### API testing — the full router, exactly as production runs it

Handler/API tests mount the **full chi router with its whole middleware stack** via `httptest` — so routing, request-validation middleware, auth resolution, envelopes, and error mapping are all exercised as production runs them, not bypassed by calling a handler function directly:

```go
req := httptest.NewRequest(http.MethodGet, "/items?status=ACTIVE", nil)
rec := httptest.NewRecorder()
newRouter(deps).ServeHTTP(rec, req) // REAL router — middleware included
// assert status, the DxResponse/DxPagedResponse envelope shape, error URNs
```

Per the standards, **every endpoint has handler tests** — a reviewer will look. Assert observable behavior: status codes, the response envelope, error taxonomy (401 without token, 403 without rights, 400 on bad input) — not private plumbing. The two spec tests (spec-loads, spec-covers-all-routes) pin the OpenAPI contract at build time.

### Elasticsearch testing — mock the transport, or point at a cluster

ES code follows the same skip-when-absent rule, at two levels:

- **Unit — no server.** Inject `esclient.Config.Transport` (an `http.RoundTripper`) that returns canned JSON, or point `Addresses` at an `httptest.Server`. The one gotcha: the official client verifies the product, so the fake response **must set `X-Elastic-Product: Elasticsearch`** or `esclient.New` rejects the connection. This tests query construction and hit decoding with zero infrastructure.
- **Integration — real cluster.** Set `ES_TEST_ADDR` to a running cluster; the `repository/integration_test.go` suite runs against it and skips otherwise. (`dxtest/containers` has Postgres and Redis helpers but **no ES helper today** — ES integration uses the env-addr pattern.) See [Elasticsearch](../module-3-advanced/elasticsearch) §Testing.

### Contract — the migration's safety net

While a Go service replaces a Java one, the question is never "is the Go code nice?" but "**do clients notice?**" Contract checks answer it: same request → both stacks → compare responses (shape, fields, semantics — modulo the documented envelope deltas in `CLIENT-CONTRACT-CHANGES.md`). Today this is scripted/manual during porting rather than a permanent suite. If you port a legacy endpoint, a contract comparison belongs in your PR description.

### E2E — `make dev-demo` as executable spec

Read the smoke test's source (it's in the orchestrator repo) and you'll find it *is* the architecture diagrams, asserted: health across the fleet; JWT issuance; gateway enforcement (no token → 401, token without rights → 403, with → 200); proxying to upstreams; and the full **PAP → RMQ → PDP propagation** — create policy, allowed within ~3s; delete, denied within ~3s. That's why "dev-demo green" gates image promotion: it's the minimum claim "the platform works" reduced to a script.

E2E's role is *critical paths only* — it tells you **that** something broke, rarely **what**; the lower levels exist to answer *what*.

### What the platform still doesn't do

Honest inventory: no property-based/fuzz-testing tradition, no load-testing harness in the standard kit, contract checks not yet automated, and no ES testcontainers helper (env-addr instead). As with library gaps, absence isn't an invitation to freelance — it's context, and potentially your contribution.

:::info[Platform connection]
`claude-docs/TESTING.md` is the authoritative version of this page — read it now. The testcontainers machinery is `dx-common-go/dxtest/containers` (read `postgres.go` — the `sync.Once` container + DSN-fallback + skip logic is the whole strategy in one file) and the shared fixture schema is `dxtest/fixtures`. Live examples: any `*_integration_test.go` under a service's `internal/repository/postgres/`, and `database/postgres/repository/repository_integration_test.go` in the library itself. The capstone requires all of: unit tests, a repository integration test via `dxtest/containers`, an API test through the full router, and a dev-demo-style smoke script.
:::

## Exercises

1. Write a repository integration test for `dx-scratch-go` using `containers.Postgres(t, containers.WithMigrations(...))` — cover the soft-delete filter and paging. Run it twice in a row to prove your cleanup; then run it with `DX_TEST_PG_DSN` pointed at the local stack and confirm identical results.
2. Migration test: add a `down`/`up` round-trip test for your newest migration. Then break the `down` on purpose and watch the round-trip catch it.
3. API test: mount the full router and assert the three auth outcomes (401/403/200) and the paginated envelope shape for a list endpoint. Add the spec-covers-all-routes test.
4. Mocking judgment: for a service method that calls a repository and a publisher, write it with hand fakes for *both* ports. Then explain in one sentence why you would **not** fake the DB to test the repository itself.
5. ES unit test: test a search method with a `Transport` mock returning canned hits (remember the product header), asserting the decoded `[]T` and total — no cluster.
6. Write `smoke.sh` for your scratch service: boots it (compose), waits for ready, runs five curl checks including one auth-denial and one write-then-read path, exits nonzero on any failure. Under 50 lines.

## Check yourself

- What does "green for a developer with nothing but the repo" require of every infra-dependent test?
- What does `containers.Postgres(t, WithMigrations(...))` give you that a hand-rolled `pgxpool.New(dsn)` does not?
- When do you fake a dependency, and when do you use the real thing? Why is the database always the real thing for a repository test?
- What header must an Elasticsearch transport-mock response set, and why?
- What exactly does `dev-demo` prove, and what does it deliberately not tell you?

## References

- [Martin Fowler: Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)
- [Testcontainers for Go](https://golang.testcontainers.org/)
- Platform: `claude-docs/TESTING.md` (required); `dx-common-go/dxtest/{containers,fixtures}`; `dx-common-go/FRAMEWORK.md` §6; the `dev-demo` target's source; `CLIENT-CONTRACT-CHANGES.md`
