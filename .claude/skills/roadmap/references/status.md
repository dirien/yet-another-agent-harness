<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Status mode — detailed workflow

Persona: Analyst (factual, concise, no recommendations).

## Constraint

This mode is **read-only**. Do not modify any files.

## Step 1: Gather data

```bash
gh issue list --state open \
  --json number,title,labels,milestone | jq length
gh api repos/{owner}/{repo}/milestones \
  --jq '.[] | {title, open_issues, closed_issues}'
```

Read ROADMAP.md to understand the documented state.

## Step 2: Present summary tables

```text
## Roadmap Status

| Metric              | Value |
|---------------------|-------|
| Open Issues         | X     |
| Current Milestone   | vX.Y  |
| Milestone Progress  | X/Y   |
| Issues in ROADMAP   | X     |

### By Category

| Category      | Open | Done |
|---------------|------|------|
| Enhancement   | X    | X    |
| Bug           | X    | X    |
| Documentation | X    | X    |

### By Level of Effort

| Effort    | Open | Done | Unestimated |
|-----------|------|------|-------------|
| Small     | X    | X    |             |
| Medium    | X    | X    |             |
| Large     | X    | X    |             |
| **Total** | X    | X    | X           |
```

## Step 3: Flag obvious discrepancies

Only flag clearly visible issues without deep analysis:

- "Note: X issues in ROADMAP.md marked done are still open in GitHub"
- "Note: X open issues have no effort estimate"
- "Note: X issues in GitHub are not referenced in ROADMAP.md"
