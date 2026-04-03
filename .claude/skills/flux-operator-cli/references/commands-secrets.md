# Create Secret Commands

All secret commands follow this pattern:

```bash
flux-operator create secret <type> <name> [flags]
```

### Common Flags (all secret types)

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace for the secret |
| `--annotation` | Comma-separated annotations: `key1=val1,key2=val2` |
| `--label` | Comma-separated labels: `key1=val1,key2=val2` |
| `--immutable` | Create an immutable secret |
| `--export` | Output YAML instead of creating in-cluster (for GitOps) |

---

## basic-auth

Create a secret with username/password credentials.

```bash
flux-operator create secret basic-auth <name> \
  --username=<user> \
  --password=<pass>
```

| Flag | Description |
|------|-------------|
| `--username` | Username (required) |
| `--password` | Password (required; mutually exclusive with `--password-stdin`) |
| `--password-stdin` | Read password from stdin |

### Usage

```bash
# Create basic auth secret
flux-operator create secret basic-auth git-auth \
  --username=bot --password=ghp_xxxx -n flux-system

# Export as YAML
flux-operator create secret basic-auth git-auth \
  --username=bot --password=ghp_xxxx --export

# Read password from stdin
echo "ghp_xxxx" | flux-operator create secret basic-auth git-auth \
  --username=bot --password-stdin -n flux-system
```

---

## githubapp

Create a secret for GitHub App authentication.

```bash
flux-operator create secret githubapp <name> \
  --app-id=<id> \
  --app-installation-id=<id> \
  --app-private-key-file=<path>
```

| Flag | Description |
|------|-------------|
| `--app-id` | GitHub App ID (required) |
| `--app-installation-id` | GitHub App installation ID (required) |
| `--app-private-key-file` | Path to private key PEM file (required) |
| `--app-base-url` | GitHub Enterprise base URL (optional) |

### Usage

```bash
# Standard GitHub
flux-operator create secret githubapp gh-app \
  --app-id=12345 \
  --app-installation-id=67890 \
  --app-private-key-file=private-key.pem \
  -n flux-system

# GitHub Enterprise
flux-operator create secret githubapp gh-app \
  --app-id=12345 \
  --app-installation-id=67890 \
  --app-private-key-file=private-key.pem \
  --app-base-url=https://github.example.com \
  -n flux-system
```

---

## proxy

Create a secret for HTTP/HTTPS proxy configuration.

```bash
flux-operator create secret proxy <name> \
  --address=<url>
```

| Flag | Description |
|------|-------------|
| `--address` | Proxy URL (required) |
| `--username` | Proxy username (optional) |
| `--password` | Proxy password (optional) |
| `--password-stdin` | Read password from stdin |

### Usage

```bash
# Proxy without auth
flux-operator create secret proxy corp-proxy \
  --address=http://proxy.corp.example.com:3128 \
  -n flux-system

# Proxy with auth
flux-operator create secret proxy corp-proxy \
  --address=http://proxy.corp.example.com:3128 \
  --username=proxyuser --password=proxypass \
  -n flux-system
```

---

## registry

Create a secret for container registry authentication.

```bash
flux-operator create secret registry <name> \
  --server=<host> \
  --username=<user> \
  --password=<pass>
```

| Flag | Description |
|------|-------------|
| `--server` | Registry server hostname (required) |
| `--username` | Registry username (required) |
| `--password` | Registry password (required; mutually exclusive with `--password-stdin`) |
| `--password-stdin` | Read password from stdin |

### Usage

```bash
# GitHub Container Registry
flux-operator create secret registry ghcr-auth \
  --server=ghcr.io \
  --username=bot \
  --password=ghp_xxxx \
  -n flux-system

# Docker Hub
flux-operator create secret registry dockerhub \
  --server=docker.io \
  --username=myuser \
  --password-stdin \
  -n flux-system
```

---

## sops

Create a secret for SOPS decryption (age or GPG keys).

```bash
flux-operator create secret sops <name> [flags]
```

| Flag | Description |
|------|-------------|
| `--age-key-file` | Path to age key file (repeatable for multiple keys) |
| `--gpg-key-file` | Path to GPG key file (repeatable for multiple keys) |
| `--age-key-stdin` | Read age key from stdin |
| `--gpg-key-stdin` | Read GPG key from stdin |

### Usage

```bash
# Single age key
flux-operator create secret sops sops-age \
  --age-key-file=age.key \
  -n flux-system

# Multiple age keys
flux-operator create secret sops sops-age \
  --age-key-file=key1.txt \
  --age-key-file=key2.txt \
  -n flux-system

# GPG key
flux-operator create secret sops sops-gpg \
  --gpg-key-file=private.gpg \
  -n flux-system

# Age key from stdin
cat age.key | flux-operator create secret sops sops-age \
  --age-key-stdin \
  -n flux-system
```

---

## ssh

Create a secret for SSH authentication (Git over SSH).

```bash
flux-operator create secret ssh <name> \
  --private-key-file=<path> \
  --knownhosts-file=<path>
```

| Flag | Description |
|------|-------------|
| `--private-key-file` | Path to SSH private key (required) |
| `--public-key-file` | Path to SSH public key (optional) |
| `--knownhosts-file` | Path to known_hosts file (required) |
| `--password` | Passphrase for the private key (optional) |
| `--password-stdin` | Read passphrase from stdin |

### Usage

```bash
# Standard SSH key
flux-operator create secret ssh git-ssh \
  --private-key-file=~/.ssh/id_ed25519 \
  --knownhosts-file=~/.ssh/known_hosts \
  -n flux-system

# With passphrase
flux-operator create secret ssh git-ssh \
  --private-key-file=~/.ssh/id_rsa \
  --knownhosts-file=~/.ssh/known_hosts \
  --password=mypassphrase \
  -n flux-system

# Export for GitOps
flux-operator create secret ssh git-ssh \
  --private-key-file=id_ed25519 \
  --knownhosts-file=known_hosts \
  --export > secret.yaml
```

---

## tls

Create a secret with TLS certificate and key.

```bash
flux-operator create secret tls <name> [flags]
```

| Flag | Description |
|------|-------------|
| `--tls-crt-file` | Path to TLS certificate file |
| `--tls-key-file` | Path to TLS private key file |
| `--ca-crt-file` | Path to CA certificate file |

### Usage

```bash
# TLS cert + key
flux-operator create secret tls my-tls \
  --tls-crt-file=cert.pem \
  --tls-key-file=key.pem \
  -n flux-system

# With CA certificate
flux-operator create secret tls my-tls \
  --tls-crt-file=cert.pem \
  --tls-key-file=key.pem \
  --ca-crt-file=ca.pem \
  -n flux-system
```

---

## web-config

Create a secret for web UI OIDC configuration.

```bash
flux-operator create secret web-config <name> \
  --base-url=<url> \
  --issuer-url=<url>
```

| Flag | Default | Description |
|------|---------|-------------|
| `--base-url` | — | Base URL of the web UI (required) |
| `--provider` | `OIDC` | Auth provider |
| `--issuer-url` | — | OIDC issuer URL (required for OIDC provider) |
| `--client-id` | — | OIDC client ID |
| `--client-secret` | — | OIDC client secret |
| `--client-secret-stdin` | — | Read client secret from stdin |
| `--client-secret-rnd` | — | Generate a random client secret |

### Usage

```bash
# OIDC config
flux-operator create secret web-config web-oidc \
  --base-url=https://flux.example.com \
  --issuer-url=https://auth.example.com \
  --client-id=flux-ui \
  --client-secret=my-secret \
  -n flux-system

# With random client secret
flux-operator create secret web-config web-oidc \
  --base-url=https://flux.example.com \
  --issuer-url=https://auth.example.com \
  --client-id=flux-ui \
  --client-secret-rnd \
  -n flux-system
```
