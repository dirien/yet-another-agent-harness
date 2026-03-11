# Pulumi Azure Best Practices (TypeScript)

## Provider Selection

**IMPORTANT: Always use `azure-native` provider first.** The `azure-native` provider is auto-generated from Azure Resource Manager APIs and provides:
- 100% API coverage
- Same-day updates for new Azure features
- Full ARM template parity

Only use `@pulumi/azure` (classic) for resources not yet in azure-native or for specific legacy scenarios.

```typescript
// PREFERRED: azure-native provider
import * as azure from "@pulumi/azure-native";

// FALLBACK: classic provider (only when needed)
import * as azureClassic from "@pulumi/azure";
```

## Provider Configuration

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as azure from "@pulumi/azure-native";

// Use ESC for credentials via OIDC
// ESC handles AZURE_CLIENT_ID, AZURE_TENANT_ID, AZURE_SUBSCRIPTION_ID

// Multi-subscription deployments
const prodProvider = new azure.Provider("prod", {
    subscriptionId: "prod-subscription-id",
});

const devProvider = new azure.Provider("dev", {
    subscriptionId: "dev-subscription-id",
});
```

## Essential Resources

### Resource Group

```typescript
import * as azure from "@pulumi/azure-native";

const resourceGroup = new azure.resources.ResourceGroup("rg", {
    resourceGroupName: `rg-${pulumi.getProject()}-${pulumi.getStack()}`,
    location: "westeurope",
    tags: {
        Environment: pulumi.getStack(),
        ManagedBy: "Pulumi",
    },
});
```

### Storage Account - Secure Pattern

```typescript
import * as azure from "@pulumi/azure-native";

const storageAccount = new azure.storage.StorageAccount("storage", {
    resourceGroupName: resourceGroup.name,
    accountName: `st${pulumi.getProject()}${pulumi.getStack()}`.replace(/-/g, "").substring(0, 24),
    location: resourceGroup.location,

    // Use Standard_ZRS for zone redundancy
    sku: {
        name: azure.storage.SkuName.Standard_ZRS,
    },
    kind: azure.storage.Kind.StorageV2,

    // Security settings
    allowBlobPublicAccess: false,
    minimumTlsVersion: azure.storage.MinimumTlsVersion.TLS1_2,
    enableHttpsTrafficOnly: true,

    // Network security
    networkRuleSet: {
        defaultAction: azure.storage.DefaultAction.Deny,
        bypass: azure.storage.Bypass.AzureServices,
        virtualNetworkRules: [{
            virtualNetworkResourceId: subnet.id,
        }],
    },

    // Encryption
    encryption: {
        services: {
            blob: { enabled: true, keyType: azure.storage.KeyType.Account },
            file: { enabled: true, keyType: azure.storage.KeyType.Account },
        },
        keySource: azure.storage.KeySource.Microsoft_Storage,
    },

    tags: {
        Environment: pulumi.getStack(),
    },
});

// Blob container with private access
const container = new azure.storage.BlobContainer("data", {
    resourceGroupName: resourceGroup.name,
    accountName: storageAccount.name,
    containerName: "data",
    publicAccess: azure.storage.PublicAccess.None,
});

// Lifecycle management
const managementPolicy = new azure.storage.ManagementPolicy("lifecycle", {
    resourceGroupName: resourceGroup.name,
    accountName: storageAccount.name,
    managementPolicyName: "default",
    policy: {
        rules: [{
            name: "archive-old-blobs",
            enabled: true,
            type: azure.storage.RuleType.Lifecycle,
            definition: {
                actions: {
                    baseBlob: {
                        tierToCool: { daysAfterModificationGreaterThan: 30 },
                        tierToArchive: { daysAfterModificationGreaterThan: 90 },
                        delete: { daysAfterModificationGreaterThan: 365 },
                    },
                },
                filters: {
                    blobTypes: ["blockBlob"],
                },
            },
        }],
    },
});
```

### Virtual Network

```typescript
import * as azure from "@pulumi/azure-native";

const vnet = new azure.network.VirtualNetwork("vnet", {
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: `vnet-${pulumi.getStack()}`,
    location: resourceGroup.location,
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    tags: {
        Environment: pulumi.getStack(),
    },
});

// Public subnet (for load balancers, bastion)
const publicSubnet = new azure.network.Subnet("public", {
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: vnet.name,
    subnetName: "public",
    addressPrefix: "10.0.1.0/24",
});

// Private subnet (for apps)
const privateSubnet = new azure.network.Subnet("private", {
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: vnet.name,
    subnetName: "private",
    addressPrefix: "10.0.2.0/24",
    serviceEndpoints: [
        { service: "Microsoft.Storage" },
        { service: "Microsoft.Sql" },
        { service: "Microsoft.KeyVault" },
    ],
    delegations: [{
        name: "webapp",
        serviceName: "Microsoft.Web/serverFarms",
    }],
});

// Database subnet
const dbSubnet = new azure.network.Subnet("database", {
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: vnet.name,
    subnetName: "database",
    addressPrefix: "10.0.3.0/24",
    delegations: [{
        name: "flexibleServers",
        serviceName: "Microsoft.DBforPostgreSQL/flexibleServers",
    }],
});
```

### Azure SQL Database

```typescript
import * as azure from "@pulumi/azure-native";

const sqlServer = new azure.sql.Server("sql", {
    resourceGroupName: resourceGroup.name,
    serverName: `sql-${pulumi.getProject()}-${pulumi.getStack()}`,
    location: resourceGroup.location,

    // Use AAD authentication
    administrators: {
        administratorType: azure.sql.AdministratorType.ActiveDirectory,
        azureADOnlyAuthentication: true,
        login: "sql-admins",
        sid: aadGroupId,
        tenantId: tenantId,
    },

    // Security
    minimalTlsVersion: "1.2",
    publicNetworkAccess: azure.sql.ServerPublicNetworkAccess.Disabled,

    tags: {
        Environment: pulumi.getStack(),
    },
});

// Private endpoint for secure access
const privateEndpoint = new azure.network.PrivateEndpoint("sql-pe", {
    resourceGroupName: resourceGroup.name,
    privateEndpointName: "pe-sql",
    location: resourceGroup.location,
    subnet: { id: dbSubnet.id },
    privateLinkServiceConnections: [{
        name: "sql",
        privateLinkServiceId: sqlServer.id,
        groupIds: ["sqlServer"],
    }],
});

const database = new azure.sql.Database("db", {
    resourceGroupName: resourceGroup.name,
    serverName: sqlServer.name,
    databaseName: "app",

    // Serverless for dev, provisioned for prod
    sku: pulumi.getStack() === "prod" ? {
        name: "GP_Gen5",
        tier: "GeneralPurpose",
        family: "Gen5",
        capacity: 2,
    } : {
        name: "GP_S_Gen5",
        tier: "GeneralPurpose",
        family: "Gen5",
        capacity: 1,
    },

    autoPauseDelay: pulumi.getStack() !== "prod" ? 60 : -1,

    // Zone redundancy for prod
    zoneRedundant: pulumi.getStack() === "prod",

    // Backup
    requestedBackupStorageRedundancy: pulumi.getStack() === "prod"
        ? azure.sql.BackupStorageRedundancy.Geo
        : azure.sql.BackupStorageRedundancy.Local,

    tags: {
        Environment: pulumi.getStack(),
    },
});
```

### App Service / Function App

```typescript
import * as azure from "@pulumi/azure-native";

const appServicePlan = new azure.web.AppServicePlan("plan", {
    resourceGroupName: resourceGroup.name,
    name: `plan-${pulumi.getStack()}`,
    location: resourceGroup.location,
    kind: "Linux",
    reserved: true,
    sku: {
        name: pulumi.getStack() === "prod" ? "P1v3" : "B1",
        tier: pulumi.getStack() === "prod" ? "PremiumV3" : "Basic",
    },
});

const webapp = new azure.web.WebApp("app", {
    resourceGroupName: resourceGroup.name,
    name: `app-${pulumi.getProject()}-${pulumi.getStack()}`,
    location: resourceGroup.location,
    serverFarmId: appServicePlan.id,

    siteConfig: {
        linuxFxVersion: "NODE|18-lts",
        alwaysOn: pulumi.getStack() === "prod",
        http20Enabled: true,
        minTlsVersion: azure.web.SupportedTlsVersions.SupportedTlsVersions_1_2,
        ftpsState: azure.web.FtpsState.Disabled,

        // VNet integration
        vnetRouteAllEnabled: true,
    },

    httpsOnly: true,

    // Managed identity
    identity: {
        type: azure.web.ManagedServiceIdentityType.SystemAssigned,
    },

    // VNet integration
    virtualNetworkSubnetId: privateSubnet.id,

    tags: {
        Environment: pulumi.getStack(),
    },
});

// App settings from Key Vault
const appSettings = new azure.web.WebAppApplicationSettings("settings", {
    resourceGroupName: resourceGroup.name,
    name: webapp.name,
    properties: {
        "WEBSITE_NODE_DEFAULT_VERSION": "~18",
        "DB_CONNECTION_STRING": pulumi.interpolate`@Microsoft.KeyVault(VaultName=${keyVault.name};SecretName=db-connection-string)`,
    },
});
```

### Key Vault

```typescript
import * as azure from "@pulumi/azure-native";

const keyVault = new azure.keyvault.Vault("kv", {
    resourceGroupName: resourceGroup.name,
    vaultName: `kv-${pulumi.getProject()}-${pulumi.getStack()}`.substring(0, 24),
    location: resourceGroup.location,

    properties: {
        tenantId: tenantId,
        sku: {
            family: azure.keyvault.SkuFamily.A,
            name: azure.keyvault.SkuName.Standard,
        },

        // Use RBAC instead of access policies
        enableRbacAuthorization: true,

        // Security
        enableSoftDelete: true,
        softDeleteRetentionInDays: pulumi.getStack() === "prod" ? 90 : 7,
        enablePurgeProtection: pulumi.getStack() === "prod",

        // Network
        networkAcls: {
            defaultAction: azure.keyvault.NetworkRuleAction.Deny,
            bypass: azure.keyvault.NetworkRuleBypassOptions.AzureServices,
            virtualNetworkRules: [{
                id: privateSubnet.id,
            }],
        },
    },

    tags: {
        Environment: pulumi.getStack(),
    },
});

// Grant web app access to Key Vault via RBAC
const kvSecretUser = new azure.authorization.RoleAssignment("kv-secret-user", {
    principalId: webapp.identity.apply(id => id!.principalId),
    principalType: azure.authorization.PrincipalType.ServicePrincipal,
    roleDefinitionId: `/subscriptions/${subscriptionId}/providers/Microsoft.Authorization/roleDefinitions/4633458b-17de-408a-b874-0445c86b69e6`, // Key Vault Secrets User
    scope: keyVault.id,
});
```

### AKS - Kubernetes

```typescript
import * as azure from "@pulumi/azure-native";

const aks = new azure.containerservice.ManagedCluster("aks", {
    resourceGroupName: resourceGroup.name,
    resourceName: `aks-${pulumi.getStack()}`,
    location: resourceGroup.location,

    // Use managed identity
    identity: {
        type: azure.containerservice.ResourceIdentityType.SystemAssigned,
    },

    // AAD integration
    aadProfile: {
        managed: true,
        enableAzureRBAC: true,
    },

    // Network configuration
    networkProfile: {
        networkPlugin: azure.containerservice.NetworkPlugin.Azure,
        networkPolicy: azure.containerservice.NetworkPolicy.Calico,
        serviceCidr: "10.96.0.0/16",
        dnsServiceIP: "10.96.0.10",
    },

    // Node pools
    agentPoolProfiles: [{
        name: "system",
        mode: azure.containerservice.AgentPoolMode.System,
        count: 3,
        vmSize: "Standard_D4s_v3",
        vnetSubnetID: privateSubnet.id,
        availabilityZones: ["1", "2", "3"],
        enableAutoScaling: true,
        minCount: 3,
        maxCount: 5,
    }],

    // Security
    apiServerAccessProfile: {
        enablePrivateCluster: true,
    },
    disableLocalAccounts: true,

    // Addons
    addonProfiles: {
        azurePolicy: { enabled: true },
        omsAgent: {
            enabled: true,
            config: {
                logAnalyticsWorkspaceResourceID: logAnalytics.id,
            },
        },
    },

    autoUpgradeProfile: {
        upgradeChannel: azure.containerservice.UpgradeChannel.Stable,
    },

    tags: {
        Environment: pulumi.getStack(),
    },
});
```

## Security Best Practices

### Managed Identity

```typescript
// Always prefer managed identity over service principals
const identity = new azure.managedidentity.UserAssignedIdentity("app-identity", {
    resourceGroupName: resourceGroup.name,
    resourceName: `id-${pulumi.getProject()}-${pulumi.getStack()}`,
    location: resourceGroup.location,
});

// Assign roles using RBAC
const storageBlobReader = new azure.authorization.RoleAssignment("blob-reader", {
    principalId: identity.principalId,
    principalType: azure.authorization.PrincipalType.ServicePrincipal,
    roleDefinitionId: `/subscriptions/${subscriptionId}/providers/Microsoft.Authorization/roleDefinitions/2a2b9908-6ea1-4ae2-8e65-a410df84e7d1`, // Storage Blob Data Reader
    scope: storageAccount.id,
});
```

### Azure Policy

```typescript
// Use built-in policies via Azure Policy service
// Or create custom policies
const policy = new azure.authorization.PolicyAssignment("require-tags", {
    policyAssignmentName: "require-tags",
    scope: resourceGroup.id,
    policyDefinitionId: "/providers/Microsoft.Authorization/policyDefinitions/96670d01-0a4d-4649-9c89-2d3abc0a5025", // Require tag
    parameters: {
        tagName: { value: "Environment" },
    },
});
```

## Monitoring

```typescript
const logAnalytics = new azure.operationalinsights.Workspace("logs", {
    resourceGroupName: resourceGroup.name,
    workspaceName: `log-${pulumi.getStack()}`,
    location: resourceGroup.location,
    sku: { name: azure.operationalinsights.WorkspaceSkuNameEnum.PerGB2018 },
    retentionInDays: 30,
});

const appInsights = new azure.insights.Component("insights", {
    resourceGroupName: resourceGroup.name,
    resourceName: `appi-${pulumi.getStack()}`,
    location: resourceGroup.location,
    applicationType: azure.insights.ApplicationType.Web,
    kind: "web",
    workspaceResourceId: logAnalytics.id,
});
```

## Tagging Strategy

```typescript
// Standard tags for all Azure resources
const defaultTags = {
    Environment: pulumi.getStack(),
    Project: pulumi.getProject(),
    ManagedBy: "Pulumi",
    CostCenter: "engineering",
    Owner: "platform-team",
};

// Apply via stack transformation
pulumi.runtime.registerStackTransformation((args) => {
    if (args.props["tags"] !== undefined) {
        args.props["tags"] = { ...args.props["tags"], ...defaultTags };
    }
    return { props: args.props, opts: args.opts };
});
```
