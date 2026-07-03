---
title: Pointers & Memory Basics
sidebar_label: Pointers & Memory
description: What pointers are in Go, when to use them, and why there's no pointer arithmetic to fear.
---

# Pointers & Memory Basics

## Learning objectives

- Use `&` and `*` confidently; know what `new` does and why you'll rarely write it.
- Decide when a function should take or return a pointer.
- Understand nil-pointer dereferences and how to avoid the common ones.
- Know just enough about stack vs heap to stop worrying (escape analysis does the deciding).

## Prerequisites

- [Structs & Methods](structs-methods)

## Time estimate

**2.5 hours**

## Concepts

### Pointers without the fear

If C pointers scarred you: Go pointers have **no arithmetic**, are **type-safe**, and are **garbage-collected**. A pointer is just "a reference to a value someone else may also be looking at."

```go
x := 10
p := &x        // p is *int, pointing at x
*p = 20        // dereference and write; x is now 20

q := new(int)  // allocates a zeroed int, returns *int
// idiomatic Go rarely uses new; &T{} covers structs
```

### When to use a pointer

| Situation | Pointer? |
|---|---|
| Function must modify the caller's value | Yes |
| Struct is large (copying is wasteful) | Yes |
| Struct contains a mutex, pool, or connection | Yes — copying those is a bug |
| "Absent vs zero" must be distinguishable (e.g. optional JSON field) | Yes (`*string`, `*int`) |
| Small immutable value (`time.Time`, small structs) | No — pass by value |

Slices, maps, channels, and functions are already reference-like — you almost never need `*[]int` or `*map[K]V`.

```go
// Optional fields in an API payload — nil means "not provided"
type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
```

### nil and how it bites

Dereferencing a nil pointer panics. The platform defends with early returns:

```go
func (s *Service) Rename(p *Policy, name string) error {
	if p == nil {
		return errors.New("nil policy")
	}
	p.Name = name
	return nil
}
```

A method with a pointer receiver **can** be called on a nil receiver — the panic happens only when it touches a field. Some types exploit this deliberately (a nil logger that no-ops), but treat that as an advanced pattern, not a default.

### Stack, heap, and escape analysis

You do not choose where a value lives; the **compiler** does. If a value's lifetime provably ends with the function, it goes on the stack (cheap). If it "escapes" — returned by pointer, captured by a closure, stored in a longer-lived struct — it moves to the heap (garbage-collected, slightly costlier).

```go
func makeUser() *User {
	u := User{Name: "ada"} // escapes: we return its address
	return &u              // perfectly legal in Go — no dangling pointer
}
```

Returning a pointer to a local variable is safe and idiomatic. Escape analysis and its performance implications return in [Memory & Performance](../module-2-intermediate/memory-performance); for now the takeaway is: **write for correctness and clarity first** — the compiler's allocation decisions are visible later via `go build -gcflags=-m` when you actually need them.

:::info[Platform connection]
DX constructor functions all return pointers — `postgres.NewPolicyRepo(pool)` returns `*PolicyRepo`, handlers hold `*Service`, and `main.go` wires a graph of pointers at boot. Optional-field pointers (`*string`) appear throughout request DTOs where "field omitted" must differ from "field set to empty". You'll never see pointer arithmetic, and you'll never manually free anything.
:::

## Exercises

1. Write `func swap(a, b *int)` and demonstrate it. Then explain why `func swap(a, b int)` can't work.
2. Write an `UpdateDataset(req UpdateRequest)` function where each non-nil field of the request overwrites the dataset — the pointer-as-optional pattern.
3. Trigger a nil-pointer panic on purpose, read the stack trace top to bottom, then guard it with an early return.
4. Run `go build -gcflags=-m` on a file containing `makeUser` above and find the "escapes to heap" line.

## Check yourself

- Why can Go safely return `&localVariable` when C can't?
- Name three types you should (almost) never take a pointer to.
- How do you represent "this optional field was not provided" in a JSON request struct?

## References

- [A Tour of Go — Pointers](https://go.dev/tour/moretypes/1)
- [Go by Example: Pointers](https://gobyexample.com/pointers)
- [Go FAQ — When are function parameters passed by value?](https://go.dev/doc/faq#pass_by_value)
