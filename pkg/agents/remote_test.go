package agents

import (
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/gitcache"
)

func TestStripFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no frontmatter",
			input: "# Hello\n\nSome content.",
			want:  "# Hello\n\nSome content.",
		},
		{
			name:  "with frontmatter",
			input: "---\nname: test\nmodel: sonnet\n---\n\n# Hello\n\nSome content.",
			want:  "# Hello\n\nSome content.",
		},
		{
			name:  "frontmatter only",
			input: "---\nname: test\n---\n",
			want:  "",
		},
		{
			name:  "unclosed frontmatter",
			input: "---\nname: test\n# Hello",
			want:  "---\nname: test\n# Hello",
		},
		{
			name:  "empty content after frontmatter",
			input: "---\nname: test\n---\n\n",
			want:  "",
		},
		{
			name:  "crlf frontmatter",
			input: "---\r\nname: test\r\n---\r\n\r\n# Hello\r\n",
			want:  "# Hello\n",
		},
		{
			name:  "no frontmatter with CRLF",
			input: "# Hello\r\nContent\r\n",
			want:  "# Hello\nContent\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripFrontmatter(tt.input)
			if got != tt.want {
				t.Errorf("stripFrontmatter() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewRemoteAgent(t *testing.T) {
	a := NewRemoteAgent(
		"test-agent", "A test agent",
		"github.com/example/repo@v1.0.0", "agents/test.md",
		WithModel("opus"),
		WithTools("Read, Grep"),
	)

	if a.Name() != "test-agent" {
		t.Errorf("Name() = %q, want %q", a.Name(), "test-agent")
	}
	if a.Description() != "A test agent" {
		t.Errorf("Description() = %q, want %q", a.Description(), "A test agent")
	}
	if a.Model() != "opus" {
		t.Errorf("Model() = %q, want %q", a.Model(), "opus")
	}
	if a.Tools() != "Read, Grep" {
		t.Errorf("Tools() = %q, want %q", a.Tools(), "Read, Grep")
	}
	if a.Uses() != "github.com/example/repo@v1.0.0" {
		t.Errorf("Uses() = %q, want %q", a.Uses(), "github.com/example/repo@v1.0.0")
	}
	if a.Subpath() != "agents/test.md" {
		t.Errorf("Subpath() = %q, want %q", a.Subpath(), "agents/test.md")
	}
}

func TestNewRemoteAgent_WithAdvanced(t *testing.T) {
	a := NewRemoteAgent(
		"test-agent", "A test agent",
		"github.com/example/repo@v1.0.0", "agents/test.md",
		WithAdvanced(AgentAdvanced{
			PermissionMode: "acceptEdits",
			MaxTurns:       10,
		}),
	)

	adv := a.Advanced()
	if adv.PermissionMode != "acceptEdits" {
		t.Errorf("Advanced().PermissionMode = %q, want %q", adv.PermissionMode, "acceptEdits")
	}
	if adv.MaxTurns != 10 {
		t.Errorf("Advanced().MaxTurns = %d, want %d", adv.MaxTurns, 10)
	}
}

func TestNewRemoteAgent_DefaultAdvanced(t *testing.T) {
	a := NewRemoteAgent(
		"test-agent", "A test agent",
		"github.com/example/repo@v1.0.0", "agents/test.md",
	)

	adv := a.Advanced()
	if adv.PermissionMode != "" {
		t.Errorf("Advanced().PermissionMode = %q, want empty", adv.PermissionMode)
	}
}

func TestParseUses(t *testing.T) {
	tests := []struct {
		uses    string
		wantURL string
		wantRef string
		wantErr bool
	}{
		{"github.com/owner/repo@v1.0.0", "https://github.com/owner/repo.git", "v1.0.0", false},
		{"github.com/owner/repo@abc123", "https://github.com/owner/repo.git", "abc123", false},
		{"github.com/owner/repo", "", "", true},
		{"github.com/owner/repo@", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.uses, func(t *testing.T) {
			url, ref, err := gitcache.ParseUses(tt.uses)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseUses(%q) error = %v, wantErr %v", tt.uses, err, tt.wantErr)
			}
			if url != tt.wantURL {
				t.Errorf("ParseUses(%q) url = %q, want %q", tt.uses, url, tt.wantURL)
			}
			if ref != tt.wantRef {
				t.Errorf("ParseUses(%q) ref = %q, want %q", tt.uses, ref, tt.wantRef)
			}
		})
	}
}
