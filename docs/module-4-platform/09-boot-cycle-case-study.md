---
title: Case Study — Restructuring a Boot Cycle
sidebar_label: Case Study — Boot Cycle
description: How dx-user-go's 363-line main.go became a readable, documented boot cycle without changing behaviour, and how the change was proven.
---

# Case study: restructuring a boot cycle

**Status: In development.** The refactor is on dx-user-go branch [`refactor/main-boot-cycle`](https://github.com/datakaveri/dx-user-go/tree/refactor/main-boot-cycle) (commit `8061f89`). It is tested but not yet merged to `dev`. The bootstrap behaviour it relies on is **Implemented** in `dx-common-go/platform/bootstrap`.

Read [Service anatomy](./04-service-anatomy.md) and [Bootstrap and typed configuration](../platform/bootstrap-config.md) first.

## Learning goals

After this lesson you can:

- separate what the shared bootstrap owns from what a service's own composition root owns;
- restructure a composition root without changing behaviour, and say which properties had to stay fixed;
- spot the typed-nil interface trap in wiring code;
- choose a verification ladder that proves the change, not just that it compiles.

## 1. Two layers

| Layer | Owns | Where |
|---|---|---|
| Shared bootstrap | Boot mode, signals, config, logger, tracing, migrations, Postgres, start order of servers and workers, shutdown order, exit code | `dx-common-go/platform/bootstrap/bootstrap.go` (`run`, `serve`) |
| Service composition root | Every other dependency, the domain services, and the HTTP/gRPC interfaces | `dx-user-go/cmd/server` |

The refactor changed **only** the second layer. The shared bootstrap was untouched, so the start-up and shutdown order stayed the same.

## 2. Before

`dx-user-go/cmd/server/main.go` held everything in 363 lines:

- the `bootstrap.Spec`;
- `wire()`: about 190 lines building roughly 15 components (token issuer, RabbitMQ publishers, credits client, catalogue client, Keycloak clients, seven domain services, mailers, audit publisher, workload verifier, gRPC lookups, router);
- a `grpcDeps` struct that `wire()` filled **as a side effect** for `Spec.GRPC` to read later;
- helpers.

Comments explained decisions mostly through ticket codes (`F39`, `SB-1`, `P2-1`). A reader could not answer three basic questions:

- What starts, in what order?
- Which dependencies are optional?
- What closes at shutdown, in what order?

## 3. After

The design follows the UNIX habit of small parts, each doing one job, joined in one visible place:

| File | One job |
|---|---|
| `main.go` | Only the `Spec`. `Wire` calls `build()`, then `newHTTPHandler()`; `GRPC` calls `newGRPCServer()`. |
| `components.go` | `build()`: one numbered list of everything the service is made of, returning a `components` value |
| `broker.go`, `identity.go`, `credits.go`, `catalogue.go`, `mail.go`, `audit.go` | One outside system each: open it, register its closer and probe, return it |
| `http.go`, `grpc.go` | Turn `components` into an interface. They only read it. |
| `doc.go`, `docs/BOOT_CYCLE.md` | The boot and shutdown cycle in plain language, a dependency table, boot modes and exit codes |

Open `components.go` first. Its `build()` comment is the whole boot in ten lines.

## 4. What had to stay fixed

A restructure is only safe if you know which properties are behaviour. Here there were five:

1. **Shutdown-hook order.** Bootstrap runs closers in reverse registration order. The order rabbitmq → credits → catalogue → auditing therefore had to be kept exactly. Moving a step in `build()` silently reorders shutdown.
2. **Every log and error string,** byte-identical. Operators and alerts search for them. This was checked by diffing all string literals before and after.
3. **Required versus optional.** Each dependency must still be either fatal or a warning, exactly as before.
4. **Routes, middleware order, error mappers, auth config.**
5. **"Off" means a true nil.** See the next section.

Two internal clean-ups were allowed because nothing outside the process can observe them:

- one stateless Keycloak client instead of five;
- one workload signer instead of two, so credits and catalogue now share a token cache.

## 5. The typed-nil trap

Several dependencies are passed around as interfaces: `service.CreditProvisioner`, `delegation.Mailer`, `credential.Revocations`, `bootstrap.Servable`. In Go, an interface that holds a nil pointer is **not** equal to `nil`:

~~~go
func openCreditsWrong(...) service.CreditProvisioner {
    var cc *credits.Client // nil
    return cc              // non-nil interface holding a nil pointer
}

if credits != nil { // true!
    credits.Provision(ctx, ...) // nil pointer dereference at run time
}
~~~

The original code avoided this by declaring interface-typed variables and assigning only on success. The refactor keeps that guarantee by returning the interface type itself, with a bare `nil` when the dependency is off. `newGRPCServer` returns `bootstrap.Servable`, not `*grpcserver.Server`, for the same reason. Bootstrap treats a nil `Servable` as "no gRPC server".

## 6. How the change was proven

Compiling proves little in wiring code, because most failures only show at start-up or on the first real request. The change was accepted after each rung of the ladder in [Verify a boot cycle](../operations/verify-boot-cycle.md):

| Rung | Result |
|---|---|
| Unit | build, vet, all tests, plus 4 new helper tests |
| Boot modes | `config-check` exits 0; an unknown mode exits 1; with no database the old and new binaries fail identically |
| Local stack | healthy readiness, clean SIGTERM shutdown, exit 0 |
| Real environment | 47 recorded calls across all 7 controllers against the hosted dev environment, including create, update and delete on test records the run created itself |
| Everything on | gRPC, strict workload identity, credits, catalogue and e-mail all built; clean shutdown with every closer registered |
| Reject path | direct HTTP and gRPC calls without a gateway credential refused |
| Accept path | the same calls through a local dx-gateway-go answered 200 with the right user |
| Outbound path | a request that makes dx-user-go call dx-catalogue-go succeeded; with the catalogue stopped, the same request failed |

Two gaps remain, both outside the refactor:

- the credits call is untested, because its only trigger changes a real user's role;
- the catalogue's owner-lookup call is untested, because no available role allows its scope.

## 7. Findings the restructure exposed

Making the boot readable surfaced real issues that the old file hid. They are recorded as gaps, not fixed in the refactor:

- **Optional broker controls readiness.** RabbitMQ is optional by design, but its check is registered with `App.Probe`, which is critical. A broker outage would make every replica unready. The direction is `Health.AddOptional`; this is deferred.
- **Shutdown can exceed the grace period.** Each bootstrap shutdown phase has its own 20 s budget, about 60 s in total against Kubernetes' 30 s default.
- **Lazy gRPC clients hide wrong addresses until the first request.** The failure then reached the client as a 500 that included an internal address.

## Exercises

1. **Read the order.** From `build()`'s comment alone, write the shutdown order of dx-user-go's closers, including the two bootstrap registers itself (`postgres`, then `observability`). Check your answer against `docs/BOOT_CYCLE.md`.
2. **Find the trap.** In `credits.go`, change `openCredits` to return `*credits.Client`, and keep the call site in `build()` unchanged. Predict what happens when `credits.enabled: false`, then write a test that catches it.
3. **Move a step.** Swap steps 4 and 5 in `build()`. Which observable behaviour changed? How would you have caught it in review?
4. **Apply it.** Pick another service whose `main.go` is over 200 lines. List its dependencies with required/optional and closer order *before* touching code. That table is the behaviour you must keep.

## Check yourself

- Why may `Spec.GRPC` read what `Spec.Wire` built, and what would break if bootstrap called them in the other order?
- Why is "every log string byte-identical" part of behaviour?
- Which rung of the ladder proved the workload signer is *used*, not just built?
- Why is a stopped-dependency negative control stronger than a successful call alone?
