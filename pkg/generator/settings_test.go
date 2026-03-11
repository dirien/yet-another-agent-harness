package generator_test

import (
	"encoding/json"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func TestGenerateClaudeSettings_Empty(t *testing.T) {
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := generator.GenerateClaudeSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if _, ok := out["$schema"]; !ok {
		t.Error("expected $schema field in output")
	}
}

func TestGenerateClaudeSettings_WithModel(t *testing.T) {
	cfg := &schema.HarnessConfig{
		Version:  "1",
		Settings: &schema.Settings{Model: "opus"},
	}
	data, err := generator.GenerateClaudeSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out generator.ClaudeSettings
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out.Model != "opus" {
		t.Errorf("Model: got %q, want %q", out.Model, "opus")
	}
}

func TestGenerateClaudeSettings_WithHooks(t *testing.T) {
	cfg := &schema.HarnessConfig{
		Version: "1",
		Hooks: schema.HooksConfig{
			schema.HookPostToolUse: {
				{
					Matcher: `^Edit$`,
					Hooks:   []schema.HookHandler{{Type: schema.HookTypeCommand, Command: "yaah hook PostToolUse"}},
				},
			},
		},
	}
	data, err := generator.GenerateClaudeSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	hooksMap, ok := out["hooks"].(map[string]any)
	if !ok {
		t.Fatal("expected hooks map in output")
	}
	if _, ok := hooksMap["PostToolUse"]; !ok {
		t.Error("expected PostToolUse key in hooks")
	}
}

func TestGenerateClaudeSettings_WithMCP(t *testing.T) {
	cfg := &schema.HarnessConfig{
		Version: "1",
		MCP: &schema.MCPConfig{
			Servers: []schema.MCPServer{
				{Name: "test-server", Transport: schema.MCPTransportStdio, Command: "npx", Args: []string{"test"}},
			},
		},
	}
	data, err := generator.GenerateClaudeSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	mcpMap, ok := out["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("expected mcpServers map in output")
	}
	if _, ok := mcpMap["test-server"]; !ok {
		t.Error("expected test-server key in mcpServers")
	}
}

func TestGenerateLSPConfig_Nil(t *testing.T) {
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := generator.GenerateLSPConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil data for empty LSP config")
	}
}

func TestGenerateLSPConfig_WithServers(t *testing.T) {
	cfg := &schema.HarnessConfig{
		Version: "1",
		LSP: &schema.LSPConfig{
			Servers: []schema.LSPServer{
				{ID: "gopls", Command: []string{"gopls", "serve"}, ExtensionToLanguage: map[string]string{".go": "go"}},
			},
		},
	}
	data, err := generator.GenerateLSPConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	goplsEntry, ok := out["gopls"].(map[string]any)
	if !ok {
		t.Fatal("expected gopls key in output")
	}
	cmd, ok := goplsEntry["command"].([]any)
	if !ok || len(cmd) != 2 {
		t.Errorf("expected command [gopls serve], got %v", goplsEntry["command"])
	}
}
