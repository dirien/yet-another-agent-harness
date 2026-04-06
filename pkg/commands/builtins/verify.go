package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// VerifyCommand implements post-execution verification via /yaah:verify.
type VerifyCommand struct{}

// NewVerifyCommand creates a new VerifyCommand.
func NewVerifyCommand() *VerifyCommand { return &VerifyCommand{} }

func (c *VerifyCommand) Name() string { return "yaah/verify" }
func (c *VerifyCommand) Description() string {
	return "Verify implementation against plan requirements"
}
func (c *VerifyCommand) ArgumentHint() string { return "[phase-number]" }
func (c *VerifyCommand) Model() string        { return "" }
func (c *VerifyCommand) AllowedTools() string { return "" }
func (c *VerifyCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *VerifyCommand) Content() string {
	return `# /yaah:verify — Post-Execution Verification

## When to use
When the user runs ` + "`/yaah:verify [phase]`" + ` after ` + "`/yaah:execute`" + ` completes.

## Goal
Produce a structured verification report that proves (or disproves) that the implementation satisfies the plan and requirements.

## Steps

### 1. Load plan and execution artifacts
- Read all ` + "`PLAN.md`" + ` files for the target phase
- Read all ` + "`SUMMARY.md`" + ` files produced by ` + "`/yaah:execute`" + `
- Collect every task with its Files, Action, and Verify fields

### 2. Three-level artifact checks
For each artifact claimed in the plan:

**Level 1 — Existence** (Glob):
- Confirm every file listed in ` + "`files_modified`" + ` exists on disk

**Level 2 — Content** (Read + Grep):
- Verify the file contains the expected symbols, exports, or implementations
- Check that interfaces are implemented, functions are defined, configs are present

**Level 3 — Wiring** (Grep for imports and registrations):
- Confirm new code is actually imported and wired into the application
- Check registry registrations, route mounts, flag definitions, etc.

### 3. Stub detection
Search for implementation stubs that indicate incomplete work:
- ` + "`TODO`" + ` and ` + "`FIXME`" + ` comments in modified files
- ` + "`panic(\"not implemented\")`" + ` or ` + "`panic(\"TODO\")`" + `
- Functions that ` + "`return nil`" + ` with no logic
- Empty function bodies ` + "`{}`" + ` in non-trivial functions

Report every stub as a GAP.

### 4. Build verification
Auto-detect and run the project build:
- Go: ` + "`go build -v ./...`" + `
- Node.js: ` + "`npm run build`" + ` or ` + "`yarn build`" + `
- Rust: ` + "`cargo build`" + `
- Python: ` + "`python -m py_compile`" + ` on modified files
- Java/Gradle: ` + "`./gradlew build`" + `
- Java/Maven: ` + "`mvn compile`" + `

Record exit code and stderr output.

### 5. Test verification
Auto-detect and run the test suite:
- Go: ` + "`go test -race ./...`" + `
- Node.js: ` + "`npm test`" + `
- Rust: ` + "`cargo test`" + `
- Python: ` + "`pytest`" + `
- Java/Gradle: ` + "`./gradlew test`" + `
- Java/Maven: ` + "`mvn test`" + `

Record pass/fail/skip counts.

### 6. Requirements traceability
For each REQ-ID in scope for this phase:
- Trace: REQ-ID → task → artifact → verified?
- Status: ` + "`SATISFIED`" + ` (all checks pass) or ` + "`GAP`" + ` (any check fails)

### 7. Write VERIFICATION.md
Write to ` + "`.planning/phases/{NN}-{slug}/VERIFICATION.md`" + ` with YAML frontmatter:
` + "```" + `yaml
---
phase: 1
verified: 2024-01-15T10:30:00Z
status: passed   # passed | failed | partial
score: 12/14     # checks passed / total checks
---
` + "```" + `

Report sections:
- **Artifact Checks**: table of file, level, result
- **Stub Report**: list of stubs found (empty = clean)
- **Build**: exit code, duration, key output
- **Tests**: pass/fail/skip counts, any failures
- **Requirements Traceability**: REQ-ID, status, notes

### 8. Update STATE.md and commit
Set ` + "`status: verified`" + ` (or ` + "`failed`" + `), update ` + "`last_updated`" + `.
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): verification report phase {N}\"`" + `

## Rules
- All three artifact check levels must pass for an artifact to be SATISFIED
- A single stub is sufficient to mark that task as GAP
- Build or test failure marks the entire phase as ` + "`failed`" + `
- Never skip traceability — every in-scope REQ-ID must have a status
`
}
