<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Configuration Management for Agent-Ready Go

## Priority Order

Apply configuration in this order — later sources override earlier ones:

1. **Default value** — compiled into the binary
2. **Config file** — optional, for complex multi-field configuration
3. **Environment variable** — the primary mechanism for agents and containers
4. **Explicit flag** — highest precedence, passed directly by the caller

For most services and CLI tools, environment variables alone are sufficient.
Reach for a config file only when you have more than ~10 settings.

## Struct-Based Config

Define a typed struct with `env` tags and parse it with
`github.com/caarlos0/env/v11` — a lightweight library with no external
dependencies and good error messages:

```bash
go get github.com/caarlos0/env/v11
```

```go
package config

import (
    "errors"
    "fmt"
    "time"

    "github.com/caarlos0/env/v11"
)

type Config struct {
    Port     int           `env:"PORT"          envDefault:"8080"`
    LogLevel string        `env:"LOG_LEVEL"     envDefault:"info"`
    DBURL    string        `env:"DATABASE_URL"  required:"true"`
    Timeout  time.Duration `env:"TIMEOUT"       envDefault:"30s"`
    Debug    bool          `env:"DEBUG"         envDefault:"false"`
}

func Load() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("loading config: %w", err)
    }
    return cfg, nil
}
```

All fields are type-safe — durations and booleans are parsed automatically.
Missing `required` fields produce a clear error message that agents can act on.

## Validation

Add a `Validate` method to catch logical errors that the `env` parser cannot
detect (range checks, cross-field constraints, etc.):

```go
func (c *Config) Validate() error {
    var errs []error
    if c.Port < 1 || c.Port > 65535 {
        errs = append(errs, fmt.Errorf("PORT %d out of range [1-65535]", c.Port))
    }
    if c.DBURL == "" {
        errs = append(errs, errors.New("DATABASE_URL is required"))
    }
    return errors.Join(errs...)
}
```

## Config in main()

Load and validate configuration at the very top of `main`. Fail fast with exit
code `2` so the process never starts partially configured — agents see an
immediate, unambiguous failure:

```go
package main

import (
    "fmt"
    "os"

    "github.com/example/myapp/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "config error: %v\n", err)
        os.Exit(2)
    }
    if err := cfg.Validate(); err != nil {
        fmt.Fprintf(os.Stderr, "invalid config: %v\n", err)
        os.Exit(2)
    }

    // cfg is valid — pass it to the rest of the application.
    if err := run(cfg); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
```

## When to Use Viper

[Viper](https://github.com/spf13/viper) is appropriate when you need **all
three** of: config files (TOML/YAML/JSON), environment variables, **and** CLI
flags in a single unified system — typically large CLIs or daemon processes with
rich configuration.

For simpler cases, prefer the struct-based approach above:

| Scenario | Recommendation |
| -------- | -------------- |
| Container / 12-factor app (env vars only) | `caarlos0/env` struct |
| CLI with a few flags and env var overrides | `pflag` / `cobra` + env defaults |
| Daemon with a config file + env overrides | Viper |
| Multi-command CLI with subcommand configs | Viper |

Viper example (minimal):

```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.AutomaticEnv() // env vars override file values

if err := viper.ReadInConfig(); err != nil {
    // config file is optional — only fail on parse errors
    var notFound viper.ConfigFileNotFoundError
    if !errors.As(err, &notFound) {
        return fmt.Errorf("reading config: %w", err)
    }
}

port := viper.GetInt("port")
```
