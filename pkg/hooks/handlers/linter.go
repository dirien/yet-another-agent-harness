package handlers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// Step is a single command in a profile pipeline (fix -> format -> check).
type Step struct {
	Label string   // Human-readable label for output ("fix", "format", "check").
	Cmd   []string // Command + args. The file path is appended as the last arg.
	// If AppendFile is false, the file path is NOT appended (for whole-project linters).
	AppendFile bool
	// If FailBlocks is true, a non-zero exit from this step blocks the edit.
	FailBlocks bool
}

// Profile defines a complete lint pipeline for a set of file extensions.
type Profile struct {
	Name       string   // Profile name (e.g. "ruff", "golangci-lint", "biome").
	Extensions []string // File extensions this profile covers (e.g. ".go", ".py").
	Steps      []Step   // Ordered pipeline of commands to run.
}

// --- Built-in profiles ---

// Ruff returns a lint profile for Python using ruff.
func Ruff() Profile {
	return Profile{
		Name:       "ruff",
		Extensions: []string{".py"},
		Steps: []Step{
			{Label: "fix", Cmd: []string{"ruff", "check", "--fix"}, AppendFile: true},
			{Label: "format", Cmd: []string{"ruff", "format"}, AppendFile: true},
			{Label: "check", Cmd: []string{"ruff", "check"}, AppendFile: true, FailBlocks: true},
		},
	}
}

// GolangCILint returns a lint profile for Go using golangci-lint.
func GolangCILint() Profile {
	return Profile{
		Name:       "golangci-lint",
		Extensions: []string{".go"},
		Steps: []Step{
			{Label: "format", Cmd: []string{"gofmt", "-w"}, AppendFile: true},
			{Label: "lint", Cmd: []string{"golangci-lint", "run", "--fix"}, AppendFile: false, FailBlocks: true},
		},
	}
}

// GoVet returns a minimal Go lint profile (no golangci-lint needed).
func GoVet() Profile {
	return Profile{
		Name:       "go-vet",
		Extensions: []string{".go"},
		Steps: []Step{
			{Label: "format", Cmd: []string{"gofmt", "-w"}, AppendFile: true},
			{Label: "vet", Cmd: []string{"go", "vet", "./..."}, AppendFile: false, FailBlocks: true},
		},
	}
}

// Biome returns a lint profile for JS/TS using Biome.
func Biome() Profile {
	return Profile{
		Name:       "biome",
		Extensions: []string{".ts", ".tsx", ".js", ".jsx", ".json"},
		Steps: []Step{
			{Label: "check", Cmd: []string{"npx", "@biomejs/biome", "check", "--fix"}, AppendFile: true, FailBlocks: true},
		},
	}
}

// Prettier returns a lint profile for JS/TS using Prettier.
func Prettier() Profile {
	return Profile{
		Name:       "prettier",
		Extensions: []string{".ts", ".tsx", ".js", ".jsx", ".json", ".css", ".md"},
		Steps: []Step{
			{Label: "format", Cmd: []string{"npx", "prettier", "--write"}, AppendFile: true},
		},
	}
}

// RustFmt returns a lint profile for Rust.
func RustFmt() Profile {
	return Profile{
		Name:       "rustfmt",
		Extensions: []string{".rs"},
		Steps: []Step{
			{Label: "format", Cmd: []string{"rustfmt"}, AppendFile: true},
			{Label: "check", Cmd: []string{"cargo", "clippy", "--fix", "--allow-dirty"}, AppendFile: false, FailBlocks: true},
		},
	}
}

// --- Linter handler ---

var (
	_ hooks.Handler   = (*Linter)(nil)
	_ hooks.FileAware = (*Linter)(nil)
)

// Linter is a PostToolUse handler that enforces code standards via pluggable profiles.
// Implements: Handler, FileAware.
type Linter struct {
	profiles []Profile
	extIndex map[string]*Profile // extension -> first matching profile
}

// NewLinter creates a Linter with no profiles. Use AddProfile or NewLinterWith.
func NewLinter() *Linter {
	return &Linter{
		extIndex: make(map[string]*Profile),
	}
}

// NewLinterWith creates a Linter pre-loaded with the given profiles.
func NewLinterWith(profiles ...Profile) *Linter {
	t := NewLinter()
	for _, e := range profiles {
		t.AddProfile(e)
	}
	return t
}

// AddProfile registers a lint profile. Later profiles for the same extension override earlier ones.
func (t *Linter) AddProfile(e Profile) {
	t.profiles = append(t.profiles, e)
	for _, ext := range e.Extensions {
		stored := t.profiles[len(t.profiles)-1]
		t.extIndex[ext] = &stored
	}
}

// Profiles returns all registered lint profiles.
func (t *Linter) Profiles() []Profile { return t.profiles }

func (t *Linter) Name() string { return "linter" }

func (t *Linter) Events() []schema.HookEvent {
	return []schema.HookEvent{schema.HookPostToolUse}
}

func (t *Linter) Match() *regexp.Regexp {
	return editWriteMatch
}

func (t *Linter) FileExtensions() []string {
	exts := make([]string, 0, len(t.extIndex))
	for ext := range t.extIndex {
		exts = append(exts, ext)
	}
	return exts
}

// LintFile runs the lint pipeline for the given file path.
// If profileName is non-empty, only that profile is used; otherwise the profile
// is selected by file extension. Returns the combined lint output and whether
// the lint blocked (i.e. a FailBlocks step failed).
func (t *Linter) LintFile(ctx context.Context, filePath, profileName, cwd string) (output string, blocked bool, err error) {
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return "", false, fmt.Errorf("file not found: %s", filePath)
	}

	var profile *Profile
	if profileName != "" {
		for i := range t.profiles {
			if t.profiles[i].Name == profileName {
				profile = &t.profiles[i]
				break
			}
		}
		if profile == nil {
			return "", false, fmt.Errorf("unknown profile: %s", profileName)
		}
	} else {
		ext := filepath.Ext(filePath)
		profile = t.extIndex[ext]
		if profile == nil {
			return "", false, fmt.Errorf("no lint profile for extension %q", ext)
		}
	}

	var msgs []string
	for _, step := range profile.Steps {
		args := make([]string, len(step.Cmd))
		copy(args, step.Cmd)
		if step.AppendFile {
			args = append(args, filePath)
		}

		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		if cwd != "" {
			cmd.Dir = cwd
		}
		out, cmdErr := cmd.CombinedOutput()
		outStr := strings.TrimSpace(string(out))

		if cmdErr != nil {
			msgs = append(msgs, fmt.Sprintf("[%s/%s] %s", profile.Name, step.Label, outStr))
			if step.FailBlocks {
				return strings.Join(msgs, "\n"), true, nil
			}
		} else if outStr != "" {
			msgs = append(msgs, fmt.Sprintf("[%s/%s] %s", profile.Name, step.Label, outStr))
		}
	}

	return strings.Join(msgs, "\n"), false, nil
}

func (t *Linter) Execute(ctx context.Context, input *hooks.Input) (*hooks.Result, error) {
	fp := input.FilePath()
	if fp == "" {
		return nil, nil
	}
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return nil, nil
	}

	ext := filepath.Ext(fp)
	profile, ok := t.extIndex[ext]
	if !ok {
		return nil, nil
	}

	var msgs []string
	for _, step := range profile.Steps {
		args := make([]string, len(step.Cmd))
		copy(args, step.Cmd)
		if step.AppendFile {
			args = append(args, fp)
		}

		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Dir = input.Cwd
		out, err := cmd.CombinedOutput()
		outStr := strings.TrimSpace(string(out))

		if err != nil {
			msgs = append(msgs, fmt.Sprintf("[%s/%s] %s", profile.Name, step.Label, outStr))
			if step.FailBlocks {
				return &hooks.Result{
					Error: strings.Join(msgs, "\n"),
					Block: true,
				}, nil
			}
		} else if outStr != "" {
			msgs = append(msgs, fmt.Sprintf("[%s/%s] %s", profile.Name, step.Label, outStr))
		}
	}

	if len(msgs) > 0 {
		return &hooks.Result{Output: strings.Join(msgs, "\n")}, nil
	}
	return nil, nil
}
