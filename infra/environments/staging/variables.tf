variable "environment" {
  description = "Environment name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

# VPC Variables
variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
}

variable "private_subnet_cidrs" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
}

# DynamoDB
variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode"
  type        = string
  default     = "PAY_PER_REQUEST"
}

# API Gateway
variable "api_throttle_burst_limit" {
  description = "API throttle burst limit"
  type        = number
}

variable "api_throttle_rate_limit" {
  description = "API throttle rate limit"
  type        = number
}

variable "websocket_throttle_burst_limit" {
  description = "WebSocket throttle burst limit"
  type        = number
}

variable "websocket_throttle_rate_limit" {
  description = "WebSocket throttle rate limit"
  type        = number
}

# Lambda
variable "lambda_memory" {
  description = "Lambda memory allocation"
  type        = number
}

variable "booking_service_zip" {
  description = "Path to booking service ZIP"
  type        = string
}

variable "provider_service_zip" {
  description = "Path to provider service ZIP"
  type        = string
}

variable "payment_service_zip" {
  description = "Path to payment service ZIP"
  type        = string
}

variable "sos_service_zip" {
  description = "Path to SOS service ZIP"
  type        = string
}

variable "authorizer_zip" {
  description = "Path to authorizer ZIP"
  type        = string
}

variable "shared_layer_zip" {
  description = "Path to shared layer ZIP"
  type        = string
}

# Monitoring
variable "alert_email" {
  description = "Email for alerts"
  type        = string
}

variable "log_retention_days" {
  description = "Log retention in days"
  type        = number
}

# Tags
variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
}

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # ElastiCache (Redis) Variables
# variable "redis_node_type" {
#   description = "Redis node type"
#   type        = string
#   default     = "cache.t3.micro"
# }
#
# variable "redis_num_nodes" {
#   description = "Number of Redis nodes"
#   type        = number
#   default     = 1
# }
#
# variable "redis_auth_token" {
#   description = "Redis AUTH token"
#   type        = string
#   sensitive   = true
#   default     = ""
# }

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # RDS (Aurora PostgreSQL) Variables
# variable "db_master_username" {
#   description = "Master database username"
#   type        = string
#   sensitive   = true
# }
#
# variable "db_master_password" {
#   description = "Master database password"
#   type        = string
#   sensitive   = true
# }
#
# variable "db_instance_class" {
#   description = "Database instance class"
#   type        = string
#   default     = "db.t4g.micro"
# }
#
# variable "db_instance_count" {
#   description = "Number of database instances"
#   type        = number
#   default     = 1
# }
