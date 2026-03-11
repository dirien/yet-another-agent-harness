package hooks_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// stubHandler is a minimal Handler for testing dispatch logic.
type stubHandler struct {
	name      string
	events    []schema.HookEvent
	match     *regexp.Regexp
	result    *hooks.Result
	execError error
	called    bool
}

func (s *stubHandler) Name() string               { return s.name }
func (s *stubHandler) Events() []schema.HookEvent { return s.events }
func (s *stubHandler) Match() *regexp.Regexp      { return s.match }
func (s *stubHandler) Execute(_ context.Context, _ *hooks.Input) (*hooks.Result, error) {
	s.called = true
	return s.result, s.execError
}

func TestRegistry_RegisterAndHandlers(t *testing.T) {
	r := hooks.NewRegistry()
	if got := len(r.Handlers()); got != 0 {
		t.Fatalf("expected 0 handlers, got %d", got)
	}

	h := &stubHandler{name: "test"}
	r.Register(h)

	if got := len(r.Handlers()); got != 1 {
		t.Fatalf("expected 1 handler, got %d", got)
	}
	if r.Handlers()[0].Name() != "test" {
		t.Fatalf("expected handler name 'test', got %q", r.Handlers()[0].Name())
	}
}

func TestRegistry_DispatchByEvent(t *testing.T) {
	tests := []struct {
		name       string
		events     []schema.HookEvent
		dispatch   schema.HookEvent
		wantCalled bool
	}{
		{"matching event", []schema.HookEvent{schema.HookPreToolUse}, schema.HookPreToolUse, true},
		{"non-matching event", []schema.HookEvent{schema.HookPreToolUse}, schema.HookPostToolUse, false},
		{"multiple events match", []schema.HookEvent{schema.HookPreToolUse, schema.HookPostToolUse}, schema.HookPostToolUse, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := hooks.NewRegistry()
			h := &stubHandler{
				name:   "h",
				events: tt.events,
				result: &hooks.Result{Output: "ok"},
			}
			r.Register(h)

			input := &hooks.Input{}
			_, err := r.Dispatch(context.Background(), tt.dispatch, input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h.called != tt.wantCalled {
				t.Fatalf("expected called=%v, got %v", tt.wantCalled, h.called)
			}
		})
	}
}

func TestRegistry_DispatchByToolName(t *testing.T) {
	r := hooks.NewRegistry()
	h := &stubHandler{
		name:   "edit-only",
		events: []schema.HookEvent{schema.HookPostToolUse},
		match:  regexp.MustCompile(`^Edit$`),
		result: &hooks.Result{Output: "matched"},
	}
	r.Register(h)

	t.Run("matching tool", func(t *testing.T) {
		input := &hooks.Input{ToolName: "Edit"}
		results, err := r.Dispatch(context.Background(), schema.HookPostToolUse, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("non-matching tool", func(t *testing.T) {
		h.called = false
		input := &hooks.Input{ToolName: "Bash"}
		results, err := r.Dispatch(context.Background(), schema.HookPostToolUse, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})
}

func TestCombineResults(t *testing.T) {
	tests := []struct {
		name      string
		results   []*hooks.Result
		wantBlock bool
		wantOut   string
		wantErr   string
	}{
		{
			name:    "empty",
			results: nil,
		},
		{
			name: "single output",
			results: []*hooks.Result{
				{Output: "hello"},
			},
			wantOut: "hello",
		},
		{
			name: "block propagates",
			results: []*hooks.Result{
				{Output: "ok"},
				{Error: "bad", Block: true},
			},
			wantBlock: true,
			wantOut:   "ok",
			wantErr:   "bad",
		},
		{
			name: "multiple outputs joined",
			results: []*hooks.Result{
				{Output: "a"},
				{Output: "b"},
			},
			wantOut: "a\nb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			combined := hooks.CombineResults(tt.results)
			if combined.Block != tt.wantBlock {
				t.Errorf("Block: got %v, want %v", combined.Block, tt.wantBlock)
			}
			if combined.Output != tt.wantOut {
				t.Errorf("Output: got %q, want %q", combined.Output, tt.wantOut)
			}
			if combined.Error != tt.wantErr {
				t.Errorf("Error: got %q, want %q", combined.Error, tt.wantErr)
			}
		})
	}
}
