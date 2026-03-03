# TowCommand Serverless Infrastructure Migration

## Executive Summary

Successfully removed all non-serverless AWS resources from TowCommand infrastructure to eliminate idle costs. The infrastructure now incurs near-zero charges when not in use.

**Estimated Monthly Cost Savings: ~$118+**

## Resources Removed

### 1. Database & Cache Layer
- **RDS PostgreSQL** (~$36/month) - Removed from all environments
- **ElastiCache Redis** (~$18/month) - Removed from all environments

### 2. Network Infrastructure
- **NAT Gateways** (~$32/month each × 2) - Removed from VPC
- **Elastic IPs** - Associated with NAT gateways

## Files Modified

### Infrastructure (Terraform)

#### Environment Configuration
1. `infra/environments/dev/main.tf` - Removed RDS/ElastiCache modules
2. `infra/environments/staging/main.tf` - Removed RDS/ElastiCache modules
3. `infra/environments/prod/main.tf` - Removed RDS/ElastiCache modules

#### Environment Variables
4. `infra/environments/dev/variables.tf` - Removed redis_* and db_* variables
5. `infra/environments/staging/variables.tf` - Removed redis_* and db_* variables
6. `infra/environments/prod/variables.tf` - Removed redis_* and db_* variables

#### Environment Outputs
7. `infra/environments/dev/outputs.tf` - Removed RDS/Redis outputs
8. `infra/environments/staging/outputs.tf` - Removed RDS/Redis outputs
9. `infra/environments/prod/outputs.tf` - Removed RDS/Redis outputs

#### Modules
10. `infra/modules/vpc/main.tf` - Removed NAT gateways, kept VPC endpoints
11. `infra/modules/lambda/main.tf` - Removed Redis endpoint env vars
12. `infra/modules/lambda/variables.tf` - Removed redis_endpoint variable
13. `infra/modules/lambda/iam.tf` - Removed ElastiCache IAM policies

### Application Code

#### Analytics Service
14. `services/analytics/src/lib/pg-client.ts` - Removed PostgreSQL Pool, added TODO for DynamoDB/Athena
15. `services/analytics/src/handler.ts` - Updated to note PostgreSQL removal
16. `services/analytics/src/queries/demand-heatmap.ts` - Removed SQL, added DynamoDB TODO
17. `services/analytics/src/queries/provider-performance.ts` - Removed SQL, added DynamoDB TODO
18. `services/analytics/src/queries/revenue-report.ts` - Removed SQL, added DynamoDB TODO

### Configuration Files

19. `.env.example` - Removed REDIS_* and PG_* variables
20. `docker-compose.yml` - Removed postgres and redis services

## Architecture Changes

### Before (Hybrid)
```
Lambda → RDS (PostgreSQL)      [~$36/month]
      → ElastiCache (Redis)     [~$18/month]
      → NAT Gateway → Internet  [~$64/month]
```

### After (Serverless)
```
Lambda → DynamoDB (via VPC Endpoint)  [Pay-per-request]
      → S3 (via VPC Endpoint)         [Pay-per-storage]
      → CloudWatch Logs (VPC Endpoint) [Pay-per-ingestion]
```

## VPC Configuration

### Kept
- VPC and subnets (minimal cost ~$0.05/day)
- Internet Gateway (no charge)
- VPC Endpoints for:
  - S3 Gateway Endpoint
  - DynamoDB Gateway Endpoint
  - CloudWatch Logs Interface Endpoint
- VPC Flow Logs

### Removed
- NAT Gateways (primary cost driver)
- Elastic IPs associated with NAT
- NAT Gateway routes in private subnets

### Why VPC is Still Needed
1. API Gateway private integrations (future feature)
2. VPC endpoint security for AWS service access
3. Network isolation and security groups

## Analytics Service Migration

### Current State
- PostgreSQL queries removed
- TODO comments with implementation guidance

### Next Steps
Choose one serverless approach:

#### Option 1: DynamoDB (Recommended for MVP)
- Store metrics in DynamoDB tables
- Use Global Secondary Indexes (GSI) for time-based queries
- Update via EventBridge streams
- Query in Lambda for real-time dashboards
- Minimal latency, familiar SQL-like interface (PartiQL)

**Example Schema:**
```
AnalyticsMetrics (PK: metric_id, SK: timestamp)
- Demand heatmap: metric_id = "demand#{grid_cell_id}", timestamp = ISO string
- Provider stats: metric_id = "provider#{provider_id}", timestamp = ISO string
- Revenue: metric_id = "revenue#{day}", timestamp = ISO string
```

#### Option 2: Amazon Athena + S3
- Archive events to S3 (via S3 sink in EventBridge)
- Query with Athena (SQL interface)
- Pay only per query executed
- Good for ad-hoc analytics and reports

#### Option 3: Aurora Serverless v2
- If ACID transactions become critical
- Auto-scales to zero when idle
- Similar to Aurora but with better scaling
- More cost-effective than RDS for variable workloads

#### Option 4: Amazon QuickSight
- Visualization layer on top of any data source
- No database needed for dashboards
- Self-service BI

## Migration Checklist

- [ ] Remove PostgreSQL logic from analytics service
- [ ] Implement DynamoDB schema for analytics
- [ ] Update EventBridge rules to write to DynamoDB
- [ ] Rewrite analytics queries using DynamoDB PartiQL or AWS SDK
- [ ] Test analytics dashboards with DynamoDB
- [ ] Remove local postgres from developer setup docs
- [ ] Update CI/CD to not require postgres containers
- [ ] Performance test DynamoDB queries vs original SQL

## Testing the New Architecture

### Lambda Connectivity
```bash
# Test DynamoDB access from Lambda
aws lambda invoke --function-name towcommand-booking-dev \
  --payload '{"test": true}' response.json

# Check CloudWatch Logs for connectivity
aws logs tail /aws/lambda/towcommand-booking-dev --follow
```

### VPC Endpoint Verification
```bash
# Verify VPC endpoints are functional
aws ec2 describe-vpc-endpoints --filter "Name=vpc-id,Values=<vpc-id>"
```

## Rollback Plan

If you need to re-add RDS/Redis:

1. Restore from Terraform state history
2. Recreate RDS/ElastiCache modules
3. Restore security groups and IAM policies
4. Repoint Lambda environment variables
5. Restore analytics to use PostgreSQL

**Note:** Terraform state history will contain the full configuration.

## Monitoring & Alerts

### CloudWatch Metrics to Monitor
- Lambda execution duration (may increase without Redis cache)
- DynamoDB consumed capacity (if migrating analytics)
- VPC endpoint packet flows

### Cost Monitoring
Track AWS Cost Explorer for:
- Reduction in NAT Gateway charges
- Reduction in RDS charges
- Reduction in ElastiCache charges

Expected to see charges drop by ~$118/month.

## Documentation Updates Needed

- [ ] ARCHITECTURE.md - Update to reflect serverless-only design
- [ ] DEVELOPMENT.md - Remove docker postgres setup
- [ ] README - Note about serverless design and cost benefits
- [ ] Setup guides - Remove postgres/redis installation steps

## Questions & Support

### What if Lambda needs caching?
- Use DynamoDB with TTL
- Use ElastiCache for specific features (optional, only when needed)
- Cache at application level in Lambda

### What about analytics performance?
- DynamoDB: Sub-millisecond for single item lookup, microsecond latency for GSI queries
- Athena: Query latency seconds (good for scheduled reports, not real-time)
- Consider hybrid: DynamoDB for real-time, Athena for historical analysis

### Is the VPC still needed?
- Yes, for future API Gateway private integrations
- For network security and isolation
- Cost is minimal (~$0.05/day)
- Can be simplified further if APIs are public-only

## References

- AWS Lambda Best Practices: https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html
- VPC Endpoints: https://docs.aws.amazon.com/vpc/latest/privatelink/endpoint-services.html
- DynamoDB Analytics: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html
- Aurora Serverless v2: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraMigrationGuide/Aurora.Serverless.html
