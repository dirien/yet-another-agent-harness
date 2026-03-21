# Components

yaah organizes everything into components. Each one follows the same pattern: implement an interface, register it, let yaah generate the output.

## Hooks

Hooks run shell commands when Claude Code fires lifecycle events (18 in total: `SessionStart`, `PreToolUse`, `PostToolUse`, etc.). yaah ships five handlers.

### Linter

Runs formatters and linters after every file edit. You pick which tools apply to which file extensions through profiles.

```go
h.Hooks().Register(handlers.NewLinterWith(
    handlers.GolangCILint(),   // .go -- gofmt + golangci-lint
    handlers.Ruff(),           // .py -- ruff check + format
    handlers.Prettier(),       // .js/.ts/.css/.md -- prettier
    handlers.TypeScript(),     // .ts/.tsx -- tsc --noEmit (type checking)
    handlers.Biome(),          // .ts/.tsx/.js/.jsx -- biome
    handlers.RustFmt(),        // .rs -- rustfmt + clippy
    handlers.GoVet(),          // .go -- gofmt + go vet (minimal)
))
```

> **Note:** ESLint is intentionally not included as a built-in profile because it requires per-project configuration (`.eslintrc`, plugins). `tsc --noEmit` runs per-file and works with or without a `tsconfig.json`.

> **Multi-profile support:** Multiple profiles can match the same file extension. For example, registering both `Prettier()` and `TypeScript()` means `.ts` files get formatted by Prettier first, then type-checked by tsc — both run in registration order.

You can also write your own profiles:

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

Blocks dangerous shell commands before they execute. Catches `rm -rf /`, `git push --force main`, `git reset --hard`, `DROP TABLE`, and similar patterns by default.

```go
guard := handlers.NewCommandGuard()
guard.Block(`kubectl delete namespace`, "namespace deletion")
h.Hooks().Register(guard)
```

### SecretScanner

Scans edited files for hardcoded credentials: AWS keys, GitHub PATs, OpenAI/Anthropic API keys, private keys, auth tokens, Slack tokens. Blocks the edit when it finds something.

```go
scanner := handlers.NewSecretScanner()
scanner.AddPattern(`my-custom-key-[0-9a-f]{32}`, "custom API key")
h.Hooks().Register(scanner)
```

### CommentChecker

Flags placeholder comments like `TODO: implement`, `FIXME`, `HACK`, `// ...`, and `// your code here`. Also catches non-English comments in edited files. Blocks the edit so the agent must fix the issue.

```go
checker := handlers.NewCommentChecker()
checker.AddPattern(`(?i)//\s*HACK`)
h.Hooks().Register(checker)
```

### SessionLogger

Writes session start/stop events with timestamps to a log file.

```go
h.Hooks().Register(handlers.NewSessionLogger("/var/log/claude-sessions"))
```

### Middleware chains

Chain multiple handlers into a sequential pipeline. Each link receives the previous link's result, enabling composed workflows like "scan for secrets, then suggest remediation if blocked."

```go
import "github.com/dirien/yet-another-agent-harness/pkg/hooks"

chain := hooks.NewChain("secret-remediation",
    []schema.HookEvent{schema.HookPostToolUse},
    regexp.MustCompile(`(?i)^(Edit|Write|MultiEdit)$`),
    hooks.HandlerLink(handlers.NewSecretScanner()),
    hooks.OnBlock(func(ctx context.Context, input *hooks.Input, prev *hooks.Result) (*hooks.Result, error) {
        prev.Output += "\n\nRemediation: Move the secret to an env var or secrets manager."
        return prev, nil
    }),
)
h.Hooks().Register(chain)
```

Available combinators:

| Combinator    | Description                                                    |
| ------------- | -------------------------------------------------------------- |
| `HandlerLink` | Wraps a `Handler` as a chain link                              |
| `OnBlock`     | Runs only when the previous link blocked (`result.Block=true`) |
| `OnError`     | Runs only when the previous link returned an error             |
| `Transform`   | Runs unconditionally, receives the previous result             |

Chains implement the `Handler` interface, so they can be registered like any other handler or nested inside other chains.

### Writing your own

Implement the `hooks.Handler` interface:

```go
type Handler interface {
    Name() string
    Events() []schema.HookEvent
    Match() *regexp.Regexp
    Execute(ctx context.Context, input *hooks.Input) (*hooks.Result, error)
}
```

Register it with `h.Hooks().Register(myHandler)`. The generated `settings.json` wires it to run via `yaah hook <event>`.

### HookHandler fields

| Field            | Type              | Description                                                     |
| ---------------- | ----------------- | --------------------------------------------------------------- |
| `type`           | string            | `command`, `http`, `prompt`, or `agent`                         |
| `command`        | string            | Shell command to run (`type=command`)                           |
| `url`            | string            | HTTP endpoint to POST to (`type=http`)                          |
| `prompt`         | string            | LLM prompt to evaluate (`type=prompt` or `agent`)               |
| `timeout`        | int               | Timeout in seconds (defaults: command=600, prompt=30, agent=60) |
| `statusMessage`  | string            | Custom spinner message while the hook runs                      |
| `once`           | bool              | Run only once per session (skills only)                         |
| `async`          | bool              | Run in background without blocking (`type=command` only)        |
| `headers`        | map[string]string | HTTP headers (`type=http` only)                                 |
| `allowedEnvVars` | []string          | Env var names for header interpolation (`type=http` only)       |
| `model`          | string            | Model identifier (`type=prompt` or `agent`)                     |

## MCP servers

Four built-in providers, plus a generic `Custom` for anything else.

```go
import "github.com/dirien/yet-another-agent-harness/pkg/mcp/providers"

h.MCP().Register(providers.NewContext7())
h.MCP().Register(providers.NewPulumi())
h.MCP().Register(providers.NewNotion("ntn-your-api-token"))
h.MCP().Register(providers.NewYaah())  // self-referencing yaah MCP server
h.MCP().Register(providers.NewCustom(schema.MCPServer{
    Name:      "my-server",
    Transport: schema.MCPTransportStdio,
    Command:   "my-mcp-binary",
    Args:      []string{"--mode", "production"},
}))
```

### yaah MCP server

yaah includes a built-in MCP server (`yaah serve`) that exposes its capabilities as tools Claude Code can call directly. The server uses the [official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk) and communicates over stdio.

| Tool                 | Description                                         |
| -------------------- | --------------------------------------------------- |
| `yaah_scan_secrets`  | Scan a file for hardcoded secrets and credentials   |
| `yaah_lint`          | Run lint checks on a file using configured profiles |
| `yaah_check_command` | Check whether a shell command is safe to run        |
| `yaah_doctor`        | Run health checks and report missing dependencies   |
| `yaah_session_info`  | Query session history or get server info            |

The MCP server reuses the same handler instances registered in the harness, so lint profiles, secret patterns, and command guard rules all apply.

#### Project-level discovery

`yaah generate` produces a `.mcp.json` file at the project root. Claude Code auto-discovers this file and connects to the servers listed in it. No manual configuration needed — just run `yaah generate` and start a new Claude Code session.

```json
{
  "mcpServers": {
    "yaah": {
      "type": "stdio",
      "command": "yaah",
      "args": ["serve"]
    }
  }
}
```

The `.mcp.json` file is separate from `.claude/settings.json`. It is Claude Code's official mechanism for project-level MCP server registration.

Transport options: `stdio`, `sse`, `streamable-http`, `http`, `websocket`.

### OAuth

Remote MCP servers that need OAuth:

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

yaah enables LSP servers through `enabledPlugins` in `settings.json`, pointing at the official Claude Code marketplace. Providers that implement `MarketplaceProvider` get added automatically. `yaah doctor` checks that the actual binaries exist on your machine.

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

Additional providers (yaml-ls, pulumi, pulumi-yaml) exist but aren't in the official marketplace, so they're left out of defaults.

## Skills

Skills produce `SKILL.md` files under `.claude/skills/`. Three are built in: `commit` (git commit workflow), `pr` (pull request creation), `review` (code review checklist).

```go
import "github.com/dirien/yet-another-agent-harness/pkg/skills/builtins"

h.Skills().Register(builtins.NewCommitSkill())
h.Skills().Register(builtins.NewPRSkill())
h.Skills().Register(builtins.NewReviewSkill())
```

### Remote skills

Pin a skill to a remote git repo. Works like GitHub Actions refs:

```go
import "github.com/dirien/yet-another-agent-harness/pkg/skills"

h.Skills().Register(skills.NewRemoteSkill(
    "team-standards",
    "Shared coding standards",
    "github.com/myorg/claude-skills@v1.2.0",
    "skills/standards/SKILL.md",
))
```

The ref can be a tag, branch, or commit SHA. Skills are cached in `~/.yaah/cache/skills/` and invalidated when the ref changes. Set `YAAH_HOME` to override the cache directory. Multi-file skills with reference directories implement `skills.SkillWithFiles`.

### Default remote skills

yaah ships 27 remote skills from four repos:

**pulumi/agent-skills** -- Pulumi IaC authoring and migration:

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

**dirien/claude-skills** -- Pulumi language-specific:

| Skill               | Description                              |
| ------------------- | ---------------------------------------- |
| `pulumi-typescript` | Pulumi TypeScript IaC with ESC and OIDC  |
| `pulumi-go`         | Pulumi Go IaC with ESC and OIDC          |
| `pulumi-python`     | Pulumi Python IaC with ESC and OIDC      |
| `pulumi-neo`        | Pulumi Neo conversational infrastructure |
| `pulumi-cli`        | Pulumi CLI command reference             |

**jeffallan/claude-skills** -- Development and operations:

| Skill                   | Description                                               |
| ----------------------- | --------------------------------------------------------- |
| `golang-pro`            | Go concurrency, microservices, gRPC, performance          |
| `kubernetes-specialist` | K8s deployments, Helm, RBAC, NetworkPolicies              |
| `devops-engineer`       | CI/CD, Docker, Kubernetes, Terraform, GitOps              |
| `python-pro`            | Python 3.11+, type safety, async, pytest, ruff            |
| `typescript-pro`        | Advanced TypeScript types, generics, tRPC                 |
| `csharp-developer`      | C# .NET 8+, ASP.NET Core, Blazor, EF Core                 |
| `javascript-pro`        | ES2023+ JavaScript, async/await, ESM, Node.js             |
| `cli-developer`         | CLI tools, argument parsing, shell completions            |
| `sre-engineer`          | SLOs, error budgets, incident response, capacity planning |
| `the-fool`              | Devil's advocate, pre-mortems, red teaming                |
| `architecture-designer` | System architecture, ADRs, scalability                    |
| `spring-boot-engineer`  | Spring Boot 3.x, Security 6, JPA, WebFlux                 |
| `code-reviewer`         | Code review for bugs, security, performance               |

**netresearch/agent-rules-skill** -- AGENTS.md generation:

| Skill         | Description                                                          |
| ------------- | -------------------------------------------------------------------- |
| `agent-rules` | Generate and maintain AGENTS.md files following agents.md convention |

### Skill frontmatter

Skills can implement `SkillWithFrontmatter` to inject YAML frontmatter into the generated `SKILL.md`:

```go
type SkillWithFrontmatter interface {
    Skill
    Frontmatter() SkillFrontmatter
}
```

| Field                      | Type   | Description                                 |
| -------------------------- | ------ | ------------------------------------------- |
| `argument-hint`            | string | Usage hint shown during autocomplete        |
| `disable-model-invocation` | bool   | Prevent Claude from auto-loading this skill |
| `user-invocable`           | \*bool | Show in /menu (default true)                |
| `allowed-tools`            | string | Comma-separated tool allowlist              |
| `model`                    | string | Model override when skill is active         |
| `context`                  | string | Set to `fork` for subagent execution        |
| `agent`                    | string | Subagent type when `context=fork`           |

### Writing your own

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

Agents produce markdown files with YAML frontmatter under `.claude/agents/`. Three ship as built-ins:

- **executor** (sonnet) -- single-task, strict verification
- **librarian** (haiku) -- research only, read-only tools
- **reviewer** (opus) -- code review, read-only tools

```go
import agentpkg "github.com/dirien/yet-another-agent-harness/pkg/agents"

h.Agents().Register(agentpkg.NewExecutor())
h.Agents().Register(agentpkg.NewLibrarian())
h.Agents().Register(agentpkg.NewReviewer())
```

### Remote agents

Pin an agent to a remote git repo, same syntax as remote skills:

```go
import agentpkg "github.com/dirien/yet-another-agent-harness/pkg/agents"

h.Agents().Register(agentpkg.NewRemoteAgent(
    "my-custom-agent",
    "What this agent does",
    "github.com/owner/repo@v1.0.0",
    "agents/my-agent.md",
    agentpkg.WithModel("sonnet"),
    agentpkg.WithTools("Read, Grep, Glob"),
))
```

The ref can be a tag, branch, or commit SHA. Agents are cached in `~/.yaah/cache/agents/` and invalidated when the ref changes. If the remote markdown file contains YAML frontmatter, it is stripped so that yaah can re-generate clean frontmatter from the agent's metadata fields.

Options:

| Option           | Description                             |
| ---------------- | --------------------------------------- |
| `WithModel()`    | Model override (sonnet/opus/haiku)      |
| `WithTools()`    | Comma-separated tool allowlist          |
| `WithAdvanced()` | Advanced frontmatter fields (see below) |

### Default remote agents

yaah ships 9 remote agents from [msitarzewski/agency-agents](https://github.com/msitarzewski/agency-agents):

| Agent                            | Model  | Description                                                  |
| -------------------------------- | ------ | ------------------------------------------------------------ |
| `agency-ai-engineer`             | sonnet | AI/ML engineering, model integration, LLM pipelines          |
| `agency-backend-architect`       | sonnet | Backend system design, API architecture, scalability         |
| `agency-security-engineer`       | sonnet | Security analysis, threat modeling, vulnerability assessment |
| `agency-code-reviewer`           | sonnet | Structured code review (read-only tools)                     |
| `agency-software-architect`      | opus   | System architecture, design patterns, technical decisions    |
| `agency-devops-automator`        | sonnet | CI/CD pipelines, infrastructure automation, deployments      |
| `agency-sre`                     | sonnet | Site reliability, observability, incident response           |
| `agency-api-tester`              | sonnet | API testing, contract validation, endpoint coverage          |
| `agency-performance-benchmarker` | sonnet | Performance profiling, load testing, optimization            |

### Advanced frontmatter

Implement `AgentWithAdvanced` for extra fields:

```go
type AgentWithAdvanced interface {
    Agent
    Advanced() AgentAdvanced
}
```

| Field             | Type           | Description                                                      |
| ----------------- | -------------- | ---------------------------------------------------------------- |
| `disallowedTools` | string         | Comma-separated tool denylist                                    |
| `permissionMode`  | string         | `default`, `acceptEdits`, `dontAsk`, `bypassPermissions`, `plan` |
| `maxTurns`        | int            | Maximum agentic turns before stopping                            |
| `skills`          | []string       | Skills to preload into agent context                             |
| `mcpServers`      | map[string]any | MCP servers for this agent                                       |
| `hooks`           | HooksConfig    | Lifecycle hooks (PreToolUse/PostToolUse/Stop)                    |
| `memory`          | string         | Persistent memory scope: `user`, `project`, `local`              |
| `background`      | bool           | Run as background task                                           |
| `isolation`       | string         | Isolation mode: `worktree`                                       |

## Commands

Slash commands produce markdown files under `.claude/commands/`. Implement `commands.Command`:

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

### Advanced frontmatter

Implement `CommandWithAdvanced` for extra fields:

```go
type CommandWithAdvanced interface {
    Command
    Advanced() CommandAdvanced
}
```

| Field                      | Type   | Description                                   |
| -------------------------- | ------ | --------------------------------------------- |
| `disable-model-invocation` | bool   | Prevent Claude from auto-loading this command |
| `user-invocable`           | \*bool | Show in /menu (default true)                  |
| `context`                  | string | Set to `fork` for subagent execution          |
| `agent`                    | string | Subagent type when `context=fork`             |

## Plugins

The `schema.Plugin` struct matches the official Claude Code `plugin.json` spec.

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

## Session store

The session store provides a persistent audit trail for every Claude Code session. It records hook events, tool calls, blocked commands, file modifications, and security findings to JSON files.

### How it works

1. Claude Code sets a `CLAUDE_SESSION_ID` environment variable for each session
2. Every hook event flows through the harness's `HandleHookEvent` method
3. The harness loads (or creates) a session file at `.claude/sessions/<session-id>.json`
4. It records the event details and persists back to disk with atomic writes (temp file + rename)

### What gets tracked

| Field            | Type             | Description                                       |
| ---------------- | ---------------- | ------------------------------------------------- |
| `id`             | string           | Session identifier from Claude Code               |
| `started_at`     | timestamp        | When the session started                          |
| `last_event_at`  | timestamp        | When the last hook event fired                    |
| `event_count`    | int              | Total number of hook events                       |
| `tool_calls`     | []ToolCallRecord | Every tool invocation with name, input, timestamp |
| `blocked_calls`  | []ToolCallRecord | Commands blocked by the guard, with reasons       |
| `files_modified` | []string         | Absolute paths of files that were edited          |
| `findings`       | []Finding        | Secrets detected, lint issues, security findings  |

### CLI commands

```bash
yaah session list              # List recent sessions (sorted by last activity)
yaah session show <id>         # Full details: tool calls, blocked commands, findings
yaah session clean             # Remove sessions older than 7 days
```

### MCP access

The `yaah_session_info` MCP tool lets Claude Code query session history mid-conversation:

```
# Called by Claude Code via MCP:
yaah_session_info {}                          # Returns server info
yaah_session_info {"session_id": "abc123"}    # Returns full session details
```

### Programmatic access

```go
store := harness.SessionStore()

// Load a specific session
sess, err := store.Load("session-id")

// List all sessions
sessions, err := store.List()

// Save a session
err := store.Save(sess)

// Clean up old sessions
deleted, err := store.Cleanup(7 * 24 * time.Hour)
```

### Security

Session IDs are validated to prevent path traversal attacks. IDs containing path separators (`/`, `\`) or special values (`.`, `..`) are rejected.
