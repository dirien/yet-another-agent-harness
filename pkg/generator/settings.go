package generator

import (
	"encoding/json"
	"fmt"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// ClaudeSettings is the output format for .claude/settings.json.
type ClaudeSettings struct {
	Schema                string                          `json:"$schema,omitempty"`
	Model                 string                          `json:"model,omitempty"`
	AlwaysThinkingEnabled *bool                           `json:"alwaysThinkingEnabled,omitempty"`
	AutoUpdatesChannel    string                          `json:"autoUpdatesChannel,omitempty"`
	EffortLevel           string                          `json:"effortLevel,omitempty"`
	StatusLine            *schema.StatusLine              `json:"statusLine,omitempty"`
	Env                   map[string]string               `json:"env,omitempty"`
	TeammateMode          string                          `json:"teammateMode,omitempty"`
	Hooks                 map[string][]claudeSettingsHook `json:"hooks,omitempty"`
	McpServers            map[string]claudeMCPServer      `json:"mcpServers,omitempty"`

	// Model & Performance
	AvailableModels         []string `json:"availableModels,omitempty"`
	FastMode                *bool    `json:"fastMode,omitempty"`
	FastModePerSessionOptIn *bool    `json:"fastModePerSessionOptIn,omitempty"`

	// Permissions & Security
	Permissions                     *schema.PermissionsConfig `json:"permissions,omitempty"`
	Sandbox                         string                    `json:"sandbox,omitempty"`
	AllowManagedPermissionRulesOnly *bool                     `json:"allowManagedPermissionRulesOnly,omitempty"`

	// Hooks & Automation
	DisableAllHooks       *bool `json:"disableAllHooks,omitempty"`
	AllowManagedHooksOnly *bool `json:"allowManagedHooksOnly,omitempty"`

	// Git & Attribution
	Attribution            string `json:"attribution,omitempty"`
	IncludeGitInstructions *bool  `json:"includeGitInstructions,omitempty"`

	// Authentication
	ApiKeyHelper string `json:"apiKeyHelper,omitempty"`

	// UI & Behavior
	Language                   string   `json:"language,omitempty"`
	OutputStyle                string   `json:"outputStyle,omitempty"`
	ShowTurnDuration           *bool    `json:"showTurnDuration,omitempty"`
	SpinnerVerbs               []string `json:"spinnerVerbs,omitempty"`
	SpinnerTipsEnabled         *bool    `json:"spinnerTipsEnabled,omitempty"`
	SpinnerTipsOverride        []string `json:"spinnerTipsOverride,omitempty"`
	PrefersReducedMotion       *bool    `json:"prefersReducedMotion,omitempty"`
	TerminalProgressBarEnabled *bool    `json:"terminalProgressBarEnabled,omitempty"`

	// Plugins
	EnabledPlugins          schema.EnabledPluginsMap `json:"enabledPlugins,omitempty"`
	PluginConfigs           map[string]any           `json:"pluginConfigs,omitempty"`
	ExtraKnownMarketplaces  []string                 `json:"extraKnownMarketplaces,omitempty"`
	StrictKnownMarketplaces []string                 `json:"strictKnownMarketplaces,omitempty"`
	SkippedMarketplaces     []string                 `json:"skippedMarketplaces,omitempty"`
	SkippedPlugins          []string                 `json:"skippedPlugins,omitempty"`
	BlockedMarketplaces     []string                 `json:"blockedMarketplaces,omitempty"`
	PluginTrustMessage      string                   `json:"pluginTrustMessage,omitempty"`

	// MCP Management
	EnableAllProjectMcpServers *bool    `json:"enableAllProjectMcpServers,omitempty"`
	EnabledMcpjsonServers      []string `json:"enabledMcpjsonServers,omitempty"`
	DisabledMcpjsonServers     []string `json:"disabledMcpjsonServers,omitempty"`
	AllowedMcpServers          []string `json:"allowedMcpServers,omitempty"`
	DeniedMcpServers           []string `json:"deniedMcpServers,omitempty"`
	AllowManagedMcpServersOnly *bool    `json:"allowManagedMcpServersOnly,omitempty"`

	// Organization
	CompanyAnnouncements  []string `json:"companyAnnouncements,omitempty"`
	CleanupPeriodDays     int      `json:"cleanupPeriodDays,omitempty"`
	PlansDirectory        string   `json:"plansDirectory,omitempty"`
	AutoMemoryEnabled     *bool    `json:"autoMemoryEnabled,omitempty"`
	SkipWebFetchPreflight *bool    `json:"skipWebFetchPreflight,omitempty"`

	// File & Directory
	FileSuggestion        string   `json:"fileSuggestion,omitempty"`
	RespectGitignore      *bool    `json:"respectGitignore,omitempty"`
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`
}

type claudeSettingsHook struct {
	Matcher string                      `json:"matcher,omitempty"`
	Hooks   []claudeSettingsHookHandler `json:"hooks"`
}

type claudeSettingsHookHandler struct {
	Type           string            `json:"type"`
	Command        string            `json:"command,omitempty"`
	URL            string            `json:"url,omitempty"`
	Prompt         string            `json:"prompt,omitempty"`
	Timeout        int               `json:"timeout,omitempty"`
	StatusMessage  string            `json:"statusMessage,omitempty"`
	Once           bool              `json:"once,omitempty"`
	Async          bool              `json:"async,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	AllowedEnvVars []string          `json:"allowedEnvVars,omitempty"`
	Model          string            `json:"model,omitempty"`
}

type claudeMCPServer struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	OAuth   *schema.MCPOAuth  `json:"oauth,omitempty"`
}

// GenerateClaudeSettings converts a HarnessConfig into Claude Code's settings.json format.
func GenerateClaudeSettings(cfg *schema.HarnessConfig) ([]byte, error) {
	out := ClaudeSettings{
		Schema: "https://json.schemastore.org/claude-code-settings.json",
	}

	if s := cfg.Settings; s != nil {
		// Core
		out.Model = s.Model
		out.AlwaysThinkingEnabled = s.AlwaysThinkingEnabled
		out.AutoUpdatesChannel = s.AutoUpdatesChannel
		out.EffortLevel = s.EffortLevel
		out.StatusLine = s.StatusLine
		out.Env = s.Env
		out.TeammateMode = s.TeammateMode
		// Model & Performance
		out.AvailableModels = s.AvailableModels
		out.FastMode = s.FastMode
		out.FastModePerSessionOptIn = s.FastModePerSessionOptIn
		// Permissions & Security
		out.Permissions = s.Permissions
		out.Sandbox = s.Sandbox
		out.AllowManagedPermissionRulesOnly = s.AllowManagedPermissionRulesOnly
		// Hooks & Automation
		out.DisableAllHooks = s.DisableAllHooks
		out.AllowManagedHooksOnly = s.AllowManagedHooksOnly
		// Git & Attribution
		out.Attribution = s.Attribution
		out.IncludeGitInstructions = s.IncludeGitInstructions
		// Authentication
		out.ApiKeyHelper = s.ApiKeyHelper
		// UI & Behavior
		out.Language = s.Language
		out.OutputStyle = s.OutputStyle
		out.ShowTurnDuration = s.ShowTurnDuration
		out.SpinnerVerbs = s.SpinnerVerbs
		out.SpinnerTipsEnabled = s.SpinnerTipsEnabled
		out.SpinnerTipsOverride = s.SpinnerTipsOverride
		out.PrefersReducedMotion = s.PrefersReducedMotion
		out.TerminalProgressBarEnabled = s.TerminalProgressBarEnabled
		// Plugins
		out.EnabledPlugins = s.EnabledPlugins
		out.PluginConfigs = s.PluginConfigs
		out.ExtraKnownMarketplaces = s.ExtraKnownMarketplaces
		out.StrictKnownMarketplaces = s.StrictKnownMarketplaces
		out.SkippedMarketplaces = s.SkippedMarketplaces
		out.SkippedPlugins = s.SkippedPlugins
		out.BlockedMarketplaces = s.BlockedMarketplaces
		out.PluginTrustMessage = s.PluginTrustMessage
		// MCP Management
		out.EnableAllProjectMcpServers = s.EnableAllProjectMcpServers
		out.EnabledMcpjsonServers = s.EnabledMcpjsonServers
		out.DisabledMcpjsonServers = s.DisabledMcpjsonServers
		out.AllowedMcpServers = s.AllowedMcpServers
		out.DeniedMcpServers = s.DeniedMcpServers
		out.AllowManagedMcpServersOnly = s.AllowManagedMcpServersOnly
		// Organization
		out.CompanyAnnouncements = s.CompanyAnnouncements
		out.CleanupPeriodDays = s.CleanupPeriodDays
		out.PlansDirectory = s.PlansDirectory
		out.AutoMemoryEnabled = s.AutoMemoryEnabled
		out.SkipWebFetchPreflight = s.SkipWebFetchPreflight
		// File & Directory
		out.FileSuggestion = s.FileSuggestion
		out.RespectGitignore = s.RespectGitignore
		out.AdditionalDirectories = s.AdditionalDirectories
	}

	if cfg.Hooks != nil {
		out.Hooks = make(map[string][]claudeSettingsHook)
		for event, rules := range cfg.Hooks {
			var hookRules []claudeSettingsHook
			for _, rule := range rules {
				var handlers []claudeSettingsHookHandler
				for _, h := range rule.Hooks {
					handlers = append(handlers, claudeSettingsHookHandler{
						Type:           string(h.Type),
						Command:        h.Command,
						URL:            h.URL,
						Prompt:         h.Prompt,
						Timeout:        h.Timeout,
						StatusMessage:  h.StatusMessage,
						Once:           h.Once,
						Async:          h.Async,
						Headers:        h.Headers,
						AllowedEnvVars: h.AllowedEnvVars,
						Model:          h.Model,
					})
				}
				hookRules = append(hookRules, claudeSettingsHook{
					Matcher: rule.Matcher,
					Hooks:   handlers,
				})
			}
			out.Hooks[string(event)] = hookRules
		}
	}

	if cfg.MCP != nil {
		out.McpServers = make(map[string]claudeMCPServer)
		for _, srv := range cfg.MCP.Servers {
			out.McpServers[srv.Name] = claudeMCPServer{
				Type:    string(srv.Transport),
				Command: srv.Command,
				Args:    srv.Args,
				URL:     srv.URL,
				Env:     srv.Env,
				Headers: srv.Headers,
				OAuth:   srv.OAuth,
			}
		}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal settings: %w", err)
	}
	return data, nil
}

// mcpJSON is the format for .mcp.json (project-level MCP server config).
type mcpJSON struct {
	McpServers map[string]claudeMCPServer `json:"mcpServers"`
}

// GenerateMCPJSON produces .mcp.json content from the MCP config.
// This is the project-level MCP config that Claude Code discovers automatically.
func GenerateMCPJSON(cfg *schema.HarnessConfig) ([]byte, error) {
	if cfg.MCP == nil || len(cfg.MCP.Servers) == 0 {
		return nil, nil
	}

	out := mcpJSON{
		McpServers: make(map[string]claudeMCPServer),
	}
	for _, srv := range cfg.MCP.Servers {
		out.McpServers[srv.Name] = claudeMCPServer{
			Type:    string(srv.Transport),
			Command: srv.Command,
			Args:    srv.Args,
			URL:     srv.URL,
			Env:     srv.Env,
			Headers: srv.Headers,
			OAuth:   srv.OAuth,
		}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal mcp.json: %w", err)
	}
	return data, nil
}

// claudeLSPServer is the per-server format for lsp.json.
type claudeLSPServer struct {
	Command               []string          `json:"command"`
	ExtensionToLanguage   map[string]string `json:"extensionToLanguage"`
	Env                   map[string]string `json:"env,omitempty"`
	Transport             string            `json:"transport,omitempty"`
	InitializationOptions map[string]any    `json:"initializationOptions,omitempty"`
	Settings              map[string]any    `json:"settings,omitempty"`
	WorkspaceFolder       string            `json:"workspaceFolder,omitempty"`
	StartupTimeout        int               `json:"startupTimeout,omitempty"`
	ShutdownTimeout       int               `json:"shutdownTimeout,omitempty"`
	RestartOnCrash        *bool             `json:"restartOnCrash,omitempty"`
	MaxRestarts           int               `json:"maxRestarts,omitempty"`
}

// GenerateLSPConfig produces lsp.json content from LSP server definitions.
func GenerateLSPConfig(cfg *schema.HarnessConfig) ([]byte, error) {
	if cfg.LSP == nil || len(cfg.LSP.Servers) == 0 {
		return nil, nil
	}

	out := make(map[string]claudeLSPServer)
	for _, srv := range cfg.LSP.Servers {
		out[srv.ID] = claudeLSPServer{
			Command:               srv.Command,
			ExtensionToLanguage:   srv.ExtensionToLanguage,
			Env:                   srv.Env,
			Transport:             srv.Transport,
			InitializationOptions: srv.InitializationOptions,
			Settings:              srv.Settings,
			WorkspaceFolder:       srv.WorkspaceFolder,
			StartupTimeout:        srv.StartupTimeout,
			ShutdownTimeout:       srv.ShutdownTimeout,
			RestartOnCrash:        srv.RestartOnCrash,
			MaxRestarts:           srv.MaxRestarts,
		}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal lsp config: %w", err)
	}
	return data, nil
}
