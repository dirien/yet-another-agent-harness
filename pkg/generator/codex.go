package generator

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	toml "github.com/pelletier/go-toml/v2"
)

// CodexGenerator produces configuration for OpenAI Codex CLI.
type CodexGenerator struct{}

func (g *CodexGenerator) Target() schema.TargetAgent { return schema.TargetCodex }

// codexConfig is the output format for .codex/config.toml.
type codexConfig struct {
	MCPServers map[string]codexMCPServer `toml:"mcp_servers,omitempty"`
	Notify     []string                  `toml:"notify,omitempty"`
	Features   *codexFeatures            `toml:"features,omitempty"`
}

type codexMCPServer struct {
	Command string            `toml:"command,omitempty"`
	Args    []string          `toml:"args,omitempty"`
	URL     string            `toml:"url,omitempty"`
	Env     map[string]string `toml:"env,omitempty"`
}

type codexFeatures struct {
	CodexHooks bool `toml:"codex_hooks"`
}

func (g *CodexGenerator) GenerateSettings(cfg *schema.HarnessConfig) ([]byte, error) {
	out := codexConfig{
		Notify:   []string{"yaah", "hook", "Notification"},
		Features: &codexFeatures{CodexHooks: true},
	}

	if cfg.MCP != nil && len(cfg.MCP.Servers) > 0 {
		out.MCPServers = make(map[string]codexMCPServer)
		for _, srv := range cfg.MCP.Servers {
			out.MCPServers[srv.Name] = codexMCPServer{
				Command: srv.Command,
				Args:    srv.Args,
				URL:     srv.URL,
				Env:     srv.Env,
			}
		}
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(out); err != nil {
		return nil, fmt.Errorf("marshal codex config.toml: %w", err)
	}
	return buf.Bytes(), nil
}

func (g *CodexGenerator) SettingsPath() string { return ".codex/config.toml" }

// GenerateMCP returns nil — MCP is embedded in config.toml.
func (g *CodexGenerator) GenerateMCP(_ *schema.HarnessConfig) ([]byte, error) {
	return nil, nil
}

func (g *CodexGenerator) MCPPath() string { return "" }

// codexHooksConfig is the output format for .codex/hooks.json.
type codexHooksConfig struct {
	Hooks map[string][]codexHookGroup `json:"hooks"`
}

type codexHookGroup struct {
	Hooks []codexHookEntry `json:"hooks"`
}

type codexHookEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

func (g *CodexGenerator) GenerateHooks(cfg *schema.HarnessConfig) ([]byte, error) {
	hooksMap := make(map[string][]codexHookGroup)

	if cfg.Hooks != nil {
		for event := range cfg.Hooks {
			codexEvent := CodexEventName(event)
			if codexEvent == "" || codexEvent == "notify" {
				continue // notify is handled via config.toml, not hooks.json
			}
			if _, exists := hooksMap[codexEvent]; exists {
				continue
			}
			hooksMap[codexEvent] = []codexHookGroup{
				{
					Hooks: []codexHookEntry{
						{
							Type:    "command",
							Command: fmt.Sprintf("yaah hook %s", event),
							Timeout: 10,
						},
					},
				},
			}
		}
	}

	if len(hooksMap) == 0 {
		return nil, nil
	}

	out := codexHooksConfig{Hooks: hooksMap}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal codex hooks.json: %w", err)
	}
	return data, nil
}

func (g *CodexGenerator) HooksPath() string    { return ".codex/hooks.json" }
func (g *CodexGenerator) SkillsDir() string    { return ".agents/skills" }
func (g *CodexGenerator) AgentsDir() string    { return "" }
func (g *CodexGenerator) AgentFileExt() string { return "" }
func (g *CodexGenerator) CommandsDir() string  { return "" }
