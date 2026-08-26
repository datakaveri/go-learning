---
title: Environment
sidebar_label: Environment
description: Install Go 1.25, verify the toolchain, and run the local Data Exchange stack.
---

# Environment

## Outcomes

You will install the toolchain, run a Go test, start the platform, and know where each repository lives.

## Toolchain

Install:

- Go 1.25 or newer;
- Git;
- Docker Desktop with Docker Compose v2;
- GNU Make;
- an editor with gopls.

Verify:

~~~bash
go version
go env GOMOD GOWORK GOPATH
gofmt -h
go test -h
docker compose version
git --version
make --version
~~~

Configure the editor to use gopls, format on save, and organize imports. Do not install a separate formatter that fights gofmt.

## First Go module

~~~bash
mkdir hello-dx
cd hello-dx
go mod init example.com/hello-dx
go mod edit -go=1.25
~~~

Create a function and table-driven test, then run:

~~~bash
gofmt -w .
go test ./...
go vet ./...
~~~

## Platform workspace

In the orchestration repository:

~~~bash
make dev-clone
make dev-up
make dev-init-dbs
make dev-demo
~~~

Service repositories are cloned inside the workspace because the Compose build context needs them together. They remain independent Git repositories.

Useful daily commands:

~~~bash
make dev-pull
make dev-status
make dev-logs SVC=dx-gateway-go
make dev-token
make dev-down
~~~

## Checkpoint

- What is the difference between GOPATH, GOMOD, and GOWORK?
- Why does gofmt have no project-specific style configuration?
- Can you run make dev-demo and identify the first failing component from logs?
- Can you show which Git repository owns a service file?

Do not continue until go test and the local platform health checks run successfully.
