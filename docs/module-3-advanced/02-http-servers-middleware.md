---
title: HTTP Servers and Middleware
sidebar_label: HTTP & Middleware
description: Typed handlers, route declarations, middleware order, timeouts, streams, and shutdown.
---

# HTTP servers and middleware

## Typed handler

~~~go
func (h *Handler) Get(
    ctx context.Context,
    req GetRequest,
) (GetResponse, error) {
    widget, err := h.widgets.Get(ctx, req.Actor.Subject, req.ID)
    if err != nil {
        return GetResponse{}, fmt.Errorf("get widget: %w", err)
    }
    return toResponse(widget), nil
}
~~~

GetRequest embeds httpx.Actor and binds its ID from the route. Business code sees context and typed values, not ResponseWriter or router context.

## Declarative route

~~~go
httpx.GET(
    "/{id}",
    httpx.Handle(handler.Get),
    httpx.OpID("getWidget"),
)
~~~

Route options explicitly select Public, Optional identity, Roles, Streaming, and operation ID. A route table can be compared to OpenAPI in a test.

## Standard middleware

NewRouter owns:

1. panic recovery;
2. trace context;
3. request ID;
4. trusted client address;
5. structured request logging;
6. CORS;
7. service middleware;
8. authentication and role gates;
9. per-route timeout and compression.

Health, metrics, and docs are mounted outside the business API base. Streaming routes bypass request timeout and compression because both can break SSE or large streaming bodies.

## Server lifecycle

platform/bootstrap starts net/http with typed server timeouts, reacts to SIGINT/SIGTERM, drains requests, stops workers, and closes dependencies. Handlers must pass the request context to every blocking call.

## Middleware rules

- A middleware either enriches context, enforces policy, observes, or transforms; keep one responsibility.
- Never trust identity or forwarding headers before the verification middleware.
- Avoid reading an unbounded body in middleware.
- Preserve status and interfaces needed by streaming responses.
- Do not use global mutable state.

## Exercise

Create public `/healthz/live`, protected `/v1/widgets/{id}`, optional `/v1/widgets` search, and streaming `/v1/events` routes. Write tests proving which handler runs for anonymous and forged identity requests.

## Check yourself

- Why are operational endpoints outside auth?
- What breaks when compression wraps SSE?
- Which layer renders an application error?
- Why are routes data rather than scattered registration calls?
