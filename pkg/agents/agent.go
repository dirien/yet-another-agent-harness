package agents

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// Agent is the interface for custom Claude Code agent definitions.
type Agent interface {
	// Name returns the agent identifier (invoked as @agent-name).
	Name() string

	// Description returns when this agent should be used.
	Description() string

	// Model returns the model override (haiku/sonnet/opus/inherit), or empty for default.
	Model() string

	// Tools returns the comma-separated tool allowlist, or empty for all tools.
	Tools() string

	// Content returns the full markdown body (system prompt for the agent).
	Content() string
}

// AgentWithAdvanced is an optional interface for agents that need advanced
// frontmatter fields beyond name/description/model/tools.
type AgentWithAdvanced interface {
	Agent
	Advanced() AgentAdvanced
}

// AgentWithSource is an optional interface for agents loaded from a remote
// git repository, providing provenance information.
type AgentWithSource interface {
	Agent
	Uses() string    // e.g. "github.com/owner/repo@ref"
	Subpath() string // path within the repo
}

// AgentAdvanced holds optional advanced frontmatter fields for an agent.
type AgentAdvanced struct {
	DisallowedTools string
	PermissionMode  string // default/acceptEdits/dontAsk/bypassPermissions/plan
	MaxTurns        int
	Skills          []string
	McpServers      map[string]any
	Hooks           schema.HooksConfig
	Memory          string // user/project/local
	Background      bool
	Isolation       string // worktree
}

// Registry holds all registered agents.
type Registry struct {
	agents []Agent
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers an agent.
func (r *Registry) Register(a Agent) {
	r.agents = append(r.agents, a)
}

// Agents returns all registered agents.
func (r *Registry) Agents() []Agent {
	return r.agents
}

// ToConfig converts agents into an AgentsConfig for generation.
func (r *Registry) ToConfig() *schema.AgentsConfig {
	cfg := &schema.AgentsConfig{}
	for _, a := range r.agents {
		agent := schema.Agent{
			Name:        a.Name(),
			Description: a.Description(),
			Source:      ".claude/agents/" + a.Name() + ".md",
			Model:       a.Model(),
			Tools:       a.Tools(),
		}
		if adv, ok := a.(AgentWithAdvanced); ok {
			opts := adv.Advanced()
			agent.DisallowedTools = opts.DisallowedTools
			agent.PermissionMode = opts.PermissionMode
			agent.MaxTurns = opts.MaxTurns
			agent.Skills = opts.Skills
			agent.McpServers = opts.McpServers
			agent.Hooks = opts.Hooks
			agent.Memory = opts.Memory
			agent.Background = opts.Background
			agent.Isolation = opts.Isolation
		}
		cfg.Agents = append(cfg.Agents, agent)
	}
	return cfg
}
