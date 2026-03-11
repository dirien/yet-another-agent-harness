package schema

// HarnessConfig is the root configuration composing every concern.
type HarnessConfig struct {
	Version     string             `json:"version"`
	Settings    *Settings          `json:"settings,omitempty"`
	Hooks       HooksConfig        `json:"hooks,omitempty"`
	MCP         *MCPConfig         `json:"mcp,omitempty"`
	LSP         *LSPConfig         `json:"lsp,omitempty"`
	Skills      *SkillsConfig      `json:"skills,omitempty"`
	Agents      *AgentsConfig      `json:"agents,omitempty"`
	Commands    *CommandsConfig    `json:"commands,omitempty"`
	Plugins     *PluginsConfig     `json:"plugins,omitempty"`
	Permissions *PermissionsConfig `json:"permissions,omitempty"`
}

// PermissionRule defines an allow or deny rule for a tool pattern.
type PermissionRule struct {
	Pattern string `json:"pattern"`
	Allow   bool   `json:"allow"`
}

// PermissionsConfig holds tool permission rules.
type PermissionsConfig struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}
