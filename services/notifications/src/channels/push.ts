import { SNSClient, PublishCommand } from '@aws-sdk/client-sns';

const snsClient = new SNSClient({ region: process.env.AWS_REGION ?? 'ap-southeast-1' });

export class PushChannel {
  private platformAppArn: string;

  constructor() {
    this.platformAppArn = process.env.SNS_PLATFORM_APP_ARN ?? '';
  }

  async send(userId: string, title: string, message: string, data?: Record<string, string>): Promise<string> {
    try {
      // In production, look up the user's device endpoint ARN from DynamoDB
      // For now, publish to a topic-based approach
      const endpointArn = await this.getEndpointArn(userId);
      if (!endpointArn) {
        console.log(`No push endpoint for user ${userId}, skipping push notification`);
        return '';
      }

      const payload = {
        default: message,
        GCM: JSON.stringify({
          notification: { title, body: message },
          data: { ...data, title, body: message },
        }),
        APNS: JSON.stringify({
          aps: { alert: { title, body: message }, sound: 'default' },
          ...data,
        }),
      };

      const response = await snsClient.send(new PublishCommand({
        TargetArn: endpointArn,
        Message: JSON.stringify(payload),
        MessageStructure: 'json',
      }));

      return response.MessageId ?? '';
    } catch (error) {
      console.error(`Push notification error for user ${userId}:`, error);
      return '';
    }
  }

  private async getEndpointArn(userId: string): Promise<string | null> {
    // In production: query DynamoDB for user's registered device endpoint
    // For MVP: return null to gracefully skip push
    return null;
  }
}
