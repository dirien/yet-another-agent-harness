<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Upgrade report template

Use this structure when presenting the upgrade summary to the user.

## Pre-upgrade plan (presented for approval)

```text
Dependency Upgrade Plan

Ecosystem: [Go / Node.js / Python / Rust / .NET / Mixed]
Outdated packages: X
Vulnerable packages: X

| # | Package | Current | Target | Priority | Breaking? |
|---|---------|---------|--------|----------|-----------|
| 1 | ...     | ...     | ...    | CRITICAL | No        |
| 2 | ...     | ...     | ...    | HIGH     | Yes       |
| 3 | ...     | ...     | ...    | MEDIUM   | No        |

Proceed with upgrades? (All / Select by number / Skip)
```

## Post-upgrade summary

```markdown
# Dependency Upgrade Report

Generated: YYYY-MM-DD
Ecosystem: [Go / Node.js / Python / Rust / .NET / Mixed]

## Results

### Upgraded Successfully

| Package | From | To | Tests |
| ------- | ---- | -- | ----- |
| ...     | ...  | ...| PASS  |

### Failed (Rolled Back)

| Package | From | Target | Failure Reason |
| ------- | ---- | ------ | -------------- |
| ...     | ...  | ...    | Test failure in X |

### Skipped

| Package | From | Target | Reason |
| ------- | ---- | ------ | ------ |
| ...     | ...  | ...    | Breaking change, needs migration |

## Remaining Vulnerabilities

| Package | CVE | Severity | Why Not Fixed |
| ------- | --- | -------- | ------------- |
| ...     | ... | ...      | Requires major version upgrade |

## Lockfile Status

| File | Updated | Verified |
| ---- | ------- | -------- |
| ...  | Yes/No  | Yes/No   |

## Next Steps

- [ ] [Any manual migration steps needed for skipped upgrades]
- [ ] [Remaining vulnerability remediation]
```
