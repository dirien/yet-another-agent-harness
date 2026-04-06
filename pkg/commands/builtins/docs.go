package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// DocsCommand implements documentation generation via /yaah:docs.
type DocsCommand struct{}

// NewDocsCommand creates a new DocsCommand.
func NewDocsCommand() *DocsCommand { return &DocsCommand{} }

func (c *DocsCommand) Name() string { return "yaah/docs" }
func (c *DocsCommand) Description() string {
	return "Generate or update project documentation from codebase analysis"
}
func (c *DocsCommand) ArgumentHint() string { return "[--force] [doc-type]" }
func (c *DocsCommand) Model() string        { return "" }
func (c *DocsCommand) AllowedTools() string { return "" }
func (c *DocsCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *DocsCommand) Content() string {
	return `# /yaah:docs — Documentation Generation

## When to use
When the user runs ` + "`/yaah:docs`" + ` to generate or refresh project documentation from codebase analysis.

## Goal
Produce accurate, source-verified documentation. Never guess. Mark every undiscoverable claim.

## Steps

### 1. Detect project type
Identify the primary project type by reading source files:
- **API**: has route definitions, handler files, OpenAPI spec, or HTTP server setup
- **CLI**: has ` + "`cmd/`" + ` directory, cobra/click/argparse usage, or flag definitions
- **Library**: no main entrypoint, exports a public API surface
- **Monorepo**: multiple packages or workspaces with independent entrypoints
- **Web app**: has frontend framework (React, Vue, Svelte, etc.) or static site config

Tag the detected type with its source (e.g., "CLI (from cmd/root.go)").

### 2. Inventory existing docs
- Glob for ` + "`*.md`" + `, ` + "`docs/**`" + `, ` + "`README*`" + `
- For each file, check for the ` + "`<!-- generated-by: yaah-docs -->`" + ` marker
- Files with the marker: safe to regenerate
- Files without the marker: prompt user before overwriting (unless ` + "`--force`" + ` is set)

### 3. Select doc set by project type
Always generate:
- ` + "`README.md`" + ` — project overview, quickstart, badges
- ` + "`ARCHITECTURE.md`" + ` — component map, data flow, key decisions

Conditionally generate:
- API: ` + "`docs/api.md`" + ` (endpoint reference), ` + "`docs/errors.md`" + `
- CLI: ` + "`docs/commands.md`" + ` (flag reference per subcommand)
- Library: ` + "`docs/usage.md`" + ` (integration guide), ` + "`docs/examples.md`" + `
- Monorepo: ` + "`docs/packages.md`" + ` (package inventory and dependencies)
- Web app: ` + "`docs/deployment.md`" + `, ` + "`docs/configuration.md`" + `

If a specific ` + "`doc-type`" + ` argument was provided, generate only that document.

### 4. Wave-based generation
**Wave 1** (foundation — must complete before Wave 2):
- README.md and ARCHITECTURE.md
- Spawn in parallel using Agent(subagent_type="doc-writer", run_in_background=true)

**Wave 2** (supplementary — may reference Wave 1 outputs):
- All conditional docs selected in step 3
- Spawn in parallel using Agent(subagent_type="doc-writer", run_in_background=true)

Each doc-writer agent receives:
- The target file path
- The project type and detected facts (tagged with sources)
- Instruction to add ` + "`<!-- generated-by: yaah-docs -->`" + ` as the first HTML comment

### 5. Verification pass
For each generated document:
- Extract factual claims (version numbers, file paths, command syntax)
- Cross-check each claim against source files
- Replace unverifiable claims with ` + "`<!-- VERIFY: <claim> -->`" + ` markers
- Report the count of VERIFY markers to the user

### 6. Commit
Run: ` + "`git add docs/ README.md ARCHITECTURE.md && git commit -m \"docs: regenerate project documentation\"`" + `

## Rules
- NEVER guess: if a fact cannot be found in source files, use ` + "`<!-- VERIFY: ... -->`" + ` instead
- ALWAYS add the ` + "`<!-- generated-by: yaah-docs -->`" + ` marker to every generated file
- NEVER overwrite files without the marker unless ` + "`--force`" + ` is passed
- Tag every detected fact with its source file
- Wave 2 must not start until Wave 1 docs are written to disk
`
}
