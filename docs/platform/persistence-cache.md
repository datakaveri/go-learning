---
title: SQL, transactions, migrations, and cache
description: Use the shared PostgreSQL abstraction and cache safely.
---

# SQL, transactions, migrations, and cache

**Status: Implemented platform APIs; service adoption is partial.** PostgreSQL is the transactional default, PostGIS extends it for spatial ownership, Redis is a cache/ephemeral coordination dependency, and Elasticsearch/S3 use focused foundation packages outside `platform/`.

## PostgreSQL lifecycle

Declare `platform/database/sql.Config` through bootstrap. The platform opens a bounded pool, registers health, exposes `App.DB` and `App.Tx`, and closes it after work drains. Tune pool size from concurrency and database limits; do not copy a fleet-wide number blindly.

The exported `sql.DB`, `Tx`, `Querier`, `Rows`, and `Row` interfaces avoid driver coupling. Use `platform/database/sql/pgx` only for a named capability such as `COPY` or `LISTEN/NOTIFY`, and isolate that import in the adapter.

## Transactions

`sql.Manager.Do` commits on nil and rolls back on error or panic. The transaction travels on context, so repositories use the same query code inside and outside it. Nested `Do` joins the outer transaction. `DoRetry` retries a complete, replay-safe callback for serialization/deadlock failures.

```go
func (s *Service) Create(ctx context.Context, cmd CreateWidget) (domain.Widget, error) {
    widget, err := domain.NewWidget(cmd.ID, cmd.Name, cmd.Owner)
    if err != nil {
        return domain.Widget{}, err
    }

    err = s.tx.Do(ctx, func(txCtx context.Context) error {
        if err := s.widgets.Create(txCtx, widget); err != nil {
            return err
        }
        return events.Publish(txCtx, s.outbox, widgetCreated, WidgetCreated{
            ID: widget.ID(), Owner: widget.Owner(),
        }, events.WithCorrelationID(cmd.CorrelationID))
    })
    if err != nil {
        return domain.Widget{}, fmt.Errorf("create widget: %w", err)
    }
    s.dispatcher.Kick()
    return widget, nil
}
```

The outbox table name is a compile-time constant. Never derive an identifier from user input.

## Repository choices

Use parameterized explicit queries, the generic `sql.Repo[T]`/query DSL for conventional CRUD, or sqlc-generated code that accepts the platform `Querier`. The choice stays inside the adapter. Domain and application packages never import pgx or generated SQL structs.

Use `SELECT … FOR UPDATE` for competing transitions. Choose an isolation level deliberately. A serialization retry is correct only if the whole callback can run again safely.

## Migrations

- Embed ordered, immutable migration files in the service binary.
- Use a per-service history table.
- Run one `migrate-only` actor before application rollout.
- Test from an empty database and from the last released schema.
- Make destructive/long-running changes expand-contract, observable, and reversible.
- Readiness means the schema needed by the binary exists; application replicas do not race to apply DDL.

## Cache

`platform/cache.Store` and the memory/Redis adapters provide typed cache mechanics. The service owns key namespaces, TTL, invalidation, stampede behavior, and whether stale data is safe.

Use cache-aside for derived reads:

1. Read a versioned key.
2. On miss, load from the owner.
3. Store with bounded TTL.
4. Invalidate on the authoritative event and tolerate duplicate invalidation.

A cache failure must either fall back to the owner or deny if the cached data is security-critical and no safe authoritative path exists. Never treat “cache miss” as authorization allow.

## Other stores

- PostgreSQL/PostGIS: transactional, relational, and spatial truth owned by one service.
- Elasticsearch: searchable projection. The owner defines mappings/templates, index aliases, reindex and reconciliation. It is not a transaction coordinator.
- Redis: bounded cache, rate-limit state, short-lived approval/session coordination when explicitly designed. Not durable truth.
- S3-compatible storage: object bytes. The file service owns bucket/key namespaces, metadata links, retention, checksums, processing/quarantine, and short-lived capabilities.

## Failure and observability

Expose pool saturation/acquisition time, slow queries, transaction retries, migration version, cache hit/miss/error, and invalidation lag. Backups are incomplete until restore is tested. Search indexes and cache projections need rebuild/reconciliation procedures.

## Common misuse

- Cross-service joins into another owner’s tables.
- Running DDL from every replica.
- Publishing an event after commit with no outbox.
- Retrying one statement inside a failed transaction.
- Building SQL from policy-provided text.
- Caching an allow decision past its expiry/revision.
- Treating Elasticsearch or Redis as a hidden business source of truth.

