variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "api_throttle_burst_limit" {
  description = "API Gateway throttle burst limit"
  type        = number
  default     = 5000
}

variable "api_throttle_rate_limit" {
  description = "API Gateway throttle rate limit"
  type        = number
  default     = 2000
}

variable "websocket_throttle_burst_limit" {
  description = "WebSocket throttle burst limit"
  type        = number
  default     = 5000
}

variable "websocket_throttle_rate_limit" {
  description = "WebSocket throttle rate limit"
  type        = number
  default     = 2000
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 14
}

variable "cognito_user_pool_arn" {
  description = "ARN of the Cognito User Pool"
  type        = string
}

variable "cognito_client_id" {
  description = "Cognito client ID"
  type        = string
}

variable "cognito_user_pool_endpoint" {
  description = "Cognito User Pool endpoint"
  type        = string
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
