output "alerts_topic_arn" {
  description = "SNS topic ARN for alerts"
  value       = aws_sns_topic.alerts.arn
}

output "sos_dashboard_name" {
  description = "Name of the SOS CloudWatch dashboard"
  value       = aws_cloudwatch_dashboard.sos.dashboard_name
}

output "error_alarm_arn" {
  description = "ARN of the Lambda error alarm"
  value       = aws_cloudwatch_metric_alarm.lambda_errors.arn
}

output "xray_error_group_arn" {
  description = "ARN of the X-Ray error group"
  value       = aws_xray_group.errors.arn
}
