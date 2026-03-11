package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	harnesspkg "github.com/dirien/yet-another-agent-harness/pkg/harness"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/mcpserver"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	"github.com/spf13/cobra"
)

// Set via -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// newHarness creates the default Harness with all built-in handlers.
// In a consumer repo, you'd replace this with your own setup.
func newHarness() *harnesspkg.Harness {
	return harnesspkg.NewWithDefaults(harnesspkg.AllDefaults())
}

func main() {
	root := &cobra.Command{
		Use:   "yaah",
		Short: "yaah — yet another agent harness for Claude Code",
	}

	root.AddCommand(
		generateCmd(),
		hookCmd(),
		infoCmd(),
		versionCmd(),
		doctorCmd(),
		serveCmd(),
		sessionCmd(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func generateCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate .claude/ directory from built-in defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := newHarness()
			cfg := p.GenerateConfig()

			if outputDir == "" {
				outputDir = "."
			}
			if err := p.WriteAll(outputDir); err != nil {
				return err
			}

			data, err := generator.GenerateClaudeSettings(cfg)
			if err != nil {
				return err
			}

			claudeDir := filepath.Join(outputDir, ".claude")
			if err := os.MkdirAll(claudeDir, 0o755); err != nil {
				return fmt.Errorf("create .claude dir: %w", err)
			}

			outPath := filepath.Join(claudeDir, "settings.json")
			if err := os.WriteFile(outPath, data, 0o644); err != nil {
				return fmt.Errorf("write settings: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated %s\n", outPath)

			// Generate .mcp.json for project-level MCP server discovery.
			mcpData, err := generator.GenerateMCPJSON(cfg)
			if err != nil {
				return err
			}
			if mcpData != nil {
				mcpPath := filepath.Join(outputDir, ".mcp.json")
				if err := os.WriteFile(mcpPath, mcpData, 0o644); err != nil {
					return fmt.Errorf("write .mcp.json: %w", err)
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated %s\n", mcpPath)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output base directory (default: current directory)")
	return cmd
}

// hookCmd is the runtime dispatcher: `yaah hook <event>` reads stdin and
// dispatches through all registered handlers for that event.
func hookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook <event>",
		Short: "Run as a Claude Code hook — dispatches to registered handlers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			event := schema.HookEvent(args[0])

			input, err := hooks.ReadInput(os.Stdin)
			if err != nil {
				return err
			}

			p := newHarness()
			if err := p.HandleHookEvent(context.Background(), event, input); err != nil {
				if errors.Is(err, harnesspkg.ErrHookBlocked) {
					os.Exit(2)
				}
				return err
			}
			return nil
		},
	}
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "yaah %s (commit: %s, built: %s)\n", version, commit, date)
		},
	}
}

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show all registered handlers, providers, and skills",
		Run: func(cmd *cobra.Command, args []string) {
			p := newHarness()
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), p.Summary())
		},
	}
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the yaah MCP server over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			p := newHarness()
			srv := mcpserver.New(p)
			return srv.Start(ctx)
		},
	}
}

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the health of your yaah setup and report missing dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			p := newHarness()
			issues := 0

			_, _ = fmt.Fprintln(out, "yaah doctor")
			_, _ = fmt.Fprintln(out, "===========")
			_, _ = fmt.Fprintf(out, "Version: %s (commit: %s, built: %s)\n\n", version, commit, date)

			// 1. yaah binary
			_, _ = fmt.Fprintln(out, "Binary:")
			if path, err := exec.LookPath("yaah"); err == nil {
				_, _ = fmt.Fprintf(out, "  ✓ yaah                     %s\n", path)
			} else {
				_, _ = fmt.Fprintln(out, "  ✗ yaah                     not in PATH (running from local build)")
				issues++
			}

			// 2. git (needed for remote skills)
			if path, err := exec.LookPath("git"); err == nil {
				_, _ = fmt.Fprintf(out, "  ✓ git                      %s\n", path)
			} else {
				_, _ = fmt.Fprintln(out, "  ✗ git                      not found (required for remote skills)")
				issues++
			}
			_, _ = fmt.Fprintln(out)

			// 3. .claude/settings.json validation
			_, _ = fmt.Fprintln(out, "Config:")
			settingsPath := filepath.Join(".claude", "settings.json")
			if data, err := os.ReadFile(settingsPath); err == nil {
				if json.Valid(data) {
					_, _ = fmt.Fprintf(out, "  ✓ %-25s valid JSON\n", settingsPath)
				} else {
					_, _ = fmt.Fprintf(out, "  ✗ %-25s invalid JSON\n", settingsPath)
					issues++
				}
			} else {
				_, _ = fmt.Fprintf(out, "  ✗ %-25s not found (run 'yaah generate')\n", settingsPath)
				issues++
			}
			_, _ = fmt.Fprintln(out)

			// 4. LSP servers
			lspResults := p.LSP().CheckAll()
			if len(lspResults) > 0 {
				_, _ = fmt.Fprintln(out, "LSP Servers:")
				_, _ = fmt.Fprint(out, lsp.FormatCheckResults(lspResults))
				for _, cr := range lspResults {
					if !cr.Installed {
						issues++
					}
				}
				_, _ = fmt.Fprintln(out)
			}

			// 5. MCP servers
			mcpProviders := p.MCP().Providers()
			if len(mcpProviders) > 0 {
				_, _ = fmt.Fprintln(out, "MCP Servers:")
				for _, prov := range mcpProviders {
					srv := prov.Server()
					if srv.Command == "" {
						_, _ = fmt.Fprintf(out, "  ✓ %-24s %s (remote)\n", prov.Name(), srv.URL)
						continue
					}
					if path, err := exec.LookPath(srv.Command); err == nil {
						_, _ = fmt.Fprintf(out, "  ✓ %-24s %s\n", prov.Name(), path)
					} else {
						_, _ = fmt.Fprintf(out, "  ✗ %-24s %s not found\n", prov.Name(), srv.Command)
						issues++
					}
				}
				_, _ = fmt.Fprintln(out)
			}

			// 6. Lint tool binaries
			seen := make(map[string]bool)
			var lintChecks []struct{ name, bin string }
			for _, h := range p.Hooks().Handlers() {
				linter, ok := h.(*handlers.Linter)
				if !ok {
					continue
				}
				for _, prof := range linter.Profiles() {
					for _, step := range prof.Steps {
						bin := step.Cmd[0]
						if seen[bin] {
							continue
						}
						seen[bin] = true
						lintChecks = append(lintChecks, struct{ name, bin string }{
							name: fmt.Sprintf("%s/%s", prof.Name, step.Label),
							bin:  bin,
						})
					}
				}
			}
			if len(lintChecks) > 0 {
				_, _ = fmt.Fprintln(out, "Lint Tools:")
				for _, lc := range lintChecks {
					if path, err := exec.LookPath(lc.bin); err == nil {
						_, _ = fmt.Fprintf(out, "  ✓ %-24s %s\n", lc.name, path)
					} else {
						_, _ = fmt.Fprintf(out, "  ✗ %-24s %s not found\n", lc.name, lc.bin)
						issues++
					}
				}
				_, _ = fmt.Fprintln(out)
			}

			// 7. Summary
			_, _ = fmt.Fprintln(out, "---")
			if issues == 0 {
				_, _ = fmt.Fprintln(out, "All checks passed!")
			} else {
				_, _ = fmt.Fprintf(out, "%d issue(s) found. Install missing tools to get full functionality.\n", issues)
			}
			return nil
		},
	}
}

// sessionCmd is the parent command for session management subcommands.
func sessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage Claude Code session state",
	}

	cmd.AddCommand(
		sessionListCmd(),
		sessionShowCmd(),
		sessionCleanCmd(),
	)

	return cmd
}

func sessionListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List recent sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := newHarness()
			sessions, err := p.SessionStore().List()
			if err != nil {
				return fmt.Errorf("list sessions: %w", err)
			}

			if len(sessions) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No sessions found.")
				return nil
			}

			// Sort by LastEventAt descending (most recent first).
			sort.Slice(sessions, func(i, j int) bool {
				ti := sessions[i].LastEventAt
				if ti.IsZero() {
					ti = sessions[i].StartedAt
				}
				tj := sessions[j].LastEventAt
				if tj.IsZero() {
					tj = sessions[j].StartedAt
				}
				return ti.After(tj)
			})

			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(out, "%-40s %-22s %-8s %-22s\n", "ID", "STARTED", "EVENTS", "LAST EVENT")
			_, _ = fmt.Fprintf(out, "%-40s %-22s %-8s %-22s\n", "---", "---", "---", "---")
			for _, sess := range sessions {
				started := sess.StartedAt.Format(time.RFC3339)
				lastEvent := "-"
				if !sess.LastEventAt.IsZero() {
					lastEvent = sess.LastEventAt.Format(time.RFC3339)
				}
				_, _ = fmt.Fprintf(out, "%-40s %-22s %-8d %-22s\n",
					sess.ID, started, sess.EventCount, lastEvent)
			}
			return nil
		},
	}
}

func sessionShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show full details for a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := newHarness()
			sess, err := p.SessionStore().Load(args[0])
			if err != nil {
				return fmt.Errorf("load session: %w", err)
			}

			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(out, "Session:      %s\n", sess.ID)
			_, _ = fmt.Fprintf(out, "Started:      %s\n", sess.StartedAt.Format(time.RFC3339))
			if !sess.LastEventAt.IsZero() {
				_, _ = fmt.Fprintf(out, "Last Event:   %s\n", sess.LastEventAt.Format(time.RFC3339))
			}
			_, _ = fmt.Fprintf(out, "Event Count:  %d\n", sess.EventCount)

			if len(sess.ToolCalls) > 0 {
				_, _ = fmt.Fprintf(out, "\nTool Calls (%d):\n", len(sess.ToolCalls))
				for _, tc := range sess.ToolCalls {
					blocked := ""
					if tc.Blocked {
						blocked = " [BLOCKED]"
					}
					_, _ = fmt.Fprintf(out, "  %s  %-20s %s%s\n",
						tc.Timestamp.Format(time.RFC3339), tc.ToolName, tc.Input, blocked)
				}
			}

			if len(sess.BlockedCalls) > 0 {
				_, _ = fmt.Fprintf(out, "\nBlocked Calls (%d):\n", len(sess.BlockedCalls))
				for _, bc := range sess.BlockedCalls {
					_, _ = fmt.Fprintf(out, "  %s  %-20s reason=%s\n",
						bc.Timestamp.Format(time.RFC3339), bc.ToolName, bc.Reason)
				}
			}

			if len(sess.FilesModified) > 0 {
				_, _ = fmt.Fprintf(out, "\nFiles Modified (%d):\n", len(sess.FilesModified))
				for _, f := range sess.FilesModified {
					_, _ = fmt.Fprintf(out, "  %s\n", f)
				}
			}

			if len(sess.Findings) > 0 {
				_, _ = fmt.Fprintf(out, "\nFindings (%d):\n", len(sess.Findings))
				for _, f := range sess.Findings {
					_, _ = fmt.Fprintf(out, "  [%s] %s %s:%d %s\n",
						f.Severity, f.Type, f.File, f.Line, f.Message)
				}
			}

			return nil
		},
	}
}

func sessionCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove sessions older than 7 days",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := newHarness()
			deleted, err := p.SessionStore().Cleanup(7 * 24 * time.Hour)
			if err != nil {
				return fmt.Errorf("cleanup sessions: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removed %d session(s) older than 7 days.\n", deleted)
			return nil
		},
	}
}
