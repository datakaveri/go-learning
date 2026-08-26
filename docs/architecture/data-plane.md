---
title: Data Plane architecture
description: Ingestion, files, NGSI-LD, OGC, search, spatial/temporal persistence, authorization, and recovery.
---

# Data Plane architecture

**Status: Partially implemented.** File and OGC capabilities are substantial but incomplete; NGSI-LD/temporal coverage and public authorization are incomplete; carried decisions and typed obligation enforcement are planned.

The Data Plane owns bytes, spatial/temporal/query projections, ingestion state, and safe query execution. The Control Plane owns identities, catalogue metadata, grants, and entitlements. Their integration is through authenticated APIs and events, not shared policy or business databases.

## Components

| Component | Owns | Store/infrastructure | Status |
|---|---|---|---|
| `dx-files-connect-api-go` | File metadata, upload/download capabilities, processing jobs, object namespace | PostgreSQL, Redis, S3-compatible storage, RabbitMQ | Partially implemented |
| `dx-dataplane-rs-go` | NGSI-LD entity/temporal query execution and index interpretation | Elasticsearch | Partially implemented; not publicly routed |
| `dx-dataplane-ogc-go` | OGC Features/Tiles/Coverages/Processes/Records operations | PostgreSQL/PostGIS; job state | Partially implemented |
| `dx-subscription-go` | Subscription definition and continued delivery state | PostgreSQL/events | Partially implemented integration |

## File ingestion workflow

1. Provider authenticates through gateway; file service verifies workload and subject.
2. Application authorizes create/upload on the target catalogue/resource/organisation.
3. Service creates owned metadata and a bounded upload capability or streams the body under limits.
4. Object is written to the owned bucket/key namespace with checksum/content metadata.
5. A processing job is claimed idempotently and moves through validation/quarantine/derived-output states.
6. Metadata state and processing events are durable; consumers notify/audit/index as contracts require.
7. A file is discoverable/downloadable only when its state and authorization permit it.

Failures: partial upload is expired/cleaned; checksum mismatch or malware goes to quarantine; worker crash resumes from durable state; duplicate completion is idempotent; object/metadata inconsistency is reconciled.

## Discovery and retrieval

Catalogue owns searchable dataset metadata. Data Plane services resolve a governed resource to their owned query/index/object representation. Query parsers normalize and validate inputs before authorization binding. Responses use NGSI-LD/OGC-native representations through `HandleRaw` where the standard owns the shape.

## NGSI-LD and search

`dx-dataplane-rs-go` queries Elasticsearch-backed entity/temporal indexes. The service owns index mapping, alias, query builder, pagination/limits, and rebuild/reconciliation. Current index naming/configuration remains a blocker in the broader capability. It must not accept raw policy DSL or expose a public route before authorization can filter safely.

## OGC and spatial/temporal behavior

`dx-dataplane-ogc-go` uses PostGIS for OGC collection/feature/tile/coverage/process/record paths. It owns CRS/filter parsing, allowlisted functions, parameterized SQL/PostGIS construction, spatial indexes, output negotiation, result/job bounds, and standards conformance. EDR, Maps, and DGGS remain outside current delivered coverage and must be status-labelled if documented.

## Authorization

Current services still contain explicit service-level checks in paths where the target carried decision is unavailable. The target is [carried authorization](../integrations/data-plane-authorization.mdx): a PEP obtains a composite decision; Data Plane verifies/binds it and translates typed obligations. It never calls the PDP synchronously per query or interprets executable policy text.

Until the contract is implemented:

- keep unsafe public routes disabled;
- preserve explicit, reviewed in-service authorization for reachable paths;
- deny unsupported constraints/filters;
- keep query/body/result/time bounds;
- audit subject/actor/org/resource/query class/result without protected payloads.

## Retry and recovery

- Ingestion commands use idempotency keys and durable job state.
- Database/object/index writes define compensation/reconciliation rather than pretending to share a transaction.
- Search/index/object projections have watermarks, lag metrics, replay/rebuild tools, and ownership.
- Query reads do not retry indefinitely; respect request deadline and avoid amplification.
- Background work uses bounded concurrency/leases and stops on cancellation/loss.
- Subscription delivery stores cursor/effect idempotency and stops on authorization expiry/revocation.

## Signals

Observe upload bytes/time/failure, processing queue/age/state/retry, object integrity, query class/duration/result buckets, datastore pool/search health, spatial/temporal bounds, decision/obligation enforcement, rejected unsupported filters, subscription cursor/lag, reconciliation drift, and protected-access audit correlation.

