import { BaseRepository } from './base.repo';
import { providerKeys, toProviderItem, providerDocKeys, toProviderDocItem } from '../entities/provider';
import type { Provider, ProviderDoc } from '@towcommand/core';

export class ProviderRepository extends BaseRepository {
  async getById(providerId: string): Promise<Provider | null> {
    const { PK, SK } = providerKeys(providerId);
    return this.getItem<Provider>(PK, SK);
  }

  async create(provider: Provider, city = 'NCR'): Promise<void> {
    await this.putItem(toProviderItem(provider, city));
  }

  async update(providerId: string, updates: Partial<Provider>): Promise<void> {
    const { PK, SK } = providerKeys(providerId);
    await this.updateItem(PK, SK, { ...updates, updatedAt: new Date().toISOString() });
  }

  async getByTierAndCity(tier: string, city: string, limit = 20): Promise<Provider[]> {
    return this.query<Provider>({
      IndexName: 'GSI3-ProviderByTier',
      KeyConditionExpression: 'GSI3PK = :pk',
      ExpressionAttributeValues: { ':pk': `TIER#${tier}#${city}` },
      ScanIndexForward: false,
      Limit: limit,
    });
  }

  async uploadDoc(doc: ProviderDoc): Promise<void> {
    await this.putItem(toProviderDocItem(doc));
  }

  async getDocs(providerId: string): Promise<ProviderDoc[]> {
    return this.query<ProviderDoc>({
      KeyConditionExpression: 'PK = :pk AND begins_with(SK, :sk)',
      ExpressionAttributeValues: { ':pk': `PROV#${providerId}`, ':sk': 'DOC#' },
    });
  }
}
