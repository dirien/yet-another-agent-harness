package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var (
	_ hooks.Handler   = (*CommentChecker)(nil)
	_ hooks.FileAware = (*CommentChecker)(nil)
)

// CommentChecker is a PostToolUse handler that detects non-English or
// AI-generated comments (like "TODO: implement" placeholders).
// Implements: Handler, FileAware.
type CommentChecker struct {
	patterns []*regexp.Regexp
}

// NewCommentChecker creates a comment checker with default patterns for lazy/placeholder comments.
func NewCommentChecker() *CommentChecker {
	return &CommentChecker{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)//\s*TODO:?\s*(implement|add|fix|write|handle)\s*(this|here|later)?\.?\s*$`),
			regexp.MustCompile(`(?i)#\s*TODO:?\s*(implement|add|fix|write|handle)\s*(this|here|later)?\.?\s*$`),
			regexp.MustCompile(`(?i)(//|#)\s*FIXME\s*$`),
			regexp.MustCompile(`(?i)(//|#)\s*HACK\s*$`),
			regexp.MustCompile(`(?i)(//|#)\s*\.\.\.`),
			regexp.MustCompile(`(?i)(//|#)\s*placeholder`),
			regexp.MustCompile(`(?i)(//|#)\s*your code here`),
		},
	}
}

// AddPattern adds a custom regex pattern to flag.
func (c *CommentChecker) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile pattern: %w", err)
	}
	c.patterns = append(c.patterns, re)
	return nil
}

func (c *CommentChecker) Name() string { return "comment-checker" }

func (c *CommentChecker) Events() []schema.HookEvent {
	return []schema.HookEvent{schema.HookPostToolUse}
}

func (c *CommentChecker) Match() *regexp.Regexp {
	return editWriteMatch
}

func (c *CommentChecker) FileExtensions() []string {
	return []string{".py", ".go", ".ts", ".tsx", ".js", ".jsx", ".rs", ".java", ".kt"}
}

func (c *CommentChecker) Execute(_ context.Context, input *hooks.Input) (*hooks.Result, error) {
	fp := input.FilePath()
	if fp == "" {
		return nil, nil
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, nil
	}
	defer func() { _ = f.Close() }()

	var findings []string
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		for _, p := range c.patterns {
			if p.MatchString(trimmed) {
				findings = append(findings, fmt.Sprintf("  %s:%d: placeholder comment: %s", fp, lineNum, trimmed))
				break
			}
		}

		if isComment(trimmed) {
			commentText := extractCommentText(trimmed)
			if hasNonLatin(commentText) {
				findings = append(findings, fmt.Sprintf("  %s:%d: non-English comment: %s", fp, lineNum, trimmed))
			}
		}
	}

	if len(findings) > 0 {
		return &hooks.Result{
			Error: fmt.Sprintf("comment-checker found %d issue(s):\n%s", len(findings), strings.Join(findings, "\n")),
		}, nil
	}
	return nil, nil
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#")
}

func extractCommentText(line string) string {
	if strings.HasPrefix(line, "//") {
		return strings.TrimSpace(line[2:])
	}
	if strings.HasPrefix(line, "#") {
		return strings.TrimSpace(line[1:])
	}
	return line
}

func hasNonLatin(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
