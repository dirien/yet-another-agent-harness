package plugins

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

var _ MarketplacePlugin = (*Codex)(nil)

// Codex provides the OpenAI Codex plugin for Claude Code.
// It integrates Codex as a code review and task delegation agent.
// Source: https://github.com/openai/codex-plugin-cc
type Codex struct{}

// NewCodex creates a new Codex plugin.
func NewCodex() *Codex { return &Codex{} }

func (c *Codex) Plugin() schema.Plugin {
	return schema.Plugin{
		Name:        "codex",
		Version:     "1.0.0",
		Description: "OpenAI Codex integration for code review and task delegation",
		Author: schema.PluginAuthor{
			Name: "OpenAI",
			URL:  "https://openai.com",
		},
		Homepage:   "https://github.com/openai/codex-plugin-cc",
		Repository: "https://github.com/openai/codex-plugin-cc",
		License:    "Apache-2.0",
		Keywords:   []string{"codex", "openai", "code-review", "delegation"},
		Commands: []string{
			"commands/review.md",
			"commands/adversarial-review.md",
			"commands/rescue.md",
			"commands/setup.md",
			"commands/status.md",
			"commands/result.md",
			"commands/cancel.md",
		},
		Agents: []string{"agents/"},
		Skills: []string{
			"skills/codex-cli-runtime",
			"skills/codex-result-handling",
			"skills/gpt-5-4-prompting",
		},
	}
}

func (c *Codex) MarketplaceKey() string { return "codex@openai-codex" }
