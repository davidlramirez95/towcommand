# Booking Service Lambda
resource "aws_lambda_function" "booking_service" {
  filename         = var.booking_service_zip
  function_name    = "towcommand-booking-${var.environment}"
  role             = aws_iam_role.lambda_booking.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 30
  memory_size      = var.lambda_memory
  source_code_hash = filebase64sha256(var.booking_service_zip)

  environment {
    variables = {
      ENVIRONMENT    = var.environment
      TABLE_NAME     = var.dynamodb_table_name
      EVENT_BUS_NAME = var.event_bus_name
      # TODO: Uncomment when ElastiCache/RDS is provisioned
      # REDIS_ENDPOINT = var.redis_endpoint
      # RDS_ENDPOINT   = var.rds_endpoint
    }
  }

  tracing_config {
    mode = "Active"
  }

  tags = var.tags
}

# Provider Service Lambda
resource "aws_lambda_function" "provider_service" {
  filename         = var.provider_service_zip
  function_name    = "towcommand-provider-${var.environment}"
  role             = aws_iam_role.lambda_provider.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 30
  memory_size      = var.lambda_memory
  source_code_hash = filebase64sha256(var.provider_service_zip)

  environment {
    variables = {
      ENVIRONMENT    = var.environment
      TABLE_NAME     = var.dynamodb_table_name
      EVENT_BUS_NAME = var.event_bus_name
      # TODO: Uncomment when ElastiCache/RDS is provisioned
      # REDIS_ENDPOINT = var.redis_endpoint
      # RDS_ENDPOINT   = var.rds_endpoint
    }
  }

  tracing_config {
    mode = "Active"
  }

  tags = var.tags
}

# Payment Service Lambda
resource "aws_lambda_function" "payment_service" {
  filename         = var.payment_service_zip
  function_name    = "towcommand-payment-${var.environment}"
  role             = aws_iam_role.lambda_payment.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 30
  memory_size      = var.lambda_memory
  source_code_hash = filebase64sha256(var.payment_service_zip)

  environment {
    variables = {
      ENVIRONMENT    = var.environment
      TABLE_NAME     = var.dynamodb_table_name
      EVENT_BUS_NAME = var.event_bus_name
      # TODO: Uncomment when ElastiCache/RDS is provisioned
      # REDIS_ENDPOINT = var.redis_endpoint
      # RDS_ENDPOINT   = var.rds_endpoint
    }
  }

  tracing_config {
    mode = "Active"
  }

  tags = var.tags
}

# SOS Service Lambda
resource "aws_lambda_function" "sos_service" {
  filename         = var.sos_service_zip
  function_name    = "towcommand-sos-${var.environment}"
  role             = aws_iam_role.lambda_sos.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 30
  memory_size      = var.lambda_memory
  source_code_hash = filebase64sha256(var.sos_service_zip)

  environment {
    variables = {
      ENVIRONMENT    = var.environment
      TABLE_NAME     = var.dynamodb_table_name
      EVENT_BUS_NAME = var.event_bus_name
      # TODO: Uncomment when ElastiCache/RDS is provisioned
      # REDIS_ENDPOINT = var.redis_endpoint
      # RDS_ENDPOINT   = var.rds_endpoint
    }
  }

  tracing_config {
    mode = "Active"
  }

  tags = var.tags
}

# Authorizer Lambda
resource "aws_lambda_function" "authorizer" {
  filename         = var.authorizer_zip
  function_name    = "towcommand-authorizer-${var.environment}"
  role             = aws_iam_role.lambda_authorizer.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 5
  memory_size      = 128
  source_code_hash = filebase64sha256(var.authorizer_zip)

  environment {
    variables = {
      ENVIRONMENT = var.environment
    }
  }

  tracing_config {
    mode = "Active"
  }

  tags = var.tags
}
