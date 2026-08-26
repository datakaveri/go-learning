---
title: 3. Persistence, cache, events, and workers
description: Add service-owned PostgreSQL, migrations, transactions, outbox, consumers, and bounded work.
---

# 3. Persistence, cache, events, and workers

**Status: Shared APIs implemented; tutorial adoption required.** Elasticsearch, S3, and some service adapters remain outside `platform/*`; do not invent platform imports.

## Step 1 — own the store

Purpose: make the service the only writer and contract owner for its business state.

Files: `db/migrations/*.up.sql`, `*.down.sql` where policy permits, `internal/repository/postgres`, migration tests.

```sql
CREATE TABLE widgets (
    id uuid PRIMARY KEY,
    owner_id text NOT NULL,
    organisation_id text NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE example_outbox (
    id text PRIMARY KEY,
    topic text NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL,
    sent_at timestamptz,
    attempts integer NOT NULL DEFAULT 0,
    claimed_by text,
    claimed_until timestamptz,
    claim_token text
);
```

Use the exact outbox schema expected by the current platform migrations/tests; the abbreviated DDL above is illustrative and must be checked against `platform/events` before use.

Common mistake: reading another service’s table or running DDL from every replica.

Verification: migrate empty database and last released schema using the same image. Expected: one migration history table and schema compatible with the binary.

## Step 2 — repository adapter

Use `platform/database/sql` interfaces. Explicit parameterized SQL is always acceptable; `sql.Repo[T]` is useful for conventional rows; sqlc is useful for a stable query set. Keep row structs and mapping inside the adapter.

```go
func (r *Repository) ByID(ctx context.Context, id string) (domain.Widget, error) {
    row, err := dxsql.SQLOne[widgetRow](ctx, dxsql.Conn(ctx, r.db), `
        SELECT id, owner_id, organisation_id, name, created_at
          FROM widgets
         WHERE id = $1`, id)
    if err != nil {
        return domain.Widget{}, mapRepositoryError(err)
    }
    return row.domain(), nil
}
```

Common mistake: calling the raw pool inside `Manager.Do`; it escapes the ambient transaction.

Verification: repository integration tests for not-found, uniqueness, cancellation, nullable values, ordering/pagination, and transaction rollback.

## Step 3 — transaction and outbox

Purpose: prevent a committed change with a missing fact or a fact for rolled-back state.

Files: application use case, event contract, outbox migration, dispatcher wiring.

Inside one `app.Tx.Do`, write the row and `events.Publish` to `events.Outbox`. After commit, call `Dispatcher.Kick`. Register `Dispatcher.Run` with `app.Background`.

Common mistake: publish before/after commit without the outbox, or retry a partial nested operation.

Verification: integration test forces domain insert failure and outbox failure; neither can commit alone. Expected: one durable change and one outbox event, with duplicate downstream delivery tolerated.

## Step 4 — event consumer

Purpose: maintain a local projection or respond to another owner’s fact.

Files: `internal/worker/<event>_consumer.go`, processed-event/idempotency migration, reconciliation job, contract tests.

Classify messages as success, transient retry, drop, or quarantine. Mark an event processed in the same transaction as the local projection. Preserve original event ID on replay.

Verification: duplicate, reordered, unsupported-version, poison, broker reconnect, DLQ/replay, and reconciliation tests. Expected: duplicates cause one business effect; broker restart resumes consumption.

## Step 5 — cache only derived data

Use `platform/cache` with a namespaced versioned key and bounded TTL. Define the authoritative fallback and invalidation event. Security-sensitive cached allows must be revision/expiry bound; a miss or cache outage never becomes allow.

Verification: hit, miss, stale, corrupt value, Redis outage, invalidation duplicate/loss/reconciliation. Expected: correctness does not depend on cache availability unless the service explicitly fails closed.

## Step 6 — workers and schedules

Use `App.Go`/`Background`, `platform/lease`, `sql.Manager.Lock`, `SKIP LOCKED`, or `platform/executor` according to ownership. Bound concurrency/batches/timeouts/retries. Stop on context cancellation or lease loss; drain before infrastructure closes.

Verification: crash mid-item, lease loss, duplicate ownership, cancellation, slow dependency, retry exhaustion, and graceful shutdown tests. Expected: no orphan goroutine or unowned high-risk side effect.

## Store selection

- PostgreSQL: transactional state.
- PostGIS: spatial state owned by the service.
- Elasticsearch: search/query projection with mapping/reindex/reconciliation.
- Redis: bounded cache or explicitly ephemeral coordination.
- S3-compatible storage: objects; metadata/ownership remains in an owning service.

Record backup, restore, retention, reindex/rebuild, and ownership before adding any store. Next: [security and gateway](./04-security-gateway.md).

