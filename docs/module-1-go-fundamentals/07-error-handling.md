---
title: Error Handling
sidebar_label: Error Handling
description: Errors as values — wrapping with %w, errors.Is/As, sentinel vs typed errors, and the platform's non-negotiable rules.
---

# Error Handling

## Learning objectives

- Treat errors as values: create, return, wrap, and inspect them.
- Wrap with `fmt.Errorf("context: %w", err)` and match with `errors.Is` / `errors.As`.
- Choose between sentinel errors and typed errors.
- Apply the platform's two hard rules: **wrap with context** and **log or return, never both**.
- Know what panic is for (almost nothing).

## Prerequisites

- [Interfaces](interfaces) — `error` is just an interface.

## Time estimate

**4 hours** — this page matters more than any other in Module 1.

## Concepts

### error is an interface

```go
type error interface {
	Error() string
}
```

No exceptions, no try/catch. A function that can fail returns `error` as its last value, and the caller decides — immediately, explicitly — what to do:

```go
policy, err := repo.FindByID(ctx, id)
if err != nil {
	return nil, fmt.Errorf("find policy %s: %w", id, err)
}
```

That `if err != nil` block is not noise — it is the control-flow of failure, visible at every call site. Embrace it.

### Wrapping with %w — the platform's first rule

`fmt.Errorf` with the `%w` verb wraps the original error inside a new one, adding context while preserving the cause:

```go
if err := json.Unmarshal(body, &req); err != nil {
	return fmt.Errorf("decode create-policy request: %w", err)
}
```

Every hop up the stack adds one clause of context. The final log line reads like a story:

```
handle create policy: decode create-policy request: unexpected end of JSON input
```

Rules for the message: lower-case, no trailing punctuation, describe **what you were doing**, don't repeat what the caller already knows.

### Inspecting: errors.Is and errors.As

Because errors are wrapped in chains, never compare with `==` or match on strings. Use the two chain-walking helpers:

```go
// errors.Is — "is this (or anything it wraps) that sentinel?"
if errors.Is(err, pgx.ErrNoRows) {
	return nil, ErrPolicyNotFound
}

// errors.As — "is there a *ValidationError anywhere in the chain? give it to me"
var ve *ValidationError
if errors.As(err, &ve) {
	fmt.Println(ve.Field)
}
```

### Sentinel errors vs typed errors

**Sentinel** — a package-level value, for conditions callers test for:

```go
var ErrNotFound = errors.New("policy not found")
```

**Typed** — a struct implementing `error`, for errors carrying data:

```go
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Reason)
}
```

If your custom error wraps another, implement `Unwrap() error` so `errors.Is/As` can keep walking the chain.

### Log or return — never both

The platform's second hard rule. If you log an error *and* return it, every layer above does the same and one failure produces five stack-flavored log lines. Decide at each level:

- **Can't handle it here?** Wrap and return. No log.
- **This is the top (handler, worker loop)?** Log once, with full context, and translate to a response.

```go
// service layer — wrap and return
if err := s.store.Insert(ctx, p); err != nil {
	return fmt.Errorf("insert policy: %w", err)
}

// handler layer — the top: log once, respond once
if err := s.CreatePolicy(ctx, req); err != nil {
	logger.Error("create policy failed", zap.Error(err))
	dxerrors.WriteError(w, err)
	return
}
```

### panic — reserved for the impossible

`panic` is for programmer errors that make continuing meaningless (corrupted invariants, impossible states), not for expected failures. You will essentially never write `panic` in service code; you *will* rely on the HTTP middleware `Recoverer` converting a panicking handler into a 500 instead of a crashed process. `recover` belongs in framework code, not yours.

:::info[Platform connection]
`dx-common-go/errors` defines the platform error taxonomy — `NewValidation`, `NewNotFound`, `NewConflict`, etc. — and each maps to an HTTP status and an IUDX problem URN (like `urn:dx:as:InvalidParamValue`) in the response body. Your service code wraps low-level errors upward until the handler translates them via `dxerrors.WriteError`. Module 3's [REST API page](../module-3-advanced/rest-api-development) covers the taxonomy; the wrapping discipline you learn here is what makes it work.
:::

## Exercises

1. Write `readConfig(path string)` that can fail three ways (missing file, bad YAML, invalid value), wrapping each with context. From `main`, print the full chain and test `errors.Is(err, os.ErrNotExist)`.
2. Build `*ValidationError` with `Unwrap()`; wrap it two levels deep; retrieve it with `errors.As` and print its fields.
3. Find the bug: a function that does `return nil, err` where `err` was declared as `var err *ValidationError` and never assigned — connect this to the typed-nil trap from [Interfaces](interfaces).
4. Take a small program that logs errors at three layers and refactor it to log-once-at-the-top. Compare the log output before and after a simulated failure.

## Check yourself

- What does `%w` do that `%v` doesn't?
- When do you use `errors.Is` vs `errors.As`?
- Why is "log and return" banned?
- Where is `panic` acceptable?

## References

- [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) — required reading
- [Effective Go — Errors](https://go.dev/doc/effective_go#errors)
- [Uber Go Style Guide — Errors](https://github.com/uber-go/guide/blob/master/style.md#errors)
- Platform: `dx-common-go/errors` package; `claude-docs/GO-SERVICE-STANDARDS.md` (error conventions)
