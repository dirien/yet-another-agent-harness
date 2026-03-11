# Pulumi AWS Best Practices (Python)

## Provider Configuration

```python
import pulumi
import pulumi_aws as aws

# Use ESC for credentials via OIDC - avoid static credentials

# Multi-region deployments with explicit providers
us_east_1 = aws.Provider("us-east-1", region="us-east-1")
eu_west_1 = aws.Provider("eu-west-1", region="eu-west-1")

# Use provider with resources (ACM certs for CloudFront must be in us-east-1)
certificate = aws.acm.Certificate(
    "cert",
    domain_name="example.com",
    validation_method="DNS",
    opts=pulumi.ResourceOptions(provider=us_east_1),
)
```

## Essential Resources

### S3 - Secure Bucket Pattern

```python
import pulumi
import pulumi_aws as aws

bucket = aws.s3.Bucket(
    "data-bucket",
    # Enable versioning for data protection
    versioning={"enabled": True},
    # Server-side encryption
    server_side_encryption_configuration={
        "rule": {
            "apply_server_side_encryption_by_default": {
                "sse_algorithm": "aws:kms",
                "kms_master_key_id": kms_key.arn,
            },
            "bucket_key_enabled": True,
        },
    },
    tags={
        "Environment": pulumi.get_stack(),
        "ManagedBy": "Pulumi",
    },
)

# Block public access
public_access_block = aws.s3.BucketPublicAccessBlock(
    "block-public",
    bucket=bucket.id,
    block_public_acls=True,
    block_public_policy=True,
    ignore_public_acls=True,
    restrict_public_buckets=True,
)

# Lifecycle rules for cost optimization
lifecycle_config = aws.s3.BucketLifecycleConfigurationV2(
    "lifecycle",
    bucket=bucket.id,
    rules=[{
        "id": "transition-to-ia",
        "status": "Enabled",
        "transitions": [
            {"days": 30, "storage_class": "STANDARD_IA"},
            {"days": 90, "storage_class": "GLACIER"},
        ],
    }],
)
```

### VPC - Network Pattern

```python
import pulumi_awsx as awsx

# Use awsx for simplified VPC creation
vpc = awsx.ec2.Vpc(
    "main-vpc",
    cidr_block="10.0.0.0/16",
    number_of_availability_zones=3,
    nat_gateways=awsx.ec2.NatGatewayConfigurationArgs(
        strategy=awsx.ec2.NatGatewayStrategy.ONE_PER_AZ,
    ),
    subnet_specs=[
        awsx.ec2.SubnetSpecArgs(type=awsx.ec2.SubnetType.PUBLIC, cidr_mask=24),
        awsx.ec2.SubnetSpecArgs(type=awsx.ec2.SubnetType.PRIVATE, cidr_mask=24),
        awsx.ec2.SubnetSpecArgs(type=awsx.ec2.SubnetType.ISOLATED, cidr_mask=24),
    ],
    tags={"Name": "main-vpc"},
)

pulumi.export("vpc_id", vpc.vpc_id)
pulumi.export("public_subnet_ids", vpc.public_subnet_ids)
pulumi.export("private_subnet_ids", vpc.private_subnet_ids)
```

### RDS - Database Pattern

```python
import pulumi
import pulumi_aws as aws

config = pulumi.Config()
is_prod = pulumi.get_stack() == "prod"

db_subnet_group = aws.rds.SubnetGroup(
    "db-subnets",
    subnet_ids=vpc.isolated_subnet_ids,
)

database = aws.rds.Instance(
    "postgres",
    engine="postgres",
    engine_version="15.4",
    instance_class="db.t3.medium",
    allocated_storage=20,
    max_allocated_storage=100,

    db_name="myapp",
    username="admin",
    password=config.require_secret("db_password"),

    db_subnet_group_name=db_subnet_group.name,
    vpc_security_group_ids=[db_security_group.id],

    # High availability
    multi_az=is_prod,

    # Backup and maintenance
    backup_retention_period=7,
    backup_window="03:00-04:00",
    maintenance_window="Mon:04:00-Mon:05:00",

    # Security
    storage_encrypted=True,
    deletion_protection=is_prod,

    # Performance insights
    performance_insights_enabled=True,
    performance_insights_retention_period=7,

    skip_final_snapshot=not is_prod,
    final_snapshot_identifier=f"{pulumi.get_project()}-final-snapshot",

    tags={
        "Environment": pulumi.get_stack(),
    },
)
```

### Lambda Function

```python
import pulumi
import pulumi_aws as aws

role = aws.iam.Role(
    "lambda-role",
    assume_role_policy="""{
        "Version": "2012-10-17",
        "Statement": [{
            "Action": "sts:AssumeRole",
            "Principal": {"Service": "lambda.amazonaws.com"},
            "Effect": "Allow"
        }]
    }""",
)

aws.iam.RolePolicyAttachment(
    "lambda-basic",
    role=role.name,
    policy_arn=aws.iam.ManagedPolicy.AWS_LAMBDA_BASIC_EXECUTION_ROLE,
)

function = aws.lambda_.Function(
    "api-handler",
    runtime="nodejs18.x",
    handler="index.handler",
    code=pulumi.FileArchive("./lambda"),
    role=role.arn,
    memory_size=512,
    timeout=30,
    environment={
        "variables": {
            "TABLE_NAME": dynamo_table.name,
            "STAGE": pulumi.get_stack(),
        },
    },
    tracing_config={"mode": "Active"},
)
```

## Security Best Practices

### IAM - Least Privilege

```python
import pulumi
import pulumi_aws as aws

custom_policy = aws.iam.Policy(
    "app-policy",
    policy=bucket.arn.apply(lambda arn: f"""{{
        "Version": "2012-10-17",
        "Statement": [{{
            "Effect": "Allow",
            "Action": ["s3:GetObject", "s3:PutObject"],
            "Resource": ["{arn}/*"]
        }}]
    }}"""),
)
```

### Secrets Manager

```python
import pulumi
import pulumi_aws as aws

is_prod = pulumi.get_stack() == "prod"

secret = aws.secretsmanager.Secret(
    "api-key",
    name=f"{pulumi.get_stack()}/api-key",
    recovery_window_in_days=30 if is_prod else 0,
)
```

## Auto-Scaling

```python
import pulumi_aws as aws

scaling_target = aws.appautoscaling.Target(
    "ecs-scaling",
    max_capacity=10,
    min_capacity=2,
    resource_id=pulumi.Output.concat("service/", cluster.name, "/", service.name),
    scalable_dimension="ecs:service:DesiredCount",
    service_namespace="ecs",
)

scaling_policy = aws.appautoscaling.Policy(
    "cpu-scaling",
    policy_type="TargetTrackingScaling",
    resource_id=scaling_target.resource_id,
    scalable_dimension=scaling_target.scalable_dimension,
    service_namespace=scaling_target.service_namespace,
    target_tracking_scaling_policy_configuration={
        "target_value": 70,
        "predefined_metric_specification": {
            "predefined_metric_type": "ECSServiceAverageCPUUtilization",
        },
        "scale_in_cooldown": 300,
        "scale_out_cooldown": 60,
    },
)
```

## Tagging Strategy

```python
import pulumi

# Register transformation to add tags to all resources
def add_tags(args: pulumi.ResourceTransformationArgs):
    if args.props.get("tags") is not None:
        args.props["tags"]["Environment"] = pulumi.get_stack()
        args.props["tags"]["Project"] = pulumi.get_project()
        args.props["tags"]["ManagedBy"] = "Pulumi"
    return pulumi.ResourceTransformationResult(args.props, args.opts)

pulumi.runtime.register_stack_transformation(add_tags)
```
