package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ lsp.Provider = (*YamlLS)(nil)

// YamlLS provides the YAML language server (useful for Pulumi YAML).
type YamlLS struct{}

func NewYamlLS() *YamlLS { return &YamlLS{} }

func (y *YamlLS) Name() string { return "yaml-ls" }

func (y *YamlLS) Server() schema.LSPServer {
	return schema.LSPServer{
		ID:                  "yaml-ls",
		Command:             []string{"yaml-language-server", "--stdio"},
		ExtensionToLanguage: map[string]string{".yaml": "yaml", ".yml": "yaml"},
	}
}

func (y *YamlLS) InstallHint() string {
	return "npm install -g yaml-language-server"
}
