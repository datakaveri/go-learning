---
title: Interfaces
sidebar_label: Interfaces
description: Implicit satisfaction, consumer-owned ports, type assertions, and interface design.
---

# Interfaces

## Outcomes

You will define small interfaces at the point of use, inject implementations, and avoid interfaces that only mirror a concrete type.

## Implicit satisfaction

~~~go
type Clock interface {
    Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
    return time.Now().UTC()
}

var _ Clock = SystemClock{}
~~~

There is no implements declaration. A type satisfies an interface by having its method set. The compile-time assertion is optional but useful at adapter boundaries.

## Define the port beside its consumer

~~~go
package service

type WidgetStore interface {
    Get(context.Context, string) (Widget, error)
    Save(context.Context, Widget) error
}

type Service struct {
    widgets WidgetStore
}

func New(widgets WidgetStore) *Service {
    return &Service{widgets: widgets}
}
~~~

The application names only the behavior it needs. A PostgreSQL adapter and a test fake can both satisfy it without the service importing either.

## Method sets

A value of type T has methods with value receivers. A pointer *T has methods with both pointer and value receivers. If an interface requires a method implemented only on *T, pass a pointer.

Use a pointer receiver when a method mutates the value, the type is large, or consistency requires it. Do not mix receivers without a reason.

## Assertions and type switches

~~~go
checker, ok := dependency.(health.Checker)
if !ok {
    return errors.New("dependency has no health check")
}

switch value := input.(type) {
case string:
    return parseString(value)
case []byte:
    return parseBytes(value)
default:
    return fmt.Errorf("unsupported input %T", input)
}
~~~

Use comma-ok unless a failed assertion is a proven programmer invariant. Avoid interface{} or any when a type parameter or small interface can express the contract.

## Interface design rules

- Keep interfaces small and cohesive.
- Accept interfaces; return useful concrete types or existing contracts.
- Do not create an interface merely to mock a concrete dependency.
- Include context on operations that can block.
- Make lifecycle explicit; an owner closes the resource.
- Do not hide vendor-specific options behind generic maps.

## Platform connection

dx-common-go uses this design throughout: cache.Store, events.Bus, sql.DB, health.Checker, and service-defined repositories. Read one interface and its primary adapter, then identify which package owns construction.

## Exercise

Implement an in-memory WidgetStore and a service that rejects duplicate names. Write tests using only the interface. Then list the smallest methods a PostgreSQL adapter actually needs.

## Check yourself

- Does *T satisfy every interface T satisfies?
- Who should own a repository interface: its consumer or its driver?
- When is a type assertion appropriate?
- Why is returning a huge vendor-shaped interface a coupling problem?
