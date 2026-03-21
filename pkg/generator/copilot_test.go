package generator_test

import (
	"encoding/json"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func TestCopilot_GenerateSettings_Nil(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	data, err := gen.GenerateSettings(&schema.HarnessConfig{Version: "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil settings (Copilot has no settings file)")
	}
}

func TestCopilot_GenerateMCP(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	cfg := &schema.HarnessConfig{
		Version: "1",
		MCP: &schema.MCPConfig{
			Servers: []schema.MCPServer{
				{Name: "test-srv", Transport: schema.MCPTransportStdio, Command: "npx", Args: []string{"test"}},
			},
		},
	}
	data, err := gen.GenerateMCP(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil MCP data")
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	mcpServers, ok := out["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("expected mcpServers in output")
	}
	srv, ok := mcpServers["test-srv"].(map[string]any)
	if !ok {
		t.Fatal("expected test-srv in mcpServers")
	}
	if srv["type"] != "stdio" {
		t.Errorf("type: got %q, want %q", srv["type"], "stdio")
	}
}

func TestCopilot_GenerateMCP_EnvPassthrough(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	cfg := &schema.HarnessConfig{
		Version: "1",
		MCP: &schema.MCPConfig{
			Servers: []schema.MCPServer{
				{
					Name:      "with-env",
					Transport: schema.MCPTransportStdio,
					Command:   "node",
					Args:      []string{"server.js"},
					Env:       map[string]string{"API_KEY": "secret", "DEBUG": "true"},
				},
			},
		},
	}
	data, err := gen.GenerateMCP(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil MCP data")
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	mcpServers := out["mcpServers"].(map[string]any)
	srv := mcpServers["with-env"].(map[string]any)

	env, ok := srv["env"].(map[string]any)
	if !ok {
		t.Fatal("expected env map in server entry")
	}
	if env["API_KEY"] != "secret" {
		t.Errorf("env[API_KEY]: got %q, want %q", env["API_KEY"], "secret")
	}
	if env["DEBUG"] != "true" {
		t.Errorf("env[DEBUG]: got %q, want %q", env["DEBUG"], "true")
	}
}

func TestCopilot_GenerateMCP_Empty(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	data, err := gen.GenerateMCP(&schema.HarnessConfig{Version: "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil MCP data for empty config")
	}
}

func TestCopilot_GenerateHooks(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	cfg := &schema.HarnessConfig{
		Version: "1",
		Hooks: schema.HooksConfig{
			schema.HookSessionStart: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook SessionStart"}}},
			},
			schema.HookPreToolUse: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook PreToolUse"}}},
			},
			schema.HookPostToolUse: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook PostToolUse"}}},
			},
			schema.HookUserPromptSubmit: {
				{Hooks: []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook UserPromptSubmit"}}},
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

	// Check version field.
	version, ok := out["version"].(float64)
	if !ok || version != 1 {
		t.Errorf("version: got %v, want 1", out["version"])
	}

	hooksMap, ok := out["hooks"].(map[string]any)
	if !ok {
		t.Fatal("expected hooks map in output")
	}

	expectedEvents := []string{"sessionStart", "preToolUse", "postToolUse", "userPromptSubmitted"}
	for _, event := range expectedEvents {
		entries, ok := hooksMap[event].([]any)
		if !ok {
			t.Errorf("expected %q in hooks", event)
			continue
		}
		entry, ok := entries[0].(map[string]any)
		if !ok {
			t.Errorf("expected hook entry for %q", event)
			continue
		}
		if entry["type"] != "command" {
			t.Errorf("%s type: got %q, want %q", event, entry["type"], "command")
		}
		if bash, ok := entry["bash"].(string); !ok || bash == "" {
			t.Errorf("%s bash: got %q, want non-empty", event, entry["bash"])
		}
		if timeout, ok := entry["timeoutSec"].(float64); !ok || timeout != 30 {
			t.Errorf("%s timeoutSec: got %v, want 30", event, entry["timeoutSec"])
		}
	}
}

func TestCopilot_GenerateHooks_Empty(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	data, err := gen.GenerateHooks(&schema.HarnessConfig{Version: "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil hooks data for empty config")
	}
}

func TestCopilot_Paths(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	if gen.MCPPath() != ".copilot/mcp-config.json" {
		t.Errorf("MCPPath: got %q", gen.MCPPath())
	}
	if gen.HooksPath() != ".github/hooks/hooks.json" {
		t.Errorf("HooksPath: got %q", gen.HooksPath())
	}
	if gen.SkillsDir() != ".github/skills" {
		t.Errorf("SkillsDir: got %q", gen.SkillsDir())
	}
	if gen.AgentsDir() != ".github/agents" {
		t.Errorf("AgentsDir: got %q", gen.AgentsDir())
	}
	if gen.AgentFileExt() != ".agent.md" {
		t.Errorf("AgentFileExt: got %q", gen.AgentFileExt())
	}
	if gen.CommandsDir() != "" {
		t.Errorf("CommandsDir: got %q, want empty", gen.CommandsDir())
	}
}
