---
title: Generics
sidebar_label: Generics
description: Type parameters, constraints, inference, zero values, and platform examples.
---

# Generics

## Outcomes

You will recognize when a type parameter removes duplication without erasing domain meaning.

## Functions and types

~~~go
func Map[A, B any](in []A, fn func(A) B) []B {
    out := make([]B, 0, len(in))
    for _, value := range in {
        out = append(out, fn(value))
    }
    return out
}

type Page[T any] struct {
    Items []T
    Total int64
}
~~~

The compiler often infers type arguments from ordinary arguments. Constraints describe supported operations; use comparable for map keys or equality and a type set only when algorithms require operators.

## Zero values

~~~go
func First[T any](values []T) (T, bool) {
    var zero T
    if len(values) == 0 {
        return zero, false
    }
    return values[0], true
}
~~~

Return an explicit boolean or error when the zero value is ambiguous.

## Platform examples

- platform/http.Handler[Req, Res] keeps handlers typed.
- platform/paging.Page[T] provides one client paging shape.
- platform/database/sql.Repo[T] scans and persists row types.
- platform/events.Topic[T] binds event name, version, and payload type.
- cache.GetOrLoad[T] returns typed cached values.
- bootstrap.Spec[C] carries a typed service config.

Each generic solves one repeated mechanical pattern while preserving an explicit semantic boundary.

## When not to use generics

Prefer an interface when implementations vary by behavior. Prefer a concrete domain function when two types merely look structurally similar. Avoid a generic repository that grows joins, workflows, validation, authorization, and auditing; those are different concerns.

## Exercise

Implement Map and Filter. Then use paging.MapPage to convert storage rows to API responses while preserving metadata. Compare the code and decide which implementation belongs in application code.

## Check yourself

- How do constraints differ from interfaces used as values?
- When does the compiler infer T?
- Why does Page[T] help more than a result typed any?
- What is the warning sign that a generic abstraction is swallowing domain behavior?
