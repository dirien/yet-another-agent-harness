<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Priority matrix

How to categorize and prioritize dependency upgrades.

## Priority levels

| Priority | Criteria | Action |
| -------- | -------- | ------ |
| CRITICAL | Known CVE, security advisory, actively exploited | Upgrade immediately |
| HIGH | Major version behind, dependency EOL, high-severity advisory | Upgrade this session |
| MEDIUM | Minor version behind, feature improvements available | Upgrade if time permits |
| LOW | Patch version behind, docs/minor bug fixes only | Batch with other work |

## Signals for each level

**CRITICAL signals:**

- govulncheck, npm audit, pip-audit, or cargo audit flags a
  vulnerability
- CVE exists with CVSS score >= 7.0
- Dependency maintainer issued a security advisory
- Known exploit in the wild

**HIGH signals:**

- Dependency is 2+ major versions behind current
- Maintainer announced end-of-life or deprecated the version
- Significant performance or stability fixes in newer version
- Breaking changes accumulate (harder to upgrade later)

**MEDIUM signals:**

- Minor version features or performance improvements available
- New APIs that would simplify existing code
- Deprecation warnings for APIs currently in use

**LOW signals:**

- Only patch-level changes (typo fixes, doc updates, minor bugs)
- Changes do not affect features used by the project
- Cosmetic or internal refactoring in the dependency

## Breaking change assessment

Before upgrading across a major version boundary, check:

1. **CHANGELOG or release notes** — search for "BREAKING",
   "migration", "upgrade guide"
2. **Major version boundary** — v1 → v2, v2 → v3, etc.
3. **Go module path changes** — module paths include major version
   (e.g., `github.com/foo/bar/v2` → `github.com/foo/bar/v3`)
4. **Node peer dependency conflicts** — major upgrades may break
   peer dependency requirements
5. **Python deprecation removals** — features deprecated in minor
   versions are often removed in major versions
6. **API surface changes** — renamed functions, changed signatures,
   removed exports

## Upgrade ordering

When multiple upgrades are needed:

1. Security fixes first (CRITICAL, then HIGH)
2. Dependencies of dependencies (transitive deps that unblock
   others)
3. Smallest changes first within the same priority (patch before
   minor before major)
4. Most-used dependency last (highest risk of breakage, benefits
   from other upgrades being stable first)
