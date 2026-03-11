package skills

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// Skill is the interface for generating SKILL.md files that Claude Code can use.
type Skill interface {
	// Name returns the skill identifier (used as directory name under .claude/skills/).
	Name() string

	// Description returns a human-readable description for the skill.
	Description() string

	// Content returns the full SKILL.md markdown content (body after frontmatter).
	Content() string

	// Source returns the skill source (local or remote).
	Source() schema.SkillSource
}

// SkillWithFrontmatter is an optional interface for skills that provide
// YAML frontmatter metadata (allowed-tools, model, context, etc.).
type SkillWithFrontmatter interface {
	Skill
	Frontmatter() SkillFrontmatter
}

// SkillFrontmatter holds optional YAML frontmatter fields for a skill.
type SkillFrontmatter struct {
	ArgumentHint           string
	DisableModelInvocation bool
	UserInvocable          *bool // nil = default (true)
	AllowedTools           string
	Model                  string
	Context                string // "fork" for subagent execution
	AgentType              string // subagent type when context=fork
}

// SkillWithFiles is an optional interface for skills that provide additional
// files beyond SKILL.md (e.g. a references/ folder).
type SkillWithFiles interface {
	Skill
	// ExtraFiles returns a map of relative path → content for additional files
	// that should be written alongside SKILL.md.
	ExtraFiles() map[string]string
}

// Registry holds all registered skills.
type Registry struct {
	skills []Skill
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers a skill.
func (r *Registry) Register(s Skill) {
	r.skills = append(r.skills, s)
}

// Skills returns all registered skills.
func (r *Registry) Skills() []Skill {
	return r.skills
}
