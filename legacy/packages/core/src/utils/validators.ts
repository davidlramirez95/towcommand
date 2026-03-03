import { z } from 'zod';
import { ServiceType, WeightClass, BookingStatus, PaymentMethod } from '../types';

export const geoLocationSchema = z.object({
  lat: z.number().min(-90).max(90),
  lng: z.number().min(-180).max(180),
  address: z.string().optional(),
});

export const phPhoneSchema = z.string().regex(
  /^\+63[0-9]{10}$/,
  'Must be a valid Philippine phone number (+63XXXXXXXXXX)',
);

export const createBookingSchema = z.object({
  pickupLocation: geoLocationSchema,
  dropoffLocation: geoLocationSchema,
  vehicleId: z.string().min(1),
  serviceType: z.nativeEnum(ServiceType),
  paymentMethodId: z.string().min(1),
  estimateId: z.string().min(1),
  notes: z.string().max(500).optional(),
});

export const estimateBookingSchema = z.object({
  pickupLocation: geoLocationSchema,
  dropoffLocation: geoLocationSchema,
  serviceType: z.nativeEnum(ServiceType),
  weightClass: z.nativeEnum(WeightClass),
});

export const updateLocationSchema = z.object({
  lat: z.number().min(4.5).max(21.5),
  lng: z.number().min(116.0).max(127.0),
  heading: z.number().min(0).max(360),
  speed: z.number().min(0).max(200),
});

export const updateStatusSchema = z.object({
  status: z.nativeEnum(BookingStatus),
  metadata: z.record(z.unknown()).optional(),
});

export const verifyOTPSchema = z.object({
  type: z.enum(['PICKUP', 'DROPOFF']),
  code: z.string().length(6),
  providerLocation: geoLocationSchema,
});

export const submitRatingSchema = z.object({
  bookingId: z.string().min(1),
  rating: z.number().int().min(1).max(5),
  comment: z.string().max(1000).optional(),
  tags: z.array(z.string()).optional(),
});

export const providerRegistrationSchema = z.object({
  name: z.string().min(2).max(100),
  phone: phPhoneSchema,
  email: z.string().email(),
  truckType: z.string(),
  plateNumber: z.string().regex(/^[A-Z0-9]{3,4}-?[A-Z0-9]{3,4}$/i, 'Invalid plate number format'),
  ltoRegistration: z.string().min(1),
  maxWeightCapacityKg: z.number().positive(),
  serviceAreas: z.array(z.string()).min(1),
});

export const initiatePaymentSchema = z.object({
  bookingId: z.string().min(1),
  method: z.nativeEnum(PaymentMethod),
  amount: z.number().positive(),
});

export type CreateBookingInput = z.infer<typeof createBookingSchema>;
export type EstimateBookingInput = z.infer<typeof estimateBookingSchema>;
export type UpdateLocationInput = z.infer<typeof updateLocationSchema>;
export type VerifyOTPInput = z.infer<typeof verifyOTPSchema>;
export type SubmitRatingInput = z.infer<typeof submitRatingSchema>;
export type ProviderRegistrationInput = z.infer<typeof providerRegistrationSchema>;
