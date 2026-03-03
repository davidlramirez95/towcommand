import { SNSClient, PublishCommand } from '@aws-sdk/client-sns';

export class SMSChannel {
  private client: SNSClient;

  constructor() {
    this.client = new SNSClient({ region: process.env.AWS_REGION || 'ap-southeast-1' });
  }

  async send(phoneNumber: string, message: string): Promise<string> {
    // TODO: Implement SMS sending via SNS
    // - Validate phone number format
    // - Send via SNS Publish
    // - Log message ID

    try {
      const command = new PublishCommand({
        Message: message,
        PhoneNumber: phoneNumber,
      });
      const response = await this.client.send(command);
      return response.MessageId ?? '';
    } catch (error) {
      console.error('SMS send error:', error);
      throw error;
    }
  }
}
