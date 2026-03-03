output "table_name" {
  description = "Name of the DynamoDB table"
  value       = aws_dynamodb_table.main.name
}

output "table_arn" {
  description = "ARN of the DynamoDB table"
  value       = aws_dynamodb_table.main.arn
}

output "stream_arn" {
  description = "ARN of the DynamoDB Streams"
  value       = aws_dynamodb_table.main.stream_arn
}
