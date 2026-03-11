package schema

// StatusLineType defines how the status line is rendered.
type StatusLineType string

const (
	StatusLineCommand StatusLineType = "command"
	StatusLineStatic  StatusLineType = "static"
)

// StatusLine configures the Claude Code status line.
type StatusLine struct {
	Type    StatusLineType `json:"type"              jsonschema:"enum=command,enum=static"`
	Command string         `json:"command,omitempty" jsonschema:"description=Shell command producing status line output"`
	Text    string         `json:"text,omitempty"    jsonschema:"description=Static text for the status line"`
	Padding int            `json:"padding,omitempty" jsonschema:"description=Padding around the status line"`
}

// EnabledPluginsMap controls which Claude Code plugins are enabled.
type EnabledPluginsMap map[string]bool

// Settings holds top-level Claude Code settings.
// Fields align with the official Claude Code settings.json specification.
type Settings struct {
	// Core
	Model                 string            `json:"model,omitempty"                        jsonschema:"description=Default model (opus/sonnet/haiku)"`
	AlwaysThinkingEnabled *bool             `json:"alwaysThinkingEnabled,omitempty"        jsonschema:"description=Enable extended thinking"`
	AutoUpdatesChannel    string            `json:"autoUpdatesChannel,omitempty"           jsonschema:"description=Update channel (latest/stable)"`
	EffortLevel           string            `json:"effortLevel,omitempty"                  jsonschema:"description=Reasoning effort (low/medium/high),enum=low,enum=medium,enum=high"`
	StatusLine            *StatusLine       `json:"statusLine,omitempty"                   jsonschema:"description=Status line configuration"`
	Env                   map[string]string `json:"env,omitempty"                          jsonschema:"description=Environment variables for Claude Code"`
	TeammateMode          string            `json:"teammateMode,omitempty"                 jsonschema:"description=Agent team mode (tmux/etc)"`

	// Model & Performance
	AvailableModels         []string `json:"availableModels,omitempty"              jsonschema:"description=List of available models to choose from"`
	FastMode                *bool    `json:"fastMode,omitempty"                     jsonschema:"description=Enable fast mode for quicker responses"`
	FastModePerSessionOptIn *bool    `json:"fastModePerSessionOptIn,omitempty"      jsonschema:"description=Require per-session opt-in for fast mode"`

	// Permissions & Security
	Permissions                     *PermissionsConfig `json:"permissions,omitempty"                       jsonschema:"description=Tool permission rules (allow/deny)"`
	Sandbox                         string             `json:"sandbox,omitempty"                           jsonschema:"description=Sandbox mode (auto/permissive)"`
	AllowManagedPermissionRulesOnly *bool              `json:"allowManagedPermissionRulesOnly,omitempty"   jsonschema:"description=Only allow managed permission rules"`

	// Hooks & Automation
	DisableAllHooks       *bool `json:"disableAllHooks,omitempty"               jsonschema:"description=Disable all hooks globally"`
	AllowManagedHooksOnly *bool `json:"allowManagedHooksOnly,omitempty"         jsonschema:"description=Only allow managed hooks"`

	// Git & Attribution
	Attribution            string `json:"attribution,omitempty"                  jsonschema:"description=Git attribution mode"`
	IncludeGitInstructions *bool  `json:"includeGitInstructions,omitempty"       jsonschema:"description=Include git instructions in context"`

	// Authentication
	ApiKeyHelper string `json:"apiKeyHelper,omitempty"                 jsonschema:"description=Command to retrieve API key"`

	// UI & Behavior
	Language                   string   `json:"language,omitempty"                     jsonschema:"description=Language for Claude responses"`
	OutputStyle                string   `json:"outputStyle,omitempty"                  jsonschema:"description=Output formatting style"`
	ShowTurnDuration           *bool    `json:"showTurnDuration,omitempty"             jsonschema:"description=Show turn duration in output"`
	SpinnerVerbs               []string `json:"spinnerVerbs,omitempty"                 jsonschema:"description=Custom spinner verb list"`
	SpinnerTipsEnabled         *bool    `json:"spinnerTipsEnabled,omitempty"           jsonschema:"description=Enable tips during spinner"`
	SpinnerTipsOverride        []string `json:"spinnerTipsOverride,omitempty"          jsonschema:"description=Override default spinner tips"`
	PrefersReducedMotion       *bool    `json:"prefersReducedMotion,omitempty"         jsonschema:"description=Reduce animation and motion"`
	TerminalProgressBarEnabled *bool    `json:"terminalProgressBarEnabled,omitempty"   jsonschema:"description=Enable terminal progress bar"`

	// Plugins
	EnabledPlugins          EnabledPluginsMap `json:"enabledPlugins,omitempty"                jsonschema:"description=Plugin enable/disable map"`
	PluginConfigs           map[string]any    `json:"pluginConfigs,omitempty"                 jsonschema:"description=Per-plugin configuration"`
	ExtraKnownMarketplaces  []string          `json:"extraKnownMarketplaces,omitempty"        jsonschema:"description=Additional marketplace URLs"`
	StrictKnownMarketplaces []string          `json:"strictKnownMarketplaces,omitempty"       jsonschema:"description=Managed marketplace allowlist"`
	SkippedMarketplaces     []string          `json:"skippedMarketplaces,omitempty"           jsonschema:"description=Marketplaces to skip"`
	SkippedPlugins          []string          `json:"skippedPlugins,omitempty"                jsonschema:"description=Plugins to skip"`
	BlockedMarketplaces     []string          `json:"blockedMarketplaces,omitempty"           jsonschema:"description=Managed blocked marketplaces"`
	PluginTrustMessage      string            `json:"pluginTrustMessage,omitempty"            jsonschema:"description=Custom trust message for plugins"`

	// MCP Management
	EnableAllProjectMcpServers *bool    `json:"enableAllProjectMcpServers,omitempty"   jsonschema:"description=Auto-enable all project MCP servers"`
	EnabledMcpjsonServers      []string `json:"enabledMcpjsonServers,omitempty"        jsonschema:"description=Explicitly enabled .mcp.json servers"`
	DisabledMcpjsonServers     []string `json:"disabledMcpjsonServers,omitempty"       jsonschema:"description=Explicitly disabled .mcp.json servers"`
	AllowedMcpServers          []string `json:"allowedMcpServers,omitempty"            jsonschema:"description=Allowlist of MCP server names"`
	DeniedMcpServers           []string `json:"deniedMcpServers,omitempty"             jsonschema:"description=Denylist of MCP server names"`
	AllowManagedMcpServersOnly *bool    `json:"allowManagedMcpServersOnly,omitempty"   jsonschema:"description=Only allow managed MCP servers"`

	// Organization
	CompanyAnnouncements  []string `json:"companyAnnouncements,omitempty"         jsonschema:"description=Company-wide announcements"`
	CleanupPeriodDays     int      `json:"cleanupPeriodDays,omitempty"            jsonschema:"description=Days before cleaning up old sessions"`
	PlansDirectory        string   `json:"plansDirectory,omitempty"               jsonschema:"description=Directory for plan files"`
	AutoMemoryEnabled     *bool    `json:"autoMemoryEnabled,omitempty"            jsonschema:"description=Enable automatic memory"`
	SkipWebFetchPreflight *bool    `json:"skipWebFetchPreflight,omitempty"        jsonschema:"description=Skip web fetch preflight checks"`

	// File & Directory
	FileSuggestion        string   `json:"fileSuggestion,omitempty"               jsonschema:"description=File suggestion mode"`
	RespectGitignore      *bool    `json:"respectGitignore,omitempty"             jsonschema:"description=Respect .gitignore when searching"`
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"        jsonschema:"description=Additional directories to include"`
}
