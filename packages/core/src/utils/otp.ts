import { randomInt } from 'crypto';

export function generateOTP(length = 6): string {
  const max = Math.pow(10, length);
  const min = Math.pow(10, length - 1);
  return String(randomInt(min, max));
}

export function isOTPExpired(expiresAt: string): boolean {
  return new Date(expiresAt) < new Date();
}

export const OTP_CONFIG = {
  length: 6,
  expiryMinutes: 15,
  maxAttempts: 3,
  cooldownMinutes: 5,
};
