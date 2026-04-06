package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// AutonomousCommand runs the full workflow without human intervention via /yaah:autonomous.
type AutonomousCommand struct{}

// NewAutonomousCommand creates a new AutonomousCommand.
func NewAutonomousCommand() *AutonomousCommand { return &AutonomousCommand{} }

func (c *AutonomousCommand) Name() string        { return "yaah/autonomous" }
func (c *AutonomousCommand) Description() string {
	return "Run the full workflow autonomously for a phase"
}
func (c *AutonomousCommand) ArgumentHint() string { return "<phase-number> [--dry-run]" }
func (c *AutonomousCommand) Model() string        { return "" }
func (c *AutonomousCommand) AllowedTools() string { return "" }
func (c *AutonomousCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *AutonomousCommand) Content() string {
	return `# /yaah:autonomous — Autonomous Phase Execution

## When to use
When the user runs ` + "`/yaah:autonomous <phase-number> [--dry-run]`" + ` to execute a complete phase lifecycle without pausing for user input.

## Behavior
- Runs the full lifecycle for the specified phase in sequence, without human intervention at each step.
- If ` + "`--dry-run`" + ` is provided: describe what would happen at each step without actually executing anything.

## Lifecycle steps

### Step 1: Discuss — ` + "`/yaah:discuss {N}`" + `
- Run the discuss workflow for phase N.
- For every decision presented, auto-select the recommended default option.
- Record all auto-selected decisions in STATE.md under "Autonomous Decisions".

### Step 2: Plan — ` + "`/yaah:plan {N}`" + `
- Run the plan workflow for phase N.
- Do not pause for user review of the generated plan.

### Step 3: Execute — ` + "`/yaah:execute {N}`" + `
- Run the execute workflow for all waves in phase N.
- Do not pause between waves.

### Step 4: Verify — ` + "`/yaah:verify {N}`" + `
- Run the verify workflow for phase N.
- NEVER skip this step.

### Step 5: Evaluate result
- If verification passes: report success and print a completion summary.
- If verification fails: attempt one fix cycle (see Fix Cycle below).

## Fix cycle (one attempt only)
1. Analyze verification failures to identify gaps.
2. Re-plan only the failing areas (do not re-plan the entire phase).
3. Re-execute the affected tasks.
4. Re-verify.
5. If verification now passes: report success.
6. If verification still fails: STOP immediately and report all remaining failures to the user.

## Safeguards
- NEVER skip verification — it is mandatory.
- Stop after 2 fix cycles maximum (one autonomous, then report to user).
- Stop immediately if architectural changes are needed — these require human judgment.
- Record every autonomous decision in STATE.md with a timestamp and rationale.

## Dry-run mode
When ` + "`--dry-run`" + ` is specified:
- Print a detailed plan of what each step would do.
- Show which decisions would be auto-selected in the discuss step.
- Show the expected wave structure from the plan step.
- Do NOT modify any files or run any commands.

## Output format
After completion (or failure), print a summary:

` + "```" + `
# Autonomous Execution Summary — Phase {N}

## Steps Completed
- Discuss: {result}
- Plan: {result}
- Execute: {waves completed}
- Verify: {passed / failed}
- Fix cycle: {not needed / attempted — {result}}

## Autonomous Decisions Made
- D-01: <decision> — auto-selected: <option> (reason: recommended default)
- ...

## Outcome
{SUCCESS: phase {N} complete} or {STOPPED: {reason, remaining failures}}
` + "```" + `

## Rules
- NEVER make architectural decisions autonomously — stop and ask the user.
- NEVER skip verification, even if execution appears clean.
- NEVER run more than 2 fix cycles — escalate to the user instead.
- All autonomous decisions must be recorded in STATE.md before proceeding.
`
}
