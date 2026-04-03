<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Upgrade commands per ecosystem

Per-language commands for listing outdated packages, upgrading,
verifying lockfiles, rolling back, and scanning for vulnerabilities.

## Go

```bash
# List outdated
go list -u -m all 2>/dev/null | grep '\['

# Upgrade specific package
go get package/name@latest          # latest version
go get package/name@v1.2.3          # specific version
go get package/name@v1              # latest v1.x

# Tidy and verify
go mod tidy
go mod verify

# Test
go test ./...
go build ./...

# Rollback
git checkout -- go.mod go.sum

# Vulnerability scan
govulncheck ./...
```

## Node.js

```bash
# List outdated
npm outdated

# Upgrade (respects semver in package.json)
npm update package-name

# Upgrade to latest (ignores semver range)
npm install package-name@latest

# Check available major updates
npx npm-check-updates

# Apply major updates to package.json
npx npm-check-updates -u
npm install

# Test
npm test
npm run build

# Rollback
git checkout -- package.json package-lock.json
npm install

# Vulnerability scan
npm audit
```

## Python

```bash
# List outdated
pip list --outdated

# Upgrade specific package
pip install --upgrade package-name

# Update requirements.txt after upgrade
pip freeze > requirements.txt

# Poetry projects
poetry show --outdated
poetry update package-name

# Pipenv projects
pipenv update package-name

# Test
pytest
python -m pytest

# Rollback (choose the appropriate lockfile)
git checkout -- requirements.txt
# or: git checkout -- poetry.lock
# or: git checkout -- Pipfile.lock

# Vulnerability scan
pip-audit
```

## Rust

```bash
# List outdated (requires cargo-outdated)
cargo outdated

# Update specific package (within Cargo.toml constraints)
cargo update -p package-name

# For major version changes, edit Cargo.toml first, then:
cargo update

# Test
cargo test
cargo build

# Rollback
git checkout -- Cargo.lock

# Vulnerability scan (requires cargo-audit)
cargo audit
```

## .NET

```bash
# List outdated
dotnet list package --outdated

# Upgrade specific package
dotnet add package PackageName --version X.Y.Z

# Upgrade to latest
dotnet add package PackageName

# Test
dotnet test
dotnet build

# Rollback
git checkout -- *.csproj

# Vulnerability scan
dotnet list package --vulnerable
```

## Lockfile verification

After upgrades, verify lockfile integrity:

| Ecosystem | Lockfile | Verify command |
| --------- | -------- | -------------- |
| Go | `go.sum` | `go mod verify` |
| Node | `package-lock.json` | `npm ci` (clean install from lockfile) |
| Python | `requirements.txt` | `pip install -r requirements.txt --dry-run` |
| Rust | `Cargo.lock` | `cargo build` (validates lockfile) |
| .NET | N/A (no lockfile) | `dotnet restore` |
