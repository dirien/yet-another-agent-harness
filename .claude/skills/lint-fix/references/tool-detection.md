<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Tool detection

Detect project linting, formatting, and validation tools by scanning
for configuration files, Makefile targets, and package manager scripts.
Use the first match found in each ecosystem.

## Makefile detection

If `Makefile` exists, discover relevant targets:

```bash
make -qp 2>/dev/null | grep -E '^[a-z].*:' | cut -d: -f1 | sort -u
```

Look for these common targets: `lint`, `test`, `validate`, `format`,
`fmt`, `gen`, `generate`, `build`, `check`, `vet`.

## Node.js / TypeScript

If `package.json` exists:

```bash
jq -r '.scripts | keys[]' package.json 2>/dev/null
```

Common scripts: `lint`, `test`, `format`, `prettier`, `eslint`,
`build`, `typecheck`, `validate`.

Standalone configs (check existence):

- `.eslintrc*`, `eslint.config.*` → `npx eslint .`
- `.prettierrc*`, `prettier.config.*` → `npx prettier --check .`
- `biome.json`, `biome.jsonc` → `npx biome check .`
- `tsconfig.json` → `npx tsc --noEmit`
- `.stylelintrc*` → `npx stylelint "**/*.css"`

## Go

If `go.mod` exists:

- `golangci-lint run` (if `golangci-lint` installed or
  `.golangci.yml` exists)
- `gofmt -s -l .` (always available with Go)
- `go vet ./...` (always available with Go)
- `go test ./...` (always available with Go)
- `go build ./...` (always available with Go)

## Python

If `pyproject.toml`, `setup.py`, or `requirements.txt` exists:

- `ruff check .` (if `ruff` installed or `[tool.ruff]` in
  pyproject.toml)
- `black --check .` (if `black` installed or `[tool.black]` in
  pyproject.toml)
- `flake8` (if `flake8` installed or `.flake8` exists)
- `mypy .` (if `mypy` installed or `[tool.mypy]` in pyproject.toml)
- `pytest` (if `pytest` installed or `[tool.pytest]` in
  pyproject.toml)

## .NET

If `*.csproj` or `*.sln` exists:

- `dotnet format --verify-no-changes`
- `dotnet build`
- `dotnet test`

## Rust

If `Cargo.toml` exists:

- `cargo clippy -- -D warnings`
- `cargo fmt --check`
- `cargo test`
- `cargo build`

## Markdown (always check)

- `markdownlint` (if installed or `.markdownlint*` config exists)
- `commitlint` (if installed or `.commitlintrc*` config exists)

## Pipeline ordering

Once tools are detected, order them:

1. **Generation** — `make gen`, code generators (must run first,
   output may overwrite manual fixes)
2. **Formatting** — prettier, black, gofmt, rustfmt, dotnet format
   (fix formatting before linting)
3. **Linting** — eslint, golangci-lint, ruff, clippy, markdownlint
4. **Type checking** — tsc, mypy (after lint fixes)
5. **Validation** — `make validate`, schema checks
6. **Testing** — go test, pytest, jest, cargo test, dotnet test
7. **Building** — go build, cargo build, dotnet build (final check)

## Fallback

If no tools are detected, report to the user:

```text
No linting or validation tools detected in this project.

Checked: Makefile, package.json, go.mod, pyproject.toml,
*.csproj, Cargo.toml, standalone linter configs.

Would you like to set up linting for this project?
```
