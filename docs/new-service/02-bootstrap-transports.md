---
title: 2. Bootstrap, configuration, HTTP, and gRPC
description: Compose a service with verified shared platform APIs and transport contracts.
---

# 2. Bootstrap, configuration, HTTP, and gRPC

**Status: HTTP/bootstrap implemented; internal gRPC partially implemented.** The code shapes shown for current packages are source-verified. A generated service registrar remains contract-specific.

## Step 1 — typed configuration

Purpose: make an invalid environment fail before it accepts traffic.

Files: `internal/config/config.go`, `configs/config.yaml`, config tests.

Embed `platform/config.Base`, add only service dependencies, and implement `Validate`. Include PostgreSQL, workload verifier, workload issuer, current authorization client, broker/cache only when needed. Do not add secrets to YAML.

```go
type Config struct {
    platformconfig.Base `mapstructure:",squash"`
    Postgres            dxsql.Config             `mapstructure:"postgres"`
    WorkloadVerifier    workload.VerifierConfig  `mapstructure:"workload_verifier"`
    WorkloadIssuer      issuer.Config            `mapstructure:"workload_issuer"`
    Authorization       fga.Config               `mapstructure:"authorization"`
}
```

Common mistake: duplicating `server`, `openapi`, or log fields already in `config.Base`; using a secret default; allowing empty workload enforcement.

Verification:

```bash
DX_BOOT_MODE=config-check go run ./cmd/server
```

Expected: rendered valid settings exit zero without serving; missing destination credentials or database settings fail with a named field.

## Step 2 — composition root

Purpose: declare dependencies while the platform owns order and shutdown.

Files: `cmd/server/main.go`, embedded migration package.

```go
bootstrap.Run(bootstrap.Spec[serviceconfig.Config]{
    Name: "dx-example-go",
    Config: serviceconfig.Options(),
    Deps: func(cfg *serviceconfig.Config) bootstrap.Deps {
        return bootstrap.Deps{
            Migrations: bootstrap.Migrations(migrations.FS, ".", "schema_migrations_example"),
            Postgres:   bootstrap.Required(cfg.Postgres),
        }
    },
    Wire: wire,
})
```

`wire` creates adapters, application service, handlers, supervised work, and router. It contains no business decision.

Common mistake: creating a pool/logger/signal handler manually, or marking a required database/PDP optional.

Verification: unit-test constructors; run config-check and migrate-only roles. Expected: one migration actor and one coordinated lifecycle.

## Step 3 — HTTP contract and handler

Purpose: expose public behavior without leaking router or datastore concerns into use cases.

Files: `api/openapi/openapi.yaml`, `internal/transport/http/handler.go`, `routes.go`.

```go
type CreateRequest struct {
    httpx.Actor
    Name string `json:"name" validate:"required,min=1,max=120"`
}

func (h *Handler) Create(ctx context.Context, req CreateRequest) (httpx.Created[Response], error) {
    widget, err := h.app.Create(ctx, application.CreateCommand{
        Subject: req.Subject,
        Name:    req.Name,
    })
    if err != nil {
        return httpx.Created[Response]{}, err
    }
    return httpx.Created[Response]{Value: toResponse(widget)}, nil
}
```

```go
func Routes(h *Handler) httpx.RouteSet {
    return httpx.Routes("/widgets",
        httpx.POST("/", httpx.Handle(h.Create), httpx.OpID("createWidget")),
        httpx.GET("/{id}", httpx.Handle(h.Get), httpx.OpID("getWidget")),
    )
}
```

Use `HandleRaw` only for a native standard, webhook, file, redirect, or stream. For an anonymous-or-authenticated read, use `Optional()` plus `HandleOptional`; an invalid supplied token remains an error.

Verification: direct handler tests plus router/OpenAPI drift tests. Expected: binding/validation errors never invoke the application service; every route has a matching operation ID.

## Step 4 — identity middleware

Create the workload verifier from config. Build the router with `middleware.Resolve`, trusted subject projection, the verifier, and direct JWT config only when direct access is an approved deployment path. The gateway-facing destination must use `workload_verifier.enforcement: required` and allow only intended callers/subject asserters.

Security warning: middleware order matters—workload verification precedes represented-subject resolution. Never read internal identity headers manually.

Verification: workload/subject-asserter negative matrix. Expected: absent/invalid/wrong-audience/disallowed caller requests stop before the handler.

## Step 5 — internal gRPC when required

Purpose: provide an explicit internal contract rather than another public-style HTTP endpoint.

Files: `api/proto/example/v1/example.proto`, generated code, `internal/transport/grpc/server.go`, caller contract tests.

```proto
syntax = "proto3";
package cdpg.example.v1;
option go_package = "github.com/datakaveri/dx-example-go/api/proto/example/v1;examplev1";

service ExampleService {
  rpc GetWidget(GetWidgetRequest) returns (GetWidgetResponse);
}
```

Register the generated service through `platform/grpc/server.New` and return it from `bootstrap.Spec.GRPC`. Apply workload credentials with the destination audience at callers. The shared surface currently provides unary interception; a streaming RPC needs a reviewed platform extension.

Common mistake: adding internal HTTP because one current service has it, or serving gRPC without workload interceptors/deadlines.

Verification: protobuf generation/compatibility, in-process client/server test, wrong-audience and deadline tests. Expected: only allowlisted workloads reach the method.

## Checkpoint

The service should now build, validate config, serve health plus protected typed HTTP, and optionally serve authenticated unary gRPC. It still must not be exposed before persistence, events, authorization, and negative tests are complete. Next: [persistence and events](./03-persistence-events.md).

