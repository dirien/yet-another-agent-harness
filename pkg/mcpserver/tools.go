package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
