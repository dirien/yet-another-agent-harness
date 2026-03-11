# Pulumi TypeScript Reference

## Project Setup

### Package.json

```json
{
    "name": "my-pulumi-project",
    "main": "index.ts",
    "devDependencies": {
        "@types/node": "^20",
        "typescript": "^5"
    },
    "dependencies": {
        "@pulumi/pulumi": "^3",
        "@pulumi/aws": "^6"
    }
}
```

### tsconfig.json

```json
{
    "compilerOptions": {
        "target": "ES2020",
        "module": "commonjs",
        "lib": ["ES2020"],
        "strict": true,
        "moduleResolution": "node",
        "esModuleInterop": true,
        "skipLibCheck": true,
        "forceConsistentCasingInFileNames": true,
        "outDir": "./bin"
    },
    "include": ["**/*.ts"],
    "exclude": ["node_modules"]
}
```

### Pulumi.yaml

```yaml
name: my-project
description: My Pulumi project
runtime:
  name: nodejs
  options:
    typescript: true
```

## Working with Inputs and Outputs

### Output Transformations

```typescript
import * as pulumi from "@pulumi/pulumi";

// Single output transformation
const bucketUrl = bucket.bucketDomainName.apply(
    domain => `https://${domain}`
);

// Multiple outputs
const combined = pulumi.all([resource1.id, resource2.arn]).apply(
    ([id, arn]) => ({ id, arn })
);

// Interpolation
const message = pulumi.interpolate`Resource ${resource.id} created`;

// Conditional output
const endpoint = pulumi.output(resource.endpoint).apply(
    ep => ep || "default-endpoint"
);
```

### Output to Promise (for async operations)

```typescript
// Only use outside of resource declarations
export = async () => {
    const stackRef = new pulumi.StackReference("myorg/other/prod");

    // Safe to await in async export function
    const value = await stackRef.getOutputValue("someOutput");

    return { value };
};
```

## Common Patterns

### Resource Options

```typescript
const resource = new aws.s3.Bucket("bucket", {}, {
    // Explicit dependencies
    dependsOn: [otherResource],

    // Parent for component resources
    parent: this,

    // Protect from deletion
    protect: true,

    // Ignore changes to specific properties
    ignoreChanges: ["tags"],

    // Custom provider
    provider: customProvider,

    // Resource aliases for refactoring
    aliases: [{ name: "old-bucket-name" }],

    // Delete before replace
    deleteBeforeReplace: true,

    // Custom timeouts
    customTimeouts: {
        create: "30m",
        update: "30m",
        delete: "30m",
    },

    // Force replacement when specific properties change
    replaceOnChanges: ["tags.Environment"],

    // Skip delete when parent resource is deleted (useful for Kubernetes)
    deletedWith: parentResource,
});
```

### Resource Hooks

```typescript
// Define lifecycle hooks for custom logic
const resource = new aws.s3.Bucket("bucket", {}, {
    hooks: {
        beforeCreate: async (args) => {
            console.log(`Creating resource: ${args.urn}`);
        },
        afterCreate: async (args) => {
            console.log(`Created with outputs: ${JSON.stringify(args.outputs)}`);
        },
    },
});
```

> **Note:** Hooks require `--run-program` flag when running `pulumi destroy`. Not supported in YAML.

### Provider Configuration

```typescript
// Explicit provider for multi-region
const usEast1 = new aws.Provider("us-east-1", {
    region: "us-east-1",
});

const certificate = new aws.acm.Certificate("cert", {
    domainName: "example.com",
}, { provider: usEast1 });
```

### Async/Await Patterns

```typescript
export = async () => {
    // Fetch data before creating resources
    const zones = await aws.getAvailabilityZones({
        state: "available",
    });

    // Use in resource creation
    const subnets = zones.names.map((az, i) =>
        new aws.ec2.Subnet(`subnet-${i}`, {
            vpcId: vpc.id,
            availabilityZone: az,
            cidrBlock: `10.0.${i}.0/24`,
        })
    );

    return {
        subnetIds: subnets.map(s => s.id),
    };
};
```

## Type Safety

### Strict Resource Typing

```typescript
import * as aws from "@pulumi/aws";

// Type-safe inputs
const bucketArgs: aws.s3.BucketArgs = {
    versioning: {
        enabled: true,
    },
    serverSideEncryptionConfiguration: {
        rule: {
            applyServerSideEncryptionByDefault: {
                sseAlgorithm: "AES256",
            },
        },
    },
};

const bucket = new aws.s3.Bucket("bucket", bucketArgs);
```

### Configuration Types

```typescript
interface AppConfig {
    port: number;
    replicas: number;
    image: string;
    environment: Record<string, string>;
}

const config = new pulumi.Config("app");
const appConfig = config.requireObject<AppConfig>("settings");
```

## ESM Support

For native ES modules:

```json
// package.json
{
    "type": "module"
}
```

```json
// tsconfig.json
{
    "compilerOptions": {
        "module": "nodenext",
        "moduleResolution": "nodenext"
    }
}
```

```yaml
# Pulumi.yaml (for older @pulumi/pulumi versions)
runtime:
  name: nodejs
  options:
    nodeargs: "--loader ts-node/esm --no-warnings"
```
