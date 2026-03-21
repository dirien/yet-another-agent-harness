package schema

import "fmt"

// TargetAgent identifies a supported coding agent.
type TargetAgent string

const (
	TargetClaude   TargetAgent = "claude"
	TargetOpenCode TargetAgent = "opencode"
	TargetCodex    TargetAgent = "codex"
	TargetCopilot  TargetAgent = "copilot"
)

// AllTargets returns every supported target agent.
func AllTargets() []TargetAgent {
	return []TargetAgent{TargetClaude, TargetOpenCode, TargetCodex, TargetCopilot}
}

// ValidateTarget parses a string into a TargetAgent or returns an error.
func ValidateTarget(s string) (TargetAgent, error) {
	switch TargetAgent(s) {
	case TargetClaude, TargetOpenCode, TargetCodex, TargetCopilot:
		return TargetAgent(s), nil
	default:
		return "", fmt.Errorf("unknown target agent %q (valid: claude, opencode, codex, copilot)", s)
	}
}
