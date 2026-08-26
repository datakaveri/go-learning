---
title: Go service architecture standard
description: Required service layout, dependency direction, ownership, lifecycle, and test seams.
---

# Go service architecture standard

**Status: Partially implemented across the fleet; required for new services.** Several current repositories use `internal/api` and `internal/service` as names for transport and application layers. New services use the clearer layout below. Those names are a target convention, not a claim that every current service already has them.

![Required Go service layers](/img/architecture/service-layering.svg)

## Reference layout

```text
cmd/server/                 composition root only
internal/transport/http/    public or service REST adapters
internal/transport/grpc/    internal protobuf adapters
internal/application/       use cases and application ports
internal/domain/            entities, value objects, invariants, domain errors
internal/repository/        persistence adapters; interfaces normally live at consumers
internal/integration/       clients for other services and external systems
internal/worker/            event consumers, outbox dispatchers, scheduled work
internal/config/            typed service configuration and validation
db/migrations/              embedded, versioned service-owned DDL
api/openapi/                public REST contract
api/proto/                  internal gRPC contract
configs/                    non-secret local defaults
tests/                      integration, contract, security, and end-to-end fixtures
```

## Dependency rule

| Layer | May depend on | Must not depend on |
|---|---|---|
| Domain | Standard library and small domain value packages | HTTP, gRPC, SQL drivers, brokers, caches, platform bootstrap |
| Application | Domain and consumer-owned interfaces | Concrete database, router, broker, Keycloak, OpenFGA, or vendor clients |
| Transport | Application ports, request/response types, platform transport adapters | Repository implementation or direct datastore queries |
| Repository/integration | Application/domain ports and infrastructure SDKs | Transport packages or cross-service business orchestration |
| Worker | Application ports, event contracts, lifecycle primitives | Unbounded goroutines, direct writes that bypass use cases |
| Config | Typed values and validation | Live clients or side effects |
| `cmd/server` | Every adapter needed to compose the service | Business rules |

Dependencies point inward. A domain package is useful in a unit test with no network, database, environment, or framework. An application use case is testable with small fakes declared at its consumer boundary.

## Layer responsibilities

### Domain

Own invariants, state transitions, value validation, and domain errors. Use constructors that return errors when invalid state is possible. Do not add JSON, SQL, OpenFGA, or HTTP meanings to a domain type merely for adapter convenience.

### Application

Orchestrate a complete use case: load state, authorize the object/action, call domain behavior, commit changes and events, and return a transport-neutral result. The application layer owns transaction boundaries because it knows which operations form one business action.

```go
type WidgetRepository interface {
    Create(context.Context, domain.Widget) error
    ByID(context.Context, string) (domain.Widget, error)
}

type TransactionManager interface {
    Do(context.Context, func(context.Context) error, ...dxsql.TxOption) error
}
```

The interface is declared where it is consumed. A repository does not publish a broad interface for every potential caller.

### Adapters

HTTP and gRPC adapters bind and validate transport input, obtain the verified subject from context, call one application operation, and map output. Repositories translate domain operations to parameterized datastore calls. Integrations translate an owned application port to a remote contract. None of these layers decides business policy independently.

### Composition root

`cmd/server/main.go` declares `bootstrap.Spec`, dependencies, wiring, workers, and closers. `bootstrap.Run` owns config load, logger creation, dependency readiness, HTTP/gRPC serving, signals, ordered drain, worker cancellation, and infrastructure close. Do not install a second signal handler.

## Transactions and events

The application layer calls `platform/database/sql.Manager.Do`. Repositories obtain the ambient transaction through the context. A domain change and its `platform/events.Outbox` write occur inside the same callback. After commit, a supervised dispatcher publishes the event; consumers remain idempotent because delivery is at least once.

## Configuration and lifecycle

- Embed `platform/config.Base` in the service config.
- Use `config.Load` or an explicit `Spec.Load` wrapper when compatibility aliases genuinely exist.
- Validate all required and security-sensitive values at startup.
- Declare required dependencies with `bootstrap.Required`; use `bootstrap.Degrade` only when the service has a correct documented degraded mode.
- Register fail-fast essential loops with `App.Go`, reconnecting broker consumers with `App.Background`, cleanup with `App.Closer`, and extra readiness checks with `App.Probe`.
- Use `platform/executor` for lifecycle-owned per-request background work. Never launch an unmanaged goroutine from a handler.

## Test seams

Test the domain without mocks, application use cases with small fakes, HTTP/gRPC adapters as contract boundaries, and infrastructure adapters against real disposable dependencies. Run migrations from empty and from the previous version. Verify cancellation and shutdown, not just happy-path output.

## Current examples

- `dx-registry-go` demonstrates declarative bootstrap and a conventional Control Plane service.
- `dx-catalogue-go` demonstrates a current internal gRPC surface and Elasticsearch integration.
- `dx-audit-go` demonstrates a supervised event consumer.
- `dx-files-connect-api-go` demonstrates object storage and worker lifecycles, while also showing why mixed older and platform packages must not be copied wholesale.
- `dx-dataplane-ogc-go` demonstrates `HandleRaw` for OGC-owned representations.

Use these repositories as evidence for a specific pattern, not as templates to copy in full. The [new-service tutorial](../new-service/quick-start.md) is the normative composition.

