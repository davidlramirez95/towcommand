export const TABLE_CONFIG = {
  tableName: process.env.DYNAMODB_TABLE_NAME ?? 'TowCommand-dev',
  partitionKey: 'PK',
  sortKey: 'SK',
  gsis: {
    GSI1: { pk: 'GSI1PK', sk: 'GSI1SK', name: 'GSI1-UserJobs' },
    GSI2: { pk: 'GSI2PK', sk: 'GSI2SK', name: 'GSI2-StatusJobs' },
    GSI3: { pk: 'GSI3PK', sk: 'GSI3SK', name: 'GSI3-ProviderByTier' },
    GSI4: { pk: 'GSI4PK', sk: 'GSI4SK', name: 'GSI4-DisputeByStatus' },
    GSI5: { pk: 'GSI5PK', sk: 'GSI5SK', name: 'GSI5-PhoneIndex' },
  },
} as const;

export const KEY_PREFIXES = {
  USER: 'USER#',
  PROVIDER: 'PROV#',
  BOOKING: 'BOOK#',
  JOB: 'JOB#',
  VEHICLE: 'VEH#',
  PAYMENT: 'PAY#',
  TRANSACTION: 'TXN#',
  OTP: 'OTP#',
  SOS: 'SOS#',
  CHAT: 'CHAT#',
  MESSAGE: 'MSG#',
  RATING: 'RATING',
  DISPUTE: 'DISPUTE#',
  AUDIT: 'AUDIT#',
  REGION: 'REGION#',
  STATUS: 'STATUS#',
  TIER: 'TIER#',
  PHONE: 'PHONE#',
  EMAIL: 'EMAIL#',
  EVIDENCE: 'EVIDENCE#',
  MEDIA: 'MEDIA#',
  DOC: 'DOC#',
  SUKI: 'SUKI',
  PROFILE: 'PROFILE',
  DETAILS: 'DETAILS',
  AREA: 'AREA',
} as const;

export function buildKey(prefix: string, id: string): string {
  return `${prefix}${id}`;
}
