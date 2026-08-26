---
title: Observability
sidebar_label: Observability
description: Logs, metrics, traces, health, audit, SLOs, and asynchronous work.
---

# Observability

## Four distinct signals

- Logs explain individual events with structured context.
- Metrics show aggregate rates, latency, errors, and saturation.
- Traces connect work across HTTP and broker boundaries.
- Audit records answer who performed a material action on what and with what outcome.

Do not use one signal as a substitute for all others.

## HTTP and dependency signals

platform/http establishes request IDs, traces, structured request logs, and recovery. platform/observability/health exposes liveness and readiness. Adapters measure dependency operation and result.

Track:

- requests by route template, method, status, and latency;
- authorization allow/deny and verification failures;
- database pool, query, lock, and migration signals;
- event confirmation, backlog, retry, drop, and dead-letter count;
- cache hit/miss/error;
- search and object-store duration and failure;
- jobs by queued, active, retrying, failed, and completed state.

Never put raw URL IDs, users, errors, or unbounded client values into metric labels.

## Health

Liveness does not call remote dependencies. Readiness fails for mandatory dependencies and reports optional degradation. A green process with a broken required database should not receive traffic.

## SLO thinking

Define a user-visible indicator, target, and window. An availability SLI may be the proportion of valid requests served successfully under a latency threshold, excluding explicit client errors. Use error budget to guide release risk.

## Exercise

Create an observability plan for dx-files-connect-api-go multipart completion: logs, metrics, spans, audit, readiness, dashboard, and three alerts. Include object-store and processing-job failure.

## Check yourself

- Why is resource ID a bad metric label?
- What belongs in audit but not logs?
- Why must liveness ignore PostgreSQL?
- Which signal shows an event consumer backlog trend?
