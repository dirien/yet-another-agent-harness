# Build, Diff, and Patch Commands

These commands work offline — no cluster access required.

---

## build instance

Generate Kubernetes manifests from a FluxInstance definition.

```bash
flux-operator build instance -f <path-to-fluxinstance.yaml>
```

### Flags

| Flag | Description |
|------|-------------|
| `-f, --file` | Path to FluxInstance YAML file (required) |
| `--distribution-artifact` | OCI artifact URL for Flux distribution; overrides what's in the FluxInstance spec |

### Usage

```bash
# Basic build
flux-operator build instance -f flux-instance.yaml

# Override the distribution artifact
flux-operator build instance -f flux-instance.yaml \
  --distribution-artifact oci://ghcr.io/controlplaneio-fluxcd/flux-operator-manifests:v2.5
```

The output is the rendered Kubernetes manifests that the operator would apply.

---

## build rset

Generate Kubernetes manifests from a ResourceSet definition.

```bash
flux-operator build rset -f <path-to-resourceset.yaml>
```

### Flags

| Flag | Description |
|------|-------------|
| `-f, --file` | Path to ResourceSet YAML file (required) |
| `--inputs-from` | Path to a YAML file containing ResourceSet inputs |
| `--inputs-from-provider` | Path to a ResourceSetInputProvider YAML of static type |

### Usage

```bash
# Basic build
flux-operator build rset -f resourceset.yaml

# With inputs from a file
flux-operator build rset -f resourceset.yaml --inputs-from inputs.yaml

# With inputs from a static provider definition
flux-operator build rset -f resourceset.yaml \
  --inputs-from-provider provider.yaml
```

---

## diff yaml

Compare two YAML files and produce an RFC 6902 JSON patch describing the differences.

```bash
flux-operator diff yaml <source> <target>
```

### Arguments

Both `<source>` and `<target>` accept:
- Local file paths
- Remote URLs: GitHub raw files, GitLab files, GitHub Gists, OCI artifact URLs

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | `json-patch-yaml` | Output format: `json-patch-yaml` or `json-patch` |

### Behavior

- Ignores metadata and status fields — focuses on semantic differences in specs
- Produces RFC 6902 compliant patches

### Usage

```bash
# Compare two local files
flux-operator diff yaml old-instance.yaml new-instance.yaml

# Compare local against remote
flux-operator diff yaml flux-instance.yaml \
  https://raw.githubusercontent.com/org/repo/main/flux-instance.yaml

# Output as JSON
flux-operator diff yaml old.yaml new.yaml -o json-patch
```

---

## patch instance

Generate kustomize patches for upgrading Flux controllers. Modifies the FluxInstance
YAML in-place, appending patches to `.spec.kustomize.patches`.

```bash
flux-operator patch instance -f <path-to-fluxinstance.yaml>
```

### What it does

1. Fetches CRD schemas from both the current and target Flux versions
2. Computes JSON patches for changed CRDs
3. Generates Deployment image patches for controller version updates
4. Appends all patches to `.spec.kustomize.patches` in the FluxInstance YAML
5. Automatically replaces any previously generated patches (safe to run repeatedly)

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-f, --filename` | — | Path to FluxInstance YAML (required); use `-` for stdin |
| `-v, --version` | `main` | Target Flux version; accepts `main`, `v2.<minor>`, or just `<minor>` |
| `-r, --registry` | — | Override the container registry for Flux images |
| `-c, --components` | from spec | Comma-separated list of controllers to patch (defaults to `.spec.components`) |

### Usage

```bash
# Patch to latest main
flux-operator patch instance -f flux-instance.yaml

# Patch to a specific version
flux-operator patch instance -f flux-instance.yaml -v v2.5

# Patch specific components with custom registry
flux-operator patch instance -f flux-instance.yaml \
  -v v2.5 \
  -r my-registry.example.com/flux \
  -c source-controller,kustomize-controller

# Read from stdin
cat flux-instance.yaml | flux-operator patch instance -f -
```
