<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Structured Logging for Agent-Ready Go

Choose **one** library. All four support JSON output. `log/slog` is the
recommended default for new projects (stdlib, Go 1.21+); Zap is fastest for
performance-critical code; ZeroLog is zero-alloc; Logrus is best for existing
codebases.

---

## log/slog (Recommended — stdlib, Go 1.21+)

No installation required — `log/slog` is part of the Go standard library.

### slog setup

```go
package logger

import (
    "log/slog"
    "os"
    "strings"

    "golang.org/x/term"
)

func New() *slog.Logger {
    levelVar := &slog.LevelVar{}
    if raw := os.Getenv("LOG_LEVEL"); raw != "" {
        switch strings.ToLower(raw) {
        case "debug":
            levelVar.Set(slog.LevelDebug)
        case "warn":
            levelVar.Set(slog.LevelWarn)
        case "error":
            levelVar.Set(slog.LevelError)
        default:
            levelVar.Set(slog.LevelInfo)
        }
    }

    opts := &slog.HandlerOptions{Level: levelVar}
    var h slog.Handler
    if term.IsTerminal(int(os.Stderr.Fd())) && os.Getenv("LOG_FORMAT") != "json" {
        h = slog.NewTextHandler(os.Stderr, opts)
    } else {
        h = slog.NewJSONHandler(os.Stderr, opts)
    }
    return slog.New(h)
}
```

### slog usage

```go
log := logger.New()
slog.SetDefault(log)

log.Info("server started", slog.Int("port", 8080), slog.String("env", "prod"))
log.Error("request failed", slog.String("path", r.URL.Path), "err", err)
```

JSON output example:

```json
{"time":"2024-03-25T10:12:32.123Z","level":"INFO","msg":"server started","port":8080,"env":"prod"}
```

---

## Zap (Recommended for performance-critical code)

```bash
go get go.uber.org/zap
```

### Zap setup

```go
package logger

import (
    "fmt"
    "os"
    "strings"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "golang.org/x/term"
)

func New() *zap.Logger {
    var cfg zap.Config
    if term.IsTerminal(int(os.Stderr.Fd())) && os.Getenv("LOG_FORMAT") != "json" {
        cfg = zap.NewDevelopmentConfig()
    } else {
        cfg = zap.NewProductionConfig()
    }

    if raw := os.Getenv("LOG_LEVEL"); raw != "" {
        var lvl zapcore.Level
        if err := lvl.UnmarshalText([]byte(strings.ToLower(raw))); err == nil {
            cfg.Level.SetLevel(lvl)
        }
    }

    logger, err := cfg.Build()
    if err != nil {
        panic(fmt.Sprintf("failed to build logger: %v", err))
    }
    return logger
}
```

### Zap usage

```go
log := logger.New()
defer log.Sync()

log.Info("server started", zap.Int("port", 8080), zap.String("env", "prod"))
log.Error("request failed", zap.Error(err), zap.String("path", r.URL.Path))
```

JSON output example:

```json
{"level":"info","ts":1711330352.123,"msg":"server started","port":8080,"env":"prod"}
```

---

## ZeroLog (Zero allocation)

```bash
go get github.com/rs/zerolog
```

### ZeroLog setup

```go
package logger

import (
    "os"
    "strings"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "golang.org/x/term"
)

func Init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    if raw := os.Getenv("LOG_LEVEL"); raw != "" {
        switch strings.ToLower(raw) {
        case "debug":
            zerolog.SetGlobalLevel(zerolog.DebugLevel)
        case "warn":
            zerolog.SetGlobalLevel(zerolog.WarnLevel)
        case "error":
            zerolog.SetGlobalLevel(zerolog.ErrorLevel)
        default:
            zerolog.SetGlobalLevel(zerolog.InfoLevel)
        }
    }

    if term.IsTerminal(int(os.Stderr.Fd())) && os.Getenv("LOG_FORMAT") != "json" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }
    // else: default JSON to stderr
}
```

### ZeroLog usage

```go
logger.Init()
log.Info().Int("port", 8080).Str("env", "prod").Msg("server started")
log.Error().Err(err).Str("path", r.URL.Path).Msg("request failed")
```

---

## Logrus (Existing/legacy codebases)

```bash
go get github.com/sirupsen/logrus
```

### Logrus setup

```go
package logger

import (
    "os"

    "github.com/sirupsen/logrus"
    "golang.org/x/term"
)

func Init() {
    if !term.IsTerminal(int(os.Stderr.Fd())) || os.Getenv("LOG_FORMAT") == "json" {
        logrus.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
    }
    logrus.SetOutput(os.Stderr)

    if raw := os.Getenv("LOG_LEVEL"); raw != "" {
        lvl, err := logrus.ParseLevel(raw)
        if err == nil {
            logrus.SetLevel(lvl)
        }
    }
}
```

### Logrus usage

```go
logger.Init()
logrus.WithFields(logrus.Fields{"port": 8080, "env": "prod"}).Info("server started")
logrus.WithError(err).WithField("path", r.URL.Path).Error("request failed")
```

---

## Logger-from-context Pattern

Store an enriched logger in context so handlers automatically carry request-scoped
fields (request ID, user ID, trace ID) without explicit parameter threading:

```go
package logger

import (
    "context"
    "log/slog"
)

type loggerKey struct{}

// WithLogger stores l in ctx.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext retrieves the logger stored by WithLogger.
// Falls back to slog.Default() when none is present.
func FromContext(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}
```

Usage in a handler:

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log := logger.FromContext(r.Context())
    log.Info("handling request", slog.String("method", r.Method))
}
```

---

## Base Fields Pattern

Attach service name, version, and environment to every log line so log
aggregators can filter without parsing message content:

```go
log := slog.New(h).With(
    slog.String("service", "order-service"),
    slog.String("version", version),
    slog.String("env", os.Getenv("APP_ENV")),
)
slog.SetDefault(log)
```

Every subsequent call to `slog.Info(...)` automatically includes those fields.

---

## Common Requirements (All Libraries)

1. **Never** use `fmt.Println` for operational logs — agents cannot correlate
   unstructured lines with structured events.
2. Log to **stderr** by default; keep stdout clean for machine-readable output.
3. Use structured fields, not interpolated strings:
   - ❌ `log.Infof("user %s logged in at %s", userID, time)`
   - ✅ `log.Info("user logged in", slog.String("user_id", userID), slog.Time("at", t))`
4. Support `LOG_LEVEL` env var: `debug`, `info`, `warn`, `error`.
5. Add request/trace IDs as log fields for HTTP/gRPC handlers.

## Log Level Guidelines

| Level   | When to use                                     |
| ------- | ----------------------------------------------- |
| `debug` | Detailed internal state, disabled in production |
| `info`  | Normal operations agents should observe         |
| `warn`  | Recoverable issues, degraded behavior           |
| `error` | Failures requiring attention                    |
| `fatal` | Only in init paths — never in request handlers  |
