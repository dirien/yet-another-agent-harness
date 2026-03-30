package harness

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dirien/yet-another-agent-harness/pkg/agents"
	"github.com/dirien/yet-another-agent-harness/pkg/commands"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/mcp"
	"github.com/dirien/yet-another-agent-harness/pkg/plugins"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	"github.com/dirien/yet-another-agent-harness/pkg/session"
	"github.com/dirien/yet-another-agent-harness/pkg/skills"
)

// ErrHookBlocked is returned by HandleHookEvent when a handler signals
// that the action should be blocked (e.g. a dangerous command was caught).
var ErrHookBlocked = errors.New("hook blocked the action")

// Harness is the top-level runtime wiring together every component.
type Harness struct {
	hooks        *hooks.Registry
	mcp          *mcp.Registry
	lsp          *lsp.Registry
	skills       *skills.Registry
	agents       *agents.Registry
	commands     *commands.Registry
	plugins      *plugins.Registry
	settings     *schema.Settings
	sessionStore *session.Store
}

// New creates a new Harness with empty registries.
func New() *Harness {
	return &Harness{
		hooks:        hooks.NewRegistry(),
		mcp:          mcp.NewRegistry(),
		lsp:          lsp.NewRegistry(),
		skills:       skills.NewRegistry(),
		agents:       agents.NewRegistry(),
		commands:     commands.NewRegistry(),
		plugins:      plugins.NewRegistry(),
		sessionStore: session.NewStore(filepath.Join(".claude", "sessions")),
	}
}

// Hooks returns the hook registry for registering handlers.
func (p *Harness) Hooks() *hooks.Registry { return p.hooks }

// MCP returns the MCP registry for registering providers.
func (p *Harness) MCP() *mcp.Registry { return p.mcp }

// LSP returns the LSP registry for registering language server providers.
func (p *Harness) LSP() *lsp.Registry { return p.lsp }

// Skills returns the skills registry for registering skills.
func (p *Harness) Skills() *skills.Registry { return p.skills }

// Agents returns the agent registry for registering agents.
func (p *Harness) Agents() *agents.Registry { return p.agents }

// Commands returns the command registry for registering commands.
func (p *Harness) Commands() *commands.Registry { return p.commands }

// Plugins returns the plugin registry for registering plugins.
func (p *Harness) Plugins() *plugins.Registry { return p.plugins }

// SessionStore returns the session store for querying session state.
func (p *Harness) SessionStore() *session.Store { return p.sessionStore }

// SetSettings sets the base Claude Code settings.
func (p *Harness) SetSettings(s *schema.Settings) { p.settings = s }

// HandleHookEvent dispatches a hook event through all enlisted handlers.
// It also records the event in the session store when a session ID is present.
// Returns ErrHookBlocked if any handler signals that the action should be blocked.
func (p *Harness) HandleHookEvent(ctx context.Context, event schema.HookEvent, input *hooks.Input) error {
	// Load session state (if we have a session ID).
	var sess *session.Session
	if input.SessionID != "" {
		var err error
		sess, err = p.sessionStore.Load(input.SessionID)
		if err != nil {
			// Non-fatal: log and continue without session tracking.
			_, _ = fmt.Fprintf(os.Stderr, "session load: %v\n", err)
			sess = nil
		}
	}

	// Record the tool call before dispatching.
	now := time.Now().UTC()
	var record *session.ToolCallRecord
	if sess != nil && input.ToolName != "" {
		record = &session.ToolCallRecord{
			Timestamp: now,
			ToolName:  input.ToolName,
			Input:     summarizeToolInput(input),
		}
	}

	// Dispatch to handlers.
	results, err := p.hooks.Dispatch(ctx, event, input)
	if err != nil {
		return err
	}

	combined := hooks.CombineResults(results)
	blocked := combined.Block

	// Update session state.
	if sess != nil {
		sess.LastEventAt = now
		sess.EventCount++

		if record != nil {
			if blocked {
				record.Blocked = true
				record.Reason = combined.Error
				sess.BlockedCalls = append(sess.BlockedCalls, *record)
			}
			sess.ToolCalls = append(sess.ToolCalls, *record)
		}

		// Track file modifications from Edit/Write tools.
		if fp := input.FilePath(); fp != "" && !blocked {
			if input.ToolName == "Edit" || input.ToolName == "Write" || input.ToolName == "MultiEdit" {
				if !containsString(sess.FilesModified, fp) {
					sess.FilesModified = append(sess.FilesModified, fp)
				}
			}
		}

		// Save session (best-effort).
		if saveErr := p.sessionStore.Save(sess); saveErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "session save: %v\n", saveErr)
		}
	}

	if combined.Output != "" {
		_, _ = fmt.Fprint(os.Stdout, combined.Output)
	}
	if combined.Error != "" {
		_, _ = fmt.Fprint(os.Stderr, combined.Error)
	}
	if blocked {
		return ErrHookBlocked
	}
	return nil
}

// summarizeToolInput returns a short summary of the tool input for session recording.
func summarizeToolInput(input *hooks.Input) string {
	if len(input.ToolInput) == 0 {
		return ""
	}

	// Try to extract file_path first.
	if fp := input.FilePath(); fp != "" {
		return fp
	}

	// Try to extract bash command.
	if cmd := input.BashCommand(); cmd != "" {
		if len(cmd) > 120 {
			return cmd[:120] + "..."
		}
		return cmd
	}

	// Fall back to a truncated JSON representation.
	raw := string(input.ToolInput)
	if len(raw) > 120 {
		return raw[:120] + "..."
	}
	return raw
}

// containsString checks if a string slice contains the given value.
func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// GenerateConfig builds a complete HarnessConfig from all registered components.
func (p *Harness) GenerateConfig() *schema.HarnessConfig {
	cfg := &schema.HarnessConfig{
		Version:  "1",
		Settings: p.settings,
	}

	// Hooks: one rule per event, pointing to `yaah hook <event>`.
	handlers := p.hooks.Handlers()
	if len(handlers) > 0 {
		cfg.Hooks = make(schema.HooksConfig)
		seen := make(map[schema.HookEvent]map[string]bool)

		for _, handler := range handlers {
			for _, event := range handler.Events() {
				if seen[event] == nil {
					seen[event] = make(map[string]bool)
				}
				matcher := ""
				if m := handler.Match(); m != nil {
					matcher = m.String()
				}
				if seen[event][matcher] {
					continue
				}
				seen[event][matcher] = true

				cfg.Hooks[event] = append(cfg.Hooks[event], schema.HookRule{
					Matcher: matcher,
					Hooks: []schema.HookHandler{{
						Type:    schema.HookTypeCommand,
						Command: fmt.Sprintf("yaah hook %s", event),
					}},
				})
			}
		}
	}

	// MCP from providers.
	if providers := p.mcp.Providers(); len(providers) > 0 {
		cfg.MCP = p.mcp.ToConfig()
	}

	// LSP: add marketplace plugins to enabledPlugins.
	if providers := p.lsp.Providers(); len(providers) > 0 {
		for _, prov := range providers {
			if mp, ok := prov.(lsp.MarketplaceProvider); ok {
				if cfg.Settings == nil {
					cfg.Settings = &schema.Settings{}
				}
				if cfg.Settings.EnabledPlugins == nil {
					cfg.Settings.EnabledPlugins = schema.EnabledPluginsMap{}
				}
				cfg.Settings.EnabledPlugins[mp.MarketplaceKey()] = true
			}
		}
	}

	// Plugins: add marketplace plugins to enabledPlugins.
	if registered := p.plugins.Plugins(); len(registered) > 0 {
		for _, plug := range registered {
			if mp, ok := plug.(plugins.MarketplacePlugin); ok {
				if cfg.Settings == nil {
					cfg.Settings = &schema.Settings{}
				}
				if cfg.Settings.EnabledPlugins == nil {
					cfg.Settings.EnabledPlugins = schema.EnabledPluginsMap{}
				}
				cfg.Settings.EnabledPlugins[mp.MarketplaceKey()] = true
			}
		}
	}

	// Skills.
	if registeredSkills := p.skills.Skills(); len(registeredSkills) > 0 {
		cfg.Skills = &schema.SkillsConfig{}
		for _, d := range registeredSkills {
			source := d.Source()
			if source.Path == "" && source.Uses == "" {
				source.Path = filepath.Join(".claude", "skills", d.Name(), "SKILL.md")
			}
			skill := schema.Skill{
				Name:        d.Name(),
				Description: d.Description(),
				Source:      source,
				Enabled:     true,
			}
			if swf, ok := d.(skills.SkillWithFrontmatter); ok {
				fm := swf.Frontmatter()
				skill.ArgumentHint = fm.ArgumentHint
				skill.DisableModelInvocation = fm.DisableModelInvocation
				skill.UserInvocable = fm.UserInvocable
				skill.AllowedTools = fm.AllowedTools
				skill.Model = fm.Model
				skill.Context = fm.Context
				skill.AgentType = fm.AgentType
			}
			if swm, ok := d.(skills.SkillWithMetadata); ok {
				m := swm.Metadata()
				skill.Category = m.Category
				skill.Tags = m.Tags
				skill.Risk = m.Risk
				skill.Tier = m.Tier
			}
			cfg.Skills.Skills = append(cfg.Skills.Skills, skill)
		}
	}

	// Agents.
	if registeredAgents := p.agents.Agents(); len(registeredAgents) > 0 {
		cfg.Agents = p.agents.ToConfig()
	}

	// Commands.
	if registeredCommands := p.commands.Commands(); len(registeredCommands) > 0 {
		cfg.Commands = &schema.CommandsConfig{}
		for _, c := range registeredCommands {
			cmd := schema.Command{
				Name:         c.Name(),
				Description:  c.Description(),
				Source:       filepath.Join(".claude", "commands", c.Name()+".md"),
				ArgumentHint: c.ArgumentHint(),
				Model:        c.Model(),
				AllowedTools: c.AllowedTools(),
			}
			if adv, ok := c.(commands.CommandWithAdvanced); ok {
				opts := adv.Advanced()
				cmd.DisableModelInvocation = opts.DisableModelInvocation
				cmd.UserInvocable = opts.UserInvocable
				cmd.Context = opts.Context
				cmd.AgentType = opts.AgentType
			}
			cfg.Commands.Commands = append(cfg.Commands.Commands, cmd)
		}
	}

	return cfg
}

// TargetWriter provides the directory layout for writing skills, agents, and commands.
// This is satisfied by the AgentGenerator implementations in the generator package.
type TargetWriter interface {
	SkillsDir() string
	AgentsDir() string
	AgentFileExt() string
	CommandsDir() string
}

// AgentToolsFormatter is an optional interface that TargetWriter implementations
// can satisfy to customize how the "tools" frontmatter is rendered in agent files.
// By default, tools are written as a comma-separated string (Claude Code format).
type AgentToolsFormatter interface {
	// FormatAgentTools converts a comma-separated tool allowlist (e.g. "Read, Grep, Glob")
	// into the target-specific frontmatter representation.
	FormatAgentTools(tools string) string
}

// WriteAll writes all generated files (skills, agents, commands, LSP) to baseDir
// using the default Claude Code layout.
func (p *Harness) WriteAll(baseDir string) error {
	return p.WriteAllForTarget(baseDir, &defaultClaudeTarget{})
}

// defaultClaudeTarget provides the Claude Code directory layout for backward compatibility.
type defaultClaudeTarget struct{}

func (t *defaultClaudeTarget) SkillsDir() string    { return ".claude/skills" }
func (t *defaultClaudeTarget) AgentsDir() string    { return ".claude/agents" }
func (t *defaultClaudeTarget) AgentFileExt() string { return ".md" }
func (t *defaultClaudeTarget) CommandsDir() string  { return ".claude/commands" }

// WriteAllForTarget writes skills, agents, and commands using the directory layout
// provided by the given TargetWriter.
func (p *Harness) WriteAllForTarget(baseDir string, tw TargetWriter) error {
	skillsDir := tw.SkillsDir()
	agentsDir := tw.AgentsDir()
	agentExt := tw.AgentFileExt()
	commandsDir := tw.CommandsDir()

	// Skills.
	for _, d := range p.skills.Skills() {
		content := d.Content()
		if content == "" {
			continue
		}
		dir := filepath.Join(baseDir, skillsDir, d.Name())
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create skill dir %s: %w", d.Name(), err)
		}

		// Build skill markdown with optional frontmatter.
		skillMD := buildSkillMarkdown(d)
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
			return fmt.Errorf("write skill %s: %w", d.Name(), err)
		}

		// Write extra files (e.g. references/) if the skill provides them.
		if swf, ok := d.(skills.SkillWithFiles); ok {
			for relPath, fileContent := range swf.ExtraFiles() {
				fp := filepath.Join(dir, relPath)
				if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
					return fmt.Errorf("create dir for %s/%s: %w", d.Name(), relPath, err)
				}
				if err := os.WriteFile(fp, []byte(fileContent), 0o644); err != nil {
					return fmt.Errorf("write %s/%s: %w", d.Name(), relPath, err)
				}
			}
		}

		source := d.Source()
		_, _ = fmt.Fprintf(os.Stderr, "  skill    %-18s", d.Name())
		if source.Uses != "" {
			_, _ = fmt.Fprintf(os.Stderr, " (remote: %s)", source.Uses)
		}
		if swf, ok := d.(skills.SkillWithFiles); ok && len(swf.ExtraFiles()) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, " +%d files", len(swf.ExtraFiles()))
		}
		_, _ = fmt.Fprintln(os.Stderr)
	}

	// Agents (skip if target doesn't support file-based agents).
	if agentsDir != "" {
		for _, a := range p.agents.Agents() {
			dir := filepath.Join(baseDir, agentsDir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("create agents dir: %w", err)
			}
			content := buildAgentMarkdown(a, tw)
			path := filepath.Join(dir, a.Name()+agentExt)
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return fmt.Errorf("write agent %s: %w", a.Name(), err)
			}
			_, _ = fmt.Fprintf(os.Stderr, "  agent    %-18s model=%s", a.Name(), a.Model())
			if src, ok := a.(agents.AgentWithSource); ok && src.Uses() != "" {
				_, _ = fmt.Fprintf(os.Stderr, " (remote: %s)", src.Uses())
			}
			_, _ = fmt.Fprintln(os.Stderr)
		}
	}

	// LSP: enable marketplace plugins via enabledPlugins in settings.json.
	// Binary availability is reported via doctor.
	if lspProviders := p.lsp.Providers(); len(lspProviders) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "  lsp      enabledPlugins (%d servers)\n", len(lspProviders))
		results := p.lsp.CheckAll()
		_, _ = fmt.Fprint(os.Stderr, lsp.FormatCheckResults(results))
	}

	// Commands (skip if target doesn't support commands).
	if commandsDir != "" {
		for _, c := range p.commands.Commands() {
			dir := filepath.Join(baseDir, commandsDir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("create commands dir: %w", err)
			}
			content := buildCommandMarkdown(c)
			path := filepath.Join(dir, c.Name()+".md")
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return fmt.Errorf("write command %s: %w", c.Name(), err)
			}
			_, _ = fmt.Fprintf(os.Stderr, "  command  %-18s\n", c.Name())
		}
	}

	return nil
}

func buildSkillMarkdown(s skills.Skill) string {
	var b strings.Builder

	// Check if skill has frontmatter.
	var hasFrontmatter bool
	var fm skills.SkillFrontmatter
	if swf, ok := s.(skills.SkillWithFrontmatter); ok {
		fm = swf.Frontmatter()
		hasFrontmatter = fm.ArgumentHint != "" || fm.DisableModelInvocation ||
			fm.UserInvocable != nil || fm.AllowedTools != "" ||
			fm.Model != "" || fm.Context != "" || fm.AgentType != ""
	}

	content := s.Content()

	// Skip adding frontmatter if the content already has its own.
	contentHasFrontmatter := strings.HasPrefix(strings.TrimSpace(content), "---")

	if !contentHasFrontmatter {
		// Emit frontmatter — some agents (e.g. Codex) require it.
		b.WriteString("---\n")
		_, _ = fmt.Fprintf(&b, "name: %s\n", s.Name())
		_, _ = fmt.Fprintf(&b, "description: \"%s\"\n", s.Description())
		if hasFrontmatter {
			if fm.ArgumentHint != "" {
				_, _ = fmt.Fprintf(&b, "argument-hint: %s\n", fm.ArgumentHint)
			}
			if fm.DisableModelInvocation {
				b.WriteString("disable-model-invocation: true\n")
			}
			if fm.UserInvocable != nil && !*fm.UserInvocable {
				b.WriteString("user-invocable: false\n")
			}
			if fm.AllowedTools != "" {
				_, _ = fmt.Fprintf(&b, "allowed-tools: %s\n", fm.AllowedTools)
			}
			if fm.Model != "" {
				_, _ = fmt.Fprintf(&b, "model: %s\n", fm.Model)
			}
			if fm.Context != "" {
				_, _ = fmt.Fprintf(&b, "context: %s\n", fm.Context)
			}
			if fm.AgentType != "" {
				_, _ = fmt.Fprintf(&b, "agent: %s\n", fm.AgentType)
			}
		}
		b.WriteString("---\n\n")
	}

	b.WriteString(content)
	return b.String()
}

func buildAgentMarkdown(a agents.Agent, tw TargetWriter) string {
	var b strings.Builder
	b.WriteString("---\n")
	_, _ = fmt.Fprintf(&b, "name: %s\n", a.Name())
	_, _ = fmt.Fprintf(&b, "description: \"%s\"\n", a.Description())
	if a.Model() != "" {
		_, _ = fmt.Fprintf(&b, "model: %s\n", a.Model())
	}
	if a.Tools() != "" {
		if f, ok := tw.(AgentToolsFormatter); ok {
			b.WriteString(f.FormatAgentTools(a.Tools()))
		} else {
			_, _ = fmt.Fprintf(&b, "tools: %s\n", a.Tools())
		}
	}

	// Advanced fields.
	if adv, ok := a.(agents.AgentWithAdvanced); ok {
		opts := adv.Advanced()
		if opts.DisallowedTools != "" {
			_, _ = fmt.Fprintf(&b, "disallowedTools: %s\n", opts.DisallowedTools)
		}
		if opts.PermissionMode != "" {
			_, _ = fmt.Fprintf(&b, "permissionMode: %s\n", opts.PermissionMode)
		}
		if opts.MaxTurns > 0 {
			_, _ = fmt.Fprintf(&b, "maxTurns: %d\n", opts.MaxTurns)
		}
		if len(opts.Skills) > 0 {
			_, _ = fmt.Fprintf(&b, "skills: %s\n", strings.Join(opts.Skills, ", "))
		}
		if opts.Memory != "" {
			_, _ = fmt.Fprintf(&b, "memory: %s\n", opts.Memory)
		}
		if opts.Background {
			b.WriteString("background: true\n")
		}
		if opts.Isolation != "" {
			_, _ = fmt.Fprintf(&b, "isolation: %s\n", opts.Isolation)
		}
	}

	b.WriteString("---\n\n")
	b.WriteString(a.Content())
	return b.String()
}

func buildCommandMarkdown(c commands.Command) string {
	var b strings.Builder
	b.WriteString("---\n")
	_, _ = fmt.Fprintf(&b, "description: \"%s\"\n", c.Description())
	if c.ArgumentHint() != "" {
		_, _ = fmt.Fprintf(&b, "argument-hint: %s\n", c.ArgumentHint())
	}
	if c.Model() != "" {
		_, _ = fmt.Fprintf(&b, "model: %s\n", c.Model())
	}
	if c.AllowedTools() != "" {
		_, _ = fmt.Fprintf(&b, "allowed-tools: %s\n", c.AllowedTools())
	}

	// Advanced fields.
	if adv, ok := c.(commands.CommandWithAdvanced); ok {
		opts := adv.Advanced()
		if opts.DisableModelInvocation {
			b.WriteString("disable-model-invocation: true\n")
		}
		if opts.UserInvocable != nil && !*opts.UserInvocable {
			b.WriteString("user-invocable: false\n")
		}
		if opts.Context != "" {
			_, _ = fmt.Fprintf(&b, "context: %s\n", opts.Context)
		}
		if opts.AgentType != "" {
			_, _ = fmt.Fprintf(&b, "agent: %s\n", opts.AgentType)
		}
	}

	b.WriteString("---\n\n")
	b.WriteString(c.Content())
	return b.String()
}

// Summary returns a human-readable summary of all registered components.
func (p *Harness) Summary() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Hooks:    %d handlers", len(p.hooks.Handlers())))
	lines = append(lines, fmt.Sprintf("MCP:      %d providers", len(p.mcp.Providers())))
	lines = append(lines, fmt.Sprintf("LSP:      %d servers", len(p.lsp.Providers())))
	lines = append(lines, fmt.Sprintf("Skills:   %d skills", len(p.skills.Skills())))
	lines = append(lines, fmt.Sprintf("Agents:   %d agents", len(p.agents.Agents())))
	lines = append(lines, fmt.Sprintf("Commands: %d commands", len(p.commands.Commands())))
	lines = append(lines, fmt.Sprintf("Plugins:  %d plugins", len(p.plugins.Plugins())))

	if handlers := p.hooks.Handlers(); len(handlers) > 0 {
		lines = append(lines, "\nHandlers:")
		for _, s := range handlers {
			events := make([]string, len(s.Events()))
			for i, e := range s.Events() {
				events[i] = string(e)
			}
			match := "*"
			if m := s.Match(); m != nil {
				match = m.String()
			}
			lines = append(lines, fmt.Sprintf("  %-20s events=[%s] match=%s", s.Name(), strings.Join(events, ","), match))
		}
	}

	if providers := p.mcp.Providers(); len(providers) > 0 {
		lines = append(lines, "\nMCP Providers:")
		for _, e := range providers {
			srv := e.Server()
			lines = append(lines, fmt.Sprintf("  %-20s transport=%s", e.Name(), srv.Transport))
		}
	}

	if lspProviders := p.lsp.Providers(); len(lspProviders) > 0 {
		lines = append(lines, "\nLSP Servers:")
		results := p.lsp.CheckAll()
		for _, cr := range results {
			if cr.Installed {
				lines = append(lines, fmt.Sprintf("  %-20s ✓ %s", cr.Name, cr.BinaryPath))
			} else {
				lines = append(lines, fmt.Sprintf("  %-20s ✗ not found → %s", cr.Name, cr.InstallHint))
			}
		}
	}

	if registeredSkills := p.skills.Skills(); len(registeredSkills) > 0 {
		lines = append(lines, "\nSkills:")
		for _, d := range registeredSkills {
			lines = append(lines, fmt.Sprintf("  %-20s %s", d.Name(), d.Description()))
		}
	}

	if registeredAgents := p.agents.Agents(); len(registeredAgents) > 0 {
		lines = append(lines, "\nAgents:")
		for _, a := range registeredAgents {
			lines = append(lines, fmt.Sprintf("  %-20s model=%s", a.Name(), a.Model()))
		}
	}

	if registeredCommands := p.commands.Commands(); len(registeredCommands) > 0 {
		lines = append(lines, "\nCommands:")
		for _, c := range registeredCommands {
			lines = append(lines, fmt.Sprintf("  %-20s %s", c.Name(), c.Description()))
		}
	}

	if registeredPlugins := p.plugins.Plugins(); len(registeredPlugins) > 0 {
		lines = append(lines, "\nPlugins:")
		for _, plug := range registeredPlugins {
			meta := plug.Plugin()
			key := ""
			if mp, ok := plug.(plugins.MarketplacePlugin); ok {
				key = fmt.Sprintf(" (marketplace: %s)", mp.MarketplaceKey())
			}
			lines = append(lines, fmt.Sprintf("  %-20s %s%s", meta.Name, meta.Description, key))
		}
	}

	return strings.Join(lines, "\n")
}
