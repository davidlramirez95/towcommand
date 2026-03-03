import { KEY_PREFIXES, buildKey } from '../table-design';
import type { Provider, ProviderDoc } from '@towcommand/core';

export function providerKeys(providerId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.PROVIDER, providerId),
    SK: KEY_PREFIXES.PROFILE,
  };
}

export function providerGSI3Keys(trustTier: string, city: string, rating: number) {
  return {
    GSI3PK: `${KEY_PREFIXES.TIER}${trustTier}#${city}`,
    GSI3SK: String(rating).padStart(5, '0'),
  };
}

export function providerDocKeys(providerId: string, docType: string) {
  return {
    PK: buildKey(KEY_PREFIXES.PROVIDER, providerId),
    SK: buildKey(KEY_PREFIXES.DOC, docType),
  };
}

export function toProviderItem(provider: Provider, city = 'NCR') {
  return {
    ...providerKeys(provider.providerId),
    ...providerGSI3Keys(provider.trustTier, city, provider.rating),
    entityType: 'Provider',
    ...provider,
  };
}

export function toProviderDocItem(doc: ProviderDoc) {
  return {
    ...providerDocKeys(doc.providerId, doc.docType),
    entityType: 'ProviderDoc',
    ...doc,
  };
}
