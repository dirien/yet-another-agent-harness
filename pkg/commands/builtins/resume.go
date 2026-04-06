package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ResumeCommand restores from a paused session via /yaah:resume.
type ResumeCommand struct{}

// NewResumeCommand creates a new ResumeCommand.
func NewResumeCommand() *ResumeCommand { return &ResumeCommand{} }

func (c *ResumeCommand) Name() string        { return "yaah/resume" }
func (c *ResumeCommand) Description() string { return "Resume work from a previous session handoff" }
func (c *ResumeCommand) ArgumentHint() string { return "" }
func (c *ResumeCommand) Model() string        { return "" }
func (c *ResumeCommand) AllowedTools() string { return "" }
func (c *ResumeCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ResumeCommand) Content() string {
	return `# /yaah:resume — Resume from Paused Session

## When to use
When the user runs ` + "`/yaah:resume`" + ` after a previous ` + "`/yaah:pause`" + `.

## Prerequisites
- ` + "`.planning/HANDOFF.md`" + ` must exist

## Steps

### 1. Load handoff context
- Read ` + "`.planning/HANDOFF.md`" + `
- If it does not exist, report "No paused session found" and suggest ` + "`/yaah:next`" + `

### 2. Verify git state
- Run ` + "`git status`" + ` to check for unexpected uncommitted changes
- Run ` + "`git log --oneline -5`" + ` to confirm branch and recent commits
- If uncommitted changes exist that were NOT noted in HANDOFF.md, warn the user before proceeding

### 3. Restore planning context
- Read ` + "`.planning/STATE.md`" + ` for current phase and status
- Read ` + "`.planning/ROADMAP.md`" + ` for phase overview and completion markers
- Read the active ` + "`PLAN.md`" + ` referenced in HANDOFF.md

### 4. Report restored state
Print a summary:
- Phase and status
- What was in progress when paused
- Uncommitted changes (if any)
- The exact next action to take

### 5. Suggest the next action
Based on the "Next Steps" section of HANDOFF.md, state the single most immediate action.
Use the same format as ` + "`/yaah:next`" + ` — one recommendation, one sentence of context.

### 6. Delete HANDOFF.md
After successful restore, remove the handoff file and commit:
- ` + "`git rm .planning/HANDOFF.md && git commit -m \"docs(planning): resume session\"`" + `

## Rules
- NEVER resume into an ambiguous state — verify git matches what HANDOFF.md described
- ALWAYS delete HANDOFF.md after a successful resume to avoid stale state
- If HANDOFF.md is corrupt or unreadable, fall back to ` + "`/yaah:next`" + ` behavior
`
}
