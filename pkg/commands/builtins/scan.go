package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// ScanCommand scans the codebase for patterns and issues via /yaah:scan.
type ScanCommand struct{}

// NewScanCommand creates a new ScanCommand.
func NewScanCommand() *ScanCommand { return &ScanCommand{} }

func (c *ScanCommand) Name() string        { return "yaah/scan" }
func (c *ScanCommand) Description() string {
	return "Scan codebase for patterns, issues, and improvement opportunities"
}
func (c *ScanCommand) ArgumentHint() string { return "[--security | --quality | --deps | --all]" }
func (c *ScanCommand) Model() string        { return "" }
func (c *ScanCommand) AllowedTools() string { return "" }
func (c *ScanCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *ScanCommand) Content() string {
	return `# /yaah:scan — Codebase Scan

## When to use
When the user runs ` + "`/yaah:scan [--security | --quality | --deps | --all]`" + ` to audit the codebase.

## Behavior
- No flag or ` + "`--all`" + `: run all scans (security, quality, dependencies, patterns).
- ` + "`--security`" + `: focus only on security issues.
- ` + "`--quality`" + `: focus only on code quality issues.
- ` + "`--deps`" + `: focus only on dependency issues.

## Scans to perform

### Security
Look for:
- Hardcoded secrets, API keys, passwords, or tokens in source files
- SQL injection vulnerabilities (string concatenation in queries)
- XSS patterns (unescaped user input rendered as HTML)
- Unsafe deserialization
- Missing authentication or authorization checks on sensitive endpoints

### Quality
Look for:
- TODO/FIXME count and concentration by file
- Cyclomatic complexity hotspots (functions with deeply nested logic)
- Dead code (unexported/unreachable functions, unused variables)
- Duplicated logic across files
- Empty or swallowed error handlers (` + "`err != nil` ignored, `catch {}` empty`" + `)

### Dependencies
Look for:
- Outdated dependencies (compare declared versions against latest known)
- Known vulnerabilities (cross-reference CVE data where available)
- Unused dependencies declared in the manifest
- License compliance issues (GPL in a proprietary project, etc.)

### Patterns
Look for:
- Architectural patterns in use (registry, factory, observer, etc.)
- Inconsistencies in how patterns are applied across packages
- Anti-patterns (god objects, leaky abstractions, circular dependencies)

## Output format
Produce a scan report:

` + "```" + `
# Codebase Scan Report

## Security ({count} findings)
- [CRITICAL] <finding>: <file>:<line>
- [WARNING]  <finding>: <file>:<line>
- [INFO]     <finding>: <file>:<line>

## Quality ({count} findings)
- [WARNING]  <finding>: <file>:<line>
- [INFO]     <finding>: <file>:<line>

## Dependencies ({count} findings)
- [CRITICAL] <dep>: <issue>
- [WARNING]  <dep>: <issue>

## Patterns
- <pattern name>: <observation>
- <anti-pattern>: <location and recommendation>

## Summary
Total findings — critical: N, warning: N, info: N
` + "```" + `

## Rules
- Report only what is directly observable — never speculate
- Severity: CRITICAL = exploitable or data-loss risk; WARNING = likely problem; INFO = improvement opportunity
- Group findings by file when a single file has multiple issues
- Omit sections that have zero findings
`
}
