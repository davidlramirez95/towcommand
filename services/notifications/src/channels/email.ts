import { SESClient, SendEmailCommand } from '@aws-sdk/client-ses';

export class EmailChannel {
  private client: SESClient;
  private fromEmail: string;

  constructor() {
    this.fromEmail = process.env.FROM_EMAIL ?? 'noreply@towcommand.ph';
    this.client = new SESClient({ region: process.env.AWS_REGION ?? 'ap-southeast-1' });
  }

  async send(toEmail: string, subject: string, htmlBody: string): Promise<string> {
    try {
      const command = new SendEmailCommand({
        Source: this.fromEmail,
        Destination: { ToAddresses: [toEmail] },
        Message: {
          Subject: { Data: subject, Charset: 'UTF-8' },
          Body: {
            Html: {
              Data: this.wrapInTemplate(htmlBody),
              Charset: 'UTF-8',
            },
          },
        },
      });

      const response = await this.client.send(command);
      console.log(`Email sent to ${toEmail}: ${response.MessageId}`);
      return response.MessageId ?? '';
    } catch (error) {
      console.error(`Email send error to ${toEmail}:`, error);
      return '';
    }
  }

  private wrapInTemplate(content: string): string {
    return `
      <!DOCTYPE html>
      <html>
      <head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
      <body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; color: #333;">
        <div style="background: #1a73e8; color: white; padding: 20px; border-radius: 8px 8px 0 0; text-align: center;">
          <h1 style="margin: 0; font-size: 24px;">TowCommand PH</h1>
          <p style="margin: 5px 0 0; opacity: 0.9;">Ang Grab ng Towing</p>
        </div>
        <div style="padding: 20px; border: 1px solid #e0e0e0; border-top: none; border-radius: 0 0 8px 8px;">
          ${content}
        </div>
        <p style="text-align: center; color: #999; font-size: 12px; margin-top: 20px;">
          TowCommand Philippines Inc. | support@towcommand.ph
        </p>
      </body>
      </html>
    `;
  }
}
