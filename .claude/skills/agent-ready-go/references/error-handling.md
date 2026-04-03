<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Error Handling for Agent-Ready Go

Agents determine success or failure by exit codes and structured output.
Swallowed errors, panics, and ambiguous exit codes are the main cause of agents
getting stuck in retry loops.

## Wrapping Errors (Required)

Always wrap errors with context so the full call chain is traceable:

```go
// ❌ Bad — loses context
return err

// ✅ Good — agents can trace the full failure chain
return fmt.Errorf("loading user %s: %w", userID, err)
```

Check wrapped error types with `errors.As` / `errors.Is`:

```go
var notFound *NotFoundError
if errors.As(err, &notFound) {
    // handle specifically
}

if errors.Is(err, ErrPermissionDenied) {
    os.Exit(4)
}
```

Never use `== err` for sentinel error comparison — it breaks with wrapped errors.

## Structured Error Types

For CLIs and APIs used by agents, define typed errors with machine-readable
fields:

```go
type AppError struct {
    Code    string `json:"code"`             // machine-readable: "NOT_FOUND"
    Message string `json:"message"`          // human-readable
    Field   string `json:"field,omitempty"`  // validation context
    Err     error  `json:"-"`                // wrapped cause
}

func (e *AppError) Error() string { return fmt.Sprintf("[%s] %s", e.Code, e.Message) }
func (e *AppError) Unwrap() error { return e.Err }

// Constructors
func NotFound(resource, id string) *AppError {
    return &AppError{Code: "NOT_FOUND", Message: fmt.Sprintf("%s %q not found", resource, id)}
}

func InvalidInput(field, reason string) *AppError {
    return &AppError{Code: "INVALID_INPUT", Message: reason, Field: field}
}
```

JSON output agents can parse:

```json
{"code": "NOT_FOUND", "message": "user \"abc123\" not found"}
```

## Combining Multiple Errors

Use `errors.Join` (Go 1.20+) to collect validation errors rather than returning
on the first failure. Agents get a complete picture of what is wrong:

```go
var errs []error
if name == "" {
    errs = append(errs, errors.New("name is required"))
}
if email == "" {
    errs = append(errs, errors.New("email is required"))
}
if err := errors.Join(errs...); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

## Concurrent Error Propagation

Use `errgroup` when running independent goroutines — it captures the first error
and cancels remaining work via context:

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return fetchUser(ctx, id) })
g.Go(func() error { return fetchOrders(ctx, id) })
if err := g.Wait(); err != nil {
    return fmt.Errorf("parallel fetch: %w", err)
}
```

Agents detect failure via exit codes. Use them consistently:

| Code | Meaning                        |
| ---- | ------------------------------ |
| `0`  | Success                        |
| `1`  | General / unspecified error    |
| `2`  | Misuse / invalid arguments     |
| `3`  | Not found                      |
| `4`  | Permission denied              |
| `5`  | Timeout / connection failure   |

Pattern for CLI entry point:

```go
func main() {
    if err := run(); err != nil {
        var appErr *AppError
        if errors.As(err, &appErr) {
            fmt.Fprintf(os.Stderr, "error: %s\n", appErr.Message)
            switch appErr.Code {
            case "NOT_FOUND":
                os.Exit(3)
            case "PERMISSION_DENIED":
                os.Exit(4)
            case "TIMEOUT":
                os.Exit(5)
            default:
                os.Exit(1)
            }
        }
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
```

**Never** call `os.Exit` in business logic — only in `main()`. Return errors
up the stack.

**Never** use `panic` outside of `init()` or package-level setup. Panics
produce unstructured stack traces that agents cannot parse.

## Context Propagation

Thread `context.Context` through all I/O-bound functions so agents and
orchestrators can cancel long-running operations:

```go
// ❌ Bad — no cancellation support
func FetchUser(id string) (*User, error) {
    return db.Query("SELECT * FROM users WHERE id = ?", id)
}

// ✅ Good — cancellation propagates
func FetchUser(ctx context.Context, id string) (*User, error) {
    return db.QueryContext(ctx, "SELECT * FROM users WHERE id = ?", id)
}
```

Rules:

- Context is always the **first** parameter.
- Never store context in a struct.
- Never pass `nil` — use `context.Background()` or `context.TODO()` as roots.
- Always check `ctx.Err()` in long loops:

```go
for _, item := range items {
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("cancelled after %d items: %w", processed, err)
    }
    // process item
}
```

## HTTP Error Responses

Return structured JSON errors for all non-2xx responses so agents can parse
failure reasons without screen-scraping:

```go
func writeError(w http.ResponseWriter, status int, err *AppError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(err)
}
```

## Panic Recovery Middleware

Catch unexpected panics at the HTTP boundary so the process keeps serving and
the agent receives a structured error response instead of a connection reset:

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                writeError(w, http.StatusInternalServerError,
                    &AppError{Code: "INTERNAL_ERROR", Message: "an unexpected error occurred"})
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```
