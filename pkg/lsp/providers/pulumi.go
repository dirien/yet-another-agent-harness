package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.Provider = (*PulumiLSP)(nil)

// PulumiLSP provides the Pulumi YAML language server.
type PulumiLSP struct{}

func NewPulumiLSP() *PulumiLSP { return &PulumiLSP{} }

func (p *PulumiLSP) Name() string { return "pulumi" }

func (p *PulumiLSP) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "pulumi",
		Command:             []string{"pulumi", "lsp"},
		ExtensionToLanguage: map[string]string{".pp": "pulumi"},
	}
}

func (p *PulumiLSP) InstallHint() string {
	return "curl -fsSL https://get.pulumi.com | sh"
}
