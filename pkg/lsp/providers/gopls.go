package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.MarketplaceProvider = (*Gopls)(nil)

// Gopls provides the Go language server.
type Gopls struct{}

func NewGopls() *Gopls { return &Gopls{} }

func (g *Gopls) Name() string { return "gopls" }

func (g *Gopls) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "gopls",
		Command:             []string{"gopls", "serve"},
		ExtensionToLanguage: map[string]string{".go": "go"},
	}
}

func (g *Gopls) InstallHint() string {
	return "go install golang.org/x/tools/gopls@latest"
}

func (g *Gopls) MarketplaceKey() string { return "gopls-lsp@claude-plugins-official" }
