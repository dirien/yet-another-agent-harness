package hooks

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// Registry holds all registered handlers and dispatches events to them.
type Registry struct {
	handlers []Handler
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register enlists a handler into the registry.
func (c *Registry) Register(s Handler) {
	c.handlers = append(c.handlers, s)
}

// Dispatch finds all handlers matching the event and input, executes them,
// and returns the combined results.
func (c *Registry) Dispatch(ctx context.Context, event schema.HookEvent, input *Input) ([]*Result, error) {
	var results []*Result

	for _, s := range c.handlers {
		if !c.matchesEvent(s, event) {
			continue
		}

		if m := s.Match(); m != nil && input.ToolName != "" {
			if !m.MatchString(input.ToolName) {
				continue
			}
		}

		if fa, ok := s.(FileAware); ok {
			fp := input.FilePath()
			if fp != "" {
				ext := filepath.Ext(fp)
				exts := fa.FileExtensions()
				if exts != nil && !containsExt(exts, ext) {
					continue
				}
			}
		}

		result, err := s.Execute(ctx, input)
		if err != nil {
			return results, fmt.Errorf("handler %s: %w", s.Name(), err)
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

// Handlers returns all enlisted handlers.
func (c *Registry) Handlers() []Handler {
	return c.handlers
}

// CombineResults merges multiple results into a single Result for output.
func CombineResults(results []*Result) *Result {
	var outputs, errors []string
	block := false

	for _, v := range results {
		if v.Output != "" {
			outputs = append(outputs, v.Output)
		}
		if v.Error != "" {
			errors = append(errors, v.Error)
		}
		if v.Block {
			block = true
		}
	}

	return &Result{
		Output: strings.Join(outputs, "\n"),
		Error:  strings.Join(errors, "\n"),
		Block:  block,
	}
}

func (c *Registry) matchesEvent(s Handler, event schema.HookEvent) bool {
	for _, e := range s.Events() {
		if e == event {
			return true
		}
	}
	return false
}

func containsExt(exts []string, ext string) bool {
	for _, e := range exts {
		if e == ext {
			return true
		}
	}
	return false
}
