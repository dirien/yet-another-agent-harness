# Configuration

## Using yaah

There are two ways to use yaah.

### CLI (built-in defaults)

The fastest path. yaah has all its defaults compiled into the binary:

```bash
cd your-repo
yaah generate
```

That's it. yaah uses its built-in `AllDefaults()` to generate the full `.claude/` directory with every handler, skill, LSP server, and MCP provider enabled.

### Go library (full control)

For teams that want programmatic control, import yaah as a Go library. Use `DefaultOptions` to pick exactly which components to enable:

```go
// Pick only what you need
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

Or go fully custom by registering components one by one:

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

Run it with `go run ./cmd/your-setup/` whenever you change your config.

## Settings

`schema.Settings` maps directly to the official Claude Code `settings.json` spec. Here are the field groups:

| Group                    | Fields                                                                                                                                                                       |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Core                     | `model`, `alwaysThinkingEnabled`, `autoUpdatesChannel`, `effortLevel`, `statusLine`, `env`, `teammateMode`                                                                   |
| Model and performance    | `availableModels`, `fastMode`, `fastModePerSessionOptIn`                                                                                                                     |
| Permissions and security | `permissions`, `sandbox`, `allowManagedPermissionRulesOnly`                                                                                                                  |
| Hooks and automation     | `disableAllHooks`, `allowManagedHooksOnly`                                                                                                                                   |
| Git and attribution      | `attribution`, `includeGitInstructions`                                                                                                                                      |
| Authentication           | `apiKeyHelper`                                                                                                                                                               |
| UI and behavior          | `language`, `outputStyle`, `showTurnDuration`, `spinnerVerbs`, `spinnerTipsEnabled`, `spinnerTipsOverride`, `prefersReducedMotion`, `terminalProgressBarEnabled`             |
| Plugins                  | `enabledPlugins`, `pluginConfigs`, `extraKnownMarketplaces`, `strictKnownMarketplaces`, `skippedMarketplaces`, `skippedPlugins`, `blockedMarketplaces`, `pluginTrustMessage` |
| MCP management           | `enableAllProjectMcpServers`, `enabledMcpjsonServers`, `disabledMcpjsonServers`, `allowedMcpServers`, `deniedMcpServers`, `allowManagedMcpServersOnly`                       |
| Organization             | `companyAnnouncements`, `cleanupPeriodDays`, `plansDirectory`, `autoMemoryEnabled`, `skipWebFetchPreflight`                                                                  |
| File and directory       | `fileSuggestion`, `respectGitignore`, `additionalDirectories`                                                                                                                |

## CLI reference

```
yaah generate              # Generate .claude/ with built-in defaults
yaah generate -o ./out     # Output to a different directory
yaah hook <event>          # Runtime hook dispatcher (called by Claude Code)
yaah info                  # Show all registered components
yaah doctor                # Health check: validates binaries and config
yaah version               # Print version, commit, and build date
```

## Architecture

yaah uses an interface + registry pattern. Each domain has an interface for individual components and a registry that holds them:

| Domain   | Interface          | Registry            | What it does                                 |
| -------- | ------------------ | ------------------- | -------------------------------------------- |
| Hooks    | `hooks.Handler`    | `hooks.Registry`    | Run code on Claude Code lifecycle events     |
| MCP      | `mcp.Provider`     | `mcp.Registry`      | Configure MCP servers                        |
| LSP      | `lsp.Provider`     | `lsp.Registry`      | Configure LSP servers with binary validation |
| Skills   | `skills.Skill`     | `skills.Registry`   | Generate SKILL.md files                      |
| Agents   | `agents.Agent`     | `agents.Registry`   | Generate agent markdown files                |
| Commands | `commands.Command` | `commands.Registry` | Generate slash command files                 |
| Plugins  | `plugins.Plugin`   | `plugins.Registry`  | Generate plugin packages                     |

The `Harness` struct in `pkg/harness/` wires all registries together. Call `GenerateConfig()` to build the config and `WriteAll()` to write output files.

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
pkg/generator/         settings.json generation
cmd/yaah/              CLI entry point
internal/cli/          CLI utilities
```

## What gets generated

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

MCP servers go inline in `settings.json` under the `mcpServers` key. LSP servers are enabled via `enabledPlugins` in `settings.json`, referencing official Claude Code marketplace plugins (e.g. `gopls-lsp@claude-plugins-official`).
