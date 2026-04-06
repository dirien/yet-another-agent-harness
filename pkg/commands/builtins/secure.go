package builtins

import "github.com/dirien/yet-another-agent-harness/pkg/commands"

// SecureCommand runs security analysis for a phase via /yaah:secure.
type SecureCommand struct{}

// NewSecureCommand creates a new SecureCommand.
func NewSecureCommand() *SecureCommand { return &SecureCommand{} }

func (c *SecureCommand) Name() string { return "yaah/secure" }
func (c *SecureCommand) Description() string {
	return "Security threat modeling and vulnerability analysis for a phase"
}
func (c *SecureCommand) ArgumentHint() string { return "<phase-number>" }
func (c *SecureCommand) Model() string        { return "" }
func (c *SecureCommand) AllowedTools() string { return "" }
func (c *SecureCommand) Advanced() commands.CommandAdvanced {
	return commands.CommandAdvanced{Context: "fork"}
}

func (c *SecureCommand) Content() string {
	return `# /yaah:secure — Security Threat Modeling and Vulnerability Analysis

## When to use
When the user runs ` + "`/yaah:secure <phase>`" + ` to perform security analysis on the code produced in a planning phase.

## Steps

### 1. Load phase context
- Read ` + "`.planning/phases/{NN}-{slug}/`" + ` to locate all PLAN.md files for the phase
- Parse each PLAN.md frontmatter to identify ` + "`files_modified`" + ` and ` + "`files_created`"  + `
- Read every source file listed across all plans

### 2. Perform STRIDE threat analysis
Evaluate each component of the code against the STRIDE model:

- **Spoofing**: Are there authentication gaps? Can an actor impersonate another?
- **Tampering**: Is user input validated? Is data integrity enforced at boundaries?
- **Repudiation**: Is audit logging present for sensitive operations?
- **Information Disclosure**: Are secrets, tokens, or PII exposed in logs, errors, or responses?
- **Denial of Service**: Are there rate limits, timeouts, or resource caps in place?
- **Elevation of Privilege**: Are authorization checks present and consistent?

### 3. Check OWASP Top 10
Scan the source files for evidence of:
- A01 Broken Access Control — missing or inconsistent authorization
- A02 Cryptographic Failures — weak algorithms, plaintext secrets, missing TLS
- A03 Injection — SQL, command, LDAP, or template injection vectors
- A04 Insecure Design — missing threat model, unsafe defaults
- A05 Security Misconfiguration — debug flags, open CORS, default credentials
- A06 Vulnerable Components — outdated or known-vulnerable dependencies
- A07 Auth and Session Management Failures — weak session tokens, missing expiry
- A08 Software and Data Integrity Failures — unsigned artifacts, unsafe deserialization
- A09 Logging and Monitoring Failures — missing alerts, insufficient audit trail
- A10 Server-Side Request Forgery — unvalidated URL inputs reaching internal services

### 4. Targeted scans
- Scan for hardcoded secrets: API keys, passwords, tokens (look for ` + "`=` patterns near `key`, `token`, `secret`, `password`" + `)
- SQL injection: string concatenation into query variables
- XSS: unescaped user input rendered to HTML
- CSRF: state-changing endpoints missing token validation

### 5. Write security report
Write threat model to ` + "`.planning/phases/{NN}-{slug}/SECURITY.md`" + `:

` + "```" + `markdown
# Security Analysis — Phase {N}

## Threat Model (STRIDE)
| Threat | Status | Mitigation |
|--------|--------|------------|
| Spoofing | {OK / AT RISK / UNMITIGATED} | {description or "none"} |
| Tampering | {OK / AT RISK / UNMITIGATED} | {description or "none"} |
| Repudiation | {OK / AT RISK / UNMITIGATED} | {description or "none"} |
| Information Disclosure | {OK / AT RISK / UNMITIGATED} | {description or "none"} |
| Denial of Service | {OK / AT RISK / UNMITIGATED} | {description or "none"} |
| Elevation of Privilege | {OK / AT RISK / UNMITIGATED} | {description or "none"} |

## Vulnerabilities Found
- CRITICAL: {file:line} — {description}
- HIGH: {file:line} — {description}
- MEDIUM: {file:line} — {description}
- LOW: {file:line} — {description}

(List "None found" per severity level if clean.)

## OWASP Top 10 Coverage
| Category | Status | Notes |
|----------|--------|-------|
...

## Recommendations
1. {Concrete remediation step — include file and function if applicable}
...
` + "```" + `

### 6. Commit
Run: ` + "`git add .planning/ && git commit -m \"docs(planning): security analysis phase {N}\"`" + `

## Rules
- NEVER mark a threat as OK without evidence in the source files
- Tag every finding with file path and line number
- If no vulnerabilities are found in a category, explicitly state "None found"
- Do not modify source files — only write the SECURITY.md report
`
}
