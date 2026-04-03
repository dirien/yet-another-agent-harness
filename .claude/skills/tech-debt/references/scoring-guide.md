<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Scoring guide

Framework for quantifying technical debt findings and prioritizing
remediation.

## Debt score (0-10 per category)

| Score | Level | Description |
| ----- | ----- | ----------- |
| 0-2 | Minimal | Well-maintained, minor improvements only |
| 3-5 | Manageable | Some debt, not blocking progress |
| 6-8 | Significant | Actively slowing development |
| 9-10 | Critical | Blocking progress or creating production risk |

## Impact levels

| Impact | Definition | Example |
| ------ | ---------- | ------- |
| LOW | Minor inconvenience | Inconsistent naming conventions |
| MEDIUM | Slows development | Missing tests require manual QA |
| HIGH | Blocks progress | Circular dependencies prevent new features |
| CRITICAL | Production risk | No error handling in payment flow |

## Effort levels

| Effort | Time | Example |
| ------ | ---- | ------- |
| S | < 1 day | Add missing test for a function |
| M | 1-5 days | Refactor a god object into 3 modules |
| L | 1-2 weeks | Add CI/CD pipeline from scratch |
| XL | > 2 weeks | Rewrite a subsystem to remove circular deps |

## Priority calculation

Priority = Impact / Effort. Higher impact with lower effort = higher
priority.

| Priority | Criteria | Action |
| -------- | -------- | ------ |
| P1 | HIGH/CRITICAL impact + S/M effort | Fix this sprint |
| P2 | HIGH impact + L effort, or MEDIUM + S/M | Fix next sprint |
| P3 | MEDIUM impact + L effort | Fix this quarter |
| P4 | LOW impact, any effort | Backlog |

## TECH_DEBT.md template

```markdown
# Technical Debt Report

Generated: YYYY-MM-DD

## Executive Summary

Overall health score: X/90 (sum of 9 category scores)
Categories assessed: X
P1 items requiring immediate attention: X

## Scorecard

| Category | Score | Impact | Top Issue |
| -------- | ----- | ------ | --------- |
| Code quality | X/10 | ... | ... |
| Linting | X/10 | ... | ... |
| Architecture | X/10 | ... | ... |
| Testing | X/10 | ... | ... |
| Documentation | X/10 | ... | ... |
| Infrastructure | X/10 | ... | ... |
| Security | X/10 | ... | ... |
| Dependencies | X/10 | ... | ... |
| Design principles | X/10 | ... | ... |

## Detailed Findings

### [Category Name] (Score: X/10)

#### Finding 1: [Title]
- **Impact**: [LOW/MEDIUM/HIGH/CRITICAL]
- **Effort**: [S/M/L/XL]
- **Priority**: [P1/P2/P3/P4]
- **Evidence**: [specific files, metrics, or patterns found]
- **Remediation**: [concrete action to take]

## Remediation Roadmap

### This Sprint (P1)
- [ ] [Action item with specific files/scope]

### Next Sprint (P2)
- [ ] [Action item]

### This Quarter (P3)
- [ ] [Action item]

### Backlog (P4)
- [ ] [Action item]

## Success Metrics
- [Metric]: [current value] → [target value]
```
