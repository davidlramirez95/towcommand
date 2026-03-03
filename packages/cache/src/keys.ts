export const CACHE_KEYS = {
  providerLocation: (id: string) => `provider:location:${id}`,
  providerAvailable: (city: string) => `provider:available:${city}`,
  providerGeo: (city: string) => `provider:geo:${city}`,
  jobLock: (jobId: string) => `job:lock:${jobId}`,
  otp: (jobId: string, type: string) => `otp:${jobId}:${type}`,
  userClaims: (sub: string) => `user:claims:${sub}`,
  estimate: (hash: string) => `estimate:${hash}`,
  rateLimit: (userId: string) => `rate:user:${userId}`,
  wsConnection: (userId: string) => `ws:connection:${userId}`,
  surgeMultiplier: (region: string) => `surge:${region}`,
} as const;

export const CACHE_TTL = {
  providerLocation: 60,
  otp: 900,
  userClaims: 300,
  estimate: 600,
  rateLimit: 60,
  jobLock: 30,
} as const;
