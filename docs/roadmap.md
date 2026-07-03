---
id: roadmap
title: Roadmap & Milestones
slug: /roadmap
sidebar_label: Roadmap & Milestones
description: The 12-week plan — weekly schedule, milestone gates, and skip-ahead checkpoints.
---

# Roadmap & Milestones

This roadmap turns the curriculum into a **self-paced 12-week plan** at ~8–10 hours/week. It is a default, not a contract: the milestones are what matter, the weeks are suggestions. Full-time learners can roughly halve the calendar.

## The milestones

Each milestone has a **gate** — a concrete "you can now…" statement. Don't move on until you can honestly check the gate; everything after it assumes you can.

| Milestone | Weeks | Gate — "you can now…" |
|---|---|---|
| **M0** Environment & orientation | 1 | Run `make dev-demo` green locally and explain the two-track (Go/Java) architecture in two minutes |
| **M1** Go fundamentals | 2–4 | Write idiomatic Go using interfaces, wrapped errors, and generics; pass the M1 checkpoint |
| **M2** Intermediate Go | 5–7 | Build a concurrent, configurable, structured-logged, tested CLI/worker with clean shutdown |
| **M3** Advanced / service engineering | 8–10 | Build a containerized HTTP service with chi + pgx + RabbitMQ that passes all CI gates |
| **M4** Platform internalization | 11–12 | Pass the GO-SERVICE-STANDARDS self-review; find your way around `dx-common-go` without help |
| **M5** Capstone & first contribution | end of 12 (+1–2 wks) | Capstone service reviewed against the standards checklist; first real PR merged to a service's `dev` |

## Weekly schedule (part-time, ~9 h/week)

| Week | Content | Hours |
|---|---|---|
| 1 | [Environment setup](module-0-setup/environment) · [Platform orientation](module-0-setup/platform-orientation) | 6–8 |
| 2 | [Syntax, variables & types](module-1-go-fundamentals/syntax-variables-types) · [Control flow & functions](module-1-go-fundamentals/control-flow-functions) · [Collections](module-1-go-fundamentals/collections) | 9 |
| 3 | [Structs & methods](module-1-go-fundamentals/structs-methods) · [Pointers](module-1-go-fundamentals/pointers-memory-basics) · [Interfaces](module-1-go-fundamentals/interfaces) | 9 |
| 4 | [Error handling](module-1-go-fundamentals/error-handling) · [Packages & modules](module-1-go-fundamentals/packages-modules) · [Generics](module-1-go-fundamentals/generics) · [Reflection](module-1-go-fundamentals/reflection-and-when-not) | 9 |
| 5 | [Concurrency](module-2-intermediate/concurrency) · [Context](module-2-intermediate/context) | 9 |
| 6 | [Dependency injection](module-2-intermediate/dependency-injection) · [Configuration](module-2-intermediate/configuration) · [Logging](module-2-intermediate/logging) | 8 |
| 7 | [Testing](module-2-intermediate/testing) · [Benchmarking & profiling](module-2-intermediate/benchmarking-profiling) · [Memory & performance](module-2-intermediate/memory-performance) · M2 mini-project | 9 |
| 8 | [Project structure](module-3-advanced/project-structure) · [HTTP & middleware](module-3-advanced/http-servers-middleware) · [REST API development](module-3-advanced/rest-api-development) | 10 |
| 9 | [AuthN & AuthZ](module-3-advanced/authn-authz) · [Database patterns](module-3-advanced/database-patterns) · [Transactions](module-3-advanced/transactions) | 10 |
| 10 | [RabbitMQ & events](module-3-advanced/event-driven-rabbitmq) · [Workers & cron](module-3-advanced/workers-cron) · [Distributed systems](module-3-advanced/distributed-systems) · [Observability](module-3-advanced/observability) · [Containers & K8s](module-3-advanced/containers-kubernetes) · [CI/CD](module-3-advanced/cicd) | 10 |
| 11 | [Architecture deep dive](module-4-platform/architecture-deep-dive) · [Repo & workflow](module-4-platform/repo-structure-workflow) · [dx-common-go tour](module-4-platform/dx-common-go-tour) · [Service anatomy](module-4-platform/service-anatomy) | 9 |
| 12 | [Standards checklist](module-4-platform/standards-checklist) · [Security model](module-4-platform/security-model) · [Testing strategy](module-4-platform/testing-strategy) · [Deployment](module-4-platform/deployment) | 8 |
| 13–14 | [Capstone service](capstone/capstone-service) · [First contribution](capstone/first-contribution) | 10–16 |

Weeks 10 and 12 look dense; several of those pages are shorter reads with hands-on work you've partly done already. Time estimates on each page are the authoritative numbers.

## Skip-ahead checkpoints

You do **not** have to sit through material you already know. Each checkpoint below is a self-test; if you pass it comfortably, skip the module and move on. (Module 0 and Modules 4–5 are mandatory for everyone — they're platform-specific.)

### Checkpoint M1 → skip Module 1 if you can…

- Explain when a method needs a pointer receiver, and what happens when a value with value receivers is stored in an interface.
- Write a function returning `(T, error)`, wrap the error with `fmt.Errorf("...: %w", err)`, and match it upstream with `errors.Is` / `errors.As`.
- Explain the difference between a nil slice, an empty slice, and why appending to a shared slice is dangerous.
- Write a small generic function with a constraint, e.g. `func Map[T, U any](xs []T, f func(T) U) []U`.
- Explain what `internal/` does in a module and what a `replace` directive is for.

### Checkpoint M2 → skip Module 2 if you can…

- Explain why a fire-and-forget `go func()` is banned in this codebase and show three legitimate ways to wait for a goroutine.
- Use `context.Context` correctly: first parameter, cancellation propagation, `context.WithTimeout`, and why you never store a context in a struct.
- Write a table-driven test with subtests and run one case by name.
- Read a CPU profile with `go tool pprof` well enough to find an obvious hot spot.
- Explain constructor-based dependency injection and why the platform doesn't use a DI framework.

### Checkpoint M3 → skip Module 3 if you can…

- Sketch the standard DX middleware order (RequestID → RealIP → Logger → CORS → Compression → Recoverer → Timeout) and explain why Recoverer sits where it does.
- Write a parameterized pgx query, run two statements atomically in a transaction with a deferred rollback, and explain the transactional outbox pattern.
- Declare a durable RabbitMQ queue with a DLX/DLQ and explain at-least-once delivery and why consumers must be idempotent.
- Explain what a readiness probe should check that a liveness probe must not.
- Get a Go service through a multi-stage Dockerfile and describe the platform's PR gate (`go build`, `go test`, `gofmt -l`, `go vet`, `golangci-lint run`).

## Role-based entry points

| You are… | Your path |
|---|---|
| New engineer, little Go | Module 0 → 1 → 2 → 3 → 4 → 5, in order |
| Backend dev (Java/Python/Node), new to Go | Module 0 → Module 1 at express pace (focus: pointers, interfaces, errors, generics — the parts unlike your language) → 2 → 3 → 4 → 5 |
| Experienced Go dev, new to DX | Module 0 → pass all three checkpoints → Module 4 → 5. Skim Module 3's REST/RabbitMQ pages for platform conventions |
| Java-track DX engineer moving to Go | Module 0 (fast — you know the stack) → Module 1 express → 2 → 3 (compare with the Vert.x equivalents as you go) → 4 → 5 |

## Tracking progress

Keep it simple: copy this list into your notes and check things off.

- [ ] M0 gate: `make dev-demo` green, can explain the architecture
- [ ] M1 checkpoint passed
- [ ] M2 mini-project done (concurrent worker CLI)
- [ ] M3 service built and containerized, CI gates pass
- [ ] M4 standards self-review passed
- [ ] M5 capstone reviewed
- [ ] First PR merged 🎉

Your onboarding buddy or team lead reviews the M3 service, the capstone, and your first PR — those three reviews are where most of the feedback (and learning) happens.
