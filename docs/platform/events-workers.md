---
title: Events, outbox, consumers, workers, and schedules
description: Reliable asynchronous work with platform events, AMQP, leases, retries, and shutdown.
---

# Events, outbox, consumers, workers, and schedules

**Status: Implemented shared primitives; fleet adoption is partial.** Delivery is at least once. Exactly-once business effects require idempotency at the side-effect owner.

![Transactional outbox and event consumption](/img/architecture/outbox-events.svg)

## Domain facts and commands

Publish a past-tense domain fact after an owned state transition, for example `policy.created.v1`. A command asks one owner to do work and needs explicit destination/timeout semantics. Do not disguise a synchronous dependency as an event, or use a fact name that instructs every consumer to perform the same business side effect.

## Typed events

```go
type WidgetCreated struct {
    WidgetID string `json:"widgetId"`
    OwnerID  string `json:"ownerId"`
}

var widgetCreated = events.NewTopic[WidgetCreated]("widget.created").V(1)
```

The platform envelope carries stable event ID, type, version, occurrence time, correlation ID, and payload. The producer owns the schema and compatibility. Additive fields are preferred; incompatible meaning requires a new version and transition plan.

## Producer flow

1. Validate and authorize the command.
2. Begin `sql.Manager.Do`.
3. Change domain state.
4. Write the typed event through `events.Publish` to `events.Outbox` in the same context.
5. Commit.
6. Call `Dispatcher.Kick`; the supervised dispatcher publishes with confirms.
7. Mark sent only while holding the outbox claim lease.

The dispatcher runs with `app.Background("outbox", func(ctx) error { return dispatcher.Run(ctx, interval) })`. Broker outage delays delivery without losing the committed fact. Alert on oldest pending age and repeated publish failures.

## Consumer flow

```go title="Verified API shape; business handler simplified"
err := events.Subscribe(ctx, bus, widgetCreated, func(ctx context.Context, event events.Event, body WidgetCreated) error {
    if seen.AlreadyProcessed(ctx, event.ID) {
        return nil
    }
    return tx.Do(ctx, func(txCtx context.Context) error {
        if err := projection.Upsert(txCtx, body); err != nil {
            return classify(err)
        }
        return seen.MarkProcessed(txCtx, event.ID)
    })
})
```

The exact idempotency adapter is service-owned; `platform/idempotency` supplies durable mechanics for side-effect owners. `events.ErrDrop` classifies an intentionally discarded message; `events.ErrQuarantine` sends a poison message to quarantine/dead-letter handling. Other transient errors are retried under bounded AMQP policy.

## Retry, DLQ, replay, and reconciliation

- Acknowledge only after the local transaction/effect completes.
- Retry transient dependency errors with bounded exponential backoff and jitter.
- Quarantine invalid or unsupported messages with reason, event ID, schema version, and safe metadata.
- Dead-letter after the configured attempt budget.
- Replay preserves the original event ID so consumer idempotency still works.
- Reconciliation compares authoritative state to the projection and repairs missing/divergent entries. Replay is not a substitute for reconciliation.
- Consumer reconnect must be supervised and observable; a broker restart must not silently stop consumption.

## Workers and schedules

![Worker lifecycle](/img/architecture/worker-lifecycle.svg)

- Use `App.Go` for an essential loop whose failure should stop the service.
- Use `App.Background` for a reconnecting consumer/dispatcher.
- Use `platform/lease` for durable exclusive ownership across replicas. Renew before expiry, stop immediately on `ErrLost` or interruption, and release on clean completion.
- Use `sql.Manager.Lock` for a short singleton scheduled action when PostgreSQL session ownership is appropriate.
- Use `SELECT … FOR UPDATE SKIP LOCKED` for queue-shaped rows so replicas partition work without leader election.
- Use `platform/executor` when a request starts work that must be tracked, canceled, drained, and surfaced.

Every loop bounds concurrency, batch size, per-item duration, retry attempts, and shutdown wait. A ticker is stopped. A channel producer owns closure. A handler never launches a raw goroutine.

## Signals

Publish counters and latency for publish/consume/result/retry, current worker concurrency, queue and DLQ depth, outbox count/oldest age, consumer disconnect duration, lease renew/loss/interruption, schedule delay, replay outcome, and reconciliation drift. Trace event ID and correlation across producer and consumer.

## Common misuse

- Direct publish inside the transaction: the broker cannot participate in the database commit.
- Commit then publish: a crash loses the fact.
- Acknowledge before the effect is durable.
- Generate a new event ID during replay.
- Infinite retry of a permanent schema/validation failure.
- Use an in-memory lock for work shared across replicas.
- Retry an approved high-risk side effect blindly; reconcile its owner first.

