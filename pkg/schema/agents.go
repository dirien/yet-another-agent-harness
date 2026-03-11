package schema

// AgentPermission controls what a spawned agent is allowed to do.
type AgentPermission string

const (
	AgentPermAsk   AgentPermission = "ask"
	AgentPermAllow AgentPermission = "allow"
	AgentPermDeny  AgentPermission = "deny"
)

// Agent describes a custom agent definition (written as a markdown file with YAML frontmatter).
type Agent struct {
	Name            string          `json:"name"                      jsonschema:"description=Agent identifier (invoked as @agent-name)"`
	Description     string          `json:"description,omitempty"     jsonschema:"description=What this agent does"`
	Source          string          `json:"source"                    jsonschema:"description=Path to the agent markdown file"`
	Model           string          `json:"model,omitempty"           jsonschema:"description=Model override (sonnet/opus/haiku/inherit)"`
	Tools           string          `json:"tools,omitempty"           jsonschema:"description=Comma-separated tool allowlist (e.g. Read Write Bash(*))"`
	DisallowedTools string          `json:"disallowedTools,omitempty" jsonschema:"description=Comma-separated tool denylist"`
	PermissionMode  string          `json:"permissionMode,omitempty"  jsonschema:"description=Permission mode: default/acceptEdits/dontAsk/bypassPermissions/plan"`
	MaxTurns        int             `json:"maxTurns,omitempty"        jsonschema:"description=Maximum agentic turns before stopping"`
	Skills          []string        `json:"skills,omitempty"          jsonschema:"description=Skills to preload into agent context"`
	McpServers      map[string]any  `json:"mcpServers,omitempty"      jsonschema:"description=MCP servers for this agent"`
	Hooks           HooksConfig     `json:"hooks,omitempty"           jsonschema:"description=Lifecycle hooks (PreToolUse/PostToolUse/Stop)"`
	Memory          string          `json:"memory,omitempty"          jsonschema:"description=Persistent memory scope: user/project/local"`
	Background      bool            `json:"background,omitempty"      jsonschema:"description=Run as background task"`
	Isolation       string          `json:"isolation,omitempty"       jsonschema:"description=Isolation mode: worktree"`
	Permission      AgentPermission `json:"permission,omitempty"      jsonschema:"enum=ask,enum=allow,enum=deny,default=ask"`
}

// AgentsConfig holds all agent definitions.
type AgentsConfig struct {
	Agents []Agent `json:"agents" jsonschema:"description=List of custom agents"`
}

// Command describes a custom slash command (written as a markdown file with YAML frontmatter).
type Command struct {
	Name                   string      `json:"name"                          jsonschema:"description=Command name (invoked as /name)"`
	Description            string      `json:"description,omitempty"         jsonschema:"description=Shown in command listing"`
	Source                 string      `json:"source"                        jsonschema:"description=Path to the command markdown file"`
	ArgumentHint           string      `json:"argumentHint,omitempty"        jsonschema:"description=Usage hint (e.g. <path> [options])"`
	Model                  string      `json:"model,omitempty"               jsonschema:"description=Model override"`
	AllowedTools           string      `json:"allowedTools,omitempty"        jsonschema:"description=Comma-separated tool allowlist"`
	DisableModelInvocation bool        `json:"disableModelInvocation,omitempty" jsonschema:"description=Prevent Claude from auto-loading this command"`
	UserInvocable          *bool       `json:"userInvocable,omitempty"       jsonschema:"description=Show in /menu (default true)"`
	Context                string      `json:"context,omitempty"             jsonschema:"description=Set to fork for subagent execution"`
	AgentType              string      `json:"agent,omitempty"               jsonschema:"description=Subagent type when context=fork"`
	CommandHooks           HooksConfig `json:"hooks,omitempty"               jsonschema:"description=Lifecycle hooks scoped to this command"`
}

// CommandsConfig holds all command definitions.
type CommandsConfig struct {
	Commands []Command `json:"commands" jsonschema:"description=List of custom slash commands"`
}

// PluginAuthor holds plugin author metadata.
type PluginAuthor struct {
	Name  string `json:"name"            jsonschema:"description=Author name"`
	Email string `json:"email,omitempty" jsonschema:"description=Author email"`
	URL   string `json:"url,omitempty"   jsonschema:"description=Author URL"`
}

// Plugin describes a Claude Code plugin definition (plugin.json).
// Matches the official .claude-plugin/plugin.json specification.
type Plugin struct {
	Name         string       `json:"name"                    jsonschema:"description=Plugin identifier"`
	Version      string       `json:"version"                 jsonschema:"description=Semver version"`
	Description  string       `json:"description,omitempty"   jsonschema:"description=Plugin description"`
	Author       PluginAuthor `json:"author"                  jsonschema:"description=Plugin author"`
	Homepage     string       `json:"homepage,omitempty"      jsonschema:"description=Homepage URL"`
	Repository   string       `json:"repository,omitempty"    jsonschema:"description=Git repository URL"`
	License      string       `json:"license,omitempty"       jsonschema:"description=License identifier (MIT/Apache-2.0/etc)"`
	Keywords     []string     `json:"keywords,omitempty"      jsonschema:"description=Discovery keywords"`
	Commands     []string     `json:"commands,omitempty"      jsonschema:"description=Paths to command markdown files"`
	Agents       []string     `json:"agents,omitempty"        jsonschema:"description=Paths to agent markdown files or agent directory"`
	Skills       []string     `json:"skills,omitempty"        jsonschema:"description=Paths to skill directories"`
	Hooks        HooksConfig  `json:"hooks,omitempty"         jsonschema:"description=Hook registrations"`
	McpServers   string       `json:"mcpServers,omitempty"    jsonschema:"description=Path to .mcp.json or inline MCP config"`
	LspServers   string       `json:"lspServers,omitempty"    jsonschema:"description=Path to .lsp.json or inline LSP config"`
	OutputStyles string       `json:"outputStyles,omitempty"  jsonschema:"description=Path to output styles directory"`
}

// PluginsConfig holds plugin definitions for generation.
type PluginsConfig struct {
	Plugins []Plugin `json:"plugins" jsonschema:"description=List of plugins to generate"`
}
