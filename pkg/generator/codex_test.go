package generator_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	toml "github.com/pelletier/go-toml/v2"
)

func TestCodex_GenerateSettings_TOML(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := gen.GenerateSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := toml.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid TOML: %v\ndata: %s", err, data)
	}

	// Verify notify field.
	notify, ok := out["notify"]
	if !ok {
		t.Fatal("expected notify field in TOML output")
	}
	notifySlice, ok := notify.([]any)
	if !ok || len(notifySlice) != 3 {
		t.Errorf("notify: got %v, want [yaah hook Notification]", notify)
	}

	// Verify features.codex_hooks.
	features, ok := out["features"].(map[string]any)
	if !ok {
		t.Fatal("expected features table in TOML output")
	}
	if features["codex_hooks"] != true {
		t.Errorf("features.codex_hooks: got %v, want true", features["codex_hooks"])
	}
}

func TestCodex_GenerateSettings_WithMCP(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	cfg := &schema.HarnessConfig{
		Version: "1",
		MCP: &schema.MCPConfig{
			Servers: []schema.MCPServer{
				{Name: "test-srv", Transport: schema.MCPTransportStdio, Command: "npx", Args: []string{"test"}},
			},
		},
	}
	data, err := gen.GenerateSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := toml.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid TOML: %v", err)
	}

	mcpServers, ok := out["mcp_servers"].(map[string]any)
	if !ok {
		t.Fatal("expected mcp_servers in TOML output")
	}
	if _, ok := mcpServers["test-srv"]; !ok {
		t.Error("expected test-srv in mcp_servers")
	}
}

func TestCodex_GenerateHooks(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	cfg := &schema.HarnessConfig{
		Version: "1",
		Hooks: schema.HooksConfig{
			schema.HookSessionStart: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook SessionStart"}}},
			},
			schema.HookStop: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook Stop"}}},
			},
		},
	}
	data, err := gen.GenerateHooks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil hooks data")
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	hooksMap, ok := out["hooks"].(map[string]any)
	if !ok {
		t.Fatal("expected hooks map in output")
	}
	if _, ok := hooksMap["SessionStart"]; !ok {
		t.Error("expected SessionStart in hooks")
	}
	if _, ok := hooksMap["Stop"]; !ok {
		t.Error("expected Stop in hooks")
	}

	content := string(data)
	if !strings.Contains(content, "yaah hook SessionStart") {
		t.Error("expected yaah hook SessionStart command")
	}
}

func TestCodex_GenerateHooks_NoSupported(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	cfg := &schema.HarnessConfig{
		Version: "1",
		Hooks: schema.HooksConfig{
			// PreToolUse is not supported by Codex.
			schema.HookPreToolUse: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook PreToolUse"}}},
			},
		},
	}
	data, err := gen.GenerateHooks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil hooks data for unsupported events")
	}
}

func TestCodex_Paths(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	if gen.SkillsDir() != ".agents/skills" {
		t.Errorf("SkillsDir: got %q", gen.SkillsDir())
	}
	if gen.AgentsDir() != "" {
		t.Errorf("AgentsDir: got %q, want empty", gen.AgentsDir())
	}
	if gen.CommandsDir() != "" {
		t.Errorf("CommandsDir: got %q, want empty", gen.CommandsDir())
	}
	if gen.MCPPath() != "" {
		t.Errorf("MCPPath: got %q, want empty", gen.MCPPath())
	}
}
