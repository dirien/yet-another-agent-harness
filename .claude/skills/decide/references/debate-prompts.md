<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Debate prompt templates

Templates for the three-agent adversarial debate. Adapt placeholders
(`[bracketed text]`) to the specific decision. Include all gathered
context in each agent prompt — agents run in isolation.

## Round 1: Position papers

### Advocate A

```text
You are a senior advisor. Your job is to STEELMAN Position A:
[position description].

CONTEXT:
[Insert all gathered context — project state, competitive landscape,
technical constraints, user requirements]

HARD CONSTRAINTS:
[Insert constraints]

RESEARCH REQUIREMENT:
Before writing your position, research externally to find at least 2
sources that support your position. Look for: real-world examples,
market data, pricing benchmarks, case studies, or documented outcomes.
Cite sources inline as [Source: URL or description].

YOUR TASK:
Write a comprehensive position paper (1500-2000 words) that:
1. Argues WHY Position A is superior
2. Details the technical/business architecture
3. Explains the specific implementation approach
4. Proposes a concrete plan with costs/timelines
5. Identifies the target audience and why they would choose this
6. Addresses risks and mitigations
7. Makes the business/technical case with specific numbers
8. Cites at least 2 external sources found via research

Be specific. Use real numbers, real costs, real timelines.
No hand-waving. Unsourced claims must be explicitly marked as
assumptions.

Title: "POSITION PAPER: [Position A Name]"
```

### Advocate B

Same structure as Advocate A, but steelmanning Position B. Include the
same research requirement — must find at least 2 external sources
supporting Position B. Mirror the same sections so papers are directly
comparable.

### Moderator

```text
You are a senior advisor acting as MODERATOR and DEVIL'S ADVOCATE.
Your job is to critically examine BOTH positions using external
evidence.

THE TWO POSITIONS:
Position A: [description]
Position B: [description]

CONTEXT:
[Same context as the advocates]

COUNTER-EVIDENCE RESEARCH MANDATE:
Before writing your analysis, research externally to actively seek
evidence that CHALLENGES each position. For each position, find at
least 2 sources that undermine, contradict, or complicate the
advocate's likely arguments. Look for: failed examples of similar
strategies, counter-data, hidden costs, documented pitfalls, or
competing explanations. Cite sources inline as
[Source: URL or description].

YOUR TASK:
Write a comprehensive analysis (2000-2500 words) that:
1. Lists the 10 hardest questions (5 per position) and answers each
2. Presents counter-evidence found via research for each position
3. Identifies the 5 biggest risks for each approach
4. Proposes a HYBRID model that takes the best of both
5. Recommends a specific plan with justification
6. Addresses at least 5 specific tough questions relevant to this
   decision
7. Reaches a preliminary recommendation given the constraints

All claims must cite sources. Unsourced claims must be marked as
assumptions.

Title: "MODERATOR ANALYSIS: [Decision Topic]"
```

## Round 1: Synthesis

After all three agents complete, build two tables:

**Agreements** — points where all three agents align:

| Topic | Consensus |
|-------|-----------|
| ...   | ...       |

**Disagreements** — points of tension:

| Topic | Advocate A | Advocate B | Moderator |
|-------|------------|------------|-----------|
| ...   | ...        | ...        | ...       |

Add brief commentary on the 2-3 most consequential tensions, then
proceed to Round 2.

## Round 2: Forced convergence

### Revised Advocate A

```text
You are Agent-1, the advocate for Position A. This is ROUND 2.
You must respond to the strongest challenges from Agent-2 and
Agent-3.

CHALLENGES FROM AGENT-2:
[Extract the 3-4 strongest arguments against Position A from
Advocate B's paper]

CHALLENGES FROM AGENT-3:
[Extract the 3-4 hardest questions and risks the Moderator
identified for Position A]

RESEARCH: If a challenge cites evidence you can verify or counter,
research externally to do so. Cite sources inline.

YOUR TASK (500-800 words):
1. CONCEDE the points that are genuinely correct (be honest)
2. DEFEND where your position is still strongest, citing evidence
3. Propose your REVISED position incorporating valid criticism
4. State specifically what you changed and why

The goal is convergence, not winning. If the other side is right
about something, say so.
```

### Revised Advocate B

Same structure as Revised Advocate A, reversed — responds to
Advocate A's strongest points and the Moderator's challenges to
Position B.

### Consensus Moderator

```text
You are the Moderator. This is ROUND 2. Based on the full debate,
drive toward FINAL CONSENSUS.

DEBATE SUMMARY:
[Full summary of Round 1 findings — agreements, disagreements,
key tensions]

RESEARCH: If any remaining tension can be resolved by verifying a
factual claim, look it up. Cite sources in the consensus document.

YOUR TASK — Write the FINAL CONSENSUS DOCUMENT (1000-1500 words):
1. State the agreed model — name it, describe it clearly
2. Define the exact structure (tiers, architecture, plan) with
   specifics
3. Resolve each remaining tension with a decision + 1-sentence
   justification
4. Specify the execution timeline
5. List 3 key metrics to track from day one
6. Name the single biggest risk and mitigation strategy
7. Provide a one-paragraph elevator pitch

Title: "CONSENSUS: [Decision Topic]"
```

## Consensus document format

The Moderator's Round 2 output should follow this structure:

```text
CONSENSUS: [Decision Topic]

## Agreed Model
[Name and clear description of the consensus approach]

## Structure
[Tiers, architecture, components, or plan details with specifics]

## Resolved Tensions
| Tension | Decision | Justification |
|---------|----------|---------------|
| ... | ... | ... |

## Execution Timeline
[Phased timeline with milestones]

## Key Metrics
1. [Metric] — [target and measurement method]
2. [Metric] — [target and measurement method]
3. [Metric] — [target and measurement method]

## Biggest Risk
[Risk description] — Mitigation: [strategy]

## Elevator Pitch
[One paragraph summarizing the consensus for a non-technical audience]
```

## Domain adaptation

Adapt the prompt specifics to the decision domain:

- **Pricing decisions**: include cost models, margins, competitor
  pricing, willingness-to-pay research, revenue projections
- **Architecture decisions**: include performance benchmarks, scaling
  characteristics, migration costs, team expertise, maintenance burden
- **Hiring decisions**: include role definitions, team composition,
  growth trajectory, compensation benchmarks, cultural fit criteria
- **Strategy decisions**: include market analysis, competitive
  positioning, resource allocation, timeline to impact, reversibility
