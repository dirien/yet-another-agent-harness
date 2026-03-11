package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.MarketplaceProvider = (*CSharp)(nil)

// CSharp provides the C# language server.
type CSharp struct{}

func NewCSharp() *CSharp { return &CSharp{} }

func (c *CSharp) Name() string { return "csharp" }

func (c *CSharp) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "csharp",
		Command:             []string{"csharp-ls"},
		ExtensionToLanguage: map[string]string{".cs": "csharp"},
	}
}

func (c *CSharp) InstallHint() string {
	return "dotnet tool install -g csharp-ls"
}

func (c *CSharp) MarketplaceKey() string { return "csharp-lsp@claude-plugins-official" }
