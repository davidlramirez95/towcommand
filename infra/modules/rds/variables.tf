variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "private_subnet_cidrs" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "towcommand"
}

variable "db_master_username" {
  description = "Master database username"
  type        = string
  sensitive   = true
}

variable "db_master_password" {
  description = "Master database password"
  type        = string
  sensitive   = true
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "15.2"
}

variable "db_instance_class" {
  description = "Database instance class"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_instance_count" {
  description = "Number of database instances"
  type        = number
  default     = 2
}

variable "db_max_connections" {
  description = "Maximum database connections"
  type        = string
  default     = "100"
}

variable "db_backup_retention" {
  description = "Database backup retention period (days)"
  type        = number
  default     = 7
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 14
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
