---
title: Background Workers & Cron Jobs
sidebar_label: Workers & Cron
description: Supervised background loops, one-shot jobs under Kubernetes CronJob, and advisory locks for singletons.
---

# Background Workers & Cron Jobs

## Learning objectives

- Choose the right shape for background work: embedded supervised loop vs one-shot CronJob binary.
- Supervise workers with errgroup and stop them via context on shutdown.
- Make scheduled jobs idempotent and safe to run twice.
- Guard single-instance work with Postgres advisory locks.

## Prerequisites

- [Concurrency](../module-2-intermediate/concurrency), [Context](../module-2-intermediate/context), [RabbitMQ & Events](event-driven-rabbitmq)

## Time estimate

**3.5 hours**

## Concepts

### Two shapes, one decision rule

| Work | Shape | Examples |
|---|---|---|
| **Continuous** — always on, reacts to a stream | Goroutine *inside* the service, supervised, context-cancelled | outbox dispatcher, RMQ consumer, audit-event drainer |
| **Scheduled** — runs at a time, does a batch, exits | **Separate one-shot binary**, scheduled by Kubernetes CronJob | expiry sweeps, digest emails, cleanup, report generation |

The platform's rule: **scheduled work is a one-shot binary under K8s CronJob**, not a `time.Ticker` inside a service. Reasons: a service with three replicas would run the "daily" job three times; an in-process cron dies and resurrects with its pod invisibly; and CronJob gives you scheduling, retry policy, visibility (`kubectl get jobs`), and resource isolation for free. Don't import a Go cron library to rebuild what the orchestrator already does.

### Continuous workers — supervised, always

A background loop nobody watches is a silent failure: the goroutine dies, the service keeps serving, and the outbox quietly grows for a week. The platform requires supervision — the errgroup shape from [Concurrency](../module-2-intermediate/concurrency), applied at service scale:

```go
g, ctx := errgroup.WithContext(rootCtx)

g.Go(func() error { return outboxDispatcher.Run(ctx) })
g.Go(func() error { return auditConsumer.Run(ctx) })
g.Go(func() error { return httpserver.Run(ctx, cfg, router, logger) })

if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
	logger.Fatal("component failed", zap.Error(err))
}
```

One component's fatal error cancels `ctx`, every sibling winds down, the process exits, and the orchestrator restarts it *visibly* — crash-and-restart beats limp-along-silently. Each `Run(ctx)` is the ticker-plus-`ctx.Done()` select loop you built in Module 2. (The platform review flagged an unsupervised worker in the files service as a defect — the standard exists because the failure mode is real.)

### One-shot jobs — idempotent or dangerous

A CronJob binary is just `main` that does the batch and exits nonzero on failure (Kubernetes handles retries/alerting). The non-negotiable property is **idempotence**, because schedulers guarantee *at least* once, not exactly once — a retried or overlapping run must be harmless:

```go
// Good: the WHERE clause makes re-running a no-op
res, err := pool.Exec(ctx,
	`UPDATE subscriptions SET status = 'expired'
	  WHERE expires_at < now() AND status = 'active'`)
log.Info("expiry sweep", zap.Int64("rows", res.RowsAffected()))
```

State-driven ("expire everything past due that isn't expired yet") is idempotent by construction; action-driven ("send an email per row I selected") needs a sent-marker written in the same transaction as the action's record.

### Advisory locks — at most one, cheaply

Two replicas polling one outbox, or an overlapping cron run, both want the same guarantee: **at most one instance works at a time**. Postgres **advisory locks** are the platform's answer — a lock on an arbitrary ID, no table required, auto-released if the holder's connection dies (so a crashed holder can't wedge the system):

```go
var got bool
err := pool.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, jobLockID).Scan(&got)
if !got {
	log.Info("another instance holds the lock; exiting")
	return nil // clean no-op — this is success, not failure
}
defer pool.Exec(ctx, `SELECT pg_advisory_unlock($1)`, jobLockID)
// ... do the singleton work ...
```

The standards prescribe exactly this for continuous singletons (an outbox dispatcher that must not run twice) and for crons that might overlap their own next run.

:::info[Platform connection]
Live examples of the continuous shape: `dx-acl-go`'s outbox dispatcher and the consumer loops in `dx-audit-go`/`dx-notification-go` — all context-cancelled goroutines started from `main`. The scheduled shape is the platform's stated direction (GO-SERVICE-STANDARDS: "scheduled work = one-shot binary via K8s CronJob"): several legacy crons (competition scheduling, subscription expiry) still live in Python/Java and are queued to be ported into exactly the pattern you just learned — quite possibly by you, in [First Contribution](../capstone/first-contribution).
:::

## Exercises

1. Restructure `dx-scratch-go`'s `main` so HTTP server + outbox dispatcher run under one errgroup with the fail-together semantics above. Kill the dispatcher (inject an error after 10s) and verify the whole process exits.
2. Write `cmd/expire-notes/main.go`: one-shot binary marking old notes archived, idempotent by construction, correct exit codes. Run it thrice; second and third runs touch zero rows.
3. Add `pg_try_advisory_lock` to the one-shot. Run two copies concurrently (`./expire & ./expire`); exactly one works, the other exits cleanly.
4. Write (don't apply) the Kubernetes CronJob YAML for it: schedule, `concurrencyPolicy: Forbid`, `backoffLimit`. Explain what `Forbid` adds on top of your advisory lock — and why you still want both.

## Check yourself

- Why is `time.Ticker`-in-the-service the wrong home for a daily job on this platform?
- What failure mode does supervision prevent, and why is crash-and-restart preferable to it?
- Your cron ran twice due to a node failure — what property makes that a non-event?
- Why do advisory locks self-release, and why does that matter?

## References

- [Kubernetes CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
- [PostgreSQL advisory locks](https://www.postgresql.org/docs/current/explicit-locking.html#ADVISORY-LOCKS)
- [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)
- Platform: GO-SERVICE-STANDARDS.md (workers & jobs); `dx-acl-go` outbox dispatcher
