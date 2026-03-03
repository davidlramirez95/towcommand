// Notification event detail schemas for EventBridge
// Zod schemas can be added here as needed for event validation

export interface NotificationDetail {
  userId: string;
  type: 'sms' | 'email' | 'push' | 'in_app';
  title: string;
  message: string;
  metadata?: Record<string, unknown>;
  sentAt: string;
}

export interface EmailNotificationDetail extends NotificationDetail {
  type: 'email';
  email: string;
  subject: string;
}

export interface PushNotificationDetail extends NotificationDetail {
  type: 'push';
  deviceToken: string;
}
