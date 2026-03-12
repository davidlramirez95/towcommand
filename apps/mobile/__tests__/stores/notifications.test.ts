/**
 * Notifications Store Tests — 2nd Order Logic
 *
 * 2nd order concern: Notifications accumulate over time. Without a cap,
 * memory grows unbounded on long sessions. The store caps at 50, but
 * does markRead work on nonexistent IDs? Does markAllRead handle 100
 * items? Does addNotification correctly increment unreadCount?
 */
import { useNotificationStore } from '@/stores/notifications';

beforeEach(() => {
  useNotificationStore.getState().reset();
});

describe('notifications store', () => {
  const makeNotification = (id: string) => ({
    id,
    title: `Notification ${id}`,
    body: `Body for ${id}`,
    receivedAt: Date.now(),
    read: false,
  });

  it('initial state has empty notifications and zero unread', () => {
    const state = useNotificationStore.getState();
    expect(state.notifications).toEqual([]);
    expect(state.unreadCount).toBe(0);
    expect(state.pushToken).toBeNull();
  });

  it('addNotification prepends and increments unreadCount', () => {
    useNotificationStore.getState().addNotification(makeNotification('n1'));
    useNotificationStore.getState().addNotification(makeNotification('n2'));

    const state = useNotificationStore.getState();
    expect(state.notifications).toHaveLength(2);
    expect(state.notifications[0].id).toBe('n2'); // most recent first
    expect(state.unreadCount).toBe(2);
  });

  it('addNotification caps at 50 (memory leak prevention)', () => {
    for (let i = 0; i < 60; i++) {
      useNotificationStore.getState().addNotification(makeNotification(`n${i}`));
    }

    const state = useNotificationStore.getState();
    expect(state.notifications).toHaveLength(50);
    // Most recent should be first
    expect(state.notifications[0].id).toBe('n59');
  });

  it('markRead sets read=true and decrements unreadCount', () => {
    useNotificationStore.getState().addNotification(makeNotification('n1'));
    useNotificationStore.getState().addNotification(makeNotification('n2'));
    expect(useNotificationStore.getState().unreadCount).toBe(2);

    useNotificationStore.getState().markRead('n1');

    const state = useNotificationStore.getState();
    const n1 = state.notifications.find((n) => n.id === 'n1');
    expect(n1?.read).toBe(true);
    expect(state.unreadCount).toBe(1);
  });

  it('markRead on nonexistent ID is a no-op (no crash, no count change)', () => {
    useNotificationStore.getState().addNotification(makeNotification('n1'));

    expect(() => {
      useNotificationStore.getState().markRead('nonexistent');
    }).not.toThrow();

    expect(useNotificationStore.getState().unreadCount).toBe(1);
  });

  it('markRead on already-read notification doesnt double-decrement', () => {
    useNotificationStore.getState().addNotification(makeNotification('n1'));
    useNotificationStore.getState().markRead('n1');
    expect(useNotificationStore.getState().unreadCount).toBe(0);

    // Mark same one again
    useNotificationStore.getState().markRead('n1');
    expect(useNotificationStore.getState().unreadCount).toBe(0); // not -1
  });

  it('markAllRead resets all notifications and unread count', () => {
    for (let i = 0; i < 10; i++) {
      useNotificationStore.getState().addNotification(makeNotification(`n${i}`));
    }
    expect(useNotificationStore.getState().unreadCount).toBe(10);

    useNotificationStore.getState().markAllRead();

    const state = useNotificationStore.getState();
    expect(state.unreadCount).toBe(0);
    expect(state.notifications.every((n) => n.read)).toBe(true);
    expect(state.notifications).toHaveLength(10); // items still exist, just read
  });

  it('setPushToken stores the FCM/APNs token', () => {
    useNotificationStore.getState().setPushToken('ExponentPushToken[abc123]');
    expect(useNotificationStore.getState().pushToken).toBe('ExponentPushToken[abc123]');
  });

  it('reset clears everything including pushToken', () => {
    useNotificationStore.getState().setPushToken('token-123');
    useNotificationStore.getState().addNotification(makeNotification('n1'));

    useNotificationStore.getState().reset();

    const state = useNotificationStore.getState();
    expect(state.notifications).toEqual([]);
    expect(state.unreadCount).toBe(0);
    expect(state.pushToken).toBeNull();
  });
});
