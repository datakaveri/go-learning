package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/datakaveri/dx-common-go/auth/fga"
	"github.com/datakaveri/dx-common-go/metrics"
	"github.com/datakaveri/dx-common-go/platform/bootstrap"
	platformevents "github.com/datakaveri/dx-common-go/platform/events"
	"github.com/datakaveri/dx-common-go/platform/events/amqp"
	httpx "github.com/datakaveri/dx-common-go/platform/http"
	"github.com/datakaveri/dx-common-go/platform/http/middleware"
	"github.com/datakaveri/dx-common-go/platform/security/workload"
	"github.com/datakaveri/dx-common-go/platform/security/workload/issuer"

	exampledb "github.com/datakaveri/dx-example-go/db"
	"github.com/datakaveri/dx-example-go/internal/application"
	"github.com/datakaveri/dx-example-go/internal/config"
	"github.com/datakaveri/dx-example-go/internal/integration/authz"
	exampleevents "github.com/datakaveri/dx-example-go/internal/integration/events"
	"github.com/datakaveri/dx-example-go/internal/repository/postgres"
	httptransport "github.com/datakaveri/dx-example-go/internal/transport/http"
)

func main() {
	bootstrap.Run(bootstrap.Spec[config.Config]{
		Name:   "dx-example-go",
		Config: config.Options(),
		Deps: func(cfg *config.Config) bootstrap.Deps {
			return bootstrap.Deps{
				Migrations: bootstrap.Migrations(exampledb.Migrations, "migrations", "schema_migrations_example"),
				Postgres:   bootstrap.Required(cfg.Postgres),
			}
		},
		Wire: wire,
	})
}

func wire(_ context.Context, app *bootstrap.App[config.Config]) (http.Handler, error) {
	verifier, err := workload.FromConfig(app.Cfg.WorkloadVerifier)
	if err != nil {
		return nil, fmt.Errorf("workload verifier: %w", err)
	}

	workloadSource, err := issuer.New(app.Cfg.WorkloadIssuer)
	if err != nil {
		return nil, fmt.Errorf("workload issuer: %w", err)
	}
	authzConfig := app.Cfg.Authorization
	authzConfig.Workload = workloadSource
	authzClient, err := fga.New(authzConfig)
	if err != nil {
		return nil, fmt.Errorf("authorization client: %w", err)
	}

	bus, err := amqp.Open(app.Cfg.RabbitMQ)
	if err != nil {
		return nil, fmt.Errorf("event bus: %w", err)
	}
	bus.WithLogger(app.Log)
	app.Closer("event bus", func(context.Context) error { return bus.Close() })

	outbox := platformevents.NewOutbox(app.DB, "example_outbox")
	dispatcher := platformevents.NewDispatcher(outbox, bus, 100, eventLogger{app.Log})
	app.Background("outbox dispatcher", func(ctx context.Context) error {
		return dispatcher.Run(ctx, time.Second)
	})

	service, err := application.New(
		postgres.NewWidgetRepository(app.DB), authz.New(authzClient),
		exampleevents.NewWriter(outbox, dispatcher), app.Tx, time.Now,
	)
	if err != nil {
		return nil, err
	}
	handler := httptransport.NewHandler(service)

	return httpx.NewRouter(httpx.RouterSpec{
		Base: "/iudx/example/v1", URNs: httptransport.URNs,
		Health: app.Health, Metrics: metrics.Handler(), Logger: app.Log,
		Auth: httpx.AuthSpec{Authenticate: middleware.Resolve(middleware.AuthConfig{
			TrustSubjectHeaders: true,
			Workload:            verifier,
			JWT:                 app.Cfg.JWT,
		})},
	}, httptransport.Routes(handler)), nil
}

type eventLogger struct{ log *zap.Logger }

func (l eventLogger) Warn(msg string, fields ...any)  { l.log.Sugar().Warnw(msg, fields...) }
func (l eventLogger) Error(msg string, fields ...any) { l.log.Sugar().Errorw(msg, fields...) }
