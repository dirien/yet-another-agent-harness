package generator

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// ClaudeGenerator produces configuration for Claude Code.
type ClaudeGenerator struct{}

func (g *ClaudeGenerator) Target() schema.TargetAgent { return schema.TargetClaude }

func (g *ClaudeGenerator) GenerateSettings(cfg *schema.HarnessConfig) ([]byte, error) {
	return GenerateClaudeSettings(cfg)
}

func (g *ClaudeGenerator) SettingsPath() string { return ".claude/settings.json" }

func (g *ClaudeGenerator) GenerateMCP(cfg *schema.HarnessConfig) ([]byte, error) {
	return GenerateMCPJSON(cfg)
}

func (g *ClaudeGenerator) MCPPath() string { return ".mcp.json" }

// GenerateHooks returns nil because Claude Code hooks are embedded in settings.json.
func (g *ClaudeGenerator) GenerateHooks(_ *schema.HarnessConfig) ([]byte, error) {
	return nil, nil
}

func (g *ClaudeGenerator) HooksPath() string    { return "" }
func (g *ClaudeGenerator) SkillsDir() string    { return ".claude/skills" }
func (g *ClaudeGenerator) AgentsDir() string    { return ".claude/agents" }
func (g *ClaudeGenerator) AgentFileExt() string { return ".md" }
func (g *ClaudeGenerator) CommandsDir() string  { return ".claude/commands" }
