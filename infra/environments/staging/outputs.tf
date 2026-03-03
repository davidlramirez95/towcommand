# VPC Outputs
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = module.vpc.private_subnet_ids
}

# DynamoDB Outputs
output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = module.dynamodb.table_name
}

output "dynamodb_table_arn" {
  description = "DynamoDB table ARN"
  value       = module.dynamodb.table_arn
}

# Cognito Outputs
output "cognito_user_pool_id" {
  description = "Cognito user pool ID"
  value       = module.cognito.user_pool_id
}

output "cognito_mobile_client_id" {
  description = "Cognito mobile client ID"
  value       = module.cognito.mobile_client_id
}

output "cognito_identity_pool_id" {
  description = "Cognito identity pool ID"
  value       = module.cognito.identity_pool_id
}

# API Gateway Outputs
output "rest_api_endpoint" {
  description = "REST API endpoint"
  value       = module.api_gateway.rest_api_invoke_url
}

output "websocket_api_endpoint" {
  description = "WebSocket API endpoint"
  value       = module.api_gateway.websocket_api_endpoint
}

# EventBridge Outputs
output "event_bus_name" {
  description = "EventBridge event bus name"
  value       = module.eventbridge.event_bus_name
}

# S3 Outputs
output "evidence_bucket_name" {
  description = "Evidence S3 bucket name"
  value       = module.s3.evidence_bucket_name
}

# Lambda Outputs
output "booking_function_arn" {
  description = "Booking Lambda function ARN"
  value       = module.lambda.booking_function_arn
}

output "provider_function_arn" {
  description = "Provider Lambda function ARN"
  value       = module.lambda.provider_function_arn
}

output "payment_function_arn" {
  description = "Payment Lambda function ARN"
  value       = module.lambda.payment_function_arn
}

output "sos_function_arn" {
  description = "SOS Lambda function ARN"
  value       = module.lambda.sos_function_arn
}

# Monitoring Outputs
output "alerts_topic_arn" {
  description = "SNS topic for alerts"
  value       = module.monitoring.alerts_topic_arn
}

output "sos_dashboard_name" {
  description = "SOS CloudWatch dashboard name"
  value       = module.monitoring.sos_dashboard_name
}

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # ElastiCache Outputs
# output "redis_endpoint" {
#   description = "Redis primary endpoint"
#   value       = module.elasticache.redis_endpoint
# }

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # RDS Outputs
# output "rds_cluster_endpoint" {
#   description = "RDS cluster endpoint"
#   value       = module.rds.cluster_endpoint
# }
#
# output "rds_reader_endpoint" {
#   description = "RDS cluster reader endpoint"
#   value       = module.rds.cluster_reader_endpoint
# }
