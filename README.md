# yaah — yet another agent harness

A Go toolkit for managing Claude Code configuration across repositories. Define hooks, MCP servers, LSP servers, skills, agents, commands, and plugins in Go, then generate `.claude/settings.json` and all supporting files from code.

## Why

Claude Code settings are scattered across JSON files, markdown skills, agent definitions, and hook scripts. When you work across multiple repos, keeping these in sync is tedious. yaah gives you a single Go program that generates everything.

You write Go code describing your setup. yaah turns it into the `.claude/` directory structure that Claude Code expects.

## Quickstart

### Install

```bash
go install github.com/dirien/yet-another-agent-harness/cmd/yaah@latest
```

### Generate with defaults

The fastest way to get started — generates a full `.claude/` setup with all built-in handlers, skills, agents, LSP servers, and MCP providers:

```bash
cd your-repo
yaah generate --from-code
```

This gives you:

- **Linter** — golangci-lint (Go), ruff (Python), prettier (JS/TS/CSS/MD)
- **CommandGuard** — blocks `rm -rf /`, force-push to main, `DROP TABLE`, etc.
- **SecretScanner** — catches AWS keys, GitHub PATs, API keys, private keys
- **CommentChecker** — flags `TODO: implement`, `// ...`, non-English comments
- **SessionLogger** — audit trail of session start/stop events
- **MCP** — Context7 for library documentation, Pulumi for AI-powered infrastructure
- **LSP** — gopls, pyright, typescript, csharp (via official marketplace `enabledPlugins`)
- **Skills** — commit, PR, review workflows + remote Pulumi skills + 13 development skills
- **Agents** — executor (sonnet), librarian (haiku), reviewer (opus)

### Check your setup

```bash
yaah doctor
```

Reports which binaries are installed and which are missing, with install hints:

```
yaah doctor
===========
Version: dev (commit: none, built: unknown)

Binary:
  ✓ yaah                     /usr/local/bin/yaah
  ✓ git                      /usr/bin/git

Config:
  ✓ .claude/settings.json     valid JSON

LSP Servers:
  ✓ gopls                    /usr/local/bin/gopls
  ✓ pyright                  /usr/local/bin/pyright-langserver
  ✓ typescript               /usr/local/bin/typescript-language-server
  ✗ csharp                   not found → dotnet tool install -g csharp-ls

MCP Servers:
  ✓ context7                 /usr/local/bin/npx

Lint Tools:
  ✓ golangci-lint/format     /usr/local/bin/gofmt
  ✓ golangci-lint/lint       /usr/local/bin/golangci-lint

---
1 issue(s) found. Install missing tools to get full functionality.
```

### Cherry-pick what you need

Don't want everything? Use `DefaultOptions` to pick components:

```go
opts := harness.DefaultOptions{
    EnableCommandGuard:   true,
    EnableSecretScanner:  true,
    LintProfiles:         []handlers.Profile{handlers.GolangCILint()},
    EnableCommitSkill:    true,
    EnableGopls:          true,
    Settings: &schema.Settings{
        Model:       "sonnet",
        EffortLevel: "medium",
    },
}
h := harness.NewWithDefaults(opts)
```

### From JSON config

If you prefer configuration over code:

```bash
yaah init          # creates a starter yaah.json
yaah generate      # generates .claude/ from yaah.json
```

### From Go code (full control)

For maximum control, write a Go program that registers components directly:

```go
package main

import (
    "fmt"
    "os"

    "github.com/dirien/yet-another-agent-harness/pkg/generator"
    "github.com/dirien/yet-another-agent-harness/pkg/harness"
    "github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
    "github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func main() {
    thinking := true
    h := harness.New()

    h.SetSettings(&schema.Settings{
        Model:                 "opus",
        AlwaysThinkingEnabled: &thinking,
        EffortLevel:           "high",
    })

    h.Hooks().Register(handlers.NewLinterWith(
        handlers.GolangCILint(),
    ))
    h.Hooks().Register(handlers.NewCommandGuard())
    h.Hooks().Register(handlers.NewSecretScanner())

    cfg := h.GenerateConfig()
    data, _ := generator.GenerateClaudeSettings(cfg)
    os.MkdirAll(".claude", 0o755)
    os.WriteFile(".claude/settings.json", data, 0o644)
    h.WriteAll(".")

    fmt.Println("Generated .claude/")
}
```

Run it with `go run ./cmd/your-setup/` whenever you change your configuration.

## What it generates

```
.claude/
├── settings.json          # Settings, hooks, MCP servers, enabledPlugins (LSP)
├── skills/
│   ├── commit/SKILL.md    # Skill definitions
│   ├── pr/SKILL.md
│   └── review/SKILL.md
├── agents/
│   ├── executor.md        # Agent definitions with YAML frontmatter
│   ├── librarian.md
│   └── reviewer.md
└── commands/
    └── deploy.md          # Slash command definitions
```

MCP servers are written inline in `settings.json` under the `mcpServers` key. LSP servers are enabled via `enabledPlugins` in `settings.json`, referencing the official Claude Code marketplace plugins (e.g. `gopls-lsp@claude-plugins-official`).

## Architecture

yaah follows an **interface + registry** pattern. Each domain has an interface for individual components and a registry that holds them:

| Domain   | Interface          | Registry            | What it does                                 |
| -------- | ------------------ | ------------------- | -------------------------------------------- |
| Hooks    | `hooks.Handler`    | `hooks.Registry`    | Run code on Claude Code lifecycle events     |
| MCP      | `mcp.Provider`     | `mcp.Registry`      | Configure MCP servers                        |
| LSP      | `lsp.Provider`     | `lsp.Registry`      | Configure LSP servers with binary validation |
| Skills   | `skills.Skill`     | `skills.Registry`   | Generate SKILL.md files                      |
| Agents   | `agents.Agent`     | `agents.Registry`   | Generate agent markdown files                |
| Commands | `commands.Command` | `commands.Registry` | Generate slash command files                 |
| Plugins  | `plugins.Plugin`   | `plugins.Registry`  | Generate plugin packages                     |

The `Harness` struct in `pkg/harness/` wires all registries together and provides `GenerateConfig()` and `WriteAll()` to produce output files. Every concrete implementation has a compile-time interface assertion (`var _ Interface = (*Type)(nil)`).

## Hooks

Hooks run shell commands in response to Claude Code lifecycle events (18 total: `SessionStart`, `PreToolUse`, `PostToolUse`, etc.). yaah ships with five built-in handlers:

### Linter

Runs formatters and linters after every file edit. Profile-based, so you pick which tools to run per file extension.

```go
h.Hooks().Register(handlers.NewLinterWith(
    handlers.GolangCILint(),   // .go → gofmt + golangci-lint
    handlers.Ruff(),           // .py → ruff check + format
    handlers.Prettier(),       // .js/.ts/.css/.md → prettier
    handlers.Biome(),          // .ts/.tsx/.js/.jsx → biome
    handlers.RustFmt(),        // .rs → rustfmt + clippy
    handlers.GoVet(),          // .go → gofmt + go vet (minimal)
))
```

Custom profiles:

```go
handlers.Profile{
    Name:       "mypy",
    Extensions: []string{".py"},
    Steps: []handlers.Step{
        {Label: "typecheck", Cmd: []string{"mypy"}, AppendFile: true, FailBlocks: true},
    },
}
```

### CommandGuard

Blocks dangerous shell commands before they run (`PreToolUse` on Bash). Default patterns catch `rm -rf /`, `git push --force main`, `git reset --hard`, `DROP TABLE`, and similar.

```go
guard := handlers.NewCommandGuard()
guard.Block(`kubectl delete namespace`, "namespace deletion")
h.Hooks().Register(guard)
```

### SecretScanner

Scans edited files for hardcoded credentials: AWS keys, GitHub PATs, OpenAI/Anthropic API keys, private keys, auth tokens, Slack tokens, and more. Blocks the edit if found.

```go
scanner := handlers.NewSecretScanner()
scanner.AddPattern(`my-custom-key-[0-9a-f]{32}`, "custom API key")
h.Hooks().Register(scanner)
```

### CommentChecker

Flags placeholder comments (`TODO: implement`, `// ...`, `// your code here`) and non-English comments in edited files.

```go
checker := handlers.NewCommentChecker()
checker.AddPattern(`(?i)//\s*HACK\s*$`)
h.Hooks().Register(checker)
```

### SessionLogger

Writes session start/stop events with timestamps to a log file.

```go
h.Hooks().Register(handlers.NewSessionLogger("/var/log/claude-sessions"))
```

### Writing your own handler

Implement `hooks.Handler`:

```go
type Handler interface {
    Name() string
    Events() []schema.HookEvent
    Match() *regexp.Regexp
    Execute(ctx context.Context, input *hooks.Input) (*hooks.Result, error)
}
```

Register it: `h.Hooks().Register(myHandler)`. The generated `settings.json` will wire it to run via `yaah hook <event>`.

### HookHandler fields

`HookHandler` supports the following fields. All optional fields default to zero/false if omitted.

| Field            | Type              | Description                                                         |
| ---------------- | ----------------- | ------------------------------------------------------------------- |
| `type`           | string            | `command`, `http`, `prompt`, or `agent`                             |
| `command`        | string            | Shell command to run (`type=command`)                               |
| `url`            | string            | HTTP endpoint to POST to (`type=http`)                              |
| `prompt`         | string            | LLM prompt to evaluate (`type=prompt` or `agent`)                   |
| `timeout`        | int               | Timeout in **seconds** (defaults: command=600, prompt=30, agent=60) |
| `statusMessage`  | string            | Custom spinner message displayed while the hook runs                |
| `once`           | bool              | Run only once per session (skills only)                             |
| `async`          | bool              | Run in background without blocking (`type=command` only)            |
| `headers`        | map[string]string | HTTP headers to include (`type=http` only)                          |
| `allowedEnvVars` | []string          | Env var names allowed in header interpolation (`type=http` only)    |
| `model`          | string            | Model identifier for LLM evaluation (`type=prompt` or `agent`)      |

## MCP servers

```go
import "github.com/dirien/yet-another-agent-harness/pkg/mcp/providers"

h.MCP().Register(providers.NewContext7())
h.MCP().Register(providers.NewPulumi())
h.MCP().Register(providers.NewNotion("ntn-your-api-token"))
h.MCP().Register(providers.NewCustom(schema.MCPServer{
    Name:      "my-server",
    Transport: schema.MCPTransportStdio,
    Command:   "my-mcp-binary",
    Args:      []string{"--mode", "production"},
}))
```

### Transport types

`MCPTransport` supports five values: `stdio`, `sse`, `streamable-http`, `http`, and `websocket`.

### OAuth support

Remote MCP servers that require OAuth authentication use the `MCPOAuth` struct:

```go
schema.MCPServer{
    Name:      "my-remote-server",
    Transport: schema.MCPTransportSSE,
    URL:       "https://api.example.com/mcp",
    OAuth: &schema.MCPOAuth{
        ClientID:              "my-client-id",
        CallbackPort:          8085,
        AuthServerMetadataURL: "https://auth.example.com/.well-known/oauth-authorization-server",
    },
}
```

## LSP servers

yaah enables LSP servers via `enabledPlugins` in the generated `settings.json`, referencing the official Claude Code marketplace. Providers that implement `MarketplaceProvider` are added automatically. yaah also validates that server binaries exist in `$PATH` via `yaah doctor`.

Default LSP servers (all from `claude-plugins-official`):

| Server     | Marketplace key                          | Install                                                |
| ---------- | ---------------------------------------- | ------------------------------------------------------ |
| gopls      | `gopls-lsp@claude-plugins-official`      | `go install golang.org/x/tools/gopls@latest`           |
| pyright    | `pyright-lsp@claude-plugins-official`    | `pip install pyright`                                  |
| typescript | `typescript-lsp@claude-plugins-official` | `npm install -g typescript-language-server typescript` |
| csharp     | `csharp-lsp@claude-plugins-official`     | `dotnet tool install -g csharp-ls`                     |

```go
import lspproviders "github.com/dirien/yet-another-agent-harness/pkg/lsp/providers"

h.LSP().Register(lspproviders.NewGopls())
h.LSP().Register(lspproviders.NewPyright())
h.LSP().Register(lspproviders.NewTypeScript())
h.LSP().Register(lspproviders.NewCSharp())
```

Additional LSP servers (yaml-ls, pulumi, pulumi-yaml) are available as providers but not in the official marketplace, so they are not included in defaults.

## Skills

Skills generate `SKILL.md` files under `.claude/skills/`. Built-in skills: `commit` (git commit workflow), `pr` (pull request creation), `review` (code review checklist).

```go
import "github.com/dirien/yet-another-agent-harness/pkg/skills/builtins"

h.Skills().Register(builtins.NewCommitSkill())
h.Skills().Register(builtins.NewPRSkill())
h.Skills().Register(builtins.NewReviewSkill())
```

### Remote skills

Pin a skill to a remote git repo, GitHub Actions style:

```go
import "github.com/dirien/yet-another-agent-harness/pkg/skills"

h.Skills().Register(skills.NewRemoteSkill(
    "team-standards",
    "Shared coding standards",
    "github.com/myorg/claude-skills@v1.2.0",
    "skills/standards/SKILL.md",
))
```

The ref can be a tag, branch, or commit SHA. Skills are cached in `~/.yaah/cache/skills/` and invalidated when the ref changes. Override the cache directory with `YAAH_HOME`. Multi-file skills (with reference directories) implement `skills.SkillWithFiles`.

### Default remote skills

yaah ships with remote skills from three repositories:

**pulumi/agent-skills** — Pulumi IaC authoring and migration:

| Skill                        | Description                              |
| ---------------------------- | ---------------------------------------- |
| `pulumi-best-practices`      | Reliable program patterns                |
| `pulumi-component`           | ComponentResource authoring              |
| `pulumi-automation-api`      | Automation API best practices            |
| `pulumi-esc`                 | Environments, secrets, and configuration |
| `pulumi-terraform-to-pulumi` | Convert Terraform to Pulumi              |
| `pulumi-cdk-to-pulumi`       | Convert AWS CDK to Pulumi                |
| `cloudformation-to-pulumi`   | Convert CloudFormation to Pulumi         |
| `pulumi-arm-to-pulumi`       | Convert Azure ARM/Bicep to Pulumi        |

**dirien/claude-skills** — Pulumi language-specific skills:

| Skill               | Description                              |
| ------------------- | ---------------------------------------- |
| `pulumi-typescript` | Pulumi TypeScript IaC with ESC and OIDC  |
| `pulumi-go`         | Pulumi Go IaC with ESC and OIDC          |
| `pulumi-python`     | Pulumi Python IaC with ESC and OIDC      |
| `pulumi-neo`        | Pulumi Neo conversational infrastructure |
| `pulumi-cli`        | Pulumi CLI command reference             |

**jeffallan/claude-skills** — Development and operations skills:

| Skill                   | Description                                                               |
| ----------------------- | ------------------------------------------------------------------------- |
| `golang-pro`            | Go concurrent patterns, microservices, gRPC, and performance optimization |
| `kubernetes-specialist` | Kubernetes deployments, Helm, RBAC, NetworkPolicies, and multi-cluster    |
| `devops-engineer`       | CI/CD pipelines, Docker, Kubernetes, Terraform, and GitOps                |
| `python-pro`            | Python 3.11+ with type safety, async, pytest, and ruff                    |
| `typescript-pro`        | Advanced TypeScript types, generics, tRPC, and monorepo setup             |
| `csharp-developer`      | C# .NET 8+, ASP.NET Core, Blazor, EF Core, and MediatR                    |
| `javascript-pro`        | Modern ES2023+ JavaScript, async/await, ESM, and Node.js                  |
| `cli-developer`         | CLI tools with argument parsing, completions, and cross-platform support  |
| `sre-engineer`          | SLOs, error budgets, incident response, and capacity planning             |
| `the-fool`              | Devil's advocate, pre-mortems, red teaming, and assumption auditing       |
| `architecture-designer` | System architecture, ADRs, trade-offs, and scalability planning           |
| `spring-boot-engineer`  | Spring Boot 3.x, Spring Security 6, JPA, WebFlux, and Spring Cloud        |
| `code-reviewer`         | Code review for bugs, security, performance, and maintainability          |

### Skill frontmatter

Skills can optionally implement the `SkillWithFrontmatter` interface to inject YAML frontmatter into the generated `SKILL.md`:

```go
type SkillWithFrontmatter interface {
    Skill
    Frontmatter() SkillFrontmatter
}
```

Available frontmatter fields in `SkillFrontmatter`:

| Field                      | Type   | Description                                 |
| -------------------------- | ------ | ------------------------------------------- |
| `argument-hint`            | string | Usage hint shown during autocomplete        |
| `disable-model-invocation` | bool   | Prevent Claude from auto-loading this skill |
| `user-invocable`           | \*bool | Show in /menu (default true)                |
| `allowed-tools`            | string | Comma-separated tool allowlist              |
| `model`                    | string | Model override when skill is active         |
| `context`                  | string | Set to `fork` for subagent execution        |
| `agent`                    | string | Subagent type when `context=fork`           |

### Writing your own skill

Implement `skills.Skill`:

```go
type Skill interface {
    Name() string
    Description() string
    Content() string
    Source() schema.SkillSource
}
```

## Agents

Agents generate markdown files with YAML frontmatter under `.claude/agents/`. Three built-in agents:

- **executor** (sonnet) — single-task agent with strict verification
- **librarian** (haiku) — research-only, read-only tools
- **reviewer** (opus) — code review with read-only tools

```go
import agentpkg "github.com/dirien/yet-another-agent-harness/pkg/agents"

h.Agents().Register(agentpkg.NewExecutor())
h.Agents().Register(agentpkg.NewLibrarian())
h.Agents().Register(agentpkg.NewReviewer())
```

### Advanced agent frontmatter

Agents can optionally implement `AgentWithAdvanced` to set additional frontmatter fields:

```go
type AgentWithAdvanced interface {
    Agent
    Advanced() AgentAdvanced
}
```

Available fields in `AgentAdvanced`:

| Field             | Type           | Description                                                         |
| ----------------- | -------------- | ------------------------------------------------------------------- |
| `disallowedTools` | string         | Comma-separated tool denylist                                       |
| `permissionMode`  | string         | `default`, `acceptEdits`, `dontAsk`, `bypassPermissions`, or `plan` |
| `maxTurns`        | int            | Maximum agentic turns before stopping                               |
| `skills`          | []string       | Skills to preload into agent context                                |
| `mcpServers`      | map[string]any | MCP servers for this agent                                          |
| `hooks`           | HooksConfig    | Lifecycle hooks (PreToolUse/PostToolUse/Stop)                       |
| `memory`          | string         | Persistent memory scope: `user`, `project`, or `local`              |
| `background`      | bool           | Run as background task                                              |
| `isolation`       | string         | Isolation mode: `worktree`                                          |

## Commands

Slash commands generate markdown files under `.claude/commands/`. Implement `commands.Command`:

```go
type Command interface {
    Name() string
    Description() string
    ArgumentHint() string
    Model() string
    AllowedTools() string
    Content() string
}
```

### Advanced command frontmatter

Commands can optionally implement `CommandWithAdvanced` to set additional frontmatter fields:

```go
type CommandWithAdvanced interface {
    Command
    Advanced() CommandAdvanced
}
```

Available fields in `CommandAdvanced`:

| Field                      | Type   | Description                                   |
| -------------------------- | ------ | --------------------------------------------- |
| `disable-model-invocation` | bool   | Prevent Claude from auto-loading this command |
| `user-invocable`           | \*bool | Show in /menu (default true)                  |
| `context`                  | string | Set to `fork` for subagent execution          |
| `agent`                    | string | Subagent type when `context=fork`             |

## Settings

`schema.Settings` maps directly to the official Claude Code `settings.json` specification. Key field groups:

| Group                  | Fields                                                                                                                                                                       |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Core                   | `model`, `alwaysThinkingEnabled`, `autoUpdatesChannel`, `effortLevel`, `statusLine`, `env`, `teammateMode`                                                                   |
| Model & Performance    | `availableModels`, `fastMode`, `fastModePerSessionOptIn`                                                                                                                     |
| Permissions & Security | `permissions`, `sandbox`, `allowManagedPermissionRulesOnly`                                                                                                                  |
| Hooks & Automation     | `disableAllHooks`, `allowManagedHooksOnly`                                                                                                                                   |
| Git & Attribution      | `attribution`, `includeGitInstructions`                                                                                                                                      |
| Authentication         | `apiKeyHelper`                                                                                                                                                               |
| UI & Behavior          | `language`, `outputStyle`, `showTurnDuration`, `spinnerVerbs`, `spinnerTipsEnabled`, `spinnerTipsOverride`, `prefersReducedMotion`, `terminalProgressBarEnabled`             |
| Plugins                | `enabledPlugins`, `pluginConfigs`, `extraKnownMarketplaces`, `strictKnownMarketplaces`, `skippedMarketplaces`, `skippedPlugins`, `blockedMarketplaces`, `pluginTrustMessage` |
| MCP Management         | `enableAllProjectMcpServers`, `enabledMcpjsonServers`, `disabledMcpjsonServers`, `allowedMcpServers`, `deniedMcpServers`, `allowManagedMcpServersOnly`                       |
| Organization           | `companyAnnouncements`, `cleanupPeriodDays`, `plansDirectory`, `autoMemoryEnabled`, `skipWebFetchPreflight`                                                                  |
| File & Directory       | `fileSuggestion`, `respectGitignore`, `additionalDirectories`                                                                                                                |

## Plugins

The `schema.Plugin` struct matches the official Claude Code `plugin.json` specification. `Commands` is `[]string` (paths to command markdown files). Key fields:

| Field          | Type        | Description                                      |
| -------------- | ----------- | ------------------------------------------------ |
| `name`         | string      | Plugin identifier                                |
| `version`      | string      | Semver version                                   |
| `commands`     | []string    | Paths to command markdown files                  |
| `agents`       | []string    | Paths to agent markdown files or agent directory |
| `skills`       | []string    | Paths to skill directories                       |
| `hooks`        | HooksConfig | Hook registrations                               |
| `mcpServers`   | string      | Path to `.mcp.json` or inline MCP config         |
| `lspServers`   | string      | Path to `.lsp.json` or inline LSP config         |
| `outputStyles` | string      | Path to output styles directory                  |

## CLI reference

```
yaah generate                  # Generate .claude/ from yaah.json
yaah generate -c path/to.json  # Generate from a specific config file
yaah generate --from-code      # Generate from Go-registered components
yaah generate -o ./out         # Output to a different directory
yaah init                      # Create a starter yaah.json
yaah schema                    # Print the JSON Schema for yaah.json
yaah schema -o schema.json     # Write schema to file
yaah hook <event>              # Runtime hook dispatcher (called by Claude Code)
yaah info                      # Show all registered components
yaah doctor                    # Health check — validates binaries and config
yaah version                   # Print version, commit, and build date
```

## Project structure

```
pkg/schema/            Data types, one file per concern
pkg/hooks/             Handler interface + Registry
pkg/hooks/handlers/    Linter, CommandGuard, SecretScanner, CommentChecker, SessionLogger
pkg/mcp/               MCP Provider interface + Registry
pkg/mcp/providers/     Context7, Pulumi, Notion, Custom
pkg/lsp/               LSP Provider + MarketplaceProvider interfaces + Registry + binary validation
pkg/lsp/providers/     Gopls, Pyright, TypeScript, CSharp (marketplace) + YamlLS, PulumiLSP, PulumiYAML, Custom
pkg/skills/            Skill interface + Registry + RemoteSkill + SkillWithFrontmatter
pkg/skills/builtins/   CommitSkill, PRSkill, ReviewSkill
pkg/agents/            Agent interface + Registry + AgentWithAdvanced + Executor, Librarian, Reviewer
pkg/commands/          Command interface + Registry + CommandWithAdvanced
pkg/plugins/           Plugin interface + Registry
pkg/harness/           Harness (top-level wiring) + defaults
pkg/generator/         JSON Schema + settings.json + .lsp.json generation
cmd/yaah/              CLI entry point
internal/cli/          Config file discovery and loading
```

## Defaults

`harness.NewWithDefaults(harness.AllDefaults())` gives you everything out of the box:

- Linter with GolangCILint, Ruff, and Prettier profiles
- CommandGuard blocking dangerous shell commands
- CommentChecker for placeholder comments
- SecretScanner for hardcoded credentials
- SessionLogger for audit trails
- Context7 and Pulumi MCP servers
- 4 LSP servers via official marketplace (gopls, pyright, typescript, csharp)
- Commit, PR, and Review skills
- 13 remote Pulumi skills (best practices, components, ESC, migrations)
- 13 remote development skills from jeffallan/claude-skills (Go, Python, TypeScript, K8s, DevOps, SRE, and more)
- Executor, Librarian, and Reviewer agents

## License

[MIT](LICENSE)
