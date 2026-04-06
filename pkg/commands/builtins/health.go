package builtins

// HealthCommand validates .planning/ integrity via /yaah:health.
// It is lightweight and does not implement CommandWithAdvanced.
type HealthCommand struct{}

// NewHealthCommand creates a new HealthCommand.
func NewHealthCommand() *HealthCommand { return &HealthCommand{} }

func (c *HealthCommand) Name() string { return "yaah/health" }
func (c *HealthCommand) Description() string {
	return "Validate planning directory integrity and consistency"
}
func (c *HealthCommand) ArgumentHint() string { return "" }
func (c *HealthCommand) Model() string        { return "" }
func (c *HealthCommand) AllowedTools() string { return "Read, Glob, Grep" }

func (c *HealthCommand) Content() string {
	return `# /yaah:health — Planning Directory Health Check

## When to use
When the user runs ` + "`/yaah:health`" + ` to validate the consistency and integrity of the ` + "`.planning/`" + ` directory.

## Steps

### 1. Check .planning/ exists
- Glob for ` + "`.planning/`" + `.
- If it does not exist: report FAIL and suggest running ` + "`/yaah:init`" + `. Stop.

### 2. Validate required top-level files
Check each of the following files exists and is non-empty:
- ` + "`.planning/PROJECT.md`" + `
- ` + "`.planning/REQUIREMENTS.md`" + `
- ` + "`.planning/ROADMAP.md`" + `
- ` + "`.planning/STATE.md`" + `
- ` + "`.planning/config.json`" + `

Report PASS or FAIL for each file.

### 3. Validate STATE.md frontmatter
- Read ` + "`.planning/STATE.md`" + `.
- Confirm the YAML frontmatter block (between ` + "`---`" + ` delimiters) is present and parseable.
- Check required fields exist: ` + "`milestone`" + `, ` + "`phase`" + `, ` + "`status`" + `, ` + "`last_updated`" + `.
- Report PASS or FAIL with the specific missing or malformed fields.

### 4. Validate ROADMAP.md phase numbering
- Read ` + "`.planning/ROADMAP.md`" + `.
- Extract all phase numbers from headings (e.g., "Phase 1", "Phase 2").
- Confirm they are sequential starting from 1 with no gaps or duplicates.
- Report PASS or FAIL with details if gaps or duplicates are found.

### 5. Check REQ-ID coverage
- Grep ` + "`.planning/REQUIREMENTS.md`" + ` to collect all REQ-IDs (pattern: ` + "`REQ-\\d+`" + `).
- Grep all ` + "`PLAN.md`" + ` files under ` + "`.planning/phases/`" + ` for each REQ-ID.
- Flag any REQ-ID that does not appear in at least one PLAN.md.
- Report PASS if all are referenced, FAIL with list of orphaned REQ-IDs.

### 6. Check for orphaned phase directories
- Glob ` + "`.planning/phases/*/`" + ` to list all phase directories.
- Extract phase numbers from directory names.
- Cross-reference with phase numbers in ROADMAP.md.
- Flag any directory whose phase number does not appear in ROADMAP.md.
- Report PASS or FAIL with list of orphaned directories.

### 7. Validate PLAN.md frontmatter
- Glob all ` + "`PLAN.md`" + ` files under ` + "`.planning/phases/`" + `.
- For each, check the YAML frontmatter contains required fields: ` + "`phase`" + `, ` + "`plan`" + `, ` + "`wave`" + `, ` + "`files_modified`" + `.
- Report PASS or FAIL per file, listing any missing fields.

### 8. Report summary
Print a structured summary:
` + "```" + `
Health Check Results
====================
Checks passed: {X}/{total}

PASS  .planning/ exists
PASS  PROJECT.md present and non-empty
FAIL  config.json missing
...

Issues found:
- config.json: file missing — run /yaah:init to regenerate
- REQ-003: not referenced in any PLAN.md — add to relevant phase plan
` + "```" + `

### 9. Suggest fixes
For each issue found, suggest the specific action to resolve it:
- Missing file → suggest the command or manual step to create it
- Orphaned REQ-ID → suggest adding it to the relevant phase PLAN.md
- Orphaned directory → suggest adding it to ROADMAP.md or deleting it
- Invalid frontmatter → show the required YAML structure

## Rules
- Report ALL issues found — do not stop at the first failure
- Never modify any files — this is a read-only diagnostic command
- Always show the {checks passed}/{total checks} score at the top of the report
`
}
