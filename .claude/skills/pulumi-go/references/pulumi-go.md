# Pulumi Go Reference

## Project Setup

### go.mod

```go
module myproject

go 1.21

require (
    github.com/pulumi/pulumi-aws/sdk/v6 v6.0.0
    github.com/pulumi/pulumi/sdk/v3 v3.0.0
)
```

### Pulumi.yaml

```yaml
name: my-project
description: My Pulumi project
runtime: go
```

### main.go Template

```go
package main

import (
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Your infrastructure code here
        return nil
    })
}
```

## Working with Inputs and Outputs

### Output Types

```go
// String output
var bucketName pulumi.StringOutput = bucket.ID()

// Int output
var port pulumi.IntOutput = pulumi.Int(8080).ToIntOutput()

// Array output
var subnetIds pulumi.StringArrayOutput

// Map output
var tags pulumi.StringMapOutput
```

### Apply Transformations

```go
// Single output transformation
url := bucket.BucketDomainName.ApplyT(func(domain string) string {
    return "https://" + domain
}).(pulumi.StringOutput)

// Multiple outputs
combined := pulumi.All(bucket.ID(), bucket.Arn).ApplyT(
    func(args []interface{}) map[string]string {
        return map[string]string{
            "id":  args[0].(string),
            "arn": args[1].(string),
        }
    },
)

// With error handling
result := bucket.ID().ApplyT(func(id string) (string, error) {
    if id == "" {
        return "", fmt.Errorf("empty bucket ID")
    }
    return strings.ToUpper(id), nil
}).(pulumi.StringOutput)
```

### Input Types

```go
// Creating inputs
stringInput := pulumi.String("value")
intInput := pulumi.Int(42)
boolInput := pulumi.Bool(true)

// Array inputs
stringArrayInput := pulumi.StringArray{
    pulumi.String("a"),
    pulumi.String("b"),
}

// Map inputs
stringMapInput := pulumi.StringMap{
    "key1": pulumi.String("value1"),
    "key2": pulumi.String("value2"),
}
```

## Configuration

```go
import "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Default namespace (project name)
        cfg := config.New(ctx, "")

        // Named namespace
        awsCfg := config.New(ctx, "aws")

        // Required values
        region := awsCfg.Require("region")

        // Optional with default
        instanceType := cfg.Get("instanceType")
        if instanceType == "" {
            instanceType = "t3.small"
        }

        // Secrets (automatically decrypted)
        apiKey := cfg.RequireSecret("apiKey")

        // Structured configuration
        type DatabaseConfig struct {
            Host string `json:"host"`
            Port int    `json:"port"`
        }
        var dbConfig DatabaseConfig
        cfg.RequireObject("database", &dbConfig)

        return nil
    })
}
```

## Provider Configuration

```go
import "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Explicit provider for different region
        usEast1, err := aws.NewProvider(ctx, "us-east-1", &aws.ProviderArgs{
            Region: pulumi.String("us-east-1"),
        })
        if err != nil {
            return err
        }

        // Use provider with resource
        cert, err := acm.NewCertificate(ctx, "cert", &acm.CertificateArgs{
            DomainName: pulumi.String("example.com"),
        }, pulumi.Provider(usEast1))
        if err != nil {
            return err
        }

        return nil
    })
}
```

## Async Data Sources

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Get availability zones
        zones, err := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{
            State: pulumi.StringRef("available"),
        })
        if err != nil {
            return err
        }

        // Use zones
        for i, az := range zones.Names {
            _, err := ec2.NewSubnet(ctx, fmt.Sprintf("subnet-%d", i), &ec2.SubnetArgs{
                VpcId:            vpc.ID(),
                AvailabilityZone: pulumi.String(az),
                CidrBlock:        pulumi.Sprintf("10.0.%d.0/24", i),
            })
            if err != nil {
                return err
            }
        }

        return nil
    })
}
```

## Error Handling Best Practices

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Always wrap errors with context
        bucket, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{})
        if err != nil {
            return fmt.Errorf("failed to create S3 bucket: %w", err)
        }

        // Check multiple resources
        resources := []struct {
            name string
            fn   func() error
        }{
            {"vpc", func() error { _, err := createVpc(ctx); return err }},
            {"subnet", func() error { _, err := createSubnet(ctx); return err }},
        }

        for _, r := range resources {
            if err := r.fn(); err != nil {
                return fmt.Errorf("failed to create %s: %w", r.name, err)
            }
        }

        return nil
    })
}
```

## Building and Running

```bash
# Build before running (recommended for faster execution)
go build -o $(basename $(pwd))
pulumi up

# Or let Pulumi build automatically
pulumi up  # Runs: go build -o pulumi-bin
```
