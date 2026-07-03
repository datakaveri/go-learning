---
title: Concurrency
sidebar_label: Concurrency
description: Goroutines, channels, select, sync — worker pools, errgroup, and the no-fire-and-forget rule.
---

# Concurrency

## Learning objectives

- Start goroutines responsibly: every goroutine has a known exit condition and someone waiting for it.
- Communicate with channels; coordinate with `select`; know the platform's buffered-channel policy.
- Protect shared state with `sync.Mutex` and know when a mutex beats a channel.
- Build the two workhorse patterns: worker pool and `errgroup` fan-out.
- Detect races with `-race`.

## Prerequisites

- Module 1 complete — especially [Interfaces](../module-1-go-fundamentals/interfaces) and [Error Handling](../module-1-go-fundamentals/error-handling).

## Time estimate

**6 hours** — the heart of Module 2. Do all the exercises.

## Concepts

### Goroutines — cheap, but never free

```go
go process(job) // runs concurrently; costs ~kilobytes
```

A goroutine is cheap to start and impossible to stop *from outside* — it must return on its own. That leads to the platform's first concurrency rule:

> **No fire-and-forget goroutines.** Every `go` statement must answer: *how does this goroutine exit, and who notices if it fails?*

The three legitimate answers: it's bounded by a `sync.WaitGroup`/`errgroup` the caller waits on; it exits when a channel closes; or it exits when its `context` is cancelled ([next page](context)). A `go` statement with none of these is a leak — and a review rejection.

### Channels

A channel is a typed conduit; sends and receives synchronize goroutines:

```go
results := make(chan Result)    // unbuffered: send blocks until receive

go func() {
	results <- compute()        // hands off directly
}()
r := <-results
```

Closing a channel signals "no more values" — `for v := range ch` exits when `ch` closes. **Only the sender closes**; closing twice or sending on a closed channel panics.

Platform buffer policy (from the Uber guide): **buffer size is 0 or 1**. An unbuffered channel is a synchronization point you can reason about; a buffer of 1 is a hand-off slot. Anything larger must be justified in a comment — big buffers usually hide backpressure problems until production finds them.

### select

`select` waits on multiple channel operations — the enabler of timeouts and cancellation:

```go
select {
case job := <-jobs:
	handle(job)
case <-ctx.Done():
	return ctx.Err()
}
```

### Mutexes — for state, channels — for handoff

```go
type Cache struct {
	mu   sync.Mutex // note: value, not *sync.Mutex; zero value works
	data map[string]string
}

func (c *Cache) Get(k string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[k]
	return v, ok
}
```

Rule of thumb: **channels move ownership of data between goroutines; mutexes guard data that stays put.** A map read/written from multiple goroutines without a lock is a crash waiting to happen — literally; the runtime detects concurrent map writes and aborts.

Always develop and test with the race detector: `go test -race ./...` — it turns "works on my machine" data races into hard failures.

### Worker pool

Bounded concurrency — N workers draining one job channel:

```go
func pool(ctx context.Context, jobs <-chan Job, workers int) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return nil // channel closed: done
					}
					if err := handle(ctx, job); err != nil {
						return fmt.Errorf("job %s: %w", job.ID, err)
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}
	return g.Wait() // first error cancels ctx → all workers wind down
}
```

`golang.org/x/sync/errgroup` is the platform's preferred wrapper over raw `WaitGroup`: it propagates the first error and cancels the shared context so siblings stop early. Note how every exit path is explicit — that's the no-fire-and-forget rule in action.

:::info[Platform connection]
DX services embed exactly these shapes: `dx-acl-go`'s outbox dispatcher is a single supervised goroutine polling on an interval and exiting on context cancellation; `dx-common-go`'s `AuditWorker` drains a small buffered channel of audit events in a background goroutine; RabbitMQ consumers loop over a delivery channel and wind down when it closes. The GO-SERVICE-STANDARDS checklist requires every background loop to be *supervised* (errgroup-style) and *context-cancelled on shutdown* — you now know precisely what both words mean.
:::

## Exercises

1. Start 10 goroutines printing their index; make `main` wait correctly with a `WaitGroup`. Then break it (capture the loop variable wrong pre-Go-1.22 style, forget `Add`) and observe.
2. Build a pipeline: generator → squarer → printer, connected by channels, terminated by closing. Then kill it early with a `done` channel and verify no goroutine leaks (`runtime.NumGoroutine()`).
3. Write the concurrent-map crash, then fix it twice: once with a mutex, once by giving one goroutine sole ownership and a channel of updates. Run both under `-race`.
4. Implement the worker pool above for real work (e.g. fetching N URLs with per-job timeout) and make one job fail — verify errgroup cancels the rest.

## Mini-project — concurrent file hasher (3 h)

A CLI `hashall <dir>` that walks a directory tree, hashes every file with SHA-256 using a worker pool sized to `runtime.NumCPU()`, and prints `path → hash` plus a summary (files, bytes, elapsed). Requirements: clean shutdown on Ctrl-C (you'll wire real signal handling after the [Context](context) page), `-race`-clean, buffer sizes 0 or 1 with any exception justified in a comment. Keep it — you'll profile it in [Benchmarking & Profiling](benchmarking-profiling).

## Check yourself

- The three legitimate exit strategies for a goroutine?
- Why is a buffer of 500 a design smell?
- Channel or mutex: a counter incremented by 20 goroutines? A job queue feeding 5 workers?
- What does errgroup add over WaitGroup?

## References

- [A Tour of Go — Concurrency](https://go.dev/tour/concurrency/1)
- [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go Blog: Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines) — required reading
- [errgroup docs](https://pkg.go.dev/golang.org/x/sync/errgroup)
- [Uber Go Style Guide — Channel Size is One or None](https://github.com/uber-go/guide/blob/master/style.md#channel-size-is-one-or-none)
