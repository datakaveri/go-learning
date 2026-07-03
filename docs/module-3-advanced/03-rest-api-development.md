---
title: REST API Development
sidebar_label: REST API Development
description: The DX API contract — /iudx/v2 paths, response envelopes, error URNs, pagination, and OpenAPI-first validation.
---

# REST API Development

## Learning objectives

- Design REST endpoints under the platform's `/iudx/v2/<domain>` path contract.
- Return the standard success and error envelopes; map errors to the URN taxonomy.
- Implement pagination, filtering, and sorting the platform way.
- Validate requests against an embedded OpenAPI 3 spec.

## Prerequisites

- [HTTP & Middleware](http-servers-middleware), [Error Handling](../module-1-go-fundamentals/error-handling), [Generics](../module-1-go-fundamentals/generics)

## Time estimate

**5 hours**

## Concepts

### The path contract

Every Go service exposes its API under a versioned, domain-scoped base path:

```
/iudx/v2/policies
/iudx/v2/files/upload
/iudx/v2/community/discussions/{id}
```

One legacy exception (`dx-acl-go` keeps `/iudx/acl/apd/v2` for client compatibility) proves the rule: **the contract belongs to the clients** — the Go rewrite preserves legacy paths and shapes wherever clients depend on them, and every deliberate delta is documented for consumers. Resource naming is standard REST: plural nouns, no verbs in paths, nesting only where a real ownership relation exists.

### One envelope for every response

Clients should parse one shape everywhere. Success:

```json
{
  "type": "urn:dx:acl:Success",
  "title": "Success",
  "results": [ { "policyId": "…", "status": "active" } ]
}
```

Errors carry a machine-readable **URN type** from the platform taxonomy plus a human `detail`:

```json
{
  "type": "urn:dx:acl:InvalidParamValue",
  "title": "Validation Error",
  "detail": "expiresAt must be in the future"
}
```

In code you never hand-build these. The handler pattern is decode → call service → translate:

```go
func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var req CreatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dxerrors.WriteError(w, dxerrors.NewValidation("invalid request body"))
		return
	}
	policy, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.log.Error("create policy failed", zap.Error(err))
		dxerrors.WriteError(w, err) // taxonomy type in the chain → status + URN
		return
	}
	dxresp.WriteCreated(w, policy)
}
```

`dxerrors.WriteError` walks the error chain (your `%w` discipline from Module 1 paying off): a wrapped `NewNotFound` becomes 404 + `…:ResourceNotFound`; anything unrecognized becomes a safe 500 without leaking internals. The service layer returns taxonomy errors; the handler stays a translator.

### Pagination, filtering, sorting

List endpoints accept `limit`/`offset` (bounded — a client may not ask for a million rows), whitelisted filters, and whitelisted sort keys. The response is the generic paged envelope — `Page[T]`/`DxPagedResponse[T]` from the [Generics](../module-1-go-fundamentals/generics) page — carrying `totalCount`, `limit`, `offset` alongside `results`. `dx-common-go`'s `request` package (`request.Builder`) parses and bounds these query parameters so every service interprets `?limit=25&offset=50&sort=-createdAt` identically. The whitelist matters for more than UX: sort keys reach SQL `ORDER BY`, and [Database Patterns](database-patterns) explains why unvalidated identifiers must never do that.

### OpenAPI-first validation

The embedded spec (from [Project Structure](project-structure)) isn't documentation-only — it's enforced. `dx-common-go/openapi` middleware validates every request against the spec before your handler runs:

```go
r.Use(dxopenapi.ValidationMiddleware(spec, cfg)) // after StandardStack, before auth groups
```

Malformed bodies, missing required fields, wrong enum values → rejected with a uniform validation error, **before** handler code. Consequences: handlers stay thin (structural validation is done; only business rules remain), and the spec can't rot (the `TestSpecCoversAllRoutes` test from [Testing](../module-2-intermediate/testing) fails when routes and spec diverge). Swagger UI serves the same spec at `/docs`.

:::info[Platform connection]
The full standard for a DX endpoint, assembled: path under `/iudx/v2/<domain>`, envelope via `response.ServiceWriter("urn:dx:<svc>:")` (each service writes its own URN prefix), errors via the `dx-common-go/errors` taxonomy, pagination via `request.Builder`, spec embedded and validated, `/docs` serving Swagger UI. Read one real handler — `dx-acl-go/internal/api/handler.go` — and you'll recognize every line from this page.
:::

## Exercises

*(Against your `dx-scratch-go` notes service, stack running.)*

1. Move your routes under `/iudx/v2/notes` and adopt envelope responses — write tiny `writeSuccess`/`writeError` helpers mimicking the platform shapes (or import `dx-common-go` directly if your module setup allows).
2. Implement `GET /iudx/v2/notes?limit=&offset=&sort=` with bounds (max limit 100, default 20), a sort-key whitelist, and the paged envelope. Add table-driven handler tests for the boundary cases.
3. Write an OpenAPI spec for the notes API, embed it, and wire kin-openapi validation middleware. Prove a missing required field never reaches your handler.
4. Add the two spec tests: embedded-spec-loads and spec-covers-all-routes. Add a route without updating the spec; watch the second fail.
5. Compare a real response: `curl` a Go-track endpoint through the gateway (Module 0 token) and match the envelope fields against what you built.

## Check yourself

- Who owns the API contract — the service or its clients — and what follows for a rewrite?
- Trace a `NewNotFound` from repository to the client's JSON: which layers touch it and what does each do?
- Why must sort keys be whitelisted?
- What breaks first when someone adds an endpoint and forgets the spec?

## References

- [OpenAPI 3 spec](https://spec.openapis.org/oas/v3.0.3) · [kin-openapi](https://pkg.go.dev/github.com/getkin/kin-openapi)
- [Google API design guide](https://cloud.google.com/apis/design) — general REST taste
- Platform: `dx-common-go/response`, `dx-common-go/errors`, `dx-common-go/request`; `claude-docs/CLIENT-CONTRACT-CHANGES.md` (the deltas clients must know)
