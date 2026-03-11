# Pulumi AWS Best Practices (TypeScript)

## Provider Configuration

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

// Use ESC for credentials via OIDC - avoid static credentials
// ESC environment handles AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN

// Multi-region deployments with explicit providers
const usEast1 = new aws.Provider("us-east-1", { region: "us-east-1" });
const euWest1 = new aws.Provider("eu-west-1", { region: "eu-west-1" });

// Use provider with resources
const certificate = new aws.acm.Certificate("cert", {
    domainName: "example.com",
    validationMethod: "DNS",
}, { provider: usEast1 }); // ACM certs for CloudFront must be in us-east-1
```

## Essential Resources

### S3 - Secure Bucket Pattern

```typescript
import * as aws from "@pulumi/aws";

const bucket = new aws.s3.Bucket("data-bucket", {
    // Enable versioning for data protection
    versioning: { enabled: true },

    // Server-side encryption
    serverSideEncryptionConfiguration: {
        rule: {
            applyServerSideEncryptionByDefault: {
                sseAlgorithm: "aws:kms",
                kmsMasterKeyId: kmsKey.arn,
            },
            bucketKeyEnabled: true,
        },
    },

    // Block public access
    // Note: Use aws.s3.BucketPublicAccessBlock for newer approach
    tags: {
        Environment: pulumi.getStack(),
        ManagedBy: "Pulumi",
    },
});

// Explicitly block public access
const publicAccessBlock = new aws.s3.BucketPublicAccessBlock("block-public", {
    bucket: bucket.id,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
});

// Lifecycle rules for cost optimization
const lifecycleConfig = new aws.s3.BucketLifecycleConfigurationV2("lifecycle", {
    bucket: bucket.id,
    rules: [{
        id: "transition-to-ia",
        status: "Enabled",
        transitions: [{
            days: 30,
            storageClass: "STANDARD_IA",
        }, {
            days: 90,
            storageClass: "GLACIER",
        }],
    }],
});
```

### Lambda - Function Patterns

```typescript
import * as aws from "@pulumi/aws";
import * as pulumi from "@pulumi/pulumi";

// Inline function with CallbackFunction (quick prototyping)
const inlineFunction = new aws.lambda.CallbackFunction("processor", {
    callback: async (event: aws.s3.BucketEvent) => {
        console.log("Processing:", JSON.stringify(event));
        return { statusCode: 200 };
    },
    runtime: aws.lambda.Runtime.NodeJS18dX,
    memorySize: 256,
    timeout: 30,
});

// Production function with proper IAM
const role = new aws.iam.Role("lambda-role", {
    assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({
        Service: "lambda.amazonaws.com",
    }),
});

new aws.iam.RolePolicyAttachment("lambda-basic", {
    role: role.name,
    policyArn: aws.iam.ManagedPolicy.AWSLambdaBasicExecutionRole,
});

const prodFunction = new aws.lambda.Function("api-handler", {
    runtime: aws.lambda.Runtime.NodeJS18dX,
    handler: "index.handler",
    code: new pulumi.asset.AssetArchive({
        ".": new pulumi.asset.FileArchive("./lambda"),
    }),
    role: role.arn,
    memorySize: 512,
    timeout: 30,
    environment: {
        variables: {
            TABLE_NAME: dynamoTable.name,
            STAGE: pulumi.getStack(),
        },
    },
    // Enable X-Ray tracing
    tracingConfig: { mode: "Active" },
});
```

### VPC - Network Pattern

```typescript
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";

// Use awsx for simplified VPC creation
const vpc = new awsx.ec2.Vpc("main-vpc", {
    cidrBlock: "10.0.0.0/16",
    numberOfAvailabilityZones: 3,
    natGateways: {
        strategy: awsx.ec2.NatGatewayStrategy.OnePerAz, // HA
    },
    subnetSpecs: [
        { type: awsx.ec2.SubnetType.Public, cidrMask: 24 },
        { type: awsx.ec2.SubnetType.Private, cidrMask: 24 },
        { type: awsx.ec2.SubnetType.Isolated, cidrMask: 24 }, // For databases
    ],
    tags: { Name: "main-vpc" },
});

export const vpcId = vpc.vpcId;
export const publicSubnetIds = vpc.publicSubnetIds;
export const privateSubnetIds = vpc.privateSubnetIds;
```

### RDS - Database Pattern

```typescript
import * as aws from "@pulumi/aws";

const dbSubnetGroup = new aws.rds.SubnetGroup("db-subnets", {
    subnetIds: vpc.isolatedSubnetIds, // Use isolated subnets
});

const dbSecurityGroup = new aws.ec2.SecurityGroup("db-sg", {
    vpcId: vpc.vpcId,
    ingress: [{
        protocol: "tcp",
        fromPort: 5432,
        toPort: 5432,
        securityGroups: [appSecurityGroup.id], // Only from app
    }],
    egress: [{
        protocol: "-1",
        fromPort: 0,
        toPort: 0,
        cidrBlocks: ["0.0.0.0/0"],
    }],
});

const database = new aws.rds.Instance("postgres", {
    engine: "postgres",
    engineVersion: "15.4",
    instanceClass: "db.t3.medium",
    allocatedStorage: 20,
    maxAllocatedStorage: 100, // Auto-scaling

    dbName: "myapp",
    username: "admin",
    password: config.requireSecret("dbPassword"), // From ESC

    dbSubnetGroupName: dbSubnetGroup.name,
    vpcSecurityGroupIds: [dbSecurityGroup.id],

    // High availability
    multiAz: pulumi.getStack() === "prod",

    // Backup and maintenance
    backupRetentionPeriod: 7,
    backupWindow: "03:00-04:00",
    maintenanceWindow: "Mon:04:00-Mon:05:00",

    // Security
    storageEncrypted: true,
    deletionProtection: pulumi.getStack() === "prod",

    // Performance insights
    performanceInsightsEnabled: true,
    performanceInsightsRetentionPeriod: 7,

    skipFinalSnapshot: pulumi.getStack() !== "prod",
    finalSnapshotIdentifier: `${pulumi.getProject()}-final-snapshot`,

    tags: {
        Environment: pulumi.getStack(),
    },
});
```

### ECS Fargate - Container Pattern

```typescript
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";

const cluster = new aws.ecs.Cluster("app-cluster", {
    settings: [{
        name: "containerInsights",
        value: "enabled",
    }],
});

const service = new awsx.ecs.FargateService("api", {
    cluster: cluster.arn,
    networkConfiguration: {
        subnets: vpc.privateSubnetIds,
        securityGroups: [serviceSecurityGroup.id],
    },
    desiredCount: 2,
    taskDefinitionArgs: {
        container: {
            name: "api",
            image: ecrImage.imageUri,
            cpu: 256,
            memory: 512,
            essential: true,
            portMappings: [{
                containerPort: 8080,
                targetGroup: targetGroup,
            }],
            environment: [
                { name: "STAGE", value: pulumi.getStack() },
            ],
            secrets: [
                {
                    name: "DB_PASSWORD",
                    valueFrom: dbPasswordSecret.arn,
                },
            ],
            logConfiguration: {
                logDriver: "awslogs",
                options: {
                    "awslogs-group": logGroup.name,
                    "awslogs-region": aws.config.region,
                    "awslogs-stream-prefix": "api",
                },
            },
        },
    },
});
```

## Security Best Practices

### IAM - Least Privilege

```typescript
// Use managed policies where appropriate
const readOnlyPolicy = aws.iam.ManagedPolicy.AmazonS3ReadOnlyAccess;

// Create custom policies for specific needs
const customPolicy = new aws.iam.Policy("app-policy", {
    policy: pulumi.jsonStringify({
        Version: "2012-10-17",
        Statement: [{
            Effect: "Allow",
            Action: [
                "s3:GetObject",
                "s3:PutObject",
            ],
            Resource: [
                pulumi.interpolate`${bucket.arn}/*`,
            ],
            Condition: {
                StringEquals: {
                    "s3:x-amz-server-side-encryption": "aws:kms",
                },
            },
        }],
    }),
});
```

### Secrets Manager

```typescript
const secret = new aws.secretsmanager.Secret("api-key", {
    name: `${pulumi.getStack()}/api-key`,
    recoveryWindowInDays: pulumi.getStack() === "prod" ? 30 : 0,
});

// Rotate secrets automatically
const rotation = new aws.secretsmanager.SecretRotation("api-key-rotation", {
    secretId: secret.id,
    rotationLambdaArn: rotationLambda.arn,
    rotationRules: {
        automaticallyAfterDays: 30,
    },
});
```

## Cost Optimization

```typescript
// Use Savings Plans / Reserved capacity for predictable workloads
// Configure auto-scaling
const scalingTarget = new aws.appautoscaling.Target("ecs-scaling", {
    maxCapacity: 10,
    minCapacity: 2,
    resourceId: pulumi.interpolate`service/${cluster.name}/${service.name}`,
    scalableDimension: "ecs:service:DesiredCount",
    serviceNamespace: "ecs",
});

const scalingPolicy = new aws.appautoscaling.Policy("cpu-scaling", {
    policyType: "TargetTrackingScaling",
    resourceId: scalingTarget.resourceId,
    scalableDimension: scalingTarget.scalableDimension,
    serviceNamespace: scalingTarget.serviceNamespace,
    targetTrackingScalingPolicyConfiguration: {
        targetValue: 70,
        predefinedMetricSpecification: {
            predefinedMetricType: "ECSServiceAverageCPUUtilization",
        },
        scaleInCooldown: 300,
        scaleOutCooldown: 60,
    },
});

// S3 Intelligent Tiering for unknown access patterns
const intelligentTiering = new aws.s3.BucketIntelligentTieringConfiguration("tiering", {
    bucket: bucket.id,
    tierings: [
        { accessTier: "ARCHIVE_ACCESS", days: 90 },
        { accessTier: "DEEP_ARCHIVE_ACCESS", days: 180 },
    ],
});
```

## Monitoring & Observability

```typescript
// CloudWatch Alarms
const cpuAlarm = new aws.cloudwatch.MetricAlarm("high-cpu", {
    comparisonOperator: "GreaterThanThreshold",
    evaluationPeriods: 2,
    metricName: "CPUUtilization",
    namespace: "AWS/ECS",
    period: 300,
    statistic: "Average",
    threshold: 80,
    alarmActions: [snsTopic.arn],
    dimensions: {
        ClusterName: cluster.name,
        ServiceName: service.name,
    },
});

// X-Ray tracing enabled on Lambda (see above)

// CloudWatch Logs with retention
const logGroup = new aws.cloudwatch.LogGroup("app-logs", {
    name: `/aws/ecs/${pulumi.getStack()}/app`,
    retentionInDays: 30,
});
```

## Tagging Strategy

```typescript
import * as pulumi from "@pulumi/pulumi";

// Register transformation to add tags to all resources
pulumi.runtime.registerStackTransformation((args) => {
    if (args.props["tags"] !== undefined) {
        args.props["tags"] = {
            ...args.props["tags"],
            Environment: pulumi.getStack(),
            Project: pulumi.getProject(),
            ManagedBy: "Pulumi",
            CostCenter: "engineering",
        };
    }
    return { props: args.props, opts: args.opts };
});
```
