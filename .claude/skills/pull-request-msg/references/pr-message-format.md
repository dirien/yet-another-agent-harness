<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# PR_MESSAGE.md format

```text
type(scope): brief description

## Summary

- Expand on the commit subject with additional context
- Do NOT simply restate the commit subject line
- Focus on user-visible functionality and purpose
- Wrap long lines at ~70 characters for readability

## Test plan

- [x] Completed validation steps
- [ ] Pending items if any

## Changes

### New files

- `path/to/file.go` - Brief description of purpose

### Modified files

- `path/to/existing.go` - What was changed

### Housekeeping (optional)

- Any cleanup, renames, or organizational changes

Closes #<issue-number>
```

## Commit types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Formatting, no code change
- `refactor` - Code restructuring without behavior change
- `perf` - Performance improvement
- `test` - Adding/updating tests
- `chore` - Maintenance tasks, dependencies
- `ci` - CI/CD changes
