export interface DriverArrivingContext {
  customerName: string;
  driverName: string;
  plateNumber: string;
  eta: string;
  driverPhone: string;
}

export function renderDriverArrivingSMS(context: DriverArrivingContext): string {
  return `Hi ${context.customerName}, ${context.driverName} is arriving soon in ${context.plateNumber}. ETA: ${context.eta}`;
}

export function renderDriverArrivingEmail(context: DriverArrivingContext): string {
  // TODO: Implement email template rendering
  return `
    <h2>Driver Arriving</h2>
    <p>Hello ${context.customerName},</p>
    <p>Your driver <strong>${context.driverName}</strong> is arriving shortly.</p>
    <p>Vehicle: ${context.plateNumber}</p>
    <p>ETA: ${context.eta}</p>
    <p>Contact: ${context.driverPhone}</p>
  `;
}
