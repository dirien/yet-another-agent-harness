# Pulumi AWS Best Practices (Go)

## Provider Configuration

```go
package main

import (
    "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Use ESC for credentials via OIDC - avoid static credentials

        // Multi-region deployments with explicit providers
        usEast1, err := aws.NewProvider(ctx, "us-east-1", &aws.ProviderArgs{
            Region: pulumi.String("us-east-1"),
        })
        if err != nil {
            return err
        }

        // Use provider with resources (ACM certs for CloudFront must be in us-east-1)
        _, err = acm.NewCertificate(ctx, "cert", &acm.CertificateArgs{
            DomainName:       pulumi.String("example.com"),
            ValidationMethod: pulumi.String("DNS"),
        }, pulumi.Provider(usEast1))

        return err
    })
}
```

## Essential Resources

### S3 - Secure Bucket Pattern

```go
bucket, err := s3.NewBucket(ctx, "data-bucket", &s3.BucketArgs{
    // Enable versioning for data protection
    Versioning: &s3.BucketVersioningArgs{
        Enabled: pulumi.Bool(true),
    },
    // Server-side encryption
    ServerSideEncryptionConfiguration: &s3.BucketServerSideEncryptionConfigurationArgs{
        Rule: &s3.BucketServerSideEncryptionConfigurationRuleArgs{
            ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs{
                SseAlgorithm:   pulumi.String("aws:kms"),
                KmsMasterKeyId: kmsKey.Arn,
            },
            BucketKeyEnabled: pulumi.Bool(true),
        },
    },
    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
        "ManagedBy":   pulumi.String("Pulumi"),
    },
})

// Block public access
_, err = s3.NewBucketPublicAccessBlock(ctx, "block-public", &s3.BucketPublicAccessBlockArgs{
    Bucket:                bucket.ID(),
    BlockPublicAcls:       pulumi.Bool(true),
    BlockPublicPolicy:     pulumi.Bool(true),
    IgnorePublicAcls:      pulumi.Bool(true),
    RestrictPublicBuckets: pulumi.Bool(true),
})
```

### VPC - Network Pattern

```go
import "github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"

vpc, err := ec2.NewVpc(ctx, "main-vpc", &ec2.VpcArgs{
    CidrBlock:                 pulumi.StringRef("10.0.0.0/16"),
    NumberOfAvailabilityZones: pulumi.IntRef(3),
    NatGateways: &ec2.NatGatewayConfigurationArgs{
        Strategy: ec2.NatGatewayStrategyOnePerAz,
    },
    SubnetSpecs: []ec2.SubnetSpecArgs{
        {Type: ec2.SubnetTypePublic, CidrMask: pulumi.IntRef(24)},
        {Type: ec2.SubnetTypePrivate, CidrMask: pulumi.IntRef(24)},
        {Type: ec2.SubnetTypeIsolated, CidrMask: pulumi.IntRef(24)},
    },
    Tags: pulumi.StringMap{"Name": pulumi.String("main-vpc")},
})

ctx.Export("vpcId", vpc.VpcId)
ctx.Export("publicSubnetIds", vpc.PublicSubnetIds)
ctx.Export("privateSubnetIds", vpc.PrivateSubnetIds)
```

### RDS - Database Pattern

```go
dbSubnetGroup, _ := rds.NewSubnetGroup(ctx, "db-subnets", &rds.SubnetGroupArgs{
    SubnetIds: vpc.IsolatedSubnetIds,
})

isProd := ctx.Stack() == "prod"

database, err := rds.NewInstance(ctx, "postgres", &rds.InstanceArgs{
    Engine:            pulumi.String("postgres"),
    EngineVersion:     pulumi.String("15.4"),
    InstanceClass:     pulumi.String("db.t3.medium"),
    AllocatedStorage:  pulumi.Int(20),
    MaxAllocatedStorage: pulumi.Int(100),

    DbName:   pulumi.String("myapp"),
    Username: pulumi.String("admin"),
    Password: cfg.RequireSecret("dbPassword"),

    DbSubnetGroupName:   dbSubnetGroup.Name,
    VpcSecurityGroupIds: pulumi.StringArray{dbSecurityGroup.ID()},

    MultiAz: pulumi.Bool(isProd),

    BackupRetentionPeriod: pulumi.Int(7),
    BackupWindow:         pulumi.String("03:00-04:00"),
    MaintenanceWindow:    pulumi.String("Mon:04:00-Mon:05:00"),

    StorageEncrypted:   pulumi.Bool(true),
    DeletionProtection: pulumi.Bool(isProd),

    PerformanceInsightsEnabled:         pulumi.Bool(true),
    PerformanceInsightsRetentionPeriod: pulumi.Int(7),

    SkipFinalSnapshot:       pulumi.Bool(!isProd),
    FinalSnapshotIdentifier: pulumi.Sprintf("%s-final-snapshot", ctx.Project()),

    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})
```

### Lambda Function

```go
role, _ := iam.NewRole(ctx, "lambda-role", &iam.RoleArgs{
    AssumeRolePolicy: pulumi.String(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Action": "sts:AssumeRole",
            "Principal": {"Service": "lambda.amazonaws.com"},
            "Effect": "Allow"
        }]
    }`),
})

_, _ = iam.NewRolePolicyAttachment(ctx, "lambda-basic", &iam.RolePolicyAttachmentArgs{
    Role:      role.Name,
    PolicyArn: iam.ManagedPolicyAWSLambdaBasicExecutionRole,
})

function, err := lambda.NewFunction(ctx, "api-handler", &lambda.FunctionArgs{
    Runtime: pulumi.String("nodejs18.x"),
    Handler: pulumi.String("index.handler"),
    Code:    pulumi.NewFileArchive("./lambda"),
    Role:    role.Arn,
    MemorySize: pulumi.Int(512),
    Timeout:    pulumi.Int(30),
    Environment: &lambda.FunctionEnvironmentArgs{
        Variables: pulumi.StringMap{
            "TABLE_NAME": dynamoTable.Name,
            "STAGE":      pulumi.String(ctx.Stack()),
        },
    },
    TracingConfig: &lambda.FunctionTracingConfigArgs{
        Mode: pulumi.String("Active"),
    },
})
```

## Security Best Practices

### IAM - Least Privilege

```go
customPolicy, err := iam.NewPolicy(ctx, "app-policy", &iam.PolicyArgs{
    Policy: pulumi.All(bucket.Arn).ApplyT(func(args []interface{}) string {
        bucketArn := args[0].(string)
        return fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [{
                "Effect": "Allow",
                "Action": ["s3:GetObject", "s3:PutObject"],
                "Resource": ["%s/*"]
            }]
        }`, bucketArn)
    }).(pulumi.StringOutput),
})
```

### Secrets Manager

```go
secret, _ := secretsmanager.NewSecret(ctx, "api-key", &secretsmanager.SecretArgs{
    Name: pulumi.Sprintf("%s/api-key", ctx.Stack()),
    RecoveryWindowInDays: pulumi.Int(func() int {
        if ctx.Stack() == "prod" {
            return 30
        }
        return 0
    }()),
})
```

## Tagging Strategy

```go
// Register transformation to add tags to all resources
ctx.RegisterStackTransformation(func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
    if args.Props["tags"] != nil {
        tags := args.Props["tags"].(pulumi.StringMap)
        tags["Environment"] = pulumi.String(ctx.Stack())
        tags["Project"] = pulumi.String(ctx.Project())
        tags["ManagedBy"] = pulumi.String("Pulumi")
        args.Props["tags"] = tags
    }
    return &pulumi.ResourceTransformationResult{
        Props: args.Props,
        Opts:  args.Opts,
    }
})
```
