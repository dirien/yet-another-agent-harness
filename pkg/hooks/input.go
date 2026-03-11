package hooks

import (
	"encoding/json"
	"fmt"
	"io"
)

// Input is the JSON payload Claude Code passes on stdin to hook commands.
type Input struct {
	SessionID      string          `json:"session_id"`
	TranscriptPath string          `json:"transcript_path"`
	Cwd            string          `json:"cwd"`
	PermissionMode string          `json:"permission_mode"`
	HookEventName  string          `json:"hook_event_name"`
	ToolName       string          `json:"tool_name,omitempty"`
	ToolInput      json.RawMessage `json:"tool_input,omitempty"`
	ToolResponse   json.RawMessage `json:"tool_response,omitempty"`
}

// toolInputFile extracts the file_path from tool_input JSON.
type toolInputFile struct {
	FilePath string `json:"file_path"`
}

// toolInputBash extracts the command from Bash tool_input JSON.
type toolInputBash struct {
	Command string `json:"command"`
}

// ReadInput parses the hook input from the given reader (typically os.Stdin).
func ReadInput(r io.Reader) (*Input, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}
	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("parse hook input: %w", err)
	}
	return &input, nil
}

// FilePath extracts the file_path from the tool_input JSON.
func (h *Input) FilePath() string {
	if len(h.ToolInput) == 0 {
		return ""
	}
	var ti toolInputFile
	if err := json.Unmarshal(h.ToolInput, &ti); err != nil {
		return ""
	}
	return ti.FilePath
}

// BashCommand extracts the command string from Bash tool_input JSON.
func (h *Input) BashCommand() string {
	if len(h.ToolInput) == 0 {
		return ""
	}
	var ti toolInputBash
	if err := json.Unmarshal(h.ToolInput, &ti); err != nil {
		return ""
	}
	return ti.Command
}
