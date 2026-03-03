output "booking_function_arn" {
  description = "ARN of the booking Lambda function"
  value       = aws_lambda_function.booking_service.arn
}

output "provider_function_arn" {
  description = "ARN of the provider Lambda function"
  value       = aws_lambda_function.provider_service.arn
}

output "payment_function_arn" {
  description = "ARN of the payment Lambda function"
  value       = aws_lambda_function.payment_service.arn
}

output "sos_function_arn" {
  description = "ARN of the SOS Lambda function"
  value       = aws_lambda_function.sos_service.arn
}

output "authorizer_function_arn" {
  description = "ARN of the authorizer Lambda function"
  value       = aws_lambda_function.authorizer.arn
}

output "shared_layer_arn" {
  description = "ARN of the shared Lambda layer"
  value       = aws_lambda_layer_version.shared.arn
}
