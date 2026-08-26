---
title: CDPG Go style and engineering rules
description: Enforceable conventions for secure, maintainable, observable CDPG Go services.
---

# CDPG Go style and engineering rules

**Status: Implemented as the standard; fleet conformance is partial.** `gofmt`, build, vet, tests, race tests, security checks, and contract checks enforce behavior; prose alone is not a gate.

## API and package design

- Package names are short, lowercase, singular where natural, and describe a capability. Avoid `util`, `common`, `helpers`, and `misc`.
- File names are lowercase with underscores only when readability improves (`access_request.go`). Tests end in `_test.go`; generated files carry the generator header.
- Export the smallest stable surface. Accept interfaces and return concrete types unless substitution is a real consumer need.
- Declare interfaces at the consumer boundary. Keep them narrow enough to fake by hand.
- Constructors validate required dependencies. Functional options are for optional behavior with safe defaults, not required fields.

Preferred:

```go
type GrantStore interface {
    ByID(context.Context, string) (Grant, error)
}

func NewService(store GrantStore, authorizer Authorizer) (*Service, error) {
    if store == nil || authorizer == nil {
        return nil, errors.New("grant service: dependencies are required")
    }
    return &Service{store: store, authorizer: authorizer}, nil
}
```

Avoid a package-wide service locator, an interface that mirrors every repository method, or a constructor that succeeds with a nil security dependency.

## Context, errors, and telemetry

- `context.Context` is the first parameter and is never stored in a struct.
- Pass the request context to all network, database, lock, and wait operations.
- Create errors with actionable context and wrap the cause with `%w`.
- Use `errors.Is`/`errors.As`; do not inspect error strings.
- Handle an error once. Return it or log it, not both, unless adding an explicit boundary event.
- Use current platform structured logging (`zap`) and OpenTelemetry/Prometheus integrations. Do not introduce a second logger into a service.
- Metrics describe behavior, not high-cardinality identities. Put IDs in traces or structured logs under the data-handling policy.

Preferred:

```go
widget, err := repo.ByID(ctx, id)
if err != nil {
    return Widget{}, fmt.Errorf("get widget %q: %w", id, err)
}
```

Avoid `log.Printf`, swallowed errors, raw error text in an HTTP response, or subject/resource IDs as Prometheus labels.

## Concurrency and lifecycle

- The creator owns a goroutine’s cancellation and completion.
- Use `errgroup`, `bootstrap.App.Go`, `App.Background`, or `platform/executor`; never fire-and-forget from a handler.
- Keep channels directional, close them only from the sender, and document ownership.
- Bound worker concurrency, queue capacity, retries, and per-item time.
- Protect shared state with the simplest appropriate primitive. Do not hold a lock during network I/O.
- Run `go test -race ./...`; add `goleak` checks for components that own goroutines.

## Database and time

- Parameterize every value. Interpolate only compile-time allowlisted identifiers.
- Put transaction boundaries in application use cases and use `sql.Manager.Do`/`DoRetry`.
- Use `SELECT … FOR UPDATE` for conflicting state transitions; use advisory locks or durable leases for singleton ownership.
- Map no-row, uniqueness, serialization, and connection errors to stable domain/platform errors.
- Represent nullable database values explicitly. Avoid sentinel strings or zero timestamps.
- Store and compare UTC; inject a clock where time controls a business decision. Serialize RFC 3339 timestamps.
- IDs are opaque. Validate UUID/URN syntax at the edge and never derive tenant authorization from string shape alone.

## HTTP, gRPC, JSON, and protobuf

- Use typed platform handlers and declarative routes. `HandleRaw` is reserved for standards-native or streaming responses.
- Keep handlers thin: bind, call, map.
- JSON field names are stable `lowerCamelCase`; reject unknown fields where the contract requires strictness.
- Protobuf fields are additive; never reuse a number. Reserve deleted fields and names.
- Set request/body limits and timeouts. A stream must be explicitly marked and cancellation-aware.

## Configuration and secrets

- Embed `config.Base`, apply defaults, bind file then environment, validate at startup.
- Environment keys are unprefixed platform keys such as `SERVER_PORT` and `POSTGRES_DSN`; `DX_BOOT_MODE` is the explicit operational-role exception.
- No secret has a usable default. Secrets come from the deployment secret mechanism, not repository YAML or command examples.
- Never log tokens, credentials, approval codes, private keys, presigned URLs, or unredacted protected attributes.
- Use constant-time comparisons for secret material when a library does not already verify it.

## Security

- Authenticate at the edge and verify workload identity again at every service boundary.
- Treat represented subject and workload as separate identities.
- Authorize every object/state transition; role membership alone is not resource access.
- Default deny on missing mapping, invalid decision, unsupported obligation, unavailable required dependency, or expired context.
- Validate all external input, normalize before authorization, parameterize datastore queries, and constrain paths/bucket keys to owned namespaces.
- Use vetted cryptographic and JWT libraries. Pin allowed algorithms, issuer, audience, time bounds, and key source.

## Testing and dependencies

- Use table-driven unit tests, explicit fixtures/builders, and fakes at consumer interfaces.
- Prefer integration tests with real PostgreSQL, Redis, RabbitMQ, OpenFGA, Elasticsearch, or S3-compatible infrastructure when adapter semantics matter.
- Fuzz parsers, validators, URNs, route matching, obligation translators, and event decoders.
- Mock behavior you own, not an entire third-party SDK surface.
- Keep dependencies few, maintained, licensed, pinned, scanned, and justified. Shared platform capability wins over a new one-off library.
- Generated code is reproducible in CI and never manually patched.

## Required commands

```bash
gofmt -w .
go build ./...
go vet ./...
go test ./...
go test -race ./...
govulncheck ./...
```

Repositories may add stricter linters, OpenAPI/protobuf generation checks, integration suites, and container scans. They may not remove these gates without a recorded exception.

