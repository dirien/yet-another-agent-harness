<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Analysis patterns

What to investigate in each category that tech-debt analyzes directly.
For delegated categories (code quality, linting, security, dependencies),
see the referenced skills.

## Architecture analysis

Investigate the codebase for structural issues:

**Circular dependencies**:

- Check import/require patterns for cycles
- Go: `go list -f '{{.ImportPath}} {{.Imports}}' ./...` and trace
- Node: check for A imports B imports A patterns
- A cycle between 2 packages is worse than a cycle through 5

**God files**:

- Find files > 500 lines of code (excluding tests)
- Check if they mix multiple responsibilities
- A 600-line file with one clear purpose is fine; a 300-line file
  with 5 unrelated functions is debt

**Missing abstractions**:

- Concrete types passed through 3+ layers without an interface
- Direct dependency on implementation details
- Swapping an implementation requires changing multiple call sites

**Layer violations**:

- HTTP/API handlers containing business logic or database queries
- Database queries outside the data layer
- Presentation logic in domain models

**Coupling indicators**:

- Files importing > 10 packages (high fan-out)
- Functions with > 7 parameters (data clump)
- Changes to one file requiring changes in 5+ others (shotgun
  surgery)

## Test analysis

**Coverage metrics**:

```bash
# Go
go test -coverprofile=coverage.out ./... 2>/dev/null
go tool cover -func=coverage.out | grep total

# Python
pytest --cov --cov-report=term-missing 2>/dev/null

# Node
npx jest --coverage 2>/dev/null
```

**Test/code ratio**:

- Count test files vs code files
- Healthy ratio: > 0.5 (1 test file per 2 code files)
- Alarm: < 0.2

**Critical path coverage**:

- Identify entry points (main, handlers, API routes)
- Trace the critical path through the code
- Check if each step in the critical path has test coverage

**Test anti-patterns**:

- Tests with no assertions (they always pass)
- Tests using `sleep` for synchronization (flaky)
- Tests sharing mutable state (order-dependent)
- Tests testing implementation details (brittle)
- Test files with no `_test` suffix or `test_` prefix (hidden)

## Documentation analysis

**Essential files check**:

| File | Purpose | Weight |
| ---- | ------- | ------ |
| README.md | Project overview, setup, usage | Critical |
| CONTRIBUTING.md | How to contribute | Important |
| CHANGELOG.md | Version history | Important |
| Architecture docs | System design | Valuable |
| API docs | Interface contracts | Valuable (if API exists) |
| ADRs | Decision records | Valuable |

**API documentation coverage**:

- List all exported/public functions, types, interfaces
- Check which have documentation comments
- Ratio of documented exports / total exports

**Stale documentation**:

- References to files or functions that no longer exist
- Setup instructions that don't match current dependencies
- API examples using deprecated endpoints or parameters

## Infrastructure analysis

**CI/CD assessment**:

| Check | Where to look |
| ----- | ------------- |
| CI config exists | `.github/workflows/`, `.gitlab-ci.yml`, `Jenkinsfile`, `.circleci/` |
| Tests run in CI | Grep CI config for `test` commands |
| Linting in CI | Grep CI config for `lint` commands |
| Auto-deploy | Grep for deploy/release steps |
| Branch protection | `gh api repos/{owner}/{repo}/branches/main/protection` |

**Deployment readiness**:

- Dockerfile or container config exists
- Environment configuration via env vars (not hardcoded)
- Health check endpoint (for services)
- Deployment documentation or scripts

**Observability**:

- Structured logging configured (not just `fmt.Println`)
- Error tracking integration (Sentry, Bugsnag, etc.)
- Metrics endpoint or monitoring config
- Alert configuration
