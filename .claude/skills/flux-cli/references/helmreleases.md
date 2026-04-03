# HelmRelease Commands

HelmReleases tell Flux how to install and manage Helm charts.

---

## flux create helmrelease

```bash
flux create helmrelease <name> [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--chart` | — | Helm chart name or path (required unless `--chart-ref` is set) |
| `--source` | — | Chart source: `<kind>/<name>.<namespace>`. Kind: `HelmRepository`, `GitRepository`, `Bucket` |
| `--chart-ref` | — | Reference to a HelmChart or OCIRepository: `<kind>/<name>.<namespace>` |
| `--chart-version` | — | Chart version (semver range accepted; ignored for Git sources) |
| `--chart-interval` | — | Interval to check for new chart versions |
| `--release-name` | — | Custom Helm release name (defaults to HelmRelease name) |
| `--target-namespace` | — | Namespace to install the chart into |
| `--storage-namespace` | — | Namespace for storing Helm release metadata |
| `--create-target-namespace` | `false` | Create target namespace if it doesn't exist |
| `--values` | — | Local path to values.yaml file(s) (repeatable) |
| `--values-from` | — | K8s object reference for values: `Secret/<name>` or `ConfigMap/<name>` (repeatable) |
| `--depends-on` | — | Dependencies: `<name>` or `<namespace>/<name>` (repeatable) |
| `--crds` | — | CRD upgrade policy: `Skip`, `Create`, `CreateReplace` |
| `--service-account` | — | Service account for impersonation |
| `--kubeconfig-secret-ref` | — | Secret with kubeconfig for remote cluster |
| `--reconcile-strategy` | `ChartVersion` | Strategy: `Revision` or `ChartVersion` |
| `--interval` | `1m` | Reconciliation interval |
| `--export` | — | Output YAML instead of applying |
| `--label` | — | Labels as `key=value` (repeatable) |

### Examples

```bash
# From a HelmRepository
flux create helmrelease nginx \
  --source=HelmRepository/bitnami \
  --chart=nginx \
  --chart-version=">=15.0.0" \
  --values=./nginx-values.yaml

# From a GitRepository (chart in repo)
flux create helmrelease my-app \
  --source=GitRepository/my-app \
  --chart=./charts/my-app

# From a Bucket
flux create helmrelease my-app \
  --source=Bucket/my-bucket \
  --chart=./charts/my-app

# With values from a Secret
flux create helmrelease my-app \
  --source=HelmRepository/my-charts \
  --chart=my-app \
  --values-from=Secret/my-app-values

# Multiple values files
flux create helmrelease my-app \
  --source=HelmRepository/my-charts \
  --chart=my-app \
  --values=./base-values.yaml \
  --values=./prod-values.yaml

# Custom release name and target namespace
flux create helmrelease my-app \
  --source=HelmRepository/my-charts \
  --chart=my-app \
  --release-name=my-custom-release \
  --target-namespace=apps \
  --create-target-namespace

# With dependencies
flux create helmrelease my-app \
  --source=HelmRepository/my-charts \
  --chart=my-app \
  --depends-on=cert-manager \
  --depends-on=kube-system/external-dns

# CRD handling
flux create helmrelease cert-manager \
  --source=HelmRepository/jetstack \
  --chart=cert-manager \
  --crds=CreateReplace

# Cross-namespace source reference
flux create helmrelease my-app \
  --source=HelmRepository/shared-charts.shared-ns \
  --chart=my-app

# From an OCIRepository
flux create helmrelease my-app \
  --chart-ref=OCIRepository/my-oci-source

# Export YAML
flux create helmrelease my-app \
  --source=HelmRepository/my-charts \
  --chart=my-app \
  --export > helmrelease.yaml
```

---

## flux debug helmrelease

Debug a failing HelmRelease — shows detailed status, last attempted values,
and error messages.

```bash
flux debug helmrelease <name> [-n <namespace>]
```

---

## Managing HelmReleases

```bash
# List
flux get helmreleases [-A] [-w] [--no-header]

# Filter
flux get helmreleases --status-selector ready=false
flux get helmreleases -l app=frontend

# Force reconciliation
flux reconcile helmrelease <name>

# Suspend/resume
flux suspend helmrelease <name>
flux suspend helmrelease --all
flux resume helmrelease <name>

# Export
flux export helmrelease <name>
flux export helmrelease --all > all-hr.yaml

# Delete
flux delete helmrelease <name>
```
