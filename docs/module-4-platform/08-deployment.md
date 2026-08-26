---
title: Platform Deployment
sidebar_label: Deployment
description: Shared Helm chart, Argo CD applications, environment values, promotion, and operations.
---

# Platform deployment

**Status: Superseded as the canonical procedure by [Local development and GitOps deployment](../operations/local-deployment.md).** This lesson remains a compact GitOps reading exercise.

## GitOps topology

~~~mermaid
flowchart LR
    Source[Service source] --> Image[Immutable image digest]
    Values[Environment values] --> Chart[Shared Helm chart]
    Image --> Values
    Chart --> Manifests[Kubernetes manifests]
    AppSet[Argo CD ApplicationSet] --> App[Service application]
    Manifests --> App
    App --> Cluster
~~~

dx-gitops owns the running desired state. The shared chart can render Deployment, Service, probes, resources, autoscaling, disruption budget, network policy, ingress, external secret, and optional migration job.

## Change path

1. Build, test, scan, and sign an image.
2. Render chart plus service/environment values.
3. Review compatibility, secrets, network, probes, and resources.
4. Update an immutable image digest.
5. Sync to development and run smoke/contract checks.
6. Promote through environments with review.
7. Observe SLO, readiness, error, dependency, and backlog signals.

Development and staging may auto-sync. Production promotion is manual. A complete automated image-build-to-production-promotion workflow remains planned.

## Schema sequencing

Use expand-and-contract. One PreSync migration Job runs the same image and must finish before new serving replicas require the schema. Application replicas do not race to apply DDL.

## Rollback

Rollback selects a previous image/config digest. It does not reverse data mutation or destructive schema change. Every deployment plan states compatibility window, abort threshold, and forward-recovery path.

## Exercise

Trace one service from Dockerfile to image digest, environment values, rendered Deployment, ExternalSecret, probes, network policy, and Argo CD application. Identify what evidence is automatic and what remains manual.

## Check yourself

- Which repository is deployment source of truth?
- Why should production use an image digest?
- When is a migration Job safe?
- What cannot an application rollback undo?
