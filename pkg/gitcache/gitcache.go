// Package gitcache provides shared git clone + cache logic for fetching
// files from pinned remote repositories. Both RemoteSkill and RemoteAgent
// use this package to avoid duplicating the clone/cache/marker flow.
package gitcache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

// HomeDir returns the yaah home directory (~/.yaah).
// All cache and config files live under this path.
// Override with YAAH_HOME.
func HomeDir() string {
	if dir := os.Getenv("YAAH_HOME"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yaah")
}

// FetchFile clones a remote git repo (or uses a cached clone) and returns the
// content of a single file at subpath. The kind parameter ("skills" or "agents")
// controls the cache subdirectory.
func FetchFile(uses, subpath, kind string) (string, error) {
	content, _, err := FetchFileWithExtras(uses, subpath, kind, false)
	return content, err
}

// FetchFileWithExtras clones a remote git repo (or uses a cached clone) and
// returns the content of a single file at subpath plus any extra sibling files
// when scanExtras is true. The kind parameter controls the cache subdirectory.
func FetchFileWithExtras(uses, subpath, kind string, scanExtras bool) (string, map[string]string, error) {
	repoURL, ref, err := ParseUses(uses)
	if err != nil {
		return "", nil, err
	}

	cacheKey := fmt.Sprintf("%x", sha256.Sum256([]byte(kind+":"+uses+":"+subpath)))
	cacheDir := filepath.Join(HomeDir(), "cache", kind, cacheKey)

	markerFile := filepath.Join(cacheDir, ".yaah-ref")
	repoDir := filepath.Join(cacheDir, "repo")

	// Check if already cached with the right ref.
	if data, err := os.ReadFile(markerFile); err == nil && strings.TrimSpace(string(data)) == ref {
		content, err := os.ReadFile(filepath.Join(repoDir, subpath))
		if err == nil {
			var extras map[string]string
			if scanExtras {
				extras = ScanExtraFiles(repoDir, subpath)
			}
			return string(content), extras, nil
		}
	}

	// Clone fresh.
	_ = os.RemoveAll(cacheDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("create cache dir: %w", err)
	}

	if err := gitClone(repoURL, ref, repoDir); err != nil {
		return "", nil, err
	}

	// Write ref marker.
	_ = os.WriteFile(markerFile, []byte(ref), 0o644)

	content, err := os.ReadFile(filepath.Join(repoDir, subpath))
	if err != nil {
		return "", nil, fmt.Errorf("read %s in repo: %w", subpath, err)
	}

	var extras map[string]string
	if scanExtras {
		extras = ScanExtraFiles(repoDir, subpath)
	}
	return string(content), extras, nil
}

// gitClone performs a shallow clone using a 3-tier strategy:
//  1. git clone --depth=1 --branch <ref> --single-branch  (fast; works for tags and branches)
//  2. git init + git fetch --depth=1 origin <sha> + git checkout FETCH_HEAD
//     (efficient; works for SHAs on servers that support fetching by SHA)
//  3. git clone --no-checkout + git checkout <ref>
//     (full download; works everywhere but retrieves entire history)
func gitClone(repoURL, ref, repoDir string) error {
	// Tier 1: shallow clone by branch/tag name.
	cmd := exec.Command("git", "clone", "--depth=1", "--branch", ref, "--single-branch", repoURL, repoDir)
	if _, err := cmd.CombinedOutput(); err == nil {
		return nil
	}

	// Tier 2: shallow fetch by SHA (server must support uploadpack.allowReachableSHA1InWant
	// or uploadpack.allowAnySHA1InWant; works on GitHub and most modern hosts).
	_ = os.RemoveAll(repoDir)
	if mkErr := os.MkdirAll(repoDir, 0o755); mkErr == nil {
		initOut, initErr := exec.Command("git", "init", repoDir).CombinedOutput()
		if initErr == nil {
			addOut, addErr := exec.Command("git", "-C", repoDir, "remote", "add", "origin", repoURL).CombinedOutput()
			if addErr == nil {
				fetchOut, fetchErr := exec.Command("git", "-C", repoDir, "fetch", "--depth=1", "origin", ref).CombinedOutput()
				if fetchErr == nil {
					checkoutOut, checkoutErr := exec.Command("git", "-C", repoDir, "checkout", "FETCH_HEAD").CombinedOutput()
					if checkoutErr == nil {
						return nil
					}
					_ = checkoutOut // captured but not needed on success path
				}
				_ = fetchOut
			}
			_ = addOut
		}
		_ = initOut
	}

	// Tier 3: full clone without checkout, then check out the exact ref.
	// This downloads the entire history but works on any server.
	_ = os.RemoveAll(repoDir)
	if out, err := exec.Command("git", "clone", "--no-checkout", repoURL, repoDir).CombinedOutput(); err != nil {
		return fmt.Errorf("git clone --no-checkout %s failed: %s", repoURL, string(out))
	}
	if out, err := exec.Command("git", "-C", repoDir, "checkout", ref).CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout %s failed: %s", ref, string(out))
	}
	return nil
}

// ScanExtraFiles walks the directory containing the target file and collects
// all sibling files and subdirectory files (e.g. references/), returning them
// as a map of relative path to content. The target file itself is excluded.
func ScanExtraFiles(repoDir, subpath string) map[string]string {
	targetFile := filepath.Join(repoDir, subpath)
	targetDir := filepath.Dir(targetFile)
	targetBase := filepath.Base(targetFile)

	extras := make(map[string]string)
	_ = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(targetDir, path)
		if err != nil || rel == targetBase {
			return nil
		}
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
