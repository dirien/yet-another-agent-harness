# OCI Artifact & Image Automation Commands

## Table of Contents

1. [Push Artifact](#push-artifact)
2. [Pull Artifact](#pull-artifact)
3. [Tag Artifact](#tag-artifact)
4. [List Artifacts](#list-artifacts)
5. [Diff Artifact](#diff-artifact)
6. [Image Repository](#image-repository)
7. [Image Policy](#image-policy)
8. [Image Update Automation](#image-update-automation)

---

## Push Artifact

Create a tarball from a directory or file and push it to an OCI registry.

```bash
flux push artifact oci://<registry>/<name>:<tag> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `-f, --path` | Path to directory or file to package (required) |
| `--source` | Source URL (usually Git repo URL) |
| `--revision` | Source revision: `<branch\|tag>@sha1:<commit-sha>` |
| `-a, --annotations` | Custom OCI annotations: `key=value` (repeatable) |
| `--provider` | Auth provider: `generic`, `aws`, `azure`, `gcp` |
| `--creds` | Registry credentials: `<username>[:<password>]` |
| `-o, --output` | Output format: `json`, `yaml` |
| `--reproducible` | Ensure reproducible digests |
| `--ignore-paths` | Paths to ignore (.gitignore format) |
| `--insecure-registry` | Allow push without TLS |
| `--debug` | Show underlying library logs |

### Examples

```bash
# Push manifests
flux push artifact oci://ghcr.io/org/my-app:v1.0.0 \
  --path=./deploy \
  --source=https://github.com/org/my-app \
  --revision=main@sha1:abc1234

# Push with annotations
flux push artifact oci://ghcr.io/org/my-app:v1.0.0 \
  --path=./deploy \
  --annotations="org.opencontainers.image.description=My App manifests"

# Push to ECR
flux push artifact oci://123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.0.0 \
  --path=./deploy \
  --provider=aws
```

---

## Pull Artifact

Download an OCI artifact to a local directory.

```bash
flux pull artifact oci://<registry>/<name>:<tag> --output=<dir> [flags]
```

### Examples

```bash
flux pull artifact oci://ghcr.io/org/my-app:v1.0.0 --output=./downloaded
```

---

## Tag Artifact

Add a tag to an existing OCI artifact.

```bash
flux tag artifact oci://<registry>/<name>:<source-tag> --tag=<new-tag> [flags]
```

### Examples

```bash
# Promote a version to latest
flux tag artifact oci://ghcr.io/org/my-app:v1.0.0 --tag=latest

# Promote to a staging tag
flux tag artifact oci://ghcr.io/org/my-app:v1.0.0 --tag=staging
```

---

## List Artifacts

List all tags/versions of an OCI artifact.

```bash
flux list artifacts oci://<registry>/<name> [flags]
```

### Examples

```bash
flux list artifacts oci://ghcr.io/org/my-app
```

---

## Diff Artifact

Compare a local directory against an OCI artifact.

```bash
flux diff artifact oci://<registry>/<name>:<tag> --path=<local-dir> [flags]
```

### Examples

```bash
flux diff artifact oci://ghcr.io/org/my-app:v1.0.0 --path=./deploy
```

---

## Image Repository

Set up image scanning to detect new container image tags.

```bash
flux create image repository <name> [flags]
```

### Key Flags

| Flag | Description |
|------|-------------|
| `--image` | Container image to scan: `<registry>/<name>` (required) |
| `--scan-timeout` | Timeout for image scanning |
| `--secret-ref` | Secret with registry credentials |
| `--cert-secret-ref` | Secret with TLS certificates |
| `--provider` | Auth provider: `generic`, `aws`, `azure`, `gcp` |
| `--interval` | Scan interval |
| `--export` | Output YAML |

### Examples

```bash
# Scan Docker Hub image
flux create image repository my-app \
  --image=docker.io/org/my-app \
  --interval=5m

# Scan private registry
flux create image repository my-app \
  --image=ghcr.io/org/my-app \
  --secret-ref=ghcr-auth \
  --interval=5m
```

---

## Image Policy

Define which image tags to track and select for updates.

```bash
flux create image policy <name> [flags]
```

### Key Flags

| Flag | Description |
|------|-------------|
| `--image-ref` | Name of the ImageRepository to reference (required) |
| `--select-semver` | Semver range to filter tags (e.g. `>=1.0.0 <2.0.0`) |
| `--select-alpha` | Alphabetical ordering for tag selection |
| `--select-numeric` | Numeric ordering for tag selection |
| `--filter-regex` | Regex to filter tags |
| `--filter-extract` | Regex extraction for tag comparison |
| `--interval` | Reconciliation interval |
| `--export` | Output YAML |

### Examples

```bash
# Track semver tags
flux create image policy my-app \
  --image-ref=my-app \
  --select-semver=">=1.0.0"

# Track numeric tags
flux create image policy my-app \
  --image-ref=my-app \
  --select-numeric=asc \
  --filter-regex="^main-(?P<ts>[0-9]+)" \
  --filter-extract='$ts'
```

---

## Image Update Automation

Automatically commit image tag updates back to Git.

```bash
flux create image update <name> [flags]
```

### Key Flags

| Flag | Description |
|------|-------------|
| `--git-repo-ref` | GitRepository source reference (required) |
| `--git-repo-path` | Path within the repo to update |
| `--git-repo-namespace` | Namespace of the GitRepository |
| `--checkout-branch` | Branch to checkout |
| `--push-branch` | Branch to push updates to |
| `--author-name` | Git author name for commits |
| `--author-email` | Git author email for commits |
| `--commit-template` | Template for commit messages |
| `--interval` | Reconciliation interval |
| `--export` | Output YAML |

### Examples

```bash
# Auto-update images in a repo
flux create image update my-update \
  --git-repo-ref=my-app \
  --git-repo-path=./deploy \
  --checkout-branch=main \
  --push-branch=main \
  --author-name=flux \
  --author-email=flux@example.com

# Push to a separate branch (for PRs)
flux create image update my-update \
  --git-repo-ref=my-app \
  --git-repo-path=./deploy \
  --checkout-branch=main \
  --push-branch=flux-image-updates \
  --author-name=flux
```

---

## Querying Image Resources

```bash
# List all image automation objects
flux get images all [-A]

# By type
flux get images repository [-A]
flux get images policy [-A]
flux get images update [-A]

# Force reconciliation
flux reconcile image repository <name>
flux reconcile image policy <name>
flux reconcile image update <name>

# Suspend/resume
flux suspend image repository <name>
flux resume image repository <name>
```
