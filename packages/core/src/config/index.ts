/**
 * Configuration class - loads environment variables
 * Following 12-Factor: Config stored in environment
 * Pattern adapted from gutguard-ai Config singleton
 */

export interface CognitoConfig {
  userPoolId: string;
  clientId: string;
  region: string;
}

export interface DynamoDBConfig {
  tableName: string;
  endpoint?: string;
}

export interface EventBridgeConfig {
  eventBusName: string;
}

export interface S3Config {
  evidenceBucket: string;
  uploadsBucket: string;
}

export interface PaymentConfig {
  paymongoSecretKey: string;
  paymongoPublicKey: string;
}

export interface LoggingConfig {
  level: string;
}

export interface AppConfig {
  stage: string;
  region: string;
  cognito: CognitoConfig;
  dynamodb: DynamoDBConfig;
  eventbridge: EventBridgeConfig;
  s3: S3Config;
  payment: PaymentConfig;
  logging: LoggingConfig;
}

export class Config implements AppConfig {
  public readonly stage: string;
  public readonly region: string;
  public readonly cognito: CognitoConfig;
  public readonly dynamodb: DynamoDBConfig;
  public readonly eventbridge: EventBridgeConfig;
  public readonly s3: S3Config;
  public readonly payment: PaymentConfig;
  public readonly logging: LoggingConfig;

  private static instance: Config | null = null;

  constructor(env: NodeJS.ProcessEnv = process.env) {
    this.stage = env.STAGE || env.ENVIRONMENT || 'dev';
    this.region = env.AWS_REGION || 'ap-southeast-1';

    this.cognito = {
      userPoolId: env.COGNITO_USER_POOL_ID || '',
      clientId: env.COGNITO_CLIENT_ID || '',
      region: env.COGNITO_REGION || this.region,
    };

    this.dynamodb = {
      tableName: env.DYNAMODB_TABLE_NAME || `TowCommand-${this.stage}`,
      endpoint: env.DYNAMODB_ENDPOINT || undefined,
    };

    this.eventbridge = {
      eventBusName: env.EVENT_BUS_NAME || `towcommand-${this.stage}`,
    };

    this.s3 = {
      evidenceBucket: env.S3_EVIDENCE_BUCKET || `towcommand-evidence-${this.stage}`,
      uploadsBucket: env.S3_UPLOADS_BUCKET || `towcommand-uploads-${this.stage}`,
    };

    this.payment = {
      paymongoSecretKey: env.PAYMONGO_SECRET_KEY || '',
      paymongoPublicKey: env.PAYMONGO_PUBLIC_KEY || '',
    };

    this.logging = {
      level: env.LOG_LEVEL || 'info',
    };
  }

  public static getInstance(env?: NodeJS.ProcessEnv): Config {
    if (!Config.instance) {
      Config.instance = new Config(env);
    }
    return Config.instance;
  }

  public static resetInstance(): void {
    Config.instance = null;
  }

  /** Check if running in local development */
  public get isLocal(): boolean {
    return this.stage === 'dev' && !!this.dynamodb.endpoint;
  }

  /** Check if running in production */
  public get isProduction(): boolean {
    return this.stage === 'prod';
  }
}

export const config = Config.getInstance();
export default config;
