package generator

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// AgentGenerator produces configuration files for a specific coding agent.
type AgentGenerator interface {
	// Target returns which agent this generator targets.
	Target() schema.TargetAgent

	// GenerateSettings produces the main settings/config file content.
	// Returns nil if the agent has no single settings file.
	GenerateSettings(cfg *schema.HarnessConfig) ([]byte, error)

	// SettingsPath returns the relative path for the settings file.
	// Empty string means no settings file.
	SettingsPath() string

	// GenerateMCP produces a separate MCP config file.
	// Returns nil if MCP config is embedded in the settings file.
	GenerateMCP(cfg *schema.HarnessConfig) ([]byte, error)

	// MCPPath returns the relative path for the MCP config file.
	MCPPath() string

	// GenerateHooks produces a separate hooks config file.
	// Returns nil if hooks are embedded in the settings file.
	GenerateHooks(cfg *schema.HarnessConfig) ([]byte, error)

	// HooksPath returns the relative path for the hooks file.
	HooksPath() string

	// SkillsDir returns the relative path to the skills directory.
	SkillsDir() string

	// AgentsDir returns the relative path to the agents directory.
	// Empty string means file-based agents are unsupported.
	AgentsDir() string

	// AgentFileExt returns the file extension for agent files (e.g. ".md" or ".agent.md").
	AgentFileExt() string

	// CommandsDir returns the relative path to the commands directory.
	// Empty string means commands are unsupported.
	CommandsDir() string
}
