# SOS Dashboard
resource "aws_cloudwatch_dashboard" "sos" {
  dashboard_name = "towcommand-sos-${var.environment}"

  dashboard_body = jsonencode({
    widgets = [
      {
        type = "metric"
        properties = {
          metrics = [
            ["AWS/Lambda", "Invocations", { stat = "Sum" }],
            [".", "Errors", { stat = "Sum" }],
            [".", "Duration", { stat = "Average" }]
          ]
          period = 300
          stat   = "Average"
          region = var.aws_region
          title  = "SOS Lambda Metrics"
        }
      },
      {
        type = "log"
        properties = {
          query   = "fields @timestamp, @message | stats count() by @message"
          region  = var.aws_region
          title   = "SOS Logs"
        }
      }
    ]
  })
}

# Lambda Errors Alarm
resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "towcommand-lambda-errors-${var.environment}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Sum"
  threshold           = "5"
  alarm_description   = "Alert when Lambda errors exceed threshold"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  tags = var.tags
}

# DynamoDB Throttling Alarm
resource "aws_cloudwatch_metric_alarm" "ddb_throttle" {
  alarm_name          = "towcommand-ddb-throttle-${var.environment}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "ConsumedWriteCapacityUnits"
  namespace           = "AWS/DynamoDB"
  period              = "300"
  statistic           = "Sum"
  threshold           = "80"
  alarm_description   = "Alert when DynamoDB write capacity is nearly exhausted"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  tags = var.tags
}

# Payment Processing Failures
resource "aws_cloudwatch_metric_alarm" "payment_failures" {
  alarm_name          = "towcommand-payment-failures-${var.environment}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "PaymentFailures"
  namespace           = "TowCommand"
  period              = "300"
  statistic           = "Sum"
  threshold           = "10"
  alarm_description   = "Alert when payment processing failures exceed threshold"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  treat_missing_data = "notBreaching"

  tags = var.tags
}

# API Gateway 5xx Errors
resource "aws_cloudwatch_metric_alarm" "api_5xx_errors" {
  alarm_name          = "towcommand-api-5xx-${var.environment}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = "300"
  statistic           = "Sum"
  threshold           = "10"
  alarm_description   = "Alert on API Gateway 5xx errors"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  tags = var.tags
}

# SNS Topic for Alerts
resource "aws_sns_topic" "alerts" {
  name = "towcommand-alerts-${var.environment}"

  tags = var.tags
}

resource "aws_sns_topic_subscription" "alerts_email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# Log Groups for Services
resource "aws_cloudwatch_log_group" "booking_service" {
  name              = "/aws/lambda/towcommand-booking-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "provider_service" {
  name              = "/aws/lambda/towcommand-provider-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "payment_service" {
  name              = "/aws/lambda/towcommand-payment-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "sos_service" {
  name              = "/aws/lambda/towcommand-sos-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# Metric Filters for SOS Activation
resource "aws_cloudwatch_log_metric_filter" "sos_activation" {
  name           = "sos-activation-count"
  log_group_name = aws_cloudwatch_log_group.sos_service.name
  filter_pattern = "[timestamp, request_id, level = \"INFO\", msg = \"SOS Activated\"]"

  metric_transformation {
    name      = "SOSActivationCount"
    namespace = "TowCommand"
    value     = "1"
  }
}

# Metric Filters for Payment Processing
resource "aws_cloudwatch_log_metric_filter" "payment_success" {
  name           = "payment-success-count"
  log_group_name = aws_cloudwatch_log_group.payment_service.name
  filter_pattern = "[timestamp, request_id, level = \"INFO\", msg = \"Payment Processed\"]"

  metric_transformation {
    name      = "PaymentSuccessCount"
    namespace = "TowCommand"
    value     = "1"
  }
}

resource "aws_cloudwatch_log_metric_filter" "payment_failures" {
  name           = "payment-failure-count"
  log_group_name = aws_cloudwatch_log_group.payment_service.name
  filter_pattern = "[timestamp, request_id, level = \"ERROR\", msg = \"Payment Failed\"]"

  metric_transformation {
    name      = "PaymentFailures"
    namespace = "TowCommand"
    value     = "1"
  }
}

# Composite Alarm for Overall Health
resource "aws_cloudwatch_composite_alarm" "overall_health" {
  alarm_name          = "towcommand-overall-health-${var.environment}"
  alarm_description   = "Overall health check for TowCommand"
  actions_enabled     = true
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]

  alarm_rule = join(" OR ", [
    "arn:aws:cloudwatch:${var.aws_region}:${data.aws_caller_identity.current.account_id}:alarm:${aws_cloudwatch_metric_alarm.lambda_errors.alarm_name}",
    "arn:aws:cloudwatch:${var.aws_region}:${data.aws_caller_identity.current.account_id}:alarm:${aws_cloudwatch_metric_alarm.ddb_throttle.alarm_name}",
    "arn:aws:cloudwatch:${var.aws_region}:${data.aws_caller_identity.current.account_id}:alarm:${aws_cloudwatch_metric_alarm.api_5xx_errors.alarm_name}"
  ])

  tags = var.tags
}

data "aws_caller_identity" "current" {}
