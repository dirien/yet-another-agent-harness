package builtins

// ProgressCommand shows detailed project progress via /yaah:progress.
// It is lightweight and does not implement CommandWithAdvanced.
type ProgressCommand struct{}

// NewProgressCommand creates a new ProgressCommand.
func NewProgressCommand() *ProgressCommand { return &ProgressCommand{} }

func (c *ProgressCommand) Name() string        { return "yaah/progress" }
func (c *ProgressCommand) Description() string { return "Show detailed project progress with metrics" }
func (c *ProgressCommand) ArgumentHint() string { return "" }
func (c *ProgressCommand) Model() string        { return "" }
func (c *ProgressCommand) AllowedTools() string { return "Read, Glob" }

func (c *ProgressCommand) Content() string {
	return `# /yaah:progress — Project Progress Report

## When to use
When the user runs ` + "`/yaah:progress`" + ` to get a detailed view of where the project stands.

## Steps

### 1. Read planning documents
- Read ` + "`.planning/STATE.md`" + ` for current phase, status, and milestone.
- Read ` + "`.planning/ROADMAP.md`" + ` for all phases and their statuses.
- Read ` + "`.planning/REQUIREMENTS.md`" + ` to count total and satisfied REQ-IDs.

### 2. Scan phase directories
- Glob ` + "`.planning/phases/*/`" + ` to find all phase directories.
- For each phase directory, check which artifacts exist:
  - ` + "`PLAN.md`" + ` → phase has been planned
  - ` + "`SUMMARY.md`" + ` → phase has been executed (count plan files for X/N plans done)
  - ` + "`VERIFICATION.md`" + ` → phase has been verified
  - ` + "`SUMMARY.md`" + ` timestamps → used for timing metrics if available

### 3. Count requirement coverage
- Extract all REQ-IDs from ` + "`.planning/REQUIREMENTS.md`" + ` in v1 (Current Scope).
- Glob all ` + "`VERIFICATION.md`" + ` files and check which REQ-IDs are marked SATISFIED.
- Calculate: {satisfied REQ-IDs} / {total v1 REQ-IDs}.

### 4. Build progress bar per phase
For each phase, render a 12-block progress bar based on status:
- NOT STARTED: ` + "`░░░░░░░░░░░░`" + `
- PLANNED (has PLAN.md): ` + "`░░░░░░░░░░░░ PLANNED`" + `
- EXECUTING (has SUMMARY.md but not VERIFICATION.md): partial fill based on plans complete
- VERIFIED (has VERIFICATION.md): ` + "`████████████ VERIFIED ✓`" + `

### 5. Print the progress report
Output in this format:
` + "```" + `
Milestone: {version from STATE.md}
Progress: {completed phases}/{total phases} ({percentage}%)

Phase 1: {name} ████████████ VERIFIED ✓
Phase 2: {name} ██████░░░░░░ EXECUTING ({X}/{N} plans)
Phase 3: {name} ░░░░░░░░░░░░ PLANNED
Phase 4: {name} ░░░░░░░░░░░░ NOT STARTED

Quick tasks: {count from .planning/quick/} completed
Current: Phase {N}, Plan {X} of {total}
Next: /yaah:execute {N} (continue Wave {W})
` + "```" + `

### 6. Include timing metrics
If any ` + "`SUMMARY.md`" + ` files contain timestamp headers or "Completed:" fields:
- Calculate elapsed time per phase.
- Show average time per phase if 2+ phases are complete.
- Show estimated remaining time based on average (if applicable).

### 7. Show requirement coverage
Append at the end:
` + "```" + `
Requirement coverage: {satisfied}/{total} REQ-IDs satisfied
` + "```" + `

## Rules
- Never modify any files — this is a read-only reporting command
- If ` + "`.planning/`" + ` does not exist, recommend ` + "`/yaah:init`" + ` and stop
- Always show the Next recommendation at the end of the report
- Percentage is calculated as: (completed phases / total phases) * 100, rounded to nearest integer
`
}
