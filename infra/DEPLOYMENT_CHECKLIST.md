# TowCommand Terraform Deployment Checklist

Complete checklist for deploying TowCommand infrastructure across dev, staging, and production environments.

## Pre-Deployment Setup

### AWS Account Prerequisites
- [ ] AWS account created and credentials configured
- [ ] Appropriate IAM permissions for account owner
- [ ] CloudTrail enabled for audit logging
- [ ] AWS Config enabled for compliance monitoring
- [ ] Budget alerts configured

### Backend Infrastructure (One-time setup)

#### Create S3 Buckets for State Storage
```bash
# Development
aws s3api create-bucket --bucket towcommand-terraform-state-dev --region us-east-1
aws s3api put-bucket-versioning --bucket towcommand-terraform-state-dev --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket towcommand-terraform-state-dev --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
aws s3api put-public-access-block --bucket towcommand-terraform-state-dev --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

# Staging
aws s3api create-bucket --bucket towcommand-terraform-state-staging --region us-east-1
aws s3api put-bucket-versioning --bucket towcommand-terraform-state-staging --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket towcommand-terraform-state-staging --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
aws s3api put-public-access-block --bucket towcommand-terraform-state-staging --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

# Production
aws s3api create-bucket --bucket towcommand-terraform-state-prod --region us-east-1
aws s3api put-bucket-versioning --bucket towcommand-terraform-state-prod --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket towcommand-terraform-state-prod --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
aws s3api put-public-access-block --bucket towcommand-terraform-state-prod --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"
```

#### Create DynamoDB Tables for State Locking
```bash
# Development
aws dynamodb create-table \
  --table-name terraform-lock-dev \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1

# Staging
aws dynamodb create-table \
  --table-name terraform-lock-staging \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1

# Production
aws dynamodb create-table \
  --table-name terraform-lock-prod \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

- [ ] S3 buckets created with encryption and versioning
- [ ] DynamoDB tables created for state locking
- [ ] S3 blocks public access
- [ ] Verify buckets are in us-east-1 region

## Environment Configuration

### Development Environment Setup
- [ ] Copy `infra/environments/dev/terraform.tfvars.example` to `terraform.tfvars`
- [ ] Set environment = "dev"
- [ ] Configure VPC CIDR: 10.0.0.0/16
- [ ] Set db_master_password (use random secure password)
- [ ] Set redis_auth_token (or leave empty for no auth)
- [ ] Set alert_email to team inbox
- [ ] Configure Lambda ZIP paths (or use placeholders)
- [ ] Review and adjust other variables as needed

### Staging Environment Setup
- [ ] Copy `infra/environments/staging/terraform.tfvars.example` to `terraform.tfvars`
- [ ] Set environment = "staging"
- [ ] Configure VPC CIDR: 10.1.0.0/16
- [ ] Set db_master_password
- [ ] Set redis_auth_token
- [ ] Set alert_email
- [ ] Configure Lambda ZIP paths
- [ ] Review capacity values (2 AZs, medium instances)

### Production Environment Setup
- [ ] Copy `infra/environments/prod/terraform.tfvars.example` to `terraform.tfvars`
- [ ] Set environment = "prod"
- [ ] Configure VPC CIDR: 10.2.0.0/16
- [ ] Set db_master_password (use KMS-encrypted value)
- [ ] Set redis_auth_token (use KMS-encrypted value)
- [ ] Set alert_email to operations team
- [ ] Configure Lambda ZIP paths
- [ ] Review capacity values (3 AZs, large instances)

### Lambda Deployment Packages
- [ ] Create `packages/booking-service.zip`
- [ ] Create `packages/provider-service.zip`
- [ ] Create `packages/payment-service.zip`
- [ ] Create `packages/sos-service.zip`
- [ ] Create `packages/authorizer.zip`
- [ ] Create `packages/shared-layer.zip`
- [ ] Verify all packages are Node.js compatible
- [ ] Ensure packages include node_modules for dependencies
- [ ] Test packages locally before deploying

## Pre-Deployment Validation

### Terraform Syntax Check
```bash
cd infra/environments/{dev,staging,prod}
terraform init
terraform validate
```
- [ ] No syntax errors reported
- [ ] All module references valid
- [ ] All variables properly typed

### Format Check
```bash
terraform fmt -recursive infra/
```
- [ ] All files properly formatted

### Documentation Review
- [ ] Reviewed README.md for overview
- [ ] Reviewed STRUCTURE.md for detailed architecture
- [ ] Understood module dependencies
- [ ] Understood variable hierarchy
- [ ] Reviewed security settings for each environment

## Deployment Execution

### Development Environment Deployment

```bash
cd infra/environments/dev

# Initialize Terraform with backend
terraform init

# Create deployment plan
terraform plan -out=tfplan

# Review plan output
# Check resources to be created
# Verify values from terraform.tfvars
```
- [ ] Terraform init completed successfully
- [ ] terraform plan shows expected resources
- [ ] VPC created with 2 AZs
- [ ] DynamoDB table created (pay-per-request)
- [ ] Cognito resources created
- [ ] API Gateway REST and WebSocket created
- [ ] Lambda functions created with ARM64
- [ ] ElastiCache Redis created
- [ ] RDS Aurora cluster created
- [ ] S3 buckets created with proper encryption
- [ ] Monitoring dashboard and alarms created

```bash
# Apply infrastructure
terraform apply tfplan

# Save outputs
terraform output > outputs.json
```
- [ ] terraform apply completed successfully
- [ ] All resources created
- [ ] Outputs saved to outputs.json
- [ ] Verify in AWS Console:
  - [ ] VPC created
  - [ ] Subnets visible in Console
  - [ ] DynamoDB table accessible
  - [ ] Cognito user pool created
  - [ ] API Gateway endpoints working
  - [ ] Lambda functions appear in Console
  - [ ] CloudWatch log groups created

### Staging Environment Deployment

```bash
cd infra/environments/staging

# Initialize Terraform
terraform init

# Create deployment plan
terraform plan -out=tfplan

# Review plan
```
- [ ] Check resource counts
- [ ] Verify capacity increases from dev
- [ ] Confirm 2 instances/nodes for HA
- [ ] Review VPC CIDR (10.1.0.0/16)

```bash
# Apply infrastructure
terraform apply tfplan

# Save outputs
terraform output > outputs.json
```
- [ ] All resources created
- [ ] 2 RDS instances in different AZs
- [ ] 2 ElastiCache nodes
- [ ] NAT gateways per AZ
- [ ] Increased capacity from dev configuration

### Production Environment Deployment

```bash
cd infra/environments/prod

# Initialize Terraform
terraform init

# Create deployment plan
terraform plan -out=tfplan

# Review plan carefully
```
- [ ] Verify resource counts match expected
- [ ] Check high-capacity instance types
- [ ] Confirm 3 AZ deployment
- [ ] Review 365-day log retention
- [ ] Verify backup configurations

```bash
# Apply infrastructure
terraform apply tfplan

# Save outputs
terraform output > outputs.json

# Create state backup
aws s3 cp terraform.tfstate s3://towcommand-terraform-state-prod/backup/
```
- [ ] All production resources created
- [ ] 3 RDS instances across AZs
- [ ] 3 ElastiCache nodes with Multi-AZ
- [ ] 3 availability zones configured
- [ ] Enhanced monitoring enabled
- [ ] Alarms configured with SNS
- [ ] State backup completed

## Post-Deployment Testing

### Network Connectivity Tests
- [ ] VPC created with correct CIDR blocks
- [ ] Public subnets have internet access
- [ ] Private subnets route through NAT
- [ ] VPC endpoints working (S3, DynamoDB, Logs)
- [ ] Security groups configured correctly
- [ ] Route tables properly configured

### Database Tests
```bash
# RDS Connection Test
psql -h {rds_endpoint} -U admin -d towcommand

# Redis Connection Test
redis-cli -h {redis_endpoint} ping
```
- [ ] RDS cluster accessible
- [ ] Can create test table in RDS
- [ ] Redis cluster responds to PING
- [ ] Can set/get values in Redis
- [ ] Encryption enabled on both

### Cognito Tests
- [ ] User pool created and accessible
- [ ] Mobile app client configured
- [ ] Identity pool connected
- [ ] Can initiate sign-up flow
- [ ] Phone and email login options work

### API Gateway Tests
```bash
# REST API Health Check
curl https://{api_endpoint}/health

# WebSocket Connection
wcat ws://{websocket_endpoint}
```
- [ ] REST API health check returns 200
- [ ] WebSocket endpoint responds
- [ ] Cognito authorization working
- [ ] Throttling limits in place

### Lambda Tests
- [ ] Functions appear in Console
- [ ] Can invoke test payloads
- [ ] CloudWatch logs appear
- [ ] Environment variables set correctly
- [ ] X-Ray traces visible

### Monitoring Tests
- [ ] CloudWatch dashboards created
- [ ] Can view dashboard metrics
- [ ] Alarms created and active
- [ ] SNS topic configured
- [ ] Test email sent from SNS topic
- [ ] X-Ray groups and rules created

### S3 Tests
- [ ] Evidence bucket accessible
- [ ] Object Lock enabled and enforced
- [ ] Lifecycle policies configured
- [ ] Versioning enabled
- [ ] Encryption working
- [ ] Logging bucket receives logs

## Backup and Documentation

### Generate Documentation
```bash
# Create architecture diagram
terraform graph | dot -Tsvg > architecture.svg

# Export state
terraform state list > resource_list.txt

# Create outputs file
terraform output -json > outputs.json
```
- [ ] Architecture diagrams generated
- [ ] Resource list documented
- [ ] Outputs saved for reference

### Create Runbooks
- [ ] Scaling runbook for each resource
- [ ] Disaster recovery procedures
- [ ] Monitoring and alerting guide
- [ ] Security incident response guide
- [ ] Rollback procedures

### Team Communication
- [ ] Notify team of new endpoints
- [ ] Share API documentation
- [ ] Share Cognito pool configuration
- [ ] Share database access details
- [ ] Document cost breakdown

## Post-Deployment Monitoring Setup

### CloudWatch Configuration
- [ ] Verify SOS dashboard exists
- [ ] Confirm Lambda error alarms active
- [ ] Check DynamoDB throttle alarms
- [ ] Verify API 5xx error alarms
- [ ] Check payment failure tracking

### Log Aggregation
- [ ] Confirm log groups created for:
  - [ ] API Gateway access logs
  - [ ] Lambda function logs
  - [ ] RDS logs
  - [ ] ElastiCache logs
  - [ ] VPC Flow Logs

### X-Ray Setup
- [ ] X-Ray sampling enabled
- [ ] Error group created
- [ ] Throttling group created
- [ ] Insight rule active
- [ ] Can see traces in console

## Cleanup and Deprovisioning (When needed)

### Before Destroying
```bash
# Backup current state
aws s3 cp terraform.tfstate s3://towcommand-terraform-state-{env}/backups/

# Create snapshots
aws rds create-db-cluster-snapshot --db-cluster-identifier towcommand-{env} --db-cluster-snapshot-identifier final-snapshot

# Export DynamoDB data
aws dynamodb export-table-to-point-in-time --table-arn {arn} --s3-bucket {bucket}
```
- [ ] State files backed up
- [ ] RDS snapshots created
- [ ] DynamoDB data exported
- [ ] S3 data backed up
- [ ] Confirmation received from team

### Destroy Infrastructure
```bash
cd infra/environments/{env}
terraform destroy -auto-approve
```
- [ ] Destruction approved by team lead
- [ ] All resources destroyed
- [ ] S3 state bucket cleaned up
- [ ] DynamoDB lock table removed
- [ ] VPC endpoints deleted
- [ ] Verify in AWS Console

## Ongoing Maintenance

### Weekly Tasks
- [ ] Review CloudWatch alarms
- [ ] Check log aggregation
- [ ] Monitor costs
- [ ] Review X-Ray traces for errors

### Monthly Tasks
- [ ] Review security group rules
- [ ] Check backup status
- [ ] Rotate credentials if needed
- [ ] Update Terraform providers
- [ ] Review capacity utilization

### Quarterly Tasks
- [ ] Review and update documentation
- [ ] Disaster recovery drill
- [ ] Capacity planning review
- [ ] Security audit
- [ ] Cost optimization review

### Annually
- [ ] Major infrastructure review
- [ ] Update to latest Terraform version
- [ ] AWS service updates assessment
- [ ] Architecture redesign if needed

## Support and Escalation

### Common Issues

#### Terraform State Lock
```bash
# View locks
aws dynamodb scan --table-name terraform-lock-{env}

# Force unlock
terraform force-unlock {LOCK_ID}
```

#### Lambda Permission Denied
- Check IAM role permissions
- Verify DynamoDB/EventBridge ARNs
- Review encryption key access

#### RDS Connection Issues
- Check security group rules
- Verify subnet routing
- Confirm database credentials
- Check performance insights

#### API Gateway Errors
- Check CloudWatch logs
- Verify Cognito authorizer
- Review throttling limits
- Check backend Lambda logs

### Emergency Contacts
- Infrastructure Team Lead: [Name/Contact]
- AWS Support Plan: [Plan Level]
- OncAll Engineer: [Rotation Schedule]

---

**Last Updated**: [Date]
**Terraform Version**: >= 1.0
**AWS Provider Version**: >= 5.0
**Next Review Date**: [Date]
