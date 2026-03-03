# X-Ray Sampling Rule
resource "aws_xray_sampling_rule" "main" {
  rule_name      = "towcommand-${var.environment}"
  priority       = 1000
  version        = 1
  reservoir_size = 1
  fixed_rate     = 0.05
  url_path       = "*"
  host           = "*"
  http_method    = "*"
  service_type   = "*"
  service_name   = "*"
  resource_arn   = "*"

  attributes = {
    Environment = var.environment
    Service     = "towcommand"
  }
}

# X-Ray Group for Errors
resource "aws_xray_group" "errors" {
  group_name        = "towcommand-errors-${var.environment}"
  filter_expression = "service(\"*\") { error = true }"
  insights_enabled  = true

  tags = var.tags
}

# X-Ray Group for Throttling
resource "aws_xray_group" "throttling" {
  group_name        = "towcommand-throttling-${var.environment}"
  filter_expression = "service(\"*\") { http.status >= 400 }"
  insights_enabled  = true

  tags = var.tags
}

# X-Ray Insights for Anomaly Detection
resource "aws_xray_insight_rule" "high_error_rate" {
  rule_name      = "high-error-rate-${var.environment}"
  filter_string  = "service(\"*\") { error = true }"
  group_arn      = aws_xray_group.errors.arn
  priority       = 100

  tags = var.tags
}
