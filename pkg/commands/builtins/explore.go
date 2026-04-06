package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ExploreCommand provides interactive codebase exploration via /yaah:explore.
type ExploreCommand struct{}

// NewExploreCommand creates a new ExploreCommand.
func NewExploreCommand() *ExploreCommand { return &ExploreCommand{} }

func (c *ExploreCommand) Name() string        { return "yaah/explore" }
func (c *ExploreCommand) Description() string { return "Interactive codebase exploration and analysis" }
func (c *ExploreCommand) ArgumentHint() string { return "[area-or-question]" }
func (c *ExploreCommand) Model() string        { return "" }
func (c *ExploreCommand) AllowedTools() string { return "" }
func (c *ExploreCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ExploreCommand) Content() string {
	return `# /yaah:explore — Interactive Codebase Exploration

## When to use
When the user runs ` + "`/yaah:explore [area-or-question]`" + ` to understand the codebase or a specific area within it.

## Behavior
- If an argument is given: investigate that specific area or answer that question in depth.
- If no argument is given: provide a broad codebase overview.

## Steps

### 1. Map the project structure
- Read package managers: go.mod, package.json, Cargo.toml, pyproject.toml, pom.xml, build.gradle
- Run ` + "`find . -maxdepth 3 -type f | head -80`" + ` for structure overview
- Identify key directories, entry points, and their purposes

### 2. Identify the tech stack and dependencies
- Enumerate external dependencies and their roles
- Note the language, framework, build system, and test framework in use

### 3. Find the main abstractions
- Locate primary interfaces, key types, and dominant design patterns
- Note how components are wired together (dependency injection, registries, factories, etc.)

### 4. Discover tests and their coverage areas
- Find all test files and what they exercise
- Note any integration or end-to-end tests vs unit tests

### 5. Identify configuration points
- Environment variables, config files, feature flags
- Default values and where they are set

## Output format
Produce a structured exploration report:

` + "```" + `
# Codebase Exploration

## Structure
{directory tree with purpose annotations}

## Entry Points
{main files, CLI commands, API routes}

## Key Abstractions
{interfaces, types, patterns}

## Dependencies
{external deps with purpose}

## Test Coverage
{test files, what they cover}

## Configuration
{env vars, config files, feature flags}
` + "```" + `

## Persistence
If ` + "`.planning/`" + ` exists, save the report to ` + "`.planning/research/exploration.md`" + `.

## Rules
- Tag every finding with its source file (e.g., "registry pattern (pkg/hooks/registry.go)")
- When a specific area or question was provided, focus the report on that area first, then summarize the rest
- NEVER guess — only report what is directly observable in the codebase
`
}
