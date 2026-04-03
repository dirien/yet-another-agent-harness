<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# HTTP Service Patterns for Agent-Ready Go

## Graceful Shutdown

Services that don't handle SIGTERM drop in-flight requests when containers
restart. Use `signal.NotifyContext` to listen for OS signals, then call
`srv.Shutdown` to drain connections before exiting:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := &http.Server{
        Addr:    ":8080",
        Handler: buildRouter(),
    }

    // Start server in a goroutine so main can wait for shutdown signal.
    go func() {
        slog.Info("server listening", slog.String("addr", srv.Addr))
        if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            fmt.Fprintf(os.Stderr, "server error: %v\n", err)
            os.Exit(1)
        }
    }()

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    <-ctx.Done()
    slog.Info("shutdown signal received")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        fmt.Fprintf(os.Stderr, "graceful shutdown failed: %v\n", err)
        os.Exit(1)
    }
    slog.Info("server stopped cleanly")
}
```

## Health Check Endpoints

Expose `/healthz` (liveness) and `/readyz` (readiness) on every HTTP service.
Kubernetes and load balancers use these to route traffic:

- **liveness** (`/healthz`) — is the process alive and not deadlocked?
- **readiness** (`/readyz`) — is the process ready to serve traffic (DB
  connected, caches warm)?

```go
type HealthResponse struct {
    Status  string            `json:"status"`           // "ok" | "degraded" | "down"
    Version string            `json:"version"`
    Checks  map[string]string `json:"checks,omitempty"` // "db": "ok", "cache": "degraded"
}

func livenessHandler(version string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(HealthResponse{
            Status:  "ok",
            Version: version,
        })
    }
}

func readinessHandler(db *sql.DB, version string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        checks := map[string]string{}
        status := "ok"

        if err := db.PingContext(r.Context()); err != nil {
            checks["db"] = "down"
            status = "down"
        } else {
            checks["db"] = "ok"
        }

        code := http.StatusOK
        if status == "down" {
            code = http.StatusServiceUnavailable
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        _ = json.NewEncoder(w).Encode(HealthResponse{
            Status:  status,
            Version: version,
            Checks:  checks,
        })
    }
}
```

Register them on a dedicated internal mux so they are never blocked by
authentication middleware:

```go
mux.HandleFunc("/healthz", livenessHandler(version))
mux.HandleFunc("/readyz", readinessHandler(db, version))
```

## Request ID Middleware

Propagate `X-Request-ID` through every request so logs, traces, and error
responses can be correlated across services. Generate a UUID if the header is
absent:

```go
import (
    "context"
    "log/slog"
    "net/http"

    "github.com/google/uuid"
)

type requestIDKey struct{}

func RequestIDMiddleware(log *slog.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = uuid.New().String()
        }
        w.Header().Set("X-Request-ID", id)
        ctx := context.WithValue(r.Context(), requestIDKey{}, id)
        reqLog := log.With(slog.String("request_id", id))
        next.ServeHTTP(w, WithLogger(ctx, reqLog))
    })
}

func RequestIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey{}).(string); ok {
        return id
    }
    return ""
}
```

Install:

```bash
go get github.com/google/uuid
```

## OpenTelemetry Tracing

Add distributed tracing so agents and operators can follow a request across
service boundaries.

Install dependencies:

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
```

Initialize the tracer provider at startup:

```go
import (
    "context"
    "log/slog"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
    exp, err := otlptracegrpc.New(ctx)
    if err != nil {
        return nil, err
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(newResource(serviceName)),
    )
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})
    return tp.Shutdown, nil
}
```

Instrument a handler with a span:

```go
tracer := otel.Tracer("order-service")

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "HandleCreateOrder")
    defer span.End()

    // Attach trace/span IDs to log lines for correlation.
    sc := span.SpanContext()
    logger.FromContext(ctx).Info("handling request",
        slog.String("trace_id", sc.TraceID().String()),
        slog.String("span_id", sc.SpanID().String()),
    )

    // ... handler logic using ctx
}
```
