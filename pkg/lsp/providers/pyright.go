package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.MarketplaceProvider = (*Pyright)(nil)

// Pyright provides the Python language server.
type Pyright struct{}

func NewPyright() *Pyright { return &Pyright{} }

func (p *Pyright) Name() string { return "pyright" }

func (p *Pyright) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "pyright",
		Command:             []string{"pyright-langserver", "--stdio"},
		ExtensionToLanguage: map[string]string{".py": "python", ".pyi": "python"},
	}
}

func (p *Pyright) InstallHint() string {
	return "pip install pyright"
}

func (p *Pyright) MarketplaceKey() string { return "pyright-lsp@claude-plugins-official" }
