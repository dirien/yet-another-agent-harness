package gitcache

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseUses verifies that ParseUses correctly splits a "host/owner/repo@ref"
// string into a full HTTPS clone URL and a ref, and returns errors for malformed
// inputs.
func TestParseUses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		uses        string
		wantURL     string
		wantRef     string
		wantErrFrag string // non-empty means we expect an error containing this substring
	}{
		{
			name:    "semver tag",
			uses:    "github.com/owner/repo@v1.0.0",
			wantURL: "https://github.com/owner/repo.git",
			wantRef: "v1.0.0",
		},
		{
			name:    "commit sha",
			uses:    "github.com/owner/repo@abc123",
			wantURL: "https://github.com/owner/repo.git",
			wantRef: "abc123",
		},
		{
			name:        "missing at-ref",
			uses:        "github.com/owner/repo",
			wantErrFrag: "missing @ref",
		},
		{
			name:        "empty ref after at",
			uses:        "github.com/owner/repo@",
			wantErrFrag: "empty ref after @",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotURL, gotRef, err := ParseUses(tc.uses)

			if tc.wantErrFrag != "" {
				if err == nil {
					t.Fatalf("ParseUses(%q): expected error containing %q, got nil", tc.uses, tc.wantErrFrag)
				}
				if !containsStr(err.Error(), tc.wantErrFrag) {
					t.Fatalf("ParseUses(%q): error %q does not contain %q", tc.uses, err.Error(), tc.wantErrFrag)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseUses(%q): unexpected error: %v", tc.uses, err)
			}
			if gotURL != tc.wantURL {
				t.Errorf("ParseUses(%q): repoURL = %q, want %q", tc.uses, gotURL, tc.wantURL)
			}
			if gotRef != tc.wantRef {
				t.Errorf("ParseUses(%q): ref = %q, want %q", tc.uses, gotRef, tc.wantRef)
			}
		})
	}
}

// TestScanExtraFiles verifies that ScanExtraFiles correctly collects sibling
// and subdirectory files while excluding the target file itself and any
// dot-prefixed (hidden) files.
func TestScanExtraFiles(t *testing.T) {
	t.Parallel()

	t.Run("returns sibling and subdir files, excludes target and hidden", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		// Target file — must be excluded from extras.
		writeFile(t, dir, "agent.md", "# agent")

		// Sibling plain file — must be included.
		writeFile(t, dir, "helper.txt", "helper content")

		// Hidden file — must be excluded.
		writeFile(t, dir, ".hidden", "secret")

		// File inside a subdirectory — must be included.
		writeFile(t, filepath.Join(dir, "refs"), "doc.md", "doc content")

		extras := ScanExtraFiles(dir, "agent.md")

		if extras == nil {
			t.Fatal("ScanExtraFiles returned nil, want a non-nil map")
		}

		wantEntries := map[string]string{
			"helper.txt":  "helper content",
			"refs/doc.md": "doc content",
		}

		if len(extras) != len(wantEntries) {
			t.Errorf("extras has %d entries, want %d: %v", len(extras), len(wantEntries), extras)
		}

		for relPath, wantContent := range wantEntries {
			gotContent, ok := extras[relPath]
			if !ok {
				t.Errorf("extras missing key %q", relPath)
				continue
			}
			if gotContent != wantContent {
				t.Errorf("extras[%q] = %q, want %q", relPath, gotContent, wantContent)
			}
		}

		// Explicitly assert excluded items are absent.
		for _, excluded := range []string{"agent.md", ".hidden"} {
			if _, ok := extras[excluded]; ok {
				t.Errorf("extras should not contain %q", excluded)
			}
		}
	})

	t.Run("only target file returns nil", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFile(t, dir, "agent.md", "# agent")

		extras := ScanExtraFiles(dir, "agent.md")

		if extras != nil {
			t.Errorf("ScanExtraFiles with only target file: want nil, got %v", extras)
		}
	})
}

// TestHomeDir verifies that HomeDir returns the value of YAAH_HOME when set,
// and falls back to ~/.yaah when the variable is unset.
func TestHomeDir(t *testing.T) {
	// Not parallel: modifies the process environment.

	t.Run("YAAH_HOME override", func(t *testing.T) {
		custom := "/tmp/my-custom-yaah"
		t.Setenv("YAAH_HOME", custom)

		got := HomeDir()
		if got != custom {
			t.Errorf("HomeDir() = %q, want %q", got, custom)
		}
	})

	t.Run("fallback to ~/.yaah when YAAH_HOME is unset", func(t *testing.T) {
		// Ensure the variable is absent for this sub-test.
		t.Setenv("YAAH_HOME", "")

		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot determine user home dir:", err)
		}
		want := filepath.Join(home, ".yaah")

		got := HomeDir()
		if got != want {
			t.Errorf("HomeDir() = %q, want %q", got, want)
		}
	})
}

// --- helpers -----------------------------------------------------------------

// writeFile creates intermediate directories and writes content to
// filepath.Join(dir, name).
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	full := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", full, err)
	}
}

// containsStr reports whether s contains substr.
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
