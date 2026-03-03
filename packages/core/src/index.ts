// Core types and constants
export * from './types';
export * from './constants';
export * from './errors';

// Config (gutguard-ai pattern)
export { Config, config } from './config/index.js';
export type { AppConfig, CognitoConfig, DynamoDBConfig, EventBridgeConfig, S3Config, PaymentConfig, LoggingConfig } from './config/index.js';

// Feature Flags (gutguard-ai pattern)
export { getFeatureFlags, isFeatureEnabled, resetFeatureFlagCache } from './config/feature-flags.js';
export type { FeatureFlags, UserTier } from './config/feature-flags.js';

// Response Builder (gutguard-ai pattern)
export { ResponseBuilder, success, created, error, errors } from './utils/response.js';
export type { ErrorBody, ErrorCode } from './utils/response.js';

// Logger (gutguard-ai pattern)
export { Logger, ChildLogger, logger } from './utils/logger.js';
export type { LogLevel, LogMeta, LogEntry } from './utils/logger.js';

// Handler Base (gutguard-ai pattern)
export { BaseHandler } from './utils/handler.js';
export type { HandlerFunction, HandlerMap } from './utils/handler.js';

// Other utilities
export * from './utils/geo';
export * from './utils/pricing';
export * from './utils/otp';
export * from './utils/validators';
