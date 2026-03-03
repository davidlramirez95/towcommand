// Provider event detail schemas for EventBridge
// Zod schemas can be added here as needed for event validation

export interface ProviderOnlineDetail {
  providerId: string;
  location: { lat: number; lng: number };
  onlineAt: string;
}

export interface ProviderOfflineDetail {
  providerId: string;
  offlineAt: string;
}

export interface ProviderVerifiedDetail {
  providerId: string;
  verifiedAt: string;
  trustTier: string;
}
