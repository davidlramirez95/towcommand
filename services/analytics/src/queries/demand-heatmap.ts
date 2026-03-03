import { getDocClient, getTableName } from '@towcommand/db';
import { QueryCommand } from '@aws-sdk/lib-dynamodb';

export interface DemandHeatmapCell {
  gridCellId: string;
  hour: number;
  bookingCount: number;
  demandLevel: 'low' | 'medium' | 'high' | 'surge';
}

export async function getDemandHeatmap(date: string, hour?: number): Promise<DemandHeatmapCell[]> {
  const client = getDocClient();
  const tableName = getTableName();

  let keyCondition: string;
  const expressionValues: Record<string, unknown> = {
    ':pk': `ANALYTICS#HEATMAP#${date}`,
  };

  if (hour !== undefined) {
    const hourStr = String(hour).padStart(2, '0');
    keyCondition = 'PK = :pk AND begins_with(SK, :sk)';
    expressionValues[':sk'] = `HOUR#${hourStr}#CELL#`;
  } else {
    keyCondition = 'PK = :pk';
  }

  const result = await client.send(new QueryCommand({
    TableName: tableName,
    KeyConditionExpression: keyCondition,
    ExpressionAttributeValues: expressionValues,
  }));

  return (result.Items ?? []).map((item) => {
    const count = (item.bookingCount as number) ?? 0;
    return {
      gridCellId: (item.gridCellId as string) ?? '',
      hour: (item.hour as number) ?? 0,
      bookingCount: count,
      demandLevel: getDemandLevel(count),
    };
  });
}

function getDemandLevel(count: number): 'low' | 'medium' | 'high' | 'surge' {
  if (count >= 20) return 'surge';
  if (count >= 10) return 'high';
  if (count >= 5) return 'medium';
  return 'low';
}

export async function getDemandSummary(date: string): Promise<{
  totalBookings: number;
  peakHour: number;
  hotspots: Array<{ gridCellId: string; count: number }>;
}> {
  const cells = await getDemandHeatmap(date);

  const hourTotals: Record<number, number> = {};
  const cellTotals: Record<string, number> = {};

  for (const cell of cells) {
    hourTotals[cell.hour] = (hourTotals[cell.hour] ?? 0) + cell.bookingCount;
    cellTotals[cell.gridCellId] = (cellTotals[cell.gridCellId] ?? 0) + cell.bookingCount;
  }

  let peakHour = 0;
  let peakCount = 0;
  for (const [hour, count] of Object.entries(hourTotals)) {
    if (count > peakCount) {
      peakHour = parseInt(hour, 10);
      peakCount = count;
    }
  }

  const hotspots = Object.entries(cellTotals)
    .map(([gridCellId, count]) => ({ gridCellId, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10);

  const totalBookings = cells.reduce((sum, c) => sum + c.bookingCount, 0);

  return { totalBookings, peakHour, hotspots };
}
