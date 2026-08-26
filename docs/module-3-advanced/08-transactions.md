---
title: Transactions and Outbox
sidebar_label: Transactions
description: Context-propagated transactions, retries, locks, and atomic event publication.
---

# Transactions and outbox

## Context propagation

~~~go
err := txManager.Do(ctx, func(txCtx context.Context) error {
    order, err := orders.Create(txCtx, input)
    if err != nil {
        return fmt.Errorf("create order: %w", err)
    }
    if err := eventOutbox.OrderCreated(txCtx, order); err != nil {
        return fmt.Errorf("write order event: %w", err)
    }
    return nil
})
~~~

Repositories use the transaction carried by txCtx. Passing the original ctx accidentally runs outside the transaction. Nested Do joins the outer transaction and only the outer call commits.

## Retry

DoRetry retries serialization failure and deadlock with jittered backoff. Its callback can run multiple times:

- do not send email, HTTP, or broker messages inside it;
- generate stable IDs outside or make generation deterministic;
- place external effects in an outbox;
- preserve context deadlines.

## Isolation and locks

Use the weakest isolation that preserves the invariant. Row locks protect contested records. FOR UPDATE SKIP LOCKED lets replicas claim different queue rows. Manager.Lock provides a non-blocking PostgreSQL advisory lock for truly singleton work.

## Outbox

platform/events.Outbox stores the event envelope in the same transaction as domain state. Dispatcher publishes pending events later with broker confirmation.

This yields atomic state plus pending publication, not exactly-once delivery. Consumers deduplicate by event ID.

## Exercise

Implement “approve access request”: lock the pending request, reject a repeated decision, update status, create policy state, and write a policy event in one transaction. List what the consumer must do on redelivery.

## Check yourself

- Why is a broker publish inside a SQL transaction unsafe?
- Which context must repositories receive?
- When is DoRetry unsafe?
- What problem does SKIP LOCKED solve?
