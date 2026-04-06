package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ShipCommand creates a PR from verified phase work via /yaah:ship.
type ShipCommand struct{}

// NewShipCommand creates a new ShipCommand.
func NewShipCommand() *ShipCommand { return &ShipCommand{} }

func (c *ShipCommand) Name() string        { return "yaah/ship" }
func (c *ShipCommand) Description() string { return "Create a pull request from verified phase work" }
func (c *ShipCommand) ArgumentHint() string { return "[phase-number]" }
func (c *ShipCommand) Model() string        { return "" }
func (c *ShipCommand) AllowedTools() string { return "" }
func (c *ShipCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ShipCommand) Content() string {
	return `# /yaah:ship — Ship Verified Work

## When to use
When the user runs ` + "`/yaah:ship`" + ` or ` + "`/yaah:ship <N>`" + ` after ` + "`/yaah:verify`" + ` passes.

## Prerequisites
- VERIFICATION.md must exist with status "passed" for the target phase
- If no argument, use current phase from STATE.md

## Steps

### 1. Validate readiness
- Read VERIFICATION.md — confirm status is "passed"
- If status is "gaps_found", STOP and suggest ` + "`/yaah:verify {N}`" + ` first
- Read all SUMMARY.md files to collect changes

### 2. Prepare branch
- Check current git branch
- If on main/master, create a feature branch: ` + "`git checkout -b phase-{N}-{slug}`" + `
- If already on a feature branch, use it

### 3. Generate PR description
Analyze all commits from the phase:
- Run ` + "`git log --oneline main..HEAD`" + ` (or appropriate base)
- Read PLAN.md files for original objectives
- Read SUMMARY.md files for what was actually done
- Read VERIFICATION.md for test results

### 4. Create PR
Use ` + "`gh pr create`" + ` with:
- Title: concise summary under 70 chars
- Body structured as:
  ## Summary
  {1-3 bullet points from phase objectives}

  ## Changes
  {list of key changes from SUMMARYs}

  ## Verification
  - Build: {PASS/FAIL}
  - Tests: {PASS/FAIL}
  - Artifact checks: {score}%

  ## Phase
  Phase {N}: {name} from ` + "`.planning/ROADMAP.md`" + `

### 5. Update state
- Update STATE.md: phase status → "shipped", PR URL recorded
- Update ROADMAP.md: mark phase as shipped with date
- Commit: ` + "`docs(planning): ship phase {N}`" + `

### 6. Report
Print PR URL and suggest next step.

## Rules
- NEVER ship unverified work — VERIFICATION.md must show "passed"
- NEVER force push
- Include all phase commits in the PR
- Return the PR URL to the user
`
}
