package harness_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/harness"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// blockingHandler always blocks.
type blockingHandler struct{}

func (b *blockingHandler) Name() string { return "blocker" }
func (b *blockingHandler) Events() []schema.HookEvent {
	return []schema.HookEvent{schema.HookPreToolUse}
}
func (b *blockingHandler) Match() *regexp.Regexp { return nil }
func (b *blockingHandler) Execute(_ context.Context, _ *hooks.Input) (*hooks.Result, error) {
	return &hooks.Result{Error: "blocked!", Block: true}, nil
}

func TestHarness_New(t *testing.T) {
	h := harness.New()
	if h.Hooks() == nil {
		t.Fatal("Hooks() returned nil")
	}
	if h.MCP() == nil {
		t.Fatal("MCP() returned nil")
	}
	if h.LSP() == nil {
		t.Fatal("LSP() returned nil")
	}
	if h.Skills() == nil {
		t.Fatal("Skills() returned nil")
	}
	if h.Agents() == nil {
		t.Fatal("Agents() returned nil")
	}
	if h.Commands() == nil {
		t.Fatal("Commands() returned nil")
	}
	if h.Plugins() == nil {
		t.Fatal("Plugins() returned nil")
	}
}

func TestHarness_HandleHookEvent_Block(t *testing.T) {
	h := harness.New()
	h.Hooks().Register(&blockingHandler{})

	err := h.HandleHookEvent(context.Background(), schema.HookPreToolUse, &hooks.Input{})
	if !errors.Is(err, harness.ErrHookBlocked) {
		t.Fatalf("expected ErrHookBlocked, got %v", err)
	}
}

func TestHarness_HandleHookEvent_NoBlock(t *testing.T) {
	h := harness.New()
	err := h.HandleHookEvent(context.Background(), schema.HookPreToolUse, &hooks.Input{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHarness_GenerateConfig(t *testing.T) {
	h := harness.New()
	h.SetSettings(&schema.Settings{Model: "opus"})

	cfg := h.GenerateConfig()
	if cfg.Version != "1" {
		t.Errorf("Version: got %q, want %q", cfg.Version, "1")
	}
	if cfg.Settings == nil {
		t.Fatal("Settings is nil")
	}
	if cfg.Settings.Model != "opus" {
		t.Errorf("Model: got %q, want %q", cfg.Settings.Model, "opus")
	}
}

func TestHarness_Summary(t *testing.T) {
	h := harness.New()
	s := h.Summary()
	if s == "" {
		t.Fatal("Summary() returned empty string")
	}
}

func TestHarness_DefaultsIncludeWorkflowCommands(t *testing.T) {
	h := harness.NewWithDefaults(harness.AllDefaults())
	cmds := h.Commands().Commands()
	if len(cmds) < 29 {
		t.Errorf("expected at least 29 commands, got %d", len(cmds))
	}
	names := make(map[string]bool)
	for _, c := range cmds {
		names[c.Name()] = true
	}
	for _, want := range []string{
		"yaah/init", "yaah/discuss", "yaah/plan", "yaah/execute", "yaah/verify", "yaah/docs", "yaah/next", "yaah/quick",
		"yaah/ship", "yaah/pause", "yaah/resume", "yaah/complete-milestone", "yaah/new-milestone", "yaah/settings",
		"yaah/add-phase", "yaah/insert-phase", "yaah/remove-phase", "yaah/health", "yaah/progress", "yaah/review",
		"yaah/secure", "yaah/todo", "yaah/note", "yaah/cleanup", "yaah/forensics", "yaah/explore", "yaah/scan", "yaah/import", "yaah/autonomous",
	} {
		if !names[want] {
			t.Errorf("missing command %q", want)
		}
	}
}
