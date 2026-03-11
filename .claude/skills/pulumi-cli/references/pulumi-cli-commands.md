# Pulumi CLI Commands Reference

Complete reference for Pulumi CLI commands.

## Project Initialization

### pulumi new

Create a new Pulumi project from a template.

```bash
# Interactive mode - prompts for project name, description, stack
pulumi new

# Create from specific template
pulumi new typescript
pulumi new python
pulumi new go

# Cloud-specific templates
pulumi new aws-typescript
pulumi new azure-python
pulumi new gcp-go
pulumi new kubernetes-typescript

# Non-interactive with all options specified
pulumi new typescript \
  --name my-project \
  --description "My infrastructure project" \
  --stack dev \
  --yes

# Generate in specific directory
pulumi new typescript --dir ./my-project

# List available templates
pulumi new --list
pulumi new --list | grep aws

# Use template from URL
pulumi new https://github.com/pulumi/templates/tree/master/aws-typescript
```

**Common flags:**
- `--name` - Project name (default: directory name)
- `--description` - Project description
- `--stack` - Initial stack name
- `--yes` / `-y` - Skip confirmation prompts
- `--dir` - Target directory (default: current)
- `--force` / `-f` - Overwrite existing files
- `--generate-only` - Generate project files only, don't create stack
- `--secrets-provider` - Secrets provider (default: `passphrase`, options: `awskms`, `azurekeyvault`, `gcpkms`, `hashivault`)

## Deployment Commands

### pulumi preview

Preview changes before deployment.

```bash
# Basic preview
pulumi preview

# Preview with specific stack
pulumi preview --stack dev

# JSON output for programmatic use
pulumi preview --json

# Show detailed diff of changes
pulumi preview --diff

# Preview with refresh (sync state first)
pulumi preview --refresh

# Preview specific resource
pulumi preview --target urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket

# Preview excluding specific resources
pulumi preview --target-dependents

# Suppress progress messages
pulumi preview --suppress-progress

# Set parallelism (default: infinity)
pulumi preview --parallel 10
```

**Common flags:**
- `--stack` / `-s` - Stack name
- `--config` / `-c` - Config values (key=value)
- `--diff` - Show detailed diff
- `--json` - JSON output
- `--refresh` - Refresh state before preview
- `--target` - Target specific URNs
- `--target-dependents` - Include dependents of targets
- `--replace` - Mark resources for replacement
- `--suppress-outputs` - Suppress stack outputs
- `--suppress-progress` - Suppress progress messages
- `--parallel` / `-p` - Parallelism limit

### pulumi up

Deploy infrastructure changes.

```bash
# Interactive deployment (shows preview, asks for confirmation)
pulumi up

# Non-interactive deployment
pulumi up --yes

# Deploy specific stack
pulumi up --stack prod --yes

# Skip preview entirely
pulumi up --skip-preview --yes

# Deploy with refresh
pulumi up --refresh --yes

# Target specific resources
pulumi up --target urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket --yes

# Replace specific resources (force recreation)
pulumi up --replace urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket --yes

# Continue on errors (skip failed resources)
pulumi up --continue-on-error --yes

# Show secrets in output
pulumi up --show-secrets --yes

# Set config at deploy time
pulumi up --config aws:region=us-west-2 --yes

# Policy packs enforcement
pulumi up --policy-pack /path/to/policy-pack --yes
```

**Common flags:**
- `--yes` / `-y` - Automatic approval
- `--skip-preview` - Skip preview step
- `--refresh` - Refresh state first
- `--target` - Deploy only specific URNs
- `--target-dependents` - Include dependent resources
- `--replace` - Force replacement of URNs
- `--continue-on-error` - Continue on resource failures
- `--show-secrets` - Display secrets in output
- `--message` / `-m` - Update message for history
- `--policy-pack` - Path to policy pack
- `--parallel` / `-p` - Parallelism limit
- `--suppress-outputs` - Suppress stack outputs

### pulumi destroy

Tear down infrastructure.

```bash
# Interactive destroy (shows preview, asks confirmation)
pulumi destroy

# Non-interactive destroy
pulumi destroy --yes

# Destroy specific stack
pulumi destroy --stack dev --yes

# Destroy specific resources only
pulumi destroy --target urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket --yes

# Destroy with dependents
pulumi destroy --target urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket --target-dependents --yes

# Skip refresh
pulumi destroy --skip-preview --yes

# Continue even if some resources fail to delete
pulumi destroy --continue-on-error --yes
```

**Common flags:**
- `--yes` / `-y` - Automatic approval
- `--skip-preview` - Skip preview
- `--target` - Target specific URNs
- `--target-dependents` - Include dependents
- `--continue-on-error` - Continue on failures
- `--parallel` / `-p` - Parallelism limit
- `--message` / `-m` - Destroy message for history

### pulumi refresh

Synchronize state with actual cloud resources.

```bash
# Interactive refresh
pulumi refresh

# Non-interactive refresh
pulumi refresh --yes

# Refresh specific stack
pulumi refresh --stack prod --yes

# Skip preview
pulumi refresh --skip-preview --yes

# Target specific resources
pulumi refresh --target urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket --yes

# Clear pending operations
pulumi refresh --clear-pending-creates --yes
```

**Common flags:**
- `--yes` / `-y` - Automatic approval
- `--skip-preview` - Skip preview step
- `--target` - Target specific URNs
- `--clear-pending-creates` - Clear stuck pending creates
- `--parallel` / `-p` - Parallelism limit
- `--message` / `-m` - Refresh message

## Stack Management

### pulumi stack

Manage stacks.

```bash
# Show current stack info
pulumi stack

# List all stacks
pulumi stack ls
pulumi stack ls --all  # Include stacks from all organizations

# Initialize new stack
pulumi stack init dev
pulumi stack init myorg/myproject/prod  # Fully qualified name

# Initialize with specific secrets provider
pulumi stack init dev --secrets-provider awskms://alias/my-key

# Select/switch stack
pulumi stack select dev
pulumi stack select myorg/myproject/prod

# Remove stack (must be empty)
pulumi stack rm dev
pulumi stack rm dev --yes
pulumi stack rm dev --force  # Remove even if not empty (dangerous!)

# Rename stack
pulumi stack rename staging

# Show stack outputs
pulumi stack output
pulumi stack output bucketName  # Specific output
pulumi stack output --json
pulumi stack output --show-secrets

# Show stack history
pulumi stack history
pulumi stack history --full-dates
pulumi stack history --page-size 20

# View update details
pulumi stack history --page-size 1 --json

# Tag stack
pulumi stack tag set environment production
pulumi stack tag ls
pulumi stack tag rm environment

# Change secrets provider
pulumi stack change-secrets-provider awskms://alias/new-key
```

### pulumi stack export/import

Backup and restore stack state.

```bash
# Export state to file
pulumi stack export --file state.json

# Export specific stack
pulumi stack export --stack prod --file prod-state.json

# Import state from file
pulumi stack import --file state.json

# Migrate state between backends
pulumi stack export --stack dev --file backup.json
pulumi login s3://my-bucket
pulumi stack init dev
pulumi stack import --file backup.json
```

## Configuration

### pulumi config

Manage stack configuration (basic commands, see language-specific skills for ESC).

```bash
# Set configuration value
pulumi config set aws:region us-west-2

# Set secret (encrypted)
pulumi config set dbPassword mysecret --secret

# Set from file
pulumi config set-all --path --file config.json

# Get configuration value
pulumi config get aws:region

# List all config
pulumi config

# Remove configuration
pulumi config rm aws:region

# Refresh config from ESC
pulumi config refresh

# Add ESC environment
pulumi config env add myorg/myproject-dev
pulumi config env rm myorg/myproject-dev

# Copy config between stacks
pulumi config cp --stack source --dest target
```

**Common flags:**
- `--secret` - Encrypt the value
- `--path` - Treat key as path (dotted notation)
- `--plaintext` - Print secrets in plaintext
- `--json` - JSON output

## Resource Import

### pulumi import

Import existing cloud resources into Pulumi.

```bash
# Import single resource
pulumi import aws:s3/bucket:Bucket my-bucket my-existing-bucket-id

# Import with parent
pulumi import aws:s3/bucket:Bucket my-bucket my-existing-bucket-id \
  --parent urn:pulumi:dev::myproject::custom:module:MyModule::parent

# Import without generating code
pulumi import aws:s3/bucket:Bucket my-bucket my-existing-bucket-id --skip-generate

# Import to specific file
pulumi import aws:s3/bucket:Bucket my-bucket my-bucket-id --out imports.ts

# Import from JSON file (bulk import)
pulumi import --file resources.json

# Protect imported resources
pulumi import aws:s3/bucket:Bucket my-bucket my-bucket-id --protect
```

**resources.json format:**
```json
{
  "resources": [
    {
      "type": "aws:s3/bucket:Bucket",
      "name": "my-bucket",
      "id": "my-existing-bucket-id"
    },
    {
      "type": "aws:ec2/instance:Instance",
      "name": "my-server",
      "id": "i-1234567890abcdef0"
    }
  ]
}
```

## Utility Commands

### pulumi login/logout

Manage backend authentication.

```bash
# Login to Pulumi Cloud (default)
pulumi login

# Login to specific organization
pulumi login --cloud-url https://api.pulumi.com

# Login to self-hosted backend
pulumi login s3://my-bucket
pulumi login azblob://my-container
pulumi login gs://my-bucket
pulumi login file://~/.pulumi-state

# Logout
pulumi logout
pulumi logout --all  # Logout from all backends
```

### pulumi whoami

Show current user.

```bash
pulumi whoami
pulumi whoami --verbose  # Include backend URL and organizations
```

### pulumi about

Show environment information.

```bash
pulumi about
pulumi about --json
```

### pulumi cancel

Cancel a running update.

```bash
# Cancel current update
pulumi cancel

# Cancel specific stack's update
pulumi cancel --stack dev

# Force cancel (dangerous - may corrupt state)
pulumi cancel --yes
```

### pulumi console

Open Pulumi Cloud console.

```bash
# Open current stack in browser
pulumi console

# Open specific stack
pulumi console --stack prod
```

### pulumi logs

View aggregated logs from cloud resources.

```bash
# View logs
pulumi logs

# Follow logs in real-time
pulumi logs --follow

# Filter by time
pulumi logs --since 1h
pulumi logs --since 2024-01-01T00:00:00Z

# Filter by resource
pulumi logs --resource my-function

# JSON output
pulumi logs --json
```

### pulumi watch (deprecated)

Continuous deployment mode (watches for file changes).

```bash
# Start watch mode
pulumi watch

# With specific path
pulumi watch --path ./src
```

## Policy Commands

### pulumi policy

Manage Policy Packs.

```bash
# Create new policy pack
pulumi policy new aws-typescript

# Publish policy pack
pulumi policy publish myorg

# Enable policy pack for organization
pulumi policy enable myorg/my-policy-pack latest

# Disable policy pack
pulumi policy disable myorg/my-policy-pack

# List policy packs
pulumi policy ls myorg

# Validate policy pack
pulumi policy validate-config myorg/my-policy-pack
```

## Package Commands

### pulumi package

Manage Pulumi packages and components.

```bash
# Add package from registry
pulumi package add aws
pulumi package add azure-native

# Add package from git
pulumi package add github.com/myorg/my-component
pulumi package add github.com/myorg/my-component@v1.0.0

# Add local package
pulumi package add /path/to/local/package

# Generate SDK from schema
pulumi package gen-sdk ./schema.json --language typescript

# Get schema for installed package
pulumi package get-schema aws
```

## Schema Commands

### pulumi schema

Work with Pulumi schemas.

```bash
# Check schema
pulumi schema check ./schema.json
```

## Version and Updates

```bash
# Show version
pulumi version

# Check for updates
pulumi version --check-updates

# Install specific version (via package managers)
brew install pulumi/tap/pulumi
choco install pulumi
```

## Global Flags

These flags work with most commands:

| Flag | Description |
|------|-------------|
| `--cwd` / `-C` | Run as if started in specified directory |
| `--disable-integrity-checking` | Disable integrity checking of state |
| `--emoji` | Enable emoji in output |
| `--help` / `-h` | Help for command |
| `--logflow` | Flow log output to parent process |
| `--logtostderr` | Log to stderr instead of files |
| `--non-interactive` | Disable interactive mode |
| `--profiling` | Emit CPU/memory profiles to specified directory |
| `--tracing` | Emit tracing to specified endpoint |
| `--verbose` / `-v` | Verbosity level (0-9) |

## Output Formats

Many commands support different output formats:

```bash
# JSON output
pulumi stack output --json
pulumi preview --json
pulumi stack ls --json

# Suppress specific output types
pulumi up --suppress-outputs    # Hide stack outputs
pulumi up --suppress-progress   # Hide progress messages
pulumi up --suppress-permalink  # Hide permalink to update
```
