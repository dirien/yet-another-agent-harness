# Pulumi Azure Best Practices (Python)

## Provider Selection

**IMPORTANT: Always use `azure-native` provider first.** The `azure-native` provider is auto-generated from Azure Resource Manager APIs and provides:
- 100% API coverage
- Same-day updates for new Azure features
- Full ARM template parity

Only use `pulumi-azure` (classic) for resources not yet in azure-native.

```python
# PREFERRED: azure-native provider
import pulumi_azure_native as azure

# FALLBACK: classic provider (only when needed)
import pulumi_azure as azure_classic
```

## Provider Configuration

```python
import pulumi
import pulumi_azure_native as azure

# Use ESC for credentials via OIDC

# Multi-subscription deployments
prod_provider = azure.Provider("prod", subscription_id="prod-subscription-id")
dev_provider = azure.Provider("dev", subscription_id="dev-subscription-id")
```

## Essential Resources

### Resource Group

```python
import pulumi
import pulumi_azure_native as azure

rg = azure.resources.ResourceGroup(
    "rg",
    resource_group_name=f"rg-{pulumi.get_project()}-{pulumi.get_stack()}",
    location="westeurope",
    tags={
        "Environment": pulumi.get_stack(),
        "ManagedBy": "Pulumi",
    },
)
```

### Storage Account - Secure Pattern

```python
import pulumi
import pulumi_azure_native as azure

storage_account = azure.storage.StorageAccount(
    "storage",
    resource_group_name=rg.name,
    account_name=f"st{pulumi.get_project()}{pulumi.get_stack()}".replace("-", "")[:24],
    location=rg.location,

    # Use Standard_ZRS for zone redundancy
    sku=azure.storage.SkuArgs(
        name=azure.storage.SkuName.STANDARD_ZRS,
    ),
    kind=azure.storage.Kind.STORAGE_V2,

    # Security settings
    allow_blob_public_access=False,
    minimum_tls_version=azure.storage.MinimumTlsVersion.TLS1_2,
    enable_https_traffic_only=True,

    # Network security
    network_rule_set=azure.storage.NetworkRuleSetArgs(
        default_action=azure.storage.DefaultAction.DENY,
        bypass=azure.storage.Bypass.AZURE_SERVICES,
    ),

    # Encryption
    encryption=azure.storage.EncryptionArgs(
        services=azure.storage.EncryptionServicesArgs(
            blob=azure.storage.EncryptionServiceArgs(
                enabled=True,
                key_type=azure.storage.KeyType.ACCOUNT,
            ),
        ),
        key_source=azure.storage.KeySource.MICROSOFT_STORAGE,
    ),

    tags={
        "Environment": pulumi.get_stack(),
    },
)
```

### Virtual Network

```python
import pulumi_azure_native as azure

vnet = azure.network.VirtualNetwork(
    "vnet",
    resource_group_name=rg.name,
    virtual_network_name=f"vnet-{pulumi.get_stack()}",
    location=rg.location,
    address_space=azure.network.AddressSpaceArgs(
        address_prefixes=["10.0.0.0/16"],
    ),
    tags={
        "Environment": pulumi.get_stack(),
    },
)

private_subnet = azure.network.Subnet(
    "private",
    resource_group_name=rg.name,
    virtual_network_name=vnet.name,
    subnet_name="private",
    address_prefix="10.0.2.0/24",
    service_endpoints=[
        azure.network.ServiceEndpointPropertiesFormatArgs(service="Microsoft.Storage"),
        azure.network.ServiceEndpointPropertiesFormatArgs(service="Microsoft.Sql"),
        azure.network.ServiceEndpointPropertiesFormatArgs(service="Microsoft.KeyVault"),
    ],
    private_endpoint_network_policies=azure.network.VirtualNetworkPrivateEndpointNetworkPolicies.ENABLED,
)
```

### Azure SQL Database

```python
import pulumi
import pulumi_azure_native as azure

is_prod = pulumi.get_stack() == "prod"

sql_server = azure.sql.Server(
    "sql",
    resource_group_name=rg.name,
    server_name=f"sql-{pulumi.get_project()}-{pulumi.get_stack()}",
    location=rg.location,

    # Use AAD authentication
    administrators=azure.sql.ServerExternalAdministratorArgs(
        administrator_type=azure.sql.AdministratorType.ACTIVE_DIRECTORY,
        azure_ad_only_authentication=True,
        login="sql-admins",
        sid=aad_group_id,
        tenant_id=tenant_id,
    ),

    minimal_tls_version="1.2",
    public_network_access=azure.sql.ServerPublicNetworkAccess.DISABLED,

    tags={
        "Environment": pulumi.get_stack(),
    },
)

database = azure.sql.Database(
    "db",
    resource_group_name=rg.name,
    server_name=sql_server.name,
    database_name="app",

    sku=azure.sql.SkuArgs(
        name="GP_Gen5" if is_prod else "GP_S_Gen5",
        tier="GeneralPurpose",
        family="Gen5",
        capacity=2 if is_prod else 1,
    ),

    zone_redundant=is_prod,

    tags={
        "Environment": pulumi.get_stack(),
    },
)
```

### Key Vault

```python
import pulumi
import pulumi_azure_native as azure

is_prod = pulumi.get_stack() == "prod"

key_vault = azure.keyvault.Vault(
    "kv",
    resource_group_name=rg.name,
    vault_name=f"kv-{pulumi.get_project()}-{pulumi.get_stack()}"[:24],
    location=rg.location,

    properties=azure.keyvault.VaultPropertiesArgs(
        tenant_id=tenant_id,
        sku=azure.keyvault.SkuArgs(
            family=azure.keyvault.SkuFamily.A,
            name=azure.keyvault.SkuName.STANDARD,
        ),

        # Use RBAC instead of access policies
        enable_rbac_authorization=True,

        # Security
        enable_soft_delete=True,
        soft_delete_retention_in_days=90 if is_prod else 7,
        enable_purge_protection=is_prod,

        # Network
        network_acls=azure.keyvault.NetworkRuleSetArgs(
            default_action=azure.keyvault.NetworkRuleAction.DENY,
            bypass=azure.keyvault.NetworkRuleBypassOptions.AZURE_SERVICES,
        ),
    ),

    tags={
        "Environment": pulumi.get_stack(),
    },
)
```

### App Service

```python
import pulumi
import pulumi_azure_native as azure

is_prod = pulumi.get_stack() == "prod"

app_service_plan = azure.web.AppServicePlan(
    "plan",
    resource_group_name=rg.name,
    name=f"plan-{pulumi.get_stack()}",
    location=rg.location,
    kind="Linux",
    reserved=True,
    sku=azure.web.SkuDescriptionArgs(
        name="P1v3" if is_prod else "B1",
        tier="PremiumV3" if is_prod else "Basic",
    ),
)

webapp = azure.web.WebApp(
    "app",
    resource_group_name=rg.name,
    name=f"app-{pulumi.get_project()}-{pulumi.get_stack()}",
    location=rg.location,
    server_farm_id=app_service_plan.id,

    site_config=azure.web.SiteConfigArgs(
        linux_fx_version="NODE|18-lts",
        always_on=is_prod,
        http20_enabled=True,
        min_tls_version=azure.web.SupportedTlsVersions.SUPPORTED_TLS_VERSIONS_1_2,
        ftps_state=azure.web.FtpsState.DISABLED,
        vnet_route_all_enabled=True,
    ),

    https_only=True,

    # Managed identity
    identity=azure.web.ManagedServiceIdentityArgs(
        type=azure.web.ManagedServiceIdentityType.SYSTEM_ASSIGNED,
    ),

    virtual_network_subnet_id=private_subnet.id,

    tags={
        "Environment": pulumi.get_stack(),
    },
)
```

## Security Best Practices

### Managed Identity

```python
import pulumi_azure_native as azure

# Create dedicated managed identity
identity = azure.managedidentity.UserAssignedIdentity(
    "app-identity",
    resource_group_name=rg.name,
    resource_name_=f"id-{pulumi.get_project()}-{pulumi.get_stack()}",
    location=rg.location,
)

# Assign roles using RBAC
blob_reader = azure.authorization.RoleAssignment(
    "blob-reader",
    principal_id=identity.principal_id,
    principal_type=azure.authorization.PrincipalType.SERVICE_PRINCIPAL,
    role_definition_id=f"/subscriptions/{subscription_id}/providers/Microsoft.Authorization/roleDefinitions/2a2b9908-6ea1-4ae2-8e65-a410df84e7d1",
    scope=storage_account.id,
)
```

## Monitoring

```python
import pulumi_azure_native as azure

log_analytics = azure.operationalinsights.Workspace(
    "logs",
    resource_group_name=rg.name,
    workspace_name=f"log-{pulumi.get_stack()}",
    location=rg.location,
    sku=azure.operationalinsights.WorkspaceSkuArgs(
        name=azure.operationalinsights.WorkspaceSkuNameEnum.PER_GB2018,
    ),
    retention_in_days=30,
)

app_insights = azure.insights.Component(
    "insights",
    resource_group_name=rg.name,
    resource_name_=f"appi-{pulumi.get_stack()}",
    location=rg.location,
    application_type=azure.insights.ApplicationType.WEB,
    kind="web",
    workspace_resource_id=log_analytics.id,
)
```

## Tagging Strategy

```python
import pulumi

default_tags = {
    "Environment": pulumi.get_stack(),
    "Project": pulumi.get_project(),
    "ManagedBy": "Pulumi",
    "CostCenter": "engineering",
}
```
