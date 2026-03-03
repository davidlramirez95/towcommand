resource "aws_api_gateway_authorizer" "cognito" {
  name          = "cognito-authorizer-${var.environment}"
  rest_api_id   = aws_api_gateway_rest_api.main.id
  type          = "COGNITO_USER_POOLS"
  provider_arns = [var.cognito_user_pool_arn]
  identity_source = "method.request.header.Authorization"
}

resource "aws_apigatewayv2_authorizer" "cognito_websocket" {
  api_id             = aws_apigatewayv2_api.websocket.id
  authorizer_type    = "JWT"
  identity_sources   = ["$request.header.Authorization"]
  name               = "cognito-ws-authorizer-${var.environment}"
  jwt_configuration {
    audience = [var.cognito_client_id]
    issuer   = var.cognito_user_pool_endpoint
  }
}
