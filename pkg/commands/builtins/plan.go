package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// PlanCommand implements structured phase planning via /yaah:plan.
type PlanCommand struct{}

// NewPlanCommand creates a new PlanCommand.
func NewPlanCommand() *PlanCommand { return &PlanCommand{} }

func (c *PlanCommand) Name() string { return "yaah/plan" }
func (c *PlanCommand) Description() string {
	return "Create a structured implementation plan for a project phase"
}
func (c *PlanCommand) ArgumentHint() string { return "<phase-number-or-description>" }
func (c *PlanCommand) Model() string        { return "" }
func (c *PlanCommand) AllowedTools() string { return "" }
func (c *PlanCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *PlanCommand) Content() string {
	return `# /yaah:plan — Structured Phase Planning

## When to use
When the user runs ` + "`/yaah:plan <phase>`" + ` to produce a wave-based implementation plan for a project phase.

## Goal
Goal-backward planning: define the end state first, derive artifacts, then decompose into tasks grouped into dependency waves.

## Steps

### 1. Load context
- Read all ` + "`.planning/`" + ` files: PROJECT.md, REQUIREMENTS.md, ROADMAP.md, STATE.md
- Read ` + "`.planning/phases/{NN}-{slug}/CONTEXT.md`" + ` if it exists (run ` + "`/yaah:discuss`" + ` first if missing)
- Note all D-XX decisions — they are binding constraints

### 2. Research codebase
Write ` + "`.planning/phases/{NN}-{slug}/RESEARCH.md`" + ` with:
- Every relevant file and its role, tagged with confidence:
  - ` + "`[VERIFIED]`" + ` — read and confirmed
  - ` + "`[CITED]`" + ` — referenced but not fully read
  - ` + "`[ASSUMED]`" + ` — inferred from naming or structure
- Integration points that tasks must respect
- Patterns already in use that tasks must follow

### 3. Goal-backward decomposition
1. **Define end state**: What does "done" look like for this phase? (concrete, verifiable)
2. **Identify artifacts**: What files must exist or change?
3. **Derive tasks**: What must happen to produce each artifact?
4. **Group into waves**: Assign wave numbers so that all dependencies within a wave are satisfied by prior waves

Wave assignment rule: ` + "`wave(task) = max(wave(deps)) + 1`" + `

### 4. Write PLAN.md files
Create one PLAN.md per logical group (max 3 tasks per plan, max 3 plans per phase).
Write to ` + "`.planning/phases/{NN}-{slug}/{plan-slug}/PLAN.md`" + `.

Each PLAN.md uses YAML frontmatter:
` + "```" + `yaml
---
phase: 1
plan: auth
wave: 1
depends_on: []
files_modified:
  - pkg/auth/handler.go
  - internal/cli/root.go
---
` + "```" + `

Each task section:
` + "```" + `markdown
## Task: <title>

**Files**: <comma-separated list>
**Action**: <imperative sentence — what to write or change>
**Verify**: <exact command that must exit 0>
**Done**: [ ]
` + "```" + `

### 5. Self-verify the plan
Before writing, check:
- Every REQ-ID in scope appears in at least one task
- Every Verify criterion is a concrete shell command (not "check that it works")
- No two tasks in the same wave modify the same file
- Dependency graph is acyclic
- Max 3 tasks per PLAN.md, max 3 PLAN.md files per phase

### 6. Update STATE.md and commit
Set ` + "`status: planned`" + `, update ` + "`last_updated`" + `.
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): create phase {N} plan\"`" + `

## Rules
- No vague tasks: "implement auth" is not a task; "write ` + "`pkg/auth/handler.go`" + ` implementing ` + "`Handler`" + ` interface with JWT validation" is
- Honor all CONTEXT.md decisions (D-XX) — never contradict them
- Tag all research findings with their confidence level
- Wave numbering must be globally consistent within a phase
`
}
