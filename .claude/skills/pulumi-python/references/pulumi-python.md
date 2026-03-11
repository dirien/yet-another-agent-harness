# Pulumi Python Reference

## Project Setup

### requirements.txt

```
pulumi>=3.0.0,<4.0.0
pulumi-aws>=6.0.0,<7.0.0
```

### pyproject.toml (for poetry/uv)

```toml
[project]
name = "my-pulumi-project"
version = "0.1.0"
requires-python = ">=3.9"
dependencies = [
    "pulumi>=3.0.0,<4.0.0",
    "pulumi-aws>=6.0.0,<7.0.0",
]

[tool.poetry.dependencies]
python = "^3.9"
pulumi = "^3.0.0"
pulumi-aws = "^6.0.0"
```

### Pulumi.yaml

```yaml
name: my-project
description: My Pulumi project
runtime:
  name: python
  options:
    toolchain: pip
    virtualenv: venv
```

### __main__.py Template

```python
"""Pulumi infrastructure program."""
import pulumi
import pulumi_aws as aws

# Your infrastructure code here
```

## Working with Inputs and Outputs

### Output Transformations

```python
import pulumi

# Single output transformation
url = bucket.bucket_domain_name.apply(lambda domain: f"https://{domain}")

# Multiple outputs with Output.all
combined = pulumi.Output.all(bucket.id, bucket.arn).apply(
    lambda args: {"id": args[0], "arn": args[1]}
)

# String interpolation with Output.concat
message = pulumi.Output.concat("Bucket ARN: ", bucket.arn)

# Format string style
formatted = pulumi.Output.format("Bucket {0} has ARN {1}", bucket.id, bucket.arn)

# Conditional output
endpoint = bucket.website_endpoint.apply(lambda ep: ep or "default-endpoint")
```

### Input Types

```python
import pulumi

# Basic inputs are automatically wrapped
bucket = aws.s3.Bucket("bucket", tags={"key": "value"})

# Explicit Output wrapping
bucket = aws.s3.Bucket(
    "bucket",
    tags=pulumi.Output.from_input({"key": "value"}),
)

# Using other resource outputs as inputs
policy = aws.s3.BucketPolicy(
    "policy",
    bucket=bucket.id,  # Output used as Input
)
```

## Configuration

```python
import pulumi

# Default namespace (project name)
config = pulumi.Config()

# Named namespace
aws_config = pulumi.Config("aws")

# Required values
region = aws_config.require("region")

# Optional with default
instance_type = config.get("instance_type") or "t3.small"

# Boolean
debug = config.get_bool("debug") or False

# Integer
port = config.get_int("port") or 8080

# Float
threshold = config.get_float("threshold") or 0.8

# Secrets (returns Output[str])
api_key = config.require_secret("api_key")

# Objects
db_config = config.require_object("database")
# Returns: {"host": "...", "port": 5432}

# Optional object
optional_config = config.get_object("optional") or {}
```

## Provider Configuration

```python
import pulumi
import pulumi_aws as aws

# Explicit provider for different region
us_east_1 = aws.Provider("us-east-1", region="us-east-1")

# Use provider with resource
cert = aws.acm.Certificate(
    "cert",
    domain_name="example.com",
    opts=pulumi.ResourceOptions(provider=us_east_1),
)
```

## Data Sources

```python
import pulumi_aws as aws

# Get availability zones
zones = aws.get_availability_zones(state="available")

# Use in resources
subnets = [
    aws.ec2.Subnet(
        f"subnet-{i}",
        vpc_id=vpc.id,
        availability_zone=az,
        cidr_block=f"10.0.{i}.0/24",
    )
    for i, az in enumerate(zones.names)
]

# Async data source with Output
ami = aws.ec2.get_ami_output(
    most_recent=True,
    owners=["amazon"],
    filters=[
        aws.ec2.GetAmiFilterArgs(
            name="name",
            values=["amzn2-ami-hvm-*-x86_64-gp2"],
        ),
    ],
)

instance = aws.ec2.Instance("instance", ami=ami.id, instance_type="t3.micro")
```

## Type Hints

```python
from typing import Optional, Sequence
import pulumi
import pulumi_aws as aws


def create_vpc(
    name: str,
    cidr_block: pulumi.Input[str],
    az_count: int = 2,
) -> aws.ec2.Vpc:
    """Create a VPC with the given parameters."""
    return aws.ec2.Vpc(
        name,
        cidr_block=cidr_block,
        enable_dns_hostnames=True,
        tags={"Name": name},
    )


# With return type annotation
def get_subnets(vpc: aws.ec2.Vpc) -> Sequence[aws.ec2.Subnet]:
    # ...
    pass
```

## Args Classes vs Dict Literals

```python
import pulumi_aws as aws

# Using Args classes (full type checking)
bucket = aws.s3.Bucket(
    "bucket",
    versioning=aws.s3.BucketVersioningArgs(
        enabled=True,
    ),
    server_side_encryption_configuration=aws.s3.BucketServerSideEncryptionConfigurationArgs(
        rule=aws.s3.BucketServerSideEncryptionConfigurationRuleArgs(
            apply_server_side_encryption_by_default=aws.s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs(
                sse_algorithm="AES256",
            ),
        ),
    ),
)

# Using dict literals (concise)
bucket = aws.s3.Bucket(
    "bucket",
    versioning={"enabled": True},
    server_side_encryption_configuration={
        "rule": {
            "apply_server_side_encryption_by_default": {
                "sse_algorithm": "AES256",
            },
        },
    },
)
```

## Virtual Environment Setup

### Using pip

```bash
python -m venv venv
source venv/bin/activate  # Linux/macOS
venv\Scripts\activate     # Windows
pip install -r requirements.txt
```

### Using poetry

```bash
poetry install
poetry shell
```

### Using uv

```bash
uv sync
source .venv/bin/activate
```

## Type Checking Configuration

### mypy.ini

```ini
[mypy]
python_version = 3.9
warn_return_any = True
warn_unused_configs = True
plugins = pulumi.mypy

[mypy-pulumi_aws.*]
ignore_missing_imports = True
```

### pyproject.toml (for pyright)

```toml
[tool.pyright]
pythonVersion = "3.9"
typeCheckingMode = "basic"
```
