---
title: Interfaces
sidebar_label: Interfaces
description: Implicit satisfaction, small interfaces, the comma-ok assertion, and Go's most important design rule.
---

# Interfaces

## Learning objectives

- Define and satisfy interfaces implicitly — no `implements` keyword.
- Apply the platform rules: keep interfaces small, **define them at the consumer**, accept interfaces / return concrete types.
- Use type assertions (comma-ok form only) and type switches safely.
- Explain the "typed nil" trap.

## Prerequisites

- [Structs & Methods](structs-methods), [Pointers & Memory](pointers-memory-basics)

## Time estimate

**4 hours**

## Concepts

### Implicit satisfaction

An interface is a set of method signatures. Any type with those methods satisfies it — **no declaration connects them**:

```go
type Checker interface {
	Check(ctx context.Context) error
}

type PgChecker struct{ pool *pgxpool.Pool }

func (c *PgChecker) Check(ctx context.Context) error {
	return c.pool.Ping(ctx) // *PgChecker is now a Checker. That's it.
}
```

This decouples packages: the interface's author and the implementation's author never need to know about each other.

### Small interfaces win

The stdlib's most-used interfaces have one or two methods (`io.Reader`, `io.Writer`, `fmt.Stringer`, `error`). A one-method interface is trivially satisfiable, mockable, and composable. If your interface has six methods, you've probably defined a class, not an interface.

### Define interfaces at the consumer

This is the rule that most surprises people coming from Java/C#:

> The package that **uses** the dependency defines the interface, sized to exactly what it needs. The implementing package returns a **concrete type** and doesn't know the interface exists.

```go
// package service — the CONSUMER defines what it needs:
type PolicyStore interface {
	Insert(ctx context.Context, p *domain.Policy) error
	FindByID(ctx context.Context, id string) (*domain.Policy, error)
}

type PolicyService struct{ store PolicyStore }

// package repository — the PRODUCER returns a concrete type:
func NewPolicyRepo(pool *pgxpool.Pool) *PolicyRepo { ... }
```

`*PolicyRepo` satisfies `PolicyStore` implicitly. Benefits: the service is testable with a fake store; the repository stays honest (no premature abstraction); and adding a repo method doesn't force every consumer to care. Corollary rule: **accept interfaces, return concrete types.**

### Type assertions — comma-ok, always

```go
var v any = fetch()

s, ok := v.(string) // comma-ok: safe
if !ok { /* handle */ }

s := v.(string)     // single-value form PANICS on mismatch — banned in platform code
```

Type switches handle several possibilities cleanly:

```go
switch x := v.(type) {
case string:
	return x
case int:
	return strconv.Itoa(x)
case fmt.Stringer:
	return x.String()
default:
	return fmt.Sprintf("%v", v)
}
```

### The typed-nil trap

An interface value holds a (type, value) pair. It is `== nil` only when **both** are nil:

```go
func find() *Record { return nil }

var v any = find()
fmt.Println(v == nil) // false! v holds (type=*Record, value=nil)
```

This bites hardest with `error` returns — returning a nil concrete `*MyError` through an `error` interface makes `err != nil` true. Rule: functions returning `error` return literal `nil`, not a nil concrete pointer.

### The empty interface and `any`

`any` (alias for `interface{}`) matches everything and therefore tells you nothing. It's appropriate at true dynamic boundaries (JSON of unknown shape, `map[string]any` constraints) and almost nowhere else — [Generics](generics) replaced most historical uses.

:::info[Platform connection]
The consumer-defined-interface rule **is** the DX architecture in miniature: every service layer defines the store interface it needs, and `internal/repository/postgres` provides the concrete type — which is exactly what makes handler and service tests possible without a database. In `dx-common-go`, `health.Checker` (one method) lets any dependency plug into `/healthz/ready`. And the comma-ok rule is enforced by the platform's style skill — the panicking form won't survive review.
:::

## Exercises

1. Define `type Notifier interface{ Notify(msg string) error }` in a `service` package; implement `EmailNotifier` and `SlackNotifier` in another package (returning concrete types); wire both through the service. No `implements` anywhere — verify with a compile-time assertion: `var _ service.Notifier = (*EmailNotifier)(nil)`.
2. Reproduce the typed-nil trap with an `error` return, observe the wrong behavior, then fix it.
3. Write `func describe(v any) string` with a type switch covering `string`, `bool`, `[]byte`, `fmt.Stringer`, default.
4. Take your Module-1 `Dataset` type and make it satisfy `fmt.Stringer` and `sort.Interface` (over a `[]Dataset` wrapper type).

## Mini-project — pluggable storage (2 h)

Build a tiny key-value CLI (`get`/`set`/`list`) where the command layer defines `type Store interface {...}` and two backends implement it: an in-memory map and a JSON file. Switch backends with a flag. This is Module 2's [dependency injection](../module-2-intermediate/dependency-injection) in embryo.

## Check yourself

- Where should an interface be declared, and who should return concrete types?
- Why is the one-value type assertion banned in platform code?
- Explain how an interface holding a nil pointer is not nil.
- Why are one-method interfaces preferred?

## References

- [A Tour of Go — Interfaces](https://go.dev/tour/methods/9)
- [Effective Go — Interfaces](https://go.dev/doc/effective_go#interfaces)
- [Go FAQ — nil error](https://go.dev/doc/faq#nil_error) — the typed-nil explanation
- [Google Go Style Guide — Interfaces](https://google.github.io/styleguide/go/decisions#interfaces)
