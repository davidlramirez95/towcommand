// Payment event detail schemas for EventBridge
// Zod schemas can be added here as needed for event validation

export interface PaymentInitiatedDetail {
  paymentId: string;
  bookingId: string;
  amount: number;
  currency: string;
  initiatedAt: string;
}

export interface PaymentCompletedDetail {
  paymentId: string;
  bookingId: string;
  amount: number;
  completedAt: string;
}

export interface PaymentFailedDetail {
  paymentId: string;
  bookingId: string;
  reason: string;
  failedAt: string;
}
