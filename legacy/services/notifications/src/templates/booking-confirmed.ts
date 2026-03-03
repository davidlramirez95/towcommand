export interface BookingConfirmedContext {
  bookingId: string;
  customerName: string;
  providerName: string;
  pickupAddress: string;
  estimatedArrival: string;
  totalCost: number;
  providerPhone: string;
}

export function renderBookingConfirmedSMS(context: BookingConfirmedContext): string {
  return `Hi ${context.customerName}, your booking ${context.bookingId} is confirmed. Driver ${context.providerName} will arrive in ${context.estimatedArrival}. Total: ₱${context.totalCost}`;
}

export function renderBookingConfirmedEmail(context: BookingConfirmedContext): string {
  // TODO: Implement email template rendering
  return `
    <h2>Booking Confirmed</h2>
    <p>Your booking ${context.bookingId} is confirmed.</p>
    <p>Driver: ${context.providerName}</p>
    <p>Phone: ${context.providerPhone}</p>
    <p>Estimated arrival: ${context.estimatedArrival}</p>
    <p>Total cost: ₱${context.totalCost}</p>
  `;
}
