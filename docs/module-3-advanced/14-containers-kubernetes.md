---
title: Containers and Kubernetes
sidebar_label: Containers & Kubernetes
description: Reproducible images, pod lifecycle, probes, resources, secrets, and shared Helm deployment.
---

# Containers and Kubernetes

## Image principles

Use a multi-stage build, pinned toolchain and base images, a non-root runtime user, no compiler or source in the final image, read-only filesystem where practical, and an immutable image digest.

~~~dockerfile
FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/service /service
USER nonroot:nonroot
ENTRYPOINT ["/service"]
~~~

The orchestration repository uses a workspace-level build context so the service and dx-common-go checkout are available together.

## Pod contract

- /healthz/live for process liveness;
- /healthz/ready for mandatory dependencies;
- SIGTERM triggers HTTP drain, worker stop, and dependency close;
- requests and limits are set;
- secrets arrive through secret references;
- service account and network policy are least privilege.

## Shared chart

The GitOps repository's shared Helm chart renders Deployment, Service, probes, resources, HPA, PDB, NetworkPolicy, Ingress, ExternalSecret, and an optional schema-migration Job.

Values choose image digest, config, secret reference, replicas, resource policy, ingress, and optional features. Do not fork the chart to set a service value the existing schema supports.

## Exercise

Build your service image, run it as non-root with read-only filesystem, send SIGTERM during a slow request, and confirm the request drains. Render the Helm chart and inspect probes, secret references, and network policy.

## Check yourself

- Why use an immutable digest?
- What happens when readiness fails?
- Why does the build context include more than one repository?
- What must complete before the pod exits?
