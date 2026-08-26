---
title: 1. Choose the boundary and scaffold
description: Establish ownership, contracts, repository, module, and target service layout.
---

# 1. Choose the boundary and scaffold

**Status: Implemented process.** The target directory layout is required for new services; existing fleet adoption is partial.

## Step 1 — write the service charter

Purpose: prove a new bounded context is needed before creating another deployment and datastore.

Create an ADR/service charter that answers:

- What business capability has exactly one owner?
- Which current service was checked to rule out duplication?
- Which actors invoke it and which operations are public?
- What data does it own? What data must it obtain by API/event?
- Which domain facts does it publish and consume?
- Which object/action/relation vocabulary will authorization use?
- What are its availability, consistency, privacy, and audit requirements?
- Is synchronous gRPC actually required, or is an event/projection better?
- Which target features are planned rather than needed for the first slice?

Common mistake: defining a service around a database table or framework capability. The boundary follows business ownership.

Verification: architecture and data owners approve the charter, service inventory name, and ownership map. Expected: no duplicated owner and no cross-service database join.

## Step 2 — allocate contracts and a port

Public APIs use REST/OpenAPI through the gateway. Internal synchronous APIs target gRPC/protobuf. Events carry domain facts.

Before selecting any host/container port, update the orchestration repository using `claude-docs/PORTS.md` as the authority. Do not reuse a port shown in a tutorial. Reserve:

- container HTTP port (normally the platform default unless the fleet registry says otherwise);
- optional internal gRPC port;
- local host port only when direct developer access is needed.

Files:

```text
api/openapi/openapi.yaml
api/proto/example/v1/example.proto   # only if an internal synchronous API exists
```

Define operation IDs and resource/action mapping before implementation. A planned Data Plane decision artifact or OPA input is not a service-owned contract.

Verification: OpenAPI lint and protobuf compatibility checks pass. Expected: every operation has one owner and one authentication/authorization posture.

## Step 3 — create the repository and module

```bash
mkdir dx-example-go
cd dx-example-go
git init
go mod init github.com/datakaveri/dx-example-go
go get github.com/datakaveri/dx-common-go@<reviewed-version>
```

During workspace development only:

```go title="go.mod — local workspace convenience"
replace github.com/datakaveri/dx-common-go => ../dx-common-go
```

Prefer a root `go.work` when multiple local modules are edited together. Never merge a machine-specific absolute `replace` or publish with a local replacement.

Verification:

```bash
go mod tidy
go mod verify
```

Expected: dependencies resolve from an approved source and no credentials appear in module/proxy configuration.

## Step 4 — create the layout

```text
dx-example-go/
├── cmd/server/main.go
├── internal/
│   ├── transport/http/
│   ├── transport/grpc/
│   ├── application/
│   ├── domain/
│   ├── repository/postgres/
│   ├── integration/authz/
│   ├── worker/
│   └── config/
├── db/migrations/
├── api/openapi/
├── api/proto/
├── configs/config.yaml
├── tests/
├── Dockerfile
├── go.mod
└── README.md
```

Files may be absent until the service needs that capability. Do not create empty layers or a generic `util` package.

## Step 5 — create domain and application seams

The teaching example uses a `Widget` only to demonstrate structure.

```go
type Repository interface {
    Create(context.Context, domain.Widget) error
    ByID(context.Context, string) (domain.Widget, error)
}

type Authorizer interface {
    Check(context.Context, Check) (Decision, error)
}
```

The application package declares these interfaces because it consumes them. A PostgreSQL adapter and authorization client implement them later.

Test first:

```go
func TestCreateDeniedDoesNotWrite(t *testing.T) {
    repo := &fakeRepository{}
    service := application.New(repo, denyAuthorizer{})

    _, err := service.Create(context.Background(), validCommand())

    require.ErrorIs(t, err, platformerrors.ErrForbidden)
    require.Equal(t, 0, repo.createCalls)
}
```

Verification: `go test ./internal/domain ./internal/application`. Expected: invalid state and denied authority cannot reach a repository.

## Checkpoint

Proceed only when the charter, API/event ownership, authorization mapping, port reservation process, domain invariants, and application seams are reviewed. Next: [bootstrap and transports](./02-bootstrap-transports.md).

