---
title: Distributed Systems
sidebar_label: Distributed Systems
description: Timeouts, retries, idempotency, consistency, partial failure, and resilience budgets.
---

# Distributed systems

A remote call can succeed after the caller times out. A message can be delivered twice. Clocks differ. A dependency can be half available. Design contracts around these facts.

## Timeout budget

The edge request deadline is the total budget. Each downstream call gets a smaller child timeout, leaving time to handle the result and respond. Do not reset to context.Background inside a request.

## Retry test

Retry only when:

1. the error is transient;
2. the operation is idempotent or has an idempotency key;
3. the remaining deadline can accommodate another attempt;
4. backoff and jitter prevent synchronized load;
5. retry occurs at one layer.

## Consistency choices

- Strong transaction inside one service-owned database.
- Eventual projection across service or search boundaries.
- Read-your-write from authoritative state when immediate confirmation matters.
- Saga or compensating action for multi-service workflow; no hidden cross-service transaction.

Name the boundary to the client. “Policy accepted” can be different from “policy projection observable.”

## Caches and breakers

A cache is not authoritative unless the domain explicitly says so. Circuit breakers stop repeated calls to a failing dependency; they do not invent a correct fallback. Bulkheads bound concurrent work so one dependency cannot exhaust the process.

## Idempotency

Use a stable operation key and store the result or state transition. A duplicate webhook, callback, event, or client command returns the recorded outcome rather than repeating a charge, grant, or email.

## Exercise

Design a marketplace purchase through payment callback, order state, policy grant, audit, and notification. Mark transaction boundaries, idempotency keys, retry owners, events, compensations, and client-visible states.

## Check yourself

- Can a timeout prove a remote operation failed?
- Where should retry occur?
- Which reads require authoritative state?
- What does a circuit breaker protect?
