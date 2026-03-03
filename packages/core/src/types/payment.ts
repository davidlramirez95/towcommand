export enum PaymentMethod {
  GCASH = 'gcash',
  MAYA = 'maya',
  CARD = 'card',
  CASH = 'cash',
  CORPORATE = 'corporate',
}

export enum PaymentStatus {
  PENDING = 'pending',
  HELD = 'held',
  CAPTURED = 'captured',
  REFUNDED = 'refunded',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

export interface Payment {
  paymentId: string;
  bookingId: string;
  userId: string;
  amount: number;
  currency: 'PHP';
  method: PaymentMethod;
  status: PaymentStatus;
  holdAmount?: number;
  gatewayRef?: string;
  capturedAt?: string;
  refundedAt?: string;
  refundReason?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ProviderPayout {
  payoutId: string;
  providerId: string;
  bookingId: string;
  grossAmount: number;
  commission: number;
  commissionRate: number;
  netAmount: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  payoutMethod: PaymentMethod;
  scheduledAt: string;
  completedAt?: string;
}
