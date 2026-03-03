terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "TowCommand"
      ManagedBy   = "Terraform"
    }
  }
}

# VPC Module - stripped to minimal serverless configuration (no NAT gateways)
module "vpc" {
  source = "../../modules/vpc"

  environment              = var.environment
  aws_region              = var.aws_region
  vpc_cidr                = var.vpc_cidr
  public_subnet_cidrs     = var.public_subnet_cidrs
  private_subnet_cidrs    = var.private_subnet_cidrs
  availability_zones      = var.availability_zones
  log_retention_days      = var.log_retention_days
  tags                    = var.tags
}

# DynamoDB Module
module "dynamodb" {
  source = "../../modules/dynamodb"

  environment   = var.environment
  billing_mode  = var.dynamodb_billing_mode
  tags          = var.tags
}

# Cognito Module
module "cognito" {
  source = "../../modules/cognito"

  environment = var.environment
  tags        = var.tags
}

# API Gateway Module
module "api_gateway" {
  source = "../../modules/api-gateway"

  environment                     = var.environment
  api_throttle_burst_limit        = var.api_throttle_burst_limit
  api_throttle_rate_limit         = var.api_throttle_rate_limit
  websocket_throttle_burst_limit  = var.websocket_throttle_burst_limit
  websocket_throttle_rate_limit   = var.websocket_throttle_rate_limit
  log_retention_days              = var.log_retention_days
  cognito_user_pool_arn           = module.cognito.user_pool_arn
  cognito_client_id               = module.cognito.mobile_client_id
  cognito_user_pool_endpoint      = module.cognito.user_pool_endpoint
  tags                            = var.tags
}

# EventBridge Module
module "eventbridge" {
  source = "../../modules/eventbridge"

  environment = var.environment
  tags        = var.tags
}

# S3 Module
module "s3" {
  source = "../../modules/s3"

  environment = var.environment
  tags        = var.tags
}

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # ElastiCache Module (Redis)
# module "elasticache" {
#   source = "../../modules/elasticache"
#
#   environment          = var.environment
#   vpc_id               = module.vpc.vpc_id
#   private_subnet_ids   = module.vpc.private_subnet_ids
#   private_subnet_cidrs = var.private_subnet_cidrs
#   redis_node_type      = var.redis_node_type
#   redis_num_nodes      = var.redis_num_nodes
#   redis_auth_token     = var.redis_auth_token
#   log_retention_days   = var.log_retention_days
#   tags                 = var.tags
# }

# TODO: Uncomment when feature is ready and budget allows (not serverless - costs money on AWS)
# # RDS Module (Aurora PostgreSQL - for analytics)
# module "rds" {
#   source = "../../modules/rds"
#
#   environment          = var.environment
#   vpc_id               = module.vpc.vpc_id
#   private_subnet_ids   = module.vpc.private_subnet_ids
#   private_subnet_cidrs = var.private_subnet_cidrs
#   db_master_username   = var.db_master_username
#   db_master_password   = var.db_master_password
#   db_instance_class    = var.db_instance_class
#   db_instance_count    = var.db_instance_count
#   log_retention_days   = var.log_retention_days
#   tags                 = var.tags
# }

# Lambda Module - Redis and RDS endpoints removed
module "lambda" {
  source = "../../modules/lambda"

  environment            = var.environment
  lambda_memory          = var.lambda_memory
  dynamodb_table_name    = module.dynamodb.table_name
  dynamodb_table_arn     = module.dynamodb.table_arn
  event_bus_name         = module.eventbridge.event_bus_name
  event_bus_arn          = module.eventbridge.event_bus_arn
  booking_service_zip    = var.booking_service_zip
  provider_service_zip   = var.provider_service_zip
  payment_service_zip    = var.payment_service_zip
  sos_service_zip        = var.sos_service_zip
  authorizer_zip         = var.authorizer_zip
  shared_layer_zip       = var.shared_layer_zip
  # TODO: Uncomment when ElastiCache is enabled
  # redis_endpoint         = module.elasticache.redis_endpoint
  # rds_endpoint           = module.rds.cluster_endpoint
  tags                   = var.tags
}

# Monitoring Module
module "monitoring" {
  source = "../../modules/monitoring"

  environment        = var.environment
  aws_region         = var.aws_region
  alert_email        = var.alert_email
  log_retention_days = var.log_retention_days
  tags               = var.tags
}
