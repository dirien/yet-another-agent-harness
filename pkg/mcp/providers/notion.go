package providers

import (
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var _ mcp.Provider = (*Notion)(nil)

// Notion provides the Notion MCP server for reading/writing Notion pages.
type Notion struct {
	apiToken string
}

func NewNotion(apiToken string) *Notion {
	return &Notion{apiToken: apiToken}
}

func (n *Notion) Name() string { return "notion" }

func (n *Notion) Server() schema.MCPServer {
	return schema.MCPServer{
		Name:      "notion",
		Transport: schema.MCPTransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "@notionhq/notion-mcp-server"},
		Env: map[string]string{
			"OPENAPI_MCP_HEADERS": `{"Authorization": "Bearer ` + n.apiToken + `", "Notion-Version": "2022-06-28"}`,
		},
	}
}
