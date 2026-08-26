---
title: Repositories and Workflow
sidebar_label: Repos & Workflow
description: Multi-repository ownership, orchestration layout, upstream-first shared changes, and review flow.
---

# Repositories and workflow

## Workspace layout

~~~text
cdpg/
  docker-compose.yml
  docker-compose.go-stack.yml
  Makefile
  scripts/
  dx-common-go/
  dx-gateway-go/
  dx-authz-go/
  dx-acl-go/
  dx-catalogue-go/
  ... other service repositories ...
  dx-gitops/
  cdpg-docs/
  dx-common-go-docs/
  go-learning/
~~~

The outer repository owns orchestration. Each nested service and documentation directory is its own Git repository and must be changed, reviewed, and released independently.

## Start a work session

~~~bash
make dev-pull
git -C dx-catalogue-go status --short --branch
git -C dx-common-go status --short --branch
~~~

Preserve unrelated local changes. Confirm which repository owns the file before editing.

## Upstream-first shared changes

If a service needs a missing reusable seam:

1. prove at least two consumers share the semantic need;
2. design and test the narrow dx-common-go API;
3. merge and publish an approved source reference;
4. update one service and verify the complete flow;
5. expand adoption in separate focused changes.

Service code must not temporarily copy the intended SDK implementation across repositories.

## Adding a service

Register:

1. repository and canonical structure;
2. dx-common-go dependency and source pin;
3. Compose build, environment, health, and dependency wiring;
4. gateway route and authorization;
5. database/schema and migrations;
6. GitOps values and ApplicationSet discovery;
7. port and service documentation;
8. smoke and contract tests.

## Exercise

Pick a small cross-cutting improvement. Write the SDK PR scope, first service adoption scope, merge order, compatibility evidence, and rollback. Identify which documentation site owns each change.

## Check yourself

- Does the outer orchestrator commit service source?
- Which change merges first for a new platform API?
- Where is a gateway path registered?
- Why are cross-repository changes separate PRs?
