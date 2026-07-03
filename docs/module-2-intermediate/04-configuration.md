---
title: Configuration Management
sidebar_label: Configuration
description: Env-driven config with viper and mapstructure — the path to dxconfig.LoadService[T].
---

# Configuration Management

## Learning objectives

- Model configuration as a typed struct with `mapstructure` tags.
- Layer sources the platform way: defaults → YAML file → environment variables.
- Validate configuration at boot and fail fast with a good message.
- Keep secrets out of code and out of baked files.

## Prerequisites

- [Structs & Methods](../module-1-go-fundamentals/structs-methods) (tags), [Dependency Injection](dependency-injection)

## Time estimate

**2.5 hours**

## Concepts

### Config is a struct, not a bag

```go
type Config struct {
	Server   ServerConfig     `mapstructure:"server"`
	Postgres postgres.Config  `mapstructure:"postgres"`
	RMQ      rabbitmq.Config  `mapstructure:"rmq"`
}

type ServerConfig struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}
```

Typed config means typos fail loudly, IDEs autocomplete, and each component receives only its slice (`cfg.Postgres`), per the DI page. Note the nested shared types — services reuse `dx-common-go`'s `postgres.Config` etc. rather than redeclaring host/port/user fields.

### Layered sources with viper

The platform uses [viper](https://github.com/spf13/viper) with a fixed precedence, lowest to highest:

1. **Defaults** — a `map[string]any` in code (`"server.port": 8080`), so a bare binary runs in dev.
2. **YAML file** — `configs/config.yaml`, baked into the image, holding non-secret dev defaults.
3. **Environment variables** — the production override mechanism (`SERVER_PORT=9000`); dots become underscores.

```go
func Load() (*Config, error) {
	v := viper.New()
	v.SetDefault("server.port", 8080)
	v.SetConfigFile("configs/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, cfg.Validate()
}
```

Env-var overriding is what makes the same image run in dev, staging, and production — twelve-factor style. The platform's policy is **unprefixed** env names (`POSTGRES_HOST`, not `DXACL_POSTGRES_HOST`), because each service runs in its own container with its own environment.

### Validate at boot

A config error discovered at request time is an outage; discovered at boot it's a log line. Every service implements `Validate()`:

```go
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port %d out of range", c.Server.Port)
	}
	if c.Postgres.Host == "" {
		return errors.New("postgres.host is required")
	}
	return nil
}
```

### Secrets

Two rules, no exceptions:

1. **No secrets in code or in the baked `config.yaml`.** The YAML holds an empty default; the real value arrives via environment (injected by Compose locally, by External Secrets Operator in Kubernetes).
2. **Never log config wholesale.** Log the non-sensitive summary; redact DSNs and keys.

:::info[Platform connection]
What you built by hand above is `dxconfig.LoadService[T]` from `dx-common-go/config` — the generic loader every service calls in one line:

```go
func Load() (*Config, error) {
	return dxconfig.LoadService[Config](dxconfig.ServiceOptions{
		Defaults: map[string]any{"server.port": "8080", "postgres.max_conns": 10},
	})
}
```

Same precedence (defaults → baked YAML → env), plus your `Validate()`. GO-SERVICE-STANDARDS requires this loader — hand-rolled viper code in a service is a review finding. (Historical note: `dx-catalogue-go` uses a nonstandard `DX` env prefix; it's documented as a deviation, not an example to copy.)
:::

## Exercises

1. Build the `Load()` above for a toy service with server + database sections. Prove each precedence layer: run bare (defaults), with a YAML file, and with `SERVER_PORT=9999` overriding both.
2. Add `Validate()` with three rules and confirm the process exits at boot with a helpful message when violated.
3. Add a `Password string` field: empty in YAML, injected via env, **asserted absent** from your startup log line.
4. Read `dx-common-go/config` and compare with your version — list two things the shared loader handles that yours doesn't.

## Check yourself

- Recite the three-layer precedence, lowest to highest.
- Why are env vars the production override rather than mounted files or flags?
- Where do secrets live in dev? In production?
- Why unprefixed env names on this platform?

## References

- [viper](https://github.com/spf13/viper) · [mapstructure](https://pkg.go.dev/github.com/go-viper/mapstructure/v2)
- [The Twelve-Factor App — Config](https://12factor.net/config)
- Platform: `dx-common-go/config`, `claude-docs/CONFIG.md` (Go section), GO-SERVICE-STANDARDS.md (config rules)
