package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// RemovePhaseCommand removes a future phase from the roadmap via /yaah:remove-phase.
type RemovePhaseCommand struct{}

// NewRemovePhaseCommand creates a new RemovePhaseCommand.
func NewRemovePhaseCommand() *RemovePhaseCommand { return &RemovePhaseCommand{} }

func (c *RemovePhaseCommand) Name() string { return "yaah/remove-phase" }
func (c *RemovePhaseCommand) Description() string {
	return "Remove a future phase from the roadmap"
}
func (c *RemovePhaseCommand) ArgumentHint() string { return "<phase-number>" }
func (c *RemovePhaseCommand) Model() string        { return "" }
func (c *RemovePhaseCommand) AllowedTools() string { return "" }
func (c *RemovePhaseCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *RemovePhaseCommand) Content() string {
	return `# /yaah:remove-phase — Remove a Future Phase

## When to use
When the user runs ` + "`/yaah:remove-phase <phase-number>`" + ` to remove a planned (not yet started or in-progress) phase.

## Steps

### 1. Validate the phase
- Read ` + "`.planning/ROADMAP.md`" + ` and confirm the specified phase number exists.
- If the phase does not exist, report the error and stop.

### 2. Check phase status
- Read the phase status from ROADMAP.md.
- If the phase status is ` + "`COMPLETED`" + `, ` + "`VERIFIED`" + `, ` + "`EXECUTING`" + `, or ` + "`IN PROGRESS`" + `:
  - REFUSE to remove it.
  - Explain: "Phase {N} is {status} and cannot be removed. Only future phases (NOT STARTED or PLANNED) may be removed."
  - Stop.

### 3. Ask for confirmation
Before making any changes, print:
- The phase name and description being removed.
- The list of requirements (REQ-IDs) that will be affected.
- "Are you sure you want to remove phase {N} — {name}? This cannot be undone."
- Wait for explicit user confirmation before continuing.

### 4. Remove phase from ROADMAP.md
- Delete the phase entry from ` + "`.planning/ROADMAP.md`" + `.
- Renumber all subsequent phases to fill the gap (each shifts by -1).

### 5. Handle requirements
For each REQ-ID associated with the removed phase:
- Move them to an "Out of Scope" or "v2" section in ` + "`.planning/REQUIREMENTS.md`" + `.
- Add a note: "Removed with phase {N} on {date}."
- Do NOT delete REQ-IDs — they serve as a record of deferred scope.

### 6. Handle phase directory
- Check ` + "`.planning/phases/{NN}-{slug}/`" + ` for the removed phase.
- If the directory is empty: delete it silently.
- If the directory contains artifacts (CONTEXT.md, PLAN.md, etc.):
  - Ask the user: "Phase directory contains artifacts. Delete them? (yes/no)"
  - If yes: delete the directory and all contents.
  - If no: move directory to ` + "`.planning/archive/{NN}-{slug}/`" + `.

### 7. Rename subsequent phase directories
For each phase directory that was renumbered:
- Rename ` + "`.planning/phases/{old-NN}-{slug}/`" + ` to ` + "`.planning/phases/{new-NN}-{slug}/`" + `.

### 8. Update STATE.md
- Update total phase count.
- If the removed phase was after the current phase, no status change needed.
- Update ` + "`last_updated`" + ` timestamp.

### 9. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): remove phase {N}\"`" + `

## Rules
- NEVER remove a completed or in-progress phase
- ALWAYS ask for confirmation before deleting anything
- NEVER delete REQ-IDs — move them to Out of Scope or v2
- ALWAYS renumber subsequent phases to maintain sequential numbering
`
}
