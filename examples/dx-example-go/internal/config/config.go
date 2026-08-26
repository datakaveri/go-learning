package config

import (
	"errors"
	"time"

	"github.com/datakaveri/dx-common-go/auth/fga"
	dxjwt "github.com/datakaveri/dx-common-go/auth/jwt"
	platformconfig "github.com/datakaveri/dx-common-go/platform/config"
	dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"
	"github.com/datakaveri/dx-common-go/platform/events/amqp"
	"github.com/datakaveri/dx-common-go/platform/security/workload"
	"github.com/datakaveri/dx-common-go/platform/security/workload/issuer"
)

type Config struct {
	platformconfig.Base `mapstructure:",squash"`
	Postgres            dxsql.Config            `mapstructure:"postgres"`
	RabbitMQ            amqp.Config             `mapstructure:"rabbitmq"`
	JWT                 dxjwt.Config            `mapstructure:"jwt"`
	WorkloadVerifier    workload.VerifierConfig `mapstructure:"workload_verifier"`
	WorkloadIssuer      issuer.Config           `mapstructure:"workload_issuer"`
	Authorization       fga.Config              `mapstructure:"authorization"`
}

func Options() platformconfig.Options {
	return platformconfig.Options{Defaults: map[string]any{
		"postgres.search_path":   "example",
		"postgres.max_conns":     10,
		"rabbitmq.exchange":      "example",
		"rabbitmq.exchange_type": "topic",
		"rabbitmq.prefetch":      16,
		"authorization.timeout":  2 * time.Second,
	}}
}

func (c Config) Validate() error {
	if c.Postgres.DSN == "" {
		return errors.New("postgres.dsn is required")
	}
	if c.RabbitMQ.URL == "" || c.RabbitMQ.Exchange == "" {
		return errors.New("rabbitmq.url and rabbitmq.exchange are required")
	}
	if c.Authorization.BaseURL == "" {
		return errors.New("authorization.base_url is required")
	}
	if err := c.WorkloadVerifier.Validate(); err != nil {
		return err
	}
	if !c.WorkloadIssuer.Enabled {
		return errors.New("workload_issuer.enabled must be true for authorization calls")
	}
	return c.WorkloadIssuer.Validate()
}
