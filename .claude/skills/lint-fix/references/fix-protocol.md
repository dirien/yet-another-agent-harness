<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Atomic fix protocol

Fix linting and validation issues one at a time to prevent cascade
failures. Each fix is immediately verified before moving on.

## Core rule

**Fix ONE issue, then re-verify.** Do not batch fixes. A fix that
resolves one issue may introduce new ones — catching this immediately
is cheaper than debugging a chain of cascading changes.

## Protocol

For each issue reported by a failing command:

1. **Mark the issue** — note the file, line, and rule/error
2. **Fix exactly one issue** — make the minimal change needed
3. **Re-run the failing command** — not the full pipeline, just the
   command that reported this issue
4. **Verify two things**:
   - The specific issue is resolved
   - No NEW issues were introduced by the fix
5. **If the fix created new issues**: revert the fix and try a
   different approach
6. **If the fix is clean**: move to the next issue
7. **After all issues for this command are resolved**: move to the
   next command in the pipeline

## When a fix creates more issues

```bash
# See what changed
git diff

# Revert a specific file
git checkout -- path/to/file

# If multiple files need reverting
git stash
# Try a different approach, then:
git stash pop
```

Try alternative approaches in order:

1. Different fix for the same issue
2. Fix the new issues too (if they are genuine improvements)
3. Suppress the rule for this line with a justification comment
   (last resort)

## When to stop fixing

- **Zero errors**: all commands in the pipeline pass — success
- **Circular fixes**: fix A introduces B, fix B reintroduces A —
  ask the user which to prioritize
- **Non-fixable issues**: the issue requires an architectural change
  or upstream fix — document it and skip with justification
- **False positives**: the linter is wrong — suppress with a
  justification comment

## Re-verification cadence

- After each individual fix: re-run the failing command only
- After all issues in one command are fixed: re-run that command
  once more (full pass)
- After all commands pass individually: run the full pipeline
  end-to-end as final verification
