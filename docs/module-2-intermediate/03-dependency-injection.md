---
title: Dependency Injection
sidebar_label: Dependency Injection
description: Constructor injection, manual composition, lifecycle ownership, and platform bootstrap.
---

# Dependency injection

## Outcomes

You will build an object graph without globals, keep interfaces near consumers, and give one composition root lifecycle ownership.

## Constructor injection

~~~go
type Service struct {
    widgets WidgetStore
    clock   Clock
    log     *zap.Logger
}

func New(
    widgets WidgetStore,
    clock Clock,
    log *zap.Logger,
) *Service {
    if widgets == nil {
        panic("widgets is required")
    }
    return &Service{widgets: widgets, clock: clock, log: log}
}
~~~

Required dependencies are explicit and immutable after construction. Tests provide small fakes. Avoid package-level mutable clients and service locators.

## Manual composition

Go wiring is ordinary code. That makes construction order, failure policy, and cleanup visible:

~~~go
func wire(
    _ context.Context,
    app *bootstrap.App[config.Config],
) (http.Handler, error) {
    repo := postgres.NewWidgetRepository(app.DB)
    svc := service.New(repo, systemClock{}, app.Log)
    handler := api.New(svc)

    return api.Router(app, handler), nil
}
~~~

Only cmd/server imports every layer. Domain and application packages cannot reach outward to build a dependency.

## Bootstrap owns the process

~~~go
bootstrap.Run(bootstrap.Spec[config.Config]{
    Name:   "dx-widget-go",
    Config: config.Options(),
    Deps: func(c *config.Config) bootstrap.Deps {
        return bootstrap.Deps{
            Migrations: bootstrap.Migrations(
                widgetdb.Migrations,
                "migrations",
                "schema_migrations_widget",
            ),
            Postgres: bootstrap.Required(c.Postgres),
        }
    },
    Wire: wire,
})
~~~

Bootstrap fixes startup and shutdown order. Wire constructs domain-specific clients and registers closers, probes, and workers.

## Optional dependency test

An optional dependency is valid only when the service has a correct degraded behavior. For example, a broker can be temporarily unavailable if a database transaction still commits an outbox row for later publication. “The process starts” is not proof of a valid degraded mode.

## Exercise

Refactor a handler that opens PostgreSQL directly. Define a store interface, inject it into a service, create the adapter in wire, and write a unit test with a fake. Draw which component closes the pool.

## Check yourself

- Why is cmd/server allowed to import all layers?
- Who owns closing an injected client?
- When should a constructor panic?
- What evidence justifies bootstrap.Degrade?
