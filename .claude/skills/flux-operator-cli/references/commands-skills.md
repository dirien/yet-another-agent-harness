# Skills Management Commands

The Flux Operator CLI includes commands for managing AI agent skills distributed
as OCI artifacts.

---

## skills install

Install skills from an OCI artifact repository.

```bash
flux-operator skills install <repository> [flags]
```

Skills are installed to the `.agents/skills` directory.

| Flag | Default | Description |
|------|---------|-------------|
| `--tag` | `latest` | OCI artifact tag to install |
| `--skill` | — | Install specific skill(s) only (repeatable) |
| `--agent` | — | Install specific agent(s) only (repeatable) |
| `--verify` | `true` | Verify cosign signature |
| `--verify-oidc-issuer` | — | Expected OIDC issuer for verification |
| `--verify-oidc-subject-regex` | — | Expected OIDC subject regex for verification |
| `--verify-trusted-root` | — | Path to trusted_root.json |

### Usage

```bash
# Install all skills from a repository
flux-operator skills install ghcr.io/org/flux-skills

# Install a specific tag
flux-operator skills install ghcr.io/org/flux-skills --tag v1.2.0

# Install only specific skills
flux-operator skills install ghcr.io/org/flux-skills \
  --skill deployment-helper \
  --skill troubleshooter

# Skip verification (not recommended)
flux-operator skills install ghcr.io/org/flux-skills --verify=false
```

---

## skills list

List all installed skills and their sources.

```bash
flux-operator skills list
```

No additional flags. Shows each installed skill, its source repository, and version.

---

## skills update

Check for updates and install them.

```bash
flux-operator skills update [flags]
```

| Flag | Description |
|------|-------------|
| `--verify-trusted-root` | Path to trusted_root.json for verification |
| `--dry-run` | Check for updates without installing |

### Usage

```bash
# Check and install updates
flux-operator skills update

# Dry run — see what would be updated
flux-operator skills update --dry-run
```

---

## skills uninstall

Remove installed skills.

```bash
flux-operator skills uninstall <repository> [flags]
```

| Flag | Description |
|------|-------------|
| `--all` | Uninstall all skills from all repositories |

### Usage

```bash
# Uninstall skills from a specific repository
flux-operator skills uninstall ghcr.io/org/flux-skills

# Uninstall everything
flux-operator skills uninstall --all
```

---

## skills publish

Package local skills and push them as an OCI artifact.

```bash
flux-operator skills publish <repository> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | `skills` | Local directory containing skills to package |
| `--tag` | — | OCI tag(s) to push (repeatable) |
| `--diff-tag` | — | Tag to diff against (only publish changed skills) |
| `--annotation, -a` | — | OCI annotations (repeatable, `key=value`) |
| `--sign` | `false` | Sign the artifact with cosign |
| `-o, --output` | — | Output format: `json` |

### Usage

```bash
# Publish with a version tag
flux-operator skills publish ghcr.io/org/flux-skills --tag v1.0.0

# Publish with multiple tags
flux-operator skills publish ghcr.io/org/flux-skills \
  --tag v1.0.0 --tag latest

# Publish from a custom directory
flux-operator skills publish ghcr.io/org/flux-skills \
  --path ./my-skills --tag v1.0.0

# Sign the artifact
flux-operator skills publish ghcr.io/org/flux-skills \
  --tag v1.0.0 --sign

# Only publish skills that changed since a previous tag
flux-operator skills publish ghcr.io/org/flux-skills \
  --tag v1.1.0 --diff-tag v1.0.0

# Output publish metadata as JSON
flux-operator skills publish ghcr.io/org/flux-skills \
  --tag v1.0.0 -o json
```
