---
title: Database Patterns
sidebar_label: Database Patterns
description: Service-owned data, repositories, parameterized queries, paging, and escape hatches.
---

# Database patterns

## Ownership first

One service owns a schema's writes, migrations, backup, and recovery. Other services use its API or event projection. A shared PostgreSQL server does not imply shared table ownership.

## Platform interfaces

platform/database/sql provides DB, Querier, Tx, Config, Manager, generic Repo[T], predicate/query helpers, and error mapping. Application packages define repository ports; PostgreSQL adapters use these interfaces.

~~~go
widgets := dxsql.NewRepo[widgetRow](
    db,
    dxsql.WithTable[widgetRow]("widgets"),
    dxsql.WithID[widgetRow]("widget_id"),
)

row, err := widgets.Where(
    dxsql.Eq("org_id", orgID),
    dxsql.Eq("widget_id", widgetID),
).One(ctx)
~~~

Repo derives projections from db tags and avoids SELECT *. Insert returns generated columns. Update and Delete require a predicate.

## Query selection

- Use Repo and Query for single-table CRUD, filtering, sorting, paging, and row locking.
- Use generated sqlc queries for stable complex SQL.
- Use dxsql.SQL or SQLOne for dynamic parameterized SQL with joins, CTEs, JSONB, PostGIS, or window functions.

The DSL deliberately does not grow into a complete SQL language.

## Safety

- Pass values as parameters.
- Resolve sort fields through dxsql.Sortable.
- Use deterministic ordering.
- Translate not found and constraints at the adapter boundary.
- Bound pool size and query time.
- Use FOR UPDATE SKIP LOCKED for database-backed work queues.

## Exercise

Implement list widgets by organization with allowed sort keys, stable tie-breaker, paging, and a row-to-domain mapper. Add an injection attempt as a test and verify it never becomes SQL.

## Check yourself

- Why does generic Repo use row types rather than domain entities?
- When should you choose sqlc?
- Why do Update and Delete refuse empty predicates?
- What does service-owned data prohibit?
