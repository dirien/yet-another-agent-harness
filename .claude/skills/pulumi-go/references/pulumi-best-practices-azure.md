# Pulumi Azure Best Practices (Go)

## Provider Selection

**IMPORTANT: Always use `azure-native` provider first.** The `azure-native` provider is auto-generated from Azure Resource Manager APIs and provides:
- 100% API coverage
- Same-day updates for new Azure features
- Full ARM template parity

Only use `pulumi-azure` (classic) for resources not yet in azure-native.

```go
// PREFERRED: azure-native provider
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure"

// FALLBACK: classic provider (only when needed)
import azureclassic "github.com/pulumi/pulumi-azure/sdk/v5/go/azure"
```

## Provider Configuration

```go
import (
    "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Use ESC for credentials via OIDC

        // Multi-subscription deployments
        prodProvider, _ := azure.NewProvider(ctx, "prod", &azure.ProviderArgs{
            SubscriptionId: pulumi.String("prod-subscription-id"),
        })

        return nil
    })
}
```

## Essential Resources

### Resource Group

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/resources"

rg, err := resources.NewResourceGroup(ctx, "rg", &resources.ResourceGroupArgs{
    ResourceGroupName: pulumi.Sprintf("rg-%s-%s", ctx.Project(), ctx.Stack()),
    Location:          pulumi.String("westeurope"),
    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
        "ManagedBy":   pulumi.String("Pulumi"),
    },
})
```

### Storage Account - Secure Pattern

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/storage"

storageAccount, err := storage.NewStorageAccount(ctx, "storage", &storage.StorageAccountArgs{
    ResourceGroupName: rg.Name,
    AccountName:       pulumi.String(strings.ReplaceAll(fmt.Sprintf("st%s%s", ctx.Project(), ctx.Stack()), "-", "")[:24]),
    Location:          rg.Location,

    Sku: &storage.SkuArgs{
        Name: pulumi.String(storage.SkuName_Standard_ZRS),
    },
    Kind: pulumi.String(storage.KindStorageV2),

    // Security settings
    AllowBlobPublicAccess:  pulumi.Bool(false),
    MinimumTlsVersion:      pulumi.String(storage.MinimumTlsVersion_TLS1_2),
    EnableHttpsTrafficOnly: pulumi.Bool(true),

    // Network security
    NetworkRuleSet: &storage.NetworkRuleSetArgs{
        DefaultAction: pulumi.String(storage.DefaultActionDeny),
        Bypass:        pulumi.String(storage.BypassAzureServices),
    },

    // Encryption
    Encryption: &storage.EncryptionArgs{
        Services: &storage.EncryptionServicesArgs{
            Blob: &storage.EncryptionServiceArgs{
                Enabled: pulumi.Bool(true),
                KeyType: pulumi.String(storage.KeyTypeAccount),
            },
        },
        KeySource: pulumi.String(storage.KeySource_Microsoft_Storage),
    },

    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})
```

### Virtual Network

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/network"

vnet, err := network.NewVirtualNetwork(ctx, "vnet", &network.VirtualNetworkArgs{
    ResourceGroupName:  rg.Name,
    VirtualNetworkName: pulumi.Sprintf("vnet-%s", ctx.Stack()),
    Location:           rg.Location,
    AddressSpace: &network.AddressSpaceArgs{
        AddressPrefixes: pulumi.StringArray{pulumi.String("10.0.0.0/16")},
    },
    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})

privateSubnet, _ := network.NewSubnet(ctx, "private", &network.SubnetArgs{
    ResourceGroupName:  rg.Name,
    VirtualNetworkName: vnet.Name,
    SubnetName:         pulumi.String("private"),
    AddressPrefix:      pulumi.String("10.0.2.0/24"),
    ServiceEndpoints: network.ServiceEndpointPropertiesFormatArray{
        &network.ServiceEndpointPropertiesFormatArgs{Service: pulumi.String("Microsoft.Storage")},
        &network.ServiceEndpointPropertiesFormatArgs{Service: pulumi.String("Microsoft.Sql")},
        &network.ServiceEndpointPropertiesFormatArgs{Service: pulumi.String("Microsoft.KeyVault")},
    },
    PrivateEndpointNetworkPolicies: pulumi.String(network.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled),
})
```

### Azure SQL Database

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/sql"

sqlServer, err := sql.NewServer(ctx, "sql", &sql.ServerArgs{
    ResourceGroupName: rg.Name,
    ServerName:        pulumi.Sprintf("sql-%s-%s", ctx.Project(), ctx.Stack()),
    Location:          rg.Location,

    // Use AAD authentication
    Administrators: &sql.ServerExternalAdministratorArgs{
        AdministratorType:         pulumi.String(sql.AdministratorTypeActiveDirectory),
        AzureADOnlyAuthentication: pulumi.Bool(true),
        Login:                     pulumi.String("sql-admins"),
        Sid:                       pulumi.String(aadGroupId),
        TenantId:                  pulumi.String(tenantId),
    },

    MinimalTlsVersion:   pulumi.String("1.2"),
    PublicNetworkAccess: pulumi.String(sql.ServerPublicNetworkAccessDisabled),

    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})

isProd := ctx.Stack() == "prod"

database, _ := sql.NewDatabase(ctx, "db", &sql.DatabaseArgs{
    ResourceGroupName: rg.Name,
    ServerName:        sqlServer.Name,
    DatabaseName:      pulumi.String("app"),

    Sku: func() *sql.SkuArgs {
        if isProd {
            return &sql.SkuArgs{
                Name:     pulumi.String("GP_Gen5"),
                Tier:     pulumi.String("GeneralPurpose"),
                Family:   pulumi.String("Gen5"),
                Capacity: pulumi.Int(2),
            }
        }
        return &sql.SkuArgs{
            Name:     pulumi.String("GP_S_Gen5"),
            Tier:     pulumi.String("GeneralPurpose"),
            Family:   pulumi.String("Gen5"),
            Capacity: pulumi.Int(1),
        }
    }(),

    ZoneRedundant: pulumi.Bool(isProd),

    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})
```

### Key Vault

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/keyvault"

keyVault, err := keyvault.NewVault(ctx, "kv", &keyvault.VaultArgs{
    ResourceGroupName: rg.Name,
    VaultName:         pulumi.String(fmt.Sprintf("kv-%s-%s", ctx.Project(), ctx.Stack())[:24]),
    Location:          rg.Location,

    Properties: &keyvault.VaultPropertiesArgs{
        TenantId: pulumi.String(tenantId),
        Sku: &keyvault.SkuArgs{
            Family: pulumi.String(keyvault.SkuFamilyA),
            Name:   keyvault.SkuNameStandard,
        },

        // Use RBAC instead of access policies
        EnableRbacAuthorization: pulumi.Bool(true),

        // Security
        EnableSoftDelete: pulumi.Bool(true),
        SoftDeleteRetentionInDays: pulumi.Int(func() int {
            if ctx.Stack() == "prod" {
                return 90
            }
            return 7
        }()),
        EnablePurgeProtection: pulumi.Bool(ctx.Stack() == "prod"),

        // Network
        NetworkAcls: &keyvault.NetworkRuleSetArgs{
            DefaultAction: pulumi.String(keyvault.NetworkRuleActionDeny),
            Bypass:        pulumi.String(keyvault.NetworkRuleBypassOptionsAzureServices),
        },
    },

    Tags: pulumi.StringMap{
        "Environment": pulumi.String(ctx.Stack()),
    },
})
```

## Security Best Practices

### Managed Identity

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/managedidentity"

identity, _ := managedidentity.NewUserAssignedIdentity(ctx, "app-identity", &managedidentity.UserAssignedIdentityArgs{
    ResourceGroupName: rg.Name,
    ResourceName:      pulumi.Sprintf("id-%s-%s", ctx.Project(), ctx.Stack()),
    Location:          rg.Location,
})

// Assign roles using RBAC
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/authorization"

_, _ = authorization.NewRoleAssignment(ctx, "blob-reader", &authorization.RoleAssignmentArgs{
    PrincipalId:      identity.PrincipalId,
    PrincipalType:    pulumi.String(authorization.PrincipalTypeServicePrincipal),
    RoleDefinitionId: pulumi.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/2a2b9908-6ea1-4ae2-8e65-a410df84e7d1", subscriptionId),
    Scope:            storageAccount.ID(),
})
```

## Monitoring

```go
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/operationalinsights"
import "github.com/pulumi/pulumi-azure-native-sdk/v2/go/azure/insights"

logAnalytics, _ := operationalinsights.NewWorkspace(ctx, "logs", &operationalinsights.WorkspaceArgs{
    ResourceGroupName: rg.Name,
    WorkspaceName:     pulumi.Sprintf("log-%s", ctx.Stack()),
    Location:          rg.Location,
    Sku: &operationalinsights.WorkspaceSkuArgs{
        Name: pulumi.String(operationalinsights.WorkspaceSkuNameEnumPerGB2018),
    },
    RetentionInDays: pulumi.Int(30),
})

appInsights, _ := insights.NewComponent(ctx, "insights", &insights.ComponentArgs{
    ResourceGroupName: rg.Name,
    ResourceName:      pulumi.Sprintf("appi-%s", ctx.Stack()),
    Location:          rg.Location,
    ApplicationType:   pulumi.String(insights.ApplicationTypeWeb),
    Kind:              pulumi.String("web"),
    WorkspaceResourceId: logAnalytics.ID(),
})
```

## Tagging Strategy

```go
defaultTags := pulumi.StringMap{
    "Environment": pulumi.String(ctx.Stack()),
    "Project":     pulumi.String(ctx.Project()),
    "ManagedBy":   pulumi.String("Pulumi"),
    "CostCenter":  pulumi.String("engineering"),
}
```
