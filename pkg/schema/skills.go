package schema

// SkillSource defines where a skill is loaded from — local path or remote git.
type SkillSource struct {
	// Path is a local filesystem path to the SKILL.md file.
	Path string `json:"path,omitempty" jsonschema:"description=Local path to the SKILL.md file"`

	// Uses references a remote git repo in GitHub Actions style: owner/repo@ref
	// Examples:
	//   github.com/dirien/my-skills@v1.0.0
	//   github.com/dirien/my-skills@main
	//   github.com/dirien/my-skills@abc1234
	Uses string `json:"uses,omitempty" jsonschema:"description=Remote git reference: host/owner/repo@tag-or-sha"`

	// Subpath is the path within the remote repo to the SKILL.md (default: SKILL.md in root).
	Subpath string `json:"subpath,omitempty" jsonschema:"description=Path within the remote repo to the skill file"`
}

// Skill describes a single skill configuration.
type Skill struct {
	Name                   string      `json:"name"                             jsonschema:"description=Unique skill identifier"`
	Description            string      `json:"description,omitempty"            jsonschema:"description=Human-readable description (used by Claude for auto-invocation)"`
	Source                 SkillSource `json:"source"                           jsonschema:"description=Where to load the skill from"`
	Enabled                bool        `json:"enabled"                          jsonschema:"description=Whether this skill is active,default=true"`
	ArgumentHint           string      `json:"argumentHint,omitempty"           jsonschema:"description=Usage hint shown during autocomplete (e.g. [issue-number])"`
	DisableModelInvocation bool        `json:"disableModelInvocation,omitempty" jsonschema:"description=Prevent Claude from auto-loading this skill"`
	UserInvocable          *bool       `json:"userInvocable,omitempty"          jsonschema:"description=Show in /menu (default true)"`
	AllowedTools           string      `json:"allowedTools,omitempty"           jsonschema:"description=Comma-separated tool allowlist"`
	Model                  string      `json:"model,omitempty"                  jsonschema:"description=Model override when skill is active"`
	Context                string      `json:"context,omitempty"                jsonschema:"description=Set to fork for subagent execution"`
	AgentType              string      `json:"agent,omitempty"                  jsonschema:"description=Subagent type when context=fork (Explore/Plan/etc)"`
	SkillHooks             HooksConfig `json:"hooks,omitempty"                  jsonschema:"description=Lifecycle hooks scoped to this skill"`
}

// SkillsConfig holds all skill definitions.
type SkillsConfig struct {
	Skills []Skill `json:"skills" jsonschema:"description=List of skills to register"`
}
