package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// CleanupCommand removes temporary planning artifacts via /yaah:cleanup.
type CleanupCommand struct{}

// NewCleanupCommand creates a new CleanupCommand.
func NewCleanupCommand() *CleanupCommand { return &CleanupCommand{} }

func (c *CleanupCommand) Name() string { return "yaah/cleanup" }
func (c *CleanupCommand) Description() string {
	return "Clean up temporary planning artifacts and state"
}
func (c *CleanupCommand) ArgumentHint() string { return "[--dry-run]" }
func (c *CleanupCommand) Model() string        { return "" }
func (c *CleanupCommand) AllowedTools() string { return "" }
func (c *CleanupCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *CleanupCommand) Content() string {
	return `# /yaah:cleanup — Planning Artifact Cleanup

## When to use
When the user runs ` + "`/yaah:cleanup`" + ` to remove stale or empty planning artifacts from ` + "`.planning/`" + `.

## Steps

### 1. Parse arguments
- Detect ` + "`--dry-run`" + ` flag in the argument string
- In dry-run mode: list what would be removed, make no changes

### 2. Scan .planning/ for cleanup candidates
Identify the following categories of stale artifacts:

**Empty phase directories**
- Glob ` + "`.planning/phases/*/`" + `
- A phase directory is empty if it contains no files (no CONTEXT.md, PLAN.md, SUMMARY.md, VERIFICATION.md, SECURITY.md, or RESEARCH.md)

**Stale HANDOFF.md**
- Glob ` + "`.planning/phases/*/HANDOFF.md`" + `
- Mark as stale if the file's last-modified date is older than 7 days
- Read each HANDOFF.md to include its title in the report

**Old notes**
- Glob ` + "`.planning/notes/*.md`" + `
- Mark as old if the filename date (` + "`YYYY-MM-DD`" + `) is older than 30 days from today

**Orphaned quick task records**
- Glob ` + "`.planning/quick/*.md`" + `
- A record is orphaned if its slug does not correspond to any phase directory or active ROADMAP.md phase

### 3. Report findings
Print a summary table:
` + "```" + `
Cleanup candidates:
  [empty phase]    .planning/phases/03-auth-provider/
  [stale handoff]  .planning/phases/01-core/HANDOFF.md  (14 days old)
  [old note]       .planning/notes/2026-03-01.md  (35 days old)
  [orphaned quick] .planning/quick/2026-02-15-fix-typo.md
` + "```" + `

If nothing is found, print "Nothing to clean up." and stop.

### 4. If --dry-run
Stop here. Print: "Dry run complete. Re-run without ` + "`--dry-run`" + ` to apply."

### 5. Ask for confirmation
Print the list and ask: "Remove these N items? (yes/no)"
If the user says no, abort without changes.

### 6. Remove confirmed items
Delete each listed file or directory.
Print each removal as it completes: "Removed: {path}"

### 7. Commit if changes were made
Run: ` + "`git add .planning/ && git commit -m \"chore(planning): clean up stale artifacts\"`" + `

## Rules
- NEVER remove STATE.md, PROJECT.md, REQUIREMENTS.md, or ROADMAP.md
- NEVER remove a phase directory that contains any artifact files
- Always ask for confirmation before deleting — dry-run is safe, removal is destructive
- If ` + "`.planning/`" + ` does not exist, print "No .planning/ directory found." and stop
`
}
