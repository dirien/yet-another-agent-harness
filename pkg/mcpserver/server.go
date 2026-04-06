// Package mcpserver implements a custom MCP server that exposes yaah's
// handler capabilities as MCP tools, allowing Claude Code to invoke them
// directly via the Model Context Protocol.
package mcpserver

import (
	"context"

	"github.com/dirien/yet-another-agent-harness/pkg/harness"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version is the MCP server version reported during initialization.
const Version = "0.1.0"

// Server wraps an MCP server backed by a yaah Harness.
type Server struct {
	harness   *harness.Harness
	mcpServer *mcp.Server
}

// New creates a new MCP server wired to the given Harness.
// All tools are registered during construction.
func New(h *harness.Harness) *Server {
	s := &Server{
		harness: h,
		mcpServer: mcp.NewServer(&mcp.Implementation{
			Name:    "yaah",
			Version: Version,
		}, nil),
	}
	s.registerTools()
	return s
}

// Start runs the MCP server over stdio, blocking until the context is
// cancelled or the input stream is closed.
func (s *Server) Start(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// registerTools adds all yaah MCP tools to the underlying MCP server.
func (s *Server) registerTools() {
	s.addScanSecretsTool()
	s.addLintTool()
	s.addCheckCommandTool()
	s.addDoctorTool()
	s.addSessionInfoTool()
	s.addPlanningStatusTool()
	s.addPlanningInitTool()
}
