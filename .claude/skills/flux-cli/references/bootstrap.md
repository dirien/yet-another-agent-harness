# Bootstrap Commands

Bootstrap pushes Flux manifests to a Git repository and deploys Flux on the cluster.

---

## flux bootstrap github

```bash
flux bootstrap github [flags]
```

### GitHub-Specific Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--owner` | — | GitHub user or organization (required) |
| `--repository` | — | Repository name (required) |
| `--team` | — | Team slugs with access to the repo (repeatable) |
| `--personal` | `false` | Use a personal repo (not organization) |
| `--hostname` | `github.com` | GitHub Enterprise hostname |
| `--ssh-hostname` | — | SSH hostname (for GHE with different SSH host) |
| `--reconcile` | `false` | Reconcile existing repo without pushing manifests |

### Example

```bash
flux bootstrap github \
  --owner=my-org \
  --repository=fleet-infra \
  --branch=main \
  --path=clusters/production \
  --personal \
  --token-auth
```

---

## flux bootstrap gitlab

```bash
flux bootstrap gitlab [flags]
```

### GitLab-Specific Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--owner` | — | GitLab user or group (required) |
| `--repository` | — | Repository name (required) |
| `--hostname` | `gitlab.com` | GitLab hostname |
| `--personal` | `false` | Use a personal repo |
| `--team` | — | Teams with access |
| `--deploy-token-auth` | `false` | Use deploy token for auth |
| `--read-write-key` | `false` | Use read-write deploy key |
| `--reconcile` | `false` | Reconcile without pushing |

---

## flux bootstrap gitea

```bash
flux bootstrap gitea [flags]
```

### Gitea-Specific Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--owner` | — | Gitea user or organization (required) |
| `--repository` | — | Repository name (required) |
| `--hostname` | — | Gitea hostname (required) |
| `--personal` | `false` | Use a personal repo |

---

## flux bootstrap bitbucket-server

```bash
flux bootstrap bitbucket-server [flags]
```

### Bitbucket-Specific Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--owner` | — | Bitbucket project key (required) |
| `--repository` | — | Repository slug (required) |
| `--hostname` | — | Bitbucket Server hostname (required) |
| `--group` | — | Bitbucket groups with access |

---

## flux bootstrap git

Generic bootstrap for any Git provider.

```bash
flux bootstrap git [flags]
```

### Git-Specific Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--url` | — | Git repository URL (required) |

---

## Common Bootstrap Flags (all providers)

### Git Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--branch` | `main` | Git branch |
| `--path` | — | Path within the repo for Flux manifests |
| `--author-email` | — | Email for Git commits |
| `--author-name` | `Flux` | Name for Git commits |
| `--commit-message-appendix` | — | Text appended to commit messages |

### Component Selection

| Flag | Default | Description |
|------|---------|-------------|
| `--components` | `source-controller,kustomize-controller,helm-controller,notification-controller` | Controllers to install |
| `--components-extra` | — | Additional controllers (e.g. `image-reflector-controller,image-automation-controller`) |

### SSH Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--private-key-file` | — | Path to SSH private key |
| `--ssh-key-algorithm` | `ecdsa` | Key algorithm: `rsa`, `ecdsa`, `ed25519` |
| `--ssh-ecdsa-curve` | `p384` | ECDSA curve: `p256`, `p384`, `p521` |
| `--ssh-rsa-bits` | `2048` | RSA key size |

### Authentication

| Flag | Default | Description |
|------|---------|-------------|
| `--token-auth` | `false` | Use PAT instead of SSH |
| `--ca-file` | — | TLS CA file for self-signed certs |

### Registry & Images

| Flag | Default | Description |
|------|---------|-------------|
| `--registry` | `ghcr.io/fluxcd` | Controller image registry |
| `--image-pull-secret` | — | Secret for private registries |
| `--registry-creds` | — | Credentials in `user:password` format |
| `--version` | — | Toolkit version |

### Cluster Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--cluster-domain` | `cluster.local` | Internal cluster domain |
| `--watch-all-namespaces` | `true` | Monitor all namespaces |
| `--network-policy` | `true` | Enable network policies |
| `--toleration-keys` | — | Toleration keys for scheduling |

### GPG Signing

| Flag | Description |
|------|-------------|
| `--gpg-key-id` | GPG key ID for signing |
| `--gpg-key-ring` | Path to GPG keyring |
| `--gpg-passphrase` | GPG passphrase |

### General

| Flag | Default | Description |
|------|---------|-------------|
| `--log-level` | `info` | Log verbosity: `debug`, `info`, `error` |
| `--force` | `false` | Override existing Flux installation |

---

## Workflow: Bootstrap from scratch

```bash
# 1. Export your token
export GITHUB_TOKEN=ghp_xxxx

# 2. Bootstrap
flux bootstrap github \
  --owner=my-org \
  --repository=fleet-infra \
  --branch=main \
  --path=clusters/staging \
  --personal

# 3. Verify
flux check
flux get all -A
```

## Workflow: Bootstrap with image automation

```bash
flux bootstrap github \
  --owner=my-org \
  --repository=fleet-infra \
  --branch=main \
  --path=clusters/production \
  --components-extra=image-reflector-controller,image-automation-controller \
  --read-write-key \
  --personal
```

## Workflow: Re-bootstrap (upgrade)

```bash
# Run bootstrap again with a new version — it's idempotent
flux bootstrap github \
  --owner=my-org \
  --repository=fleet-infra \
  --branch=main \
  --path=clusters/production \
  --personal
```
