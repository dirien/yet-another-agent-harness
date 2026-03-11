# Pulumi State Management

Complete reference for Pulumi state operations, recovery patterns, and best practices.

## Understanding Pulumi State

Pulumi state tracks the mapping between your program's resources and the actual cloud resources. State is stored in a backend (Pulumi Cloud, S3, Azure Blob, GCS, or local filesystem).

**Key concepts:**
- **URN (Uniform Resource Name)**: Unique identifier for each resource
- **State checkpoint**: Snapshot of all resources at a point in time
- **Pending operations**: In-flight operations that haven't completed

## State Subcommands

### pulumi state delete

Remove a resource from state without deleting the actual cloud resource.

```bash
# Delete resource from state (keeps cloud resource)
pulumi state delete 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket'

# Delete with confirmation skip
pulumi state delete 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket' --yes

# Force delete even if protected
pulumi state delete 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket' --force

# Delete from specific stack
pulumi state delete --stack prod 'urn:pulumi:prod::myproject::aws:s3/bucket:Bucket::my-bucket'
```

**Use cases:**
- Resource was deleted manually in cloud console
- Need to "adopt" resource with different configuration
- Removing orphaned state entries
- Migrating resources between stacks

### pulumi state move

Move resources between stacks.

```bash
# Move resource to another stack
pulumi state move \
  --source dev \
  --dest prod \
  'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket'

# Move multiple resources
pulumi state move \
  --source dev \
  --dest prod \
  'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::bucket-1' \
  'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::bucket-2'

# Move with confirmation skip
pulumi state move --source dev --dest prod 'urn:...' --yes
```

**Important considerations:**
- Both stacks must use same backend
- Target stack must exist
- Resource must be removed from source stack's code
- Resource must be added to target stack's code

### pulumi state rename

Rename a resource in state.

```bash
# Rename resource
pulumi state rename \
  'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::old-name' \
  'new-name'

# Rename in specific stack
pulumi state rename --stack dev 'urn:pulumi:dev::...' 'new-name'
```

**Use cases:**
- Refactoring code that changes resource names
- Correcting naming mistakes without recreating resources

### pulumi state protect / unprotect

Prevent accidental deletion of resources.

```bash
# Protect resource from deletion
pulumi state protect 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::critical-bucket'

# Unprotect resource
pulumi state unprotect 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::critical-bucket'

# Check protection status (in state export)
pulumi stack export | jq '.deployment.resources[] | select(.protect == true)'
```

**Protected resources:**
- Cannot be deleted via `pulumi destroy`
- Cannot be replaced (must be unprotected first)
- Visible in `pulumi preview` with protection indicator

### pulumi state taint / untaint

Mark resources for replacement.

```bash
# Mark resource for replacement on next update
pulumi state taint 'urn:pulumi:dev::myproject::aws:ec2/instance:Instance::web-server'

# Remove taint (cancel scheduled replacement)
pulumi state untaint 'urn:pulumi:dev::myproject::aws:ec2/instance:Instance::web-server'
```

**Use cases:**
- Force recreation of misbehaving resource
- Trigger replacement without code changes
- Recovery from corrupted resource state

### pulumi state repair

Repair invalid state checkpoints.

```bash
# Repair state
pulumi state repair

# Repair specific stack
pulumi state repair --stack dev
```

**What it repairs:**
- Removes duplicate resources
- Clears orphaned pending operations
- Fixes version mismatches

## Stack Export/Import

### Full State Backup

```bash
# Export current stack state
pulumi stack export --file state-backup.json

# Export specific stack
pulumi stack export --stack prod --file prod-backup.json

# Export with version
pulumi stack export --version 42 --file state-v42.json
```

### State Restoration

```bash
# Import state from backup
pulumi stack import --file state-backup.json

# Import to specific stack
pulumi stack import --stack dev --file state-backup.json
```

### State Migration Between Backends

```bash
# 1. Export from source backend
pulumi login https://api.pulumi.com
pulumi stack export --stack myorg/myproject/prod --file state.json

# 2. Login to target backend
pulumi login s3://my-state-bucket

# 3. Create stack and import
pulumi stack init prod
pulumi stack import --file state.json
```

## Common Recovery Patterns

### Resource Deleted Outside Pulumi

When someone deletes a resource directly in the cloud console:

```bash
# Option 1: Refresh to sync state (removes from state)
pulumi refresh --yes

# Option 2: Delete from state manually (if refresh fails)
pulumi state delete 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::deleted-bucket'

# Then run up to recreate
pulumi up --yes
```

### Stuck Pending Operations

When an update was interrupted and left pending operations:

```bash
# View pending operations
pulumi stack export | jq '.deployment.pending_operations'

# Option 1: Refresh with clear pending
pulumi refresh --clear-pending-creates --yes

# Option 2: Cancel and retry
pulumi cancel --yes
pulumi refresh --yes
pulumi up --yes

# Option 3: Manual state repair
pulumi state repair
pulumi refresh --yes
```

### Resource Exists in Cloud But Not in State

When you need to adopt existing resources:

```bash
# Option 1: Import the resource
pulumi import aws:s3/bucket:Bucket my-bucket my-existing-bucket-id

# Option 2: Refresh to detect (if resource is in code)
pulumi refresh --yes
```

### State Corruption Recovery

When state becomes corrupted:

```bash
# 1. Try automatic repair
pulumi state repair

# 2. If that fails, check for backups
pulumi stack history
pulumi stack export --version <previous-good-version> --file backup.json

# 3. Import the backup
pulumi stack import --file backup.json

# 4. Refresh to sync with reality
pulumi refresh --yes
```

### Moving Resources Between Stacks

Complete workflow for refactoring resources across stacks:

```bash
# 1. Export source stack state
pulumi stack export --stack source-stack --file source.json

# 2. Identify resource URN
cat source.json | jq '.deployment.resources[] | select(.urn | contains("my-resource"))'

# 3. Move resource
pulumi state move \
  --source source-stack \
  --dest target-stack \
  'urn:pulumi:source-stack::myproject::aws:s3/bucket:Bucket::my-bucket'

# 4. Update code in both stacks
# - Remove resource from source stack code
# - Add resource to target stack code (with same URN name)

# 5. Verify both stacks
pulumi preview --stack source-stack
pulumi preview --stack target-stack
```

### Renaming Resources Without Recreation

```bash
# 1. Get current URN
pulumi stack export | jq '.deployment.resources[].urn' | grep my-resource

# 2. Rename in state
pulumi state rename 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::old-name' 'new-name'

# 3. Update code to match new name
# Change: new aws.s3.Bucket("old-name", {...})
# To:     new aws.s3.Bucket("new-name", {...})

# 4. Preview to verify no changes
pulumi preview
```

## State Inspection

### View State Contents

```bash
# View all resources
pulumi stack export | jq '.deployment.resources[]'

# View specific resource
pulumi stack export | jq '.deployment.resources[] | select(.urn | contains("my-bucket"))'

# View resource outputs
pulumi stack export | jq '.deployment.resources[] | select(.urn | contains("my-bucket")) | .outputs'

# View resource inputs
pulumi stack export | jq '.deployment.resources[] | select(.urn | contains("my-bucket")) | .inputs'

# Count resources by type
pulumi stack export | jq '[.deployment.resources[].type] | group_by(.) | map({type: .[0], count: length})'

# Find protected resources
pulumi stack export | jq '.deployment.resources[] | select(.protect == true) | .urn'

# Find pending operations
pulumi stack export | jq '.deployment.pending_operations'
```

### URN Structure

URNs follow this format:
```
urn:pulumi:<stack>::<project>::<type>::<name>
```

Example:
```
urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::my-bucket
         │        │              │                    │
       stack   project    resource type           resource name
```

## Best Practices

### State Backup Strategy

```bash
# Before major operations, always backup
pulumi stack export --file "backup-$(date +%Y%m%d-%H%M%S).json"

# Automate in CI/CD
pulumi stack export --file "state-${GITHUB_SHA}.json"
aws s3 cp "state-${GITHUB_SHA}.json" "s3://my-backups/pulumi/"
```

### Safe State Modifications

1. **Always backup before state operations**
2. **Use `--yes` carefully** - prefer interactive confirmation
3. **Preview after state changes** - verify no unexpected diffs
4. **Test in non-production first**

### Protecting Critical Resources

```bash
# Protect production databases
pulumi state protect 'urn:pulumi:prod::myproject::aws:rds/instance:Instance::prod-db'

# Protect stateful resources
pulumi state protect 'urn:pulumi:prod::myproject::aws:efs/fileSystem:FileSystem::shared-storage'
```

### State Locking

Pulumi Cloud and most backends support state locking:
- Prevents concurrent modifications
- Automatic with Pulumi Cloud
- Configure for S3: enable DynamoDB locking

```bash
# S3 backend with locking
pulumi login 's3://my-bucket?region=us-west-2&awssdk=v2'
```

## Troubleshooting

### "resource already exists"

```bash
# Resource exists in cloud but not state - import it
pulumi import aws:s3/bucket:Bucket my-bucket existing-bucket-id

# Or delete from cloud and let Pulumi recreate
aws s3 rb s3://existing-bucket-id --force
pulumi up --yes
```

### "resource not found" during refresh

```bash
# Resource was deleted outside Pulumi
pulumi state delete 'urn:pulumi:dev::myproject::aws:s3/bucket:Bucket::missing-bucket'
```

### "pending operation" errors

```bash
# Clear pending operations
pulumi refresh --clear-pending-creates --yes

# Or cancel and retry
pulumi cancel --yes
pulumi state repair
```

### State file too large

```bash
# Check state size
pulumi stack export | wc -c

# Consider splitting into multiple stacks
# Use stack references for dependencies
```
