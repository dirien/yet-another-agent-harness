package builtins

// NextCommand implements next-step recommendation via /yaah:next.
// It is lightweight and does not implement CommandWithAdvanced.
type NextCommand struct{}

// NewNextCommand creates a new NextCommand.
func NewNextCommand() *NextCommand { return &NextCommand{} }

func (c *NextCommand) Name() string         { return "yaah/next" }
func (c *NextCommand) Description() string  { return "Auto-detect and recommend the next workflow step" }
func (c *NextCommand) ArgumentHint() string { return "" }
func (c *NextCommand) Model() string        { return "" }
func (c *NextCommand) AllowedTools() string { return "Read, Glob" }

func (c *NextCommand) Content() string {
	return `# /yaah:next — Next Step Recommendation

## When to use
When the user runs ` + "`/yaah:next`" + ` and wants to know what to do next in the workflow.

## Steps

### 1. Check for .planning/
- Glob for ` + "`.planning/`" + `
- If it does not exist: recommend ` + "`/yaah:init`" + ` and stop

### 2. Read STATE.md and ROADMAP.md
- Read ` + "`.planning/STATE.md`" + ` for current phase and status
- Read ` + "`.planning/ROADMAP.md`" + ` for phase list and completion markers

### 3. Scan phases/ for artifacts
- Glob ` + "`.planning/phases/*/`" + ` to find phase directories
- For the current phase, check which artifacts exist: CONTEXT.md, PLAN.md, SUMMARY.md, VERIFICATION.md

### 4. Determine state and recommend
| STATE.md status | Artifacts present | Recommendation |
|-----------------|-------------------|----------------|
| initialized     | none              | ` + "`/yaah:discuss {N}`" + ` |
| initialized     | CONTEXT.md        | ` + "`/yaah:plan {N}`" + ` |
| discussed       | CONTEXT.md        | ` + "`/yaah:plan {N}`" + ` |
| planned         | PLAN.md           | ` + "`/yaah:execute {N}`" + ` |
| executed        | SUMMARY.md        | ` + "`/yaah:verify {N}`" + ` |
| verified        | VERIFICATION.md   | ` + "`/yaah:plan {N+1}`" + ` or ` + "`/yaah:docs`" + ` if last phase |
| failed          | VERIFICATION.md   | Re-read VERIFICATION.md and fix gaps |

### 5. Print exactly one recommendation
Output a single line with the exact command to run next. No explanation beyond one sentence of context.

## Rules
- ONE recommendation only — do not list alternatives
- Output must be under 5 lines total
- If state is ambiguous, default to the earliest incomplete step
`
}
