import { useEffect, useRef, useCallback } from 'react';
import { WSClient } from '@/lib/ws';
import { useBookingStore } from '@/stores/booking';
import { useLocationStore } from '@/stores/location';
import { useAuthStore } from '@/stores/auth';

/**
 * Manages the WebSocket connection lifecycle and routes messages
 * to the appropriate Zustand stores.
 */
export function useWebSocket() {
  const wsRef = useRef<WSClient | null>(null);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const updateStatus = useBookingStore((s) => s.updateStatus);
  const updateProviderLocation = useBookingStore((s) => s.updateProviderLocation);
  const updateETA = useBookingStore((s) => s.updateETA);

  const handleMessage = useCallback(
    (type: string, data: Record<string, unknown>) => {
      switch (type) {
        case 'location_update':
          if (typeof data.lat === 'number' && typeof data.lng === 'number') {
            updateProviderLocation(data.lat, data.lng);
          }
          break;
        case 'booking_status':
          if (typeof data.status === 'string') {
            updateStatus(data.status as Parameters<typeof updateStatus>[0]);
          }
          break;
        case 'eta_update':
          if (typeof data.eta === 'number') {
            updateETA(data.eta);
          }
          break;
        case 'pong':
          // Heartbeat response, no action needed
          break;
        default:
          break;
      }
    },
    [updateStatus, updateProviderLocation, updateETA],
  );

  useEffect(() => {
    if (!isAuthenticated) return;

    const client = new WSClient({
      onMessage: handleMessage,
    });
    wsRef.current = client;
    client.connect();

    return () => {
      client.disconnect();
      wsRef.current = null;
    };
  }, [isAuthenticated, handleMessage]);

  const sendLocation = useCallback((lat: number, lng: number) => {
    wsRef.current?.send('location_update', { lat, lng });
  }, []);

  const sendChatMessage = useCallback((bookingId: string, message: string) => {
    wsRef.current?.send('chat_message', { bookingId, message });
  }, []);

  return {
    isConnected: wsRef.current?.isConnected ?? false,
    sendLocation,
    sendChatMessage,
  };
}
