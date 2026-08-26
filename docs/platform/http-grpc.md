---
title: HTTP, gRPC, responses, and errors
description: Use the shared transport adapters, middleware, validation, paging, and error mapping.
---

# HTTP, gRPC, responses, and errors

**Status: HTTP implemented; gRPC partially implemented.** Shared gRPC server/client foundations and unary interceptors exist, while streaming and fleet adoption remain incomplete.

## Typed HTTP handlers

`platform/http` handlers use `func(context.Context, Request) (Response, error)`. The adapter binds path/query/header/body fields, obtains the verified actor, validates, invokes, and renders.

```go
type GetRequest struct {
    httpx.Actor
    ID string `path:"id" validate:"required,uuid"`
}

func (h *Handler) Get(ctx context.Context, req GetRequest) (WidgetResponse, error) {
    widget, err := h.service.Get(ctx, req.Subject, req.ID)
    if err != nil {
        return WidgetResponse{}, err
    }
    return toResponse(widget), nil
}

routes := httpx.Routes("/widgets",
    httpx.GET("/{id}", httpx.Handle(handler.Get), httpx.OpID("getWidget")),
)
```

Use `HandleVoid` for no-content operations, `HandleOptional` with an embedded `OptionalActor` for anonymous-or-verified reads, and `HandleRaw` for blobs, redirects, SSE, OGC/NGSI-LD representations, MCP payloads, or webhooks whose format CDPG does not own.

## Router behavior

`httpx.NewRouter` installs recovery, tracing, request ID, trusted client IP handling, structured request logging, CORS, route timeouts, response compression, health, metrics, and docs. Routes are data, enabling OpenAPI drift tests.

```go
router := httpx.NewRouter(httpx.RouterSpec{
    Base:    "/iudx/example/v1",
    URNs:    httpx.URNSpace("example"),
    Health:  app.Health,
    Logger:  app.Log,
    Auth:    httpx.AuthSpec{Authenticate: authMiddleware},
    Timeout: 15 * time.Second,
}, routes)
```

Authentication middleware is injected; `platform/http` does not decide what makes a credential valid. `Public()` bypasses identity. `Optional()` resolves a valid credential when supplied and rejects an invalid one. `Streaming()` excludes a route from request timeout and response compression.

## Response and error conventions

Owned JSON endpoints use the shared success envelope. Platform error types map to stable HTTP/gRPC status and problem responses. Domain-specific mappings are registered with `WithMappers` or router `Mappers`; unclassified failures become 500 without leaking details.

Use `platform/paging` for bounded request/result types. The service still owns stable ordering and cursor semantics.

## Request validation

Use struct tags for syntactic validation and the domain/application layer for semantic invariants. OpenAPI request validation protects the published contract; it does not replace object authorization or domain validation.

## Internal gRPC

Build the server through `platform/grpc/server.New`, register generated services, attach workload verification, and return it through `bootstrap.Spec.GRPC`. Clients use the platform gRPC package to attach destination-bound workload credentials and telemetry.

```go title="Illustrative wiring; generated registrar omitted"
GRPC: func(app *bootstrap.App[serviceconfig.Config]) (bootstrap.Servable, error) {
    return grpcserver.New(app.Cfg.GRPC, grpcserver.Options{
        Log:      app.Log,
        Workload: verifier,
    }, registrar)
},
```

The exact generated registrar is service-owned. Mark a code example as pseudocode if its protobuf contract has not been created.

## Failure and testing

- Bad binding or validation: stable client error; handler is not invoked.
- Missing/invalid required identity: 401/`Unauthenticated`.
- Denial: 403/`PermissionDenied` with safe reason code.
- Panic: recovered and correlated; details remain internal.
- Client cancellation: downstream calls stop and no detached goroutine remains.
- Stream write after headers: record the failure; do not try to replace the response.

Tests call typed handlers directly, exercise the router for middleware and binding, compare declarative routes to OpenAPI, and exercise gRPC through an in-memory or local listener with real interceptors.

## Common misuse

- Reading arbitrary identity headers inside a handler.
- Returning bare JSON for an owned API because it is convenient.
- Using `HandleRaw` to bypass validation/error mapping.
- Adding service-specific global middleware without documenting order.
- Giving an internal gRPC call no deadline or wrong destination audience.
- Adding streaming RPCs before workload and telemetry streaming interceptors exist.

