package hooks

import (
	"context"
	"regexp"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ Handler = (*Chain)(nil)

// ChainFunc is a function that receives context, input, and the previous handler's result,
// and returns a modified result.
type ChainFunc func(ctx context.Context, input *Input, prev *Result) (*Result, error)

// ChainLink connects a handler or function in the chain.
// Set either Handler or Func (not both). Condition is optional and
// gates whether the link executes based on the previous result.
type ChainLink struct {
	Handler   Handler
	Func      ChainFunc
	Condition func(*Result) bool
}

// Chain implements the Handler interface by running links sequentially.
type Chain struct {
	name   string
	events []schema.HookEvent
	match  *regexp.Regexp
	links  []ChainLink
}

// NewChain creates a new Chain with the given name, events, and matcher.
func NewChain(name string, events []schema.HookEvent, match *regexp.Regexp, links ...ChainLink) *Chain {
	return &Chain{
		name:   name,
		events: events,
		match:  match,
		links:  links,
	}
}

// Name returns the chain's identifier.
func (c *Chain) Name() string { return c.name }

// Events returns the hook events this chain listens to.
// If events were explicitly set, those are returned; otherwise the union
// of all link handler events is returned.
func (c *Chain) Events() []schema.HookEvent {
	if len(c.events) > 0 {
		return c.events
	}

	seen := make(map[schema.HookEvent]bool)
	var union []schema.HookEvent
	for _, link := range c.links {
		if link.Handler == nil {
			continue
		}
		for _, ev := range link.Handler.Events() {
			if !seen[ev] {
				seen[ev] = true
				union = append(union, ev)
			}
		}
	}
	return union
}

// Match returns the chain's compiled regex matcher.
func (c *Chain) Match() *regexp.Regexp { return c.match }

// Execute runs all links in sequence, passing each result to the next link.
func (c *Chain) Execute(ctx context.Context, input *Input) (*Result, error) {
	var prev *Result

	for _, link := range c.links {
		// If link has a condition and it returns false, skip this link.
		if link.Condition != nil && !link.Condition(prev) {
			continue
		}

		var (
			result *Result
			err    error
		)

		switch {
		case link.Handler != nil:
			result, err = link.Handler.Execute(ctx, input)
		case link.Func != nil:
			result, err = link.Func(ctx, input, prev)
		default:
			// No handler or func -- skip.
			continue
		}

		if err != nil {
			return prev, err
		}
		if result != nil {
			prev = result
		}
	}

	return prev, nil
}
