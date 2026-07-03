---
title: Coding Standards & the Service Checklist
sidebar_label: Standards Checklist
description: GO-SERVICE-STANDARDS as a study guide, plus the style rules — Google primary, Uber supplementary, CDPG on top.
---

# Coding Standards & the Service Checklist

## Learning objectives

- Know the three-layer style authority: Google Go Style Guide → Uber guide → CDPG platform rules.
- Internalize the non-negotiables that fail review, with the *why* for each.
- Use GO-SERVICE-STANDARDS.md as the self-review checklist it's designed to be.

## Prerequisites

- [Service Anatomy](service-anatomy) — you've now seen every rule obeyed in real code.

## Time estimate

**3 hours**

## Concepts

### The authority stack

When style questions arise, precedence is explicit: **CDPG platform rules** (the tables below + GO-SERVICE-STANDARDS.md) override the **Uber Go Style Guide**, which supplements the **Google Go Style Guide** as the base. Mechanical matters (formatting, imports) are settled by `gofmt`/`goimports` and aren't discussed at all. This exists so review threads argue about *design*, never taste.

### The non-negotiables — with reasons

You've met every one of these in the curriculum; here they are as the single review-time list. Left column fails review; right column is why.

**Errors & control flow**

| Rule | Because |
|---|---|
| Wrap with `fmt.Errorf("context: %w", err)`; match with `errors.Is`/`As` | error chains are the debugging trail ([M1](../module-1-go-fundamentals/error-handling)) |
| Log **or** return, never both | one failure = one log line |
| Early returns; no `else` after return; happy path left-aligned | readability is a review criterion, not a preference |
| Custom error types implement `Unwrap()` | or `errors.Is/As` go blind through your type |

**Concurrency**

| Rule | Because |
|---|---|
| No fire-and-forget goroutines — errgroup/channel-close/ctx exit, always | leaks and silent failures ([M2](../module-2-intermediate/concurrency)) |
| Channel buffers 0 or 1; larger needs a justifying comment | big buffers hide backpressure |
| `ctx` first param, never stored in structs, typed keys for values | ([M2](../module-2-intermediate/context)) |
| Mutex as value field (`mu sync.Mutex`), never copy a struct holding one | copied locks guard nothing — `go vet` catches some, review the rest |

**Types & APIs**

| Rule | Because |
|---|---|
| Comma-ok for every type assertion | the one-value form panics ([M1](../module-1-go-fundamentals/interfaces)) |
| Interfaces at the consumer; accept interfaces, return concrete types | testability + honest APIs |
| Copy slices/maps at trust boundaries | aliasing is action-at-a-distance ([M1](../module-1-go-fundamentals/collections)) |
| Enums start at 1, or zero = explicit "unknown" | accidental zero values shouldn't mean something real |
| No `util`/`common`/`helpers` packages; `doc.go` on exported packages | names are design ([M1](../module-1-go-fundamentals/packages-modules)) |

**Platform integration**

| Rule | Because |
|---|---|
| Envelopes via `dxresp` writers; errors via `dxerrors` taxonomy — never hand-built JSON | one contract for every client ([M3](../module-3-advanced/rest-api-development)) |
| `StandardStack()` + OpenAPI validation + resolver, in that order | order is semantics ([M3](../module-3-advanced/http-servers-middleware)) |
| `dxconfig.LoadService[T]`, baked YAML defaults, no secrets anywhere in the repo | ([M2](../module-2-intermediate/configuration)) |
| Parameterized values, allowlisted identifiers, explicit soft-delete filters | injection + resurrection bugs ([M3](../module-3-advanced/database-patterns)) |
| `defer tx.Rollback(ctx)` on the line after `Begin` | every exit path covered ([M3](../module-3-advanced/transactions)) |
| `ReliablePublisher` + outbox for state-changing events; consumers idempotent, DLQ'd | ([M3](../module-3-advanced/event-driven-rabbitmq)) |
| zap only — no `fmt.Println`, no stdlib log | structured or it didn't happen ([M2](../module-2-intermediate/logging)) |

### The checklist as a workflow

GO-SERVICE-STANDARDS.md is organized by area (boot, config, HTTP, contract, persistence, RMQ, workers, observability, security, testing, repo shape). Use it three ways:

1. **Before a PR** — self-review your diff against the relevant sections. Cheaper than a review round-trip.
2. **Reviewing others** — cite sections, not opinions: "standards §persistence: identifiers from allowlists" ends discussions kindly.
3. **The capstone grade** — [your capstone](../capstone/capstone-service) is assessed section by section against it.

And always, the mechanical gate before any human looks:

```bash
gofmt -l .   # empty
goimports -w .
go build ./... && go vet ./... && go test ./...
golangci-lint run
```

:::info[Platform connection]
The full rule set lives in two places you should now read end to end: `claude-docs/GO-SERVICE-STANDARDS.md` (the service checklist) and the workspace's `go-style` skill (`skills/go-style/` — the same rules organized for code-time application, with a `references/cdpg-conventions.md` covering response format, middleware order, health endpoints, config, DB patterns). If you use Claude Code for development, install the skill — it applies these standards automatically as you write.
:::

## Exercises

1. Read GO-SERVICE-STANDARDS.md in full (finally, with everything to hang it on). Note the three requirements this curriculum gave the least attention — and read their sections twice.
2. Self-review `dx-scratch-go` against the full checklist. Produce a findings list with severity; fix the top five.
3. Review practice: take a diff you wrote in Module 3 and write a review of it *citing standards sections* for each finding — the skill you'll use on colleagues' PRs.
4. Calibration: find one place in a real service that technically deviates from a rule. Decide: defect worth a ticket, documented exception, or your misreading? (All three exist — knowing which is which is the mature skill.)

## Check yourself

- What's the authority order when style guides conflict?
- For any three non-negotiables: state the rule *and* the failure it prevents.
- What are the three uses of the standards doc as a workflow?

## References

- [Google Go Style Guide](https://google.github.io/styleguide/go/) · [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Effective Go](https://go.dev/doc/effective_go) · [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- Platform: `claude-docs/GO-SERVICE-STANDARDS.md`; the `go-style` skill in the workspace
