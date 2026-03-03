/**
 * API Gateway Lambda Authorizer for TowCommand
 * Validates JWT tokens from Cognito with role-based access control
 * Pattern adapted from gutguard-ai Authorizer
 */

import type {
  APIGatewayTokenAuthorizerEvent,
  APIGatewayAuthorizerResult,
  PolicyDocument,
  Statement,
} from 'aws-lambda';
import { CognitoJwtVerifier } from 'aws-jwt-verify';
import type { CognitoJwtVerifierSingleUserPool } from 'aws-jwt-verify/cognito-verifier';
import { Logger } from '@towcommand/core/utils/logger';
import { Config } from '@towcommand/core/config';

export interface RoutePermissions {
  [route: string]: string[];
}

export interface AuthorizerContext {
  [key: string]: string | number | boolean;
  userId: string;
  username: string;
  email: string;
  userType: string;
  groups: string;
}

export interface JwtPayload {
  sub: string;
  username?: string;
  email?: string;
  'cognito:groups'?: string[];
  'custom:userType'?: string;
  'custom:phoneNumber'?: string;
}

/**
 * Route-based access control for TowCommand API
 * Maps routes to allowed user groups
 */
export class RouteAccessController {
  private static readonly ROUTE_PERMISSIONS: RoutePermissions = {
    // Booking routes - customers and admins
    '/bookings': ['Customers', 'Admins', 'Dispatchers'],
    '/bookings/create': ['Customers', 'Admins'],
    '/bookings/cancel': ['Customers', 'Providers', 'Admins'],
    '/bookings/estimate': ['Customers', 'Admins'],
    '/bookings/status': ['Customers', 'Providers', 'Admins', 'Dispatchers'],

    // Provider routes
    '/providers/register': ['Providers', 'Admins'],
    '/providers/location': ['Providers'],
    '/providers/availability': ['Providers'],
    '/providers/nearby': ['Customers', 'Admins', 'Dispatchers'],
    '/providers/dashboard': ['Providers', 'Admins'],

    // User routes
    '/users/profile': ['Customers', 'Providers', 'Admins'],
    '/users/vehicles': ['Customers', 'Admins'],

    // Payment routes
    '/payments': ['Customers', 'Providers', 'Admins'],
    '/payments/initiate': ['Customers', 'Admins'],
    '/payments/webhook': ['Customers', 'Providers', 'Admins'],

    // Rating routes
    '/ratings': ['Customers', 'Providers', 'Admins'],
    '/ratings/submit': ['Customers'],

    // SOS routes - everyone
    '/sos': ['Customers', 'Providers', 'Admins', 'Dispatchers'],
    '/sos/activate': ['Customers', 'Providers'],

    // Admin routes
    '/admin/users': ['Admins'],
    '/admin/reports': ['Admins', 'Dispatchers'],
    '/admin/settings': ['Admins'],
    '/admin/disputes': ['Admins'],

    // Diagnosis routes
    '/diagnosis': ['Customers', 'Providers', 'Admins'],
  };

  public checkAccess(route: string | null, userGroups: string[]): boolean {
    if (!route) return true;

    const permissions = RouteAccessController.ROUTE_PERMISSIONS;

    // Exact match
    if (permissions[route]) {
      return userGroups.some(group => permissions[route]?.includes(group));
    }

    // Prefix match (e.g., /bookings/123 matches /bookings)
    for (const [permRoute, allowedGroups] of Object.entries(permissions)) {
      if (route.startsWith(permRoute)) {
        return userGroups.some(group => allowedGroups.includes(group));
      }
    }

    // Default: allow authenticated users
    return true;
  }

  public static getRoutePermissions(): RoutePermissions {
    return { ...RouteAccessController.ROUTE_PERMISSIONS };
  }
}

export class PolicyGenerator {
  public static generate(
    principalId: string,
    effect: 'Allow' | 'Deny',
    resource: string,
    context: AuthorizerContext = {} as AuthorizerContext
  ): APIGatewayAuthorizerResult {
    const statement: Statement = {
      Action: 'execute-api:Invoke',
      Effect: effect,
      Resource: resource,
    };

    const policyDocument: PolicyDocument = {
      Version: '2012-10-17',
      Statement: [statement],
    };

    return {
      principalId,
      policyDocument,
      context,
    };
  }

  public static extractWildcardResource(methodArn: string): string {
    return methodArn.split('/').slice(0, 2).join('/') + '/*';
  }
}

export class TokenExtractor {
  public static extract(authorizationHeader: string | undefined): string | null {
    if (!authorizationHeader) return null;
    if (authorizationHeader.startsWith('Bearer ')) {
      return authorizationHeader.substring(7);
    }
    return authorizationHeader;
  }

  public static extractRouteFromArn(methodArn: string): string | null {
    const parts = methodArn.split('/');
    if (parts.length >= 4) {
      return '/' + parts.slice(3).join('/');
    }
    return null;
  }
}

export class JwtVerifierFactory {
  private static verifier: CognitoJwtVerifierSingleUserPool<{
    userPoolId: string;
    tokenUse: 'access';
    clientId: string;
  }> | null = null;

  public static create(config: Config): CognitoJwtVerifierSingleUserPool<{
    userPoolId: string;
    tokenUse: 'access';
    clientId: string;
  }> {
    if (!JwtVerifierFactory.verifier) {
      JwtVerifierFactory.verifier = CognitoJwtVerifier.create({
        userPoolId: config.cognito.userPoolId,
        tokenUse: 'access',
        clientId: config.cognito.clientId,
      });
    }
    return JwtVerifierFactory.verifier;
  }

  public static reset(): void {
    JwtVerifierFactory.verifier = null;
  }
}

export class Authorizer {
  private readonly config: Config;
  private readonly logger: Logger;
  private readonly routeAccessController: RouteAccessController;

  constructor(
    config?: Config,
    logger?: Logger,
    routeAccessController?: RouteAccessController
  ) {
    this.config = config || Config.getInstance();
    this.logger = logger || Logger.getInstance('authorizer');
    this.routeAccessController = routeAccessController || new RouteAccessController();
  }

  public async authorize(event: APIGatewayTokenAuthorizerEvent): Promise<APIGatewayAuthorizerResult> {
    this.logger.debug('Authorizer invoked', {
      methodArn: event.methodArn,
    });

    const token = TokenExtractor.extract(event.authorizationToken);

    if (!token) {
      this.logger.warn('No token provided');
      throw new Error('Unauthorized');
    }

    try {
      const verifier = JwtVerifierFactory.create(this.config);
      const payload = await verifier.verify(token) as JwtPayload;

      const userGroups = payload['cognito:groups'] || [];
      const userType = payload['custom:userType'] || 'CUSTOMER';
      const route = TokenExtractor.extractRouteFromArn(event.methodArn);

      this.logger.info('Token verified', {
        userId: payload.sub,
        userType,
        groups: userGroups,
        route,
      });

      const hasAccess = this.routeAccessController.checkAccess(route, userGroups);

      if (!hasAccess) {
        this.logger.warn('Access denied', {
          userId: payload.sub,
          groups: userGroups,
          route,
        });
        throw new Error('Unauthorized');
      }

      const context: AuthorizerContext = {
        userId: payload.sub,
        username: payload.username || '',
        email: payload.email || '',
        userType,
        groups: JSON.stringify(userGroups),
      };

      const resourceArn = PolicyGenerator.extractWildcardResource(event.methodArn);
      return PolicyGenerator.generate(payload.sub, 'Allow', resourceArn, context);
    } catch (err) {
      const error = err as Error;
      this.logger.error('Authorization failed', { error: error.message });
      throw new Error('Unauthorized');
    }
  }
}

const authorizer = new Authorizer();

export async function handler(event: APIGatewayTokenAuthorizerEvent): Promise<APIGatewayAuthorizerResult> {
  return authorizer.authorize(event);
}
