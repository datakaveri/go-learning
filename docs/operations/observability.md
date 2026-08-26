---
title: Observability and operations
description: Minimum logs, metrics, traces, health, readiness, alerts, redaction, and shutdown signals.
---

# Observability and operations

**Status: Platform foundations implemented; fleet signal coverage is partial.** Every new service must expose the minimum set before route enablement.

## Correlation model

Preserve W3C trace context plus the platform request ID across gateway, HTTP/gRPC, database, events, and workers. Security-sensitive signals should reference:

- `request_id`, `trace_id`, and event/correlation ID;
- calling workload;
- subject and actor kind/ID;
- organisation;
- resource type/reference and action;
- authorization decision ID/revisions when the composite contract exists;
- result, stable reason, duration, and retry/reconciliation state.

Use policy-approved pseudonymization where identifiers are sensitive. Never put these IDs into unbounded metric labels.

## Structured logs

Use the platform `zap` logger, consistent field names, one log per handled failure boundary, and safe error classification. Log start/end for long jobs and security/governance changes. Do not log routine success at a volume that hides failures.

Never log:

- passwords, access/refresh/workload tokens, client secrets, or private keys;
- approval codes or one-time exchange material;
- presigned URLs or storage credentials;
- payment/provider secrets;
- protected payloads, query results, or unredacted sensitive attributes;
- raw policy inputs when they contain contextual/private attributes.

## Metrics

Every service exposes:

- HTTP/gRPC request count, duration, in-flight, status/code;
- authentication and authorization allow/deny/error by bounded reason/profile;
- dependency latency/error/breaker state and connection-pool health;
- business operation count/result without user/resource labels;
- outbox age/count, queue/DLQ/retry, consumer connection, projection lag/drift;
- cache hit/miss/error/invalidation;
- worker concurrency/lease/retry/interruption and schedule delay;
- shutdown drain duration/incomplete count.

Data Plane adds query class, bounded result-size buckets, filter/obligation result, datastore duration, and quota/metering signals without revealing queries or records.

## Traces

Span gateway routing, identity validation, authorization, application use case, datastore query group, external call, event publish/consume, and agent tool stages. Record safe identifiers/revisions and error class. Do not attach body payloads or credentials. Sampling must retain denied/high-risk/error traces under policy without becoming an exfiltration path.

## Health and readiness

- `/healthz/live`: process can make progress. Avoid remote dependency restart loops.
- `/healthz/ready`: advertised contract is safe now—schema ready, required owner dependencies reachable, required policy/model/bundle ready.
- `/metrics`: internal scrape path, not generally public through the gateway.

OPA readiness is **Planned**. When implemented, readiness must report the active verified bundle revision and refuse activation of an invalid/unverified bundle; it must not imply OPA exists today.

## Alerts and runbooks

Alert on user impact or loss of safety: error/latency SLO burn, authorization backend error, workload-token failures, outbox age, queue/DLQ, projection lag/drift, migration failure, pool exhaustion, cache outage when fail-closed, worker lease churn, agent kill/revocation failure, audit ingestion gap, backup/restore failure, and incomplete shutdown.

Each alert links to a runbook with diagnosis, safe mitigation, replay/reconcile steps, data/security impact, escalation owner, and recovery verification.

## Graceful shutdown

On SIGTERM: remove readiness/stop intake, drain HTTP/gRPC, cancel workers, stop lease renewal and acknowledge only completed messages, flush bounded telemetry/audit, close dependencies, exit before grace expiry. Test this with a request, event, transaction, SSE stream, and worker in flight.

