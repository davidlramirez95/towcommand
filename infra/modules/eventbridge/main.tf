resource "aws_cloudwatch_event_bus" "main" {
  name = "towcommand-${var.environment}"
  tags = var.tags
}

resource "aws_cloudwatch_event_rule" "booking_created" {
  name           = "booking-created-${var.environment}"
  event_bus_name = aws_cloudwatch_event_bus.main.name
  event_pattern  = jsonencode({
    source      = ["tc.booking"]
    detail-type = ["BookingCreated"]
  })
}

resource "aws_cloudwatch_event_rule" "booking_completed" {
  name           = "booking-completed-${var.environment}"
  event_bus_name = aws_cloudwatch_event_bus.main.name
  event_pattern  = jsonencode({
    source      = ["tc.booking"]
    detail-type = ["BookingCompleted"]
  })
}

resource "aws_cloudwatch_event_rule" "sos_activated" {
  name           = "sos-activated-${var.environment}"
  event_bus_name = aws_cloudwatch_event_bus.main.name
  event_pattern  = jsonencode({
    source      = ["tc.sos"]
    detail-type = ["SOSActivated"]
  })
}

resource "aws_cloudwatch_event_rule" "payment_completed" {
  name           = "payment-completed-${var.environment}"
  event_bus_name = aws_cloudwatch_event_bus.main.name
  event_pattern  = jsonencode({
    source      = ["tc.payment"]
    detail-type = ["PaymentCompleted"]
  })
}

resource "aws_cloudwatch_event_rule" "provider_online" {
  name           = "provider-online-${var.environment}"
  event_bus_name = aws_cloudwatch_event_bus.main.name
  event_pattern  = jsonencode({
    source      = ["tc.provider"]
    detail-type = ["ProviderOnline"]
  })
}
