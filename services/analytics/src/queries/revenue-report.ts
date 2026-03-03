import { getDocClient, getTableName } from '@towcommand/db';
import { QueryCommand } from '@aws-sdk/lib-dynamodb';

export interface RevenueReport {
  period: string;
  totalBookings: number;
  completedBookings: number;
  cancelledBookings: number;
  totalRevenue: number;
  averageBookingValue: number;
  paymentsProcessed: number;
  paymentVolume: number;
}

export async function getRevenueReport(startDate: string, endDate: string): Promise<RevenueReport[]> {
  const client = getDocClient();
  const tableName = getTableName();

  // Query daily summary records between startDate and endDate
  const result = await client.send(new QueryCommand({
    TableName: tableName,
    KeyConditionExpression: 'PK BETWEEN :start AND :end AND SK = :sk',
    ExpressionAttributeValues: {
      ':start': `ANALYTICS#DAILY#${startDate}`,
      ':end': `ANALYTICS#DAILY#${endDate}`,
      ':sk': 'SUMMARY',
    },
  }));

  // DynamoDB BETWEEN on PK requires a Scan; use query per day instead
  const reports: RevenueReport[] = [];
  const start = new Date(startDate);
  const end = new Date(endDate);

  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    const dateKey = d.toISOString().split('T')[0];

    const dayResult = await client.send(new QueryCommand({
      TableName: tableName,
      KeyConditionExpression: 'PK = :pk AND SK = :sk',
      ExpressionAttributeValues: {
        ':pk': `ANALYTICS#DAILY#${dateKey}`,
        ':sk': 'SUMMARY',
      },
    }));

    const item = dayResult.Items?.[0];
    if (item) {
      const completed = (item.completedBookings as number) ?? 0;
      const revenue = (item.totalRevenue as number) ?? 0;
      reports.push({
        period: dateKey,
        totalBookings: (item.totalBookings as number) ?? 0,
        completedBookings: completed,
        cancelledBookings: (item.cancelledBookings as number) ?? 0,
        totalRevenue: revenue,
        averageBookingValue: completed > 0 ? Math.round(revenue / completed) : 0,
        paymentsProcessed: (item.paymentsProcessed as number) ?? 0,
        paymentVolume: (item.paymentVolume as number) ?? 0,
      });
    }
  }

  return reports;
}
