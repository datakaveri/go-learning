---
title: Bootstrap and typed configuration
description: Compose a service, validate configuration, supervise work, and shut down safely.
---

# Bootstrap and typed configuration

**Status: Implemented.** Current APIs are `platform/bootstrap` and `platform/config`; adoption across service repositories is partial.

## What the modules own

`config.Load[T]` applies platform defaults, then an optional YAML file, then environment variables. It binds typed fields and invokes `Validate()` when implemented. `bootstrap.Run` loads config, creates the logger, applies migrations, opens declared dependencies, wires handlers, starts HTTP and optional gRPC, supervises workers, drains servers, stops work, and closes infrastructure.

They do not decide which database, broker, or business dependency a service needs.

## Typed configuration

```go
package config

import (
    "errors"

    platformconfig "github.com/datakaveri/dx-common-go/platform/config"
    dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"
    "github.com/datakaveri/dx-common-go/platform/security/workload"
)

type Config struct {
    platformconfig.Base `mapstructure:",squash"`
    Postgres dxsql.Config            `mapstructure:"postgres"`
    Workload workload.VerifierConfig `mapstructure:"workload"`
    PublicBaseURL string              `mapstructure:"public_base_url"`
}

func (c Config) Validate() error {
    if c.PublicBaseURL == "" {
        return errors.New("public_base_url is required")
    }
    return c.Workload.Validate()
}
```

Local non-secret defaults:

```yaml
server:
  port: 8080
  max_body_bytes: 1048576
postgres:
  dsn: ""
  search_path: example
workload:
  enforcement: required
  service: dx-example-go
```

Set secrets through the deployment secret mechanism. An empty DSN should fail before serving when PostgreSQL is required.

## Declarative bootstrap

```go
func main() {
    bootstrap.Run(bootstrap.Spec[serviceconfig.Config]{
        Name:    "dx-example-go",
        Version: version,
        Config: platformconfig.Options{
            Defaults: map[string]any{"server.port": 8080},
        },
        Deps: func(cfg *serviceconfig.Config) bootstrap.Deps {
            return bootstrap.Deps{
                Migrations: bootstrap.Migrations(migrations.FS, ".", "schema_migrations_example"),
                Postgres:   bootstrap.Required(cfg.Postgres),
            }
        },
        Wire: wire,
    })
}
```

`Name` and `Wire` are required. `Load` overrides the loader only when the service has a verified need. `GRPC` returns a `bootstrap.Servable` when the service exposes an internal gRPC surface.

## Lifecycle choices

| API | Use | Failure behavior |
|---|---|---|
| `bootstrap.Required(cfg)` | Service cannot answer correctly without it | Startup fails; orchestrator can retry. |
| `bootstrap.Degrade(cfg)` | A documented correct degraded mode exists | Startup continues with a nil handle; `Wire` must branch explicitly. |
| `app.Go(name, fn)` | Essential long-running work | Returned error coordinates shutdown. |
| `app.Background(name, fn)` | Reconnecting consumer or recoverable loop | Failure is logged and restarted with bounded backoff. |
| `app.Closer(name, fn)` | Owned resource not constructed by bootstrap | Runs last-in-first-out after servers and workers drain. |
| `app.Probe(name, checker)` | Service-specific readiness | Failed check removes readiness. |

Do not mark a required authorization, database, or policy dependency optional merely to make a pod start.

## Operational modes

The same image supports:

- normal serving when `DX_BOOT_MODE` is unset or `serve`;
- rendered configuration validation with `DX_BOOT_MODE=config-check`;
- a single migration actor with `DX_BOOT_MODE=migrate-only`.

An unknown mode is a startup error. GitOps uses the same image for the migration job and deployment, preventing schema/code skew.

## Failure and observability

Configuration errors name the invalid field and stop boot. Dependency readiness is registered automatically. Shutdown first stops accepting HTTP/gRPC work, then cancels and joins workers, then closes infrastructure. A service must test cancellation and drain within its Kubernetes termination grace period.

## Common misuse

- Creating a logger before config, so the configured level is ignored.
- Installing a second signal handler in `main`.
- Calling `os.Getenv` throughout business code.
- Starting a goroutine in `Wire` without registering it.
- Treating a secret as a default.
- Using `Degrade` without a correct code path and an alert.

Verification:

```bash
DX_BOOT_MODE=config-check go run ./cmd/server
go test ./internal/config ./cmd/server
```

Expected: valid rendered configuration exits zero without binding a port or dialing an application dependency; invalid security configuration exits non-zero.

