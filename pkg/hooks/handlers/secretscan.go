package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

var (
	_ hooks.Handler   = (*SecretScanner)(nil)
	_ hooks.FileAware = (*SecretScanner)(nil)
)

// SecretScanner is a PostToolUse handler that detects hardcoded secrets
// and credentials in edited files.
// Implements: Handler, FileAware.
type SecretScanner struct {
	patterns []*secretPattern
}

type secretPattern struct {
	re    *regexp.Regexp
	label string
}

// NewSecretScanner creates a secret scanner with default secret detection patterns.
func NewSecretScanner() *SecretScanner {
	v := &SecretScanner{}
	_ = v.addPattern(`(?i)(aws_access_key_id|aws_secret_access_key)\s*[=:]\s*["']?[A-Za-z0-9/+=]{20,}`, "AWS credential")
	_ = v.addPattern(`(?i)(AKIA|ABIA|ACCA|ASIA)[0-9A-Z]{16}`, "AWS access key ID")
	_ = v.addPattern(`(?i)ghp_[0-9a-zA-Z]{36}`, "GitHub personal access token")
	_ = v.addPattern(`(?i)gho_[0-9a-zA-Z]{36}`, "GitHub OAuth token")
	_ = v.addPattern(`(?i)github_pat_[0-9a-zA-Z_]{82}`, "GitHub fine-grained PAT")
	_ = v.addPattern(`sk-[0-9a-zA-Z]{20}T3BlbkFJ[0-9a-zA-Z]{20}`, "OpenAI API key")
	_ = v.addPattern(`sk-ant-[0-9a-zA-Z-]{90,}`, "Anthropic API key")
	_ = v.addPattern(`(?i)(password|passwd|secret)\s*[:=]{1,2}\s*["'][^"']{8,}["']`, "hardcoded password/secret")
	_ = v.addPattern(`(?i)bearer\s+[a-zA-Z0-9\-._~+/]+=*`, "bearer token")
	_ = v.addPattern(`-----BEGIN (RSA |EC |DSA )?PRIVATE KEY-----`, "private key")
	_ = v.addPattern(`(?i)(slack|xoxb|xoxp|xapp|xoxa)-[0-9a-zA-Z-]{10,}`, "Slack token")
	_ = v.addPattern(`(?i)SG\.[0-9a-zA-Z_-]{22}\.[0-9a-zA-Z_-]{43}`, "SendGrid API key")
	_ = v.addPattern(`(?i)hooks\.slack\.com/services/T[A-Z0-9]+/B[A-Z0-9]+/[a-zA-Z0-9]+`, "Slack webhook URL")
	return v
}

// AddPattern adds a custom secret detection regex.
func (v *SecretScanner) AddPattern(pattern, label string) error {
	return v.addPattern(pattern, label)
}

func (v *SecretScanner) addPattern(pattern, label string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile pattern %q: %w", label, err)
	}
	v.patterns = append(v.patterns, &secretPattern{re: re, label: label})
	return nil
}

func (v *SecretScanner) Name() string { return "secret-scanner" }

func (v *SecretScanner) Events() []schema.HookEvent {
	return []schema.HookEvent{schema.HookPostToolUse}
}

var editWriteMatch = regexp.MustCompile(`(?i)^(Edit|Write|MultiEdit)$`)

func (v *SecretScanner) Match() *regexp.Regexp {
	return editWriteMatch
}

func (v *SecretScanner) FileExtensions() []string {
	return nil // scan all text files
}

// SecretFinding represents a single secret detected in a file.
type SecretFinding struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Pattern  string `json:"pattern"`
}

// ScanFile scans a file for secrets and returns any findings.
// Returns nil if the file is binary or cannot be opened.
func (v *SecretScanner) ScanFile(filePath string) ([]SecretFinding, error) {
	if isLikelyBinary(filePath) {
		return nil, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var findings []SecretFinding
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, p := range v.patterns {
			if p.re.MatchString(line) {
				findings = append(findings, SecretFinding{
					FilePath: filePath,
					Line:     lineNum,
					Pattern:  p.label,
				})
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return findings, fmt.Errorf("scan file %s: %w", filePath, err)
	}

	return findings, nil
}

func (v *SecretScanner) Execute(_ context.Context, input *hooks.Input) (*hooks.Result, error) {
	fp := input.FilePath()
	if fp == "" {
		return nil, nil
	}

	findings, err := v.ScanFile(fp)
	if err != nil {
		return &hooks.Result{
			Error: fmt.Sprintf("secret-scanner: %v", err),
		}, nil
	}

	if len(findings) > 0 {
		lines := make([]string, len(findings))
		for i, f := range findings {
			lines[i] = fmt.Sprintf("  %s:%d: possible %s", f.FilePath, f.Line, f.Pattern)
		}
		return &hooks.Result{
			Error: fmt.Sprintf("secret-scanner found %d potential secret(s):\n%s\nDo NOT commit these. Use environment variables or a secrets manager instead.", len(findings), strings.Join(lines, "\n")),
			Block: true,
		}, nil
	}
	return nil, nil
}

// binaryExts lists file extensions that should be skipped by the secret scanner.
var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
	".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
	".zip": true, ".tar": true, ".gz": true, ".bz2": true,
	".pdf": true, ".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".lock": true,
}

func isLikelyBinary(path string) bool {
	return binaryExts[filepath.Ext(path)]
}
