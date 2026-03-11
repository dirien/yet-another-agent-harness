package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ mcp.Provider = (*Context7)(nil)

// Context7 provides the Context7 MCP server for up-to-date library documentation.
type Context7 struct{}

func NewContext7() *Context7 { return &Context7{} }

func (c *Context7) Name() string { return "context7" }

func (c *Context7) Server() schema.MCPServer {
	return schema.MCPServer{
		Name:      "context7",
		Transport: schema.MCPTransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "@context7/mcp"},
	}
}
