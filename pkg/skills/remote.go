package skills

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	schemapkg "github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// RemoteSkill loads a SKILL.md from a pinned remote git repo.
// Uses syntax: "github.com/owner/repo@ref" (tag, branch, or commit SHA).
// It also discovers extra files (e.g. references/) alongside SKILL.md.
type RemoteSkill struct {
	name        string
	description string
	uses        string // e.g. "github.com/dirien/my-skills@v1.0.0"
	subpath     string // path within repo, e.g. "skills/commit/SKILL.md"

	// Resolved content (lazy-loaded).
	content    string
	extraFiles map[string]string // relative path → content
	resolved   bool
}

// NewRemoteSkill creates a skill that fetches from a remote git repo.
//
//	uses:    "github.com/owner/repo@v1.0.0" or "github.com/owner/repo@sha"
//	subpath: path within repo to SKILL.md (empty = "SKILL.md")
func NewRemoteSkill(name, description, uses, subpath string) *RemoteSkill {
	if subpath == "" {
		subpath = "SKILL.md"
	}
	return &RemoteSkill{
		name:        name,
		description: description,
		uses:        uses,
		subpath:     subpath,
	}
}

var _ SkillWithFiles = (*RemoteSkill)(nil)

func (r *RemoteSkill) Name() string        { return r.name }
func (r *RemoteSkill) Description() string { return r.description }

func (r *RemoteSkill) Source() schemapkg.SkillSource {
	return schemapkg.SkillSource{
		Uses:    r.uses,
		Subpath: r.subpath,
	}
}

func (r *RemoteSkill) Content() string {
	r.resolve()
	return r.content
}

// ExtraFiles returns additional files discovered alongside SKILL.md
// (e.g. a references/ folder). Implements SkillWithFiles.
func (r *RemoteSkill) ExtraFiles() map[string]string {
	r.resolve()
	return r.extraFiles
}

func (r *RemoteSkill) resolve() {
	if r.resolved {
		return
	}
	r.resolved = true

	content, extras, err := r.fetch()
	if err != nil {
		r.content = fmt.Sprintf("# Error loading remote skill\n\nFailed to fetch %s: %s\n", r.uses, err)
		return
	}
	r.content = content
	r.extraFiles = extras
}

// ParseUses splits "github.com/owner/repo@ref" into (repoURL, ref).
func ParseUses(uses string) (repoURL, ref string, err error) {
	at := strings.LastIndex(uses, "@")
	if at < 0 {
		return "", "", fmt.Errorf("invalid uses %q: missing @ref (expected host/owner/repo@ref)", uses)
	}
	repo := uses[:at]
	ref = uses[at+1:]
	if ref == "" {
		return "", "", fmt.Errorf("invalid uses %q: empty ref after @", uses)
	}
	repoURL = "https://" + repo + ".git"
	return repoURL, ref, nil
}

func (r *RemoteSkill) fetch() (string, map[string]string, error) {
	repoURL, ref, err := ParseUses(r.uses)
	if err != nil {
		return "", nil, err
	}

	// Cache in ~/.yaah/cache/skills/<hash>
	cacheKey := fmt.Sprintf("%x", sha256.Sum256([]byte(r.uses+r.subpath)))
	cacheDir := filepath.Join(HomeDir(), "cache", "skills", cacheKey)

	// Check if already cached with the right ref.
	markerFile := filepath.Join(cacheDir, ".yaah-ref")
	repoDir := filepath.Join(cacheDir, "repo")
	if data, err := os.ReadFile(markerFile); err == nil && strings.TrimSpace(string(data)) == ref {
		content, err := os.ReadFile(filepath.Join(repoDir, r.subpath))
		if err == nil {
			extras := scanExtraFiles(repoDir, r.subpath)
			return string(content), extras, nil
		}
	}

	// Clone fresh.
	_ = os.RemoveAll(cacheDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("create cache dir: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth=1", "--branch", ref, "--single-branch", repoURL, repoDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		// If --branch fails (it doesn't work for SHAs), try clone + checkout.
		_ = os.RemoveAll(repoDir)
		cmd = exec.Command("git", "clone", "--no-checkout", repoURL, repoDir)
		if out2, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("git clone failed: %s\n%s", string(out), string(out2))
		}
		cmd = exec.Command("git", "-C", repoDir, "checkout", ref)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("git checkout %s failed: %s", ref, string(out))
		}
	}

	// Write ref marker.
	_ = os.WriteFile(markerFile, []byte(ref), 0o644)

	content, err := os.ReadFile(filepath.Join(repoDir, r.subpath))
	if err != nil {
		return "", nil, fmt.Errorf("read %s in repo: %w", r.subpath, err)
	}

	extras := scanExtraFiles(repoDir, r.subpath)
	return string(content), extras, nil
}

// scanExtraFiles walks the directory containing the SKILL.md and collects all
// sibling files and subdirectory files (e.g. references/), returning them as
// a map of relative path to content. The SKILL.md itself is excluded.
func scanExtraFiles(repoDir, subpath string) map[string]string {
	skillFile := filepath.Join(repoDir, subpath)
	skillDir := filepath.Dir(skillFile)
	skillBase := filepath.Base(skillFile)

	extras := make(map[string]string)
	_ = filepath.Walk(skillDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(skillDir, path)
		if err != nil || rel == skillBase {
			return nil
		}
		// Skip hidden files and non-text files.
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		extras[rel] = string(data)
		return nil
	})
	if len(extras) == 0 {
		return nil
	}
	return extras
}

// HomeDir returns the yaah home directory (~/.yaah).
// All cache and config files live under this path.
func HomeDir() string {
	if dir := os.Getenv("YAAH_HOME"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yaah")
}
