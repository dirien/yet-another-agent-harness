package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// InitCommand implements project onboarding via /yaah:init.
type InitCommand struct{}

// NewInitCommand creates a new InitCommand.
func NewInitCommand() *InitCommand { return &InitCommand{} }

func (c *InitCommand) Name() string { return "yaah/init" }
func (c *InitCommand) Description() string {
	return "Project onboarding: discover codebase, set vision, create roadmap"
}
func (c *InitCommand) ArgumentHint() string { return "[project description]" }
func (c *InitCommand) Model() string        { return "" }
func (c *InitCommand) AllowedTools() string { return "" }
func (c *InitCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *InitCommand) Content() string {
	return `# /yaah:init — Project Onboarding

## When to use
When the user runs /yaah:init or starts working in a new codebase for the first time.

## Prerequisites
- Confirm ` + "`.planning/`" + ` does not already exist. If it does, ask user before overwriting.

## Steps

### 1. Detect project
- Read package managers: go.mod, package.json, Cargo.toml, pyproject.toml, pom.xml, build.gradle
- Run ` + "`git log --oneline -20`" + ` to understand recent activity
- Run ` + "`find . -maxdepth 2 -type f | head -50`" + ` for structure overview
- Identify: language, framework, build system, test framework, CI/CD

### 2. Gather user context
If no argument was provided, ask the user:
- What is this project? (1-2 sentence vision)
- What are you trying to accomplish? (current goals)
- Any constraints or preferences? (tech choices, timeline, team size)

### 3. Research the codebase
Spawn 3 parallel research agents using the Agent tool with subagent_type="researcher":

Agent(subagent_type="researcher", run_in_background=true, prompt="Analyze the tech stack: dependencies, versions, framework patterns. Write findings to .planning/research/stack.md")
Agent(subagent_type="researcher", run_in_background=true, prompt="Map the architecture: directory structure, entry points, data flow. Write findings to .planning/research/architecture.md")
Agent(subagent_type="researcher", run_in_background=true, prompt="Identify pitfalls: TODO/FIXME density, test coverage gaps, dependency issues. Write findings to .planning/research/pitfalls.md")

### 4. Create .planning/ directory
Write these files:

**PROJECT.md**:
- Project name and vision (from user or README)
- Current goals
- Detected tech stack with sources
- Constraints from user input
- Team context

**REQUIREMENTS.md**:
- v1 (Current Scope) with REQ-IDs derived from goals
- v2 (Future) with deferred ideas
- Out of Scope section

**ROADMAP.md**:
- Phases with scope, requirements mapping, success criteria, and status

**STATE.md** (with YAML frontmatter):
- milestone: v1.0, phase: 0, status: initialized, last_updated timestamp
- Current Position, Decisions Made, Quick Tasks Completed sections

**config.json**:
- mode: interactive, granularity: standard, model_profile: quality
- workflow toggles for research, plan_check, discuss, auto_advance

Write research outputs to ` + "`.planning/research/stack.md`" + `, ` + "`architecture.md`" + `, ` + "`pitfalls.md`" + `.
Create empty directories: ` + "`.planning/phases/`" + `, ` + "`.planning/quick/`" + `, ` + "`.planning/notes/`" + `.

### 5. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): initialize project planning\"`" + `

### 6. Report
Print summary: project type detected, phases planned, next step recommendation.
Suggest: "Run ` + "`/yaah:discuss 1`" + ` to clarify implementation decisions, or ` + "`/yaah:plan 1`" + ` to start planning Phase 1."

## Rules
- NEVER overwrite existing ` + "`.planning/`" + ` without explicit user confirmation
- NEVER guess the project vision — ask the user
- Tag all detected facts with their source (e.g., "Go 1.25 (from go.mod)")
`
}
