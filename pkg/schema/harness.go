package schema

// HarnessConfig is the root configuration composing every concern.
type HarnessConfig struct {
	Version     string             `json:"version"                jsonschema:"description=Schema version,default=1"`
	Settings    *Settings          `json:"settings,omitempty"     jsonschema:"description=Claude Code settings"`
	Hooks       HooksConfig        `json:"hooks,omitempty"        jsonschema:"description=Hook definitions keyed by event"`
	MCP         *MCPConfig         `json:"mcp,omitempty"          jsonschema:"description=MCP server definitions"`
	LSP         *LSPConfig         `json:"lsp,omitempty"          jsonschema:"description=LSP server definitions"`
	Skills      *SkillsConfig      `json:"skills,omitempty"       jsonschema:"description=Skill definitions"`
	Agents      *AgentsConfig      `json:"agents,omitempty"       jsonschema:"description=Agent definitions"`
	Commands    *CommandsConfig    `json:"commands,omitempty"     jsonschema:"description=Custom slash command definitions"`
	Plugins     *PluginsConfig     `json:"plugins,omitempty"      jsonschema:"description=Plugin definitions"`
	Permissions *PermissionsConfig `json:"permissions,omitempty"  jsonschema:"description=Permission rules"`
}

// PermissionRule defines an allow or deny rule for a tool pattern.
type PermissionRule struct {
	Pattern string `json:"pattern" jsonschema:"description=Regex matching tool names"`
	Allow   bool   `json:"allow"   jsonschema:"description=Whether to allow or deny"`
}

// PermissionsConfig holds tool permission rules.
type PermissionsConfig struct {
	Allow []string `json:"allow,omitempty" jsonschema:"description=Tool patterns to always allow"`
	Deny  []string `json:"deny,omitempty"  jsonschema:"description=Tool patterns to always deny"`
}
