<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Generate mode — detailed workflow

Persona: Technical Architect + Product Manager hybrid (foundational,
systematic, boundary-setting).

## Pre-flight check

If CONTEXT.md or ROADMAP.md already exists, ask the user:

- **Regenerate**: Analyze fresh and overwrite existing files
- **Merge**: Preserve existing content and add missing elements
- **Abort**: Cancel and keep existing files

## CONTEXT.md generation

### Step 1: Deep codebase analysis

Investigate the codebase to understand:

- Project structure and architecture
- Core dependencies and their purposes
- Entry points and main workflows
- What the tool does (capabilities)

### Step 2: Infer boundaries

Determine what the project explicitly does NOT do:

- What external systems does it avoid calling directly?
- What state does it avoid persisting?
- What operations does it delegate vs implement?
- What domains are out of scope?

### Step 3: Research similar tools

Research externally to find similar tools and understand common
feature expectations, where this tool differs, and industry
terminology.

### Step 4: Draft CONTEXT.md

Use this template:

```markdown
# {Project} Context & Boundaries

## Core Architectural Identity

{One-paragraph description of what the project IS}

## Technical Boundaries ("Hard No's")

{Bulleted list of what the project does NOT do}

## Data Source of Truth

{Where does authoritative data come from?}

## Interaction Model

{Inbound and outbound interfaces}

## Verification

{How to check if a feature violates boundaries}
```

## ROADMAP.md generation

### Step 1: Fetch all GitHub issues

```bash
gh issue list --state all --limit 100 \
  --json number,title,body,labels,milestone,state
gh api repos/{owner}/{repo}/milestones
```

### Step 2: Categorize issues

Group by:

- Immediate focus (current milestone or high priority)
- Near-term (next milestone)
- Future vision (backlog, research, unassigned or low priority)
- Completed (closed milestones or done issues, grouped by quarter)

For each open issue, estimate LOE using the heuristics in
`references/label-config.md` and apply `effort/*` labels.

### Step 3: Identify themes

Look for patterns: feature areas, issue types (bug, enhancement,
docs), and cross-cutting concerns.

### Step 4: Draft ROADMAP.md

Use this template:

```markdown
# {Project} Strategic Roadmap

## Vision

{Brief description and link to CONTEXT.md}

## Immediate Focus (vX.Y.Z)

- [ ] #1 Issue title [S]
- [ ] #2 Issue title [M]

## Near-Term Vision (vX.Y+1.0)

- [ ] #3 Issue title [L]

## Future Vision (Long-Term)

- [ ] #4 Issue title [M]

## Completed Milestones

### YYYY-QN

- [x] #5 Issue title [S]

## Boundary Safeguards

{Hard no's from CONTEXT.md}
```

### Step 5: Write files

Write CONTEXT.md first (it constrains ROADMAP.md). Write ROADMAP.md
second. Run markdownlint on both files after creation.
