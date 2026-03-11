package mcp

import "github.com/dirien/yet-another-agent-harness/pkg/schema"

// Provider is the interface for MCP server definitions.
type Provider interface {
	// Name returns the unique server identifier used in settings.json mcpServers key.
	Name() string

	// Server returns the MCPServer schema for code generation.
	Server() schema.MCPServer
}

// Registry holds all registered providers.
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

// ToConfig converts all registered providers into an MCPConfig.
func (r *Registry) ToConfig() *schema.MCPConfig {
	cfg := &schema.MCPConfig{}
	for _, p := range r.providers {
		cfg.Servers = append(cfg.Servers, p.Server())
	}
	return cfg
}
