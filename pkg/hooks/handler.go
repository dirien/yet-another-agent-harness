package hooks

import (
	"context"
	"regexp"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// Result is what a handler returns after execution.
type Result struct {
	// Output is shown to Claude on stdout (informational).
	Output string
	// Error is shown to Claude on stderr (problems to fix).
	Error string
	// Block signals exit code 2 — blocks the action (PreToolUse) or signals failure (PostToolUse).
	Block bool
}

// Handler is the core interface every hook implementation must satisfy.
type Handler interface {
	// Name returns a unique identifier for this handler.
	Name() string

	// Events returns which hook lifecycle events this handler listens to.
	Events() []schema.HookEvent

	// Match returns a compiled regex that filters tool names (or event subtypes).
	// Return nil to match all events of the subscribed types.
	Match() *regexp.Regexp

	// Execute runs the handler's logic. Input comes from Claude Code stdin.
	// Return a Result to communicate back to Claude.
	Execute(ctx context.Context, input *Input) (*Result, error)
}

// Configurable is an optional interface handlers can implement
// to accept configuration from yaah.json.
type Configurable interface {
	// Configure is called once at startup with the raw JSON config block.
	Configure(raw map[string]any) error
}

// FileAware is a convenience interface for handlers that operate on files.
// The registry calls FileFilter before Execute to short-circuit early.
type FileAware interface {
	// FileExtensions returns the file extensions this handler cares about (e.g., ".py", ".go").
	// Return nil to accept all files.
	FileExtensions() []string
}
