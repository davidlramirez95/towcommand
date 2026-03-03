# TowCommand Infrastructure as Code

This directory contains the Terraform infrastructure for the TowCommand platform. The infrastructure is organized into reusable modules and environment-specific configurations.

## Directory Structure

```
infra/
├── modules/                 # Reusable Terraform modules
│   ├── dynamodb/           # DynamoDB table with GSIs
│   ├── cognito/            # Cognito User Pool and Identity Pool
│   ├── api-gateway/        # REST and WebSocket APIs
│   ├── lambda/             # Lambda functions and layers
│   ├── eventbridge/        # EventBridge event bus and rules
│   ├── elasticache/        # Redis cluster
│   ├── rds/                # Aurora PostgreSQL cluster
│   ├── s3/                 # S3 buckets with encryption and lifecycle
│   ├── vpc/                # VPC with public/private subnets
│   └── monitoring/         # CloudWatch dashboards, alarms, and X-Ray
├── environments/           # Environment-specific configurations
│   ├── dev/                # Development environment
│   ├── staging/            # Staging environment
│   └── prod/               # Production environment
└── README.md              # This file
```

## Modules

### DynamoDB
- Single-region NoSQL database
- 5 Global Secondary Indexes for flexible querying
- DynamoDB Streams enabled for event processing
- Point-in-time recovery and server-side encryption enabled
- Pay-per-request billing for dev, provisioned for prod

### Cognito
- User Pool with phone and email as login options
- Custom attributes: user_type, trust_tier, provider_id
- Mobile app client with SRP authentication
- Identity Pool for AWS resource access
- Device tracking with challenge-required-on-new-device

### API Gateway
- REST API with regional endpoint
- WebSocket API for real-time communication
- Cognito authorizers for both REST and WebSocket
- Throttling policies to prevent abuse
- X-Ray tracing enabled
- CloudWatch logging with detailed metrics

### Lambda
- 5 microservices: Booking, Provider, Payment, SOS, Authorizer
- ARM64 architecture for cost optimization
- Shared Lambda layer for common dependencies
- X-Ray tracing and CloudWatch logging
- Fine-grained IAM roles per function

### EventBridge
- Custom event bus for application events
- Rules for: BookingCreated, BookingCompleted, SOSActivated, PaymentCompleted, ProviderOnline
- Event schema registry for validation
- Dead-letter queue support via EventBridge

### ElastiCache
- Redis cluster for caching and real-time data
- Multi-AZ support in prod
- Automatic failover in prod
- At-rest and in-transit encryption
- CloudWatch logs for slow-log and engine logs

### RDS
- Aurora PostgreSQL cluster for relational data
- Read replicas for scalability
- Automated backups and point-in-time recovery
- Enhanced monitoring with CloudWatch
- Performance Insights enabled
- KMS encryption at rest

### S3
- Evidence bucket with Object Lock (GOVERNANCE mode, 30-day minimum)
- Automatic lifecycle policies: STANDARD_IA (30d), GLACIER (90d), DEEP_ARCHIVE (180d)
- Versioning and MFA delete protection
- Separate logging bucket with lifecycle rules
- KMS encryption at rest
- Public access blocked

### VPC
- Public and private subnets across multiple AZs
- NAT Gateways for private subnet internet access
- VPC Endpoints for S3, DynamoDB, and CloudWatch Logs
- VPC Flow Logs for security monitoring
- Internet Gateway for public internet access

### Monitoring
- CloudWatch Dashboards for SOS metrics
- Alarms for Lambda errors, DynamoDB throttling, API 5xx errors, Payment failures
- SNS topic for alert notifications
- X-Ray sampling rules and insight groups
- CloudWatch log groups for all services
- Custom metrics via log metric filters

## Environments

### Dev
- Pay-per-request DynamoDB
- Single-node Redis (t3.micro)
- Single-instance RDS (t4g.micro)
- Lower throttling limits
- 14-day log retention
- Skip final snapshot on RDS destroy

### Staging
- Provisioned DynamoDB
- 2-node Redis with failover (t3.small)
- 2-instance Aurora cluster (t4g.small)
- Medium throttling limits
- 30-day log retention
- Final snapshot taken on destroy

### Production
- Provisioned DynamoDB
- 3-node Redis with failover (r6g.xlarge)
- 3-instance Aurora cluster (r6g.xlarge)
- High throttling limits
- 365-day log retention
- Final snapshot taken on destroy
- Additional backups and availability zones

## Getting Started

### Prerequisites
- Terraform >= 1.0
- AWS CLI configured with appropriate credentials
- S3 bucket for Terraform state (create before deployment)
- DynamoDB table for state locking (create before deployment)

### Deployment Steps

1. **Create state backend resources** (one-time setup):
```bash
# Create S3 bucket for state
aws s3api create-bucket --bucket towcommand-terraform-state-dev --region us-east-1

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket towcommand-terraform-state-dev \
  --versioning-configuration Status=Enabled

# Enable encryption
aws s3api put-bucket-encryption \
  --bucket towcommand-terraform-state-dev \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'

# Create DynamoDB table for locking
aws dynamodb create-table \
  --table-name terraform-lock-dev \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

2. **Configure environment**:
```bash
cd infra/environments/dev
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

3. **Initialize and plan**:
```bash
terraform init
terraform plan -out=tfplan
```

4. **Apply infrastructure**:
```bash
terraform apply tfplan
```

5. **Retrieve outputs**:
```bash
terraform output
```

## Sensitive Variables

The following variables contain sensitive information and should NOT be committed to version control:

- `db_master_password` - RDS master password
- `redis_auth_token` - Redis AUTH token
- Database credentials

Use a `.tfvars.local` file or AWS Secrets Manager/Parameter Store for these values.

## State Management

- Terraform state is stored in S3 with encryption and versioning enabled
- State locking is managed via DynamoDB to prevent concurrent modifications
- State files contain sensitive data and should never be committed to git

## Updating Infrastructure

### Planning changes:
```bash
cd infra/environments/{dev,staging,prod}
terraform plan
```

### Applying changes:
```bash
terraform apply
```

### Destroying resources:
```bash
# Warning: This will delete all resources in the environment
terraform destroy
```

## Troubleshooting

### State lock issues:
```bash
# View locks
aws dynamodb scan --table-name terraform-lock-dev

# Force unlock (use with caution)
terraform force-unlock <LOCK_ID>
```

### Provider version conflicts:
```bash
terraform init -upgrade
```

### Permission denied errors:
- Verify AWS credentials are configured
- Check IAM permissions for the executing user/role

## Security Best Practices

1. **Never commit sensitive data** - Use `.tfvars.local` or Secrets Manager
2. **State file protection** - Always enable S3 encryption and versioning
3. **IAM roles** - Use least-privilege principle for all Lambda and service roles
4. **Encryption** - All data at rest is encrypted with KMS
5. **Network isolation** - RDS and ElastiCache in private subnets
6. **Audit logging** - VPC Flow Logs, CloudTrail, and CloudWatch Logs enabled

## Cost Optimization

- Use ARM64 Lambda for 20% cost savings
- Pay-per-request billing for dev DynamoDB
- T-series instances for non-production databases
- Lifecycle policies for S3 archival
- Reserved instances for production (if committed)

## Disaster Recovery

- RDS automated backups with 7-35 day retention
- DynamoDB point-in-time recovery enabled
- S3 versioning and Object Lock for evidence
- Multi-AZ deployment in production
- Cross-region failover (manual setup required)

## Monitoring and Alerts

- CloudWatch alarms for critical metrics
- SNS notifications to team email
- X-Ray for distributed tracing
- Custom metrics for application events
- Log retention policies for cost management

## Support and Maintenance

For issues or questions about the infrastructure:
1. Check CloudWatch Logs and CloudTrail
2. Review X-Ray traces for errors
3. Validate Terraform configuration with `terraform validate`
4. Check AWS service limits and quotas

## References

- [Terraform AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [DynamoDB Best Practices](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html)
