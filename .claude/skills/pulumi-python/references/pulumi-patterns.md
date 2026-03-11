# Pulumi Infrastructure Patterns (Python)

## Component Resources

```python
import pulumi
from pulumi_aws import ec2


class VpcArgs:
    def __init__(
        self,
        cidr_block: pulumi.Input[str],
        az_count: int = 2,
        enable_nat_gateway: bool = False,
    ):
        self.cidr_block = cidr_block
        self.az_count = az_count
        self.enable_nat_gateway = enable_nat_gateway


class Vpc(pulumi.ComponentResource):
    vpc_id: pulumi.Output[str]
    public_subnet_ids: pulumi.Output[list]
    private_subnet_ids: pulumi.Output[list]

    def __init__(
        self,
        name: str,
        args: VpcArgs,
        opts: pulumi.ResourceOptions = None,
    ):
        super().__init__("custom:network:Vpc", name, {}, opts)

        self.vpc = ec2.Vpc(
            f"{name}-vpc",
            cidr_block=args.cidr_block,
            enable_dns_hostnames=True,
            enable_dns_support=True,
            tags={"Name": name},
            opts=pulumi.ResourceOptions(parent=self),
        )

        self.vpc_id = self.vpc.id

        # Create subnets...
        self.public_subnet_ids = []
        self.private_subnet_ids = []

        self.register_outputs({
            "vpc_id": self.vpc_id,
            "public_subnet_ids": self.public_subnet_ids,
            "private_subnet_ids": self.private_subnet_ids,
        })
```

## Stack References

```python
import pulumi

# Reference networking stack
network_stack = pulumi.StackReference("myorg/networking/prod")

# Get typed outputs
vpc_id = network_stack.get_output("vpc_id")
subnet_ids = network_stack.get_output("private_subnet_ids")

# Use in resources
security_group = aws.ec2.SecurityGroup(
    "sg",
    vpc_id=vpc_id,
    # ...
)
```

## Transformations

```python
import pulumi

# Register transformation for all resources
def add_tags(args: pulumi.ResourceTransformationArgs):
    if args.props.get("tags") is not None:
        args.props["tags"]["Environment"] = pulumi.get_stack()
        args.props["tags"]["ManagedBy"] = "Pulumi"
    return pulumi.ResourceTransformationResult(args.props, args.opts)

pulumi.runtime.register_stack_transformation(add_tags)

# Enforce encryption on all S3 buckets
def enforce_s3_encryption(args: pulumi.ResourceTransformationArgs):
    if args.type_ == "aws:s3/bucket:Bucket":
        args.props["server_side_encryption_configuration"] = {
            "rule": {
                "apply_server_side_encryption_by_default": {
                    "sse_algorithm": "aws:kms",
                },
            },
        }
    return pulumi.ResourceTransformationResult(args.props, args.opts)

pulumi.runtime.register_stack_transformation(enforce_s3_encryption)
```

## Resource Protection

Prevent accidental deletion of critical resources:

```python
import pulumi
from pulumi_aws import rds, s3

# Protect resource from deletion
db = rds.Instance(
    "prod-db",
    # ... config
    opts=pulumi.ResourceOptions(protect=True),
)

# Retain resource on stack destroy (useful for data)
bucket = s3.Bucket(
    "data-bucket",
    # ... config
    opts=pulumi.ResourceOptions(retain_on_delete=True),
)
```

## Component Resource Pitfalls

Common mistakes to avoid:

1. **Forgetting `parent=self`** - Child resources won't be properly associated; deletion may fail
2. **Skipping `register_outputs()`** - Outputs won't be tracked; cross-stack references break
3. **Hardcoding values** - Reduces reusability; use args class instead
4. **Overly complex components** - Keep components focused on single logical units
5. **Missing input validation** - Validate args before creating resources

## Stack Pitfalls

1. **Hardcoding environment values** - Use configuration or ESC instead
2. **Sharing state between environments** - Each stack should have isolated state
3. **Circular stack references** - Design dependency graph carefully
4. **Unprotected production resources** - Use `protect=True` for critical infra
5. **Unencrypted secrets** - Always use `--secret` flag or ESC secrets

## CI/CD Patterns

### GitHub Actions with ESC

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: pulumi/auth-actions@v1
        with:
          organization: myorg
          requested-token-type: urn:pulumi:token-type:access_token:organization

      - name: Deploy with ESC
        run: |
          pulumi env run myorg/aws-prod -- pulumi up --yes
```

## Configuration Patterns

```python
import pulumi

config = pulumi.Config("myapp")

# Environment-aware defaults
is_prod = pulumi.get_stack() == "prod"
instance_type = config.get("instance_type") or ("t3.large" if is_prod else "t3.small")

# Structured configuration
db_config = config.require_object("database")
# Returns dict: {"host": "...", "port": 5432, "name": "..."}

# With type hints
from typing import TypedDict

class DatabaseConfig(TypedDict):
    host: str
    port: int
    name: str

db: DatabaseConfig = config.require_object("database")
```

## Testing Patterns

### Unit Testing with Mocks

```python
import unittest
import pulumi


class MyMocks(pulumi.runtime.Mocks):
    def new_resource(self, args: pulumi.runtime.MockResourceArgs):
        return [f"{args.name}_id", args.inputs]

    def call(self, args: pulumi.runtime.MockCallArgs):
        return args.inputs


pulumi.runtime.set_mocks(MyMocks())


# Import after setting mocks
from my_infrastructure import bucket


class TestInfrastructure(unittest.TestCase):
    @pulumi.runtime.test
    def test_bucket_versioning(self):
        def check_versioning(versioning):
            self.assertTrue(versioning["enabled"])

        bucket.versioning.apply(check_versioning)
```

### Integration Testing with Automation API

```python
import pulumi.automation as auto


def test_deploy():
    stack = auto.create_or_select_stack(
        stack_name="test",
        project_name="test-project",
        program=lambda: None,  # Inline or from source
    )

    # Deploy
    result = stack.up()
    assert result.summary.result == "succeeded"

    # Get outputs
    outputs = stack.outputs()
    assert "bucket_name" in outputs

    # Cleanup
    stack.destroy()
```

## Resource Options

```python
import pulumi

# Common resource options
resource = aws.s3.Bucket(
    "bucket",
    opts=pulumi.ResourceOptions(
        # Explicit dependencies
        depends_on=[other_resource],

        # Parent for components
        parent=self,

        # Protect from deletion
        protect=True,

        # Ignore changes
        ignore_changes=["tags"],

        # Custom provider
        provider=custom_provider,

        # Aliases for refactoring
        aliases=[pulumi.Alias(name="old-name")],

        # Delete before replace
        delete_before_replace=True,

        # Custom timeouts
        custom_timeouts=pulumi.CustomTimeouts(
            create="30m",
            update="30m",
            delete="30m",
        ),

        # Force replacement when specific properties change
        replace_on_changes=["tags.Environment"],

        # Skip delete when parent resource is deleted (useful for Kubernetes)
        deleted_with=parent_resource,
    ),
)
```

## Resource Hooks

```python
# Define lifecycle hooks for custom logic
def before_create(args: pulumi.HookArgs):
    print(f"Creating resource: {args.urn}")

def after_create(args: pulumi.HookArgs):
    print(f"Created with outputs: {args.outputs}")

resource = aws.s3.Bucket(
    "bucket",
    opts=pulumi.ResourceOptions(
        hooks=pulumi.LifecycleHooks(
            before_create=before_create,
            after_create=after_create,
        ),
    ),
)
```

> **Note:** Hooks require `--run-program` flag when running `pulumi destroy`. Not supported in YAML.

## Dynamic Providers

```python
import pulumi
from pulumi.dynamic import Resource, ResourceProvider, CreateResult


class MyResourceProvider(ResourceProvider):
    def create(self, props):
        # Create logic
        return CreateResult(id_="my-id", outs={"result": "created"})

    def update(self, id, old_props, new_props):
        # Update logic
        return {"result": "updated"}

    def delete(self, id, props):
        # Delete logic
        pass


class MyResource(Resource):
    result: pulumi.Output[str]

    def __init__(self, name, opts=None):
        super().__init__(MyResourceProvider(), name, {"result": None}, opts)
```
