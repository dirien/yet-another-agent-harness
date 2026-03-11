package session

import "time"

// Session tracks state across hook events for a single Claude Code session.
type Session struct {
	ID            string           `json:"id"`
	StartedAt     time.Time        `json:"started_at"`
	LastEventAt   time.Time        `json:"last_event_at"`
	EventCount    int              `json:"event_count"`
	ToolCalls     []ToolCallRecord `json:"tool_calls"`
	BlockedCalls  []ToolCallRecord `json:"blocked_calls"`
	FilesModified []string         `json:"files_modified"`
	Findings      []Finding        `json:"findings"`
}

// ToolCallRecord records a single tool invocation within a session.
type ToolCallRecord struct {
	Timestamp time.Time `json:"timestamp"`
	ToolName  string    `json:"tool_name"`
	Input     string    `json:"input"` // summary, not full input
	Blocked   bool      `json:"blocked"`
	Reason    string    `json:"reason,omitempty"`
}

// Finding records a notable observation made by a handler (e.g. a detected secret).
type Finding struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "secret", "lint", "comment"
	File      string    `json:"file"`
	Line      int       `json:"line,omitempty"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // "error", "warning", "info"
}
