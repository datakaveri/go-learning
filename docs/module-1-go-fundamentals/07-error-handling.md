---
title: Error Handling
sidebar_label: Error Handling
description: Error values, wrapping, classification, single handling, and safe client translation.
---

# Error handling

## Outcomes

You will add operation context with %w, classify errors once, preserve cancellation, and translate only at process boundaries.

## Return and wrap

~~~go
widget, err := repo.Get(ctx, id)
if err != nil {
    return Widget{}, fmt.Errorf("get widget %q: %w", id, err)
}
~~~

Wrapping preserves the cause for errors.Is and errors.As. Use lower-case operation context; avoid punctuation and redundant “failed to.”

## Inspect the chain

~~~go
if errors.Is(err, context.Canceled) {
    return err
}

var validation *ValidationError
if errors.As(err, &validation) {
    return platformerrors.Validation(validation.Message)
}
~~~

Sentinel errors represent stable categories. Typed errors carry structured fields. Do not compare error strings.

## Handle once

At each layer, choose one:

- add context and return;
- translate to a domain or platform classification and return;
- recover locally with a documented fallback;
- at the top boundary, log and render or terminate.

Do not log the same error at repository, service, handler, and gateway layers.

## Platform classification

~~~go
if name == "" {
    return platformerrors.Validation("name is required")
}

if err := repo.Delete(ctx, id); err != nil {
    return platformerrors.Wrap(
        err,
        platformerrors.CodeDatabase,
        "could not delete widget",
    )
}
~~~

platform/errors constructors include Validation, Unauthorized, Forbidden, NotFound, Conflict, Internal, BadGateway, ServiceUnavailable, TooManyRequests, Expired, Database, and MethodNotAllowed.

platform/http and platform/grpc translate the same taxonomy. Unclassified errors become a safe generic server error; the internal cause is logged. Never send err.Error to a client.

## Panic

Panic is for impossible programmer invariants or unrecoverable initialization, not ordinary input or dependency failure. platform/http recovery prevents a process crash, but a recovered panic remains a defect that needs alerting.

## Exercise

Build a three-layer call: repository returns sql.ErrNoRows, application translates it to NotFound, handler returns it unchanged. Test errors.Is on the original cause and the HTTP problem status.

## Check yourself

- What information does %w preserve?
- Where should a database constraint become Conflict?
- Why is “log and return” usually double handling?
- Which errors must retain context cancellation semantics?
