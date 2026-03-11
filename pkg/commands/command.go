package commands

// Command is the interface for custom Claude Code slash commands.
type Command interface {
	// Name returns the command name (invoked as /name).
	Name() string

	// Description returns what this command does.
	Description() string

	// ArgumentHint returns usage hint (e.g., "<path> [options]"), or empty.
	ArgumentHint() string

	// Model returns the model override, or empty for default.
	Model() string

	// AllowedTools returns comma-separated tool allowlist, or empty for all.
	AllowedTools() string

	// Content returns the full markdown body (instruction prompt).
	Content() string
}

// CommandWithAdvanced is an optional interface for commands that need advanced
// frontmatter fields beyond the basics.
type CommandWithAdvanced interface {
	Command
	Advanced() CommandAdvanced
}

// CommandAdvanced holds optional advanced frontmatter fields for a command.
type CommandAdvanced struct {
	DisableModelInvocation bool
	UserInvocable          *bool  // nil = default (true)
	Context                string // "fork" for subagent execution
	AgentType              string // subagent type when context=fork
}

// Registry holds all registered commands.
type Registry struct {
	commands []Command
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers a command.
func (r *Registry) Register(c Command) {
	r.commands = append(r.commands, c)
}

// Commands returns all registered commands.
func (r *Registry) Commands() []Command {
	return r.commands
}
