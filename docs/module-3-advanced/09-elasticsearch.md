---
title: Elasticsearch
sidebar_label: Elasticsearch
description: Search adapters, mappings, index lifecycle, projections, queries, and verification.
---

# Elasticsearch

Elasticsearch powers catalogue search and NGSI-LD entity or temporal queries. Treat index design as domain architecture, not an incidental client call.

## Package surface

dx-common-go/database/elasticsearch contains client, mapping, indexing, query, and repository packages. Service application code defines a search port; the adapter imports these packages.

## Authoritative or projected

Decide whether an index is:

- authoritative data requiring snapshots and restore;
- a projection rebuildable from service-owned state or an event replay source.

Document the decision. Do not call a projection rebuildable unless the rebuild path and time are tested.

## Mapping lifecycle

- Use explicit mappings for dates, keywords, geo fields, nested objects, and numeric units.
- Version index names or templates.
- Build a new index, backfill, verify counts and sample queries, then atomically switch an alias.
- Retain the previous index until rollback risk passes.
- Normalize index and alias naming across producer and reader configuration.

## Query safety

Bound result size, aggregation buckets, wildcard use, script access, and deep pagination. Use search_after with a stable tie-breaker for deep result traversal. Keep client-visible sort and filter fields on an allowlist.

## Consistency

Search refresh is not a transaction. A write can be committed before the document is searchable. Client workflows either tolerate that lag, read authoritative state for immediate confirmation, or wait through an explicit job.

## Exercise

Design an index for catalogue items with text, keyword, organization, tags, timestamp, and location. Define mapping, alias, query allowlist, zero-downtime change, and rebuild source.

## Check yourself

- Is every platform index a projection?
- Why is alias switching useful?
- When should you use search_after?
- What evidence catches an index-name mismatch?
