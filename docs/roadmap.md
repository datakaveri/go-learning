---
title: Curriculum Roadmap
description: Modules, milestones, pacing, and evidence of completion.
---

# Curriculum roadmap

The suggested pace is twelve weeks, but milestones—not calendar time—determine readiness.

Experienced Go engineers can instead follow the role-based [learning paths](start/learning-paths.md) and complete the [new-service tutorial](new-service/quick-start.md).

| Module | Focus | Suggested pace | Completion evidence |
|---|---|---:|---|
| 0 | Toolchain and platform map | 2–3 days | Local stack health and a correct request-flow explanation |
| 1 | Go fundamentals | 2–3 weeks | Small module with tests, interfaces, wrapped errors, and generics |
| 2 | Production Go | 2 weeks | Context-aware concurrent component with config, logs, tests, and profile |
| 3 | Service engineering | 3 weeks | Layered HTTP service with SQL, events, security, telemetry, container |
| 4 | Data Exchange platform | 2 weeks | Source walkthrough and architecture review of one service |
| 5 | Capstone and contribution | 1–2 weeks | Verified service plus reviewable contribution proposal or PR |

## Milestone 1 — Go fluency

You can explain value versus pointer receivers, choose slices and maps safely, define consumer-owned interfaces, wrap errors with %w, organize modules, and use a generic type without hiding behavior.

Verification:

~~~bash
go test ./...
go test -race ./...
go vet ./...
~~~

## Milestone 2 — Production Go

You can propagate context, stop goroutines, inject dependencies without globals, load validated config, emit structured logs, isolate tests, benchmark a claim, and read a profile before optimizing.

## Milestone 3 — Service engineering

You can:

- declare typed platform/http routes and handlers;
- model client-safe platform errors and paging;
- use platform/database/sql with transaction context;
- use platform/events and an outbox for atomic publication;
- separate mandatory and degraded dependencies;
- expose health, metrics, and trace-aware logs;
- build a secure, non-root container and render Kubernetes resources.

## Milestone 4 — Platform contributor

You can trace a request gateway → authorization → service → store → event consumers, identify data and lifecycle ownership, find the correct `dx-common-go` seam, distinguish authentication from authorization, and explain current OpenFGA checks versus planned OPA/carried decisions.

## Milestone 5 — Ship-ready

Your capstone has:

- domain/application/adapter/composition separation;
- OpenAPI-aligned route declarations;
- unit, integration, contract, race, and security evidence;
- migration and rollback notes;
- deployment values and an operator checklist.
- the completed [service review scorecard](standards/review-scorecard.md) with no blocker.

## Learning discipline

- Type every exercise; do not paste a finished solution.
- Keep a decision log: what you chose, alternatives, and evidence.
- Review race and failure behavior before style.
- Link to reference docs rather than copying their API tables into notes.
- Re-run the module checkpoint after any long break.
