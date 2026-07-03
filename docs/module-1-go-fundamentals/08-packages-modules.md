---
title: Packages & Modules
sidebar_label: Packages & Modules
description: Package design, go.mod, internal/, replace directives, and the platform's naming rules.
---

# Packages & Modules

## Learning objectives

- Organize code into packages with intention — by responsibility, not by layer-name reflex.
- Manage dependencies with modules: `go.mod`, `go.sum`, semantic versions, `go mod tidy`.
- Use `internal/` to enforce privacy and `replace` for side-by-side development.
- Apply the platform naming rules: short, lowercase, singular — and never `util`.

## Prerequisites

- [Error Handling](error-handling)

## Time estimate

**3 hours**

## Concepts

### Packages are the unit of design

A package is a directory of `.go` files sharing one `package` name. Everything in it sees everything else; the capitalized subset is the public API. Good packages are **cohesive** (one responsibility, guessable from the name) and **shallow to import** (`policy.New`, not `pkg.NewPolicyFactoryImpl`).

Platform naming rules (enforced in review):

- Short, lowercase, **singular**, no underscores: `policy`, `auth`, `outbox`.
- The package name is part of every call site — don't stutter: `policy.New()`, not `policy.NewPolicy()`.
- **No `util`, `common`, `base`, or `helpers` packages.** A grab-bag name means the contents have no home; find them one.
- Every exported package carries a `doc.go` or a package comment on its main file.

### Modules: go.mod and go.sum

A module is a versioned tree of packages — the unit of distribution:

```
module github.com/datakaveri/dx-acl-go

go 1.22

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/jackc/pgx/v5 v5.6.0
)
```

- `go get pkg@v1.2.3` adds a dependency; `go mod tidy` prunes and completes the graph.
- `go.sum` records cryptographic checksums — commit it; it's supply-chain protection, not clutter.
- Major versions ≥2 live in the import path (`chi/v5`) — that's why you see `/v5` everywhere.

### internal/ — privacy at the module boundary

Packages under `internal/` can only be imported by code in the tree rooted at `internal/`'s parent. This is a compiler guarantee, not a convention:

```
dx-acl-go/
  cmd/server/main.go          ← can import internal/*
  internal/service/policy.go  ← nobody outside dx-acl-go can import this
```

The platform layout puts **everything except `cmd/` in `internal/`** — a service's packages are not a library, and `internal/` makes that stance mechanical.

### replace — side-by-side development

The DX workspace clones every repo inside the orchestrator so services can develop against the shared library without publish-bump-update cycles:

```
require github.com/datakaveri/dx-common-go v0.0.0
replace github.com/datakaveri/dx-common-go => ../dx-common-go
```

Edit `dx-common-go` locally and every service sees the change on next build. This is why the clone-inside-the-orchestrator layout is mandatory — the relative path in `replace` must resolve.

### Import hygiene

Imports come in two groups — stdlib, then everything else — and `goimports` maintains this for you:

```go
import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)
```

Blank imports (`_ "package"`) exist for side-effect registration (database drivers, `embed`); dot imports are banned.

:::info[Platform connection]
Look at any DX service's `go.mod`: the `replace` directive to `../dx-common-go` is right there, and it's the mechanism holding the multi-repo workspace together. The no-`util` rule explains `dx-common-go`'s shape — it is not one grab-bag package but ~28 focused ones (`config`, `httpserver`, `response`, `health`, `auditing`, …), each importable on its own. You'll tour them all in [Module 4](../module-4-platform/dx-common-go-tour).
:::

## Exercises

1. Create a module with packages `store` (exported `Store`, unexported helpers) and `internal/codec`. Prove from a *second* module that you can import `store` but not `internal/codec`.
2. Split your Module-1 key-value mini-project into `cmd/kv/main.go` + `internal/store/` — the platform shape in miniature.
3. Create two local modules where one `replace`s the other; make a breaking change in the library and watch the consumer fail to compile immediately.
4. Run `go mod tidy` after deleting a dependency from your code, and diff `go.mod` before/after.

## Check yourself

- Why is `util` a banned package name?
- What exactly does `internal/` restrict, and who enforces it?
- What problem does the `replace` directive solve in the DX workspace?
- Why must `go.sum` be committed?

## References

- [Go Modules Reference](https://go.dev/ref/mod)
- [How to Write Go Code](https://go.dev/doc/code)
- [Google Go Style Guide — Package names](https://google.github.io/styleguide/go/decisions#package-names)
- Platform: `claude-docs/REPOSITORIES.md` (workspace layout), `CONTRIBUTING.md` (module conventions)
