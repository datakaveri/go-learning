---
title: Benchmarking & Profiling
sidebar_label: Benchmarking & Profiling
description: testing.B, pprof, CPU and heap profiles — measuring before optimizing.
---

# Benchmarking & Profiling

## Learning objectives

- Write benchmarks with `testing.B` and read their output (ns/op, B/op, allocs/op).
- Capture and read CPU and heap profiles with `pprof`, locally and from a running service.
- Follow the discipline: **measure → change → measure again**; never optimize from intuition.

## Prerequisites

- [Testing](testing), [Concurrency](concurrency)

## Time estimate

**3 hours**

## Concepts

### Benchmarks

A benchmark is `BenchmarkXxx(b *testing.B)`; the framework calibrates `b.N`:

```go
func BenchmarkBuildKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildKey("dx", "policy", "1234")
	}
}
```

```bash
go test -bench=. -benchmem ./...
```

```
BenchmarkBuildKey-10    12094743    98.1 ns/op    48 B/op    2 allocs/op
```

Read right to left: **allocs/op** is often the story — allocation count drives GC pressure, and reducing it is the most common real-world Go optimization. Compare implementations by benchmarking both and diffing with `benchstat` (which also tells you whether the difference is statistically meaningful — single runs lie).

Benchmark hygiene: reset the timer after expensive setup (`b.ResetTimer()`), don't let the compiler optimize your work away (assign results to a package-level sink), and run multiple times (`-count=6`) for benchstat.

### Profiling with pprof

Benchmarks tell you *how fast*; profiles tell you *where the time goes*.

```bash
go test -bench=BenchmarkBuildKey -cpuprofile=cpu.out -memprofile=mem.out
go tool pprof cpu.out
(pprof) top10        # heaviest functions
(pprof) list buildKey # annotated source, cost per line
(pprof) web          # call graph (needs graphviz)
```

The two profiles you'll actually use:

- **CPU profile** — samples the call stack ~100/s; answers "what is burning CPU?"
- **Heap profile** — records allocation sites; answers "what is creating garbage?" (`-sample_index=alloc_space` for cumulative allocations, the usual view for GC pressure).

Flame graphs (`go tool pprof -http=:8081 cpu.out`) make both readable at a glance: wide boxes = expensive; look for wide boxes you didn't expect.

### Profiling a live service

Import `net/http/pprof` and your HTTP server grows `/debug/pprof/` endpoints; grab 30 seconds of CPU from a running process:

```bash
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
go tool pprof http://localhost:8080/debug/pprof/heap
```

Also there: `goroutine` (find leaks — remember `runtime.NumGoroutine()` from the concurrency exercises), `block`, and `mutex` profiles for contention. In production, gate these endpoints away from public traffic.

### The discipline

1. Have a hypothesis and a **measurement** before touching code.
2. Change one thing.
3. Measure again, with benchstat, and keep the receipt (paste the diff in the PR).

Most "optimizations" made without step 1 make code worse and performance unchanged. The next page ([Memory & Performance](memory-performance)) gives you the vocabulary for what you'll find in profiles.

:::info[Platform connection]
DX services expose Prometheus `/metrics` for continuous observation (Module 3's [Observability](../module-3-advanced/observability)) and use pprof for incident-time deep dives. When a service misbehaves under load — a hot JSON path, an N+1 query loop, a goroutine leak in a consumer — the workflow above (grab profile from the running container, `top10`, `list`) is exactly what you'll do. Practicing it on your own code first means the first real incident isn't also your first profile.
:::

## Exercises

1. Benchmark string concatenation vs `strings.Builder` vs `fmt.Sprintf` for building a 100-part string. Explain the allocs/op differences.
2. Benchmark `append` into a zero-cap slice vs a preallocated one for 10k elements; connect the result to the slice-growth mechanics from [Collections](../module-1-go-fundamentals/collections).
3. Profile your file-hasher mini-project under a large directory: find the widest flame-graph box, decide whether it's fixable (I/O usually isn't) — and write one sentence you'd put in a PR justifying either the fix or leaving it alone.
4. Add `net/http/pprof` to any toy HTTP server, generate load with `hey` or a loop of `curl`, and capture a live 15-second CPU profile.

## Check yourself

- What do ns/op, B/op, allocs/op each tell you, and which most often points at the fix?
- Why is a single benchmark run before/after a change insufficient evidence?
- Which profile for: high CPU? rising memory? a stuck service with 40k goroutines?

## References

- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [pprof docs](https://pkg.go.dev/net/http/pprof) · [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Dave Cheney: High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
