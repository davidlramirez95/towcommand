import { fetchAuthSession } from 'aws-amplify/auth';

const WS_BASE_URL = process.env.EXPO_PUBLIC_WS_URL ?? 'wss://ws.towcommand.ph';

type MessageHandler = (type: string, data: Record<string, unknown>) => void;

interface WebSocketClientOptions {
  onMessage: MessageHandler;
  onConnected?: () => void;
  onDisconnected?: () => void;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

/**
 * WebSocket client with auto-reconnect for API Gateway WebSocket.
 * Handles Cognito JWT authentication via query parameter.
 */
export class WSClient {
  private ws: WebSocket | null = null;
  private options: Required<WebSocketClientOptions>;
  private reconnectAttempts = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private isIntentionallyClosed = false;

  constructor(options: WebSocketClientOptions) {
    this.options = {
      reconnectInterval: 2000,
      maxReconnectAttempts: 10,
      onConnected: () => {},
      onDisconnected: () => {},
      ...options,
    };
  }

  async connect(): Promise<void> {
    this.isIntentionallyClosed = false;
    this.reconnectAttempts = 0;

    try {
      const session = await fetchAuthSession();
      const token = session.tokens?.idToken?.toString();
      const url = token ? `${WS_BASE_URL}?token=${token}` : WS_BASE_URL;

      this.ws = new WebSocket(url);
      this.setupListeners();
    } catch {
      this.scheduleReconnect();
    }
  }

  disconnect(): void {
    this.isIntentionallyClosed = true;
    this.cleanup();
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
  }

  send(action: string, data: Record<string, unknown>): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ action, ...data }));
    }
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  private setupListeners(): void {
    if (!this.ws) return;

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.startHeartbeat();
      this.options.onConnected();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data as string);
        const type = msg.type ?? msg.action ?? 'unknown';
        this.options.onMessage(type, msg);
      } catch {
        // Ignore malformed messages
      }
    };

    this.ws.onclose = () => {
      this.cleanup();
      this.options.onDisconnected();
      if (!this.isIntentionallyClosed) {
        this.scheduleReconnect();
      }
    };

    this.ws.onerror = () => {
      // onclose will fire after onerror
    };
  }

  private startHeartbeat(): void {
    this.heartbeatTimer = setInterval(() => {
      this.send('ping', {});
    }, 30_000);
  }

  private cleanup(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      return;
    }

    const delay = this.options.reconnectInterval * Math.pow(2, this.reconnectAttempts);
    this.reconnectAttempts++;

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, Math.min(delay, 30_000));
  }
}
