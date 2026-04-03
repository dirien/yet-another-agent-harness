<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Markdownlint Rules Reference

Reference for markdownlint rules commonly encountered during validation.
The default config for this skill is minimal — it only overrides line length
to 120 characters. All other rules use markdownlint's built-in defaults.
Projects should customize with their own `.markdownlint.json` as needed.

## Commonly triggered rules

| Rule | Name | Description | Why It Matters |
| ---- | ---- | ----------- | -------------- |
| MD001 | heading-increment | Heading levels increment by one | Proper document structure |
| MD012 | no-multiple-blanks | No multiple consecutive blank lines | Consistent spacing |
| MD013 | line-length | Line length (default config: 120 chars) | Readable lines |
| MD022 | blanks-around-headings | Blank lines around headings | Better readability |
| MD025 | single-title | Only one top-level heading (h1) | Clear document title |
| MD031 | blanks-around-fences | Blank lines around code blocks | Visual separation |
| MD032 | blanks-around-lists | Blank lines around lists | Clear boundaries |
| MD040 | fenced-code-language | Fenced code blocks should have language | Syntax highlighting |
| MD047 | single-trailing-newline | Files end with single newline | POSIX standard |

For the full list of rules, see the
[markdownlint rules documentation](https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md).

## Auto-fixable rules

Running `markdownlint --fix` automatically resolves these rules:

MD007, MD009, MD010, MD011, MD012, MD022, MD023, MD026, MD027, MD030,
MD031, MD032, MD037, MD038, MD044, MD047, MD049, MD050

Structural issues like MD001 (heading-increment) and MD025 (single-title)
require manual fixes because the correct resolution depends on the intended
document structure.

## Common fixes

**MD001 — heading-increment:**

```text
Wrong:
# Title
### Subsection (skipped h2)

Fixed:
# Title
## Section
### Subsection
```

**MD022 — blanks-around-headings:**

```text
Wrong:
Some text
## Heading
More text

Fixed:
Some text

## Heading

More text
```

**MD032 — blanks-around-lists:**

```text
Wrong:
Some text
- item 1
- item 2
More text

Fixed:
Some text

- item 1
- item 2

More text
```

**MD047 — single-trailing-newline:**

```text
Fix: ensure file ends with exactly one newline character (no blank
lines at end, no missing newline)
```

## Default config

This is the same content as `markdownlint-config.json` in this directory.
The config is intentionally minimal — only line length is overridden.
Projects should add their own `.markdownlint.json` for style preferences
(heading style, list markers, emphasis style, etc.).

```json
{
  "MD013": { "line_length": 120 }
}
```
