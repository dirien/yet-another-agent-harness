package schema

// LSPServer describes a single LSP server configuration.
// The output format matches Claude Code's lsp.json specification.
type LSPServer struct {
	ID                    string            `json:"id"                              jsonschema:"description=Unique server identifier (e.g. gopls)"`
	Command               []string          `json:"command"                         jsonschema:"description=Command and arguments to start the server"`
	ExtensionToLanguage   map[string]string `json:"extensionToLanguage"             jsonschema:"description=Maps file extensions to language identifiers (e.g. .go -> go)"`
	Env                   map[string]string `json:"env,omitempty"                   jsonschema:"description=Environment variables passed to the server"`
	Transport             string            `json:"transport,omitempty"             jsonschema:"description=Communication transport: stdio (default) or socket"`
	InitializationOptions map[string]any    `json:"initializationOptions,omitempty"  jsonschema:"description=Options passed to the server during LSP initialize"`
	Settings              map[string]any    `json:"settings,omitempty"              jsonschema:"description=Settings passed via workspace/didChangeConfiguration"`
	WorkspaceFolder       string            `json:"workspaceFolder,omitempty"       jsonschema:"description=Workspace folder path for the server"`
	StartupTimeout        int               `json:"startupTimeout,omitempty"        jsonschema:"description=Max time to wait for server startup in milliseconds"`
	ShutdownTimeout       int               `json:"shutdownTimeout,omitempty"       jsonschema:"description=Max time to wait for graceful shutdown in milliseconds"`
	RestartOnCrash        *bool             `json:"restartOnCrash,omitempty"        jsonschema:"description=Whether to automatically restart the server if it crashes"`
	MaxRestarts           int               `json:"maxRestarts,omitempty"           jsonschema:"description=Maximum number of restart attempts before giving up"`
}

// LSPConfig holds all LSP server definitions.
type LSPConfig struct {
	Servers []LSPServer `json:"servers" jsonschema:"description=List of LSP servers to register"`
}
