---
title: Verify a boot cycle
description: A step-by-step ladder for proving a change to a service's start-up, wiring, or shutdown — from unit tests to a real environment, through the gateway.
---

# Verify a boot cycle

**Status: Implemented procedure.** First used on 2026-11-26 to accept the dx-user-go composition-root refactor; see the [case study](../module-4-platform/09-boot-cycle-case-study.md). The procedure applies to any service built on `dx-common-go/platform/bootstrap`.

Wiring code fails in ways unit tests can't see: a client that is built but never used, an "off" value that is secretly "on", a gate that rejects everything or accepts everything. Climb the ladder in order. Each rung proves something the previous one could not.

## Safety rules for shared environments

Rungs 4–8 run your local build against a shared environment, such as the hosted dev environment. Before starting:

- **Create, never modify.** Only update or delete records your run created. Name them recognisably, for example `test-<you>-<yyyymmdd>-*`, keep a list of their IDs, and delete them at the end, even if a step fails.
- **No user-changing calls on accounts you don't own.** Role grants, password changes, delegations to other people and request decisions all change real users.
- **Turn off outbound side effects you don't need,** such as e-mail (`EMAIL_ENABLED=false`), so test requests don't mail real people.
- **Check message consumers before connecting to a shared broker.** A consumer on a *shared, named* queue competes with the deployed service and takes its messages. Private, exclusive, auto-delete queues are safe. For example, `dx-common-go/auth/appid/revocation.go` declares one queue per process.
- **Change settings only in your own process.** Pass them on the command line; never edit the shared env file.
- **Keep secrets out of output.** Don't print tokens, passwords or env dumps. Mask secrets in any recorded request or response.

## Load env files literally

`set -a; source dev.env` expands `$` inside double-quoted values, so a password like `"ab$cd"` silently arrives shortened. Load the file literally instead:

~~~bash
load_env() {  # usage: load_env file.env — exports KEY=VALUE pairs without expansion
  while IFS= read -r line || [ -n "$line" ]; do
    [[ -z $line || $line == \#* || $line != *=* ]] && continue
    k=${line%%=*}; v=${line#*=}
    [[ $v == \"*\" || $v == \'*\' ]] && v=${v:1:-1}
    export "$k=$v"
  done < "$1"
}
~~~

Run the service in a subshell, so the variables don't leak into your shell:

~~~bash
( load_env ~/Downloads/dev.env; export SERVER_PORT=8191; exec ./bin/server )
~~~

## The ladder

| # | Rung | Proves |
|---|---|---|
| 1 | Unit | It compiles, and moved helpers still behave |
| 2 | Boot modes | Config and failure paths are unchanged |
| 3 | Local stack | A normal start and stop work |
| 4 | Real environment, every controller | The wiring holds against real dependencies |
| 5 | Every optional part on | Each constructor and closer runs |
| 6 | Reject path | The identity gates are wired into every interface |
| 7 | Accept path through the gateway | Trusted identity flows end to end |
| 8 | Outbound call with a negative control | The downstream clients and signer are *used*, not just built |

### 1. Unit

~~~bash
go build ./... && go vet ./... && go test ./...
gofmt -l cmd/
~~~

Add tests for anything you moved that has logic of its own. dx-user-go added four:

- path-scoped middleware;
- the fail-closed catalogue fallback;
- gRPC is off at port 0;
- a disabled signer returns `nil`.

### 2. Boot modes

~~~bash
DX_BOOT_MODE=config-check go run ./cmd/server   # expect exit 0, nothing contacted
DX_BOOT_MODE=bogus        go run ./cmd/server   # expect exit 1, "is not a boot mode"
~~~

Then start the **old** and the **new** binary with no database reachable. Both must print the same error and exit 1.

### 3. Local stack

Start only the dependencies the service needs at boot (from `cdpg-claude`):

~~~bash
docker compose -f docker-compose.yml -f docker-compose.go-stack.yml \
  --profile infra --profile apps --profile go-apps up -d --build <service>
~~~

Then check:

- `/healthz/live` and `/healthz/ready` return 200, and every check is `up`;
- `docker stop` (SIGTERM) logs `shutting down`, then `stopped`, with exit code 0.

### 4. Real environment, every controller

Run the local build with the environment's env file, then call at least one read endpoint per controller and one full create → read → update → delete cycle on records you own. Record every request and response so a reviewer can follow along. Mask secrets in the recording.

Two overrides are usually needed, and both stay inside your own process:

- **A free port,** if the default is taken.
- **The workload-identity setting** (`WORKLOAD_VERIFIER_ENFORCEMENT=off`), only for this rung: your direct calls can't carry a gateway credential. Rungs 6–7 test with it back on.

### 5. Every optional part on

Start with every optional dependency enabled:

- gRPC port set;
- `WORKLOAD_VERIFIER_ENFORCEMENT=required`;
- credits and catalogue on;
- e-mail on.

Check two things:

- **Start-up:** every port listens, and none of the "X is off" warnings appear. That proves each constructor ran.
- **Shutdown:** SIGTERM shows `grpc server stopped`, then `stopped`, exit 0, no `closer failed` lines.

Remote clients built with `grpc/client.Dial` connect lazily, so this rung proves they were *built*, not that they work. That is rung 8.

### 6. Reject path

With `WORKLOAD_VERIFIER_ENFORCEMENT=required`, call the service **directly** and expect a refusal every time:

| Request | Expect |
|---|---|
| HTTP, no credentials | 401 |
| HTTP, valid user token only | 401 (`workload credential required`) |
| HTTP, user token plus a bogus `X-DX-Workload` header | 401 |
| gRPC, no `x-dx-workload` metadata | `Unauthenticated` |
| gRPC, bogus `x-dx-workload` | `Unauthenticated` |

A 200 on any of these means the gate is not wired. For gRPC without `grpcurl`, a ten-line Go client using the generated stubs (for example `dx-common-go/grpc/userpb`) and `metadata.NewOutgoingContext` is enough. Call only read-only methods.

### 7. Accept path through the gateway

Run `dx-gateway-go` locally with the gateway's own workload credential. Give it a copy of its `configs/config.yaml` where the only change is your service's `upstream`, pointed at your local port. Then call through it:

~~~bash
TOKEN=...   # a user token for an account you're allowed to use
curl -s -o /dev/null -w '%{http_code}\n' "http://localhost:<gw-port>/iudx/v2/auth/user" -H "Authorization: Bearer $TOKEN"   # expect 200
curl -s -o /dev/null -w '%{http_code}\n' "http://localhost:<svc-port>/iudx/v2/auth/user" -H "Authorization: Bearer $TOKEN"  # expect 401
~~~

Confirm that the service's request log shows the right `user_id`. The gateway strips the user token (`GATEWAY_STRIP_AUTHORIZATION=true`), so that identity can only have come from subject headers trusted from a verified gateway workload.

The gateway reporting `not_ready` locally is expected when its permission service (dx-authz-go) is cluster-only. Routes that don't call authz still work.

### 8. Outbound call with a negative control

Pick a request that makes your service call a downstream service. For dx-user-go, creating an app credential with a `data-access` limit on a catalogue item makes it call dx-catalogue-go over gRPC, with its own workload token.

1. **Downstream running locally:** expect success. Delete what you created.
2. **Downstream stopped:** expect the same request to fail with a connection error, and nothing to be created.

Step 1 alone could pass for the wrong reason. Step 2 proves the call really happened. If the downstream also runs strict workload identity, step 1 also proves your signer's credential was accepted.

## Record the result

Write down, per rung, what ran, the result, and what could not be covered and why. "Credits untested: its only trigger changes a real user's role" is a result. A blank is not.
