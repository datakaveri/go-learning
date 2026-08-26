---
title: Packages, Modules, and Workspaces
sidebar_label: Packages & Modules
description: Package boundaries, Go modules, pinned dependencies, workspaces, and multi-repository development.
---

# Packages, modules, and workspaces

## Package boundaries

A directory is a package. Exported names start with an uppercase letter. Keep packages cohesive and name them for what they provide, not generic buckets such as utils, common, helpers, or misc.

Avoid import cycles by preserving dependency direction: domain → ports → application, with adapters and composition importing inward.

## Module basics

~~~bash
go mod init example.com/widgets
go mod edit -go=1.25
go get github.com/datakaveri/dx-common-go@<approved-commit>
go mod tidy
go list -m all
~~~

go.mod declares the module path, Go language version, and dependencies. go.sum authenticates downloaded module content; it is not a lockfile but belongs in version control.

## Workspaces for local multi-repo changes

~~~bash
cd /path/to/cdpg
go work init ./dx-common-go ./dx-widget-go
go work sync
go env GOWORK
~~~

A go.work file lets local modules resolve together without committing a replace directive. Do not let an untracked workspace hide the version used in CI or a release build.

dx-common-go currently has no published tags. Services pin an approved commit or pseudo-version. See the [SDK version policy](https://datakaveri.github.io/dx-common-go-docs/versions).

## Internal packages

Code under internal can be imported only by packages within its parent tree. Service implementations live under internal so another repository cannot accidentally depend on implementation details. Public cross-service contracts belong in a deliberate module or protocol, not a service's internal package.

## Commands

~~~bash
go list ./...
go mod why github.com/datakaveri/dx-common-go
go mod graph
go mod verify
go work use ./dx-another-service-go
~~~

## Exercise

Create three packages: domain, service, and memory adapter. Make cmd/server wire them. Then try to introduce an import from domain to the adapter, observe the design problem, and remove it.

## Check yourself

- What is the difference between a package, module, and workspace?
- Why should release evidence record the module pseudo-version?
- What does internal protect?
- When is a committed replace directive inappropriate?
