# Pulumi GCP Best Practices (Python)

## Provider Configuration

```python
import pulumi
import pulumi_gcp as gcp

# Use ESC for credentials via OIDC (Workload Identity Federation)

# Multi-project deployments
prod_provider = gcp.Provider("prod", project="my-prod-project", region="europe-west1")
dev_provider = gcp.Provider("dev", project="my-dev-project", region="europe-west1")
```

## Essential Resources

### Cloud Storage - Secure Bucket

```python
import pulumi
import pulumi_gcp as gcp

bucket = gcp.storage.Bucket(
    "data-bucket",
    name=f"{gcp.config.project}-data-{pulumi.get_stack()}",
    location="EU",
    storage_class="STANDARD",

    # Uniform bucket-level access (recommended)
    uniform_bucket_level_access=True,

    # Versioning for data protection
    versioning={"enabled": True},

    # Prevent public access
    public_access_prevention="enforced",

    # Lifecycle rules
    lifecycle_rules=[
        {
            "action": {"type": "SetStorageClass", "storage_class": "NEARLINE"},
            "condition": {"age": 30},
        },
        {
            "action": {"type": "SetStorageClass", "storage_class": "COLDLINE"},
            "condition": {"age": 90},
        },
        {
            "action": {"type": "Delete"},
            "condition": {"age": 730},
        },
    ],

    labels={
        "environment": pulumi.get_stack(),
        "managed_by": "pulumi",
    },
)

# IAM binding instead of ACLs
bucket_iam_member = gcp.storage.BucketIAMMember(
    "app-reader",
    bucket=bucket.name,
    role="roles/storage.objectViewer",
    member=service_account.email.apply(lambda email: f"serviceAccount:{email}"),
)
```

### VPC Network

```python
import pulumi_gcp as gcp

# Custom VPC (don't use default)
vpc = gcp.compute.Network(
    "vpc",
    name=f"vpc-{pulumi.get_stack()}",
    auto_create_subnetworks=False,  # Always use custom subnets
    routing_mode="REGIONAL",
)

# Private subnet
private_subnet = gcp.compute.Subnetwork(
    "private",
    name=f"subnet-private-{pulumi.get_stack()}",
    ip_cidr_range="10.0.0.0/20",
    region="europe-west1",
    network=vpc.id,

    # Enable Private Google Access
    private_ip_google_access=True,

    # Enable VPC Flow Logs
    log_config={
        "aggregation_interval": "INTERVAL_5_SEC",
        "flow_sampling": 0.5,
        "metadata": "INCLUDE_ALL_METADATA",
    },

    # Secondary ranges for GKE
    secondary_ip_ranges=[
        {"range_name": "pods", "ip_cidr_range": "10.4.0.0/14"},
        {"range_name": "services", "ip_cidr_range": "10.8.0.0/20"},
    ],
)

# Cloud NAT for private instances
router = gcp.compute.Router(
    "router",
    name=f"router-{pulumi.get_stack()}",
    region="europe-west1",
    network=vpc.id,
)

nat = gcp.compute.RouterNat(
    "nat",
    name=f"nat-{pulumi.get_stack()}",
    router=router.name,
    region="europe-west1",
    nat_ip_allocate_option="AUTO_ONLY",
    source_subnetwork_ip_ranges_to_nat="ALL_SUBNETWORKS_ALL_IP_RANGES",
    log_config={
        "enable": True,
        "filter": "ERRORS_ONLY",
    },
)
```

### Cloud SQL - PostgreSQL

```python
import pulumi
import pulumi_gcp as gcp
import pulumi_random as random

db_password = random.RandomPassword(
    "db-password",
    length=24,
    special=True,
)

is_prod = pulumi.get_stack() == "prod"

sql_instance = gcp.sql.DatabaseInstance(
    "postgres",
    name=f"sql-{pulumi.get_project()}-{pulumi.get_stack()}",
    region="europe-west1",
    database_version="POSTGRES_15",

    settings={
        "tier": "db-custom-4-16384" if is_prod else "db-f1-micro",

        # High availability for prod
        "availability_type": "REGIONAL" if is_prod else "ZONAL",

        # Backups
        "backup_configuration": {
            "enabled": True,
            "start_time": "03:00",
            "point_in_time_recovery_enabled": is_prod,
            "backup_retention_settings": {
                "retained_backups": 7,
            },
        },

        # Network - use private IP
        "ip_configuration": {
            "ipv4_enabled": False,
            "private_network": vpc.id,
            "require_ssl": True,
        },

        # Maintenance
        "maintenance_window": {
            "day": 7,
            "hour": 3,
        },

        # Security flags
        "database_flags": [
            {"name": "log_checkpoints", "value": "on"},
            {"name": "log_connections", "value": "on"},
        ],

        "user_labels": {
            "environment": pulumi.get_stack(),
        },
    },

    deletion_protection=is_prod,
)

database = gcp.sql.Database(
    "app",
    name="app",
    instance=sql_instance.name,
)

user = gcp.sql.User(
    "app-user",
    name="app",
    instance=sql_instance.name,
    password=db_password.result,
)
```

### Cloud Run

```python
import pulumi
import pulumi_gcp as gcp

is_prod = pulumi.get_stack() == "prod"

service = gcp.cloudrun.Service(
    "api",
    name=f"api-{pulumi.get_stack()}",
    location="europe-west1",

    template={
        "spec": {
            "service_account_name": service_account.email,
            "containers": [{
                "image": f"gcr.io/{gcp.config.project}/api:{image_tag}",
                "resources": {
                    "limits": {
                        "cpu": "1000m",
                        "memory": "512Mi",
                    },
                },
                "envs": [
                    {"name": "ENVIRONMENT", "value": pulumi.get_stack()},
                ],
                "ports": [{"container_port": 8080}],
            }],
            "container_concurrency": 80,
            "timeout_seconds": 300,
        },

        "metadata": {
            "annotations": {
                "autoscaling.knative.dev/minScale": "2" if is_prod else "0",
                "autoscaling.knative.dev/maxScale": "100",
            },
        },
    },

    traffics=[{
        "percent": 100,
        "latest_revision": True,
    }],

    autogenerate_revision_name=True,
)

# Public access (for APIs)
iam_member = gcp.cloudrun.IamMember(
    "public",
    service=service.name,
    location=service.location,
    role="roles/run.invoker",
    member="allUsers",
)
```

## Security Best Practices

### Service Accounts & Workload Identity

```python
import pulumi_gcp as gcp

# Create dedicated service account for each workload
sa = gcp.serviceaccount.Account(
    "app-sa",
    account_id=f"sa-app-{pulumi.get_stack()}",
    display_name="Application Service Account",
)

# Grant minimal permissions
storage_binding = gcp.projects.IAMMember(
    "storage-access",
    project=gcp.config.project,
    role="roles/storage.objectViewer",
    member=sa.email.apply(lambda email: f"serviceAccount:{email}"),
)

# Workload Identity binding for GKE
workload_identity_binding = gcp.serviceaccount.IAMBinding(
    "workload-identity",
    service_account_id=sa.name,
    role="roles/iam.workloadIdentityUser",
    members=[
        pulumi.Output.concat(
            f"serviceAccount:{gcp.config.project}.svc.id.goog[default/app]"
        ),
    ],
)
```

### Secret Manager

```python
import pulumi
import pulumi_gcp as gcp

config = pulumi.Config()

secret = gcp.secretmanager.Secret(
    "api-key",
    secret_id=f"api-key-{pulumi.get_stack()}",
    replication={"automatic": True},
    labels={
        "environment": pulumi.get_stack(),
    },
)

secret_version = gcp.secretmanager.SecretVersion(
    "api-key-v1",
    secret=secret.id,
    secret_data=config.require_secret("api_key"),
)

# Grant access to service account
secret_access = gcp.secretmanager.SecretIamMember(
    "app-access",
    secret_id=secret.secret_id,
    role="roles/secretmanager.secretAccessor",
    member=sa.email.apply(lambda email: f"serviceAccount:{email}"),
)
```

## Monitoring

```python
import pulumi_gcp as gcp

# Enable all audit logs
audit_config = gcp.projects.IAMAuditConfig(
    "audit",
    project=gcp.config.project,
    service="allServices",
    audit_log_configs=[
        {"log_type": "ADMIN_READ"},
        {"log_type": "DATA_READ"},
        {"log_type": "DATA_WRITE"},
    ],
)

# Uptime check
uptime_check = gcp.monitoring.UptimeCheckConfig(
    "api-check",
    display_name="API Health Check",
    timeout="10s",
    period="60s",
    http_check={
        "path": "/health",
        "port": 443,
        "use_ssl": True,
    },
    monitored_resource={
        "type": "uptime_url",
        "labels": {
            "project_id": gcp.config.project,
            "host": "api.example.com",
        },
    },
)
```

## Labels Strategy

```python
import pulumi
import re

# Standard labels for all GCP resources
# Note: GCP labels must be lowercase with letters, numbers, underscores, dashes
default_labels = {
    "environment": pulumi.get_stack(),
    "project": re.sub(r"[^a-z0-9-]", "-", pulumi.get_project().lower()),
    "managed_by": "pulumi",
    "cost_center": "engineering",
}
```
