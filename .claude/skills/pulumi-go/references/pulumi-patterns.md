# Pulumi Infrastructure Patterns (Go)

## Component Resources

```go
package main

import (
    "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpcArgs struct {
    CidrBlock        pulumi.StringInput
    AzCount          int
    EnableNatGateway bool
}

type Vpc struct {
    pulumi.ResourceState

    VpcId            pulumi.StringOutput   `pulumi:"vpcId"`
    PublicSubnetIds  pulumi.StringArrayOutput `pulumi:"publicSubnetIds"`
    PrivateSubnetIds pulumi.StringArrayOutput `pulumi:"privateSubnetIds"`
}

func NewVpc(ctx *pulumi.Context, name string, args *VpcArgs, opts ...pulumi.ResourceOption) (*Vpc, error) {
    component := &Vpc{}
    err := ctx.RegisterComponentResource("custom:network:Vpc", name, component, opts...)
    if err != nil {
        return nil, err
    }

    vpc, err := ec2.NewVpc(ctx, name+"-vpc", &ec2.VpcArgs{
        CidrBlock:          args.CidrBlock,
        EnableDnsHostnames: pulumi.Bool(true),
        EnableDnsSupport:   pulumi.Bool(true),
        Tags: pulumi.StringMap{
            "Name": pulumi.String(name),
        },
    }, pulumi.Parent(component))
    if err != nil {
        return nil, err
    }

    component.VpcId = vpc.ID().ToStringOutput()

    ctx.RegisterResourceOutputs(component, pulumi.Map{
        "vpcId":            component.VpcId,
        "publicSubnetIds":  component.PublicSubnetIds,
        "privateSubnetIds": component.PrivateSubnetIds,
    })

    return component, nil
}
```

## Stack References

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Reference networking stack
        networkStack, err := pulumi.NewStackReference(ctx, "myorg/networking/prod", nil)
        if err != nil {
            return err
        }

        vpcId := networkStack.GetStringOutput(pulumi.String("vpcId"))

        // Type assertion for arrays
        subnetIds := networkStack.GetOutput(pulumi.String("privateSubnetIds")).ApplyT(
            func(v interface{}) []string {
                ids := v.([]interface{})
                result := make([]string, len(ids))
                for i, id := range ids {
                    result[i] = id.(string)
                }
                return result
            },
        ).(pulumi.StringArrayOutput)

        return nil
    })
}
```

## Transformations

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Register transformation for all resources
        ctx.RegisterStackTransformation(func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
            // Add tags to all taggable resources
            if args.Props != nil {
                if tags, ok := args.Props["tags"]; ok {
                    tagMap := tags.(pulumi.Map)
                    tagMap["Environment"] = pulumi.String(ctx.Stack())
                    tagMap["ManagedBy"] = pulumi.String("Pulumi")
                    args.Props["tags"] = tagMap
                }
            }
            return &pulumi.ResourceTransformationResult{
                Props: args.Props,
                Opts:  args.Opts,
            }
        })

        return nil
    })
}
```

## Resource Protection

Prevent accidental deletion of critical resources:

```go
// Protect resource from deletion
db, err := rds.NewInstance(ctx, "prod-db", &rds.InstanceArgs{
    // ... config
}, pulumi.Protect(true))

// Retain resource on stack destroy (useful for data)
bucket, err := s3.NewBucket(ctx, "data-bucket", &s3.BucketArgs{
    // ... config
}, pulumi.RetainOnDelete(true))
```

## Component Resource Pitfalls

Common mistakes to avoid:

1. **Forgetting `pulumi.Parent(component)`** - Child resources won't be properly associated; deletion may fail
2. **Skipping `RegisterResourceOutputs()`** - Outputs won't be tracked; cross-stack references break
3. **Hardcoding values** - Reduces reusability; use args struct instead
4. **Overly complex components** - Keep components focused on single logical units
5. **Missing input validation** - Validate args before creating resources

## Stack Pitfalls

1. **Hardcoding environment values** - Use configuration or ESC instead
2. **Sharing state between environments** - Each stack should have isolated state
3. **Circular stack references** - Design dependency graph carefully
4. **Unprotected production resources** - Use `pulumi.Protect(true)` for critical infra
5. **Unencrypted secrets** - Always use `--secret` flag or ESC secrets

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

## Configuration Patterns

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "myapp")

        // Environment-aware defaults
        isProd := ctx.Stack() == "prod"

        instanceType := cfg.Get("instanceType")
        if instanceType == "" {
            if isProd {
                instanceType = "t3.large"
            } else {
                instanceType = "t3.small"
            }
        }

        // Structured configuration
        type DbConfig struct {
            Host string `json:"host"`
            Port int    `json:"port"`
            Name string `json:"name"`
        }

        var dbConfig DbConfig
        cfg.RequireObject("database", &dbConfig)

        return nil
    })
}
```

## Testing Patterns

### Unit Testing with Mocks

```go
package main

import (
    "testing"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/stretchr/testify/assert"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
    return args.Name + "_id", args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
    return args.Args, nil
}

func TestInfrastructure(t *testing.T) {
    err := pulumi.RunErr(func(ctx *pulumi.Context) error {
        // Test your infrastructure
        return nil
    }, pulumi.WithMocks("project", "stack", mocks(0)))

    assert.NoError(t, err)
}
```

### Integration Testing with Automation API

```go
package main

import (
    "context"
    "testing"

    "github.com/pulumi/pulumi/sdk/v3/go/auto"
    "github.com/stretchr/testify/assert"
)

func TestDeploy(t *testing.T) {
    ctx := context.Background()

    stack, err := auto.UpsertStackLocalSource(ctx, "test", ".")
    assert.NoError(t, err)

    // Deploy
    result, err := stack.Up(ctx)
    assert.NoError(t, err)
    assert.Equal(t, "succeeded", result.Summary.Result)

    // Cleanup
    _, err = stack.Destroy(ctx)
    assert.NoError(t, err)
}
```

## Resource Options

```go
// Common resource options
resource, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{},
    // Explicit dependencies
    pulumi.DependsOn([]pulumi.Resource{otherResource}),

    // Parent for components
    pulumi.Parent(component),

    // Protect from deletion
    pulumi.Protect(true),

    // Ignore changes
    pulumi.IgnoreChanges([]string{"tags"}),

    // Custom provider
    pulumi.Provider(customProvider),

    // Aliases for refactoring
    pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("old-name")}}),

    // Delete before replace
    pulumi.DeleteBeforeReplace(true),

    // Force replacement when specific properties change
    pulumi.ReplaceOnChanges([]string{"tags.Environment"}),

    // Skip delete when parent resource is deleted (useful for Kubernetes)
    pulumi.DeletedWith(parentResource),
)
```

## Resource Hooks

```go
// Define lifecycle hooks for custom logic
resource, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{},
    pulumi.Hooks(&pulumi.LifecycleHooks{
        BeforeCreate: func(ctx context.Context, args pulumi.HookArgs) error {
            fmt.Printf("Creating resource: %s\n", args.URN)
            return nil
        },
        AfterCreate: func(ctx context.Context, args pulumi.HookArgs) error {
            fmt.Printf("Created with outputs: %v\n", args.Outputs)
            return nil
        },
    }),
)
```

> **Note:** Hooks require `--run-program` flag when running `pulumi destroy`.
