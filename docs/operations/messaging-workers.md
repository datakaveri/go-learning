---
title: Messaging, workers, and scheduled jobs
description: Platform delivery semantics, envelopes, retries, DLQ, replay, leases, schedules, and recovery.
---

# Messaging, workers, and scheduled jobs

**Status: Shared mechanisms implemented; fleet conformance is partial.** RabbitMQ delivers at least once. Consumers and side-effect owners provide idempotency.

## Event contract

Every event has one owning producer, past-tense type, integer schema version, stable event ID, occurrence time, correlation ID, and documented payload. Define compatibility and retention. Avoid sending credentials or full protected records when references suffice.

## Delivery workflow

1. Application commits domain state plus outbox event.
2. A leased dispatcher publishes through `platform/events/amqp` with publisher confirms.
3. RabbitMQ routes by documented exchange/routing key.
4. Consumer receives, validates envelope/version, checks idempotency, and executes a local transaction/effect.
5. Consumer acknowledges only after success.
6. Transient failures retry with bounded backoff; permanent failures quarantine/dead-letter.
7. Replay preserves event identity; reconciliation compares authoritative and projected state.

## Failure classification

| Class | Examples | Action |
|---|---|---|
| Success/duplicate | Applied, already processed | Ack |
| Drop | Known irrelevant event with documented reason | Record and ack/drop |
| Quarantine | Invalid schema, unsupported version, impossible invariant | DLQ/quarantine with safe reason |
| Transient | Broker/database/network unavailable, serialization conflict | Nack/retry under budget |
| Side-effect unknown | Timeout after external provider may have accepted action | Reconcile by idempotency key before retry |

Never blindly retry payment, notification, approval, or tool side effects when the first outcome is unknown.

## Consumer reconnect and shutdown

Register reconnecting consumers with `App.Background`; the loop re-establishes channels after broker churn, publishes disconnect duration/readiness, and stops on context cancellation. Stop intake, finish/abandon safely bounded in-flight work, flush acknowledgements, then close AMQP. Tests restart RabbitMQ and terminate during an item.

## Worker ownership

- Queue work: database row claim with `FOR UPDATE SKIP LOCKED` or broker delivery.
- Exclusive long operation: durable `platform/lease` with acquire, renew, interruption, release, and immediate stop on loss.
- Short scheduled singleton: `sql.Manager.Lock` when a PostgreSQL session lock is appropriate.
- Request-originated background work: `platform/executor`, bounded and drained.

No goroutine is launched without an owner, cancellation path, error sink, concurrency bound, and shutdown join.

## Scheduled jobs

Define timezone (normally UTC), missed-run behavior, overlap policy, maximum runtime, retry, idempotency, lease, manual trigger, audit, and recovery. A ticker is not a distributed scheduler by itself.

## Required signals

Publish event count/result/latency by bounded type/version, outbox pending/oldest/attempts/lease loss, consumer connection/prefetch/unacked/retry, queue/DLQ depth, replay result, projection lag, reconciliation drift, worker concurrency/lease/interrupt, and schedule due/start/end/delay.

## Verification checklist

- [ ] State/outbox atomicity failure injection.
- [ ] Stable event ID and compatible schema fixtures.
- [ ] Duplicate and out-of-order delivery.
- [ ] Consumer reconnect after broker restart.
- [ ] Retry budget and DLQ/quarantine path.
- [ ] Replay with original ID.
- [ ] Side-effect-owner idempotency conflict.
- [ ] Reconciliation repair and drift metrics.
- [ ] Worker crash, lease expiry/loss/interruption.
- [ ] SIGTERM with in-flight work.

Current API tutorial: [Events and workers](../platform/events-workers.md).

