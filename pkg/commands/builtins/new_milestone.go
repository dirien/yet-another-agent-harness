package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// NewMilestoneCommand starts a new version cycle via /yaah:new-milestone.
type NewMilestoneCommand struct{}

// NewNewMilestoneCommand creates a new NewMilestoneCommand.
func NewNewMilestoneCommand() *NewMilestoneCommand { return &NewMilestoneCommand{} }

func (c *NewMilestoneCommand) Name() string        { return "yaah/new-milestone" }
func (c *NewMilestoneCommand) Description() string { return "Start a new version cycle with fresh goals" }
func (c *NewMilestoneCommand) ArgumentHint() string { return "<version>" }
func (c *NewMilestoneCommand) Model() string        { return "" }
func (c *NewMilestoneCommand) AllowedTools() string { return "" }
func (c *NewMilestoneCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *NewMilestoneCommand) Content() string {
	return `# /yaah:new-milestone — Start New Version Cycle

## When to use
When the user runs ` + "`/yaah:new-milestone <version>`" + ` to begin planning the next release cycle.

## Prerequisites
- A version argument is required (e.g., ` + "`v2.0.0`" + `)
- The previous milestone should be complete (STATE.md status "complete"), but this is not enforced

## Steps

### 1. Load previous milestone context
- Read ` + "`.planning/PROJECT.md`" + ` for project vision and goals
- Read ` + "`.planning/REQUIREMENTS.md`" + ` for v1 (current) and v2 (future) items
- Read ` + "`.planning/ROADMAP.md`" + ` for completed phases and any deferred work
- Note all v2 items from REQUIREMENTS.md — they become candidates for this milestone

### 2. Gather new goals from the user
Ask the user:
- What is the primary goal for {version}?
- Which v2 items from the previous milestone are in scope now?
- Are there new requirements not previously captured?
- Any constraints or changes (timeline, team, tech)?

Do NOT proceed to write files until goals are confirmed.

### 3. Reset STATE.md
Update ` + "`.planning/STATE.md`" + ` frontmatter:
- milestone: {version}
- phase: 0
- status: initialized
- last_updated: {timestamp}
- Clear "Decisions Made" and "Quick Tasks Completed" sections
- Preserve "Current Position" as a historical note

### 4. Create fresh REQUIREMENTS.md
- v1 (Current Scope): items confirmed by the user for this milestone, with new REQ-IDs
- v2 (Future): items deferred to a later milestone
- Out of Scope: explicitly excluded items
- Carry forward any unfinished v2 items from the previous milestone as new v1 candidates

### 5. Create new ROADMAP.md
- Define phases for the new milestone based on user goals and v1 requirements
- Each phase: name, scope, requirements mapping, success criteria, status: planned
- Reference the previous milestone's archive for historical context

### 6. Preserve research/
- Keep ` + "`.planning/research/`" + ` from the previous milestone — it remains valid context
- Add a note in PROJECT.md: "Research from {previous-version} preserved in research/"

### 7. Commit
- Stage and commit: ` + "`git add .planning/ && git commit -m \"docs(planning): start milestone {version}\"`" + `

### 8. Report
- Print the new milestone version
- Print the phases defined
- Suggest: "Run ` + "`/yaah:discuss 1`" + ` to clarify Phase 1 decisions"

## Rules
- NEVER start a new milestone without a version argument
- NEVER discard previous v2 items — they must appear in the new REQUIREMENTS.md
- ALWAYS preserve research/ — it represents accumulated knowledge
- Ask about goals before writing files
`
}
