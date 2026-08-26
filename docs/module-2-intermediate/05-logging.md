---
title: Structured Logging
sidebar_label: Logging
description: Events, levels, context, request IDs, traces, redaction, and one-time handling.
---

# Structured logging

## Log events, not prose

~~~go
log.Info(
    "widget created",
    zap.String("widget_id", widget.ID),
    zap.String("org_id", widget.OrgID),
    zap.String("request_id", httpx.RequestIDFrom(ctx)),
)
~~~

Stable message and field names make logs queryable. Include what an operator needs to correlate the event; avoid entire request bodies.

## Levels

- Debug: temporary diagnostic detail safe to disable.
- Info: normal lifecycle or significant domain event.
- Warn: degraded but handled condition.
- Error: operation failed and the current boundary handles or terminates it.

Do not log an error and return it unless the log adds unique evidence that will not appear at the handling boundary.

## Request path

platform/http's standard stack establishes tracing, request ID, real IP, structured request logs, CORS, auth gates, timeout, and recovery. Service handlers obtain the logger through injection and add domain fields.

Request IDs correlate within an exchange. Trace IDs connect spans across HTTP and event boundaries. Propagate correlation IDs into event envelopes.

## Redaction

Never log:

- bearer or delegated tokens;
- passwords, cookies, client secrets, workload tokens, or private keys;
- database DSNs with credentials;
- payment secrets;
- unredacted model prompts or tool results;
- sensitive user or data payloads.

Log a safe identifier, classification, byte count, or hash only when the use is documented.

## Cardinality

High-cardinality fields are acceptable in logs, but not as metric labels. Even in logs, avoid uncontrolled client strings that create cost or injection risk.

## Exercise

Add logs to a create workflow at its top handling boundary. Include operation, actor, target, outcome, request ID, and error classification. Write a test using a zap observer and assert no secret value appears.

## Check yourself

- Why does the standard HTTP stack run before service middleware?
- Where should one error be logged?
- What belongs in an audit event instead of a diagnostic log?
- Why must prompt and token data be redacted?
