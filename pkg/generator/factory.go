package generator

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// ForTarget returns the AgentGenerator for the given target agent.
func ForTarget(target schema.TargetAgent) AgentGenerator {
	switch target {
	case schema.TargetClaude:
		return &ClaudeGenerator{}
	case schema.TargetOpenCode:
		return &OpenCodeGenerator{}
	case schema.TargetCodex:
		return &CodexGenerator{}
	case schema.TargetCopilot:
		return &CopilotGenerator{}
	default:
		return &ClaudeGenerator{}
	}
}
