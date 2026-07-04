---
title: Containers & Kubernetes
sidebar_label: Containers & K8s
description: Multi-stage Docker builds for Go, the Compose dev stack, and the Kubernetes concepts a service developer needs.
---

# Containers & Kubernetes

## Learning objectives

- Write a multi-stage Dockerfile producing a small, static Go image.
- Understand the platform's build-context convention (workspace root + whitelist).
- Navigate the local Docker Compose stack confidently.
- Know the Kubernetes objects that touch your service: Deployment, probes, resources, ConfigMap/Secret, CronJob.

## Prerequisites

- [Project Structure](project-structure), [Observability](observability) (probes)

## Time estimate

**4 hours**

## Concepts

### Go was made for containers

A Go binary is statically linked and self-contained — the runtime image needs no interpreter, no framework, barely an OS. The standard two-stage build:

```dockerfile
# --- build stage: full toolchain
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download            # cached layer: re-runs only when go.mod changes
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# --- runtime stage: almost nothing
FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/server /server
COPY --from=build /src/configs /configs
USER nonroot
ENTRYPOINT ["/server"]
```

Points that matter: `CGO_ENABLED=0` forces a fully static binary; the mod-download layer makes rebuilds fast; distroless (or alpine) keeps the image tens of MB with near-zero CVE surface; `USER nonroot` because containers shouldn't run as root. Remember the embedded spec/SQL from [Project Structure](project-structure)? This is where it pays: the binary *is* the deployment artifact.

### The platform's build context convention

DX services build with **the workspace root as context** (`context: .` in Compose, from `cdpg-claude/`), not the service directory. Why: the `replace github.com/datakaveri/dx-common-go => ../dx-common-go` directive from [Packages & Modules](../module-1-go-fundamentals/packages-modules) must resolve *inside the build*, so the build context must contain both the service and `dx-common-go`. A `.dockerignore` whitelist keeps the context lean. Practical consequence: build commands run from the workspace root, and a service's Dockerfile copies paths relative to it.

### The Compose dev stack

`make dev-up` is Docker Compose running the entire platform: infra (Postgres, Redis, RabbitMQ, Elasticsearch, MinIO, Keycloak) plus both service tracks. Daily muscle memory:

```bash
docker compose ps                        # what's running, what's healthy
docker compose logs -f dx-acl-go        # follow one service's logs
docker compose up -d --build dx-acl-go  # rebuild + restart just your service
docker compose exec postgres psql -U postgres iudx_db   # poke the database
```

The edit-verify loop while developing a service: change code → `up -d --build <service>` → `logs -f` → curl through the gateway.

### Kubernetes: the developer's subset

Production runs on Kubernetes. You don't operate the cluster, but your service's manifest choices are your responsibility. The objects that matter:

- **Deployment** — desired replica count of your pod; rolling updates replace pods gradually (SIGTERM → your [graceful shutdown](http-servers-middleware) drains → new pod).
- **Probes** — your [health endpoints](observability), wired: liveness probe → `/healthz/live` (fail ⇒ restart pod), readiness probe → `/healthz/ready` (fail ⇒ remove from load balancing). You built the distinction; K8s consumes it.
- **Resources** — requests (scheduling reservation) and limits (hard cap). Go note: set `GOMAXPROCS`/`GOMEMLIMIT` in accordance with limits so the runtime knows its actual budget.
- **ConfigMap / Secret** — become the environment variables your [viper config](../module-2-intermediate/configuration) reads. The whole env-override design exists for this moment.
- **Service** — stable virtual IP + DNS name in front of your pods.
- **CronJob** — the scheduled one-shots from [Workers & Cron](workers-cron), with `concurrencyPolicy` and `backoffLimit`.

```yaml
livenessProbe:
  httpGet: {path: /healthz/live, port: 8080}
readinessProbe:
  httpGet: {path: /healthz/ready, port: 8080}
resources:
  requests: {cpu: 100m, memory: 128Mi}
  limits: {memory: 256Mi}
```

Deployment to the clusters is GitOps via ArgoCD — you change manifests in a Git repo (`dx-gitops`), the cluster converges to them. That workflow is Module 4's [Deployment](../module-4-platform/deployment) page; here you just need the object vocabulary.

:::info[Platform connection]
Every service repo carries the multi-stage Dockerfile; the Compose files in `cdpg-claude/` are the readable inventory of the whole platform (read `docker-compose.go-stack.yml` — ports, env vars, dependencies, healthchecks all in one place). The probe endpoints are the GO-SERVICE-STANDARDS health contract — a service missing `/healthz/ready` can't be deployed sanely, which is *why* the standard exists.
:::

## Exercises

1. Containerize `dx-scratch-go` with the two-stage Dockerfile. Compare image sizes: single-stage golang base vs distroless. Run it with config via `-e` env vars.
2. Wire your scratch service into a small Compose file with Postgres + healthcheck + `depends_on: condition: service_healthy`. Bring it up, kill Postgres, watch your readiness (and Compose's view of it).
3. Write the full Deployment YAML for your service — probes, resources, env from a ConfigMap — and validate it (`kubectl apply --dry-run=client -f` if you have kubectl; otherwise careful reading against the docs).
4. In the real workspace: rebuild one Go service with `docker compose up -d --build dx-registry-go`, follow its boot logs, and identify each step of the boot contract (config → logger → deps → server) as it logs.

## Check yourself

- Why two stages, and what belongs in each?
- Why does the platform build from workspace root instead of the service directory?
- What happens to traffic during a rolling update, step by step, and which two things you built make it seamless (name them)?
- Which probe restarts pods, and why must it not check Postgres?

## References

- [Docker multi-stage builds](https://docs.docker.com/build/building/multi-stage/) · [distroless](https://github.com/GoogleContainerTools/distroless)
- [Kubernetes concepts — Workloads](https://kubernetes.io/docs/concepts/workloads/)
- Platform: any service Dockerfile; `docker-compose.go-stack.yml`; `claude-docs/DOCKER.md`
