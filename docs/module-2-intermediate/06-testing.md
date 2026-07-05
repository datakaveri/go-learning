---
title: Testing
sidebar_label: Testing
description: go test, table-driven tests, subtests, httptest, fakes — the way the platform actually tests.
---

# Testing

## Learning objectives

- Write and run tests with the built-in toolchain — no framework.
- Master the platform's default shape: **table-driven tests with subtests**.
- Test HTTP handlers and full routers with `httptest`.
- Substitute dependencies with hand-written fakes behind consumer-defined interfaces.
- Measure coverage honestly and know what the CI gate requires.

## Prerequisites

- [Interfaces](../module-1-go-fundamentals/interfaces), [Dependency Injection](dependency-injection)

## Time estimate

**5 hours**

## Concepts

### The toolchain is the framework

A test is a function `TestXxx(t *testing.T)` in a `_test.go` file:

```go
func TestParsePort(t *testing.T) {
	got, err := parsePort("8080")
	if err != nil {
		t.Fatalf("parsePort: unexpected error: %v", err)
	}
	if got != 8080 {
		t.Errorf("parsePort(\"8080\") = %d, want 8080", got)
	}
}
```

`t.Errorf` records a failure and continues; `t.Fatalf` stops this test. Run with `go test ./...`, always with the race detector in CI: `go test -race ./...`.

**Platform note:** DX code uses the standard library — **no testify**. Comparisons are `!=`, `reflect`-free where possible, and `cmp.Diff` from go-cmp for structs. Assertion-library fluency you bring from other ecosystems translates to plain `if` statements here.

### Table-driven tests — the house style

One test body, many cases, each a named subtest:

```go
func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{name: "negative", n: -5, want: "negative"},
		{name: "zero", n: 0, want: "zero"},
		{name: "positive", n: 3, want: "positive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classify(tt.n); got != tt.want {
				t.Errorf("classify(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
```

Why it's the standard: adding a case is one struct literal; failures name the exact case; `go test -run TestClassify/zero` reruns one case. Include error cases in the table (`wantErr bool` or a sentinel to `errors.Is` against).

### Testing HTTP without a network

`net/http/httptest` runs handlers in-process:

```go
func TestHealthLive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()

	newRouter(deps).ServeHTTP(rec, req) // the REAL router, middleware included

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
```

The platform norm is to mount the **full chi router with its middleware stack** in tests — so tests exercise routing, validation, and middleware exactly as production does. Handler tests per endpoint are a GO-SERVICE-STANDARDS requirement.

### Fakes over mocks

Because services depend on small consumer-defined interfaces, a test double is a few lines — no mocking framework, no codegen:

```go
type fakeStore struct {
	insertErr error
	inserted  []*domain.Policy
}

func (f *fakeStore) Insert(_ context.Context, p *domain.Policy) error {
	if f.insertErr != nil {
		return f.insertErr
	}
	f.inserted = append(f.inserted, p)
	return nil
}
```

Now `NewPolicyService(&fakeStore{}, ...)` tests business logic with zero infrastructure. For the repository layer itself, unit fakes prove nothing — that layer is covered by **integration tests against real Postgres**, spun up with `dxtest/containers` (testcontainers, with a DSN fallback and an automatic skip when Docker is absent). That's Module 4's [Testing Strategy](../module-4-platform/testing-strategy) topic.

### Coverage — a flashlight, not a target

```bash
go test -cover ./...
go test -coverprofile=c.out ./... && go tool cover -html=c.out
```

Use the HTML view to find *untested branches that matter* (error paths, edge cases). The platform gates on tests passing (`go test ./...` in the PR gate) and on new endpoints having tests — not on a numeric coverage bar.

:::info[Platform connection]
Real examples worth reading now: DX services test that the **embedded OpenAPI spec parses** (`TestEmbeddedSpecLoads`) and that **every registered route appears in the spec** (`TestSpecCoversAllRoutes`) — cheap tests that pin the API contract and catch drift at build time. Note the philosophy: test observable behavior (routes, statuses, envelopes), not private plumbing.
:::

## Exercises

1. Rewrite three assert-style tests (imagine testify) as plain stdlib tests. Feel the difference; it's smaller than you think.
2. Write a table-driven test for your `parsePort`-style function with at least two error cases matched via `errors.Is`. Run a single case by name with `-run`.
3. Test your logging middleware from [Logging](logging) with `httptest`: assert status passthrough and that the request-ID header appears.
4. Extract an interface from a dependency in your file-hasher, write a fake, and test the orchestration logic without touching the filesystem.
5. Generate the HTML coverage report for your project and write down the two most *important* uncovered branches (not the two easiest).

## Check yourself

- When `t.Fatalf` vs `t.Errorf`?
- Why does the platform prefer full-router tests over calling handler functions directly?
- Why fakes instead of a mocking framework, and what makes them cheap here?
- What do the spec-coverage tests protect against?

## References

- [testing package](https://pkg.go.dev/testing) · [httptest](https://pkg.go.dev/net/http/httptest)
- [Go Wiki: Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp)
- Platform: `claude-docs/TESTING.md`; GO-SERVICE-STANDARDS.md (testing gates)
