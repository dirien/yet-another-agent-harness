# Pulumi ESC Reference

## Overview

Pulumi ESC (Environments, Secrets, and Configuration) provides centralized secrets management and orchestration for infrastructure and applications.

## Pulumi ESC Commands

### Environment Management

```bash
# Create environment
pulumi env init <org>/<project>/<env>
pulumi env init myorg/myproject/dev

# List environments
pulumi env ls
pulumi env ls myorg

# Edit environment (opens editor)
pulumi env edit <org>/<env>

# View environment definition
pulumi env get <org>/<env> --show-secrets

# Delete environment
pulumi env rm <org>/<env>

# Clone environment
pulumi env clone <org>/<source> <org>/<target>
```

### Working with Values

```bash
# Set a value
pulumi env set <org>/<env> <path> <value>
pulumi env set myorg/dev pulumiConfig.aws:region us-west-2

# Set a secret
pulumi env set <org>/<env> <path> <value> --secret

# Get a specific value
pulumi env get <org>/<env> <path>
pulumi env get myorg/dev pulumiConfig.aws:region
```

### Running Commands

```bash
# Run command with environment
pulumi env run <org>/<env> -- <command>
pulumi env run myorg/aws-dev -- pulumi up
pulumi env run myorg/aws-dev -- aws s3 ls

# Open environment (resolve and display)
pulumi env open <org>/<env>
pulumi env open <org>/<env> --format json
pulumi env open <org>/<env> --format shell
pulumi env open <org>/<env> --format dotenv
```

### Versioning

```bash
# List versions
pulumi env version ls <org>/<env>

# Tag a version
pulumi env version tag <org>/<env> <tag>
pulumi env version tag myorg/prod stable

# Rollback to version
pulumi env version rollback <org>/<env> <version>

# Diff versions
pulumi env diff <org>/<env>@<v1> <org>/<env>@<v2>
```

## ESC Environment Syntax

### Basic Structure

```yaml
imports:
  - base-environment    # Import from another environment

values:
  # Static values
  staticConfig:
    region: us-west-2
    environment: production

  # Secrets (encrypted at rest)
  secrets:
    apiKey:
      fn::secret: "my-secret-value"

  # Pulumi stack configuration
  pulumiConfig:
    aws:region: ${staticConfig.region}
    myapp:apiKey: ${secrets.apiKey}

  # Environment variables
  environmentVariables:
    AWS_REGION: ${staticConfig.region}
    MY_API_KEY: ${secrets.apiKey}
```

### Dynamic Providers

#### AWS OIDC Login

```yaml
values:
  aws:
    login:
      fn::open::aws-login:
        oidc:
          roleArn: arn:aws:iam::123456789012:role/pulumi-oidc-role
          sessionName: pulumi-${context.pulumi.user.login}
          duration: 1h
```

#### AWS Secrets Manager

```yaml
values:
  secrets:
    fn::open::aws-secrets:
      region: us-west-2
      login: ${aws.login}
      get:
        dbPassword:
          secretId: prod/db/password
        apiKey:
          secretId: prod/api/key
          versionStage: AWSCURRENT
```

#### Azure OIDC Login

```yaml
values:
  azure:
    login:
      fn::open::azure-login:
        clientId: <app-client-id>
        tenantId: <tenant-id>
        subscriptionId: <subscription-id>
        oidc: true
```

#### GCP OIDC Login

```yaml
values:
  gcp:
    login:
      fn::open::gcp-login:
        project: my-project-id
        oidc:
          workloadPoolId: pulumi-pool
          providerId: pulumi-provider
          serviceAccount: pulumi@my-project.iam.gserviceaccount.com
```

#### HashiCorp Vault

```yaml
values:
  vault:
    login:
      fn::open::vault-login:
        address: https://vault.example.com
        jwt:
          role: pulumi-role

  secrets:
    fn::open::vault-secrets:
      login: ${vault.login}
      read:
        dbCreds:
          path: secret/data/database
```

#### 1Password

```yaml
values:
  secrets:
    fn::open::1password-secrets:
      login:
        serviceAccountToken:
          fn::secret: <service-account-token>
      get:
        apiKey:
          ref: op://Vault/API Key/credential
```

### Environment Composition

```yaml
# base.yaml
values:
  common:
    region: us-west-2
    tags:
      ManagedBy: Pulumi

# dev.yaml
imports:
  - base

values:
  pulumiConfig:
    aws:region: ${common.region}
    myapp:environment: dev
    myapp:instanceType: t3.small

# prod.yaml
imports:
  - base

values:
  pulumiConfig:
    aws:region: ${common.region}
    myapp:environment: prod
    myapp:instanceType: t3.large
```

## Integration with Pulumi IaC

### Linking ESC to Stacks

```bash
# Add ESC environment to stack
pulumi config env add <org>/<env>

# Remove ESC environment
pulumi config env rm <org>/<env>

# List linked environments
pulumi config env ls
```

### Accessing Config in Code

```typescript
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

// Values from pulumiConfig block are available
const region = config.require("aws:region");
const instanceType = config.require("myapp:instanceType");
const apiKey = config.requireSecret("myapp:apiKey");
```

### Environment Variables

Environment variables defined in ESC are automatically available to your Pulumi program and any commands run with `pulumi env run`.
