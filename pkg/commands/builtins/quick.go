package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// QuickCommand implements lightweight task execution via /yaah:quick.
type QuickCommand struct{}

// NewQuickCommand creates a new QuickCommand.
func NewQuickCommand() *QuickCommand { return &QuickCommand{} }

func (c *QuickCommand) Name() string        { return "yaah/quick" }
func (c *QuickCommand) Description() string { return "Execute a task without full planning overhead" }
func (c *QuickCommand) ArgumentHint() string {
	return "<task-description> [--discuss] [--research] [--validate]"
}
func (c *QuickCommand) Model() string        { return "" }
func (c *QuickCommand) AllowedTools() string { return "" }
func (c *QuickCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *QuickCommand) Content() string {
	return `# /yaah:quick — Lightweight Task Execution

## When to use
When the user runs ` + "`/yaah:quick <task>`" + ` to complete a self-contained task without the full planning workflow.

## Goal
Just do it. Default behavior is immediate execution. Flags compose to add optional steps.

## Steps

### 1. Parse arguments and flags
- Extract the task description (everything before flags)
- Detect flags: ` + "`--discuss`" + `, ` + "`--research`" + `, ` + "`--validate`" + `
- Flags may appear anywhere in the argument string

### 2. Optional: --discuss
If ` + "`--discuss`" + ` is present:
- Analyze the task and identify the 1-2 key implementation decisions
- Present each decision with 2-3 concrete options and tradeoffs
- Confirm the user's approach before proceeding
- Do NOT ask what to build — only how

### 3. Optional: --research
If ` + "`--research`" + ` is present:
- Grep the codebase for patterns similar to the task
- Read 2-3 most relevant files to understand existing conventions
- Note: file paths to touch, patterns to follow, imports to use
- Summarize findings in 3-5 bullet points before executing

### 4. Execute the task
Perform the task directly:
- Make the necessary code changes
- Write atomic commits as logical units of work complete
- Use Conventional Commits format: ` + "`feat:`" + `, ` + "`fix:`" + `, ` + "`refactor:`" + `, etc.
- Do not create planning files — this is a quick task

### 5. Optional: --validate
If ` + "`--validate`" + ` is present:
- Auto-detect build system and run build (see ` + "`/yaah:verify`" + ` for detection logic)
- Auto-detect test runner and run tests
- Report pass/fail with output
- If validation fails, fix and re-run before completing

### 6. Record the task
Write a brief record to ` + "`.planning/quick/{YYYY-MM-DD}-{slug}.md`" + ` (create ` + "`.planning/quick/`" + ` if needed):
` + "```" + `markdown
# {task description}
Date: {date}
Flags: {flags used}
Files changed: {list}
Outcome: {one sentence}
` + "```" + `

## Rules
- Default behavior (no flags) is: execute immediately, no discussion, no research
- Flags are additive and compose: ` + "`--discuss --research --validate`" + ` all apply together
- This command is independent of the phase workflow — no ROADMAP.md updates
- Atomic commits are always required regardless of flags
- If the task touches more than 5 files, suggest using ` + "`/yaah:plan`" + ` instead
`
}
