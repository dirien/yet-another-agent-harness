package generator

import (
	"encoding/json"
	"fmt"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// CopilotGenerator produces configuration for GitHub Copilot CLI.
type CopilotGenerator struct{}

func (g *CopilotGenerator) Target() schema.TargetAgent { return schema.TargetCopilot }

// GenerateSettings returns nil — Copilot has no single settings file.
func (g *CopilotGenerator) GenerateSettings(_ *schema.HarnessConfig) ([]byte, error) {
	return nil, nil
}

func (g *CopilotGenerator) SettingsPath() string { return "" }

// copilotMCPConfig is the output format for .copilot/mcp-config.json.
type copilotMCPConfig struct {
	MCPServers map[string]copilotMCPServer `json:"mcpServers"`
}

type copilotMCPServer struct {
	Type    string            `json:"type"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

func (g *CopilotGenerator) GenerateMCP(cfg *schema.HarnessConfig) ([]byte, error) {
	if cfg.MCP == nil || len(cfg.MCP.Servers) == 0 {
		return nil, nil
	}

	out := copilotMCPConfig{
		MCPServers: make(map[string]copilotMCPServer),
	}
	for _, srv := range cfg.MCP.Servers {
		entry := copilotMCPServer{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
		}
		if srv.Command != "" {
			entry.Type = "stdio"
		} else {
			entry.Type = "http"
		}
		out.MCPServers[srv.Name] = entry
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal copilot mcp-config.json: %w", err)
	}
	return data, nil
}

func (g *CopilotGenerator) MCPPath() string { return ".copilot/mcp-config.json" }

// copilotHooksConfig is the output format for .github/hooks/hooks.json.
type copilotHooksConfig struct {
	Version int                           `json:"version"`
	Hooks   map[string][]copilotHookEntry `json:"hooks"`
}

type copilotHookEntry struct {
	Type       string `json:"type"`
	Bash       string `json:"bash"`
	TimeoutSec int    `json:"timeoutSec"`
}

func (g *CopilotGenerator) GenerateHooks(cfg *schema.HarnessConfig) ([]byte, error) {
	hooksMap := make(map[string][]copilotHookEntry)

	// Add hooks for events that are configured in the harness config.
	if cfg.Hooks != nil {
		for event := range cfg.Hooks {
			copilotEvent := CopilotEventName(event)
			if copilotEvent == "" {
				continue
			}
			if _, exists := hooksMap[copilotEvent]; exists {
				continue
			}
			hooksMap[copilotEvent] = []copilotHookEntry{
				{
					Type:       "command",
					Bash:       fmt.Sprintf("yaah hook %s", event),
					TimeoutSec: 30,
				},
			}
		}
	}

	if len(hooksMap) == 0 {
		return nil, nil
	}

	out := copilotHooksConfig{
		Version: 1,
		Hooks:   hooksMap,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal copilot hooks.json: %w", err)
	}
	return data, nil
}

func (g *CopilotGenerator) HooksPath() string    { return ".github/hooks/hooks.json" }
func (g *CopilotGenerator) SkillsDir() string    { return ".github/skills" }
func (g *CopilotGenerator) AgentsDir() string    { return ".github/agents" }
func (g *CopilotGenerator) AgentFileExt() string { return ".agent.md" }
func (g *CopilotGenerator) CommandsDir() string  { return "" }
