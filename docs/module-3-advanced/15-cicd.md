---
title: CI/CD and Supply Chain
sidebar_label: CI/CD
description: Reproducible checks, artifacts, security evidence, GitOps promotion, and rollback.
---

# CI/CD and supply chain

## Pull-request gate

A service pipeline should run:

~~~bash
gofmt -w .
git diff --exit-code
go vet ./...
go test ./...
go test -race ./...
~~~

Add lint, generated-code drift, OpenAPI/route drift, event contracts, integration tests, fuzz smoke, dependency and secret scan, image scan, and manifest rendering according to risk.

## Build evidence

Record:

- source commit and clean tree;
- Go and module versions;
- test commands and results;
- SBOM and vulnerability result;
- image digest and signature;
- configuration/schema compatibility;
- rendered manifests.

Rebuilding the same source should produce functionally identical artifacts.

## GitOps promotion

Deployment changes update the immutable digest and environment values in dx-gitops. Argo CD reconciles the reviewed desired state. Development and staging may auto-sync; production is manually promoted.

The repositories currently do not provide a complete automated image-build-to-production-promotion chain. Treat that integration as engineering work, not an existing guarantee.

## Rollout and rollback

Use readiness and bounded surge/unavailable settings. Observe SLO and dependency signals during rollout. Roll back the image and config to a known digest when safe. Schema changes require forward-compatible recovery; an image rollback cannot undo destructive DDL.

## Exercise

Design a pipeline for a service with PostgreSQL migrations and RabbitMQ events. Mark fast versus integration gates, required evidence, promotion approvals, post-deploy checks, and rollback constraints.

## Check yourself

- Why does gofmt need a clean-tree check?
- What makes an artifact traceable?
- Which repository changes the running desired state?
- Why is database rollback different from image rollback?
