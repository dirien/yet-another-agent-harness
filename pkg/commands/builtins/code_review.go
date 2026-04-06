package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// CodeReviewCommand runs structured code review via /yaah:review.
type CodeReviewCommand struct{}

// NewCodeReviewCommand creates a new CodeReviewCommand.
func NewCodeReviewCommand() *CodeReviewCommand { return &CodeReviewCommand{} }

func (c *CodeReviewCommand) Name() string        { return "yaah/review" }
func (c *CodeReviewCommand) Description() string { return "Structured code review of phase implementation" }
func (c *CodeReviewCommand) ArgumentHint() string { return "[phase-number]" }
func (c *CodeReviewCommand) Model() string        { return "" }
func (c *CodeReviewCommand) AllowedTools() string { return "" }
func (c *CodeReviewCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *CodeReviewCommand) Content() string {
	return `# /yaah:review — Structured Code Review

## When to use
When the user runs ` + "`/yaah:review [phase-number]`" + ` to perform a structured code review of a phase implementation.
If no phase number is given, use the current phase from ` + "`.planning/STATE.md`" + `.

## Steps

### 1. Load phase context
- Read ` + "`.planning/STATE.md`" + ` to determine the target phase if not specified.
- Read the ` + "`PLAN.md`" + ` for the target phase to understand planned intent.
- Read the ` + "`SUMMARY.md`" + ` for the target phase to identify files modified during execution.
- If neither PLAN.md nor SUMMARY.md exists for the phase, report an error and stop.

### 2. Identify files to review
Collect the list of modified files from (in order of preference):
1. ` + "`files_modified`" + ` frontmatter field in SUMMARY.md
2. ` + "`files_modified`" + ` frontmatter field in PLAN.md
3. Output of ` + "`git diff --name-only HEAD~1`" + ` as fallback

### 3. Review each file
For each file in the list, run a review using Agent(subagent_type="reviewer") with these five dimensions:

**Correctness**
- Does the implementation match what PLAN.md specified?
- Are edge cases handled (nil inputs, empty slices, boundary values)?
- Are error return paths complete and correct?

**Security**
- Injection risks (SQL, command, path traversal)?
- Secrets or credentials hardcoded or logged?
- OWASP Top 10 applicability (if web-facing code)?
- Input validation and sanitization present?

**Performance**
- N+1 query patterns or unnecessary repeated lookups?
- Excessive memory allocations in hot paths?
- Missing indexes or inefficient data structures?
- Unbounded loops or recursion?

**Quality**
- Naming clarity: do names reflect intent?
- Cyclomatic complexity: is any function doing too much?
- Duplication: is there copy-paste that should be extracted?
- Dead code: unreachable branches, unused variables or imports?

**Tests**
- Is test coverage adequate for the new code paths?
- Are assertions meaningful (not just "no error")?
- Are edge cases and failure paths tested?
- Are tests isolated (no shared global state)?

### 4. Classify findings by severity
Group all findings into three tiers:

` + "```" + `
CRITICAL  — Must fix before /yaah:ship (correctness bugs, security vulnerabilities)
WARNING   — Should fix soon (performance issues, missing test coverage, code quality)
SUGGESTION — Nice to have (naming, minor refactors, style)
` + "```" + `

### 5. Write REVIEW.md
Write the full review to ` + "`.planning/phases/{NN}-{slug}/REVIEW.md`" + ` with YAML frontmatter:
` + "```" + `yaml
---
phase: {N}
reviewed: {timestamp}
reviewer: yaah/review
files_reviewed: {count}
critical: {count}
warnings: {count}
suggestions: {count}
status: approved   # approved | needs-work | blocked
---
` + "```" + `

Report sections:
- **Summary**: one paragraph overview of the implementation quality
- **CRITICAL**: each finding with file, line range, explanation, and suggested fix
- **WARNING**: each finding with file, line range, and explanation
- **SUGGESTION**: brief list of improvement ideas
- **Verdict**: ` + "`APPROVED`" + ` (no criticals), ` + "`NEEDS WORK`" + ` (warnings only), or ` + "`BLOCKED`" + ` (has criticals)

### 6. Report to user
Print a concise summary:
` + "```" + `
Review complete: Phase {N} — {name}
Files reviewed: {count}
CRITICAL: {count}  WARNING: {count}  SUGGESTION: {count}
Verdict: {APPROVED | NEEDS WORK | BLOCKED}

Full report: .planning/phases/{NN}-{slug}/REVIEW.md
` + "```" + `

If CRITICAL issues were found, list each one with a suggested fix and recommend resolving them before running ` + "`/yaah:ship`" + `.

## Rules
- NEVER skip files listed in PLAN.md or SUMMARY.md — all must be reviewed
- A single CRITICAL finding sets the verdict to BLOCKED
- REVIEW.md must be written even if there are zero findings (clean bill of health)
- Do not approve a phase that has not yet been executed (no SUMMARY.md)
`
}
