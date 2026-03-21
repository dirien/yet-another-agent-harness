package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ hooks.Handler = (*CommandGuard)(nil)

// CommandGuard is a PreToolUse handler that blocks dangerous shell commands.
// Implements: Handler.
type CommandGuard struct {
	blocked []*guardRule
}

type guardRule struct {
	pattern *regexp.Regexp
	reason  string
}

// NewCommandGuard creates a CommandGuard with sensible default blocked patterns.
func NewCommandGuard() *CommandGuard {
	c := &CommandGuard{}
	c.Block(`rm\s+-rf\s+/`, "recursive delete from root")
	c.Block(`git\s+push\s+--force\s+(origin\s+)?(main|master)`, "force push to main/master")
	c.Block(`git\s+reset\s+--hard`, "hard reset discards work")
	c.Block(`DROP\s+(TABLE|DATABASE)`, "destructive SQL statement")
	c.Block(`:\s*>\s*/`, "truncate system file")
	c.Block(`mkfs\.`, "format filesystem")
	c.Block(`dd\s+if=.*of=/dev/`, "raw disk write")
	return c
}

// Block adds a pattern that will be rejected when matched against Bash commands.
func (c *CommandGuard) Block(pattern string, reason string) {
	c.blocked = append(c.blocked, &guardRule{
		pattern: regexp.MustCompile(`(?i)` + pattern),
		reason:  reason,
	})
}

func (c *CommandGuard) Name() string { return "command-guard" }

func (c *CommandGuard) Events() []schema.HookEvent {
	return []schema.HookEvent{schema.HookPreToolUse}
}

var bashMatch = regexp.MustCompile(`(?i)^Bash$`)

func (c *CommandGuard) Match() *regexp.Regexp {
	return bashMatch
}

// CommandCheckResult describes the outcome of checking a command.
type CommandCheckResult struct {
	Safe   bool   `json:"safe"`
	Reason string `json:"reason"`
}

// CheckCommand checks whether the given command string is safe.
// Returns a result indicating whether the command is safe and why.
func (c *CommandGuard) CheckCommand(cmd string) CommandCheckResult {
	for _, rule := range c.blocked {
		if rule.pattern.MatchString(cmd) {
			return CommandCheckResult{
				Safe:   false,
				Reason: fmt.Sprintf("BLOCKED: %s (matched: %s)", rule.reason, rule.pattern.String()),
			}
		}
	}
	return CommandCheckResult{Safe: true, Reason: "command is safe"}
}

func (c *CommandGuard) Execute(_ context.Context, input *hooks.Input) (*hooks.Result, error) {
	cmd := input.BashCommand()
	if cmd == "" {
		return nil, nil
	}

	var violations []string
	for _, rule := range c.blocked {
		if rule.pattern.MatchString(cmd) {
			violations = append(violations, fmt.Sprintf("BLOCKED: %s (matched: %s)", rule.reason, rule.pattern.String()))
		}
	}

	if len(violations) > 0 {
		return &hooks.Result{
			Error: strings.Join(violations, "\n"),
			Block: true,
		}, nil
	}
	return nil, nil
}
