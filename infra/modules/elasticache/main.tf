resource "aws_elasticache_subnet_group" "main" {
  name       = "towcommand-${var.environment}"
  subnet_ids = var.private_subnet_ids

  tags = var.tags
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "towcommand-redis-${var.environment}"
  engine               = "redis"
  node_type            = var.redis_node_type
  num_cache_nodes      = var.redis_num_nodes
  parameter_group_name = aws_elasticache_parameter_group.redis.name
  engine_version       = var.redis_engine_version
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [aws_security_group.elasticache.id]

  automatic_failover_enabled = var.redis_automatic_failover
  multi_az_enabled          = var.redis_multi_az

  at_rest_encryption_enabled = true
  transit_encryption_enabled = var.redis_transit_encryption
  auth_token                 = var.redis_auth_token != "" ? var.redis_auth_token : null

  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.redis_slow_log.name
    destination_type = "cloudwatch-logs"
    log_format       = "json"
    log_type         = "slow-log"
  }

  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.redis_engine_log.name
    destination_type = "cloudwatch-logs"
    log_format       = "json"
    log_type         = "engine-log"
  }

  notification_topic_arn = aws_sns_topic.elasticache_alerts.arn

  snapshot_retention_limit = var.redis_snapshot_retention
  snapshot_window          = "03:00-05:00"
  maintenance_window       = "mon:05:00-mon:07:00"

  tags = var.tags
}

resource "aws_elasticache_parameter_group" "redis" {
  family = "redis7"
  name   = "towcommand-redis-${var.environment}"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }

  parameter {
    name  = "timeout"
    value = "300"
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "redis_slow_log" {
  name              = "/aws/elasticache/towcommand-redis-${var.environment}/slow-log"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "redis_engine_log" {
  name              = "/aws/elasticache/towcommand-redis-${var.environment}/engine-log"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_sns_topic" "elasticache_alerts" {
  name = "elasticache-alerts-${var.environment}"

  tags = var.tags
}

resource "aws_security_group" "elasticache" {
  name        = "elasticache-${var.environment}"
  description = "Security group for ElastiCache"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = var.private_subnet_cidrs
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}
