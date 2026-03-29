package skills

import (
	"fmt"

	"github.com/dirien/yet-another-agent-harness/pkg/gitcache"
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
	metadata    *SkillMetadata

	// Resolved content (lazy-loaded).
	content    string
	extraFiles map[string]string // relative path → content
	resolved   bool
}

// RemoteSkillOption configures optional fields on a RemoteSkill.
type RemoteSkillOption func(*RemoteSkill)

// WithMetadata attaches catalog metadata to a remote skill.
func WithMetadata(m SkillMetadata) RemoteSkillOption {
	return func(r *RemoteSkill) {
		r.metadata = &m
	}
}

// NewRemoteSkill creates a skill that fetches from a remote git repo.
//
//	uses:    "github.com/owner/repo@v1.0.0" or "github.com/owner/repo@sha"
//	subpath: path within repo to SKILL.md (empty = "SKILL.md")
func NewRemoteSkill(name, description, uses, subpath string, opts ...RemoteSkillOption) *RemoteSkill {
	if subpath == "" {
		subpath = "SKILL.md"
	}
	r := &RemoteSkill{
		name:        name,
		description: description,
		uses:        uses,
		subpath:     subpath,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
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

// Metadata returns catalog metadata if set via WithMetadata.
// Implements SkillWithMetadata.
func (r *RemoteSkill) Metadata() SkillMetadata {
	if r.metadata != nil {
		return *r.metadata
	}
	return SkillMetadata{}
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

	content, extras, err := gitcache.FetchFileWithExtras(r.uses, r.subpath, "skills", true)
	if err != nil {
		r.content = fmt.Sprintf("# Error loading remote skill\n\nFailed to fetch %s: %s\n", r.uses, err)
		return
	}
	r.content = content
	r.extraFiles = extras
}

// ParseUses splits "github.com/owner/repo@ref" into (repoURL, ref).
// Deprecated: Use gitcache.ParseUses instead.
func ParseUses(uses string) (repoURL, ref string, err error) {
	return gitcache.ParseUses(uses)
}

// HomeDir returns the yaah home directory (~/.yaah).
// Deprecated: Use gitcache.HomeDir instead.
func HomeDir() string {
	return gitcache.HomeDir()
}
