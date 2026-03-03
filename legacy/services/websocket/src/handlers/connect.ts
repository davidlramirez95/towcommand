import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { SessionCache } from '@towcommand/cache';

const sessionCache = new SessionCache();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const connectionId = event.requestContext.connectionId as string;
    const userId = event.queryStringParameters?.userId;
    const token = event.queryStringParameters?.token;

    if (!userId || !token) {
      return { statusCode: 401, body: 'Missing userId or token' };
    }

    // Store connection mapping: userId -> connectionId
    await sessionCache.setWebSocketConnection(userId, connectionId);

    console.log(`WebSocket connected: user=${userId}, connection=${connectionId}`);

    return { statusCode: 200, body: 'Connected' };
  } catch (error) {
    console.error('WebSocket connect error:', error);
    return { statusCode: 500, body: 'Connection failed' };
  }
}
