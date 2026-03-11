package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.MarketplaceProvider = (*TypeScript)(nil)

// TypeScript provides the TypeScript/JavaScript language server.
type TypeScript struct{}

func NewTypeScript() *TypeScript { return &TypeScript{} }

func (t *TypeScript) Name() string { return "typescript" }

func (t *TypeScript) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:      "typescript",
		Command: []string{"typescript-language-server", "--stdio"},
		ExtensionToLanguage: map[string]string{
			".ts": "typescript", ".tsx": "typescriptreact",
			".js": "javascript", ".jsx": "javascriptreact",
			".mjs": "javascript", ".cjs": "javascript",
		},
	}
}

func (t *TypeScript) InstallHint() string {
	return "npm install -g typescript-language-server typescript"
}

func (t *TypeScript) MarketplaceKey() string { return "typescript-lsp@claude-plugins-official" }
