---
title: Memory Management & Performance
sidebar_label: Memory & Performance
description: The GC, escape analysis, allocation patterns — and the short list of optimizations that actually matter.
---

# Memory Management & Performance

## Learning objectives

- Describe Go's garbage collector at working-knowledge level and what "GC pressure" means.
- Use escape analysis output to explain where allocations come from.
- Apply the handful of high-value patterns: preallocation, `strings.Builder`, avoiding accidental copies and conversions.
- Recognize premature optimization and decline it.

## Prerequisites

- [Pointers & Memory Basics](../module-1-go-fundamentals/pointers-memory-basics), [Benchmarking & Profiling](benchmarking-profiling)

## Time estimate

**3 hours**

## Concepts

### The GC in three sentences

Go uses a concurrent mark-and-sweep collector that runs alongside your program with sub-millisecond pauses. You don't manage memory — you manage **allocation rate**: the more garbage you create per request, the more CPU the GC spends keeping up. Practically all Go performance work in services reduces to *allocate less on hot paths*.

### Escape analysis, revisited

From Module 1: the compiler decides stack (cheap, freed on return) vs heap (GC-managed). See its reasoning:

```bash
go build -gcflags='-m' ./... 2>&1 | grep escape
```

Common escape causes you can actually influence:

- Returning pointers to locals (fine — but know it allocates).
- Storing values into interfaces (`any` boxing) — every `zap.Any`, every `[]any` element.
- Capturing variables in closures that outlive the call.
- `fmt.Sprintf` and friends (interface parameters ⇒ boxing ⇒ allocations).

### The short list of patterns that matter

**1. Preallocate when the size is known** — the single most common review-level optimization:

```go
out := make([]Item, 0, len(rows))   // one allocation
ids := make(map[string]struct{}, n) // sized map
```

**2. Build strings with `strings.Builder`**, never `+=` in a loop (each `+=` copies everything so far).

**3. Reuse buffers on hot paths** — `sync.Pool` for per-request scratch space (bytes.Buffer, encoders). Measure first; pools add complexity and only pay off under real load.

**4. Mind hidden copies** — `[]byte(s)` and `string(b)` each copy; a range over a slice of large structs copies each element (`for i := range xs { use(&xs[i]) }` avoids it).

**5. Don't hold what you don't need** — a tiny slice sliced out of a huge buffer pins the whole backing array; `slices.Clone` the piece you keep.

### What NOT to do

- Don't sprinkle `sync.Pool`, `unsafe`, or hand-rolled memory tricks through business logic. Services are I/O-bound; the database round-trip is 1000× your struct copy.
- Don't fight the GC with giant object reuse schemes — modern Go GC is very good; your complexity budget is better spent elsewhere.
- Don't accept *any* performance claim (including your own) without a benchstat diff — the [previous page's](benchmarking-profiling) discipline.

The priority order for a slow DX endpoint is almost always: **query shape (indexes, N+1) → payload size → serialization → allocations**. Memory tricks are the last stop, not the first.

:::info[Platform connection]
The platform's style rules bake in the cheap wins — preallocation and `strings.Builder` are review comments, not optional flourishes — and zap itself is the poster child of allocation-aware design (typed fields exist precisely to avoid `fmt.Sprintf` boxing on every log line). But note where the platform *doesn't* micro-optimize: services are database- and network-bound, and the standards spend far more words on query parameterization and connection pooling than on allocation tuning. Learn the proportions along with the techniques.
:::

## Exercises

1. Run escape analysis on a file containing: a function returning `*User`, a closure capturing a local, and a `fmt.Println(x)` — map each `escapes to heap` line to its cause.
2. Fix the classic: a loop building a CSV row with `+=`. Benchmark before/after with `-benchmem` and report the allocs/op change.
3. Demonstrate backing-array pinning: read a 10 MB file, keep a 10-byte slice of it, and compare `runtime.MemStats.HeapAlloc` with and without `slices.Clone`.
4. Take the worst allocation site your profiler found in the file-hasher (from the previous page) and either fix it with a pattern above — with benchstat proof — or write the one-sentence justification for leaving it.

## Check yourself

- What does the GC cost you, and what's the one lever you control?
- Why does storing an `int` in an `any` allocate?
- Which two patterns from the short list would you check for in any PR touching a hot loop?
- Your endpoint is slow — recite the priority order before you reach for `sync.Pool`.

## References

- [A Guide to the Go Garbage Collector](https://tip.golang.org/doc/gc-guide) — skim now, return later
- [Go Wiki: CompilerOptimizations](https://go.dev/wiki/CompilerOptimizations)
- [Uber Go Style Guide — Performance](https://github.com/uber-go/guide/blob/master/style.md#performance)
