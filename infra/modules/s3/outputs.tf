output "evidence_bucket_name" {
  description = "Name of the evidence S3 bucket"
  value       = aws_s3_bucket.evidence.id
}

output "evidence_bucket_arn" {
  description = "ARN of the evidence S3 bucket"
  value       = aws_s3_bucket.evidence.arn
}

output "evidence_logs_bucket_name" {
  description = "Name of the evidence logs S3 bucket"
  value       = aws_s3_bucket.evidence_logs.id
}
