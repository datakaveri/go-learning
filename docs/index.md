---
id: index
title: Start Here
slug: /
sidebar_label: Start Here
description: The standard Go onboarding curriculum for engineers joining the Data Exchange platform.
---

# DX Go Learning Path

**Welcome.** This is the standard onboarding curriculum for every Go developer joining the **Data Exchange (DX) platform** — whether you're a new joiner or an engineer from another team who will be contributing to the Go service fleet.

The entire platform is built in Go: **15+ microservices** behind a single gateway, sharing one library (`dx-common-go`), one API contract, one security model, and one set of coding standards. This curriculum teaches you the Go language *and* those platform patterns, in an order where each topic builds on the previous one.

## What you'll be able to do at the end

- Write idiomatic, reviewable Go that passes the platform's CI gates (`go vet`, `golangci-lint`, `gofmt`).
- Build a production-shaped microservice: chi router, standard middleware, OpenAPI-validated REST API, pgx-backed persistence, RabbitMQ events, health checks, and metrics.
- Navigate the DX architecture — gateway (PEP), authorization (PDP), policy (PAP), and the domain services — and know where any given change belongs.
- Use `dx-common-go` fluently instead of reinventing what it already provides.
- Land pull requests that pass review on the first or second round.

## How this curriculum works

The path is organized into **six modules**, designed for **10–12 weeks part-time** (~8–10 hours/week). Full-time learners typically finish in 5–6 weeks.

| Module | Focus | Exercises |
|---|---|---|
| [0 — Setup & Orientation](/category/module-0-setup-orientation) | Toolchain + run the real stack locally | Guided setup |
| [1 — Go Fundamentals](/category/module-1-go-fundamentals) | The language, taught with platform idioms | Standalone programs |
| [2 — Intermediate Go](/category/module-2-intermediate-go) | Concurrency, testing, config, performance | Standalone programs |
| [3 — Advanced Go](/category/module-3-advanced-go) | Microservice engineering end to end | **Against the local stack** |
| [4 — The DX Platform](/category/module-4-the-dx-platform) | Everything platform-specific | **Against the real repos** |
| [5 — Capstone](/category/module-5-capstone) | Build a service + first real PR | The real thing |

Every topic page follows the same structure: **learning objectives → prerequisites → time estimate → concepts → platform connection → exercises → mini-project → check yourself → references.** The *Platform connection* box on each page shows you where the topic you just learned lives in the real codebase — so nothing you study is abstract for long.

## Where to start

1. **Everyone:** read the [Roadmap](/roadmap) first. It has the weekly schedule, the milestone gates, and — importantly — **skip-ahead checkpoints** so experienced developers don't sit through material they already know.
2. **New to Go (any background):** start at [Module 0](/category/module-0-setup-orientation) and go in order.
3. **Backend developer (Java/Python/Node) new to Go:** do Module 0, take Module 1 at express pace (the *check yourself* sections tell you what to skim), then proceed normally from Module 2.
4. **Experienced Go developer:** do Module 0, take the Module 1–3 checkpoint quizzes, and jump to [Module 4 — The DX Platform](/category/module-4-the-dx-platform).

## Ground rules

:::info[The platform has opinions — learn them as defaults]
DX made deliberate, uniform technology choices. Across all services: **chi** is the router, **pgx v5** is the database layer (no ORM), **zap** is the logger, **viper** (via `dx-common-go/config`) loads configuration, and **amqp091-go** talks to RabbitMQ. This curriculum teaches those tools as the default, and flags where the platform deviates from generic Go tutorials you may find elsewhere.
:::

:::tip[Do the exercises]
Reading Go is easy; Go only sticks when you write it. Every module's exercises are sized to fit the weekly schedule. The capstone assumes you did them.
:::

Ready? Head to the [Roadmap](/roadmap).
