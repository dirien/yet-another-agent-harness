# Pulumi Environment Variables

Complete reference for environment variables used with Pulumi CLI, especially for CI/CD and automation.

## Authentication

### PULUMI_ACCESS_TOKEN

Pulumi Cloud access token for authentication.

```bash
export PULUMI_ACCESS_TOKEN=pul-xxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**Usage:**
- Required for non-interactive Pulumi Cloud authentication
- Generate at https://app.pulumi.com/account/tokens
- Use organization tokens for CI/CD systems

**Best practices:**
- Store in CI/CD secrets management
- Use short-lived tokens when possible
- Rotate regularly

### PULUMI_BACKEND_URL

Override the default backend URL.

```bash
# Pulumi Cloud (default)
export PULUMI_BACKEND_URL=https://api.pulumi.com

# Self-hosted Pulumi Cloud
export PULUMI_BACKEND_URL=https://pulumi.mycompany.com

# S3 backend
export PULUMI_BACKEND_URL=s3://my-bucket

# Azure Blob backend
export PULUMI_BACKEND_URL=azblob://my-container

# GCS backend
export PULUMI_BACKEND_URL=gs://my-bucket

# Local filesystem
export PULUMI_BACKEND_URL=file://~/.pulumi-state
```

## Stack and Project

### PULUMI_STACK

Default stack for all commands.

```bash
export PULUMI_STACK=dev

# Now these are equivalent:
pulumi up --stack dev
pulumi up
```

### PULUMI_CONFIG_PASSPHRASE

Passphrase for encrypting stack secrets (when using passphrase secrets provider).

```bash
export PULUMI_CONFIG_PASSPHRASE=my-secret-passphrase
```

**Important:**
- Required for non-interactive operations with passphrase encryption
- Consider using cloud KMS instead for production

### PULUMI_CONFIG_PASSPHRASE_FILE

Read passphrase from file instead of environment variable.

```bash
export PULUMI_CONFIG_PASSPHRASE_FILE=/path/to/passphrase-file
```

## CI/CD Automation

### PULUMI_CI

Indicate running in CI environment.

```bash
export PULUMI_CI=true
```

**Effects:**
- Enables non-interactive mode
- Adjusts output formatting for CI logs
- Disables browser opening for `pulumi login`

### PULUMI_SKIP_UPDATE_CHECK

Disable update checks.

```bash
export PULUMI_SKIP_UPDATE_CHECK=true
```

**Use in CI/CD to:**
- Speed up execution
- Avoid network calls
- Prevent update prompts

### PULUMI_SKIP_CONFIRMATIONS

Skip all confirmation prompts (equivalent to `--yes`).

```bash
export PULUMI_SKIP_CONFIRMATIONS=true
```

**Warning:** Use carefully - skips safety confirmations.

## Performance

### PULUMI_PARALLEL

Control parallelism for resource operations.

```bash
# Limit to 10 concurrent operations
export PULUMI_PARALLEL=10

# Unlimited (default)
export PULUMI_PARALLEL=0

# Sequential (for debugging)
export PULUMI_PARALLEL=1
```

**When to adjust:**
- Rate limiting from cloud providers
- Debugging dependency issues
- Resource-constrained environments

### PULUMI_EXPERIMENTAL

Enable experimental features.

```bash
export PULUMI_EXPERIMENTAL=true
```

## Debugging and Logging

### PULUMI_DEBUG_COMMANDS

Enable debug output for CLI commands.

```bash
export PULUMI_DEBUG_COMMANDS=true
```

### PULUMI_DEBUG_GRPC

Debug gRPC communication with providers.

```bash
export PULUMI_DEBUG_GRPC=/path/to/grpc-log
```

### PULUMI_LOG_LEVEL

Control logging verbosity.

```bash
# Error only
export PULUMI_LOG_LEVEL=error

# Warnings and errors
export PULUMI_LOG_LEVEL=warning

# Info level (default)
export PULUMI_LOG_LEVEL=info

# Debug level
export PULUMI_LOG_LEVEL=debug
```

### PULUMI_DEBUG_PROVIDERS

Comma-separated list of providers to debug.

```bash
export PULUMI_DEBUG_PROVIDERS=aws,kubernetes
```

### PULUMI_ENABLE_LEGACY_DIFF

Enable legacy diff behavior for troubleshooting.

```bash
export PULUMI_ENABLE_LEGACY_DIFF=true
```

## Plugin Management

### PULUMI_SKIP_PROVIDER_INSTALL

Skip automatic provider plugin installation.

```bash
export PULUMI_SKIP_PROVIDER_INSTALL=true
```

### PULUMI_PLUGIN_PATH

Additional paths to search for plugins.

```bash
export PULUMI_PLUGIN_PATH=/custom/plugin/path
```

### PULUMI_PREFER_YARN

Prefer Yarn over npm for Node.js plugins.

```bash
export PULUMI_PREFER_YARN=true
```

## Cloud Provider Integration

### AWS

```bash
# Standard AWS credentials
export AWS_ACCESS_KEY_ID=AKIA...
export AWS_SECRET_ACCESS_KEY=...
export AWS_SESSION_TOKEN=...  # For temporary credentials
export AWS_REGION=us-west-2
export AWS_DEFAULT_REGION=us-west-2

# AWS profile
export AWS_PROFILE=my-profile
```

### Azure

```bash
# Service principal authentication
export ARM_CLIENT_ID=...
export ARM_CLIENT_SECRET=...
export ARM_TENANT_ID=...
export ARM_SUBSCRIPTION_ID=...

# Managed identity
export ARM_USE_MSI=true
export ARM_MSI_ENDPOINT=...
```

### Google Cloud

```bash
# Service account key file
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json

# Project and region
export GOOGLE_PROJECT=my-project
export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
```

### Kubernetes

```bash
# Kubeconfig location
export KUBECONFIG=/path/to/kubeconfig

# In-cluster configuration
export KUBERNETES_SERVICE_HOST=...
export KUBERNETES_SERVICE_PORT=...
```

## CI/CD Platform Examples

### GitHub Actions

```yaml
name: Pulumi
on: push

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: pulumi/actions@v5
        with:
          command: up
          stack-name: prod
        env:
          PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: us-west-2
```

### GitLab CI

```yaml
deploy:
  image: pulumi/pulumi:latest
  variables:
    PULUMI_ACCESS_TOKEN: $PULUMI_ACCESS_TOKEN
    PULUMI_STACK: prod
    PULUMI_CI: "true"
    PULUMI_SKIP_UPDATE_CHECK: "true"
  script:
    - pulumi login
    - pulumi stack select $PULUMI_STACK
    - pulumi up --yes
```

### Jenkins

```groovy
pipeline {
    agent any
    environment {
        PULUMI_ACCESS_TOKEN = credentials('pulumi-token')
        AWS_ACCESS_KEY_ID = credentials('aws-access-key')
        AWS_SECRET_ACCESS_KEY = credentials('aws-secret-key')
        PULUMI_CI = 'true'
        PULUMI_SKIP_UPDATE_CHECK = 'true'
    }
    stages {
        stage('Deploy') {
            steps {
                sh 'pulumi login'
                sh 'pulumi stack select prod'
                sh 'pulumi up --yes'
            }
        }
    }
}
```

### CircleCI

```yaml
version: 2.1

orbs:
  pulumi: pulumi/pulumi@2.1.0

jobs:
  deploy:
    docker:
      - image: pulumi/pulumi:latest
    environment:
      PULUMI_CI: "true"
      PULUMI_SKIP_UPDATE_CHECK: "true"
    steps:
      - checkout
      - run:
          name: Deploy
          command: |
            pulumi login
            pulumi stack select prod
            pulumi up --yes
```

### Azure DevOps

```yaml
trigger:
  - main

pool:
  vmImage: 'ubuntu-latest'

variables:
  - name: PULUMI_CI
    value: 'true'
  - name: PULUMI_SKIP_UPDATE_CHECK
    value: 'true'

steps:
  - task: Pulumi@1
    inputs:
      command: 'up'
      stack: 'prod'
      args: '--yes'
    env:
      PULUMI_ACCESS_TOKEN: $(PULUMI_ACCESS_TOKEN)
```

## Automation API

When using Pulumi Automation API, environment variables work the same way. You can also pass them programmatically:

```typescript
import * as automation from "@pulumi/pulumi/automation";

const stack = await automation.LocalWorkspace.createOrSelectStack({
    stackName: "dev",
    projectName: "my-project",
    program: async () => { /* ... */ },
}, {
    envVars: {
        AWS_REGION: "us-west-2",
        PULUMI_CONFIG_PASSPHRASE: process.env.PASSPHRASE,
    },
});
```

## Complete CI/CD Environment

Recommended environment variables for CI/CD automation:

```bash
# Authentication
export PULUMI_ACCESS_TOKEN=pul-xxx

# Automation behavior
export PULUMI_CI=true
export PULUMI_SKIP_UPDATE_CHECK=true

# Optional: secrets passphrase (if not using cloud KMS)
export PULUMI_CONFIG_PASSPHRASE=xxx

# Optional: performance tuning
export PULUMI_PARALLEL=10

# Cloud provider credentials (choose one)
# AWS
export AWS_ACCESS_KEY_ID=xxx
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_REGION=us-west-2

# OR use OIDC via Pulumi ESC (recommended)
# ESC environments inject credentials automatically
```

## Reference Table

| Variable | Description | Default |
|----------|-------------|---------|
| `PULUMI_ACCESS_TOKEN` | Pulumi Cloud authentication token | None |
| `PULUMI_BACKEND_URL` | Backend URL | `https://api.pulumi.com` |
| `PULUMI_STACK` | Default stack name | None |
| `PULUMI_CONFIG_PASSPHRASE` | Secrets encryption passphrase | None |
| `PULUMI_CI` | CI mode indicator | `false` |
| `PULUMI_SKIP_UPDATE_CHECK` | Disable update checks | `false` |
| `PULUMI_SKIP_CONFIRMATIONS` | Skip all prompts | `false` |
| `PULUMI_PARALLEL` | Operation parallelism | Unlimited |
| `PULUMI_DEBUG_COMMANDS` | Debug CLI commands | `false` |
| `PULUMI_LOG_LEVEL` | Logging verbosity | `info` |
