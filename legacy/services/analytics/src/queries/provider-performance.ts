import { getDocClient, getTableName } from '@towcommand/db';
import { ProviderRepository } from '@towcommand/db';
import { QueryCommand, ScanCommand } from '@aws-sdk/lib-dynamodb';

export interface ProviderPerformance {
  providerId: string;
  name: string;
  rating: number;
  completedTrips: number;
  totalEarnings: number;
  acceptanceRate: number;
  lastUpdated: string;
}

const providerRepo = new ProviderRepository();

export async function getProviderPerformance(providerId: string): Promise<ProviderPerformance | null> {
  const client = getDocClient();
  const tableName = getTableName();

  const result = await client.send(new QueryCommand({
    TableName: tableName,
    KeyConditionExpression: 'PK = :pk AND SK = :sk',
    ExpressionAttributeValues: {
      ':pk': `ANALYTICS#PROVIDER#${providerId}`,
      ':sk': 'STATS',
    },
  }));

  const stats = result.Items?.[0];
  const provider = await providerRepo.getById(providerId);

  if (!stats && !provider) return null;

  return {
    providerId,
    name: provider?.name ?? 'Unknown',
    rating: provider?.rating ?? 0,
    completedTrips: (stats?.completedTrips as number) ?? provider?.totalJobsCompleted ?? 0,
    totalEarnings: (stats?.totalEarnings as number) ?? 0,
    acceptanceRate: provider?.acceptanceRate ?? 0,
    lastUpdated: (stats?.lastUpdated as string) ?? provider?.updatedAt ?? '',
  };
}

export async function getTopProviders(limit = 10): Promise<ProviderPerformance[]> {
  const client = getDocClient();
  const tableName = getTableName();

  // Scan analytics provider stats (small dataset for MVP)
  const result = await client.send(new ScanCommand({
    TableName: tableName,
    FilterExpression: 'entityType = :type',
    ExpressionAttributeValues: { ':type': 'ProviderAnalytics' },
    Limit: 100,
  }));

  const items = result.Items ?? [];

  // Enrich with provider details and sort
  const performances = await Promise.all(
    items.map(async (item) => {
      const providerId = (item.PK as string).replace('ANALYTICS#PROVIDER#', '');
      const provider = await providerRepo.getById(providerId);
      return {
        providerId,
        name: provider?.name ?? 'Unknown',
        rating: provider?.rating ?? 0,
        completedTrips: (item.completedTrips as number) ?? 0,
        totalEarnings: (item.totalEarnings as number) ?? 0,
        acceptanceRate: provider?.acceptanceRate ?? 0,
        lastUpdated: (item.lastUpdated as string) ?? '',
      };
    }),
  );

  return performances
    .sort((a, b) => b.completedTrips - a.completedTrips)
    .slice(0, limit);
}
