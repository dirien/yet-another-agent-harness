# Source Commands

Sources tell Flux where to pull manifests and charts from.

---

## flux create source git

Create or update a GitRepository source.

```bash
flux create source git <name> --url=<repo-url> [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--url` | — | Git repository URL (required) |
| `--branch` | — | Git branch to track |
| `--tag` | — | Specific Git tag |
| `--tag-semver` | — | Semver range for tags (e.g. `>=1.0.0`) |
| `--commit` | — | Specific commit SHA |
| `--ref-name` | — | Git reference name |
| `-u, --username` | — | Basic auth username |
| `-p, --password` | — | Basic auth password |
| `--private-key-file` | — | SSH private key file |
| `--ssh-key-algorithm` | `ecdsa` | SSH key algorithm: `rsa`, `ecdsa`, `ed25519` |
| `--ssh-ecdsa-curve` | `p384` | ECDSA curve: `p256`, `p384`, `p521` |
| `--ssh-rsa-bits` | `2048` | RSA key size |
| `--secret-ref` | — | Existing secret with SSH or auth credentials |
| `--provider` | `generic` | Git provider: `generic`, `azure`, `github` |
| `--ca-file` | — | TLS CA file for self-signed certs |
| `--proxy-secret-ref` | — | Existing secret with proxy credentials |
| `--recurse-submodules` | `false` | Initialize and include Git submodules |
| `--ignore-paths` | — | Paths to ignore (comma-separated) |
| `--sparse-checkout-paths` | — | Paths for sparse checkout (comma-separated) |
| `-s, --silent` | `false` | Skip deploy key confirmation |
| `--interval` | `1m` | Sync interval |
| `--export` | — | Output YAML instead of applying |

### Examples

```bash
# Public repo
flux create source git my-app \
  --url=https://github.com/org/app \
  --branch=main

# Private repo with basic auth
flux create source git my-app \
  --url=https://github.com/org/app \
  --branch=main \
  --username=bot \
  --password=ghp_xxxx

# Private repo with SSH
flux create source git my-app \
  --url=ssh://git@github.com/org/app \
  --branch=main \
  --private-key-file=~/.ssh/id_ed25519

# Track semver tags
flux create source git my-app \
  --url=https://github.com/org/app \
  --tag-semver=">=1.0.0 <2.0.0"

# Export as YAML
flux create source git my-app \
  --url=https://github.com/org/app \
  --branch=main \
  --export > source.yaml
```

---

## flux create source helm

Create or update a HelmRepository source.

```bash
flux create source helm <name> --url=<chart-repo-url> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--url` | Helm repository URL (required) |
| `-u, --username` | Basic auth username |
| `-p, --password` | Basic auth password |
| `--cert-file` | TLS client certificate file |
| `--key-file` | TLS client key file |
| `--ca-file` | TLS CA certificate file |
| `--secret-ref` | Existing secret with credentials |
| `--oci-provider` | OCI auth provider for OCI Helm repos |
| `--pass-credentials` | Pass credentials to all domains |
| `--interval` | Sync interval (default `1m`) |
| `--export` | Output YAML instead of applying |

### Examples

```bash
# Public chart repo
flux create source helm bitnami \
  --url=https://charts.bitnami.com/bitnami \
  --interval=1h

# With basic auth
flux create source helm private-charts \
  --url=https://charts.example.com \
  --username=admin \
  --password=secret

# OCI-based Helm repo
flux create source helm oci-charts \
  --url=oci://ghcr.io/org/charts \
  --username=bot \
  --password=token
```

---

## flux create source oci

Create or update an OCIRepository source.

```bash
flux create source oci <name> --url=<oci-url> [flags]
```

### Key Flags

| Flag | Description |
|------|-------------|
| `--url` | OCI repository URL (required) |
| `--tag` | OCI tag to track |
| `--tag-semver` | Semver range for tags |
| `--digest` | Specific artifact digest |
| `--secret-ref` | Existing secret with registry credentials |
| `--provider` | OCI auth provider: `generic`, `aws`, `azure`, `gcp` |
| `--insecure` | Allow HTTP registry |
| `--interval` | Sync interval |
| `--export` | Output YAML instead of applying |

### Examples

```bash
# Track an OCI artifact
flux create source oci my-manifests \
  --url=oci://ghcr.io/org/manifests \
  --tag=latest

# Track with semver
flux create source oci my-manifests \
  --url=oci://ghcr.io/org/manifests \
  --tag-semver=">=1.0.0"
```

---

## flux create source bucket

Create or update a Bucket source (S3-compatible, GCS, Azure Blob).

```bash
flux create source bucket <name> --bucket-name=<name> --endpoint=<url> [flags]
```

### Key Flags

| Flag | Description |
|------|-------------|
| `--bucket-name` | Bucket name (required) |
| `--endpoint` | Bucket endpoint URL (required) |
| `--provider` | Provider: `generic`, `aws`, `azure`, `gcp` |
| `--region` | Bucket region |
| `--secret-ref` | Existing secret with access credentials |
| `--access-key` | Access key ID |
| `--secret-key` | Secret access key |
| `--insecure` | Allow HTTP |
| `--interval` | Sync interval |
| `--export` | Output YAML instead of applying |

### Examples

```bash
# S3 bucket
flux create source bucket my-bucket \
  --bucket-name=my-manifests \
  --endpoint=s3.amazonaws.com \
  --provider=aws \
  --region=us-east-1

# MinIO
flux create source bucket minio \
  --bucket-name=flux \
  --endpoint=minio.example.com \
  --access-key=admin \
  --secret-key=password \
  --insecure
```

---

## Querying Sources

```bash
# List all sources
flux get sources all [-A]

# List by type
flux get sources git [-A] [-w] [--no-header]
flux get sources helm [-A]
flux get sources oci [-A]
flux get sources bucket [-A]
flux get sources chart [-A]

# Filter by label
flux get sources git -l team=backend

# Filter by status
flux get sources git --status-selector ready=false
```

## Managing Sources

```bash
# Force reconciliation
flux reconcile source git <name>
flux reconcile source helm <name>
flux reconcile source oci <name>
flux reconcile source bucket <name>
flux reconcile source chart <name>

# Suspend/resume
flux suspend source git <name>
flux resume source git <name>

# Export
flux export source git <name>
flux export source git --all > all-git-sources.yaml

# Delete
flux delete source git <name>
```
