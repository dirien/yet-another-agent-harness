<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Label configuration

Shared label definitions and heuristics used by the sync and generate
modes. Create missing labels before syncing.

## Phase labels (mutually exclusive)

An issue should have exactly ONE phase label.

| Label | Color | Description | ROADMAP.md section |
| ----- | ----- | ----------- | ------------------ |
| `roadmap/current` | `#0E8A16` (green) | Active development | Immediate Focus |
| `roadmap/next` | `#1D76DB` (blue) | Next milestone | Near-Term Vision |
| `roadmap/future` | `#5319E7` (purple) | Later milestones | Future Vision |

## Level of effort labels (mutually exclusive)

An issue should have at most ONE effort label.

| Label | Color | Description | Guideline |
| ----- | ----- | ----------- | --------- |
| `effort/small` | `#c5def5` (light blue) | Small effort | 1-2 hours, straightforward |
| `effort/medium` | `#bfd4f2` (blue) | Medium effort | Half day to 1 day |
| `effort/large` | `#0075ca` (dark blue) | Large effort | Multi-day, complex |

## Effort estimation heuristics

When auto-estimating LOE for issues without an `effort/*` label:

| Signal | Points toward |
| ------ | ------------- |
| Single file change, clear scope | Small |
| Bug fix with known root cause | Small |
| New feature touching 2-4 files | Medium |
| Requires new tests + implementation | Medium |
| API/schema changes | Medium-Large |
| Cross-cutting refactor | Large |
| New subsystem or major feature | Large |
| Requires research or spike first | Large |
| Issue body is vague or underspecified | Medium (default) |

## Contribution labels (additive)

These can be added alongside phase labels.

| Label | Color | Description |
| ----- | ----- | ----------- |
| `community` | `#7057FF` (purple) | Suitable for external contributors |
| `cross-repo` | `#FBCA04` (yellow) | Requires changes in multiple repos |
| `spec-first` | `#D93F0B` (orange) | Blocked on spec/upstream changes |

## LOE indicators in ROADMAP.md

Append effort indicators to issue lines in ROADMAP.md:

- Format: `- [ ] #42 Issue title [S]`
- Key: `[S]` = Small, `[M]` = Medium, `[L]` = Large
- Place after the issue title, before any trailing notes
- If an issue already has an indicator, update it if the GitHub
  label differs

## Label sync behavior

During sync mode:

1. Parse ROADMAP.md to build issue-to-section mapping:
   - "Immediate Focus" → `roadmap/current`
   - "Near-Term Vision" → `roadmap/next`
   - "Future Vision" → `roadmap/future`
2. For each open issue in ROADMAP.md:
   - Fetch current labels from GitHub
   - Remove any incorrect `roadmap/*` labels
   - Add the correct `roadmap/*` label if missing
3. For each closed issue moved to "Completed Milestones":
   - Remove all `roadmap/*` labels
   - Keep `effort/*` labels for historical tracking
4. Report all label changes made

## Label creation commands

Run once if `roadmap/*` labels do not exist:

```bash
# Phase labels
gh label create "roadmap/current" --color "0E8A16" \
  --description "Active development - Immediate Focus"
gh label create "roadmap/next" --color "1D76DB" \
  --description "Next milestone - Near-Term Vision"
gh label create "roadmap/future" --color "5319E7" \
  --description "Future milestones - Future Vision"

# Effort labels
gh label create "effort/small" --color "c5def5" \
  --description "Small effort - 1-2 hours"
gh label create "effort/medium" --color "bfd4f2" \
  --description "Medium effort - half day to 1 day"
gh label create "effort/large" --color "0075ca" \
  --description "Large effort - multi-day"

# Contribution labels (optional)
gh label create "community" --color "7057FF" \
  --description "Suitable for external contributors"
gh label create "cross-repo" --color "FBCA04" \
  --description "Requires changes in multiple repos"
gh label create "spec-first" --color "D93F0B" \
  --description "Blocked on spec/upstream changes"
```
