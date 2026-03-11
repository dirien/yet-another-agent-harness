package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.Provider = (*Custom)(nil)

// Custom wraps any arbitrary LSPServer definition as a Provider.
// Use this when you need a one-off LSP server that doesn't have a dedicated type.
//
// Example:
//
//	h.LSP().Register(providers.NewCustom(
//		schema.LSPServer{
//			ID:                  "rust",
//			Command:             []string{"rust-analyzer"},
//			ExtensionToLanguage: map[string]string{".rs": "rust"},
//		},
//		"rustup component add rust-analyzer",
//	))
type Custom struct {
	server      schema.LSPServer
	installHint string
}

func NewCustom(server schema.LSPServer, installHint string) *Custom {
	return &Custom{server: server, installHint: installHint}
}

func (c *Custom) Name() string             { return c.server.ID }
func (c *Custom) Server() schema.LSPServer { return c.server }
func (c *Custom) InstallHint() string      { return c.installHint }
