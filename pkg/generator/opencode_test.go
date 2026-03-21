package generator_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func TestOpenCode_GenerateSettings_Empty(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := gen.GenerateSettings(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestOpenCode_GenerateSettings_WithMCP(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
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
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	mcpServers, ok := out["mcp"].(map[string]any)
	if !ok {
		t.Fatal("expected mcp in output")
	}
	srv, ok := mcpServers["test-srv"].(map[string]any)
	if !ok {
		t.Fatal("expected test-srv in mcp")
	}
	if srv["type"] != "local" {
		t.Errorf("type: got %q, want %q", srv["type"], "local")
	}
	cmd, ok := srv["command"].([]any)
	if !ok {
		t.Fatal("expected command to be an array")
	}
	if len(cmd) != 2 || cmd[0] != "npx" || cmd[1] != "test" {
		t.Errorf("command: got %v, want [npx test]", cmd)
	}
	if _, hasOldEnv := srv["env"]; hasOldEnv {
		t.Error("found legacy 'env' key, should be 'environment'")
	}
}

func TestOpenCode_GenerateHooks_Plugin(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := gen.GenerateHooks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil hooks data (JS plugin)")
	}

	content := string(data)
	if !strings.Contains(content, "YaahPlugin") {
		t.Error("expected YaahPlugin export in JS plugin")
	}
	if !strings.Contains(content, "tool.execute.before") {
		t.Error("expected tool.execute.before event in JS plugin")
	}
	if !strings.Contains(content, "tool.execute.after") {
		t.Error("expected tool.execute.after event in JS plugin")
	}
	if !strings.Contains(content, "session.created") {
		t.Error("expected session.created event in JS plugin")
	}
	if !strings.Contains(content, `"hook", "PreToolUse"`) {
		t.Error("expected yaah hook PreToolUse command in JS plugin")
	}
}

func TestOpenCode_MCPEmbedded(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
	cfg := &schema.HarnessConfig{Version: "1"}
	data, err := gen.GenerateMCP(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil MCP data (embedded in settings)")
	}
}

func TestOpenCode_Paths(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
	if gen.SkillsDir() != ".opencode/skills" {
		t.Errorf("SkillsDir: got %q", gen.SkillsDir())
	}
	if gen.AgentsDir() != ".opencode/agents" {
		t.Errorf("AgentsDir: got %q", gen.AgentsDir())
	}
	if gen.AgentFileExt() != ".md" {
		t.Errorf("AgentFileExt: got %q", gen.AgentFileExt())
	}
	if gen.CommandsDir() != ".opencode/commands" {
		t.Errorf("CommandsDir: got %q", gen.CommandsDir())
	}
	if gen.HooksPath() != ".opencode/plugins/yaah.js" {
		t.Errorf("HooksPath: got %q", gen.HooksPath())
	}
}
