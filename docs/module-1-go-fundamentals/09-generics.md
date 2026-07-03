---
title: Generics
sidebar_label: Generics
description: Type parameters and constraints — and how the platform uses them for the DAO and response envelopes.
---

# Generics

## Learning objectives

- Write generic functions and types with type parameters and constraints.
- Use `any`, `comparable`, and interface-based constraints (including union constraints).
- Know when generics help and when an interface (or plain code) is better.
- Read the platform's flagship generic types: `BaseDAO[T]`, `Page[T]`, `DxResponse[T]`, `LoadService[T]`.

## Prerequisites

- [Interfaces](interfaces), [Packages & Modules](packages-modules)

## Time estimate

**3.5 hours**

## Concepts

### Type parameters

```go
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

names := Map(policies, func(p Policy) string { return p.ID })
// type inference: no need to write Map[Policy, string](...)
```

`[T any]` declares a type parameter; the compiler generates and type-checks concrete instantiations. Unlike `any`-typed code there's no runtime assertion, no reflection, no chance of a type error surviving to production.

### Constraints

A constraint is an interface that limits what `T` can be — and therefore what your code may do with it:

```go
// comparable: anything usable with == (map keys, dedupe)
func Dedupe[T comparable](xs []T) []T { ... }

// union constraint: a set of concrete types
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // legal because every member of Number supports +
	}
	return total
}
```

The `~int` tilde means "any type whose underlying type is int" — so your `type Credits int` still satisfies `Number`.

### Generic types

```go
type Page[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
}

func NewPage[T any](items []T, total, limit, offset int) Page[T] { ... }
```

One definition, every entity: `Page[Policy]`, `Page[Dataset]`, `Page[AuditRecord]` — all statically typed, no casting at the call site.

### When (not) to use generics

Generics shine when the **algorithm is identical and only the type varies**: containers, slices/maps helpers, envelopes, data access plumbing. They are the wrong tool when behavior varies by type — that's an interface's job. Decision shortcut:

- Same code, many element types → **generics**.
- Different implementations behind one contract → **interface**.
- Used once, with one type → **neither**; write the plain version.

Don't pre-generalize. The platform introduced generics exactly where duplication hurt (DAO, envelopes, config loading) and nowhere else — there is deliberately **no** generic service layer.

:::info[Platform connection]
The four generic types you'll touch weekly, all in `dx-common-go`:

- **`BaseDAO[T]`** (`database/postgres/dao/base.go`) — CRUD, pagination, and transactional variants for any row-mapped struct `T`. Write a domain struct with `db` tags, get `FindByID`, `FindPage`, upsert and soft-delete support for free. Module 3's [Database Patterns](../module-3-advanced/database-patterns) is a deep dive.
- **`Page[T]`** — the pagination envelope above, for real.
- **`DxResponse[T]` / `DxPagedResponse[T]`** (`response/model.go`) — the platform's standard success envelope.
- **`config.LoadService[T]`** — loads YAML + env into *your service's* config struct type. One loader, every service.
:::

## Exercises

1. Write `Filter[T any](xs []T, keep func(T) bool) []T` and `Keys[K comparable, V any](m map[K]V) []K`.
2. Write the `Number` constraint with `~`, then `Min[T Number](xs []T) (T, error)` — returning an error for an empty slice (tie-in to [Error Handling](error-handling)).
3. Build your own `Page[T]` with a `Map` method converting `Page[T]` to `Page[U]`… and discover why methods can't introduce new type parameters. Write it as a free function instead — and remember the lesson.
4. Implement a tiny generic in-memory `Repo[T any]` with `Save(id string, v T)`, `Get(id string) (T, bool)`, `All() []T` backed by a map — a toy version of what `BaseDAO[T]` does against Postgres.

## Check yourself

- What does the `~` in a constraint mean?
- Why does `Sum[T Number]` compile but a version with `[T any]` not?
- Interface or generics: five storage backends with different behavior? A `Reverse` that works on any slice?
- Which platform types are generic, and why those?

## References

- [Go generics tutorial](https://go.dev/doc/tutorial/generics)
- [Go Blog: An Introduction to Generics](https://go.dev/blog/intro-generics)
- [Go Blog: When To Use Generics](https://go.dev/blog/when-generics) — required reading
- Platform: `dx-common-go/database/postgres/dao/base.go`, `dx-common-go/response/model.go`
