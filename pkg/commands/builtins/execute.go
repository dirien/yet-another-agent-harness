package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ExecuteCommand implements wave-based plan execution via /yaah:execute.
type ExecuteCommand struct{}

// NewExecuteCommand creates a new ExecuteCommand.
func NewExecuteCommand() *ExecuteCommand { return &ExecuteCommand{} }

func (c *ExecuteCommand) Name() string         { return "yaah/execute" }
func (c *ExecuteCommand) Description() string  { return "Execute implementation plans wave by wave" }
func (c *ExecuteCommand) ArgumentHint() string { return "<phase-number> [--wave N] [--interactive]" }
func (c *ExecuteCommand) Model() string        { return "" }
func (c *ExecuteCommand) AllowedTools() string { return "" }
func (c *ExecuteCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ExecuteCommand) Content() string {
	return `# /yaah:execute — Wave-Based Plan Execution

## When to use
When the user runs ` + "`/yaah:execute <phase>`" + ` to execute plans produced by ` + "`/yaah:plan`" + `.

## Goal
Execute all plans for a phase in wave order, spawning parallel subagents where safe, and producing a SUMMARY.md for each completed plan.

## Steps

### 1. Discover plans and handle resumption
- Glob ` + "`.planning/phases/{NN}-{slug}/*/PLAN.md`" + ` to find all plans
- For each plan, check whether a ` + "`SUMMARY.md`" + ` already exists alongside it
- Plans with an existing ` + "`SUMMARY.md`" + ` are already done — skip them (resumption)
- Group remaining plans by their ` + "`wave:`" + ` frontmatter field

### 2. Pre-execution summary
Print a table to the user:
- Wave number, plan name, files modified, dependency status
- Ask for confirmation before proceeding (unless ` + "`--interactive`" + ` is not set and the user already confirmed)

### 3. Execute wave by wave
For each wave (ascending order):

#### a. Pre-wave dependency check
- All plans from prior waves must have ` + "`SUMMARY.md`" + ` present
- If any prior-wave ` + "`SUMMARY.md`" + ` is missing, halt and report which plan is incomplete

#### b. Intra-wave conflict check
- Collect all ` + "`files_modified`" + ` from plans in this wave
- If any file appears in more than one plan, force those plans to run sequentially (in alphabetical order by plan slug) rather than in parallel

#### c. Spawn executor subagents
For plans with no file conflicts, spawn in parallel:
` + "```" + `
Agent(subagent_type="executor", run_in_background=true, prompt=<plan content>)
` + "```" + `

Each subagent receives:
- The full PLAN.md content
- The CONTEXT.md decisions (D-XX) as binding constraints
- Instruction to write atomic commits per task and produce ` + "`SUMMARY.md`" + ` on completion

#### d. Wait and collect
- Wait for all subagents in the wave to complete
- Verify each ` + "`SUMMARY.md`" + ` exists
- If a subagent failed (no ` + "`SUMMARY.md`" + `), report the failure and halt

### 4. Post-execution update
- Update ` + "`.planning/STATE.md`" + `: set ` + "`status: executed`" + `, update ` + "`last_updated`" + `
- Update ` + "`.planning/ROADMAP.md`" + `: mark phase as complete
- Run: ` + "`git add .planning/ && git commit -m \"docs(planning): mark phase {N} executed\"`" + `

### 5. Suggest next step
Print: "Phase {N} execution complete. Run ` + "`/yaah:verify {N}`" + ` to validate against requirements."

## Rules
- NEVER execute plans out of wave order — dependencies are not optional
- Each task within a plan produces an atomic commit before moving to the next task
- NEVER spawn subagents in parallel when they share files_modified — force sequential
- Use ` + "`--interactive`" + ` flag to skip subagent spawning and execute directly in the current agent
- The ` + "`--wave N`" + ` flag limits execution to a single wave number
`
}
