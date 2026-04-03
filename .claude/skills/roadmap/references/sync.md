<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Sync mode — detailed workflow

Persona: Senior Product Manager (operational, pragmatic,
prioritization-focused).

## Step 1: Gather ground truth

```bash
gh issue list --state open --limit 50 \
  --json number,title,labels,milestone,assignees,body
gh issue list --state closed --limit 20 \
  --json number,title,closedAt
gh api repos/{owner}/{repo}/milestones \
  --jq '.[] | {title, open_issues, closed_issues, due_on, state}'
```

## Step 2: Read current state

Read ROADMAP.md to understand the documented state. Read CONTEXT.md
to understand project boundaries (all recommendations must respect
these boundaries).

## Step 3: Detect structure migration

Check if ROADMAP.md uses a legacy section layout. Legacy indicators:

- Section named "Current Focus" (should be "Immediate Focus")
- Section named "Past Milestones" or "Past Milestones (Done)"
- Separate "Completed Milestones" AND "Past Milestones" sections
- Section named "Backlog / Icebox" or "Strategic Research"
- Completed sections appear BEFORE the active focus section

If any legacy indicators are found, ask the user whether to migrate
to canonical structure or skip and only sync issue states.

**Migration behavior** (if approved):

- Rename "Current Focus" → "Immediate Focus"
- Merge "Past Milestones" and "Completed Milestones" into a single
  "Completed Milestones" section
- Organize completed items under quarterly sub-headings
  (`### YYYY-QN`), most recent first
- Merge "Backlog / Icebox" and "Strategic Research" into
  "Future Vision (Long-Term)"
- Reorder to canonical order:
  1. Vision / Overview
  2. Immediate Focus
  3. Near-Term Vision
  4. Future Vision (Long-Term)
  5. Completed Milestones (with quarterly sub-headings)
  6. Boundary Safeguards (if present)
- Preserve all issue references, checkboxes, and content

## Step 4: Detect discrepancies

- Issues marked done in ROADMAP.md but still OPEN in GitHub
- Issues marked open in ROADMAP.md but CLOSED in GitHub
- Issues missing from ROADMAP.md entirely
- Milestone misalignments

## Step 5: Update ROADMAP.md

Fix all discrepancies found. Move newly closed issues to
"Completed Milestones" under the appropriate quarterly sub-heading
(e.g., `### Q1 2026`). Determine quarter from the issue's
`closedAt` date. If the quarterly sub-heading does not exist,
create it at the top of "Completed Milestones" (most recent first).
Run markdownlint after edits.

## Step 6: Estimate level of effort

For each open issue without an `effort/*` label, analyze the issue
title, body, labels, and linked PRs. Use the effort estimation
heuristics in `references/label-config.md` to determine effort.
Apply the label:

```bash
gh issue edit {number} --add-label "effort/{size}"
```

## Step 7: Sync GitHub labels

Create any missing `roadmap/*` or `effort/*` labels (see
`references/label-config.md` for creation commands). Parse
ROADMAP.md sections to determine the correct phase label per issue.
For each open issue, remove incorrect `roadmap/*` labels and add
the correct one. For closed issues, remove all `roadmap/*` labels
but keep `effort/*` labels.

## Step 8: Update ROADMAP.md with LOE indicators

For each issue line in ROADMAP.md, append the LOE indicator if not
already present. Format: `- [ ] #42 Issue title [S]`. Update
existing indicators if the GitHub label differs.

## Step 9: Generate recommendations

Analyze open issues and recommend next 3-5 priorities based on:

- Dependencies (what unblocks other work?)
- Impact vs effort (quick wins first — prioritize Small + high impact)
- Theme coherence (related issues that could be batched)
- Milestone alignment
- LOE balance (mix of S/M/L to maintain velocity)

Include a sprint capacity summary:

```text
LOE Breakdown (open issues):
  Small:  X issues (~Xh estimated)
  Medium: X issues (~Xd estimated)
  Large:  X issues (~Xd estimated)
```

## Step 10: Present summary

Show to the user:

- Structure migration changes (if performed)
- Discrepancies fixed in ROADMAP.md
- Label changes made (added/removed)
- LOE labels applied (new estimates)
- Sprint capacity summary
- Prioritized recommendations with reasoning
- Ask which issue(s) to start
