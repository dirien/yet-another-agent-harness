package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ mcp.Provider = (*Yaah)(nil)

// Yaah provides yaah itself as an MCP server, exposing built-in tools
// (secret scanning, linting, command checking, doctor, session info)
// over the Model Context Protocol via stdio transport.
type Yaah struct{}

// NewYaah creates a new Yaah MCP provider.
func NewYaah() *Yaah { return &Yaah{} }

// Name returns the provider identifier.
func (y *Yaah) Name() string { return "yaah" }

// Server returns the MCPServer definition for yaah serve.
func (y *Yaah) Server() schema.MCPServer {
	return schema.MCPServer{
		Name:      "yaah",
		Transport: schema.MCPTransportStdio,
		Command:   "yaah",
		Args:      []string{"serve"},
	}
}
