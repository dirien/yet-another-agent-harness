package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.Provider = (*PulumiYAML)(nil)

// PulumiYAML provides the Pulumi YAML language server from pulumi/pulumi-lsp.
type PulumiYAML struct{}

// NewPulumiYAML creates a new Pulumi YAML LSP provider.
func NewPulumiYAML() *PulumiYAML { return &PulumiYAML{} }

func (p *PulumiYAML) Name() string { return "pulumi-yaml" }

func (p *PulumiYAML) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "pulumi-yaml",
		Command:             []string{"pulumi-lsp"},
		ExtensionToLanguage: map[string]string{".yaml": "yaml", ".yml": "yaml"},
	}
}

func (p *PulumiYAML) InstallHint() string {
	return "go install github.com/pulumi/pulumi-lsp@latest"
}
