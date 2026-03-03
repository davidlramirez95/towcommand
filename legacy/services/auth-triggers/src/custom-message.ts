import type { CustomMessageTriggerEvent } from 'aws-lambda';

export async function handler(event: CustomMessageTriggerEvent): Promise<CustomMessageTriggerEvent> {
  try {
    const triggerSource = event.triggerSource;
    const code = event.request.codeParameter;
    const name = event.request.userAttributes.name ?? 'ka-TowCommand';

    if (triggerSource === 'CustomMessage_SignUp') {
      event.response.smsMessage =
        `[TowCommand PH] Mabuhay, ${name}! Your verification code is: ${code}. ` +
        `Valid for 15 minutes. Do not share this code.`;
      event.response.emailSubject = 'Welcome to TowCommand PH - Verify Your Account';
      event.response.emailMessage =
        `<h2>Mabuhay, ${name}!</h2>` +
        `<p>Welcome to TowCommand PH - Ang Grab ng Towing!</p>` +
        `<p>Your verification code is: <strong>${code}</strong></p>` +
        `<p>This code expires in 15 minutes.</p>` +
        `<p>If you did not create an account, please ignore this email.</p>` +
        `<p>- Team TowCommand</p>`;
    }

    if (triggerSource === 'CustomMessage_ForgotPassword') {
      event.response.smsMessage =
        `[TowCommand PH] Your password reset code is: ${code}. ` +
        `Valid for 15 minutes. If you did not request this, contact support.`;
      event.response.emailSubject = 'TowCommand PH - Password Reset';
      event.response.emailMessage =
        `<h2>Password Reset Request</h2>` +
        `<p>Hi ${name}, we received a request to reset your password.</p>` +
        `<p>Your reset code is: <strong>${code}</strong></p>` +
        `<p>This code expires in 15 minutes.</p>` +
        `<p>If you did not request this, please ignore this email and your password will remain unchanged.</p>`;
    }

    if (triggerSource === 'CustomMessage_ResendCode') {
      event.response.smsMessage =
        `[TowCommand PH] Your new verification code is: ${code}. Valid for 15 minutes.`;
      event.response.emailSubject = 'TowCommand PH - New Verification Code';
      event.response.emailMessage =
        `<p>Hi ${name}, here is your new verification code: <strong>${code}</strong></p>` +
        `<p>This code expires in 15 minutes.</p>`;
    }

    if (triggerSource === 'CustomMessage_Authentication') {
      event.response.smsMessage =
        `[TowCommand PH] Your login code is: ${code}. Do not share this with anyone.`;
    }

    return event;
  } catch (error) {
    console.error('Custom message error:', error);
    throw error;
  }
}
