package schema

// HookEvent enumerates all Claude Code hook lifecycle events.
type HookEvent string

const (
	HookSessionStart      HookEvent = "SessionStart"
	HookSessionEnd        HookEvent = "SessionEnd"
	HookInstructionsLoad  HookEvent = "InstructionsLoaded"
	HookUserPromptSubmit  HookEvent = "UserPromptSubmit"
	HookPreToolUse        HookEvent = "PreToolUse"
	HookPostToolUse       HookEvent = "PostToolUse"
	HookPostToolUseFail   HookEvent = "PostToolUseFailure"
	HookPermissionRequest HookEvent = "PermissionRequest"
	HookNotification      HookEvent = "Notification"
	HookSubagentStart     HookEvent = "SubagentStart"
	HookSubagentStop      HookEvent = "SubagentStop"
	HookStop              HookEvent = "Stop"
	HookTeammateIdle      HookEvent = "TeammateIdle"
	HookTaskCompleted     HookEvent = "TaskCompleted"
	HookConfigChange      HookEvent = "ConfigChange"
	HookWorktreeCreate    HookEvent = "WorktreeCreate"
	HookWorktreeRemove    HookEvent = "WorktreeRemove"
	HookPreCompact        HookEvent = "PreCompact"
)

// AllHookEvents returns every supported hook event.
func AllHookEvents() []HookEvent {
	return []HookEvent{
		HookSessionStart, HookSessionEnd, HookInstructionsLoad,
		HookUserPromptSubmit, HookPreToolUse, HookPostToolUse,
		HookPostToolUseFail, HookPermissionRequest, HookNotification,
		HookSubagentStart, HookSubagentStop, HookStop,
		HookTeammateIdle, HookTaskCompleted, HookConfigChange,
		HookWorktreeCreate, HookWorktreeRemove, HookPreCompact,
	}
}

// HookType defines the handler type for a hook.
type HookType string

const (
	HookTypeCommand HookType = "command"
	HookTypeHTTP    HookType = "http"
	HookTypePrompt  HookType = "prompt"
	HookTypeAgent   HookType = "agent"
)

// HookHandler is a single hook action (command, http, prompt, or agent).
type HookHandler struct {
	Type           HookType          `json:"type"                      jsonschema:"enum=command,enum=http,enum=prompt,enum=agent"`
	Command        string            `json:"command,omitempty"         jsonschema:"description=Shell command to run (type=command)"`
	URL            string            `json:"url,omitempty"             jsonschema:"description=HTTP endpoint to POST to (type=http)"`
	Prompt         string            `json:"prompt,omitempty"          jsonschema:"description=LLM prompt to evaluate (type=prompt or agent)"`
	Timeout        int               `json:"timeout,omitempty"         jsonschema:"description=Timeout in seconds (defaults: command=600 prompt=30 agent=60)"`
	StatusMessage  string            `json:"statusMessage,omitempty"   jsonschema:"description=Custom spinner message displayed while the hook runs"`
	Once           bool              `json:"once,omitempty"            jsonschema:"description=Run only once per session (skills only)"`
	Async          bool              `json:"async,omitempty"           jsonschema:"description=Run in background without blocking (type=command only)"`
	Headers        map[string]string `json:"headers,omitempty"         jsonschema:"description=HTTP headers to include (type=http only)"`
	AllowedEnvVars []string          `json:"allowedEnvVars,omitempty"  jsonschema:"description=Environment variable names allowed in header interpolation (type=http only)"`
	Model          string            `json:"model,omitempty"           jsonschema:"description=Model identifier for LLM evaluation (type=prompt or agent)"`
}

// HookRule binds a matcher pattern to one or more handlers.
type HookRule struct {
	Matcher string        `json:"matcher,omitempty" jsonschema:"description=Regex to filter when this rule fires (tool name or event subtype)"`
	Hooks   []HookHandler `json:"hooks"             jsonschema:"description=Handlers to execute when matched"`
}

// HooksConfig maps hook events to their rules.
type HooksConfig map[HookEvent][]HookRule
