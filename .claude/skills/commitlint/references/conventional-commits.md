<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Conventional Commits Reference

Full format specification for the
[Conventional Commits](https://www.conventionalcommits.org/) standard.

## Message structure

```text
<type>(<scope>): <subject>

<body>

<footer>
```

## Type (required)

Must be one of:

| Type | Description |
| ------ | ------------- |
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only changes |
| `style` | Code style changes (formatting, semicolons, etc.) |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | Performance improvement |
| `test` | Adding or updating tests |
| `build` | Changes to build system or dependencies |
| `ci` | Changes to CI configuration files and scripts |
| `chore` | Other changes that don't modify src or test files |
| `revert` | Reverts a previous commit |

## Scope (optional)

A noun describing the section of the codebase affected, in parentheses:

```text
feat(api): add login endpoint
fix(cli): resolve flag parsing error
```

- Must be lowercase
- Project-specific scopes are fine (e.g., `auth`, `ui`, `db`)

## Subject (required)

- 50 characters max (fits in GitHub PR title and `git log --oneline`)
- Start with lowercase
- No period at the end
- Use imperative mood ("add" not "added" or "adds")

## Body (optional)

- Separate from subject with a blank line
- Explain **what** and **why**, not how
- Wrap at 72 characters per line

## Footer (optional)

- Separate from body with a blank line
- Breaking changes: `BREAKING CHANGE: <description>`
- Issue references: `Fixes #123`, `Closes #456`
- Breaking change shorthand: append `!` after type/scope
  (e.g., `feat!: remove legacy API`)

## Default commitlint config

Use this `.commitlintrc.json` when no project config exists:

```json
{
  "extends": ["@commitlint/config-conventional"],
  "rules": {
    "type-enum": [
      2,
      "always",
      [
        "feat",
        "fix",
        "docs",
        "style",
        "refactor",
        "perf",
        "test",
        "build",
        "ci",
        "chore",
        "revert"
      ]
    ],
    "type-case": [2, "always", "lowercase"],
    "type-empty": [2, "never"],
    "scope-case": [2, "always", "lowercase"],
    "subject-case": [2, "always", "lowercase"],
    "subject-empty": [2, "never"],
    "subject-full-stop": [2, "never", "."],
    "header-max-length": [2, "always", 50],
    "body-leading-blank": [2, "always"],
    "body-max-line-length": [2, "always", 72],
    "footer-leading-blank": [2, "always"]
  }
}
```

## Valid examples

### Simple commit

```text
feat: add user authentication
```

### With scope

```text
fix(api): resolve null pointer in token validation
```

### With body and footer

```text
feat: add user authentication

Implement JWT-based authentication with refresh tokens.
Users can now log in with email/password and receive
access tokens that expire after 1 hour.

Fixes #123
```

### Breaking change

```text
feat!: remove legacy API endpoints

BREAKING CHANGE: The /v1/users endpoint has been removed.
Use /v2/users instead with the new authentication scheme.

Fixes #456
```

## Invalid examples

### Missing type

```text
add user authentication
```

**Error**: `type-empty` — subject may not be empty.
**Fix**: add a type prefix → `feat: add user authentication`

### Invalid type and subject too long

```text
feature: add comprehensive user authentication system with JWT tokens
```

**Errors**:

- `type-enum` — type `feature` is not allowed; use `feat`
- `header-max-length` — header is 69 characters (max 50)

**Fix**: `feat: add user authentication with JWT`
