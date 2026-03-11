package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// PRSkill generates a SKILL.md for pull request creation.
type PRSkill struct{}

func NewPRSkill() *PRSkill { return &PRSkill{} }

func (p *PRSkill) Name() string        { return "pr" }
func (p *PRSkill) Description() string  { return "Create pull requests with structured description" }
func (p *PRSkill) Source() schema.SkillSource {
	return schema.SkillSource{Path: ".claude/skills/pr/SKILL.md"}
}

func (p *PRSkill) Content() string {
	return `# /pr — Pull Request Creation

## When to use
When the user runs /pr or asks to create a pull request.

## Steps

1. Run ` + "`git log --oneline main..HEAD`" + ` to understand all commits
2. Run ` + "`git diff main...HEAD`" + ` to see the full diff
3. Check if the branch is pushed: ` + "`git rev-parse --abbrev-ref --symbolic-full-name @{u}`" + `
4. If not pushed, push with ` + "`git push -u origin HEAD`" + `
5. Create the PR using ` + "`gh pr create`" + ` with:
   - Title: short, under 70 chars
   - Body using the template:
     ## Summary
     <1-3 bullet points>

     ## Test plan
     - [ ] Tests pass
     - [ ] Manual verification

## Rules
- NEVER force push
- Analyze ALL commits, not just the latest
- Keep the title concise — details go in the body
- Return the PR URL to the user when done
`
}
