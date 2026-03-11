package schema

// MCPTransport defines how the MCP server communicates.
type MCPTransport string

const (
	MCPTransportStdio          MCPTransport = "stdio"
	MCPTransportSSE            MCPTransport = "sse"
	MCPTransportStreamableHTTP MCPTransport = "streamable-http"
	MCPTransportHTTP           MCPTransport = "http"
	MCPTransportWebSocket      MCPTransport = "websocket"
)

// MCPOAuth configures OAuth authentication for remote MCP servers.
type MCPOAuth struct {
	ClientID              string `json:"clientId,omitempty"              jsonschema:"description=OAuth client ID"`
	CallbackPort          int    `json:"callbackPort,omitempty"         jsonschema:"description=Local port for OAuth callback"`
	AuthServerMetadataURL string `json:"authServerMetadataUrl,omitempty" jsonschema:"description=OAuth server metadata discovery URL"`
}

// MCPServer describes a single MCP server configuration.
type MCPServer struct {
	Name      string            `json:"name"                jsonschema:"description=Unique server identifier"`
	Transport MCPTransport      `json:"transport"           jsonschema:"enum=stdio,enum=sse,enum=streamable-http,enum=http,enum=websocket"`
	Command   string            `json:"command,omitempty"   jsonschema:"description=Command to start the server (stdio transport)"`
	Args      []string          `json:"args,omitempty"      jsonschema:"description=Arguments for the command"`
	URL       string            `json:"url,omitempty"       jsonschema:"description=Endpoint URL (sse/streamable-http/websocket transport)"`
	Env       map[string]string `json:"env,omitempty"       jsonschema:"description=Environment variables passed to the server"`
	Headers   map[string]string `json:"headers,omitempty"   jsonschema:"description=HTTP headers (sse/streamable-http transport)"`
	OAuth     *MCPOAuth         `json:"oauth,omitempty"     jsonschema:"description=OAuth configuration for remote servers"`
}

// MCPConfig holds all MCP server definitions.
type MCPConfig struct {
	Servers []MCPServer `json:"servers" jsonschema:"description=List of MCP servers to register"`
}
