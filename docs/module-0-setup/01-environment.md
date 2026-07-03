---
title: Environment Setup
sidebar_label: Environment Setup
description: Install Go, an editor with gopls, Docker, and Git — and prove the toolchain works.
---

# Environment Setup

## Learning objectives

- Install and verify Go 1.22+, Docker Desktop, and Git.
- Set up an editor with `gopls` (the Go language server), format-on-save, and lint integration.
- Understand the handful of Go environment concepts that still matter (`GOPATH` mostly doesn't, `GOBIN` and module cache do).
- Compile, run, test, and vet a first program.

## Prerequisites

None. This is the starting line.

## Time estimate

**2–3 hours** (mostly downloads).

## Concepts

### The toolchain

Go ships as a single toolchain: compiler, formatter, test runner, vet, profiler, and module manager are all subcommands of `go`. There is no build-tool zoo — you will not need Maven/Gradle/webpack equivalents.

```bash
# macOS
brew install go

# Verify — the platform targets Go 1.22+
go version
```

The commands you'll use daily:

| Command | What it does |
|---|---|
| `go run .` | Compile and run the current package |
| `go build ./...` | Compile everything in the module |
| `go test ./...` | Run all tests |
| `go vet ./...` | Static analysis for likely bugs |
| `gofmt -l .` | List files not formatted (must be empty in CI) |
| `go mod tidy` | Sync `go.mod`/`go.sum` with actual imports |

### Editor

Any editor that speaks the Language Server Protocol works; **VS Code with the official Go extension** or **GoLand** are the common choices on the team. Make sure:

- `gopls` is enabled (VS Code's Go extension installs it automatically).
- Format-on-save is on — Go has exactly one formatting style; `gofmt` output is not negotiable and CI enforces it.
- `goimports` runs on save (manages the import block for you).

Install the linter the platform's CI runs:

```bash
brew install golangci-lint
golangci-lint --version
```

### Environment variables that still matter

Modern Go (modules era) needs almost no environment setup. Three things worth knowing:

- **Module cache** — downloaded dependencies live in `$GOPATH/pkg/mod` (default `~/go/pkg/mod`). You never edit this.
- **`GOBIN`** — where `go install` puts binaries (default `~/go/bin`). Add it to your `PATH`.
- **`GOOS`/`GOARCH`** — cross-compilation switches; the platform's Dockerfiles use these to build Linux binaries from any host.

```bash
# Add to your shell profile
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Docker and Git

The DX stack runs locally under Docker Compose, so you need **Docker Desktop** running with a reasonable resource allocation (≥8 GB memory recommended — the full stack includes PostgreSQL, Elasticsearch, Keycloak, RabbitMQ, Redis, and MinIO).

```bash
docker --version && docker compose version
git --version
```

### Hello, DX

Prove everything works end to end:

```bash
mkdir hello-dx && cd hello-dx
go mod init example.com/hello-dx
```

Create `main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Data Exchange")
}
```

```bash
go run .          # prints the greeting
go vet ./...      # no output = no findings
gofmt -l .        # no output = correctly formatted
```

:::info[Platform connection]
Those last three commands aren't a toy — `go build ./...`, `go test ./...`, `gofmt -l .`, `go vet ./...`, and `golangci-lint run` are **the literal PR gate** for every DX Go service. You'll run them hundreds of times. Build the habit on day one.
:::

## Exercises

1. Write a program that prints the Go version it was compiled with (hint: `runtime.Version()`).
2. Deliberately mis-indent your `main.go`, run `gofmt -l .` and then `gofmt -w .`, and diff the result. Notice you never argue with it.
3. Cross-compile your hello program for Linux (`GOOS=linux GOARCH=amd64 go build`) and confirm with `file` that you got an ELF binary.

## Check yourself

- What five commands make up the platform's PR gate?
- Where do downloaded module dependencies live, and do you ever edit them?
- Why is there no formatting debate in Go code review?

## References

- [Download and install Go](https://go.dev/doc/install)
- [How to Write Go Code](https://go.dev/doc/code) — modules, layout, first test
- [golangci-lint](https://golangci-lint.run/)
- Platform: `claude-docs/QUICK-START.md` in the `cdpg-claude` repo (full prerequisites list)
