import type { EventBridgeEvent } from 'aws-lambda';
import { UserRepository, ProviderRepository, BookingRepository } from '@towcommand/db';
import { SMSChannel } from './channels/sms';
import { EmailChannel } from './channels/email';
import { PushChannel } from './channels/push';
import { renderBookingConfirmedSMS, renderBookingConfirmedEmail } from './templates/booking-confirmed';
import { renderDriverArrivingSMS } from './templates/driver-arriving';
import { renderSOSAlertSMS, renderSOSAlertEmail } from './templates/sos-alert';

const smsChannel = new SMSChannel();
const emailChannel = new EmailChannel();
const pushChannel = new PushChannel();
const userRepo = new UserRepository();
const providerRepo = new ProviderRepository();
const bookingRepo = new BookingRepository();

export async function handler(event: EventBridgeEvent<string, Record<string, unknown>>): Promise<void> {
  const eventType = event['detail-type'];
  const detail = event.detail;

  try {
    console.log(`Processing notification: ${eventType}`);

    switch (eventType) {
      case 'ProviderMatched':
        await handleProviderMatched(detail);
        break;
      case 'BookingStatusChanged':
        await handleBookingStatusChanged(detail);
        break;
      case 'DriverArrived':
        await handleDriverArrived(detail);
        break;
      case 'PaymentCompleted':
        await handlePaymentCompleted(detail);
        break;
      case 'SOSActivated':
        await handleSOSActivated(detail);
        break;
      case 'UserRegistered':
        await handleUserRegistered(detail);
        break;
      default:
        console.log(`No notification handler for event: ${eventType}`);
    }
  } catch (error) {
    console.error(`Notification handler error for ${eventType}:`, error);
    // Don't throw - notification failures should not block event processing
  }
}

async function handleProviderMatched(detail: Record<string, unknown>): Promise<void> {
  const bookingId = detail.bookingId as string;
  const providerName = detail.providerName as string;
  const providerPhone = detail.providerPhone as string;
  const truckPlate = detail.truckPlate as string;
  const eta = detail.eta as number;

  const booking = await bookingRepo.getById(bookingId);
  if (!booking) return;

  const customer = await userRepo.getById(booking.customerId);
  if (!customer) return;

  const smsText = renderBookingConfirmedSMS({
    bookingId,
    customerName: customer.name,
    providerName,
    pickupAddress: booking.pickupLocation.address ?? 'your location',
    estimatedArrival: `${eta} minutes`,
    totalCost: booking.price.total,
    providerPhone,
  });

  await Promise.allSettled([
    customer.phone ? smsChannel.send(customer.phone, smsText) : Promise.resolve(),
    customer.email ? emailChannel.send(
      customer.email,
      `TowCommand: Provider Assigned - ${bookingId}`,
      renderBookingConfirmedEmail({
        bookingId,
        customerName: customer.name,
        providerName,
        pickupAddress: booking.pickupLocation.address ?? 'your location',
        estimatedArrival: `${eta} minutes`,
        totalCost: booking.price.total,
        providerPhone,
      }),
    ) : Promise.resolve(),
    pushChannel.send(customer.userId, 'Provider Assigned', `${providerName} is on the way! ETA: ${eta} min`),
  ]);
}

async function handleBookingStatusChanged(detail: Record<string, unknown>): Promise<void> {
  const bookingId = detail.bookingId as string;
  const newStatus = detail.newStatus as string;

  const booking = await bookingRepo.getById(bookingId);
  if (!booking) return;

  const customer = await userRepo.getById(booking.customerId);
  if (!customer) return;

  const statusMessages: Record<string, string> = {
    EN_ROUTE: 'Your driver is on the way!',
    ARRIVED: 'Your driver has arrived at the pickup location.',
    IN_TRANSIT: 'Your vehicle is being transported.',
    ARRIVED_DROPOFF: 'Your vehicle has arrived at the destination.',
    COMPLETED: 'Your towing service is complete. Please rate your experience!',
    CANCELLED: 'Your booking has been cancelled.',
  };

  const message = statusMessages[newStatus];
  if (!message) return;

  await Promise.allSettled([
    customer.phone ? smsChannel.send(customer.phone, `[TowCommand] ${message} Booking: ${bookingId}`) : Promise.resolve(),
    pushChannel.send(customer.userId, 'Booking Update', message),
  ]);
}

async function handleDriverArrived(detail: Record<string, unknown>): Promise<void> {
  const bookingId = detail.bookingId as string;
  const providerId = detail.providerId as string;

  const booking = await bookingRepo.getById(bookingId);
  if (!booking) return;

  const [customer, provider] = await Promise.all([
    userRepo.getById(booking.customerId),
    providerRepo.getById(providerId),
  ]);

  if (!customer || !provider) return;

  const smsText = renderDriverArrivingSMS({
    customerName: customer.name,
    driverName: provider.name,
    plateNumber: provider.plateNumber,
    eta: 'arrived',
    driverPhone: provider.phone,
  });

  await Promise.allSettled([
    customer.phone ? smsChannel.send(customer.phone, smsText) : Promise.resolve(),
    pushChannel.send(customer.userId, 'Driver Arrived', `${provider.name} (${provider.plateNumber}) has arrived!`),
  ]);
}

async function handlePaymentCompleted(detail: Record<string, unknown>): Promise<void> {
  const bookingId = detail.bookingId as string;
  const amount = detail.amount as number;

  const booking = await bookingRepo.getById(bookingId);
  if (!booking) return;

  const customer = await userRepo.getById(booking.customerId);
  if (!customer) return;

  await Promise.allSettled([
    customer.phone ? smsChannel.send(
      customer.phone,
      `[TowCommand] Payment of ₱${amount} received for booking ${bookingId}. Salamat!`,
    ) : Promise.resolve(),
    customer.email ? emailChannel.send(
      customer.email,
      `TowCommand: Payment Receipt - ${bookingId}`,
      `<h2>Payment Received</h2><p>Amount: ₱${amount}</p><p>Booking: ${bookingId}</p><p>Salamat po!</p>`,
    ) : Promise.resolve(),
  ]);
}

async function handleSOSActivated(detail: Record<string, unknown>): Promise<void> {
  const alertId = detail.alertId as string;
  const bookingId = detail.bookingId as string | undefined;
  const location = detail.location as { lat: number; lng: number; address?: string };

  const opsEmail = process.env.OPS_ALERT_EMAIL ?? 'ops@towcommand.ph';
  const opsPhone = process.env.OPS_ALERT_PHONE;

  const smsText = renderSOSAlertSMS({
    bookingId: bookingId ?? alertId,
    customerName: 'Customer',
    location: location?.address ?? `${location?.lat}, ${location?.lng}`,
    severity: 'high',
    description: 'SOS activated',
  });

  await Promise.allSettled([
    emailChannel.send(opsEmail, `URGENT: SOS Alert ${alertId}`, renderSOSAlertEmail({
      bookingId: bookingId ?? alertId,
      customerName: 'Customer',
      location: location?.address ?? `${location?.lat}, ${location?.lng}`,
      severity: 'high',
      description: 'SOS alert activated via app',
    })),
    opsPhone ? smsChannel.send(opsPhone, smsText) : Promise.resolve(),
  ]);
}

async function handleUserRegistered(detail: Record<string, unknown>): Promise<void> {
  const userId = detail.userId as string;
  const phone = detail.phone as string;

  if (phone) {
    await smsChannel.send(phone, '[TowCommand PH] Mabuhay! Welcome to TowCommand - Ang Grab ng Towing! Your account is ready.');
  }
}
