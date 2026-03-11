# yaah

**yet another agent harness**

Configure Claude Code once, use it everywhere.

## The problem

Claude Code configuration is a mess. Settings live in JSON files, skills are markdown scattered across directories, hooks are shell scripts wired by hand, MCP servers need manual JSON entries, and LSP plugins require marketplace clicks. Multiply that by the number of repos you work in. Good luck keeping any of it consistent.

## What yaah does

yaah generates the entire `.claude/` directory from Go code: settings, hooks, skills, agents, MCP servers, LSP plugins. One command, every repo, same result.

```bash
yaah generate
```

That single command gives you:

- 5 hooks out of the box: linting (golangci-lint, ruff, prettier), a command guard that blocks `rm -rf /` and friends, a secret scanner for leaked keys, a comment checker that catches `TODO: implement` placeholders, and a session logger
- MCP servers for Context7 and Pulumi, ready to go
- LSP support for Go, Python, TypeScript, and C# via the official marketplace
- 3 built-in skills (commit, PR, review) plus 26 remote skills covering Pulumi IaC, Go, Python, TypeScript, Kubernetes, DevOps, SRE, and more
- 3 agents: executor (sonnet), librarian (haiku), reviewer (opus)

Don't want all of it? Turn off what you don't need:

```go
opts := harness.DefaultOptions{
    EnableCommandGuard:  true,
    EnableSecretScanner: true,
    EnableGopls:         true,
    EnableCommitSkill:   true,
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

## How it works

yaah has a simple mental model: interfaces and registries. Each component type (hooks, MCP, LSP, skills, agents, commands) has an interface you implement and a registry you add it to. The `Harness` wires them all together and spits out the right files.

```
.claude/
├── settings.json       # hooks, MCP servers, enabledPlugins, all settings
├── skills/             # SKILL.md files (built-in + remote)
├── agents/             # agent markdown with YAML frontmatter
└── commands/           # slash command definitions
```

Write your own hook? Implement `hooks.Handler`. Custom MCP server? Implement `mcp.Provider`. Same pattern everywhere.

## Documentation

The detailed reference lives in [`docs/`](docs/):

- [Components](docs/components.md) -- hooks, MCP, LSP, skills, agents, commands, plugins
- [Configuration](docs/configuration.md) -- settings, CLI commands, project structure

## License

[MIT](LICENSE)
