---
title: Schema Migrations
sidebar_label: Schema Migrations
description: Embedded migrations, ownership modes, compatibility, backfills, and rollback planning.
---

# Schema migrations

## Two schema modes

The same service binary supports:

- schema_mode=migrate — apply its embedded ordered migrations;
- schema_mode=none — issue no DDL and use a pre-provisioned schema.

Application behavior does not branch on the mode.

## Embedded migrations

~~~go
//go:embed migrations/*.sql
var Migrations embed.FS

bootstrap.Migrations(
    Migrations,
    "migrations",
    "schema_migrations_widget",
)
~~~

Bootstrap runs migrations before opening the application pool. Use a service-specific history table when multiple services share one database server.

## Expand and contract

Rolling deployments overlap versions:

1. Expand with nullable/additive schema.
2. Deploy code that reads both forms and writes the new form.
3. Backfill in bounded, observable batches.
4. Switch reads after evidence is complete.
5. Remove the old form in a later release.

Avoid a single migration that rewrites a large table under an exclusive lock.

## Migration rules

- Only the owning service changes its schema.
- Every change is a reviewed file, never an interactive production ALTER.
- Use explicit transactions only when the database operation supports them.
- Set lock and statement timeouts.
- Separate schema expansion from large data backfills.
- Make restart behavior and partial failure clear.
- A rollback plan may be “forward fix”; do not promise impossible DDL reversal.

## Exercise

Plan adding a required normalized_name field to a large widgets table. Write expansion, dual-write, backfill, verification, read-switch, and cleanup stages with abort criteria.

## Check yourself

- What does schema_mode=none guarantee?
- Why do migrations run before repositories receive a pool?
- What makes a migration safe for rolling deploys?
- Who owns the backup and restore plan?
