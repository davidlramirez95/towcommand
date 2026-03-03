variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "lambda_memory" {
  description = "Memory allocated to Lambda functions"
  type        = number
  default     = 256
}

variable "dynamodb_table_name" {
  description = "DynamoDB table name"
  type        = string
}

variable "dynamodb_table_arn" {
  description = "DynamoDB table ARN"
  type        = string
}

variable "event_bus_name" {
  description = "EventBridge bus name"
  type        = string
}

variable "event_bus_arn" {
  description = "EventBridge bus ARN"
  type        = string
}

variable "booking_service_zip" {
  description = "Path to booking service ZIP file"
  type        = string
}

variable "provider_service_zip" {
  description = "Path to provider service ZIP file"
  type        = string
}

variable "payment_service_zip" {
  description = "Path to payment service ZIP file"
  type        = string
}

variable "sos_service_zip" {
  description = "Path to SOS service ZIP file"
  type        = string
}

variable "authorizer_zip" {
  description = "Path to authorizer ZIP file"
  type        = string
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}

# TODO: Uncomment when ElastiCache/RDS is provisioned and budget allows
# variable "redis_endpoint" {
#   description = "Redis primary endpoint"
#   type        = string
# }
#
# variable "rds_endpoint" {
#   description = "RDS cluster endpoint"
#   type        = string
# }
