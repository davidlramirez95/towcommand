import type { EventBridgeEvent } from 'aws-lambda';
import { BookingRepository, ProviderRepository } from '@towcommand/db';
import { getDocClient, getTableName, KEY_PREFIXES, buildKey } from '@towcommand/db';
import { PutCommand, UpdateCommand } from '@aws-sdk/lib-dynamodb';

const bookingRepo = new BookingRepository();
const providerRepo = new ProviderRepository();

export async function handler(event: EventBridgeEvent<string, Record<string, unknown>>): Promise<void> {
  const eventType = event['detail-type'];
  const detail = event.detail;

  try {
    console.log(`Analytics processing: ${eventType}`);

    switch (eventType) {
      case 'BookingCompleted':
        await recordBookingCompleted(detail);
        break;
      case 'BookingCancelled':
        await recordBookingCancelled(detail);
        break;
      case 'PaymentCompleted':
        await recordPaymentCompleted(detail);
        break;
      case 'BookingCreated':
        await recordBookingCreated(detail);
        break;
    }
  } catch (error) {
    console.error(`Analytics handler error for ${eventType}:`, error);
    // Don't throw - analytics failures should not block event processing
  }
}

async function recordBookingCompleted(detail: Record<string, unknown>): Promise<void> {
  const bookingId = detail.bookingId as string;
  const providerId = detail.providerId as string;
  const now = new Date();
  const dateKey = now.toISOString().split('T')[0]; // YYYY-MM-DD

  const client = getDocClient();
  const tableName = getTableName();

  // Write analytics record for this booking
  await client.send(new PutCommand({
    TableName: tableName,
    Item: {
      PK: `ANALYTICS#DAILY#${dateKey}`,
      SK: `BOOKING#${bookingId}`,
      entityType: 'AnalyticsBooking',
      bookingId,
      providerId,
      customerId: detail.customerId,
      amount: (detail.price as any)?.total ?? 0,
      serviceType: detail.serviceType,
      completedAt: now.toISOString(),
      GSI1PK: `ANALYTICS#PROVIDER#${providerId}`,
      GSI1SK: `COMPLETED#${now.toISOString()}`,
    },
  }));

  // Update daily aggregation counter
  await client.send(new UpdateCommand({
    TableName: tableName,
    Key: { PK: `ANALYTICS#DAILY#${dateKey}`, SK: 'SUMMARY' },
    UpdateExpression: 'ADD completedBookings :one, totalRevenue :amount SET entityType = :type, #date = :date',
    ExpressionAttributeNames: { '#date': 'date' },
    ExpressionAttributeValues: {
      ':one': 1,
      ':amount': (detail.price as any)?.total ?? 0,
      ':type': 'DailyAnalytics',
      ':date': dateKey,
    },
  }));

  // Update provider stats
  if (providerId) {
    await client.send(new UpdateCommand({
      TableName: tableName,
      Key: { PK: `ANALYTICS#PROVIDER#${providerId}`, SK: 'STATS' },
      UpdateExpression: 'ADD completedTrips :one, totalEarnings :amount SET entityType = :type, lastUpdated = :now',
      ExpressionAttributeValues: {
        ':one': 1,
        ':amount': (detail.price as any)?.total ?? 0,
        ':type': 'ProviderAnalytics',
        ':now': now.toISOString(),
      },
    }));
  }
}

async function recordBookingCancelled(detail: Record<string, unknown>): Promise<void> {
  const now = new Date();
  const dateKey = now.toISOString().split('T')[0];

  const client = getDocClient();
  const tableName = getTableName();

  await client.send(new UpdateCommand({
    TableName: tableName,
    Key: { PK: `ANALYTICS#DAILY#${dateKey}`, SK: 'SUMMARY' },
    UpdateExpression: 'ADD cancelledBookings :one, cancellationFees :fee SET entityType = :type, #date = :date',
    ExpressionAttributeNames: { '#date': 'date' },
    ExpressionAttributeValues: {
      ':one': 1,
      ':fee': (detail.cancellationFee as number) ?? 0,
      ':type': 'DailyAnalytics',
      ':date': dateKey,
    },
  }));
}

async function recordPaymentCompleted(detail: Record<string, unknown>): Promise<void> {
  const now = new Date();
  const dateKey = now.toISOString().split('T')[0];

  const client = getDocClient();
  const tableName = getTableName();

  await client.send(new UpdateCommand({
    TableName: tableName,
    Key: { PK: `ANALYTICS#DAILY#${dateKey}`, SK: 'SUMMARY' },
    UpdateExpression: 'ADD paymentsProcessed :one, paymentVolume :amount SET entityType = :type, #date = :date',
    ExpressionAttributeNames: { '#date': 'date' },
    ExpressionAttributeValues: {
      ':one': 1,
      ':amount': (detail.amount as number) ?? 0,
      ':type': 'DailyAnalytics',
      ':date': dateKey,
    },
  }));
}

async function recordBookingCreated(detail: Record<string, unknown>): Promise<void> {
  const now = new Date();
  const dateKey = now.toISOString().split('T')[0];
  const hour = now.getUTCHours() + 8; // PHT
  const phtHour = hour >= 24 ? hour - 24 : hour;
  const pickup = detail.pickupLocation as { lat: number; lng: number } | undefined;

  const client = getDocClient();
  const tableName = getTableName();

  // Update daily total
  await client.send(new UpdateCommand({
    TableName: tableName,
    Key: { PK: `ANALYTICS#DAILY#${dateKey}`, SK: 'SUMMARY' },
    UpdateExpression: 'ADD totalBookings :one SET entityType = :type, #date = :date',
    ExpressionAttributeNames: { '#date': 'date' },
    ExpressionAttributeValues: {
      ':one': 1,
      ':type': 'DailyAnalytics',
      ':date': dateKey,
    },
  }));

  // Update hourly demand for heatmap
  if (pickup) {
    // Grid cell: round to 2 decimal places (~1km resolution)
    const gridLat = Math.round(pickup.lat * 100) / 100;
    const gridLng = Math.round(pickup.lng * 100) / 100;
    const gridCellId = `${gridLat},${gridLng}`;

    await client.send(new UpdateCommand({
      TableName: tableName,
      Key: {
        PK: `ANALYTICS#HEATMAP#${dateKey}`,
        SK: `HOUR#${String(phtHour).padStart(2, '0')}#CELL#${gridCellId}`,
      },
      UpdateExpression: 'ADD bookingCount :one SET entityType = :type, gridCellId = :cell, #hour = :hour',
      ExpressionAttributeNames: { '#hour': 'hour' },
      ExpressionAttributeValues: {
        ':one': 1,
        ':type': 'HeatmapCell',
        ':cell': gridCellId,
        ':hour': phtHour,
      },
    }));
  }
}
