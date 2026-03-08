import { create } from 'zustand';

interface Notification {
  id: string;
  title: string;
  body: string;
  data?: Record<string, string>;
  receivedAt: number;
  read: boolean;
}

interface NotificationState {
  notifications: Notification[];
  pushToken: string | null;
  unreadCount: number;

  setPushToken: (token: string | null) => void;
  addNotification: (notification: Notification) => void;
  markRead: (id: string) => void;
  markAllRead: () => void;
  reset: () => void;
}

export const useNotificationStore = create<NotificationState>()((set, get) => ({
  notifications: [],
  pushToken: null,
  unreadCount: 0,

  setPushToken: (pushToken) => set({ pushToken }),
  addNotification: (notification) =>
    set((state) => ({
      notifications: [notification, ...state.notifications].slice(0, 50),
      unreadCount: state.unreadCount + 1,
    })),
  markRead: (id) =>
    set((state) => ({
      notifications: state.notifications.map((n) =>
        n.id === id ? { ...n, read: true } : n,
      ),
      unreadCount: Math.max(0, state.unreadCount - (state.notifications.find((n) => n.id === id && !n.read) ? 1 : 0)),
    })),
  markAllRead: () =>
    set((state) => ({
      notifications: state.notifications.map((n) => ({ ...n, read: true })),
      unreadCount: 0,
    })),
  reset: () => set({ notifications: [], pushToken: null, unreadCount: 0 }),
}));
