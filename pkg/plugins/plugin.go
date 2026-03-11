package plugins

import (
	"encoding/json"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// Plugin is the interface for generating Claude Code plugin packages.
type Plugin interface {
	// Plugin returns the plugin metadata for plugin.json generation.
	Plugin() schema.Plugin
}

// Registry holds all registered plugins.
type Registry struct {
	plugins []Plugin
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers a plugin.
func (r *Registry) Register(p Plugin) {
	r.plugins = append(r.plugins, p)
}

// Plugins returns all registered plugins.
func (r *Registry) Plugins() []Plugin {
	return r.plugins
}

// GeneratePluginJSON produces the plugin.json content for a plugin.
func GeneratePluginJSON(p Plugin) ([]byte, error) {
	return json.MarshalIndent(p.Plugin(), "", "  ")
}

// GenerateMarketplaceJSON produces a marketplace.json wrapping one or more plugins.
func GenerateMarketplaceJSON(owner schema.PluginAuthor, plugins ...Plugin) ([]byte, error) {
	type marketplaceEntry struct {
		schema.Plugin
		Source   string `json:"source"`
		Category string `json:"category,omitempty"`
	}

	type marketplace struct {
		Name    string              `json:"name"`
		Owner   schema.PluginAuthor `json:"owner"`
		Plugins []marketplaceEntry  `json:"plugins"`
	}

	m := marketplace{
		Name:  "yaah-marketplace",
		Owner: owner,
	}

	for _, p := range plugins {
		plug := p.Plugin()
		m.Plugins = append(m.Plugins, marketplaceEntry{
			Plugin:   plug,
			Source:   "./" + plug.Name,
			Category: "productivity",
		})
	}

	return json.MarshalIndent(m, "", "  ")
}
