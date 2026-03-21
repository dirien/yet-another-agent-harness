package generator

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// HookMapping maps a Claude hook event to the equivalent event name for each agent.
// An empty string means the event is not supported by that agent.
type HookMapping struct {
	Claude   string
	OpenCode string // JS plugin event name
	Codex    string // hooks.json event name
	Copilot  string // hooks.json event name
}

// hookMappings defines the cross-agent event name mapping.
var hookMappings = map[schema.HookEvent]HookMapping{
	schema.HookSessionStart: {
		Claude:   "SessionStart",
		OpenCode: "session.created",
		Codex:    "SessionStart",
		Copilot:  "sessionStart",
	},
	schema.HookSessionEnd: {
		Claude:  "SessionEnd",
		Copilot: "sessionEnd",
	},
	schema.HookUserPromptSubmit: {
		Claude:  "UserPromptSubmit",
		Copilot: "userPromptSubmitted",
	},
	schema.HookPreToolUse: {
		Claude:   "PreToolUse",
		OpenCode: "tool.execute.before",
		Copilot:  "preToolUse",
	},
	schema.HookPostToolUse: {
		Claude:   "PostToolUse",
		OpenCode: "tool.execute.after",
		Copilot:  "postToolUse",
	},
	schema.HookStop: {
		Claude:   "Stop",
		OpenCode: "stop",
		Codex:    "Stop",
		Copilot:  "agentStop",
	},
	schema.HookNotification: {
		Claude: "Notification",
		Codex:  "notify",
	},
	schema.HookSubagentStop: {
		Claude:  "SubagentStop",
		Copilot: "subagentStop",
	},
}

// CopilotEventName returns the Copilot hooks.json event name for a Claude event.
// Returns empty string if unsupported.
func CopilotEventName(event schema.HookEvent) string {
	if m, ok := hookMappings[event]; ok {
		return m.Copilot
	}
	return ""
}

// CodexEventName returns the Codex hooks.json event name for a Claude event.
// Returns empty string if unsupported.
func CodexEventName(event schema.HookEvent) string {
	if m, ok := hookMappings[event]; ok {
		return m.Codex
	}
	return ""
}

// OpenCodeEventName returns the OpenCode JS plugin event name for a Claude event.
// Returns empty string if unsupported.
func OpenCodeEventName(event schema.HookEvent) string {
	if m, ok := hookMappings[event]; ok {
		return m.OpenCode
	}
	return ""
}
