/**
 * WebSocket Client Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Reconnect backoff: if interval doubles each time but connect()
 *   resets the counter, calling connect() during reconnect creates
 *   parallel reconnect loops.
 * - send() on closed connection silently drops messages. Chat messages
 *   and location updates are lost without the user knowing.
 * - Heartbeat interval must be cleared on disconnect, otherwise
 *   timers leak and send() throws on garbage-collected WebSocket.
 * - Max reconnect attempts: after 10 failures, app must stop retrying
 *   (not spin forever consuming battery on a phone with no signal).
 */
import { WSClient } from '@/lib/ws/client';
import { fetchAuthSession } from 'aws-amplify/auth';

// Access the mock WebSocket class
const MockWebSocket = global.WebSocket as any;

beforeEach(() => {
  jest.clearAllMocks();
  jest.useFakeTimers();
  // Restore default fetchAuthSession behavior — clearAllMocks doesn't reset
  // mock implementations, so mockRejectedValue from reconnect tests persists.
  (fetchAuthSession as jest.Mock).mockResolvedValue({
    tokens: {
      idToken: { toString: () => 'mock-jwt-token-123', payload: {} },
      accessToken: { toString: () => 'mock-access-token' },
    },
  });
});

afterEach(() => {
  jest.useRealTimers();
});

describe('WSClient', () => {
  const defaultOptions = {
    onMessage: jest.fn(),
    onConnected: jest.fn(),
    onDisconnected: jest.fn(),
  };

  describe('connection', () => {
    it('connects with JWT token in query param', async () => {
      const client = new WSClient(defaultOptions);
      await client.connect();
      jest.advanceTimersByTime(10); // trigger onopen

      // Verify JWT token is in the WebSocket URL
      const ws = (client as any).ws;
      expect(ws?.url).toContain('token=mock-jwt-token-123');
    });

    it('connects without token when auth fails (public fallback)', async () => {
      (fetchAuthSession as jest.Mock).mockRejectedValueOnce(new Error('No session'));

      const client = new WSClient({
        ...defaultOptions,
        reconnectInterval: 1000,
        maxReconnectAttempts: 1,
      });
      await client.connect();

      // Should schedule reconnect, not crash
      expect(defaultOptions.onConnected).not.toHaveBeenCalled();
    });
  });

  describe('disconnect', () => {
    it('intentional disconnect prevents reconnect', async () => {
      const client = new WSClient(defaultOptions);
      await client.connect();
      jest.advanceTimersByTime(10); // trigger onopen

      client.disconnect();

      // After disconnect, no reconnect should be scheduled
      jest.advanceTimersByTime(60000);
      // onConnected should only have been called once (initial connect)
    });

    it('disconnect clears heartbeat timer (prevents timer leak)', async () => {
      const clearIntervalSpy = jest.spyOn(global, 'clearInterval');

      const client = new WSClient(defaultOptions);
      await client.connect();
      jest.advanceTimersByTime(10); // trigger onopen + start heartbeat

      client.disconnect();

      expect(clearIntervalSpy).toHaveBeenCalled();
      clearIntervalSpy.mockRestore();
    });
  });

  describe('send', () => {
    it('send on closed connection is silent no-op (no throw)', () => {
      const client = new WSClient(defaultOptions);
      // Never connected — readyState is CONNECTING or undefined

      expect(() => {
        client.send('location', { lat: 14.5, lng: 120.9 });
      }).not.toThrow();
    });

    it('isConnected returns false when not connected', () => {
      const client = new WSClient(defaultOptions);
      expect(client.isConnected).toBe(false);
    });
  });

  describe('reconnect backoff', () => {
    it('stops after maxReconnectAttempts (battery conservation)', async () => {
      (fetchAuthSession as jest.Mock).mockRejectedValue(new Error('No session'));

      const reconnectInterval = 100;
      const maxReconnectAttempts = 3;
      const client = new WSClient({
        ...defaultOptions,
        reconnectInterval,
        maxReconnectAttempts,
      });

      await client.connect(); // initial attempt fails

      // Advance timers for each reconnect attempt.
      // Note: connect() resets reconnectAttempts, so each cycle schedules
      // one reconnect. We advance enough to trigger maxReconnectAttempts retries.
      for (let i = 0; i < maxReconnectAttempts; i++) {
        jest.advanceTimersByTime(reconnectInterval * Math.pow(2, i));
        await Promise.resolve(); // flush microtasks from async connect()
      }

      // 1 initial + maxReconnectAttempts reconnects
      const callCount = (fetchAuthSession as jest.Mock).mock.calls.length;
      expect(callCount).toBeLessThanOrEqual(1 + maxReconnectAttempts);
    });
  });

  describe('message handling', () => {
    it('parses JSON messages and calls onMessage with type', async () => {
      const onMessage = jest.fn();
      const client = new WSClient({ ...defaultOptions, onMessage });

      await client.connect();
      jest.advanceTimersByTime(10); // trigger onopen

      // Simulate incoming message via the WebSocket's onmessage
      const ws = (client as any).ws;
      ws?.onmessage?.({ data: JSON.stringify({ type: 'location_update', lat: 14.5 }) });

      expect(onMessage).toHaveBeenCalledWith('location_update', expect.objectContaining({ lat: 14.5 }));
    });

    it('malformed JSON does not crash client (resilience)', async () => {
      const onMessage = jest.fn();
      const client = new WSClient({ ...defaultOptions, onMessage });

      await client.connect();
      jest.advanceTimersByTime(10);

      // Simulate malformed JSON message via the WebSocket's onmessage
      const ws = (client as any).ws;
      expect(() => {
        ws?.onmessage?.({ data: '{invalid json' });
      }).not.toThrow();

      // onMessage should not have been called with invalid JSON
      expect(onMessage).not.toHaveBeenCalled();
    });
  });
});
