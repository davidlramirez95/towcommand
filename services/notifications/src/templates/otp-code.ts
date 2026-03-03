export interface OTPCodeContext {
  code: string;
  expiryMinutes: number;
}

export function renderOTPCodeSMS(context: OTPCodeContext): string {
  return `Your TowCommand verification code is ${context.code}. Valid for ${context.expiryMinutes} minutes. Do not share this code.`;
}

export function renderOTPCodeEmail(context: OTPCodeContext): string {
  // TODO: Implement email template rendering
  return `
    <h2>Verification Code</h2>
    <p>Your TowCommand verification code is:</p>
    <h1 style="font-size: 32px; font-weight: bold;">${context.code}</h1>
    <p>This code is valid for ${context.expiryMinutes} minutes.</p>
    <p>Do not share this code with anyone.</p>
  `;
}
