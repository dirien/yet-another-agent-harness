package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// AddPhaseCommand appends a new phase to the roadmap via /yaah:add-phase.
type AddPhaseCommand struct{}

// NewAddPhaseCommand creates a new AddPhaseCommand.
func NewAddPhaseCommand() *AddPhaseCommand { return &AddPhaseCommand{} }

func (c *AddPhaseCommand) Name() string        { return "yaah/add-phase" }
func (c *AddPhaseCommand) Description() string { return "Add a new phase to the end of the roadmap" }
func (c *AddPhaseCommand) ArgumentHint() string { return "<phase-description>" }
func (c *AddPhaseCommand) Model() string        { return "" }
func (c *AddPhaseCommand) AllowedTools() string { return "" }
func (c *AddPhaseCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *AddPhaseCommand) Content() string {
	return `# /yaah:add-phase — Add New Roadmap Phase

## When to use
When the user runs ` + "`/yaah:add-phase <phase-description>`" + ` to append a new phase to the end of the roadmap.

## Steps

### 1. Read ROADMAP.md
- Read ` + "`.planning/ROADMAP.md`" + ` to understand the current phase structure.
- Count existing phases and determine the next phase number (N = last phase number + 1).

### 2. Gather phase details
If the argument does not provide enough information, ask the user:
- What is the scope of this phase? (what will be built or changed)
- What are the requirements and acceptance criteria?
- What are the success criteria that mark this phase complete?

Do NOT ask if the argument already answers all three questions.

### 3. Append to ROADMAP.md
Add the new phase entry at the end of ` + "`.planning/ROADMAP.md`" + ` with:
- Sequential phase number (N)
- Name derived from the description
- Scope, requirements mapping, success criteria, and status: ` + "`NOT STARTED`"  + `

### 4. Update REQUIREMENTS.md
If the new phase introduces new functional requirements:
- Assign new REQ-IDs (continuing from the highest existing REQ-ID)
- Add them to the v1 (Current Scope) section of ` + "`.planning/REQUIREMENTS.md`" + `

### 5. Update STATE.md
- Increment the total phase count in ` + "`.planning/STATE.md`" + ` frontmatter if that field is tracked.
- Update ` + "`last_updated`" + ` timestamp.

### 6. Create phase directory
Create the empty directory ` + "`.planning/phases/{NN}-{slug}/`" + ` where:
- ` + "`{NN}`" + ` is the zero-padded phase number (e.g., ` + "`03`" + `)
- ` + "`{slug}`" + ` is a lowercase hyphenated version of the phase name

### 7. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): add phase {N} — {name}\"`" + `

### 8. Suggest next steps
Print one of:
- ` + "`/yaah:discuss {N}`" + ` — to clarify implementation decisions for the new phase
- ` + "`/yaah:plan {N}`" + ` — to start planning the new phase immediately

## Rules
- NEVER renumber existing phases — only append at the end
- NEVER add a phase without at least one success criterion
- Phase directory must be created even if empty to reserve the slot
`
}
