---
title: Configuration
sidebar_label: Configuration
description: Typed config, precedence, validation, environment mapping, and secrets.
---

# Configuration

## Outcomes

You will load typed settings, reject invalid state at boot, and separate safe defaults from secrets.

## Typed shape

~~~go
type Config struct {
    platformconfig.Base
    Postgres dxsql.Config
    Catalogue Catalogue
    Workload workload.VerifierConfig
}

type Catalogue struct {
    BaseURL string
    Timeout time.Duration
}

func (c Config) Validate() error {
    if c.Catalogue.BaseURL == "" {
        return errors.New("catalogue.base_url is required")
    }
    return c.Workload.Validate()
}
~~~

The actual fields carry mapstructure tags. Keep units in types: time.Duration rather than an integer “seconds,” URLs rather than unrelated string fragments.

## Precedence

platform/config applies:

1. platform defaults;
2. service defaults;
3. optional config file;
4. environment variables.

Nested keys use underscores in the environment: CATALOGUE_BASE_URL and POSTGRES_MAX_CONNS. A missing production file is normal; a malformed file is an error.

## Validation

Validate required values, URL schemes, positive durations, size bounds, secret placeholders, and feature combinations. Fail before opening dependencies. Do not defer obvious misconfiguration until the first request.

## Secrets

Configuration names a secret; source control does not contain its deployment value. Use ExternalSecret or the deployment's secret provider. Do not log complete config structs because they contain DSNs, credentials, tokens, and keys.

## No hot reload

The platform treats a rollout as configuration activation. This keeps replicas on a known configuration and gives GitOps or deployment history an auditable change.

## Exercise

Add a typed outbound-client configuration with BaseURL, Timeout, and MaxRetries. Validate URL and bounds. Write table-driven tests for defaults, file override, environment override, missing file, malformed file, and invalid configuration.

## Check yourself

- Which source has highest precedence?
- Why is a missing file allowed but malformed YAML fatal?
- What makes a dependency optional rather than merely disabled?
- Why should config changes restart the process?
