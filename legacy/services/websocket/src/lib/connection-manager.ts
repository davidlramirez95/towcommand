import { ApiGatewayManagementApiClient, PostToConnectionCommand } from '@aws-sdk/client-apigatewaymanagementapi';

export class ConnectionManager {
  private client: ApiGatewayManagementApiClient;
  private endpoint: string;

  constructor() {
    this.endpoint = process.env.WEBSOCKET_ENDPOINT ?? '';
    this.client = new ApiGatewayManagementApiClient({
      endpoint: this.endpoint,
    });
  }

  async sendMessage(connectionId: string, data: unknown): Promise<void> {
    try {
      await this.client.send(
        new PostToConnectionCommand({
          ConnectionId: connectionId,
          Data: JSON.stringify(data),
        }),
      );
    } catch (error) {
      console.error(`Failed to send message to connection ${connectionId}:`, error);
      throw error;
    }
  }

  async broadcast(connectionIds: string[], data: unknown): Promise<void> {
    const promises = connectionIds.map((id) => this.sendMessage(id, data));
    await Promise.allSettled(promises);
  }
}
