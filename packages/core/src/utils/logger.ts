/**
 * Structured Logger for TowCommand Lambda functions
 * Outputs JSON for CloudWatch Logs Insights compatibility
 * Pattern adapted from gutguard-ai Logger
 */

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

export interface LogMeta {
  [key: string]: unknown;
}

export interface LogEntry {
  timestamp: string;
  level: LogLevel;
  service: string;
  message: string;
  [key: string]: unknown;
}

export class Logger {
  private static readonly LOG_LEVELS: Record<LogLevel, number> = {
    debug: 0,
    info: 1,
    warn: 2,
    error: 3,
  };

  private readonly currentLevel: number;
  private readonly serviceName: string;
  private static instance: Logger | null = null;

  constructor(
    serviceName: string = process.env.SERVICE_NAME || 'towcommand',
    level: string = process.env.LOG_LEVEL || 'info'
  ) {
    this.serviceName = serviceName;
    const normalizedLevel = level.toLowerCase() as LogLevel;
    this.currentLevel = Logger.LOG_LEVELS[normalizedLevel] ?? Logger.LOG_LEVELS.info;
  }

  public static getInstance(serviceName?: string, level?: string): Logger {
    if (!Logger.instance) {
      Logger.instance = new Logger(serviceName, level);
    }
    return Logger.instance;
  }

  public static resetInstance(): void {
    Logger.instance = null;
  }

  /** Create a child logger with additional default metadata */
  public child(defaultMeta: LogMeta): ChildLogger {
    return new ChildLogger(this, defaultMeta);
  }

  private formatLog(level: LogLevel, message: string, meta: LogMeta = {}): string {
    const entry: LogEntry = {
      timestamp: new Date().toISOString(),
      level,
      service: this.serviceName,
      message,
      ...meta,
    };
    return JSON.stringify(entry);
  }

  private shouldLog(level: LogLevel): boolean {
    return this.currentLevel <= Logger.LOG_LEVELS[level];
  }

  public debug(message: string, meta?: LogMeta): void {
    if (this.shouldLog('debug')) {
      console.log(this.formatLog('debug', message, meta));
    }
  }

  public info(message: string, meta?: LogMeta): void {
    if (this.shouldLog('info')) {
      console.log(this.formatLog('info', message, meta));
    }
  }

  public warn(message: string, meta?: LogMeta): void {
    if (this.shouldLog('warn')) {
      console.warn(this.formatLog('warn', message, meta));
    }
  }

  public error(message: string, meta?: LogMeta): void {
    if (this.shouldLog('error')) {
      console.error(this.formatLog('error', message, meta));
    }
  }
}

/** Child logger that carries default metadata with every log call */
export class ChildLogger {
  constructor(
    private readonly parent: Logger,
    private readonly defaultMeta: LogMeta
  ) {}

  public debug(message: string, meta?: LogMeta): void {
    this.parent.debug(message, { ...this.defaultMeta, ...meta });
  }

  public info(message: string, meta?: LogMeta): void {
    this.parent.info(message, { ...this.defaultMeta, ...meta });
  }

  public warn(message: string, meta?: LogMeta): void {
    this.parent.warn(message, { ...this.defaultMeta, ...meta });
  }

  public error(message: string, meta?: LogMeta): void {
    this.parent.error(message, { ...this.defaultMeta, ...meta });
  }
}

export const logger = Logger.getInstance();
export default logger;
