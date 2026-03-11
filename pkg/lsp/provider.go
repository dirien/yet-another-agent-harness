package lsp

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// Provider is the interface for LSP server definitions.
// To add a new LSP server, implement this interface and register it in the Harness.
type Provider interface {
	// Name returns the unique server identifier (e.g. "gopls").
	Name() string

	// Server returns the LSPServer schema for code generation.
	Server() schema.LSPServer

	// InstallHint returns a human-readable instruction for installing
	// the server binary (e.g. "go install golang.org/x/tools/gopls@latest").
	InstallHint() string
}

// MarketplaceProvider is an optional interface for LSP providers that are
// available in the Claude Code official marketplace. Providers implementing
// this return their enabledPlugins key (e.g. "gopls-lsp@claude-plugins-official").
type MarketplaceProvider interface {
	Provider
	MarketplaceKey() string
}

// Registry holds all registered LSP providers.
type Registry struct {
	providers []Provider
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers a provider.
func (r *Registry) Register(p Provider) {
	r.providers = append(r.providers, p)
}

// Providers returns all registered providers.
func (r *Registry) Providers() []Provider {
	return r.providers
}

// ToConfig converts all registered providers into an LSPConfig.
func (r *Registry) ToConfig() *schema.LSPConfig {
	cfg := &schema.LSPConfig{}
	for _, p := range r.providers {
		cfg.Servers = append(cfg.Servers, p.Server())
	}
	return cfg
}

// CheckResult describes whether an LSP server binary is available.
type CheckResult struct {
	Name        string
	Installed   bool
	BinaryPath  string // resolved path if installed
	InstallHint string
}

// CheckAll validates that every registered LSP server binary exists in $PATH.
func (r *Registry) CheckAll() []CheckResult {
	var results []CheckResult
	for _, p := range r.providers {
		srv := p.Server()
		if len(srv.Command) == 0 {
			continue
		}
		cr := CheckResult{
			Name:        p.Name(),
			InstallHint: p.InstallHint(),
		}
		path, err := exec.LookPath(srv.Command[0])
		if err == nil {
			cr.Installed = true
			cr.BinaryPath = path
		}
		results = append(results, cr)
	}
	return results
}

// FormatCheckResults returns a human-readable summary of LSP binary checks.
func FormatCheckResults(results []CheckResult) string {
	var b strings.Builder
	for _, cr := range results {
		if cr.Installed {
			_, _ = fmt.Fprintf(&b, "  ✓ %-24s %s\n", cr.Name, cr.BinaryPath)
		} else {
			_, _ = fmt.Fprintf(&b, "  ✗ %-24s not found → %s\n", cr.Name, cr.InstallHint)
		}
	}
	return b.String()
}
