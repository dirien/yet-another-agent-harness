package hooks_test

import (
	"strings"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(t *testing.T, input *hooks.Input)
	}{
		{
			name: "valid input",
			json: `{"session_id":"abc","cwd":"/tmp","tool_name":"Edit","tool_input":{"file_path":"/tmp/foo.go"}}`,
			check: func(t *testing.T, input *hooks.Input) {
				t.Helper()
				if input.SessionID != "abc" {
					t.Errorf("SessionID: got %q, want %q", input.SessionID, "abc")
				}
				if input.Cwd != "/tmp" {
					t.Errorf("Cwd: got %q, want %q", input.Cwd, "/tmp")
				}
				if input.ToolName != "Edit" {
					t.Errorf("ToolName: got %q, want %q", input.ToolName, "Edit")
				}
			},
		},
		{
			name:    "invalid JSON",
			json:    `{broken`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := hooks.ReadInput(strings.NewReader(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, input)
			}
		})
	}
}

func TestInput_FilePath(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"with file_path", `{"tool_input":{"file_path":"/tmp/foo.go"}}`, "/tmp/foo.go"},
		{"empty tool_input", `{}`, ""},
		{"no file_path field", `{"tool_input":{"command":"ls"}}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := hooks.ReadInput(strings.NewReader(tt.json))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := input.FilePath(); got != tt.want {
				t.Errorf("FilePath(): got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInput_BashCommand(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"with command", `{"tool_input":{"command":"ls -la"}}`, "ls -la"},
		{"empty tool_input", `{}`, ""},
		{"no command field", `{"tool_input":{"file_path":"/tmp"}}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := hooks.ReadInput(strings.NewReader(tt.json))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := input.BashCommand(); got != tt.want {
				t.Errorf("BashCommand(): got %q, want %q", got, tt.want)
			}
		})
	}
}
