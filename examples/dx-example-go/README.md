# dx-example-go

`dx-example-go` is a compile-tested teaching service for the CDPG Go developer portal. It is not a registered platform service and must not be deployed as one.

The example demonstrates domain/application seams, typed HTTP handlers, current authorization checks, platform SQL transactions, a transactional outbox, destination-bound workload identity, declarative bootstrap, and coordinated shutdown. Local unit tests need no infrastructure. Running the server needs PostgreSQL, RabbitMQ, Keycloak workload configuration, and `dx-authz-go`.

```bash
go test ./...
go test -race ./...
go vet ./...
```

The local `replace` expects this repository at `cdpg/go-learning/examples/dx-example-go` and `dx-common-go` at `cdpg/dx-common-go`. Remove it and pin a released module version in a real service.

