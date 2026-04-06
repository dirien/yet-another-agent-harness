package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// DiscussCommand implements implementation decision capture via /yaah:discuss.
type DiscussCommand struct{}

// NewDiscussCommand creates a new DiscussCommand.
func NewDiscussCommand() *DiscussCommand { return &DiscussCommand{} }

func (c *DiscussCommand) Name() string { return "yaah/discuss" }
func (c *DiscussCommand) Description() string {
	return "Capture implementation decisions before planning"
}
func (c *DiscussCommand) ArgumentHint() string { return "<phase-number>" }
func (c *DiscussCommand) Model() string        { return "" }
func (c *DiscussCommand) AllowedTools() string { return "" }
func (c *DiscussCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *DiscussCommand) Content() string {
	return `# /yaah:discuss — Capture Implementation Decisions

## When to use
When the user runs ` + "`/yaah:discuss <phase-number>`" + ` to eliminate gray areas before planning a phase.

## Goal
Surface and resolve HOW to build before writing a single task. Every gray area resolved here is one ambiguity that cannot derail execution.

## Steps

### 1. Load context
- Read ` + "`.planning/PROJECT.md`" + `, ` + "`.planning/REQUIREMENTS.md`" + `, ` + "`.planning/STATE.md`" + `
- Check for an earlier ` + "`CONTEXT.md`" + ` in the target phase directory
- Identify the phase scope from ROADMAP.md

### 2. Scout codebase for existing patterns
- Find analogous implementations that already exist
- Note reusable assets, established conventions, integration points
- Tag every finding with its source file and line range

### 3. Identify 3-5 gray areas
Focus on the phase scope. Common gray area categories:
- **Visual/UI**: Component libraries, styling approach, responsive behavior
- **API/CLI**: Interface shape, error handling strategy, pagination
- **Data**: Schema choices, migration strategy, caching layer
- **Infrastructure**: Deployment target, environment config, secrets management

Select only the areas with the highest ambiguity cost if left unresolved.

### 4. Present and discuss each area
For each gray area:
- State the question clearly
- Offer 2-3 concrete options with brief tradeoffs
- Ask the user to choose or propose an alternative
- Record the decision with rationale

### 5. Write CONTEXT.md
Write to ` + "`.planning/phases/{NN}-{slug}/CONTEXT.md`" + ` (create directory if needed).

Structure:
` + "```" + `
# Phase {N} Context

## Decisions
- D-01: <decision title> — <rationale>
- D-02: <decision title> — <rationale>

## Agent Discretion
<what the executing agent may decide independently>

## Code Context
### Reusable Assets
- <file>: <what it provides>

### Patterns
- <pattern name>: <where it's used>

### Integration Points
- <component>: <how it connects>

## Deferred Ideas
- <idea>: deferred to <phase or backlog>
` + "```" + `

### 6. Update STATE.md
Set ` + "`phase: {N}`" + `, ` + "`status: discussed`" + `, update ` + "`last_updated`" + `.

### 7. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): capture phase {N} decisions\"`" + `

## Rules
- NEVER ask WHAT to build — that lives in REQUIREMENTS.md
- Ask only HOW: implementation approach, technology choices, design patterns
- Maximum 3-5 questions per session — prioritize ruthlessly
- Always reference codebase patterns when suggesting options
- Reference decisions as D-XX in all downstream artifacts
`
}
