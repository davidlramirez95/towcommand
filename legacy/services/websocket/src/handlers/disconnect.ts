import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { SessionCache, getRedisClient } from '@towcommand/cache';

const sessionCache = new SessionCache();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const connectionId = event.requestContext.connectionId as string;

    // Scan for the userId associated with this connectionId and remove it
    // In production, you'd maintain a reverse mapping (connectionId -> userId)
    // For now, the connectionId is passed in the disconnect context
    const userId = event.queryStringParameters?.userId;

    if (userId) {
      await sessionCache.removeWebSocketConnection(userId);
    }

    console.log(`WebSocket disconnected: connection=${connectionId}`);

    return { statusCode: 200, body: 'Disconnected' };
  } catch (error) {
    console.error('WebSocket disconnect error:', error);
    return { statusCode: 500, body: 'Disconnect failed' };
  }
}
