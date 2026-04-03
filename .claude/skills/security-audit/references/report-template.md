<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Security audit report template

Use this structure when generating SECURITY_AUDIT.md.

## Template

```markdown
# Security Audit Report

Generated: YYYY-MM-DD
Scope: [full codebase / specific module / recent changes]
Tools used: [list tools that ran successfully]
Tools unavailable: [list tools that were missing]

## Executive Summary

| Severity | Count |
| -------- | ----- |
| CRITICAL | X     |
| HIGH     | X     |
| MEDIUM   | X     |
| LOW      | X     |

Overall risk assessment: [brief 1-2 sentence summary]

## Findings

### [SEVERITY] Finding title

- **Category**: [OWASP A0X / Secrets / Supply Chain / Language-specific]
- **File**: `path/to/file.ext:line`
- **Evidence**: [what was found, why it is a vulnerability, and how
  it could be exploited]
- **Mitigations present**: [any existing protections, or "none"]
- **Remediation**: [specific action to fix, with code example if
  applicable]
- **Effort**: [S/M/L]

[Repeat for each finding, ordered by severity then effort]

## Threat Model Summary

### Entry Points Assessed

| Entry Point | Type | Auth Required | Threats Identified |
| ----------- | ---- | ------------- | ------------------ |
| ...         | ...  | ...           | ...                |

### Key STRIDE Findings

| Threat | Category | DREAD Score | Severity |
| ------ | -------- | ----------- | -------- |
| ...    | ...      | ...         | ...      |

## Supply Chain Assessment

| Check | Status | Details |
| ----- | ------ | ------- |
| Lockfile committed | ... | ... |
| Dependency vulnerabilities | ... | ... |
| Dependency confusion risk | ... | ... |
| Stale dependencies | ... | ... |

## Remediation Priority

### Immediate (CRITICAL)

- [ ] [Action with file:line and specific fix]

### This Sprint (HIGH)

- [ ] [Action with file:line and specific fix]

### Next Sprint (MEDIUM)

- [ ] [Action with file:line and specific fix]

### Backlog (LOW)

- [ ] [Action with file:line and specific fix]

## Recommendations

### Quick Wins (high impact, low effort)

[List items that can be fixed in < 1 day with significant
security improvement]

### Infrastructure Improvements

[Tooling, CI/CD, monitoring recommendations]

### Process Improvements

[Code review checklists, dependency update cadence, security
training]
```

## Severity assignment guide

| Severity | Criteria | Example |
| -------- | -------- | ------- |
| CRITICAL | Exploitable now, data breach or RCE risk | Hardcoded API key in git history |
| HIGH | Exploitable with moderate effort | SQL injection in admin endpoint |
| MEDIUM | Requires specific conditions | Missing rate limiting on login |
| LOW | Theoretical risk, minimal impact | Debug logging includes request IDs |

## Report quality checklist

Before finalizing the report, verify:

- Every finding has a specific file:line reference
- Every finding has concrete remediation (not just "fix this")
- False positives from automated tools have been filtered out
- Severity assignments are consistent across findings
- The remediation priority section has actionable items
- Quick wins are called out separately for immediate impact
