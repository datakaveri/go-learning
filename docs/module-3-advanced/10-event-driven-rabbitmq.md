---
title: Event-Driven Services
sidebar_label: RabbitMQ & Events
description: Typed topics, delivery semantics, consumer groups, retries, dead letters, and idempotency.
---

# Event-driven services

## Typed contract

~~~go
type OrganizationMemberAdded struct {
    OrganizationID string
    UserID         string
}

var OrganizationMemberAddedTopic =
    events.NewTopic[OrganizationMemberAdded]("org.member.added")
~~~

Topic binds name, payload type, and version. The Event envelope adds ID, type, version, correlation ID, occurrence time, and payload.

## Publish and subscribe

~~~go
err := OrganizationMemberAddedTopic.Publish(
    ctx,
    bus,
    payload,
    events.WithCorrelationID(requestID),
)

err = OrganizationMemberAddedTopic.Subscribe(
    bus,
    "authz-sync",
    projector.Apply,
)
~~~

Members of one consumer group share work. Distinct groups each receive the event.

## Delivery semantics

RabbitMQ delivery is at least once. Publishers require confirmation. Consumers ack only after success and deduplicate by event ID or a domain key. Retriable failures use bounded backoff; exhausted messages go to a dead-letter queue.

ErrDrop acknowledges an event that can never succeed, such as malformed payload. Make drops visible; silently discarding an unsupported version can hide a broken rollout.

## Reconnect and shutdown

The AMQP adapter owns reconnect, topology, prefetch, confirmations, retry, and dead-letter configuration. Consumers run through bootstrap.Background when reconnectable failure should restart them. A consumer stops accepting delivery, finishes or safely requeues in-flight work, then closes.

## Outbox

Use platform/events.Outbox when state and event must commit atomically. Do not publish directly inside a database transaction.

## Exercise

Implement an idempotent policy projector. Deliver one event twice, inject a transient OpenFGA failure, deliver an unsupported version, and verify ack/retry/drop/dead-letter and metrics behavior.

## Check yourself

- What does a consumer group mean?
- When can a consumer ack?
- Why is exactly-once not the contract?
- What distinguishes ErrDrop from a dead-lettered transient failure?
