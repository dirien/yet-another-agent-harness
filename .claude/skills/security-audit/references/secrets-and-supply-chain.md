<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Secrets detection and supply chain security

## Secrets in current code

Scan for hardcoded credentials, API keys, tokens, and connection
strings in the current codebase.

**Automated scanning (prefer tools over grep):**

```bash
# gitleaks (recommended — low false positive rate)
gitleaks detect --no-git -v 2>&1

# trufflehog (alternative)
trufflehog filesystem . --no-update 2>&1
```

If no scanner is installed, use targeted grep patterns:

```bash
grep -rn --include="*.go" --include="*.js" --include="*.ts" \
  --include="*.py" --include="*.yaml" --include="*.yml" \
  --include="*.json" --include="*.env*" --include="*.toml" \
  --exclude-dir=node_modules --exclude-dir=vendor \
  --exclude-dir=.git \
  -E "(password|secret|api_key|apikey|private_key|token|credential)\
s*[=:]\s*['\"][^'\"]{8,}" . 2>/dev/null
```

**Reduce false positives:** Read each match. Exclude:

- Test fixtures with dummy values ("test-token", "fake-secret")
- Config templates with placeholder values ("CHANGE_ME", "xxx")
- Documentation examples
- Environment variable references (reading from env is fine)

## Secrets in git history

Current code may be clean, but secrets committed in past commits
persist in git history.

```bash
# gitleaks with git history
gitleaks detect -v 2>&1

# trufflehog with git history
trufflehog git file://. --no-update 2>&1
```

If a secret is found in history, flag it as CRITICAL — the secret
is compromised even if removed from current code. Remediation:
rotate the secret immediately.

## Supply chain security

### Dependency vulnerabilities

Run ecosystem-specific vulnerability scanners (see
`references/owasp-patterns.md`, A06 section).

### Lockfile integrity

| Ecosystem | Lockfile | Check |
| --------- | -------- | ----- |
| Node | `package-lock.json` / `yarn.lock` | Committed? Matches package.json? |
| Go | `go.sum` | Committed? `go mod verify` passes? |
| Python | `requirements.txt` / `poetry.lock` | Pinned versions? No ranges? |
| Rust | `Cargo.lock` | Committed? |

A missing or uncommitted lockfile means builds are not reproducible
and vulnerable to dependency confusion attacks.

### Dependency confusion

Check if the project uses:

- Private package registries (npmrc, pip.conf, GOPROXY)
- Namespace/scope prefixes (@org/package)
- `.npmrc` or `.pip.conf` with registry configuration

If private packages exist without namespace scoping, flag as HIGH
risk — an attacker can publish a public package with the same name
and higher version number.

### Typosquatting risk

For each dependency, check if the name is similar to a popular
package but slightly different (e.g., `lodas` vs `lodash`). Flag
any suspicious names.

### Maintenance signals

For critical dependencies, check:

- Last update date (> 2 years = flag as stale)
- Open security issues
- Maintainer count (bus factor)
- Download trends (declining = potential abandonment)
