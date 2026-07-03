---
title: Database Access Patterns
sidebar_label: Database Patterns
description: pgx v5, connection pools, parameterized SQL, the repository pattern, and the generic BaseDAO[T].
---

# Database Access Patterns

## Learning objectives

- Connect to Postgres with a `pgxpool.Pool` and understand what pooling buys.
- Write parameterized queries — and explain why identifiers can never be parameters.
- Implement the repository pattern with struct-tag row mapping.
- Use the platform's generic repository, `BaseDAO[T]`, and its query DSL.
- Handle `pgx.ErrNoRows` and soft-deletes correctly.

## Prerequisites

- [Project Structure](project-structure) (repository layer), [Generics](../module-1-go-fundamentals/generics)

## Time estimate

**5 hours**

## Concepts

### pgx and the pool

The platform talks to Postgres with **pgx v5** directly — **no ORM** (no GORM, no ent). SQL is written as SQL; Go maps rows to structs. One pool per service, created at boot (a hard dependency — `Fatal` if unreachable):

```go
pool, err := pgxpool.New(ctx, cfg.DSN())
if err != nil { ... }
defer pool.Close()

if err := pool.Ping(ctx); err != nil { ... } // fail fast, per the boot contract
```

The pool hands out connections per query and reclaims them; you tune `max_conns` in config. Every query takes `ctx` — cancellation and timeouts propagate to the database ([Context](../module-2-intermediate/context) paying off again).

### Parameterized SQL — values yes, identifiers never

```go
// Values: ALWAYS placeholders. Never string-concatenate a value into SQL.
row := pool.QueryRow(ctx,
	`SELECT id, user_id, item_id, expires_at
	   FROM policies WHERE id = $1 AND deleted_at IS NULL`, id)
```

Placeholders (`$1`) send values out-of-band — SQL injection becomes structurally impossible for values. But placeholders **cannot** carry identifiers (table names, columns, `ORDER BY` keys). Dynamic identifiers must come from a **code-side allowlist**:

```go
var sortCols = map[string]string{"created": "created_at", "name": "name"}

col, ok := sortCols[req.Sort]
if !ok {
	return dxerrors.NewValidation("invalid sort key")
}
query := fmt.Sprintf(`SELECT ... ORDER BY %s`, col) // safe: col came from OUR map
```

This is the resolution of the whitelist rule from [REST API Development](rest-api-development) — user input picks *among* identifiers you wrote; it never *is* the identifier.

### Repositories and row mapping

The repository owns all SQL for an aggregate and returns domain structs. pgx v5's collection helpers map rows by `db` tags — the reflection-inside-a-library pattern from [Module 1](../module-1-go-fundamentals/reflection-and-when-not):

```go
type Policy struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	ItemID    string     `db:"item_id"`
	ExpiresAt *time.Time `db:"expires_at"`
}

rows, err := pool.Query(ctx, `SELECT id, user_id, item_id, expires_at FROM policies WHERE user_id = $1`, uid)
if err != nil {
	return nil, fmt.Errorf("query policies: %w", err)
}
policies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Policy])
```

Absence is an error value, not an exception: single-row lookups return `pgx.ErrNoRows`, which the repository translates to the platform taxonomy — `errors.Is(err, pgx.ErrNoRows)` → `dxerrors.NewNotFound(...)`. Handlers never see pgx errors.

### BaseDAO[T] — the generic repository

CRUD, pagination, and soft-delete plumbing is identical for every entity, so `dx-common-go` provides it once, generically:

```go
dao := dao.NewBaseDAO[Policy](pool, "policies")

p, err := dao.FindByID(ctx, id)
page, err := dao.FindPage(ctx, conds, order, limit, offset) // returns Page[Policy]
```

Alongside it, a **query DSL** (`database/postgres/query`) builds `WHERE`/`ORDER BY` safely — conditions compose from typed operators, values ride as parameters, identifiers stay code-side. The division of labor:

- **BaseDAO + DSL** — standard lookups, list endpoints, CRUD. The default.
- **Raw SQL in the repository** — joins, aggregates, anything the DSL doesn't express. Fully legitimate; same parameterization rules.

There is deliberately no ORM layer on top: what you write is what runs, `EXPLAIN` output maps to your code, and the escape hatch is just… SQL.

### Soft deletes

Platform tables use `deleted_at` timestamps rather than hard `DELETE`s (audit trails, undelete). The catch: **every read must filter** `deleted_at IS NULL` — explicitly. A forgotten filter resurrects deleted data, which in a policy table is a security bug, not a cosmetic one. BaseDAO handles it for its own queries; your raw SQL must remember.

:::info[Platform connection]
`dx-common-go/database/postgres` holds all of it: the pool factory, `dao/base.go` (`BaseDAO[T]` — read it top to bottom, it's short and it's the generics page made real), and `query/` (the ConditionBuilder with its typed operators). Schema note: Go services run against the **legacy databases** with the Java stack's Flyway as schema owner — Go code issues no DDL against legacy tables, only the idempotent "schema ensure" for its own additions. `claude-docs/DATABASE.md` maps every table to its owning service.
:::

## Exercises

*(Local stack up — use its Postgres, or `docker run postgres:16` for scratch space.)*

1. Give `dx-scratch-go` a real Postgres repository: schema-ensure SQL at boot (embedded, idempotent), CRUD with `pgx.CollectRows`/`RowToStructByName`, `ErrNoRows` → NotFound translation.
2. Attack yourself: write the vulnerable string-concatenation version of a search endpoint in a throwaway branch, inject `' OR '1'='1` through curl, then fix with `$1` and watch the injection become a literal string match.
3. Implement sorted, paginated listing with the identifier-allowlist pattern; wire it to your `request.Builder`-style query parsing from the REST page.
4. Add `deleted_at` soft deletes: DELETE endpoint sets it, all reads filter it, and one test proves a soft-deleted note is invisible to every list and get.
5. Read `dao/base.go` in `dx-common-go` and write down: how it names the soft-delete column, how `WithTx` works (preview of the [next page](transactions)), and one method you'd have designed differently.

## Check yourself

- Why can `$1` carry `WHERE user_id = ?` but not `ORDER BY ?`?
- Where do pgx errors stop existing, and what replaces them?
- When BaseDAO vs raw SQL — and what rule applies to both?
- Why is a missed soft-delete filter a security bug here?

## References

- [pgx v5 docs](https://pkg.go.dev/github.com/jackc/pgx/v5) · [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
- [OWASP: SQL Injection Prevention](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html)
- Platform: `dx-common-go/database/postgres/{dao,query}`; `claude-docs/DATABASE.md`
