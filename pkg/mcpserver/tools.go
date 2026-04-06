package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------- arg structs ----------

// ScanSecretsArgs holds the arguments for the yaah_scan_secrets tool.
type ScanSecretsArgs struct {
	FilePath string `json:"file_path" jsonschema:"absolute path to the file to scan"`
}

// LintArgs holds the arguments for the yaah_lint tool.
type LintArgs struct {
	FilePath string `json:"file_path" jsonschema:"absolute path to the file to lint"`
	Profile  string `json:"profile,omitempty" jsonschema:"lint profile name (e.g. golangci-lint, ruff, prettier). If omitted, selected by file extension."`
}

// CheckCommandArgs holds the arguments for the yaah_check_command tool.
type CheckCommandArgs struct {
	Command string `json:"command" jsonschema:"the shell command to check"`
}

// DoctorArgs holds the arguments for the yaah_doctor tool (no args).
type DoctorArgs struct{}

// SessionInfoArgs holds the arguments for the yaah_session_info tool.
type SessionInfoArgs struct {
	SessionID string `json:"session_id,omitempty" jsonschema:"session ID to look up. If omitted, returns basic server info."`
}

// PlanningStatusArgs holds the arguments for the yaah_planning_status tool.
type PlanningStatusArgs struct {
	ProjectDir string `json:"project_dir,omitempty" jsonschema:"project root directory, defaults to cwd"`
}

// PlanningInitArgs holds the arguments for the yaah_planning_init tool.
type PlanningInitArgs struct {
	ProjectDir string `json:"project_dir,omitempty" jsonschema:"project root directory, defaults to cwd"`
}

// ---------- tool registration ----------

func (s *Server) addScanSecretsTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_scan_secrets",
		Description: "Scan a file for hardcoded secrets and credentials",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args ScanSecretsArgs) (*mcp.CallToolResult, any, error) {
		scanner := findSecretScanner(s.harness.Hooks())
		if scanner == nil {
			scanner = handlers.NewSecretScanner()
		}

		findings, err := scanner.ScanFile(args.FilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("scan failed: %w", err)
		}

		data, err := json.Marshal(findings)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal findings: %w", err)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addLintTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_lint",
		Description: "Run lint checks on a file using configured profiles",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args LintArgs) (*mcp.CallToolResult, any, error) {
		linter := findLinter(s.harness.Hooks())
		if linter == nil {
			return nil, nil, fmt.Errorf("no linter handler registered")
		}

		cwd := filepath.Dir(args.FilePath)
		output, blocked, err := linter.LintFile(ctx, args.FilePath, args.Profile, cwd)
		if err != nil {
			return nil, nil, fmt.Errorf("lint failed: %w", err)
		}

		type lintResult struct {
			Output  string `json:"output"`
			Blocked bool   `json:"blocked"`
		}
		data, _ := json.Marshal(lintResult{Output: output, Blocked: blocked})
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addCheckCommandTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_check_command",
		Description: "Check whether a shell command is safe to run",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args CheckCommandArgs) (*mcp.CallToolResult, any, error) {
		guard := findCommandGuard(s.harness.Hooks())
		if guard == nil {
			guard = handlers.NewCommandGuard()
		}

		result := guard.CheckCommand(args.Command)
		data, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addDoctorTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_doctor",
		Description: "Run health checks on the yaah setup and report missing dependencies",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ DoctorArgs) (*mcp.CallToolResult, any, error) {
		type checkResult struct {
			Name     string `json:"name"`
			OK       bool   `json:"ok"`
			Detail   string `json:"detail,omitempty"`
			Category string `json:"category"`
		}

		var checks []checkResult

		// Binary checks.
		if path, err := exec.LookPath("yaah"); err == nil {
			checks = append(checks, checkResult{Name: "yaah", OK: true, Detail: path, Category: "binary"})
		} else {
			checks = append(checks, checkResult{Name: "yaah", OK: false, Detail: "not in PATH", Category: "binary"})
		}

		if path, err := exec.LookPath("git"); err == nil {
			checks = append(checks, checkResult{Name: "git", OK: true, Detail: path, Category: "binary"})
		} else {
			checks = append(checks, checkResult{Name: "git", OK: false, Detail: "not found", Category: "binary"})
		}

		// Settings validation.
		settingsPath := filepath.Join(".claude", "settings.json")
		if data, err := os.ReadFile(settingsPath); err == nil {
			if json.Valid(data) {
				checks = append(checks, checkResult{Name: settingsPath, OK: true, Detail: "valid JSON", Category: "config"})
			} else {
				checks = append(checks, checkResult{Name: settingsPath, OK: false, Detail: "invalid JSON", Category: "config"})
			}
		} else {
			checks = append(checks, checkResult{Name: settingsPath, OK: false, Detail: "not found", Category: "config"})
		}

		// LSP servers.
		lspResults := s.harness.LSP().CheckAll()
		for _, cr := range lspResults {
			checks = append(checks, checkResult{
				Name:     cr.Name,
				OK:       cr.Installed,
				Detail:   cr.BinaryPath,
				Category: "lsp",
			})
		}

		// MCP servers.
		for _, prov := range s.harness.MCP().Providers() {
			srv := prov.Server()
			if srv.Command == "" {
				checks = append(checks, checkResult{Name: prov.Name(), OK: true, Detail: srv.URL, Category: "mcp"})
				continue
			}
			if path, err := exec.LookPath(srv.Command); err == nil {
				checks = append(checks, checkResult{Name: prov.Name(), OK: true, Detail: path, Category: "mcp"})
			} else {
				checks = append(checks, checkResult{Name: prov.Name(), OK: false, Detail: srv.Command + " not found", Category: "mcp"})
			}
		}

		// Lint tools.
		seen := make(map[string]bool)
		for _, h := range s.harness.Hooks().Handlers() {
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
					label := fmt.Sprintf("%s/%s", prof.Name, step.Label)
					if path, err := exec.LookPath(bin); err == nil {
						checks = append(checks, checkResult{Name: label, OK: true, Detail: path, Category: "lint"})
					} else {
						checks = append(checks, checkResult{Name: label, OK: false, Detail: bin + " not found", Category: "lint"})
					}
				}
			}
		}

		data, _ := json.Marshal(checks)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addSessionInfoTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_session_info",
		Description: "Get information about a yaah session",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args SessionInfoArgs) (*mcp.CallToolResult, any, error) {
		store := s.harness.SessionStore()

		if args.SessionID != "" {
			sess, err := store.Load(args.SessionID)
			if err != nil {
				return nil, nil, fmt.Errorf("load session: %w", err)
			}
			data, _ := json.Marshal(sess)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
			}, nil, nil
		}

		// No session ID: return basic server info.
		type serverInfo struct {
			ServerName string `json:"server_name"`
			Version    string `json:"version"`
			StartTime  string `json:"start_time"`
			Handlers   int    `json:"handlers"`
			MCPCount   int    `json:"mcp_providers"`
			LSPCount   int    `json:"lsp_servers"`
			SkillCount int    `json:"skills"`
		}
		info := serverInfo{
			ServerName: "yaah",
			Version:    Version,
			StartTime:  time.Now().UTC().Format(time.RFC3339),
			Handlers:   len(s.harness.Hooks().Handlers()),
			MCPCount:   len(s.harness.MCP().Providers()),
			LSPCount:   len(s.harness.LSP().Providers()),
			SkillCount: len(s.harness.Skills().Skills()),
		}
		data, _ := json.Marshal(info)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addPlanningStatusTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_planning_status",
		Description: "Return the current planning status for a project's .planning/ directory",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args PlanningStatusArgs) (*mcp.CallToolResult, any, error) {
		dir := args.ProjectDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return nil, nil, fmt.Errorf("get cwd: %w", err)
			}
		}

		planningDir := filepath.Join(dir, ".planning")
		if _, err := os.Stat(planningDir); os.IsNotExist(err) {
			data, _ := json.Marshal(map[string]bool{"initialized": false})
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
			}, nil, nil
		}

		type planningStatus struct {
			Initialized    bool   `json:"initialized"`
			CurrentPhase   int    `json:"currentPhase"`
			Status         string `json:"status"`
			TotalPhases    int    `json:"totalPhases"`
			PlansCreated   int    `json:"plansCreated"`
			PlansExecuted  int    `json:"plansExecuted"`
			PhasesVerified int    `json:"phasesVerified"`
			QuickTaskCount int    `json:"quickTaskCount"`
			LastUpdated    string `json:"lastUpdated"`
		}

		result := planningStatus{Initialized: true}

		// Parse YAML frontmatter from STATE.md.
		stateFile := filepath.Join(planningDir, "STATE.md")
		if raw, err := os.ReadFile(stateFile); err == nil {
			content := string(raw)
			// Extract text between the first pair of --- markers.
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				for _, line := range strings.Split(parts[1], "\n") {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "phase:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "phase:"))
						var phase int
						if _, err := fmt.Sscanf(val, "%d", &phase); err == nil {
							result.CurrentPhase = phase
						}
					} else if strings.HasPrefix(line, "status:") {
						result.Status = strings.TrimSpace(strings.TrimPrefix(line, "status:"))
					} else if strings.HasPrefix(line, "last_updated:") {
						result.LastUpdated = strings.TrimSpace(strings.TrimPrefix(line, "last_updated:"))
					}
				}
			}
		}

		// Count phase subdirectories.
		phasesDir := filepath.Join(planningDir, "phases")
		if entries, err := os.ReadDir(phasesDir); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					result.TotalPhases++
				}
			}

			// Count plan, summary, and verification files across phase subdirs.
			if plans, err := filepath.Glob(filepath.Join(phasesDir, "*", "*-PLAN.md")); err == nil {
				result.PlansCreated = len(plans)
			}
			if summaries, err := filepath.Glob(filepath.Join(phasesDir, "*", "*-SUMMARY.md")); err == nil {
				result.PlansExecuted = len(summaries)
			}
			if verifications, err := filepath.Glob(filepath.Join(phasesDir, "*", "VERIFICATION.md")); err == nil {
				result.PhasesVerified = len(verifications)
			}
		}

		// Count quick task files.
		quickDir := filepath.Join(planningDir, "quick")
		if entries, err := os.ReadDir(quickDir); err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					result.QuickTaskCount++
				}
			}
		}

		data, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func (s *Server) addPlanningInitTool() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "yaah_planning_init",
		Description: "Initialize a .planning/ directory structure in the project root",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args PlanningInitArgs) (*mcp.CallToolResult, any, error) {
		dir := args.ProjectDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return nil, nil, fmt.Errorf("get cwd: %w", err)
			}
		}

		planningDir := filepath.Join(dir, ".planning")
		if _, err := os.Stat(planningDir); err == nil {
			return nil, nil, fmt.Errorf("planning directory already exists")
		}

		dirs := []string{
			planningDir,
			filepath.Join(planningDir, "phases"),
			filepath.Join(planningDir, "quick"),
			filepath.Join(planningDir, "notes"),
			filepath.Join(planningDir, "research"),
		}

		for _, d := range dirs {
			if err := os.MkdirAll(d, 0o755); err != nil {
				return nil, nil, fmt.Errorf("create directory %s: %w", d, err)
			}
		}

		now := time.Now().UTC().Format(time.RFC3339)
		stateContent := "---\n" +
			"milestone: v1.0\n" +
			"phase: 0\n" +
			"status: initialized\n" +
			"last_updated: " + now + "\n" +
			"---\n" +
			"# Current Position\n" +
			"Initialized. Run `/init` to set up project context, or `/plan 1` to start planning.\n" +
			"# Decisions Made\n" +
			"(none yet)\n" +
			"# Quick Tasks Completed\n" +
			"(none yet)\n"

		stateFile := filepath.Join(planningDir, "STATE.md")
		if err := os.WriteFile(stateFile, []byte(stateContent), 0o644); err != nil {
			return nil, nil, fmt.Errorf("write STATE.md: %w", err)
		}

		// Build relative paths for the response.
		created := make([]string, 0, len(dirs)+1)
		for _, d := range dirs {
			rel, err := filepath.Rel(dir, d)
			if err != nil {
				rel = d
			}
			created = append(created, rel)
		}

		type initResult struct {
			Created   []string `json:"created"`
			StateFile string   `json:"state_file"`
		}
		stateRel, err := filepath.Rel(dir, stateFile)
		if err != nil {
			stateRel = stateFile
		}
		data, _ := json.Marshal(initResult{
			Created:   created,
			StateFile: stateRel,
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

// ---------- helper finders ----------

// findSecretScanner searches the hook registry for a SecretScanner handler.
func findSecretScanner(reg *hooks.Registry) *handlers.SecretScanner {
	for _, h := range reg.Handlers() {
		if ss, ok := h.(*handlers.SecretScanner); ok {
			return ss
		}
	}
	return nil
}

// findLinter searches the hook registry for a Linter handler.
func findLinter(reg *hooks.Registry) *handlers.Linter {
	for _, h := range reg.Handlers() {
		if l, ok := h.(*handlers.Linter); ok {
			return l
		}
	}
	return nil
}

// findCommandGuard searches the hook registry for a CommandGuard handler.
func findCommandGuard(reg *hooks.Registry) *handlers.CommandGuard {
	for _, h := range reg.Handlers() {
		if g, ok := h.(*handlers.CommandGuard); ok {
			return g
		}
	}
	return nil
}
