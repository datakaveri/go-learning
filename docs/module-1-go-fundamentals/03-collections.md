---
title: Collections — Slices, Maps & Strings
sidebar_label: Collections
description: Arrays, slices and their sharing semantics, maps, and how strings really work.
---

# Collections — Slices, Maps & Strings

## Learning objectives

- Explain the array/slice relationship: length, capacity, and the backing array.
- Predict when `append` shares memory and when it reallocates — and why that matters at API boundaries.
- Use maps correctly: comma-ok lookups, deletion, iteration order, and nil-map traps.
- Handle strings as UTF-8: bytes vs runes, `strings.Builder`.

## Prerequisites

- [Control Flow & Functions](control-flow-functions)

## Time estimate

**4 hours**

## Concepts

### Arrays exist; you'll use slices

An array (`[4]int`) has a fixed size that is part of its type. A **slice** (`[]int`) is a lightweight view — pointer, length, capacity — onto a backing array. Slices are what every API uses.

```go
s := []int{1, 2, 3}
s = append(s, 4)

t := s[1:3]      // shares the SAME backing array as s
t[0] = 99        // s is now [1, 99, 3, 4]
```

That sharing is the most important fact on this page. `append` may or may not allocate a new backing array depending on capacity — so two slices can alias each other invisibly:

```go
a := make([]int, 0, 4) // len 0, cap 4
a = append(a, 1, 2)
b := append(a, 3)      // fits in capacity: b shares a's array
c := append(a, 4)      // ALSO fits: c overwrites what b wrote!
_ = b[2] == c[2]       // true — both are 4. Surprise.
```

Practical rules:

- Never keep using a slice after passing it somewhere that might append to it (or vice versa).
- **Copy slices at trust boundaries** — a platform rule. If a struct stores a slice it received, copy it (`s := slices.Clone(in)`), otherwise the caller can mutate your state from a distance.
- Preallocate when you know the size: `make([]T, 0, n)` — this shows up again in [Memory & Performance](../module-2-intermediate/memory-performance).

### nil slices are fine; nil maps are not (for writes)

```go
var s []int
s = append(s, 1)       // fine — append handles nil

var m map[string]int
m["x"] = 1             // PANIC: assignment to entry in nil map
m = make(map[string]int) // must make (or literal) before writing
```

JSON note you'll hit in week one of real work: a nil slice marshals to `null`, an empty slice to `[]`. API handlers usually want `[]` — initialize accordingly.

### Maps

```go
ages := map[string]int{"ada": 36}

v, ok := ages["grace"]   // comma-ok: ok == false, v == 0
if !ok { /* absent */ }

delete(ages, "ada")

for k, v := range ages {} // ORDER IS RANDOMIZED — deliberately
```

Iteration order is randomized per run so you can't accidentally depend on it. If you need order, collect keys and sort them. Maps are not safe for concurrent writes — that's a Module 2 topic ([Concurrency](../module-2-intermediate/concurrency)).

### Strings, bytes, runes

A `string` is an immutable sequence of **bytes**, conventionally UTF-8. Indexing gives bytes; `range` decodes runes:

```go
s := "नमस्ते"
len(s)                  // 18 — bytes, not characters!
for i, r := range s {   // r is a rune (code point), i a byte offset
	fmt.Printf("%d:%c ", i, r)
}
```

Building strings in a loop? Concatenation allocates every time; use `strings.Builder`:

```go
var sb strings.Builder
for _, part := range parts {
	sb.WriteString(part)
}
result := sb.String()
```

:::info[Platform connection]
Slices and maps carry every request through a DX service: query filters are maps, result sets are slices of domain structs, and the generic pagination envelope `Page[T]` in `dx-common-go` wraps a `[]T`. The "copy at trust boundaries" rule is a review comment you *will* receive if you store a caller's slice — it's in the platform's style skill verbatim.
:::

## Exercises

1. Write `func dedupe(xs []string) []string` preserving first-seen order (use a map as the seen-set). Do not mutate the input.
2. Reproduce the aliasing surprise above, then fix it with `slices.Clone`. Explain in a comment exactly which `append` reallocated.
3. Write `func wordFreq(text string) map[string]int`, then print the results sorted by descending count (you'll need to extract and sort a key slice).
4. Write `func truncate(s string, max int) string` that truncates to at most `max` *runes* without splitting a character, appending "…" when truncated.

## Mini-project — log field extractor

Write a small program that reads lines like `level=error service=dx-acl msg="db down"` from stdin, parses each into a `map[string]string`, collects all values seen per key, and prints a summary. It's a warm-up for structured logging in Module 2, and it exercises slices, maps, and string handling together. (~1.5 h)

## Check yourself

- What three fields make up a slice header?
- Two slices share a backing array; when does writing through one become invisible to the other?
- Why is map iteration order randomized?
- `len("héllo")` — why might it not be 5?

## References

- [Go Blog: Go Slices — usage and internals](https://go.dev/blog/slices-intro) — required reading
- [Go Blog: Strings, bytes, runes and characters](https://go.dev/blog/strings)
- [Go by Example: Slices](https://gobyexample.com/slices), [Maps](https://gobyexample.com/maps)
- [Uber Go Style Guide — Copy Slices and Maps at Boundaries](https://github.com/uber-go/guide/blob/master/style.md#copy-slices-and-maps-at-boundaries)
