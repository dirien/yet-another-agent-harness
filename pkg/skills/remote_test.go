package skills_test

import (
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/skills"
)

func TestParseUses(t *testing.T) {
	tests := []struct {
		name     string
		uses     string
		wantRepo string
		wantRef  string
		wantErr  bool
	}{
		{
			name:     "tag ref",
			uses:     "github.com/owner/repo@v1.0.0",
			wantRepo: "https://github.com/owner/repo.git",
			wantRef:  "v1.0.0",
		},
		{
			name:     "sha ref",
			uses:     "github.com/owner/repo@abc1234",
			wantRepo: "https://github.com/owner/repo.git",
			wantRef:  "abc1234",
		},
		{
			name:     "branch ref",
			uses:     "github.com/dirien/my-skills@main",
			wantRepo: "https://github.com/dirien/my-skills.git",
			wantRef:  "main",
		},
		{
			name:    "missing @",
			uses:    "github.com/owner/repo",
			wantErr: true,
		},
		{
			name:    "empty ref",
			uses:    "github.com/owner/repo@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, ref, err := skills.ParseUses(tt.uses)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo: got %q, want %q", repo, tt.wantRepo)
			}
			if ref != tt.wantRef {
				t.Errorf("ref: got %q, want %q", ref, tt.wantRef)
			}
		})
	}
}

func TestNewRemoteSkill_DefaultSubpath(t *testing.T) {
	s := skills.NewRemoteSkill("test", "desc", "github.com/x/y@v1", "")
	source := s.Source()
	if source.Subpath != "SKILL.md" {
		t.Errorf("expected default subpath 'SKILL.md', got %q", source.Subpath)
	}
}

func TestNewRemoteSkill_CustomSubpath(t *testing.T) {
	s := skills.NewRemoteSkill("test", "desc", "github.com/x/y@v1", "skills/commit/SKILL.md")
	source := s.Source()
	if source.Subpath != "skills/commit/SKILL.md" {
		t.Errorf("expected subpath 'skills/commit/SKILL.md', got %q", source.Subpath)
	}
}

func TestNewRemoteSkill_NameAndDescription(t *testing.T) {
	s := skills.NewRemoteSkill("my-skill", "My skill description", "github.com/x/y@v1", "")
	if s.Name() != "my-skill" {
		t.Errorf("Name(): got %q, want %q", s.Name(), "my-skill")
	}
	if s.Description() != "My skill description" {
		t.Errorf("Description(): got %q, want %q", s.Description(), "My skill description")
	}
}

func TestHomeDir(t *testing.T) {
	t.Setenv("YAAH_HOME", "/custom/path")
	if got := skills.HomeDir(); got != "/custom/path" {
		t.Errorf("HomeDir(): got %q, want %q", got, "/custom/path")
	}
}
