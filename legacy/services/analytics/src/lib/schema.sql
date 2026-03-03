-- Analytics tables for TowCommand

CREATE TABLE IF NOT EXISTS bookings_analytics (
  booking_id VARCHAR(255) PRIMARY KEY,
  customer_id VARCHAR(255) NOT NULL,
  provider_id VARCHAR(255),
  service_type VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  pickup_location GEOGRAPHY,
  dropoff_location GEOGRAPHY,
  total_cost DECIMAL(10, 2),
  currency VARCHAR(3),
  booking_date TIMESTAMP NOT NULL,
  completion_date TIMESTAMP,
  duration_minutes INTEGER,
  distance_km DECIMAL(8, 2),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_booking_date (booking_date),
  INDEX idx_customer_id (customer_id),
  INDEX idx_provider_id (provider_id),
  INDEX idx_service_type (service_type)
);

CREATE TABLE IF NOT EXISTS provider_performance (
  provider_id VARCHAR(255) PRIMARY KEY,
  total_bookings INTEGER DEFAULT 0,
  completed_bookings INTEGER DEFAULT 0,
  cancelled_bookings INTEGER DEFAULT 0,
  average_rating DECIMAL(3, 2),
  total_earnings DECIMAL(12, 2),
  average_response_time_seconds INTEGER,
  completion_rate DECIMAL(5, 2),
  last_booking_date TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_rating (average_rating),
  INDEX idx_earnings (total_earnings)
);

CREATE TABLE IF NOT EXISTS daily_metrics (
  metric_date DATE PRIMARY KEY,
  total_bookings INTEGER DEFAULT 0,
  completed_bookings INTEGER DEFAULT 0,
  total_revenue DECIMAL(15, 2),
  active_providers INTEGER DEFAULT 0,
  active_customers INTEGER DEFAULT 0,
  average_rating DECIMAL(3, 2),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS demand_heatmap (
  grid_cell_id VARCHAR(255),
  timestamp TIMESTAMP,
  demand_level INTEGER,
  booking_count INTEGER,
  average_wait_time_seconds INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (grid_cell_id, timestamp),
  INDEX idx_timestamp (timestamp)
);
