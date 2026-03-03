output "event_bus_name" {
  description = "Name of the EventBridge event bus"
  value       = aws_cloudwatch_event_bus.main.name
}

output "event_bus_arn" {
  description = "ARN of the EventBridge event bus"
  value       = aws_cloudwatch_event_bus.main.arn
}

output "booking_created_rule_arn" {
  description = "ARN of the booking created rule"
  value       = aws_cloudwatch_event_rule.booking_created.arn
}

output "booking_completed_rule_arn" {
  description = "ARN of the booking completed rule"
  value       = aws_cloudwatch_event_rule.booking_completed.arn
}

output "sos_activated_rule_arn" {
  description = "ARN of the SOS activated rule"
  value       = aws_cloudwatch_event_rule.sos_activated.arn
}

output "payment_completed_rule_arn" {
  description = "ARN of the payment completed rule"
  value       = aws_cloudwatch_event_rule.payment_completed.arn
}

output "provider_online_rule_arn" {
  description = "ARN of the provider online rule"
  value       = aws_cloudwatch_event_rule.provider_online.arn
}

output "schema_registry_name" {
  description = "Name of the schema registry"
  value       = aws_schemas_registry.main.name
}
