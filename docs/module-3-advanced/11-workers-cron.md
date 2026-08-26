---
title: Workers and Scheduled Jobs
sidebar_label: Workers & Cron
description: Supervision, cancellation, leases, singleton work, retries, and Kubernetes scheduling.
---

# Workers and scheduled jobs

## Worker lifecycle

~~~go
app.Background("audit-consumer", consumer.Run)
app.Go("outbox-dispatcher", func(ctx context.Context) error {
    return dispatcher.Run(ctx, time.Second)
})
~~~

Background restarts an independently recoverable worker with capped backoff. Go treats failure as service failure and coordinates shutdown. Both functions must return on context cancellation.

## Choose a work model

| Work | Pattern |
|---|---|
| Continuous broker delivery | Supervised consumer |
| Durable database queue | Batch claim with FOR UPDATE SKIP LOCKED |
| Periodic idempotent maintenance | Ticker worker or scheduler |
| Independent scheduled execution | Kubernetes CronJob |
| Non-idempotent singleton | Advisory or distributed lock with bounded lease |

Prefer distributed work partitioning over leader election.

## Safe batch claim

1. Start a short transaction.
2. Select a bounded batch FOR UPDATE SKIP LOCKED.
3. Mark or lease the rows.
4. Commit.
5. Perform external work.
6. Record success or schedule retry in another short transaction.

Do not hold database locks during slow network calls.

## Retry and poison work

Persist attempt count, next-attempt time, last safe classification, and terminal state. Backoff with jitter. Make jobs idempotent. Move permanently invalid work to a visible failed state rather than retrying forever.

## Shutdown

Stop acquiring work, cancel blocking receives, finish or release in-flight leases, flush durable progress, then return. Test shutdown while every phase is active.

## Exercise

Build an email worker with a database queue, two concurrent replicas, SKIP LOCKED claim, retry schedule, terminal failure, and graceful cancellation. Prove each row is claimed by only one replica at a time.

## Check yourself

- When should a worker failure stop HTTP?
- Why not hold a row lock while sending email?
- Which jobs need a singleton lock?
- What state makes retries observable?
