// Booking event detail schemas for EventBridge
// Zod schemas can be added here as needed for event validation

export interface BookingCreatedDetail {
  bookingId: string;
  customerId: string;
  pickupLocation: { lat: number; lng: number };
  dropoffLocation: { lat: number; lng: number };
  createdAt: string;
}

export interface BookingAcceptedDetail {
  bookingId: string;
  providerId: string;
  acceptedAt: string;
}

export interface BookingStatusChangedDetail {
  bookingId: string;
  oldStatus: string;
  newStatus: string;
  changedAt: string;
}
