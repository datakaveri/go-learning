---
title: Structs & Methods
sidebar_label: Structs & Methods
description: Data modeling in Go — structs, value vs pointer receivers, embedding, and struct tags.
---

# Structs & Methods

## Learning objectives

- Define structs and initialize them the platform way (field names, omit zero values).
- Attach methods and choose correctly between value and pointer receivers.
- Use embedding for composition and explain why Go has no inheritance.
- Read struct tags — the metadata that drives JSON, config, and database mapping.

## Prerequisites

- [Collections](collections)

## Time estimate

**4 hours**

## Concepts

### Structs and initialization

```go
type Policy struct {
	ID         string
	UserID     string
	ItemID     string
	Constraint map[string]any
	ExpiresAt  time.Time
}

// Platform style: field names, zero-value fields omitted, &T{} not new(T)
p := &Policy{
	ID:     uuid.NewString(),
	UserID: userID,
	ItemID: itemID,
}
```

Positional initialization (`Policy{"a", "b", ...}`) is banned by the linters for any struct with more than a couple of fields — it silently breaks when fields are added.

### Methods and receivers

A method is a function with a **receiver**:

```go
func (p *Policy) Expired(now time.Time) bool {
	return !p.ExpiresAt.IsZero() && now.After(p.ExpiresAt)
}
```

**Value receiver** (`func (p Policy)`) gets a copy; **pointer receiver** (`func (p *Policy)`) can mutate and avoids copying. The decision rule:

1. Does the method mutate the receiver? → pointer.
2. Is the struct large, or does it contain a lock/pool/connection? → pointer.
3. Otherwise a value receiver is fine — but **don't mix**: if any method needs a pointer receiver, give all methods pointer receivers.

The subtle trap: a value stored in a variable of interface type only exposes the **value-receiver** method set. If `Expired` has a pointer receiver, `Policy{}` (the value) does not satisfy an interface requiring `Expired` — `&Policy{}` does. You'll feel this in the next page on [Interfaces](interfaces).

### Embedding — composition, not inheritance

Go has no `extends`. Instead a struct can **embed** another, promoting its fields and methods:

```go
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccessRequest struct {
	Timestamps          // embedded — no field name
	ID     string
	Status string
}

r := AccessRequest{}
_ = r.CreatedAt // promoted from Timestamps
```

This is composition with sugar: there's no polymorphic override, no `super`. When embedding a type with methods, those methods join your method set — which is how services embed shared config structs and mutexes. Embed judiciously: it's for *is-composed-of*, not for simulating class hierarchies.

### Struct tags

Tags are string metadata on fields, read by libraries via reflection:

```go
type ServerConfig struct {
	Port    int    `mapstructure:"port" json:"port"`
	Host    string `mapstructure:"host" json:"host,omitempty"`
	private string // unexported: invisible to json/mapstructure entirely
}
```

- `json:"..."` — controls JSON field names, `omitempty`, `-` to skip.
- `mapstructure:"..."` — how viper binds YAML/env config to structs (Module 2, [Configuration](../module-2-intermediate/configuration)).
- `db:"..."` — how pgx maps rows to structs (Module 3, [Database Patterns](../module-3-advanced/database-patterns)).

Only **exported** fields participate — a classic beginner bug is a struct that marshals to `{}` because its fields are lowercase.

:::info[Platform connection]
Every DX domain entity is a plain struct with tags: look at `dx-acl-go/internal/domain/policy.go` (domain entities), or any `internal/config/config.go` (a struct of `mapstructure`-tagged nested structs fed to `dxconfig.LoadService`). The generic DAO maps database rows into structs by their `db` tags via `pgx.RowToStructByName` — struct tags are literally the platform's ORM.
:::

## Exercises

1. Model a `Dataset` (id, name, provider, tags, createdAt) and a `Provider`; give `Dataset` a `String()` method and an `Age(now time.Time) time.Duration` method. Choose receivers deliberately and justify them in comments.
2. Demonstrate the mixed-receiver interface trap: define `type Named interface{ Name() string }`, implement it with a pointer receiver, and show which of `T{}` / `&T{}` assignments compile.
3. Create a `Base` struct with `ID` + timestamps, embed it in two entities, and write one function that accepts anything with promoted `ID` — notice you *can't* without an interface. Write the interface.
4. Round-trip a tagged struct through `encoding/json` — include an unexported field and an `omitempty` field, and explain the output.

## Check yourself

- Three reasons to choose a pointer receiver?
- Why must you not mix receiver kinds on one type?
- What does embedding promote, and what does it *not* do (compared to inheritance)?
- Why do unexported fields vanish from JSON?

## References

- [A Tour of Go — Methods](https://go.dev/tour/methods/1)
- [Effective Go — Embedding](https://go.dev/doc/effective_go#embedding)
- [Go by Example: Structs](https://gobyexample.com/structs), [Methods](https://gobyexample.com/methods), [Struct Embedding](https://gobyexample.com/struct-embedding)
- [Google Go Style Guide — Receiver type](https://google.github.io/styleguide/go/decisions#receiver-type)
