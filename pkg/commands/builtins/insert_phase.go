package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// InsertPhaseCommand inserts a phase between existing phases via /yaah:insert-phase.
type InsertPhaseCommand struct{}

// NewInsertPhaseCommand creates a new InsertPhaseCommand.
func NewInsertPhaseCommand() *InsertPhaseCommand { return &InsertPhaseCommand{} }

func (c *InsertPhaseCommand) Name() string { return "yaah/insert-phase" }
func (c *InsertPhaseCommand) Description() string {
	return "Insert an urgent phase between existing phases"
}
func (c *InsertPhaseCommand) ArgumentHint() string { return "<after-phase-number> <description>" }
func (c *InsertPhaseCommand) Model() string        { return "" }
func (c *InsertPhaseCommand) AllowedTools() string { return "" }
func (c *InsertPhaseCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *InsertPhaseCommand) Content() string {
	return `# /yaah:insert-phase — Insert Phase Between Existing Phases

## When to use
When the user runs ` + "`/yaah:insert-phase <after-phase-number> <description>`" + ` to insert an urgent phase after an existing one.

## WARNING
This command renumbers all subsequent phases. Inform the user before proceeding.

## Steps

### 1. Read ROADMAP.md
- Read ` + "`.planning/ROADMAP.md`" + ` to understand the full current phase structure.
- Identify the phase after which the new phase will be inserted (the "anchor" phase).
- List all phases that will be renumbered (anchor+1 onward).

### 2. Warn the user
Before making any changes, print a warning:
- "Inserting after phase {N} will renumber phases {N+1} through {last} by +1."
- List all affected phase directories that will be renamed.
- Ask the user to confirm before continuing.

### 3. Insert new phase in ROADMAP.md
- Insert the new phase entry immediately after the anchor phase.
- Assign it the number (anchor + 1).
- Renumber ALL subsequent phases (each shifts by +1).
- Set new phase status to ` + "`NOT STARTED`" + `.

### 4. Rename phase directories
For each phase directory that was renumbered:
- Rename ` + "`.planning/phases/{old-NN}-{slug}/`" + ` to ` + "`.planning/phases/{new-NN}-{slug}/`" + `.
- Preserve all existing contents inside the directory.

### 5. Update depends_on references
- Grep all ` + "`PLAN.md`" + ` files under ` + "`.planning/phases/`" + ` for ` + "`depends_on`" + ` fields.
- Update any phase number references that shifted due to renumbering.

### 6. Update REQUIREMENTS.md
- Assign new REQ-IDs (continuing from the highest existing REQ-ID) for the new phase.
- Add them to the v1 (Current Scope) section of ` + "`.planning/REQUIREMENTS.md`" + `.

### 7. Update STATE.md
- If ` + "`phase`" + ` in STATE.md frontmatter refers to a phase that was renumbered, increment it.
- Update ` + "`last_updated`" + ` timestamp.
- Update total phase count if tracked.

### 8. Create new phase directory
Create ` + "`.planning/phases/{NN}-{slug}/`" + ` for the new phase where:
- ` + "`{NN}`" + ` is the zero-padded inserted phase number
- ` + "`{slug}`" + ` is a lowercase hyphenated version of the new phase name

### 9. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): insert phase {N} — {name}\"`" + `

## Rules
- NEVER insert a phase before a completed or in-progress phase
- ALWAYS warn the user about renumbering impact before making changes
- ALL subsequent phase directories must be renamed — no orphaned numbers
- ALL depends_on references must be updated — stale references break the workflow
`
}
