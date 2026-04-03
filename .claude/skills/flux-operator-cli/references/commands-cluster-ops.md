# Cluster Operations Commands

These commands require cluster access via `~/.kube/config`.

---

## get

Retrieve Flux Operator resources and their status.

### get instance

```bash
flux-operator get instance [-n <namespace>] [-A]
```

Retrieves FluxInstance resources.

### get rset

```bash
flux-operator get rset [-n <namespace>] [-A]
```

Retrieves ResourceSet resources.

### get rsip

```bash
flux-operator get rsip [-n <namespace>] [-A]
```

Retrieves ResourceSetInputProvider resources.

### get all

```bash
flux-operator get all [-n <namespace>] [-A] [--kind <kind>] [--ready-status <status>] [-o <format>]
```

Retrieves all Flux resources with their status.

| Flag | Default | Description |
|------|---------|-------------|
| `-n, --namespace` | — | Filter by namespace |
| `-A, --all-namespaces` | — | Retrieve from all namespaces |
| `--kind` | — | Filter by resource kind (e.g. `Kustomization`, `HelmRelease`) |
| `--ready-status` | — | Filter by status: `True`, `False`, `Unknown`, `Suspended` |
| `-o, --output` | `table` | Output format: `table`, `json`, `yaml` |

### Usage

```bash
# All resources across all namespaces
flux-operator get all -A

# Only failing resources
flux-operator get all --ready-status=False

# Only HelmReleases as JSON
flux-operator get all --kind=HelmRelease -o json

# Resources in a specific namespace
flux-operator get all -n production
```

---

## export

Export Flux resources for backup or migration.

### export report

```bash
flux-operator export report [-n <namespace>] [-o <format>]
```

Exports the FluxReport containing distribution status, installed versions, and health.

### export resource

```bash
flux-operator export resource <kind>/<name> [-n <namespace>] [-o <format>]
```

Exports a specific Flux resource.

| Flag | Default | Description |
|------|---------|-------------|
| `-n, --namespace` | — | Namespace scope |
| `-o, --output` | `yaml` | Output format: `yaml`, `json` |

### Usage

```bash
# Export the flux report
flux-operator export report -n flux-system

# Export a Kustomization as YAML
flux-operator export resource Kustomization/my-app -n default

# Export a HelmRelease as JSON
flux-operator export resource HelmRelease/nginx -n default -o json
```

---

## reconcile

Trigger reconciliation of Flux resources.

### Subcommands

```bash
flux-operator reconcile instance <name> [-n <namespace>] [--wait]
flux-operator reconcile rset <name> [-n <namespace>] [--wait]
flux-operator reconcile rsip <name> [-n <namespace>] [--wait]
flux-operator reconcile resource <kind>/<name> [-n <namespace>] [--wait]
flux-operator reconcile all [-n <namespace>] [--wait]
```

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace scope |
| `--wait` | Wait for reconciliation to complete before returning |

### Usage

```bash
# Reconcile a specific Kustomization and wait
flux-operator reconcile resource Kustomization/my-app -n default --wait

# Reconcile a FluxInstance
flux-operator reconcile instance flux -n flux-system --wait

# Reconcile everything in the cluster
flux-operator reconcile all --wait
```

---

## suspend

Pause reconciliation of Flux resources. The resource stays in place but won't be
reconciled until resumed.

### Subcommands

```bash
flux-operator suspend instance <name> [-n <namespace>]
flux-operator suspend rset <name> [-n <namespace>]
flux-operator suspend rsip <name> [-n <namespace>]
flux-operator suspend resource <kind>/<name> [-n <namespace>]
```

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace scope |

### Usage

```bash
# Suspend a FluxInstance before maintenance
flux-operator suspend instance flux -n flux-system

# Suspend a specific HelmRelease
flux-operator suspend resource HelmRelease/nginx -n default
```

---

## resume

Resume reconciliation of a previously suspended resource.

### Subcommands

```bash
flux-operator resume instance <name> [-n <namespace>] [--wait]
flux-operator resume rset <name> [-n <namespace>] [--wait]
flux-operator resume rsip <name> [-n <namespace>] [--wait]
flux-operator resume resource <kind>/<name> [-n <namespace>] [--wait]
```

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace scope |
| `--wait` | Wait for the first reconciliation after resume to complete |

### Usage

```bash
# Resume and wait for reconciliation
flux-operator resume instance flux -n flux-system --wait

# Resume a HelmRelease
flux-operator resume resource HelmRelease/nginx -n default --wait
```

---

## delete

Remove Flux Operator resources from the cluster.

### Subcommands

```bash
flux-operator delete instance <name> [-n <namespace>] [--wait] [--with-suspend]
flux-operator delete rset <name> [-n <namespace>] [--wait] [--with-suspend]
flux-operator delete rsip <name> [-n <namespace>] [--wait] [--with-suspend]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-n, --namespace` | — | Namespace scope |
| `--wait` | `true` | Wait for deletion to complete |
| `--with-suspend` | `false` | Suspend the resource before deleting, leaving managed resources in-place |

### Usage

```bash
# Delete a ResourceSet
flux-operator delete rset my-rset -n default

# Delete but keep managed resources intact
flux-operator delete instance flux -n flux-system --with-suspend
```

The `--with-suspend` flag is a safety mechanism: it suspends reconciliation first, then
deletes the operator resource. The Kubernetes objects that were managed by the resource
remain untouched in the cluster.

---

## stats

Display reconciliation statistics and cumulative storage usage per source type.

```bash
flux-operator stats
```

No additional flags. Shows a summary of how many resources are reconciled, failing,
suspended, and the storage consumed by each source type.

---

## trace

Trace a Kubernetes object through the GitOps delivery pipeline.

```bash
flux-operator trace <kind>/<name> [-n <namespace>]
```

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace of the object to trace |

Identifies which Flux reconciler manages the object and traces back to the original
source (Git repository, OCI artifact, Helm chart, etc.).

### Usage

```bash
# Trace a Deployment
flux-operator trace Deployment/my-app -n default

# Trace a Service
flux-operator trace Service/frontend -n production
```

---

## tree

Visualize Flux-managed objects as a hierarchical tree.

### Subcommands

```bash
flux-operator tree rset <name> [-n <namespace>]
flux-operator tree ks <name> [-n <namespace>]
flux-operator tree hr <name> [-n <namespace>]
```

| Subcommand | What it shows |
|------------|---------------|
| `tree rset` | Objects managed by a ResourceSet |
| `tree ks` | Objects managed by a Kustomization |
| `tree hr` | Objects managed by a HelmRelease |

| Flag | Description |
|------|-------------|
| `-n, --namespace` | Namespace scope |

### Usage

```bash
# View the tree of a Kustomization
flux-operator tree ks my-app -n default

# View HelmRelease managed objects
flux-operator tree hr nginx -n default
```

---

## wait

Poll a resource until it reaches the Ready state or times out.

### Subcommands

```bash
flux-operator wait instance <name> [-n <namespace>] [--timeout <duration>]
flux-operator wait rset <name> [-n <namespace>] [--timeout <duration>]
flux-operator wait rsip <name> [-n <namespace>] [--timeout <duration>]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-n, --namespace` | — | Namespace scope |
| `--timeout` | `1m` | How long to wait before giving up (e.g. `5m`, `2m30s`) |

### Usage

```bash
# Wait up to 5 minutes for an instance to be ready
flux-operator wait instance flux -n flux-system --timeout 5m
```

---

## install

Bootstrap a cluster with the Flux Operator and a FluxInstance. Downloads the operator
distribution from `oci://ghcr.io/controlplaneio-fluxcd/flux-operator-manifests` and
installs into the `flux-system` namespace.

This command is intended for development and testing. For production, use the
Flux Operator installation guide with Helm charts.

### Steps performed

1. Install the Flux Operator
2. Install a FluxInstance
3. Configure pull secret (if needed)
4. Bootstrap from Git or OCI source
5. Configure auto-updates (if requested)

### Key Flags

| Flag | Description |
|------|-------------|
| `-f, --instance-file` | FluxInstance YAML (local path, OCI, or HTTPS URL) |
| `--instance-distribution-version` | Flux distribution version |
| `--instance-distribution-registry` | Flux distribution registry |
| `--instance-distribution-artifact` | Flux distribution OCI artifact |
| `--instance-components` | List of Flux components to install |
| `--instance-components-extra` | Additional components beyond the default set |
| `--instance-cluster-type` | Cluster type: `kubernetes`, `openshift`, `aws`, `azure`, `gcp` |
| `--instance-cluster-size` | Cluster size: `small`, `medium`, `large` |
| `--instance-cluster-domain` | Cluster domain |
| `--instance-cluster-multitenant` | Enable multitenant lockdown |
| `--instance-cluster-network-policy` | Restrict network access |
| `--instance-sync-url` | Git or OCI repository URL to sync from |
| `--instance-sync-ref` | Git ref or OCI tag |
| `--instance-sync-path` | Path within the repository to the manifests |
| `--instance-sync-creds` | Credentials in `username:token` format |
| `--instance-sync-gha-app-id` | GitHub App ID for auth |
| `--instance-sync-gha-installation-id` | GitHub App installation ID |
| `--instance-sync-gha-installation-owner` | GitHub App installation owner |
| `--instance-sync-gha-private-key-file` | GitHub App private key file |
| `--instance-sync-gha-base-url` | GitHub Enterprise base URL |
| `--auto-update` | Enable automatic updates |
| `--verify` | Verify cosign signature |
| `--certificate-identity-regexp` | Certificate identity regex for verification |
| `--certificate-oidc-issuer` | OIDC issuer for verification |
| `--trusted-root` | Path to trusted_root.json |

### Usage

```bash
# Basic install
flux-operator install

# Install with Git sync
flux-operator install \
  --instance-sync-url=https://github.com/org/fleet \
  --instance-sync-ref=main \
  --instance-sync-path=clusters/production \
  --instance-sync-creds=username:ghp_token

# Install for AWS with large cluster profile
flux-operator install \
  --instance-cluster-type=aws \
  --instance-cluster-size=large \
  --instance-cluster-multitenant

# Install from a custom FluxInstance file
flux-operator install -f custom-instance.yaml
```

---

## uninstall

Safely remove the Flux Operator and its instance from the cluster.

```bash
flux-operator -n <namespace> uninstall [--keep-namespace]
```

### Steps performed

1. Delete role bindings
2. Delete deployments
3. Remove finalizers
4. Delete CRDs
5. Delete namespace (unless `--keep-namespace`)

Does **not** delete reconciled Kubernetes objects or Helm releases — those remain
in the cluster.

| Flag | Description |
|------|-------------|
| `--keep-namespace` | Don't delete the namespace after removing the operator |

### Usage

```bash
# Full removal
flux-operator -n flux-system uninstall

# Keep the namespace for reuse
flux-operator -n flux-system uninstall --keep-namespace
```

---

## version

Display CLI, Flux Operator, and Flux distribution versions.

```bash
flux-operator version [--client]
```

| Flag | Description |
|------|-------------|
| `--client` | Show only the client (CLI) version |

### Usage

```bash
# Full version info (requires cluster access)
flux-operator version

# Client version only (no cluster needed)
flux-operator version --client
```
