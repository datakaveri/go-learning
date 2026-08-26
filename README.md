# CDPG Go Developer Portal

Canonical architecture, service-building tutorials, operational standards, and progressive Go curriculum for CDPG contributors.

The developer guide covers the Control Plane, Data Plane, Agentic Plane, gateway, identity, ACL, authorization, OpenFGA, planned OPA, shared Go platform, persistence, messaging, testing, local integration, and GitOps. `examples/dx-example-go` is a compile-tested teaching service.

## Documentation paths

- Developer Guide: architecture, new-service tutorial, integrations, standards, and operations.
- Go Curriculum: language through production service engineering.
- Teaching Service: `examples/dx-example-go`.

## Curriculum

| Module | Focus |
|---|---|
| 0 | Go 1.25 toolchain and platform orientation |
| 1 | Language fundamentals |
| 2 | Context, concurrency, injection, config, logging, testing, performance |
| 3 | HTTP, security, persistence, events, workers, observability, containers, delivery |
| 4 | Current Data Exchange architecture and dx-common-go platform APIs |
| 5 | A complete service capstone and first contribution |

Early exercises are standalone. Later exercises use the real service source and local stack.

## Build

Requires Node.js 20 or newer.

~~~bash
npm ci
npm run typecheck
npm run build
npm run start
~~~

Validate the teaching service from the CDPG workspace:

~~~bash
cd examples/dx-example-go
go test ./...
go vet ./...
go test -race ./...
~~~

The local URL is http://localhost:3000/go-learning/.

## Content rules

- Verify platform claims against source.
- Use [cdpg-docs](https://datakaveri.github.io/cdpg-docs/) for architecture and operations reference.
- Use [dx-common-go-docs](https://datakaveri.github.io/dx-common-go-docs/) for exact SDK APIs and version policy.
- Use the portal status vocabulary and never present planned behavior as operational.
- Document only the Go platform architecture; do not retain retired platform/trust material for context.
- Keep Go examples idiomatic, formatted, context-aware, and testable.
- Every lesson ends with an exercise and a self-check.
