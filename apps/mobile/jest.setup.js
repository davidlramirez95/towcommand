/* eslint-disable no-undef */

// Silence console.warn for noisy RN warnings in tests
const originalWarn = console.warn;
console.warn = (...args) => {
  if (typeof args[0] === 'string' && args[0].includes('Animated')) return;
  originalWarn(...args);
};

// Mock global fetch
global.fetch = jest.fn();

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  onopen = null;
  onclose = null;
  onmessage = null;
  onerror = null;

  constructor(url) {
    this.url = url;
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      if (this.onopen) this.onopen({});
    }, 0);
  }

  send = jest.fn();
  close = jest.fn((code, reason) => {
    this.readyState = MockWebSocket.CLOSED;
    if (this.onclose) this.onclose({ code, reason });
  });
}
global.WebSocket = MockWebSocket;
