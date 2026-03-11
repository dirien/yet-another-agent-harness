# Pulumi GCP Best Practices (TypeScript)

## Provider Configuration

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";

// Use ESC for credentials via OIDC (Workload Identity Federation)
// ESC handles GOOGLE_CREDENTIALS or Workload Identity

// Multi-project deployments
const prodProvider = new gcp.Provider("prod", {
    project: "my-prod-project",
    region: "europe-west1",
});

const devProvider = new gcp.Provider("dev", {
    project: "my-dev-project",
    region: "europe-west1",
});
```

## Essential Resources

### Cloud Storage - Secure Bucket

```typescript
import * as gcp from "@pulumi/gcp";

const bucket = new gcp.storage.Bucket("data-bucket", {
    name: `${gcp.config.project}-data-${pulumi.getStack()}`,
    location: "EU",
    storageClass: "STANDARD",

    // Uniform bucket-level access (recommended)
    uniformBucketLevelAccess: true,

    // Versioning for data protection
    versioning: { enabled: true },

    // Encryption with CMEK (optional, default is Google-managed)
    encryption: {
        defaultKmsKeyName: kmsKey.id,
    },

    // Prevent public access
    publicAccessPrevention: "enforced",

    // Lifecycle rules
    lifecycleRules: [
        {
            action: { type: "SetStorageClass", storageClass: "NEARLINE" },
            condition: { age: 30 },
        },
        {
            action: { type: "SetStorageClass", storageClass: "COLDLINE" },
            condition: { age: 90 },
        },
        {
            action: { type: "SetStorageClass", storageClass: "ARCHIVE" },
            condition: { age: 365 },
        },
        {
            action: { type: "Delete" },
            condition: { age: 730 },
        },
    ],

    labels: {
        environment: pulumi.getStack(),
        managed_by: "pulumi",
    },
});

// IAM binding instead of ACLs
const bucketIamMember = new gcp.storage.BucketIAMMember("app-reader", {
    bucket: bucket.name,
    role: "roles/storage.objectViewer",
    member: pulumi.interpolate`serviceAccount:${serviceAccount.email}`,
});
```

### VPC Network

```typescript
import * as gcp from "@pulumi/gcp";

// Custom VPC (don't use default)
const vpc = new gcp.compute.Network("vpc", {
    name: `vpc-${pulumi.getStack()}`,
    autoCreateSubnetworks: false, // Always use custom subnets
    routingMode: "REGIONAL",
});

// Private subnet
const privateSubnet = new gcp.compute.Subnetwork("private", {
    name: `subnet-private-${pulumi.getStack()}`,
    ipCidrRange: "10.0.0.0/20",
    region: "europe-west1",
    network: vpc.id,

    // Enable Private Google Access
    privateIpGoogleAccess: true,

    // Enable VPC Flow Logs
    logConfig: {
        aggregationInterval: "INTERVAL_5_SEC",
        flowSampling: 0.5,
        metadata: "INCLUDE_ALL_METADATA",
    },

    // Secondary ranges for GKE
    secondaryIpRanges: [
        { rangeName: "pods", ipCidrRange: "10.4.0.0/14" },
        { rangeName: "services", ipCidrRange: "10.8.0.0/20" },
    ],
});

// Cloud NAT for private instances
const router = new gcp.compute.Router("router", {
    name: `router-${pulumi.getStack()}`,
    region: "europe-west1",
    network: vpc.id,
});

const nat = new gcp.compute.RouterNat("nat", {
    name: `nat-${pulumi.getStack()}`,
    router: router.name,
    region: "europe-west1",
    natIpAllocateOption: "AUTO_ONLY",
    sourceSubnetworkIpRangesToNat: "ALL_SUBNETWORKS_ALL_IP_RANGES",
    logConfig: {
        enable: true,
        filter: "ERRORS_ONLY",
    },
});
```

### Cloud SQL - PostgreSQL

```typescript
import * as gcp from "@pulumi/gcp";
import * as random from "@pulumi/random";

const dbPassword = new random.RandomPassword("db-password", {
    length: 24,
    special: true,
});

const sqlInstance = new gcp.sql.DatabaseInstance("postgres", {
    name: `sql-${pulumi.getProject()}-${pulumi.getStack()}`,
    region: "europe-west1",
    databaseVersion: "POSTGRES_15",

    settings: {
        tier: pulumi.getStack() === "prod" ? "db-custom-4-16384" : "db-f1-micro",

        // High availability for prod
        availabilityType: pulumi.getStack() === "prod" ? "REGIONAL" : "ZONAL",

        // Backups
        backupConfiguration: {
            enabled: true,
            startTime: "03:00",
            pointInTimeRecoveryEnabled: pulumi.getStack() === "prod",
            backupRetentionSettings: {
                retainedBackups: 7,
            },
        },

        // Network - use private IP
        ipConfiguration: {
            ipv4Enabled: false,
            privateNetwork: vpc.id,
            requireSsl: true,
        },

        // Maintenance
        maintenanceWindow: {
            day: 7, // Sunday
            hour: 3,
        },

        // Security
        databaseFlags: [
            { name: "log_checkpoints", value: "on" },
            { name: "log_connections", value: "on" },
            { name: "log_disconnections", value: "on" },
            { name: "log_lock_waits", value: "on" },
        ],

        userLabels: {
            environment: pulumi.getStack(),
        },
    },

    deletionProtection: pulumi.getStack() === "prod",
});

const database = new gcp.sql.Database("app", {
    name: "app",
    instance: sqlInstance.name,
});

const user = new gcp.sql.User("app-user", {
    name: "app",
    instance: sqlInstance.name,
    password: dbPassword.result,
});
```

### Cloud Run

```typescript
import * as gcp from "@pulumi/gcp";

const service = new gcp.cloudrun.Service("api", {
    name: `api-${pulumi.getStack()}`,
    location: "europe-west1",

    template: {
        spec: {
            serviceAccountName: serviceAccount.email,
            containers: [{
                image: `gcr.io/${gcp.config.project}/api:${imageTag}`,
                resources: {
                    limits: {
                        cpu: "1000m",
                        memory: "512Mi",
                    },
                },
                envs: [
                    { name: "ENVIRONMENT", value: pulumi.getStack() },
                    // Use Secret Manager for secrets
                    {
                        name: "DB_PASSWORD",
                        valueFrom: {
                            secretKeyRef: {
                                name: dbPasswordSecret.secretId,
                                key: "latest",
                            },
                        },
                    },
                ],
                ports: [{ containerPort: 8080 }],
            }],
            containerConcurrency: 80,
            timeoutSeconds: 300,
        },

        metadata: {
            annotations: {
                "autoscaling.knative.dev/minScale": pulumi.getStack() === "prod" ? "2" : "0",
                "autoscaling.knative.dev/maxScale": "100",
                "run.googleapis.com/vpc-access-connector": vpcConnector.id,
                "run.googleapis.com/vpc-access-egress": "private-ranges-only",
            },
        },
    },

    traffics: [{
        percent: 100,
        latestRevision: true,
    }],

    autogenerateRevisionName: true,
});

// Public access (for APIs)
const iamMember = new gcp.cloudrun.IamMember("public", {
    service: service.name,
    location: service.location,
    role: "roles/run.invoker",
    member: "allUsers",
});
```

### GKE - Kubernetes

```typescript
import * as gcp from "@pulumi/gcp";

const cluster = new gcp.container.Cluster("gke", {
    name: `gke-${pulumi.getStack()}`,
    location: "europe-west1", // Regional cluster

    // Use VPC-native cluster
    network: vpc.name,
    subnetwork: privateSubnet.name,
    ipAllocationPolicy: {
        clusterSecondaryRangeName: "pods",
        servicesSecondaryRangeName: "services",
    },

    // Private cluster
    privateClusterConfig: {
        enablePrivateNodes: true,
        enablePrivateEndpoint: pulumi.getStack() === "prod",
        masterIpv4CidrBlock: "172.16.0.0/28",
    },

    // Workload Identity (recommended)
    workloadIdentityConfig: {
        workloadPool: `${gcp.config.project}.svc.id.goog`,
    },

    // Security
    masterAuthorizedNetworksConfig: {
        cidrBlocks: [{
            cidrBlock: "10.0.0.0/8",
            displayName: "internal",
        }],
    },

    // Cluster addons
    addonsConfig: {
        httpLoadBalancing: { disabled: false },
        horizontalPodAutoscaling: { disabled: false },
        gcePersistentDiskCsiDriverConfig: { enabled: true },
        dnsCacheConfig: { enabled: true },
    },

    // Maintenance
    maintenancePolicy: {
        dailyMaintenanceWindow: {
            startTime: "03:00",
        },
    },

    // Release channel
    releaseChannel: {
        channel: "STABLE",
    },

    // Don't create default node pool
    removeDefaultNodePool: true,
    initialNodeCount: 1,

    resourceLabels: {
        environment: pulumi.getStack(),
    },
});

// Managed node pool
const nodePool = new gcp.container.NodePool("primary", {
    name: "primary",
    cluster: cluster.name,
    location: cluster.location,

    nodeCount: 1,
    autoscaling: {
        minNodeCount: 1,
        maxNodeCount: 10,
    },

    nodeConfig: {
        machineType: "e2-standard-4",
        diskSizeGb: 100,
        diskType: "pd-ssd",

        // Security
        serviceAccount: gkeNodeSa.email,
        oauthScopes: ["https://www.googleapis.com/auth/cloud-platform"],

        // Workload Identity
        workloadMetadataConfig: {
            mode: "GKE_METADATA",
        },

        // Shielded nodes
        shieldedInstanceConfig: {
            enableSecureBoot: true,
            enableIntegrityMonitoring: true,
        },

        labels: {
            environment: pulumi.getStack(),
        },
    },

    management: {
        autoRepair: true,
        autoUpgrade: true,
    },
});
```

### Cloud Functions

```typescript
import * as gcp from "@pulumi/gcp";

// Source code bucket
const sourceBucket = new gcp.storage.Bucket("functions-source", {
    name: `${gcp.config.project}-functions-source`,
    location: "EU",
    uniformBucketLevelAccess: true,
});

const sourceArchive = new gcp.storage.BucketObject("source", {
    bucket: sourceBucket.name,
    name: `function-${Date.now()}.zip`,
    source: new pulumi.asset.FileArchive("./function"),
});

const func = new gcp.cloudfunctions.Function("processor", {
    name: `processor-${pulumi.getStack()}`,
    region: "europe-west1",
    runtime: "nodejs18",

    sourceArchiveBucket: sourceBucket.name,
    sourceArchiveObject: sourceArchive.name,
    entryPoint: "handler",

    // Use service account
    serviceAccountEmail: serviceAccount.email,

    // VPC connector for private access
    vpcConnector: vpcConnector.id,
    vpcConnectorEgressSettings: "PRIVATE_RANGES_ONLY",

    // Environment
    environmentVariables: {
        PROJECT_ID: gcp.config.project!,
        ENVIRONMENT: pulumi.getStack(),
    },

    // Secrets from Secret Manager
    secretEnvironmentVariables: [{
        key: "API_KEY",
        projectId: gcp.config.project!,
        secret: apiKeySecret.secretId,
        version: "latest",
    }],

    availableMemoryMb: 256,
    timeout: 60,
    maxInstances: 100,
    minInstances: pulumi.getStack() === "prod" ? 1 : 0,

    // Trigger configuration
    triggerHttp: true,

    labels: {
        environment: pulumi.getStack(),
    },
});
```

## Security Best Practices

### Service Accounts & Workload Identity

```typescript
// Create dedicated service account for each workload
const serviceAccount = new gcp.serviceaccount.Account("app-sa", {
    accountId: `sa-app-${pulumi.getStack()}`,
    displayName: "Application Service Account",
});

// Grant minimal permissions
const storageBinding = new gcp.projects.IAMMember("storage-access", {
    project: gcp.config.project!,
    role: "roles/storage.objectViewer",
    member: pulumi.interpolate`serviceAccount:${serviceAccount.email}`,
});

// Workload Identity binding for GKE
const workloadIdentityBinding = new gcp.serviceaccount.IAMBinding("workload-identity", {
    serviceAccountId: serviceAccount.name,
    role: "roles/iam.workloadIdentityUser",
    members: [
        pulumi.interpolate`serviceAccount:${gcp.config.project}.svc.id.goog[default/app]`,
    ],
});
```

### Secret Manager

```typescript
const secret = new gcp.secretmanager.Secret("api-key", {
    secretId: `api-key-${pulumi.getStack()}`,
    replication: {
        automatic: true,
    },
    labels: {
        environment: pulumi.getStack(),
    },
});

const secretVersion = new gcp.secretmanager.SecretVersion("api-key-v1", {
    secret: secret.id,
    secretData: config.requireSecret("apiKey"),
});

// Grant access to service account
const secretAccess = new gcp.secretmanager.SecretIamMember("app-access", {
    secretId: secret.secretId,
    role: "roles/secretmanager.secretAccessor",
    member: pulumi.interpolate`serviceAccount:${serviceAccount.email}`,
});
```

### VPC Service Controls (for sensitive data)

```typescript
const accessPolicy = new gcp.accesscontextmanager.AccessPolicy("policy", {
    parent: `organizations/${orgId}`,
    title: "Data Perimeter",
});

const servicePerimeter = new gcp.accesscontextmanager.ServicePerimeter("perimeter", {
    parent: pulumi.interpolate`accessPolicies/${accessPolicy.name}`,
    name: pulumi.interpolate`accessPolicies/${accessPolicy.name}/servicePerimeters/data_perimeter`,
    title: "Data Perimeter",
    status: {
        resources: [
            `projects/${projectNumber}`,
        ],
        restrictedServices: [
            "storage.googleapis.com",
            "bigquery.googleapis.com",
        ],
    },
});
```

## Monitoring

```typescript
// Enable all audit logs
const auditConfig = new gcp.projects.IAMAuditConfig("audit", {
    project: gcp.config.project!,
    service: "allServices",
    auditLogConfigs: [
        { logType: "ADMIN_READ" },
        { logType: "DATA_READ" },
        { logType: "DATA_WRITE" },
    ],
});

// Uptime check
const uptimeCheck = new gcp.monitoring.UptimeCheckConfig("api-check", {
    displayName: "API Health Check",
    timeout: "10s",
    period: "60s",
    httpCheck: {
        path: "/health",
        port: 443,
        useSsl: true,
    },
    monitoredResource: {
        type: "uptime_url",
        labels: {
            project_id: gcp.config.project!,
            host: "api.example.com",
        },
    },
});

// Alert policy
const alertPolicy = new gcp.monitoring.AlertPolicy("high-latency", {
    displayName: "High API Latency",
    combiner: "OR",
    conditions: [{
        displayName: "Latency > 2s",
        conditionThreshold: {
            filter: 'resource.type = "cloud_run_revision" AND metric.type = "run.googleapis.com/request_latencies"',
            comparison: "COMPARISON_GT",
            thresholdValue: 2000,
            duration: "300s",
            aggregations: [{
                alignmentPeriod: "60s",
                perSeriesAligner: "ALIGN_PERCENTILE_99",
            }],
        },
    }],
    notificationChannels: [notificationChannel.id],
});
```

## Labels Strategy

```typescript
// Standard labels for all GCP resources
const defaultLabels = {
    environment: pulumi.getStack(),
    project: pulumi.getProject().toLowerCase().replace(/[^a-z0-9-]/g, "-"),
    managed_by: "pulumi",
    cost_center: "engineering",
    team: "platform",
};

// Note: GCP labels must be lowercase and can only contain
// lowercase letters, numeric characters, underscores, and dashes
```
