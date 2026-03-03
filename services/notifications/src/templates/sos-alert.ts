export interface SOSAlertContext {
  bookingId: string;
  customerName: string;
  location: string;
  severity: 'low' | 'medium' | 'high';
  description: string;
}

export function renderSOSAlertSMS(context: SOSAlertContext): string {
  return `ALERT: SOS from booking ${context.bookingId}. Customer: ${context.customerName}. Location: ${context.location}. Severity: ${context.severity}`;
}

export function renderSOSAlertEmail(context: SOSAlertContext): string {
  // TODO: Implement email template rendering
  return `
    <h2 style="color: red;">SOS ALERT</h2>
    <p><strong>Booking ID:</strong> ${context.bookingId}</p>
    <p><strong>Customer:</strong> ${context.customerName}</p>
    <p><strong>Location:</strong> ${context.location}</p>
    <p><strong>Severity:</strong> ${context.severity}</p>
    <p><strong>Description:</strong> ${context.description}</p>
  `;
}
