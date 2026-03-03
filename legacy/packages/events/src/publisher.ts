import { EventBridgeClient, PutEventsCommand } from '@aws-sdk/client-eventbridge';
import { ulid } from 'ulid';
import type { TowCommandEvent, EventSource, EventMetadata } from '@towcommand/core';

let client: EventBridgeClient | null = null;

function getClient(): EventBridgeClient {
  if (!client) {
    client = new EventBridgeClient({
      region: process.env.AWS_REGION ?? 'ap-southeast-1',
    });
  }
  return client;
}

function getEventBusName(): string {
  return process.env.EVENT_BUS_NAME ?? `towcommand-${process.env.STAGE ?? 'dev'}`;
}

export async function publishEvent<T>(
  source: EventSource,
  detailType: string,
  detail: T,
  actor?: { userId: string; userType: string },
  correlationId?: string,
): Promise<string> {
  const eventId = ulid();
  const metadata: EventMetadata = {
    eventId,
    correlationId: correlationId ?? eventId,
    timestamp: new Date().toISOString(),
    version: '1.0',
    actor,
  };

  const event: TowCommandEvent<T> = { source, detailType, detail, metadata };

  await getClient().send(
    new PutEventsCommand({
      Entries: [
        {
          EventBusName: getEventBusName(),
          Source: source,
          DetailType: detailType,
          Detail: JSON.stringify(event),
          Time: new Date(),
        },
      ],
    }),
  );

  return eventId;
}

export async function publishBatch(
  events: Array<{ source: EventSource; detailType: string; detail: unknown }>,
): Promise<void> {
  const entries = events.map((e) => ({
    EventBusName: getEventBusName(),
    Source: e.source,
    DetailType: e.detailType,
    Detail: JSON.stringify({
      source: e.source,
      detailType: e.detailType,
      detail: e.detail,
      metadata: {
        eventId: ulid(),
        correlationId: ulid(),
        timestamp: new Date().toISOString(),
        version: '1.0',
      },
    }),
    Time: new Date(),
  }));

  // EventBridge max 10 entries per call
  for (let i = 0; i < entries.length; i += 10) {
    await getClient().send(
      new PutEventsCommand({ Entries: entries.slice(i, i + 10) }),
    );
  }
}
