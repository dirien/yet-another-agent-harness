package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ImportCommand imports an existing project into the planning workflow via /yaah:import.
type ImportCommand struct{}

// NewImportCommand creates a new ImportCommand.
func NewImportCommand() *ImportCommand { return &ImportCommand{} }

func (c *ImportCommand) Name() string        { return "yaah/import" }
func (c *ImportCommand) Description() string {
	return "Import an existing project into the planning workflow"
}
func (c *ImportCommand) ArgumentHint() string { return "" }
func (c *ImportCommand) Model() string        { return "" }
func (c *ImportCommand) AllowedTools() string { return "" }
func (c *ImportCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ImportCommand) Content() string {
	return `# /yaah:import — Import Existing Project

## When to use
When the user runs ` + "`/yaah:import`" + ` on a project that already has code but no ` + "`.planning/`" + ` directory.

## Difference from /yaah:init
- ` + "`/yaah:init`" + ` is for new projects — it captures intent and detects a stack.
- ` + "`/yaah:import`" + ` is for existing projects — it analyzes what is ALREADY BUILT and reverse-engineers requirements from the code.

## Prerequisites
- Confirm ` + "`.planning/`" + ` does not already exist. If it does, ask the user before proceeding.

## Steps

### 1. Deep codebase analysis
- Read package managers: go.mod, package.json, Cargo.toml, pyproject.toml, pom.xml, build.gradle
- Run ` + "`git log --oneline -50`" + ` to understand the full history
- Run ` + "`find . -maxdepth 3 -type f | head -100`" + ` for a structure overview
- Identify: language, framework, build system, test framework, CI/CD, deployment targets
- Map all major components, their responsibilities, and how they communicate

### 2. Reverse-engineer requirements
- For each major component or feature discovered, derive a requirement statement
- Format as: "The system shall <capability>" or "The system supports <feature>"
- Assign REQ-IDs starting at REQ-001

### 3. Map existing features to REQ-IDs
- Create a feature inventory: what the codebase currently does, validated against tests and docs

### 4. Create REQUIREMENTS.md
Write to ` + "`.planning/REQUIREMENTS.md`" + ` with two sections:
- **Implemented**: requirements already satisfied by the existing code (tagged REQ-IDs)
- **Planned**: requirements inferred from TODOs, issues, or partial implementations

### 5. Create ROADMAP.md
- Create completed phases for all existing work (status: imported)
- Each imported phase covers a logical chunk of what was already built
- Add placeholder future phases based on the Planned section of REQUIREMENTS.md

### 6. Mark completed phases as imported
- Imported phases get no PLAN.md or SUMMARY.md — just a directory with status: imported
- Create ` + "`.planning/phases/{NN}-{slug}/`" + ` directories for each imported phase

### 7. Create PROJECT.md, STATE.md, config.json
**PROJECT.md**:
- Project name and description (from README or inferred)
- Detected tech stack with sources
- Key architectural decisions already in place

**STATE.md** (with YAML frontmatter):
- milestone: v1.0, phase: imported, status: imported, last_updated timestamp
- Current Position: "Imported from existing codebase"
- Decisions Made: list key architectural decisions already present in the code

**config.json**:
- mode: interactive, granularity: standard, model_profile: quality
- workflow toggles for research, plan_check, discuss, auto_advance

### 8. Generate research/ from the deep analysis
Write to ` + "`.planning/research/`" + `:
- ` + "`stack.md`" + ` — detected stack with sources
- ` + "`architecture.md`" + ` — component map and data flow
- ` + "`pitfalls.md`" + ` — TODO/FIXME density, test coverage gaps, known issues

### 9. Ask about what to build next
After writing all files, ask the user:
- What do you want to build or improve next?
- Are there open issues, planned features, or known gaps you want to tackle?
- Any constraints (timeline, team, tech choices)?

Use the answers to populate the Planned section and future phases.

### 10. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): import existing project into planning workflow\"`" + `

### 11. Report
Print a summary:
- How many implemented requirements were found
- How many imported phases were created
- What future phases were drafted
- Suggest: "Run ` + "`/yaah:discuss 1`" + ` to plan your next phase."

## Rules
- NEVER overwrite existing ` + "`.planning/`" + ` without explicit user confirmation
- Tag every detected fact with its source (e.g., "Go 1.25 (from go.mod)")
- NEVER invent features — only document what is observably present in the code
- Mark all imported phases clearly so the team knows they were not planned through this workflow
`
}
