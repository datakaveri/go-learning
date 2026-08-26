---
title: Persistence architecture
description: Data ownership, PostgreSQL/PostGIS, Redis, Elasticsearch, S3, transactions, backup, and recovery.
---

# Persistence architecture

**Status: Partially implemented.** Shared PostgreSQL and cache APIs exist; service ownership and adoption vary. Search and object storage remain focused top-level capabilities rather than complete `platform/*` modules.

## Ownership rules

- One service owns a business dataset and is its only writer.
- Other services use an API, event projection, or exported data product—not a cross-service SQL join.
- The owner defines schema, migrations, retention, backup/restore, data classification, audit, and deletion.
- Service code is identical across deployments; DSN, search path, capacity, and feature configuration select behavior.
- A migration history table and migration actor are scoped to one service-owned schema/database.

## Store decision

| Store | Use | Do not use for |
|---|---|---|
| PostgreSQL | Durable transactional state, constraints, idempotency, outbox, leases | Search-shaped denormalized read models owned elsewhere |
| PostGIS | Owned geometry, spatial indexing/query, OGC backing state | Raw policy-provided spatial SQL |
| Elasticsearch | Catalogue/search/NGSI-LD projections and query indexes | Cross-service transaction truth |
| Redis | Bounded cache, limiter state, short-lived approvals/session coordination by design | Irreplaceable durable business state |
| S3-compatible storage | File/object bytes and derived artifacts | Authorization metadata or long-lived bearer capability inventory |

## PostgreSQL details

Use bounded pools and context-aware operations. Expose pool acquisition/usage and slow query signals. Map unique/no-row/serialization/deadlock/connection failures consistently. Use explicit nullable types and UTC timestamps.

Choose isolation based on invariants:

- `ReadCommitted` for ordinary independent operations;
- row locks for conflicting state transitions;
- `RepeatableRead` for consistent multi-read views when appropriate;
- `Serializable` plus `DoRetry` for invariants that cannot be guarded more narrowly.

`DoRetry` callbacks must be safe to replay from their start. External side effects remain outside the transaction and use outbox/idempotency.

## Atomic state and event

```go
err := tx.Do(ctx, func(txCtx context.Context) error {
    if err := orders.MarkPaid(txCtx, orderID, providerEventID); err != nil {
        return err
    }
    return events.Publish(txCtx, outbox, orderPaid, OrderPaid{
        OrderID: orderID,
    }, events.WithID(providerEventID), events.WithCorrelationID(correlationID))
})
if err != nil {
    return fmt.Errorf("record payment: %w", err)
}
dispatcher.Kick()
```

The side-effect owner records the provider event/idempotency key and result. A repeated request with the same key and different payload is a conflict, not a second effect.

## Elasticsearch

The owning service versions mappings/templates, creates aliases, bounds queries/pages, and supplies a reindex/rebuild plan. An outbox/event projection needs a watermark and reconciliation against authoritative state. Mapping/index naming must be configuration-controlled and consistent; unresolved index naming is a current limitation in `dx-dataplane-rs-go`.

## S3-compatible storage

Own bucket/key namespace, content metadata, checksums, encryption, malware/quarantine/processing state, lifecycle, retention, and delete semantics. Generate short-lived least-privilege credentials or presigned operations only after authorization. Bind capabilities to one object/action/range where available and never log them.

## Backup and recovery

Document Recovery Point Objective (RPO), Recovery Time Objective (RTO), backup schedule/encryption/retention, restore steps, ownership, and the consistency relationship among PostgreSQL, outbox, object store, and projections. Test restore into an isolated environment, then replay/reconcile derived indexes and caches.

## Required signals and tests

- pool saturation/acquire duration, slow/error/retry/lock metrics;
- migration version/failure and schema readiness;
- cache hit/miss/error/invalidation lag;
- index lag/reconciliation drift/reindex progress;
- object operation latency/error/integrity failures;
- empty/upgrade migration, rollback, concurrency, idempotency, backup/restore, and projection rebuild tests.

Detailed APIs: [SQL, transactions, migrations, and cache](../platform/persistence-cache.md).

