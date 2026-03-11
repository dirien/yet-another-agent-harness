package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ mcp.Provider = (*Custom)(nil)

// Custom wraps any arbitrary MCPServer definition as a Provider.
// Use this when you need a one-off MCP server that doesn't have a dedicated type.
type Custom struct {
	server schema.MCPServer
}

func NewCustom(server schema.MCPServer) *Custom {
	return &Custom{server: server}
}

func (c *Custom) Name() string { return c.server.Name }

func (c *Custom) Server() schema.MCPServer { return c.server }
