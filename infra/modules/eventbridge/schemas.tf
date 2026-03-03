# EventBridge schema registry for event schemas
# This helps with validation and documentation

resource "aws_schemas_registry" "main" {
  name = "towcommand-${var.environment}"

  tags = var.tags
}

# Booking event schema
resource "aws_schemas_schema" "booking_created" {
  name          = "tc.booking@BookingCreated"
  registry_name = aws_schemas_registry.main.name
  type          = "OpenApi3"
  description   = "Event published when a new booking is created"

  content = jsonencode({
    openapi = "3.0.0"
    info = {
      title       = "BookingCreated"
      version     = "1.0.0"
      description = "Event published when a new booking is created"
    }
    paths = {}
    components = {
      schemas = {
        BookingCreated = {
          type = "object"
          properties = {
            booking_id = {
              type = "string"
            }
            user_id = {
              type = "string"
            }
            location = {
              type = "object"
            }
            status = {
              type = "string"
            }
            timestamp = {
              type = "string"
              format = "date-time"
            }
          }
          required = ["booking_id", "user_id", "location", "timestamp"]
        }
      }
    }
  })
}

# Payment event schema
resource "aws_schemas_schema" "payment_completed" {
  name          = "tc.payment@PaymentCompleted"
  registry_name = aws_schemas_registry.main.name
  type          = "OpenApi3"
  description   = "Event published when payment is completed"

  content = jsonencode({
    openapi = "3.0.0"
    info = {
      title       = "PaymentCompleted"
      version     = "1.0.0"
      description = "Event published when payment is completed"
    }
    paths = {}
    components = {
      schemas = {
        PaymentCompleted = {
          type = "object"
          properties = {
            payment_id = {
              type = "string"
            }
            booking_id = {
              type = "string"
            }
            amount = {
              type = "number"
            }
            currency = {
              type = "string"
            }
            status = {
              type = "string"
            }
            timestamp = {
              type = "string"
              format = "date-time"
            }
          }
          required = ["payment_id", "booking_id", "amount", "timestamp"]
        }
      }
    }
  })
}

# SOS event schema
resource "aws_schemas_schema" "sos_activated" {
  name          = "tc.sos@SOSActivated"
  registry_name = aws_schemas_registry.main.name
  type          = "OpenApi3"
  description   = "Event published when SOS is activated"

  content = jsonencode({
    openapi = "3.0.0"
    info = {
      title       = "SOSActivated"
      version     = "1.0.0"
      description = "Event published when SOS is activated"
    }
    paths = {}
    components = {
      schemas = {
        SOSActivated = {
          type = "object"
          properties = {
            sos_id = {
              type = "string"
            }
            booking_id = {
              type = "string"
            }
            user_id = {
              type = "string"
            }
            location = {
              type = "object"
            }
            severity = {
              type = "string"
            }
            timestamp = {
              type = "string"
              format = "date-time"
            }
          }
          required = ["sos_id", "booking_id", "user_id", "timestamp"]
        }
      }
    }
  })
}
