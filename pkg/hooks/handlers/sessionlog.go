package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ hooks.Handler = (*SessionLogger)(nil)

// SessionLogger logs session start/stop events with timestamps to a file.
// Implements: Handler.
type SessionLogger struct {
	logDir string
}

// NewSessionLogger creates a session logger that writes to the given directory.
// Defaults to .claude/logs/ relative to cwd if dir is empty.
func NewSessionLogger(dir string) *SessionLogger {
	if dir == "" {
		dir = filepath.Join(".claude", "logs")
	}
	return &SessionLogger{logDir: dir}
}

func (a *SessionLogger) Name() string { return "session-logger" }

func (a *SessionLogger) Events() []schema.HookEvent {
	return []schema.HookEvent{
		schema.HookSessionStart,
		schema.HookSessionEnd,
		schema.HookStop,
	}
}

func (a *SessionLogger) Match() *regexp.Regexp { return nil }

func (a *SessionLogger) Execute(_ context.Context, input *hooks.Input) (*hooks.Result, error) {
	logDir := a.logDir
	if !filepath.IsAbs(logDir) && input.Cwd != "" {
		logDir = filepath.Join(input.Cwd, logDir)
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	logFile := filepath.Join(logDir, "sessions.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	defer func() { _ = f.Close() }()

	ts := time.Now().UTC().Format(time.RFC3339)
	entry := fmt.Sprintf("[%s] event=%s session=%s cwd=%s\n",
		ts, input.HookEventName, input.SessionID, input.Cwd)

	if _, err := f.WriteString(entry); err != nil {
		return nil, fmt.Errorf("write log: %w", err)
	}

	return nil, nil
}
