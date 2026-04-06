package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// CompleteMilestoneCommand archives the current milestone via /yaah:complete-milestone.
type CompleteMilestoneCommand struct{}

// NewCompleteMilestoneCommand creates a new CompleteMilestoneCommand.
func NewCompleteMilestoneCommand() *CompleteMilestoneCommand { return &CompleteMilestoneCommand{} }

func (c *CompleteMilestoneCommand) Name() string { return "yaah/complete-milestone" }
func (c *CompleteMilestoneCommand) Description() string {
	return "Archive current milestone, tag release, generate changelog"
}
func (c *CompleteMilestoneCommand) ArgumentHint() string { return "[version-tag]" }
func (c *CompleteMilestoneCommand) Model() string        { return "" }
func (c *CompleteMilestoneCommand) AllowedTools() string { return "" }
func (c *CompleteMilestoneCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *CompleteMilestoneCommand) Content() string {
	return `# /yaah:complete-milestone — Archive Milestone and Tag Release

## When to use
When the user runs ` + "`/yaah:complete-milestone`" + ` or ` + "`/yaah:complete-milestone <version>`" + ` after all phases are shipped.

## Prerequisites
- All phases in ROADMAP.md must have status "shipped"
- If any phase is not shipped, STOP and list the incomplete phases

## Steps

### 1. Verify all phases complete
- Read ` + "`.planning/ROADMAP.md`" + ` and check every phase for status "shipped"
- If any phase is missing "shipped" status, print the list and suggest ` + "`/yaah:ship`" + `

### 2. Determine version tag
- If an argument was provided, use it as the version tag (e.g., ` + "`v1.0.0`" + `)
- Otherwise, read the milestone from ` + "`.planning/STATE.md`" + ` frontmatter
- Confirm the tag with the user before proceeding

### 3. Generate CHANGELOG.md
Synthesize from all phase artifacts:
- Read every ` + "`.planning/phases/*/SUMMARY.md`" + ` in phase order
- Read every ` + "`.planning/phases/*/VERIFICATION.md`" + ` for test results
- Write ` + "`CHANGELOG.md`" + ` at the repository root with sections:
  ## {version} — {date}
  ### Summary
  {1-2 sentences describing the milestone goal}
  ### Changes by Phase
  {per-phase list of key changes from SUMMARYs}
  ### Verification
  {build and test results from VERIFICATIONs}

### 4. Create git tag
- Run: ` + "`git tag -a {version} -m \"Release {version}\"`" + `
- Do NOT push the tag unless the user explicitly asks

### 5. Archive .planning/ state
- Create directory: ` + "`.planning/archive/{version}/`" + `
- Copy all files from ` + "`.planning/`" + ` (except ` + "`archive/`" + ` itself) into ` + "`.planning/archive/{version}/`" + `
- This preserves the full planning history for the milestone

### 6. Update STATE.md
- Set milestone status to "complete"
- Record completion timestamp
- Reset phase to 0 in preparation for the next milestone

### 7. Commit
- Stage and commit: ` + "`git add CHANGELOG.md .planning/ && git commit -m \"chore: complete milestone {version}\"`" + `

### 8. Report
- Print the version tag created
- Print the CHANGELOG.md path
- Suggest: "Run ` + "`/yaah:new-milestone <next-version>`" + ` to start the next cycle"

## Rules
- NEVER tag an incomplete milestone — all phases must be shipped
- NEVER push the tag automatically — let the user decide when to push
- ALWAYS archive before resetting STATE.md to preserve history
`
}
