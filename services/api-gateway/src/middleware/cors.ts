import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';

export function handleCors(event: APIGatewayProxyEvent): APIGatewayProxyResult | null {
  if (event.httpMethod === 'OPTIONS') {
    return {
      statusCode: 200,
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Headers': 'Content-Type,Authorization',
        'Access-Control-Allow-Methods': 'GET,POST,PATCH,DELETE,OPTIONS',
      },
      body: '',
    };
  }
  return null;
}
