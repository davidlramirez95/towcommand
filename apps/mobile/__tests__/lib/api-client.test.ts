/**
 * API Client Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - getAuthHeader() catches ALL errors silently. If Amplify throws
 *   "token expired," the client sends an unauthenticated request.
 *   The server returns 401, user sees "unauthorized" instead of
 *   "please log in again." Test both paths.
 * - 204 returns `undefined as T`. If caller does `result.id`,
 *   that's a runtime crash. Test that callers handle it.
 * - Network errors (fetch rejects) need to become APIError instances
 *   for consistent error handling downstream.
 */
import { api, APIError } from '@/lib/api/client';
import { fetchAuthSession } from 'aws-amplify/auth';

const mockFetch = global.fetch as jest.Mock;

beforeEach(() => {
  jest.clearAllMocks();
});

describe('api client', () => {
  describe('authenticated requests', () => {
    it('attaches Bearer token from Cognito session', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ data: 'test' }),
      });

      await api.get('/test');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/test'),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer mock-jwt-token-123',
          }),
        }),
      );
    });

    it('sends request WITHOUT token when auth session fails (graceful degradation)', async () => {
      (fetchAuthSession as jest.Mock).mockRejectedValueOnce(new Error('Token expired'));

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ data: 'public' }),
      });

      const result = await api.get('/public-endpoint');

      // Should still make the request, just without auth header
      expect(mockFetch).toHaveBeenCalled();
      const headers = mockFetch.mock.calls[0][1].headers;
      expect(headers.Authorization).toBeUndefined();
      expect(result).toEqual({ data: 'public' });
    });

    it('skips auth when authenticated=false (webhook-style endpoints)', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({}),
      });

      await api.get('/health', { authenticated: false });

      expect(fetchAuthSession).not.toHaveBeenCalled();
    });
  });

  describe('HTTP methods', () => {
    beforeEach(() => {
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true }),
      });
    });

    it('GET sends no body', async () => {
      await api.get('/bookings');
      expect(mockFetch.mock.calls[0][1].method).toBe('GET');
      expect(mockFetch.mock.calls[0][1].body).toBeUndefined();
    });

    it('POST sends JSON body', async () => {
      await api.post('/bookings', { serviceType: 'FLATBED' });
      expect(mockFetch.mock.calls[0][1].method).toBe('POST');
      expect(mockFetch.mock.calls[0][1].body).toBe('{"serviceType":"FLATBED"}');
    });

    it('PUT sends JSON body', async () => {
      await api.put('/bookings/1', { status: 'CANCELLED' });
      expect(mockFetch.mock.calls[0][1].method).toBe('PUT');
    });

    it('PATCH sends JSON body', async () => {
      await api.patch('/users/1', { phone: '+639170000000' });
      expect(mockFetch.mock.calls[0][1].method).toBe('PATCH');
    });

    it('DELETE sends no body', async () => {
      await api.delete('/bookings/1');
      expect(mockFetch.mock.calls[0][1].method).toBe('DELETE');
      expect(mockFetch.mock.calls[0][1].body).toBeUndefined();
    });
  });

  describe('error handling', () => {
    it('throws APIError with structured error from server', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: () =>
          Promise.resolve({
            error: { code: 'NOT_FOUND', message: 'Booking not found' },
          }),
      });

      const error = await api.get('/bookings/999').catch((e: unknown) => e as APIError);
      expect(error).toBeInstanceOf(APIError);
      expect((error as APIError).statusCode).toBe(404);
      expect((error as APIError).code).toBe('NOT_FOUND');
      expect((error as APIError).message).toBe('Booking not found');
    });

    it('handles non-JSON error body (502 nginx HTML page)', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 502,
        statusText: 'Bad Gateway',
        json: () => Promise.reject(new Error('not JSON')),
      });

      await expect(api.get('/any')).rejects.toThrow(APIError);
    });

    it('204 response returns undefined (callers must handle)', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      const result = await api.delete('/bookings/1');
      expect(result).toBeUndefined();
    });
  });

  describe('request configuration', () => {
    it('always sets Content-Type: application/json', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({}),
      });

      await api.get('/test');

      expect(mockFetch.mock.calls[0][1].headers['Content-Type']).toBe('application/json');
    });

    it('custom headers merge with defaults', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({}),
      });

      await api.get('/test', { headers: { 'X-Custom': 'value' } });

      const headers = mockFetch.mock.calls[0][1].headers;
      expect(headers['Content-Type']).toBe('application/json');
      expect(headers['X-Custom']).toBe('value');
    });
  });
});
