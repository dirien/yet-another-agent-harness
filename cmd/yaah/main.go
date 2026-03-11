package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dirien/yet-another-agent-harness/pkg/generator"
	harnesspkg "github.com/dirien/yet-another-agent-harness/pkg/harness"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
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
