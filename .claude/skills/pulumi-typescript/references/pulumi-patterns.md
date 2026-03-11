# Pulumi Infrastructure Patterns

## Component Resources

Encapsulate related resources into reusable components:

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

export interface VpcArgs {
    cidrBlock: string;
    azCount: number;
    enableNatGateway?: boolean;
}

export class Vpc extends pulumi.ComponentResource {
    public readonly vpcId: pulumi.Output<string>;
    public readonly publicSubnetIds: pulumi.Output<string>[];
    public readonly privateSubnetIds: pulumi.Output<string>[];

    constructor(name: string, args: VpcArgs, opts?: pulumi.ComponentResourceOptions) {
        super("custom:network:Vpc", name, {}, opts);

        const vpc = new aws.ec2.Vpc(`${name}-vpc`, {
            cidrBlock: args.cidrBlock,
            enableDnsHostnames: true,
            enableDnsSupport: true,
            tags: { Name: name },
        }, { parent: this });

        this.vpcId = vpc.id;

        // Create subnets...
        this.publicSubnetIds = [];
        this.privateSubnetIds = [];

        this.registerOutputs({
            vpcId: this.vpcId,
            publicSubnetIds: this.publicSubnetIds,
            privateSubnetIds: this.privateSubnetIds,
        });
    }
}
```

## Stack References

Share outputs between stacks:

```typescript
// In networking stack - export outputs
export const vpcId = vpc.id;
export const privateSubnetIds = privateSubnets.map(s => s.id);

// In application stack - import outputs
const networkStack = new pulumi.StackReference("myorg/networking/prod");
const vpcId = networkStack.getOutput("vpcId");
const subnetIds = networkStack.getOutput("privateSubnetIds") as pulumi.Output<string[]>;
```

## Transformations

Apply transformations to all resources:

```typescript
// Add tags to all taggable resources
pulumi.runtime.registerStackTransformation((args) => {
    if (args.props["tags"]) {
        args.props["tags"] = {
            ...args.props["tags"],
            Environment: pulumi.getStack(),
            ManagedBy: "Pulumi",
            Project: pulumi.getProject(),
        };
    }
    return { props: args.props, opts: args.opts };
});

// Enforce encryption on all S3 buckets
pulumi.runtime.registerStackTransformation((args) => {
    if (args.type === "aws:s3/bucket:Bucket") {
        args.props["serverSideEncryptionConfiguration"] = {
            rule: {
                applyServerSideEncryptionByDefault: {
                    sseAlgorithm: "aws:kms",
                },
            },
        };
    }
    return { props: args.props, opts: args.opts };
});
```

## Resource Protection

Prevent accidental deletion of critical resources:

```typescript
// Protect resource from deletion
const db = new aws.rds.Instance("prod-db", {
    // ... config
}, { protect: true });

// Retain resource on stack destroy (useful for data)
const bucket = new aws.s3.Bucket("data-bucket", {
    // ... config
}, { retainOnDelete: true });
```

## Component Resource Pitfalls

Common mistakes to avoid:

1. **Forgetting `parent: this`** - Child resources won't be properly associated; deletion may fail
2. **Skipping `registerOutputs()`** - Outputs won't be tracked; cross-stack references break
3. **Hardcoding values** - Reduces reusability; use args interface instead
4. **Overly complex components** - Keep components focused on single logical units
5. **Missing input validation** - Validate args before creating resources

## Stack Pitfalls

1. **Hardcoding environment values** - Use configuration or ESC instead
2. **Sharing state between environments** - Each stack should have isolated state
3. **Circular stack references** - Design dependency graph carefully
4. **Unprotected production resources** - Use `protect: true` for critical infra
5. **Unencrypted secrets** - Always use `--secret` flag or ESC secrets

## Dynamic Providers

Create custom resources with dynamic providers:

```typescript
import * as pulumi from "@pulumi/pulumi";

const myProvider: pulumi.dynamic.ResourceProvider = {
    async create(inputs) {
        // Create logic
        return { id: "my-id", outs: { result: "created" } };
    },
    async update(id, olds, news) {
        // Update logic
        return { outs: { result: "updated" } };
    },
    async delete(id, props) {
        // Delete logic
    },
};

class MyResource extends pulumi.dynamic.Resource {
    public readonly result!: pulumi.Output<string>;

    constructor(name: string, props: {}, opts?: pulumi.CustomResourceOptions) {
        super(myProvider, name, { result: undefined, ...props }, opts);
    }
}
```

## Configuration Patterns

### Environment-Aware Configuration

```typescript
const stack = pulumi.getStack();
const isProd = stack === "prod";

const config = new pulumi.Config();
const instanceType = config.get("instanceType") || (isProd ? "t3.large" : "t3.small");
const minSize = isProd ? 3 : 1;
const maxSize = isProd ? 10 : 2;
```

### Nested Configuration

```yaml
# Pulumi.dev.yaml via ESC
values:
  pulumiConfig:
    myapp:database:
      host: localhost
      port: 5432
      name: myapp
```

```typescript
interface DbConfig {
    host: string;
    port: number;
    name: string;
}

const config = new pulumi.Config("myapp");
const dbConfig = config.requireObject<DbConfig>("database");
```

## Testing Patterns

### Unit Testing with Mocks

```typescript
import * as pulumi from "@pulumi/pulumi";

pulumi.runtime.setMocks({
    newResource: (args) => {
        return { id: `${args.name}-id`, state: args.inputs };
    },
    call: (args) => {
        return args.inputs;
    },
});

describe("Infrastructure", () => {
    let infra: typeof import("./index");

    beforeAll(async () => {
        infra = await import("./index");
    });

    it("should create a bucket with versioning", async () => {
        const versioning = await new Promise((resolve) =>
            infra.bucket.versioning.apply((v) => resolve(v))
        );
        expect(versioning).toEqual({ enabled: true });
    });
});
```

### Integration Testing with Automation API

```typescript
import { LocalWorkspace } from "@pulumi/pulumi/automation";

describe("Integration", () => {
    it("should deploy successfully", async () => {
        const stack = await LocalWorkspace.createOrSelectStack({
            stackName: "test",
            projectName: "test-project",
            program: async () => {
                // Inline program
            },
        });

        const result = await stack.up();
        expect(result.summary.result).toBe("succeeded");

        await stack.destroy();
    });
});
```

## CI/CD Patterns

### GitHub Actions with ESC

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: pulumi/auth-actions@v1
        with:
          organization: myorg
          requested-token-type: urn:pulumi:token-type:access_token:organization

      - name: Deploy with ESC
        run: |
          pulumi env run myorg/aws-prod -- pulumi up --yes
```
