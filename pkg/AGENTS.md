<!-- Managed by agent: keep sections and order; edit content, not structure. Last updated: 2026-03-23 -->

# AGENTS.md — pkg/

<!-- AGENTS-GENERATED:START overview -->
## Overview
Core Go library for yaah. Every component follows the **Interface → Registry → Generator** pattern. Registries collect typed components; generators consume them to emit per-agent config files. The `harness` package is the top-level orchestrator.
<!-- AGENTS-GENERATED:END overview -->

<!-- AGENTS-GENERATED:START filemap -->
## Key Files
```
harness/harness.go      -> Orchestrator: wires all registries, dispatches hook events
hooks/registry.go       -> Hook registry (event → handler dispatch)
hooks/handler.go        -> Handler interface contract
hooks/handlers/*.go     -> Built-in handlers: linter, command-guard, secret-scanner, comment-checker, session-logger
mcp/registry.go         -> MCP server registry
mcp/provider.go         -> Provider interface contract
mcp/providers/*.go      -> Built-in: yaah (self), context7, pulumi
lsp/registry.go         -> LSP server registry
lsp/provider.go         -> Provider interface contract
lsp/providers/*.go      -> Built-in: gopls, pyright, typescript, csharp
skills/registry.go      -> Skill registry + remote loader
skills/builtins/*.go    -> Built-in skills: commit, pr, review
agents/registry.go      -> Agent registry + remote loader
commands/registry.go    -> CLI command registry
plugins/registry.go     -> Plugin registry
generator/claude.go     -> Claude Code config generator
generator/opencode.go   -> OpenCode config generator
generator/codex.go      -> Codex CLI config generator
generator/copilot.go    -> GitHub Copilot CLI config generator
schema/settings.go      -> Settings types + defaults
session/store.go        -> Session audit log persistence
gitcache/cache.go       -> Git clone/update cache for remote resources
mcpserver/server.go     -> Built-in MCP server (stdio transport, tool handlers)
```
<!-- AGENTS-GENERATED:END filemap -->

<!-- AGENTS-GENERATED:START golden-samples -->
## Golden Samples (follow these patterns)
| For | Reference | Key patterns |
|-----|-----------|--------------|
| Hook handler | `hooks/handlers/linter.go` | Implement `Handler` interface, `CanHandle()` + `Handle()` methods |
| MCP provider | `mcp/providers/yaah.go` | Implement `Provider` interface, return `MCPServer` config |
| LSP provider | `lsp/providers/gopls.go` | Implement `Provider` interface, return `LSPServer` config |
| Config generator | `generator/claude.go` | Implement `Generator`, read from registries, write native format |
| Built-in skill | `skills/builtins/commit.go` | `Skill` struct with Name, Description, Source |
| Registry | `hooks/registry.go` | Thread-safe `Register()` + `All()` + typed iteration |
<!-- AGENTS-GENERATED:END golden-samples -->

<!-- AGENTS-GENERATED:START setup -->
## Setup & environment
- `go mod download` to install dependencies
- Go 1.25+ required (see `go.mod`)
- No external Go tools required beyond standard toolchain
<!-- AGENTS-GENERATED:END setup -->

<!-- AGENTS-GENERATED:START commands -->
## Build & tests
| Task | Command |
|------|---------|
| Vet | `go vet ./...` |
| Format | `gofmt -w .` |
| Test | `go test ./...` |
| Test (race) | `go test -race ./...` |
| Test (single pkg) | `go test -v -race ./pkg/hooks/...` |
| Build | `go build -v ./...` |
<!-- AGENTS-GENERATED:END commands -->

<!-- AGENTS-GENERATED:START code-style -->
## Code style & conventions
- Follow Go 1.25 idioms
- Use standard library over external deps when possible
- Errors: wrap with `fmt.Errorf("context: %w", err)`, lowercase no punctuation
- Naming: `camelCase` for private, `PascalCase` for exported; ID/URL/HTTP not Id/Url/Http
- Struct tags: use canonical form (json, yaml, etc.)
- Comments: complete sentences ending with period
- Package docs: first sentence summarizes purpose
- Prefer `any` over `interface{}`; use generics `[T any]` where appropriate
<!-- AGENTS-GENERATED:END code-style -->

<!-- AGENTS-GENERATED:START security -->
## Security & safety
- Validate all inputs from external sources
- Use `context.Context` for cancellation and timeouts
- Avoid goroutine leaks: always ensure termination paths
- Sensitive data: never log or include in errors
- File paths: validate and sanitize user-provided paths (especially in gitcache)
- Hook handlers are security controls (command-guard, secret-scanner) — treat changes with care
<!-- AGENTS-GENERATED:END security -->

<!-- AGENTS-GENERATED:START examples -->
## Patterns to Follow
> **Prefer looking at real code in this repo over generic examples.**
> See **Golden Samples** section above for files that demonstrate correct patterns.

Key patterns:
- **Registry pattern**: Every component type has a registry with `Register()` and `All()` methods
- **Interface contracts**: Define interfaces where consumed (`handler.go`), implement in separate files
- **Context handling**: Always pass and respect `context.Context` in I/O paths
- **Minimal deps**: Standard library first — check go.mod before adding anything
<!-- AGENTS-GENERATED:END examples -->

<!-- AGENTS-GENERATED:START help -->
## When stuck
- Check Go documentation: https://pkg.go.dev
- Review existing patterns in this codebase — every registry follows the same shape
- Check root AGENTS.md for project-wide conventions
- Run `go doc <package>` for standard library help
<!-- AGENTS-GENERATED:END help -->

<!-- AGENTS-GENERATED:START checklist -->
## PR/commit checklist
- [ ] `go test -race ./...` passes
- [ ] `go vet ./...` reports no issues
- [ ] `gofmt -w .` produces no diff
- [ ] Error messages are descriptive and wrapped with `%w`
- [ ] Public APIs have godoc comments
- [ ] `context.Context` passed and respected in all I/O paths
- [ ] New registries follow the `NewRegistry()` → `Register()` → `All()` shape
<!-- AGENTS-GENERATED:END checklist -->

## House Rules (project-specific)
- New registries must follow the same `NewRegistry()` → `Register()` → `All()` shape as existing ones
- Generators must not import other generators — they are independent output targets
- Remote resources (skills, agents) are loaded via `gitcache` — never vendor them into the repo
