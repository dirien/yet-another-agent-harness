package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ mcp.Provider = (*Pulumi)(nil)

// Pulumi provides the Pulumi MCP server for AI-powered infrastructure management.
type Pulumi struct{}

func NewPulumi() *Pulumi { return &Pulumi{} }

func (p *Pulumi) Name() string { return "pulumi" }

func (p *Pulumi) Server() schema.MCPServer {
	return schema.MCPServer{
		Name:      "pulumi",
		Transport: schema.MCPTransportHTTP,
		URL:       "https://mcp.ai.pulumi.com/mcp",
	}
}
