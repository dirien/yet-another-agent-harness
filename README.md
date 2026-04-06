# yaah

**yet another agent harness**

Configure your coding agent once, use it everywhere.

## The problem

Coding agent configuration is a mess. Settings live in JSON files, skills are markdown scattered across directories, hooks are shell scripts wired by hand, MCP servers need manual JSON entries. Multiply that by the number of repos and agents you use. Good luck keeping any of it consistent.

## What yaah does

yaah generates configuration for **four coding agents** from a single Go codebase: Claude Code, OpenCode, Codex CLI, and GitHub Copilot CLI. One command per agent, every repo, same result.

```bash
yaah generate                      # all agents
yaah generate --agent claude       # Claude Code only
yaah generate --agent opencode     # OpenCode only
yaah generate --agent codex        # Codex CLI only
yaah generate --agent copilot      # GitHub Copilot CLI only
```

That single command gives you:

- 5 hooks out of the box: linting (golangci-lint, ruff, prettier, tsc), a command guard that blocks `rm -rf /` and friends, a secret scanner for leaked keys, a comment checker that catches `TODO: implement` placeholders, and a session logger
- Middleware chains for composing handlers (e.g. secret scan + auto-remediation advice)
- MCP servers for Context7 and Pulumi, plus a built-in yaah MCP server exposing tools like secret scanning, linting, and command checking directly to the agent
- Multi-agent config generation with per-agent adaptations (MCP format, hook delivery, agent tools, skill frontmatter)
- LSP support for Go, Python, TypeScript, and C# via the official marketplace
- Session tracking that logs every tool call, blocked command, and file modification across sessions
- 3 built-in skills (commit, PR, review) plus 74 remote and workflow skills covering Pulumi IaC, Flux CD GitOps, Go, Python, TypeScript, Kubernetes, DevOps, SRE, security auditing, code quality, tech debt analysis, and more
- 16 agents: 3 built-in (executor, librarian, reviewer) plus 4 workflow agents (researcher, planner, doc-writer, verifier) and 9 remote agents from [agency-agents](https://github.com/msitarzewski/agency-agents) covering AI engineering, backend architecture, security, code review, DevOps, SRE, and testing
- Plugin support with marketplace auto-enablement — ships with [OpenAI Codex](https://github.com/openai/codex-plugin-cc) for code review and task delegation

Don't want all of it? Turn off what you don't need:

```go
opts := harness.DefaultOptions{
    EnableCommandGuard:  true,
    EnableSecretScanner: true,
    EnableGopls:         true,
    EnableCommitSkill:   true,
    EnableYaahMCP:       true,
}
h := harness.NewWithDefaults(opts)
```

## Install

### Homebrew

```bash
brew install dirien/tap/yaah
```

### Go

```bash
go install github.com/dirien/yet-another-agent-harness/cmd/yaah@latest
```

### Binary

Grab a release from [GitHub Releases](https://github.com/dirien/yet-another-agent-harness/releases). Binaries are signed with cosign and include SBOMs.

## Quick start

```bash
# Generate .claude/ with all defaults
cd your-repo
yaah generate

# Check that everything is installed
yaah doctor
```

`yaah doctor` tells you what's missing and how to install it:

```
LSP Servers:
  ✓ gopls                    /usr/local/bin/gopls
  ✓ pyright                  /usr/local/bin/pyright-langserver
  ✗ csharp                   not found → dotnet tool install -g csharp-ls
```

## Workflow

yaah includes a structured project workflow with 29 slash commands. All commands are namespaced as `/yaah:*` and are registered as explicit user commands — the model never auto-triggers them.

### Lifecycle

```
/yaah:init → /yaah:discuss → /yaah:plan → /yaah:execute → /yaah:verify → /yaah:ship
                                                                    ↓
/yaah:next ← /yaah:progress ← /yaah:complete-milestone ← /yaah:docs
```

### Command reference

| Command | Description | Subagent |
|---------|-------------|----------|
| **Core Workflow** | | |
| `/yaah:init` | Project onboarding: discover codebase, set vision, create roadmap | Yes |
| `/yaah:discuss` | Capture implementation decisions before planning | Yes |
| `/yaah:plan` | Create wave-grouped implementation plans | Yes |
| `/yaah:execute` | Execute plans wave by wave with parallel subagents | Yes |
| `/yaah:verify` | Three-level artifact verification against plan | Yes |
| `/yaah:docs` | Generate codebase-verified documentation | Yes |
| `/yaah:next` | Auto-detect and recommend next step | No |
| `/yaah:quick` | Execute a task without full planning | Yes |
| **Shipping & Milestones** | | |
| `/yaah:ship` | Create PR from verified phase work | Yes |
| `/yaah:complete-milestone` | Archive milestone, tag release, generate changelog | Yes |
| `/yaah:new-milestone` | Start new version cycle with fresh goals | Yes |
| **Session Management** | | |
| `/yaah:pause` | Save session state for later resumption | No |
| `/yaah:resume` | Resume from previous session handoff | No |
| **Phase Management** | | |
| `/yaah:add-phase` | Add phase to end of roadmap | No |
| `/yaah:insert-phase` | Insert urgent phase between existing ones | No |
| `/yaah:remove-phase` | Remove a future phase | No |
| **Quality & Security** | | |
| `/yaah:review` | Structured code review of phase implementation | Yes |
| `/yaah:secure` | STRIDE threat modeling and vulnerability analysis | Yes |
| `/yaah:health` | Validate `.planning/` integrity and consistency | No |
| **Status & Capture** | | |
| `/yaah:progress` | Detailed project progress with metrics | No |
| `/yaah:todo` | Capture, list, or complete quick todo items | No |
| `/yaah:note` | Zero-friction idea capture | No |
| **Configuration** | | |
| `/yaah:settings` | View or update workflow configuration | No |
| **Analysis & Advanced** | | |
| `/yaah:explore` | Interactive codebase exploration | Yes |
| `/yaah:scan` | Scan for security, quality, dependency issues | Yes |
| `/yaah:import` | Import existing project into planning workflow | Yes |
| `/yaah:autonomous` | Run full workflow without human intervention | Yes |
| `/yaah:forensics` | Investigate failed or stuck workflow runs | Yes |
| `/yaah:cleanup` | Clean up temporary planning artifacts | No |

### .planning/ directory structure

```
.planning/
├── PROJECT.md          # Vision, goals, tech stack, constraints
├── REQUIREMENTS.md     # Scoped requirements with REQ-IDs
├── ROADMAP.md          # Phases with scope, success criteria, status
├── STATE.md            # Current position and progress tracking
├── HANDOFF.md          # Session pause/resume state
├── TODOS.md            # Quick todo items
├── config.json         # Workflow settings
├── research/           # Project-level research
├── phases/             # Per-phase artifacts
│   └── 01-auth/
│       ├── CONTEXT.md        # Implementation decisions
│       ├── RESEARCH.md       # Phase-specific research
│       ├── 01-01-PLAN.md     # Task plan with wave grouping
│       ├── 01-01-SUMMARY.md  # Execution outcomes
│       ├── VERIFICATION.md   # Validation results
│       ├── REVIEW.md         # Code review findings
│       ├── SECURITY.md       # Threat model (STRIDE)
│       └── CHANGELOG.md      # Milestone changelog
├── quick/              # Ad-hoc task records
└── notes/              # Idea captures
```

Four specialized agents back the workflow: `researcher` (sonnet, read-only), `planner` (opus, goal-backward decomposition), `doc-writer` (sonnet, codebase-verified docs), and `verifier` (sonnet, artifact validation). Two MCP tools are also exposed: `yaah_planning_status` and `yaah_planning_init`.

See [Components documentation](docs/components.md#workflow-commands) for full details.

## Runtime features

yaah isn't just a config generator — it also runs as a runtime alongside Claude Code.

### MCP server

`yaah serve` starts an MCP server over stdio, exposing yaah's capabilities as tools that Claude Code can call directly:

| Tool                 | Description                                       |
| -------------------- | ------------------------------------------------- |
| `yaah_scan_secrets`  | Scan a file for hardcoded secrets and credentials |
| `yaah_lint`          | Run lint checks using configured profiles         |
| `yaah_check_command` | Check whether a shell command is safe to run      |
| `yaah_doctor`        | Run health checks and report missing dependencies |
| `yaah_session_info`  | Query session history or get server info          |

The MCP server is auto-discovered by Claude Code via `.mcp.json` at the project root (generated by `yaah generate`).

### Session tracking

Every hook event is recorded to `.claude/sessions/<session-id>.json`. This gives you a full audit trail of what Claude Code did in each session: tool calls, blocked commands, files modified, and security findings.

```bash
yaah session list              # List recent sessions
yaah session show <id>         # Full details for a session
yaah session clean             # Remove sessions older than 7 days
```

### Middleware chains

Compose handlers into sequential pipelines with conditional logic:

```go
chain := hooks.NewChain("secret-remediation",
    []schema.HookEvent{schema.HookPostToolUse},
    regexp.MustCompile(`(?i)^(Edit|Write)$`),
    hooks.HandlerLink(handlers.NewSecretScanner()),
    hooks.OnBlock(func(ctx context.Context, input *hooks.Input, prev *hooks.Result) (*hooks.Result, error) {
        prev.Output += "\n\nRemediation: Move the secret to an env var or secrets manager."
        return prev, nil
    }),
)
h.Hooks().Register(chain)
```

Available combinators: `HandlerLink`, `OnBlock`, `OnError`, `Transform`.

## How it works

yaah has a simple mental model: interfaces and registries. Each component type (hooks, MCP, LSP, skills, agents, commands) has an interface you implement and a registry you add it to. The `Harness` wires them all together and the per-agent generators produce the right files.

### Multi-agent output

Each agent gets files in its native format:

| Agent        | Settings                | MCP                        | Hooks                       | Skills              | Agents                      |
| ------------ | ----------------------- | -------------------------- | --------------------------- | ------------------- | --------------------------- |
| **Claude**   | `.claude/settings.json` | `.mcp.json`                | embedded in settings        | `.claude/skills/`   | `.claude/agents/*.md`       |
| **OpenCode** | `opencode.json`         | embedded (`"mcp"` key)     | `.opencode/plugins/yaah.js` | `.opencode/skills/` | `.opencode/agents/*.md`     |
| **Codex**    | `.codex/config.toml`    | embedded (`[mcp_servers]`) | `.codex/hooks.json`         | `.agents/skills/`   | not supported               |
| **Copilot**  | none                    | `.copilot/mcp-config.json` | `.github/hooks/hooks.json`  | `.github/skills/`   | `.github/agents/*.agent.md` |

Key adaptations per agent:

- **OpenCode**: MCP uses `"mcp"` key with `"local"`/`"remote"` types and `"command"` as array. Agent tools rendered as a disable-map. Hooks delivered via JS plugin (`execFileSync`).
- **Copilot**: MCP uses `"stdio"`/`"http"` transport types. Env vars passed through. Agent files use `.agent.md` extension.
- **Codex**: MCP embedded in TOML. Hooks limited to `SessionStart`/`Stop` — linting and security checks available via yaah MCP tools.
- **Claude**: Native format, all features supported.

Write your own hook? Implement `hooks.Handler`. Custom MCP server? Implement `mcp.Provider`. Same pattern everywhere.

## Documentation

The detailed reference lives in [`docs/`](docs/):

- [Components](docs/components.md) -- hooks, MCP, LSP, skills, agents, commands, plugins, MCP server, session store, middleware chains
- [Configuration](docs/configuration.md) -- settings, CLI commands, project structure

## License

[MIT](LICENSE)
