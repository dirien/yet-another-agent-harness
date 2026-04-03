# Kustomization Commands

Kustomizations tell Flux how to deploy manifests from a source onto the cluster.

---

## flux create kustomization

```bash
flux create kustomization <name> [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--source` | — | Source reference: `[<kind>/]<name>.<namespace>` (required). Kind is `GitRepository`, `OCIRepository`, or `Bucket` |
| `--path` | `./` | Path within the source to the kustomization.yaml or manifests directory |
| `--prune` | `false` | Enable garbage collection — delete objects removed from source |
| `--depends-on` | — | Dependencies that must be ready first: `<name>` or `<namespace>/<name>` (repeatable) |
| `--target-namespace` | — | Override namespace for all reconciled objects |
| `--service-account` | — | Service account for impersonation during reconciliation |
| `--decryption-provider` | — | Decryption provider: `sops` |
| `--decryption-secret` | — | Secret with OpenPGP or age keys for SOPS |
| `--health-check` | — | Workloads for health assessment: `<kind>/<name>.<namespace>` (repeatable) |
| `--health-check-timeout` | `2m` | Timeout for health checks |
| `--wait` | `false` | Enable health checking of applied resources |
| `--retry-interval` | — | Retry interval for failed reconciliation |
| `--kubeconfig-secret-ref` | — | Secret with kubeconfig for remote cluster |
| `--interval` | `1m` | Reconciliation interval |
| `--export` | — | Output YAML instead of applying |
| `--label` | — | Labels as `key=value` (repeatable) |

### Examples

```bash
# Basic deployment from a GitRepository
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy/production \
  --prune=true \
  --interval=5m

# With SOPS decryption
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy \
  --prune=true \
  --decryption-provider=sops \
  --decryption-secret=sops-age

# With dependencies
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy \
  --prune=true \
  --depends-on=infrastructure

# Health checks
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy \
  --prune=true \
  --health-check="Deployment/my-app.default" \
  --health-check-timeout=3m

# Target a remote cluster
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy \
  --kubeconfig-secret-ref=staging-kubeconfig

# Export YAML
flux create kustomization my-app \
  --source=GitRepository/my-app \
  --path=./deploy \
  --prune=true \
  --export > kustomization.yaml
```

---

## flux build kustomization

Build a Kustomization locally and output the resulting multi-doc YAML.

```bash
flux build kustomization <name> --path=<local-path> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--path` | Path to local manifests directory (required) |
| `--kustomization-file` | Path to a local Flux Kustomization YAML file |
| `--dry-run` | Run without cluster connection (variable substitutions from Secrets/ConfigMaps are skipped) |
| `-r, --recursive` | Recursively build nested Kustomizations |
| `--local-sources` | Map local paths to source references: `Kind/namespace/name=path` |
| `--ignore-paths` | Exclude files using .gitignore patterns |
| `--strict-substitute` | Fail if vars without defaults are missing |

### How it works

1. Fetches the specified Kustomization from the cluster (unless `--dry-run`)
2. Uses it to render the local kustomization.yaml
3. Outputs the resulting manifests to stdout

### Examples

```bash
# Basic build
flux build kustomization my-app --path=./deploy

# Dry run (no cluster needed)
flux build kustomization my-app \
  --path=./deploy \
  --kustomization-file=./flux/my-app.yaml \
  --dry-run

# Recursive with local sources
flux build kustomization my-app \
  --path=./deploy \
  --recursive \
  --local-sources GitRepository/flux-system/my-repo=./

# Exclude files
flux build kustomization my-app \
  --path=./deploy \
  --ignore-paths="/tests/**,*.test.yaml"
```

---

## flux diff kustomization

Build locally, perform a server-side dry-run, and print the diff against the cluster.

```bash
flux diff kustomization <name> --path=<local-path> [flags]
```

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | No differences |
| `1` | Differences found |
| `>1` | Error occurred |

### Flags

Same as `flux build kustomization`, plus:

| Flag | Default | Description |
|------|---------|-------------|
| `--progress-bar` | `true` | Show progress bar |

### Examples

```bash
# Basic diff
flux diff kustomization my-app --path=./deploy

# Use in CI (exit code 1 = changes detected)
if ! flux diff kustomization my-app --path=./deploy; then
  echo "Changes detected, review required"
fi

# With local Kustomization file
flux diff kustomization my-app \
  --path=./deploy \
  --kustomization-file=./flux/my-app.yaml

# Recursive
flux diff kustomization my-app \
  --path=./deploy \
  --recursive \
  --local-sources GitRepository/flux-system/my-repo=./
```

---

## Managing Kustomizations

```bash
# List
flux get kustomizations [-A] [-w] [--no-header]

# Filter
flux get kustomizations --status-selector ready=false
flux get kustomizations -l app=frontend

# Force reconciliation
flux reconcile kustomization <name>

# Suspend/resume
flux suspend kustomization <name>
flux suspend kustomization --all
flux resume kustomization <name>

# Export
flux export kustomization <name>
flux export kustomization --all > all-ks.yaml

# Delete
flux delete kustomization <name>
```
