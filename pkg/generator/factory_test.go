package generator_test

import (
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func TestForTarget_Claude(t *testing.T) {
	gen := generator.ForTarget(schema.TargetClaude)
	if gen.Target() != schema.TargetClaude {
		t.Errorf("Target: got %q, want %q", gen.Target(), schema.TargetClaude)
	}
	if gen.SettingsPath() != ".claude/settings.json" {
		t.Errorf("SettingsPath: got %q", gen.SettingsPath())
	}
}

func TestForTarget_OpenCode(t *testing.T) {
	gen := generator.ForTarget(schema.TargetOpenCode)
	if gen.Target() != schema.TargetOpenCode {
		t.Errorf("Target: got %q, want %q", gen.Target(), schema.TargetOpenCode)
	}
	if gen.SettingsPath() != "opencode.json" {
		t.Errorf("SettingsPath: got %q", gen.SettingsPath())
	}
}

func TestForTarget_Codex(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCodex)
	if gen.Target() != schema.TargetCodex {
		t.Errorf("Target: got %q, want %q", gen.Target(), schema.TargetCodex)
	}
	if gen.SettingsPath() != ".codex/config.toml" {
		t.Errorf("SettingsPath: got %q", gen.SettingsPath())
	}
}

func TestForTarget_Copilot(t *testing.T) {
	gen := generator.ForTarget(schema.TargetCopilot)
	if gen.Target() != schema.TargetCopilot {
		t.Errorf("Target: got %q, want %q", gen.Target(), schema.TargetCopilot)
	}
	if gen.SettingsPath() != "" {
		t.Errorf("SettingsPath: got %q, want empty", gen.SettingsPath())
	}
}

func TestForTarget_AllGenerators(t *testing.T) {
	for _, target := range schema.AllTargets() {
		gen := generator.ForTarget(target)
		if gen == nil {
			t.Errorf("ForTarget(%q) returned nil", target)
		}
		if gen.Target() != target {
			t.Errorf("ForTarget(%q).Target() = %q", target, gen.Target())
		}
	}
}
