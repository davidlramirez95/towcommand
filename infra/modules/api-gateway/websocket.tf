resource "aws_apigatewayv2_api" "websocket" {
  name          = "towcommand-ws-${var.environment}"
  protocol_type = "WEBSOCKET"
  route_selection_expression = "$request.body.action"

  tags = var.tags
}

resource "aws_apigatewayv2_stage" "websocket" {
  api_id      = aws_apigatewayv2_api.websocket.id
  name        = var.environment
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.websocket.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      status         = "$context.status"
      protocol       = "$context.protocol"
      error          = "$context.error.message"
      errorType      = "$context.error.messageString"
      connectionId   = "$context.connectionId"
      routeKey       = "$context.routeKey"
    })
  }

  default_route_settings {
    throttling_burst_limit = var.websocket_throttle_burst_limit
    throttling_rate_limit  = var.websocket_throttle_rate_limit
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "websocket" {
  name              = "/aws/apigateway/websocket/towcommand-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# Default route
resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.default.id}"
}

# Default integration (placeholder - connect to Lambda or backend)
resource "aws_apigatewayv2_integration" "default" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "MOCK"
}

# Connect route
resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.connect.id}"
}

resource "aws_apigatewayv2_integration" "connect" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "MOCK"
}

# Disconnect route
resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.disconnect.id}"
}

resource "aws_apigatewayv2_integration" "disconnect" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "MOCK"
}
