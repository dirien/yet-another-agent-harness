<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Testing for Agent-Ready Go

## Core Test Command

```bash
go test -race -count=1 -timeout=5m -coverprofile=coverage.out ./...
```

Always include `-race`. It catches data races that only appear under concurrent
load.

## go vet

Run before linting — it catches common mistakes the compiler misses:

```bash
go vet ./...
```

Common issues: unreachable code, incorrect format strings, mutex copied by value,
misuse of `sync/atomic`.

## Coverage

```bash
# Function-level summary
go tool cover -func=coverage.out

# HTML report for human inspection
go tool cover -html=coverage.out -o coverage.html

# Enforce 80% minimum threshold
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | tr -d '%')
[ $(echo "$COVERAGE >= 80" | bc -l) -eq 1 ] || { echo "Coverage ${COVERAGE}% < 80%"; exit 1; }
```

The `make cover` target in `assets/Makefile` runs all of the above.

## golangci-lint

Install:

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin
```

Run:

```bash
golangci-lint run ./...
```

Copy `assets/.golangci.yml` from this skill to the project root for a thorough
production-ready configuration.

### Key linters in the config and why they matter for agents

| Linter         | Why it matters for agents                                      |
| -------------- | -------------------------------------------------------------- |
| `errcheck`     | Unchecked errors are invisible to agents                       |
| `wrapcheck`    | Errors must be wrapped to trace failures                       |
| `errorlint`    | Prevents incorrect error comparison (breaks `errors.Is`)       |
| `govet`        | Catches suspicious constructs the compiler allows              |
| `staticcheck`  | Comprehensive static analysis — catches subtle bugs            |
| `gosec`        | Security issues agents may inadvertently introduce             |
| `contextcheck` | Ensures `context.Context` is threaded correctly                |
| `noctx`        | Ensures HTTP requests carry context                            |
| `exhaustive`   | Switch statements cover all enum values                        |
| `testifylint`  | Correct use of testify assertions                              |
| `revive`       | Comprehensive Go style rules                                   |

### Suppressing false positives

Per-line:

```go
defer f.Close() //nolint:errcheck // intentionally ignoring cleanup error
```

Per-file (top of file):

```go
//nolint:all
```

Per-rule in `.golangci.yml` under `issues.exclude-rules`:

```yaml
issues:
  exclude-rules:
    - path: internal/legacy
      linters:
        - wrapcheck
```

## Benchmarks

```go
func BenchmarkProcessOrder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        processOrder(testOrder)
    }
}
// Go 1.24+ alternative: for b.Loop() { ... }
```

Run: `go test -bench=. -benchmem ./...`

## Table-Driven Tests (Standard Pattern)

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
        {"negative", -1, -2, -3},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

## HTTP Handler Testing

Use `net/http/httptest` to test HTTP handlers without a live server. The
recorder captures status codes and body for assertions:

```go
func TestCreateUserHandler(t *testing.T) {
    svc := &mockUserService{}
    h := NewUserHandler(svc)

    body := `{"name":"alice","email":"alice@example.com"}`
    req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    h.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    var resp UserResponse
    require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
    assert.Equal(t, "alice", resp.Name)
}
```

## Integration Tests

Use build tags to separate integration tests from unit tests. Integration tests
run against real databases or external services and are excluded from `go test ./...`
by default:

```go
//go:build integration

package db_test
```

Run integration tests explicitly:

```bash
go test -tags=integration ./...
```

## CI Test Flags

- `-count=1` — disables the test result cache; ensures tests always run fresh
- `-timeout` — prevents hung tests from blocking CI indefinitely
- `-race` — enables the data race detector

```bash
go test -race -count=1 -timeout=5m -coverprofile=coverage.out ./...
```

## Test Helpers

```go
func TestSomething(t *testing.T) {
    t.Parallel() // allows concurrent execution with other parallel tests

    res := setupResource(t)
    t.Cleanup(func() { res.Close() }) // runs even if the test fails

    assertResult(t, res)
}

func assertResult(t *testing.T, res *Resource) {
    t.Helper() // failure output points to the call site, not this line
    if res == nil {
        t.Fatal("expected non-nil resource")
    }
}
```

- `t.Parallel()` — allows concurrent test execution, dramatically reduces CI time
- `t.Cleanup()` — registered cleanup runs even on failure (prefer over `defer`)
- `t.Helper()` — makes failure output point to the call site, not the helper body
  (critical for agents reading test output to locate the actual problem)

## Mock Generation

Use `mockery` to generate mocks from interfaces automatically. Add a
`go:generate` directive next to the interface definition:

```go
//go:generate mockery --name=OrderStore --output=mocks --outpkg=mocks
type OrderStore interface {
    Get(ctx context.Context, id string) (*Order, error)
    Save(ctx context.Context, o *Order) error
}
```

Install and run:

```bash
go install github.com/vektra/mockery/v2@latest
go generate ./...
```

Define interfaces for external dependencies so agents can inject test doubles:

```go
type OrderStore interface {
    Get(ctx context.Context, id string) (*Order, error)
    Save(ctx context.Context, o *Order) error
}
```

Agents can generate mock implementations automatically when interfaces are
well-defined and narrow.

## Testify (Optional)

```bash
go get github.com/stretchr/testify
```

```go
assert.Equal(t, want, got)
require.NoError(t, err) // stops test immediately on failure
assert.ErrorIs(t, err, ErrNotFound)
```
