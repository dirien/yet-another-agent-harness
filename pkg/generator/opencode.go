package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

// OpenCodeGenerator produces configuration for OpenCode.
type OpenCodeGenerator struct{}

func (g *OpenCodeGenerator) Target() schema.TargetAgent { return schema.TargetOpenCode }

// openCodeConfig is the output format for opencode.json.
type openCodeConfig struct {
	MCP map[string]openCodeMCPServer `json:"mcp,omitempty"`
}

type openCodeMCPServer struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

func (g *OpenCodeGenerator) GenerateSettings(cfg *schema.HarnessConfig) ([]byte, error) {
	out := openCodeConfig{}

	if cfg.MCP != nil && len(cfg.MCP.Servers) > 0 {
		out.MCP = make(map[string]openCodeMCPServer)
		for _, srv := range cfg.MCP.Servers {
			entry := openCodeMCPServer{
				Environment: srv.Env,
			}
			if srv.Command != "" {
				entry.Type = "local"
				entry.Command = append([]string{srv.Command}, srv.Args...)
			} else {
				entry.Type = "remote"
				entry.URL = srv.URL
			}
			out.MCP[srv.Name] = entry
		}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal opencode.json: %w", err)
	}
	return data, nil
}

func (g *OpenCodeGenerator) SettingsPath() string { return "opencode.json" }

// GenerateMCP returns nil — MCP is embedded in opencode.json.
func (g *OpenCodeGenerator) GenerateMCP(_ *schema.HarnessConfig) ([]byte, error) {
	return nil, nil
}

func (g *OpenCodeGenerator) MCPPath() string { return "" }

// GenerateHooks returns the JS plugin that bridges yaah hooks into OpenCode's plugin system.
func (g *OpenCodeGenerator) GenerateHooks(_ *schema.HarnessConfig) ([]byte, error) {
	plugin := `import { execFileSync } from "node:child_process";

// Map OpenCode camelCase args to yaah snake_case tool_input.
function toYaahInput(toolName, args) {
  const mapped = {};
  if (args) {
    for (const [k, v] of Object.entries(args)) {
      mapped[k.replace(/[A-Z]/g, (c) => "_" + c.toLowerCase())] = v;
    }
  }
  return { tool_name: toolName, tool_input: mapped };
}

// Cache tool args from before hooks so after hooks can access them.
const pendingCalls = new Map();

export const YaahPlugin = async (ctx) => {
  return {
    "tool.execute.before": async (input, output) => {
      pendingCalls.set(input.tool + ":" + Date.now(), output.args);
      try {
        execFileSync("yaah", ["hook", "PreToolUse"], {
          input: JSON.stringify(toYaahInput(input.tool, output.args)),
          cwd: ctx.directory,
          stdio: ["pipe", "pipe", "pipe"],
        });
      } catch (e) {
        if (e.status === 2) throw new Error(e.stderr?.toString() || "blocked by yaah");
      }
    },
    "tool.execute.after": async (input) => {
      // Clean up pending cache for this tool.
      for (const [key] of pendingCalls) {
        if (key.startsWith(input.tool + ":")) {
          pendingCalls.delete(key);
          break;
        }
      }
      // Wait for file to be flushed to disk before linting.
      await new Promise((r) => setTimeout(r, 150));
      try {
        execFileSync("yaah", ["hook", "PostToolUse"], {
          input: JSON.stringify(toYaahInput(input.tool, input.args || {})),
          cwd: ctx.directory,
          stdio: ["pipe", "pipe", "pipe"],
        });
      } catch (e) {
        throw new Error(e.stderr?.toString() || "blocked by yaah");
      }
    },
    event: async ({ event }) => {
      if (event.type === "session.created") {
        try {
          execFileSync("yaah", ["hook", "SessionStart"], {
            cwd: ctx.directory,
            stdio: ["pipe", "pipe", "pipe"],
          });
        } catch (e) {
          /* non-fatal */
        }
    }
  },
}};
`
	return []byte(plugin), nil
}

// openCodeAllTools lists every built-in OpenCode tool name.
var openCodeAllTools = []string{
	"bash", "edit", "glob", "grep", "list", "lsp", "patch",
	"question", "read", "skill", "todoread", "todowrite",
	"webfetch", "websearch", "write",
}

// claudeToOpenCodeTool maps Claude Code tool names (case-insensitive) to OpenCode equivalents.
var claudeToOpenCodeTool = map[string]string{
	"bash":      "bash",
	"bash(*)":   "bash",
	"edit":      "edit",
	"multiedit": "edit",
	"glob":      "glob",
	"grep":      "grep",
	"read":      "read",
	"write":     "write",
	"webfetch":  "webfetch",
	"websearch": "websearch",
	"todoread":  "todoread",
	"todowrite": "todowrite",
}

// FormatAgentTools converts a Claude-style comma-separated tool allowlist into
// OpenCode's frontmatter format: a YAML map that disables non-allowed tools.
func (g *OpenCodeGenerator) FormatAgentTools(tools string) string {
	allowed := make(map[string]bool)
	for _, t := range strings.Split(tools, ",") {
		name := strings.ToLower(strings.TrimSpace(t))
		if mapped, ok := claudeToOpenCodeTool[name]; ok {
			allowed[mapped] = true
		}
	}

	var b strings.Builder
	b.WriteString("tools:\n")
	for _, tool := range openCodeAllTools {
		if !allowed[tool] {
			_, _ = fmt.Fprintf(&b, "  %s: false\n", tool)
		}
	}
	return b.String()
}

func (g *OpenCodeGenerator) HooksPath() string    { return ".opencode/plugins/yaah.js" }
func (g *OpenCodeGenerator) SkillsDir() string    { return ".opencode/skills" }
func (g *OpenCodeGenerator) AgentsDir() string    { return ".opencode/agents" }
func (g *OpenCodeGenerator) AgentFileExt() string { return ".md" }
func (g *OpenCodeGenerator) CommandsDir() string  { return ".opencode/commands" }
