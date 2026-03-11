package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// ReviewSkill generates a SKILL.md for code review.
type ReviewSkill struct{}

func NewReviewSkill() *ReviewSkill { return &ReviewSkill{} }

func (r *ReviewSkill) Name() string        { return "review" }
func (r *ReviewSkill) Description() string  { return "Review code changes for quality, security, and correctness" }
func (r *ReviewSkill) Source() schema.SkillSource {
	return schema.SkillSource{Path: ".claude/skills/review/SKILL.md"}
}

func (r *ReviewSkill) Content() string {
	return `# /review — Code Review

## When to use
When the user runs /review or asks to review changes/PR.

## Steps

1. Identify what to review:
   - If a PR number is given: ` + "`gh pr diff <number>`" + `
   - Otherwise: ` + "`git diff`" + ` for unstaged, ` + "`git diff --staged`" + ` for staged
2. Read each changed file fully for context
3. Check for:
   - **Security**: injection, XSS, hardcoded secrets, OWASP top 10
   - **Correctness**: edge cases, off-by-one, null handling, race conditions
   - **Quality**: naming, complexity, duplication, dead code
   - **Performance**: N+1 queries, unnecessary allocations, missing indexes
   - **Tests**: adequate coverage, meaningful assertions
4. Report findings grouped by severity:
   - CRITICAL: security issues, data loss risks
   - WARNING: bugs, correctness issues
   - SUGGESTION: style, refactoring opportunities

## Rules
- Be specific: reference file:line
- Suggest fixes, don't just point out problems
- Acknowledge what's done well
- Don't nitpick formatting if a linter is configured
`
}
