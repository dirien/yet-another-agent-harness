<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Go Nolint Patterns

Scoring guidance for Phase 2 of the go-nolint-audit skill. Consult
when evaluating justification quality.

## Usually legitimate (score justification quality 1-2)

- `gochecknoglobals` on build-injected vars (`ldflags`) or immutable
  static lookup tables that never change after init
- `recvcheck` when an interface requires a specific receiver type
  (e.g., Bubble Tea requires value receivers for Init/Update/View)
- `gosec` when justification names the specific validation function
  or call site — not just "validated upstream"
- `mnd` where the number is self-documenting from immediate context:
  `wg.Add(N)` matching visible goroutines, buffer sizes with inline
  comment, ordering/precedence values
- `reassign` in test files overriding globals for test isolation
- `revive` on exported types where stuttering aids clarity, or for
  framework compatibility (e.g., Goa framework)
- `exhaustruct` on progressive initialization or optional fields at
  defaults — must specify which fields and why
- `errcheck` on best-effort drain (`resp.Body.Close`), panic recovery,
  non-critical logging — the category label must match the actual risk
- `testpackage` for white-box testing that needs access to unexported
  fields or functions
- `staticcheck` SA1019 when explicitly testing backward compatibility
  with deprecated APIs
- `protogetter` for nil vs zero-value disambiguation in protobuf
- `ireturn` / `wrapcheck` on decorator patterns returning interfaces
  for testability

## Challenge hard (score justification quality 4-5)

- `gocognit` / `funlen` with vague justifications ("complexity is
  inherent", "orchestration requires this") — these are the most
  common lazy suppressions
- `mnd` on business logic thresholds or configuration values that
  should be named constants
- `nilnil` without explaining the caller contract (what does nil,nil
  mean to the caller?)
- `errcheck` without explaining why the error is safe to ignore, or
  with category labels that don't match the risk ("best effort" on a
  database write is not the same as on a body close)
- `nestif` with "necessary nesting" — early returns almost always
  work; challenge the author to try
- `exhaustruct` with just "partial initialization" — needs to say
  what fields are left and why
- File-level `//nolint:` suppressing behavioral rules (`errcheck`,
  `gosec`) on production code. File-level on test files for
  structural rules (`testpackage`) is fine.
- Any nolint with no justification comment at all
- Bare `//nolint` (no rule specified) — suppresses everything, should
  always name the specific rule(s) being suppressed

## Red flags in justifications

These patterns in justification comments indicate lazy reasoning.
Score justification quality 4-5 when you see them:

- "inherent" / "necessary" / "required" without saying WHY
- "intentional" without explaining the INTENT
- Justification restates the rule name ("magic number is acceptable
  here")
- Justification references a removed or renamed function (stale
  comment, code evolved)
- "validated" without naming the specific validation function or
  call site
