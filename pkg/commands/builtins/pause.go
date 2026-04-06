package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// PauseCommand saves session state for later resumption via /yaah:pause.
type PauseCommand struct{}

// NewPauseCommand creates a new PauseCommand.
func NewPauseCommand() *PauseCommand { return &PauseCommand{} }

func (c *PauseCommand) Name() string        { return "yaah/pause" }
func (c *PauseCommand) Description() string { return "Save current session state for later resumption" }
func (c *PauseCommand) ArgumentHint() string { return "" }
func (c *PauseCommand) Model() string        { return "" }
func (c *PauseCommand) AllowedTools() string { return "" }
func (c *PauseCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *PauseCommand) Content() string {
	return `# /yaah:pause — Save Session State

## When to use
When the user runs ` + "`/yaah:pause`" + ` to checkpoint the current session before stopping work.

## Steps

### 1. Collect current context
- Read ` + "`.planning/STATE.md`" + ` — note current phase, status, and last action
- Read the active ` + "`PLAN.md`" + ` for the current phase to identify the in-progress task
- Run ` + "`git status`" + ` to list uncommitted changes
- Run ` + "`git diff --stat`" + ` to summarize what has changed but not been committed

### 2. Write .planning/HANDOFF.md
Create ` + "`.planning/HANDOFF.md`" + ` with the following structure:

` + "```" + `markdown
# Session Handoff
Paused: {timestamp}

## Current State
- Phase: {N} — {phase name}
- Status: {status from STATE.md}
- Last completed task: {task description}

## In Progress
- Active plan: {path to PLAN.md}
- Current task: {task description from PLAN.md}
- Progress: {what has been done so far}

## Uncommitted Changes
{output of git diff --stat, or "none" if clean}

## Next Steps
When resuming, do the following in order:
1. {first concrete action}
2. {second concrete action}
3. Run /yaah:resume to restore full context

## Notes
{any important context the next session needs to know}
` + "```" + `

### 3. Commit the handoff file
- Stage and commit: ` + "`git add .planning/HANDOFF.md && git commit -m \"docs(planning): pause session\"`" + `

### 4. Report
- Confirm HANDOFF.md was written
- Print the "Next Steps" section so the user sees what to do when they return
- Suggest: "Run ` + "`/yaah:resume`" + ` to restore this session"

## Rules
- NEVER discard uncommitted work — summarize it in HANDOFF.md so nothing is lost
- ALWAYS commit HANDOFF.md before reporting done
- If .planning/ does not exist, report that there is nothing to pause
`
}
