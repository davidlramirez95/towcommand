import type { APIGatewayRequestAuthorizerEvent, APIGatewayAuthorizerResult } from 'aws-lambda';

export interface TokenClaims {
  sub: string;
  email: string;
  phone_number?: string;
  'custom:user_type': string;
  'custom:trust_tier'?: string;
  'custom:provider_id'?: string;
}

export async function jwtAuthorizer(
  event: APIGatewayRequestAuthorizerEvent,
): Promise<APIGatewayAuthorizerResult> {
  const token = event.headers?.['authorization']?.replace('Bearer ', '') ?? '';

  if (!token) {
    return generatePolicy('anonymous', 'Deny', event.methodArn);
  }

  try {
    // In production, verify JWT with Cognito JWKS
    // For now, decode and validate structure
    const payload = decodeToken(token);

    if (!payload?.sub) {
      return generatePolicy('anonymous', 'Deny', event.methodArn);
    }

    const policy = generatePolicy(payload.sub, 'Allow', event.methodArn);
    policy.context = {
      userId: payload.sub,
      userType: payload['custom:user_type'] ?? 'customer',
      trustTier: payload['custom:trust_tier'] ?? 'basic',
      providerId: payload['custom:provider_id'] ?? '',
    };

    return policy;
  } catch {
    return generatePolicy('anonymous', 'Deny', event.methodArn);
  }
}

function decodeToken(token: string): TokenClaims | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    const payload = JSON.parse(Buffer.from(parts[1], 'base64url').toString());
    return payload as TokenClaims;
  } catch {
    return null;
  }
}

function generatePolicy(
  principalId: string,
  effect: 'Allow' | 'Deny',
  resource: string,
): APIGatewayAuthorizerResult {
  return {
    principalId,
    policyDocument: {
      Version: '2012-10-17',
      Statement: [{ Action: 'execute-api:Invoke', Effect: effect, Resource: resource }],
    },
  };
}
