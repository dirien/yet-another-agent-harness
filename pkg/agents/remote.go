package agents

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dirien/yet-another-agent-harness/pkg/gitcache"
)

// RemoteAgent loads an agent definition from a pinned remote git repo.
// Uses syntax: "github.com/owner/repo@ref" (tag, branch, or commit SHA).
// The remote markdown file becomes the agent's system prompt (Content).
// If the remote file contains YAML frontmatter (---...---), it is stripped
// so that buildAgentMarkdown can re-generate clean frontmatter from the
// agent's metadata fields.
type RemoteAgent struct {
	name        string
	description string
	model       string
	tools       string
	uses        string // e.g. "github.com/msitarzewski/agency-agents@abc123"
	subpath     string // path within repo, e.g. "engineering/engineering-ai-engineer.md"

	// Optional advanced fields.
	advanced *AgentAdvanced

	// Resolved content (lazy-loaded).
	content string
	once    sync.Once
}

// RemoteAgentOption configures optional fields on a RemoteAgent.
type RemoteAgentOption func(*RemoteAgent)

// WithModel sets the model override for the remote agent.
func WithModel(model string) RemoteAgentOption {
	return func(r *RemoteAgent) { r.model = model }
}

// WithTools sets the tool allowlist for the remote agent.
func WithTools(tools string) RemoteAgentOption {
	return func(r *RemoteAgent) { r.tools = tools }
}

// WithAdvanced sets advanced frontmatter fields for the remote agent.
func WithAdvanced(adv AgentAdvanced) RemoteAgentOption {
	return func(r *RemoteAgent) { r.advanced = &adv }
}

// NewRemoteAgent creates an agent that fetches its content from a remote git repo.
//
//	uses:    "github.com/owner/repo@v1.0.0" or "github.com/owner/repo@sha"
//	subpath: path within repo to the agent markdown file
func NewRemoteAgent(name, description, uses, subpath string, opts ...RemoteAgentOption) *RemoteAgent {
	r := &RemoteAgent{
		name:        name,
		description: description,
		uses:        uses,
		subpath:     subpath,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

var (
	_ Agent             = (*RemoteAgent)(nil)
	_ AgentWithAdvanced = (*RemoteAgent)(nil)
	_ AgentWithSource   = (*RemoteAgent)(nil)
)

func (r *RemoteAgent) Name() string        { return r.name }
func (r *RemoteAgent) Description() string { return r.description }
func (r *RemoteAgent) Model() string       { return r.model }
func (r *RemoteAgent) Tools() string       { return r.tools }

func (r *RemoteAgent) Content() string {
	r.resolve()
	return r.content
}

func (r *RemoteAgent) Advanced() AgentAdvanced {
	if r.advanced != nil {
		return *r.advanced
	}
	return AgentAdvanced{}
}

// Uses returns the remote git reference (e.g. "github.com/owner/repo@ref").
func (r *RemoteAgent) Uses() string { return r.uses }

// Subpath returns the path within the remote repo to the agent file.
func (r *RemoteAgent) Subpath() string { return r.subpath }

func (r *RemoteAgent) resolve() {
	r.once.Do(func() {
		content, err := gitcache.FetchFile(r.uses, r.subpath, "agents")
		if err != nil {
			r.content = fmt.Sprintf("# Error loading remote agent\n\nFailed to fetch %s: %s\n", r.uses, err)
			return
		}
		r.content = stripFrontmatter(content)
	})
}

// stripFrontmatter removes YAML frontmatter (---\n...\n---) from the
// beginning of a markdown file if present, returning just the body.
// This prevents double frontmatter when buildAgentMarkdown re-generates it.
func stripFrontmatter(content string) string {
	// Normalize line endings for consistent parsing.
	normalized := strings.ReplaceAll(content, "\r\n", "\n")

	if !strings.HasPrefix(normalized, "---\n") {
		return normalized
	}

	// Find the closing ---.
	rest := normalized[4:] // skip opening "---\n"
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		// Check if the closing --- is at the very end of the file.
		if strings.HasSuffix(rest, "\n---") {
			body := rest[strings.Index(rest, "\n---")+4:]
			return strings.TrimLeft(body, "\n")
		}
		return normalized // no closing ---, return as-is
	}
	return strings.TrimLeft(rest[idx+5:], "\n")
}
