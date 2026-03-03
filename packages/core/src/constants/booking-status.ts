import { BookingStatus } from '../types';

export interface StatusConfig {
  status: BookingStatus;
  label: string;
  labelFil: string;
  description: string;
  customerVisible: boolean;
  providerVisible: boolean;
  isFinal: boolean;
  canCancel: boolean;
}

export const STATUS_CONFIGS: Record<BookingStatus, StatusConfig> = {
  [BookingStatus.PENDING]: { status: BookingStatus.PENDING, label: 'Finding Provider', labelFil: 'Naghahanap ng Provider', description: 'Searching for available tow provider', customerVisible: true, providerVisible: false, isFinal: false, canCancel: true },
  [BookingStatus.MATCHED]: { status: BookingStatus.MATCHED, label: 'Provider Assigned', labelFil: 'May Provider Na', description: 'A provider has accepted your request', customerVisible: true, providerVisible: true, isFinal: false, canCancel: true },
  [BookingStatus.EN_ROUTE]: { status: BookingStatus.EN_ROUTE, label: 'On The Way', labelFil: 'Papunta Na', description: 'Provider is heading to your location', customerVisible: true, providerVisible: true, isFinal: false, canCancel: true },
  [BookingStatus.ARRIVED]: { status: BookingStatus.ARRIVED, label: 'Provider Arrived', labelFil: 'Nandito Na', description: 'Provider has arrived at pickup location', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.CONDITION_REPORT]: { status: BookingStatus.CONDITION_REPORT, label: 'Documenting Vehicle', labelFil: 'Kinukuhanan ng Litrato', description: 'Provider is documenting vehicle condition', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.OTP_VERIFIED]: { status: BookingStatus.OTP_VERIFIED, label: 'Pickup Verified', labelFil: 'Verified na ang Pickup', description: 'Pickup OTP has been verified', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.LOADING]: { status: BookingStatus.LOADING, label: 'Loading Vehicle', labelFil: 'Kinakargang Sasakyan', description: 'Vehicle is being loaded onto the truck', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.IN_TRANSIT]: { status: BookingStatus.IN_TRANSIT, label: 'In Transit', labelFil: 'Nasa Daan', description: 'Vehicle is being transported', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.ARRIVED_DROPOFF]: { status: BookingStatus.ARRIVED_DROPOFF, label: 'At Drop-off', labelFil: 'Nasa Drop-off Na', description: 'Arrived at destination', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.OTP_DROPOFF]: { status: BookingStatus.OTP_DROPOFF, label: 'Drop-off Verified', labelFil: 'Verified na ang Drop-off', description: 'Drop-off OTP has been verified', customerVisible: true, providerVisible: true, isFinal: false, canCancel: false },
  [BookingStatus.COMPLETED]: { status: BookingStatus.COMPLETED, label: 'Completed', labelFil: 'Tapos Na', description: 'Service completed successfully', customerVisible: true, providerVisible: true, isFinal: true, canCancel: false },
  [BookingStatus.CANCELLED]: { status: BookingStatus.CANCELLED, label: 'Cancelled', labelFil: 'Na-cancel', description: 'Booking was cancelled', customerVisible: true, providerVisible: true, isFinal: true, canCancel: false },
};

export const CANCELLATION_FEES: Record<string, { fee: number; formula?: string }> = {
  PENDING_FREE: { fee: 0 },
  MATCHED: { fee: 100 },
  EN_ROUTE_BASE: { fee: 200, formula: '100 + (distance_km * 30)' },
  ARRIVED: { fee: 500, formula: '500 + distance_fee' },
};
