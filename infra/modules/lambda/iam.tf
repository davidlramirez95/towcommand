# Booking Service IAM Role
resource "aws_iam_role" "lambda_booking" {
  name = "lambda-booking-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_booking_basic" {
  role       = aws_iam_role.lambda_booking.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_booking_xray" {
  role       = aws_iam_role.lambda_booking.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy" "lambda_booking_dynamodb" {
  name = "lambda-booking-dynamodb-${var.environment}"
  role = aws_iam_role.lambda_booking.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [var.dynamodb_table_arn, "${var.dynamodb_table_arn}/index/*"]
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_booking_events" {
  name = "lambda-booking-events-${var.environment}"
  role = aws_iam_role.lambda_booking.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "events:PutEvents"
        ]
        Resource = var.event_bus_arn
      }
    ]
  })
}

# Provider Service IAM Role
resource "aws_iam_role" "lambda_provider" {
  name = "lambda-provider-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_provider_basic" {
  role       = aws_iam_role.lambda_provider.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_provider_xray" {
  role       = aws_iam_role.lambda_provider.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy" "lambda_provider_dynamodb" {
  name = "lambda-provider-dynamodb-${var.environment}"
  role = aws_iam_role.lambda_provider.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [var.dynamodb_table_arn, "${var.dynamodb_table_arn}/index/*"]
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_provider_events" {
  name = "lambda-provider-events-${var.environment}"
  role = aws_iam_role.lambda_provider.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "events:PutEvents"
        ]
        Resource = var.event_bus_arn
      }
    ]
  })
}

# Payment Service IAM Role
resource "aws_iam_role" "lambda_payment" {
  name = "lambda-payment-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_payment_basic" {
  role       = aws_iam_role.lambda_payment.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_payment_xray" {
  role       = aws_iam_role.lambda_payment.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy" "lambda_payment_dynamodb" {
  name = "lambda-payment-dynamodb-${var.environment}"
  role = aws_iam_role.lambda_payment.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [var.dynamodb_table_arn, "${var.dynamodb_table_arn}/index/*"]
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_payment_events" {
  name = "lambda-payment-events-${var.environment}"
  role = aws_iam_role.lambda_payment.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "events:PutEvents"
        ]
        Resource = var.event_bus_arn
      }
    ]
  })
}

# SOS Service IAM Role
resource "aws_iam_role" "lambda_sos" {
  name = "lambda-sos-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_sos_basic" {
  role       = aws_iam_role.lambda_sos.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_sos_xray" {
  role       = aws_iam_role.lambda_sos.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy" "lambda_sos_dynamodb" {
  name = "lambda-sos-dynamodb-${var.environment}"
  role = aws_iam_role.lambda_sos.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [var.dynamodb_table_arn, "${var.dynamodb_table_arn}/index/*"]
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_sos_events" {
  name = "lambda-sos-events-${var.environment}"
  role = aws_iam_role.lambda_sos.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "events:PutEvents"
        ]
        Resource = var.event_bus_arn
      }
    ]
  })
}

# Authorizer Lambda IAM Role
resource "aws_iam_role" "lambda_authorizer" {
  name = "lambda-authorizer-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_authorizer_basic" {
  role       = aws_iam_role.lambda_authorizer.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_authorizer_xray" {
  role       = aws_iam_role.lambda_authorizer.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

# TODO: Uncomment when ElastiCache is provisioned and budget allows
# resource "aws_iam_role_policy" "lambda_booking_elasticache" {
#   name = "lambda-booking-elasticache-${var.environment}"
#   role = aws_iam_role.lambda_booking.id
#
#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect = "Allow"
#         Action = [
#           "elasticache:Connect"
#         ]
#         Resource = "*"
#       }
#     ]
#   })
# }
#
# resource "aws_iam_role_policy" "lambda_provider_elasticache" {
#   name = "lambda-provider-elasticache-${var.environment}"
#   role = aws_iam_role.lambda_provider.id
#
#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect = "Allow"
#         Action = [
#           "elasticache:Connect"
#         ]
#         Resource = "*"
#       }
#     ]
#   })
# }
