# Pulumi GCP Best Practices (Go)

## Provider Configuration

```go
package main

import (
    "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Use ESC for credentials via OIDC (Workload Identity Federation)

        // Multi-project deployments
        prodProvider, _ := gcp.NewProvider(ctx, "prod", &gcp.ProviderArgs{
            Project: pulumi.String("my-prod-project"),
            Region:  pulumi.String("europe-west1"),
        })

        return nil
    })
}
```

## Essential Resources

### Cloud Storage - Secure Bucket

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"

bucket, err := storage.NewBucket(ctx, "data-bucket", &storage.BucketArgs{
    Name:         pulumi.Sprintf("%s-data-%s", gcp.GetProject(), ctx.Stack()),
    Location:     pulumi.String("EU"),
    StorageClass: pulumi.String("STANDARD"),

    // Uniform bucket-level access (recommended)
    UniformBucketLevelAccess: pulumi.Bool(true),

    // Versioning for data protection
    Versioning: &storage.BucketVersioningArgs{
        Enabled: pulumi.Bool(true),
    },

    // Prevent public access
    PublicAccessPrevention: pulumi.String("enforced"),

    // Lifecycle rules
    LifecycleRules: storage.BucketLifecycleRuleArray{
        &storage.BucketLifecycleRuleArgs{
            Action:    &storage.BucketLifecycleRuleActionArgs{Type: pulumi.String("SetStorageClass"), StorageClass: pulumi.String("NEARLINE")},
            Condition: &storage.BucketLifecycleRuleConditionArgs{Age: pulumi.Int(30)},
        },
        &storage.BucketLifecycleRuleArgs{
            Action:    &storage.BucketLifecycleRuleActionArgs{Type: pulumi.String("SetStorageClass"), StorageClass: pulumi.String("COLDLINE")},
            Condition: &storage.BucketLifecycleRuleConditionArgs{Age: pulumi.Int(90)},
        },
        &storage.BucketLifecycleRuleArgs{
            Action:    &storage.BucketLifecycleRuleActionArgs{Type: pulumi.String("Delete")},
            Condition: &storage.BucketLifecycleRuleConditionArgs{Age: pulumi.Int(730)},
        },
    },

    Labels: pulumi.StringMap{
        "environment": pulumi.String(ctx.Stack()),
        "managed_by":  pulumi.String("pulumi"),
    },
})

// IAM binding instead of ACLs
_, _ = storage.NewBucketIAMMember(ctx, "app-reader", &storage.BucketIAMMemberArgs{
    Bucket: bucket.Name,
    Role:   pulumi.String("roles/storage.objectViewer"),
    Member: serviceAccount.Email.ApplyT(func(email string) string {
        return fmt.Sprintf("serviceAccount:%s", email)
    }).(pulumi.StringOutput),
})
```

### VPC Network

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"

// Custom VPC (don't use default)
vpc, _ := compute.NewNetwork(ctx, "vpc", &compute.NetworkArgs{
    Name:                  pulumi.Sprintf("vpc-%s", ctx.Stack()),
    AutoCreateSubnetworks: pulumi.Bool(false), // Always use custom subnets
    RoutingMode:           pulumi.String("REGIONAL"),
})

// Private subnet
privateSubnet, _ := compute.NewSubnetwork(ctx, "private", &compute.SubnetworkArgs{
    Name:        pulumi.Sprintf("subnet-private-%s", ctx.Stack()),
    IpCidrRange: pulumi.String("10.0.0.0/20"),
    Region:      pulumi.String("europe-west1"),
    Network:     vpc.ID(),

    // Enable Private Google Access
    PrivateIpGoogleAccess: pulumi.Bool(true),

    // Enable VPC Flow Logs
    LogConfig: &compute.SubnetworkLogConfigArgs{
        AggregationInterval: pulumi.String("INTERVAL_5_SEC"),
        FlowSampling:        pulumi.Float64(0.5),
        Metadata:            pulumi.String("INCLUDE_ALL_METADATA"),
    },

    // Secondary ranges for GKE
    SecondaryIpRanges: compute.SubnetworkSecondaryIpRangeArray{
        &compute.SubnetworkSecondaryIpRangeArgs{RangeName: pulumi.String("pods"), IpCidrRange: pulumi.String("10.4.0.0/14")},
        &compute.SubnetworkSecondaryIpRangeArgs{RangeName: pulumi.String("services"), IpCidrRange: pulumi.String("10.8.0.0/20")},
    },
})

// Cloud NAT for private instances
router, _ := compute.NewRouter(ctx, "router", &compute.RouterArgs{
    Name:    pulumi.Sprintf("router-%s", ctx.Stack()),
    Region:  pulumi.String("europe-west1"),
    Network: vpc.ID(),
})

_, _ = compute.NewRouterNat(ctx, "nat", &compute.RouterNatArgs{
    Name:                          pulumi.Sprintf("nat-%s", ctx.Stack()),
    Router:                        router.Name,
    Region:                        pulumi.String("europe-west1"),
    NatIpAllocateOption:           pulumi.String("AUTO_ONLY"),
    SourceSubnetworkIpRangesToNat: pulumi.String("ALL_SUBNETWORKS_ALL_IP_RANGES"),
    LogConfig: &compute.RouterNatLogConfigArgs{
        Enable: pulumi.Bool(true),
        Filter: pulumi.String("ERRORS_ONLY"),
    },
})
```

### Cloud SQL - PostgreSQL

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/sql"
import "github.com/pulumi/pulumi-random/sdk/v4/go/random"

dbPassword, _ := random.NewRandomPassword(ctx, "db-password", &random.RandomPasswordArgs{
    Length:  pulumi.Int(24),
    Special: pulumi.Bool(true),
})

isProd := ctx.Stack() == "prod"

sqlInstance, _ := sql.NewDatabaseInstance(ctx, "postgres", &sql.DatabaseInstanceArgs{
    Name:            pulumi.Sprintf("sql-%s-%s", ctx.Project(), ctx.Stack()),
    Region:          pulumi.String("europe-west1"),
    DatabaseVersion: pulumi.String("POSTGRES_15"),

    Settings: &sql.DatabaseInstanceSettingsArgs{
        Tier: pulumi.String(func() string {
            if isProd {
                return "db-custom-4-16384"
            }
            return "db-f1-micro"
        }()),

        AvailabilityType: pulumi.String(func() string {
            if isProd {
                return "REGIONAL"
            }
            return "ZONAL"
        }()),

        BackupConfiguration: &sql.DatabaseInstanceSettingsBackupConfigurationArgs{
            Enabled:                     pulumi.Bool(true),
            StartTime:                   pulumi.String("03:00"),
            PointInTimeRecoveryEnabled:  pulumi.Bool(isProd),
            BackupRetentionSettings: &sql.DatabaseInstanceSettingsBackupConfigurationBackupRetentionSettingsArgs{
                RetainedBackups: pulumi.Int(7),
            },
        },

        IpConfiguration: &sql.DatabaseInstanceSettingsIpConfigurationArgs{
            Ipv4Enabled:    pulumi.Bool(false),
            PrivateNetwork: vpc.ID(),
            RequireSsl:     pulumi.Bool(true),
        },

        MaintenanceWindow: &sql.DatabaseInstanceSettingsMaintenanceWindowArgs{
            Day:  pulumi.Int(7),
            Hour: pulumi.Int(3),
        },

        DatabaseFlags: sql.DatabaseInstanceSettingsDatabaseFlagArray{
            &sql.DatabaseInstanceSettingsDatabaseFlagArgs{Name: pulumi.String("log_checkpoints"), Value: pulumi.String("on")},
            &sql.DatabaseInstanceSettingsDatabaseFlagArgs{Name: pulumi.String("log_connections"), Value: pulumi.String("on")},
        },

        UserLabels: pulumi.StringMap{
            "environment": pulumi.String(ctx.Stack()),
        },
    },

    DeletionProtection: pulumi.Bool(isProd),
})
```

### Cloud Run

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrun"

service, _ := cloudrun.NewService(ctx, "api", &cloudrun.ServiceArgs{
    Name:     pulumi.Sprintf("api-%s", ctx.Stack()),
    Location: pulumi.String("europe-west1"),

    Template: &cloudrun.ServiceTemplateArgs{
        Spec: &cloudrun.ServiceTemplateSpecArgs{
            ServiceAccountName: serviceAccount.Email,
            Containers: cloudrun.ServiceTemplateSpecContainerArray{
                &cloudrun.ServiceTemplateSpecContainerArgs{
                    Image: pulumi.Sprintf("gcr.io/%s/api:%s", gcp.GetProject(), imageTag),
                    Resources: &cloudrun.ServiceTemplateSpecContainerResourcesArgs{
                        Limits: pulumi.StringMap{
                            "cpu":    pulumi.String("1000m"),
                            "memory": pulumi.String("512Mi"),
                        },
                    },
                    Envs: cloudrun.ServiceTemplateSpecContainerEnvArray{
                        &cloudrun.ServiceTemplateSpecContainerEnvArgs{
                            Name:  pulumi.String("ENVIRONMENT"),
                            Value: pulumi.String(ctx.Stack()),
                        },
                    },
                    Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
                        &cloudrun.ServiceTemplateSpecContainerPortArgs{ContainerPort: pulumi.Int(8080)},
                    },
                },
            },
            ContainerConcurrency: pulumi.Int(80),
            TimeoutSeconds:       pulumi.Int(300),
        },

        Metadata: &cloudrun.ServiceTemplateMetadataArgs{
            Annotations: pulumi.StringMap{
                "autoscaling.knative.dev/minScale": pulumi.String(func() string {
                    if isProd {
                        return "2"
                    }
                    return "0"
                }()),
                "autoscaling.knative.dev/maxScale": pulumi.String("100"),
            },
        },
    },

    Traffics: cloudrun.ServiceTrafficArray{
        &cloudrun.ServiceTrafficArgs{
            Percent:        pulumi.Int(100),
            LatestRevision: pulumi.Bool(true),
        },
    },

    AutogenerateRevisionName: pulumi.Bool(true),
})

// Public access (for APIs)
_, _ = cloudrun.NewIamMember(ctx, "public", &cloudrun.IamMemberArgs{
    Service:  service.Name,
    Location: service.Location,
    Role:     pulumi.String("roles/run.invoker"),
    Member:   pulumi.String("allUsers"),
})
```

## Security Best Practices

### Service Accounts & Workload Identity

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"

// Create dedicated service account for each workload
sa, _ := serviceaccount.NewAccount(ctx, "app-sa", &serviceaccount.AccountArgs{
    AccountId:   pulumi.Sprintf("sa-app-%s", ctx.Stack()),
    DisplayName: pulumi.String("Application Service Account"),
})

// Grant minimal permissions
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"

_, _ = projects.NewIAMMember(ctx, "storage-access", &projects.IAMMemberArgs{
    Project: pulumi.String(gcp.GetProject()),
    Role:    pulumi.String("roles/storage.objectViewer"),
    Member:  sa.Email.ApplyT(func(email string) string {
        return fmt.Sprintf("serviceAccount:%s", email)
    }).(pulumi.StringOutput),
})

// Workload Identity binding for GKE
_, _ = serviceaccount.NewIAMBinding(ctx, "workload-identity", &serviceaccount.IAMBindingArgs{
    ServiceAccountId: sa.Name,
    Role:             pulumi.String("roles/iam.workloadIdentityUser"),
    Members: pulumi.StringArray{
        pulumi.Sprintf("serviceAccount:%s.svc.id.goog[default/app]", gcp.GetProject()),
    },
})
```

### Secret Manager

```go
import "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/secretmanager"

secret, _ := secretmanager.NewSecret(ctx, "api-key", &secretmanager.SecretArgs{
    SecretId: pulumi.Sprintf("api-key-%s", ctx.Stack()),
    Replication: &secretmanager.SecretReplicationArgs{
        Automatic: pulumi.Bool(true),
    },
    Labels: pulumi.StringMap{
        "environment": pulumi.String(ctx.Stack()),
    },
})

_, _ = secretmanager.NewSecretVersion(ctx, "api-key-v1", &secretmanager.SecretVersionArgs{
    Secret:     secret.ID(),
    SecretData: cfg.RequireSecret("apiKey"),
})

// Grant access to service account
_, _ = secretmanager.NewSecretIamMember(ctx, "app-access", &secretmanager.SecretIamMemberArgs{
    SecretId: secret.SecretId,
    Role:     pulumi.String("roles/secretmanager.secretAccessor"),
    Member:   sa.Email.ApplyT(func(email string) string {
        return fmt.Sprintf("serviceAccount:%s", email)
    }).(pulumi.StringOutput),
})
```

## Labels Strategy

```go
// Standard labels for all GCP resources
// Note: GCP labels must be lowercase with letters, numbers, underscores, dashes
defaultLabels := pulumi.StringMap{
    "environment": pulumi.String(ctx.Stack()),
    "project":     pulumi.String(strings.ToLower(strings.ReplaceAll(ctx.Project(), "[^a-z0-9-]", "-"))),
    "managed_by":  pulumi.String("pulumi"),
    "cost_center": pulumi.String("engineering"),
}
```
