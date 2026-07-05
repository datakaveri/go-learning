---
title: Distributed Systems Concepts
sidebar_label: Distributed Systems
description: Eventual consistency, idempotency, retries with backoff, timeouts as contracts — the theory behind the platform's rules.
---

# Distributed Systems Concepts

## Learning objectives

- Reason about partial failure — the defining property of distributed systems.
- Explain delivery guarantees and why "exactly once" is a lie you compensate for with idempotency.
- Apply timeouts, retries with exponential backoff + jitter, and know where circuit breakers fit.
- Understand eventual consistency as a business decision, with the policy-propagation flow as the case study.

## Prerequisites

- [RabbitMQ & Events](event-driven-rabbitmq), [Transactions](transactions), [AuthN & AuthZ](authn-authz)

## Time estimate

**3.5 hours** — a reading-heavy page that names the principles behind rules you've already practiced.

## Concepts

### Partial failure is the whole subject

On one machine, a call either happens or the process is dead. Between machines there's a third state: **unknown**. You sent a request and got no answer — did the other side do the work? You cannot know. Every pattern on this page is a strategy for acting sensibly under that uncertainty. Fifteen services, a broker, and four datastores means the DX platform is *always* partially failing somewhere; the design goal is that nobody notices.

### Delivery guarantees, honestly

| Guarantee | Cost | Failure mode |
|---|---|---|
| At-most-once | cheap | silent loss |
| **At-least-once** | modest | **duplicates** |
| Exactly-once | not generally achievable across systems | — |

The platform chose at-least-once everywhere (outbox → confirms → redelivery), which converts the impossible "exactly once" into an engineering task you've already done: **idempotent consumers**. Generalize the principle: *any* retried operation — message handling, cron runs, HTTP retries — needs an idempotency story. The unknown-state problem above is exactly why: you must be free to try again without fear.

### Timeouts are contracts

A call without a timeout can wait forever; a blocked goroutine per stuck request means one slow dependency drains your whole service (every handler waiting on it, connection pool exhausted, healthy endpoints starved). Every network call on the platform is bounded — handler-level (the StandardStack's `Timeout`), per-call (`context.WithTimeout` into pgx and HTTP clients), and server-level (read/write timeouts).

Set them deliberately: a timeout is a **promise about your dependency's behavior**, and they compose — if your handler promises 15s and makes three sequential 10s calls, the math already failed. Budget top-down.

### Retries: backoff + jitter, or you attack yourself

A failed call is often transient — but naive retry loops turn a hiccup into an outage: a thousand clients retrying instantly and in sync is a DDoS against a recovering service. Civilized retries:

```go
delay := base // e.g. 100ms
for attempt := 0; attempt < maxAttempts; attempt++ {
	err = call(ctx)
	if err == nil || !retryable(err) {
		return err
	}
	jittered := delay/2 + time.Duration(rand.Int64N(int64(delay))) // spread the herd
	select {
	case <-time.After(jittered):
	case <-ctx.Done():
		return ctx.Err()
	}
	delay *= 2 // exponential
}
```

Rules: only retry **retryable** errors (timeouts, 503s — not validation failures); cap attempts; respect the context; and only retry operations that are safe to repeat (idempotency again — it's the same principle every time). A **circuit breaker** is the next escalation — after N consecutive failures stop calling entirely for a cool-down, failing fast instead of queueing doomed work. Know the concept; the platform doesn't yet ship a standard one (a recorded gap: HTTP resilience helpers are on the shared-library roadmap).

### Eventual consistency — the case study you're running

Create a policy on the platform and for a moment the PDP doesn't know yet:

```
policy row committed (acl) → outbox → RabbitMQ → authz consumer → OpenFGA tuple
        t=0                                                      t≈1–3s
```

Between t=0 and t≈3s, the system is **inconsistent**: the PAP says the policy exists; a gateway check still says deny. The platform accepts this because the alternative — synchronous cross-service writes — couples availability (policy creation down whenever authz is down) for a guarantee nobody needs: a consumer waiting three seconds for access to activate is fine; **revocation** lag is the direction to watch, and `make dev-demo` explicitly measures both directions.

The general lesson: consistency windows are **product decisions**, not implementation accidents. For each flow ask: who observes the lag, how long is acceptable, and in which direction does it fail safe? "Deny-by-default until the tuple lands" fails safe; the reverse would not.

### Health as a distributed-systems primitive

Liveness ("process is alive — restart me if not") vs readiness ("dependencies reachable — route traffic to me if so") is how Kubernetes makes *cluster-level* decisions from *service-level* self-knowledge. A service that reports ready while its database is down lies to the load balancer and turns one failure into many. The mechanics land in [Observability](observability); the principle belongs here.

:::info[Platform connection]
Nothing on this page is hypothetical: at-least-once + idempotency is the `authz` pipeline; timeout budgeting is the StandardStack + pgx contexts; eventual consistency is measured by `make dev-demo`'s 3-second propagation checks; fail-safe direction is OpenFGA's deny-by-default. When Module 4's architecture deep dive shows you the full picture, every arrow on the diagram will have one of this page's principles attached.
:::

## Exercises

1. Write the retry helper above as `Retry(ctx, attempts, base, fn)` with tests using a flaky fake (fails N times, then succeeds). Include a non-retryable-error case that must return immediately.
2. Measure the platform's consistency window: script create-policy → poll the check endpoint until allowed, timestamping both. Run it ten times; report min/median/max. Then do the same for revocation.
3. Break a timeout budget on purpose: handler timeout 2s, inner call sleeping 5s with its own 10s timeout. Observe which layer wins and what the client sees. Fix the budget.
4. Thought exercise (write actual answers): for the marketplace's "purchase grants dataset access" flow, what's the acceptable consistency window, who observes it, and which failure direction is safe?

## Check yourself

- What is the "third state" of a distributed call and what does it force on your design?
- Why is exactly-once delivery not a thing you can buy, and what substitutes for it?
- Why jitter, not just backoff?
- For policy propagation: which direction of staleness is dangerous and why?

## References

- [Designing Data-Intensive Applications](https://dataintensive.net/) (Kleppmann) — chapters 8–9 are this page in book form
- [AWS Builders' Library: Timeouts, retries, and backoff with jitter](https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)
- [Fallacies of distributed computing](https://en.wikipedia.org/wiki/Fallacies_of_distributed_computing)
- Platform: `claude-docs/ARCHITECTURE.md` (data flows); `make dev-demo` source
