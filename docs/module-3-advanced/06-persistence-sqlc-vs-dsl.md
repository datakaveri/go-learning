---
title: Repository DSL, sqlc, or SQL
sidebar_label: DSL vs sqlc
description: Choose the smallest persistence tool that keeps queries explicit and transaction aware.
---

# Repository DSL, sqlc, or SQL

## Three supported paths

| Query shape | Tool |
|---|---|
| Single-table CRUD, predicates, order, paging, locks | platform/database/sql Repo and Query |
| Stable complex SQL known at build time | sqlc with dxsql.Conn |
| Dynamic complex SQL using JSONB, PostGIS, CTEs, or windows | dxsql.SQL / SQLOne |

Choose by query shape, not personal preference.

## Repository DSL

~~~go
page, err := widgets.Where(
    dxsql.Eq("org_id", orgID),
    dxsql.ILike("name", searchPattern),
).Order(
    dxsql.Desc("created_at"),
    dxsql.Desc("widget_id"),
).Paged(ctx, paging.NewRequest(pageNumber, pageSize))
~~~

The DSL parameterizes values and quotes known columns. Raw accepts a controlled SQL fragment plus values; it is not permission to pass client text.

## sqlc

~~~go
queries := gen.New(dxsql.Conn(ctx, db))
rows, err := queries.TopWidgetsByOrganization(
    ctx,
    gen.TopWidgetsByOrganizationParams{OrgID: orgID},
)
~~~

Conn returns the ambient transaction when ctx carries one. Use generated types directly; do not wrap sqlc in another generic registration layer.

## Dynamic complex SQL

~~~go
rows, err := dxsql.SQL[widgetRow](
    ctx,
    dxsql.Conn(ctx, db),
    statement,
    orgID,
    limit,
)
~~~

Build dynamic identifiers only from an allowlist. Values always remain parameters.

## Review evidence

A change should state why its selected path fits. “Uses a window function with a stable shape” justifies sqlc. “Optional filters change the WHERE predicates, but the row shape is stable” can justify the DSL. Keep SQL visible enough to review indexes and plans.

## Exercise

Write one list query with the DSL, one top-N-per-group query with sqlc-style SQL, and one dynamic geospatial filter with dxsql.SQL. Explain why moving all three into one abstraction would make review harder.

## Check yourself

- How does sqlc participate in Manager.Do?
- What operations are intentionally absent from Query?
- Is dxsql.Raw safe with a client sort string?
- What evidence belongs in a persistence PR?
