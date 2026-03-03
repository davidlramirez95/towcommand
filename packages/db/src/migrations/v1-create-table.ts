import { CreateTableCommand, DynamoDBClient } from '@aws-sdk/client-dynamodb';

const TABLE_NAME = process.env.DYNAMODB_TABLE_NAME ?? 'TowCommand-dev';

export async function createTable() {
  const client = new DynamoDBClient({
    region: process.env.AWS_REGION ?? 'ap-southeast-1',
    ...(process.env.DYNAMODB_ENDPOINT ? { endpoint: process.env.DYNAMODB_ENDPOINT } : {}),
  });

  const command = new CreateTableCommand({
    TableName: TABLE_NAME,
    KeySchema: [
      { AttributeName: 'PK', KeyType: 'HASH' },
      { AttributeName: 'SK', KeyType: 'RANGE' },
    ],
    AttributeDefinitions: [
      { AttributeName: 'PK', AttributeType: 'S' },
      { AttributeName: 'SK', AttributeType: 'S' },
      { AttributeName: 'GSI1PK', AttributeType: 'S' },
      { AttributeName: 'GSI1SK', AttributeType: 'S' },
      { AttributeName: 'GSI2PK', AttributeType: 'S' },
      { AttributeName: 'GSI2SK', AttributeType: 'S' },
      { AttributeName: 'GSI3PK', AttributeType: 'S' },
      { AttributeName: 'GSI3SK', AttributeType: 'S' },
      { AttributeName: 'GSI4PK', AttributeType: 'S' },
      { AttributeName: 'GSI4SK', AttributeType: 'S' },
      { AttributeName: 'GSI5PK', AttributeType: 'S' },
      { AttributeName: 'GSI5SK', AttributeType: 'S' },
    ],
    GlobalSecondaryIndexes: [
      { IndexName: 'GSI1-UserJobs', KeySchema: [{ AttributeName: 'GSI1PK', KeyType: 'HASH' }, { AttributeName: 'GSI1SK', KeyType: 'RANGE' }], Projection: { ProjectionType: 'ALL' } },
      { IndexName: 'GSI2-StatusJobs', KeySchema: [{ AttributeName: 'GSI2PK', KeyType: 'HASH' }, { AttributeName: 'GSI2SK', KeyType: 'RANGE' }], Projection: { ProjectionType: 'ALL' } },
      { IndexName: 'GSI3-ProviderByTier', KeySchema: [{ AttributeName: 'GSI3PK', KeyType: 'HASH' }, { AttributeName: 'GSI3SK', KeyType: 'RANGE' }], Projection: { ProjectionType: 'ALL' } },
      { IndexName: 'GSI4-DisputeByStatus', KeySchema: [{ AttributeName: 'GSI4PK', KeyType: 'HASH' }, { AttributeName: 'GSI4SK', KeyType: 'RANGE' }], Projection: { ProjectionType: 'ALL' } },
      { IndexName: 'GSI5-PhoneIndex', KeySchema: [{ AttributeName: 'GSI5PK', KeyType: 'HASH' }, { AttributeName: 'GSI5SK', KeyType: 'RANGE' }], Projection: { ProjectionType: 'ALL' } },
    ],
    BillingMode: 'PAY_PER_REQUEST',
    StreamSpecification: {
      StreamEnabled: true,
      StreamViewType: 'NEW_AND_OLD_IMAGES',
    },
  });

  try {
    await client.send(command);
    console.log(`Table ${TABLE_NAME} created successfully`);
  } catch (error: unknown) {
    if ((error as { name: string }).name === 'ResourceInUseException') {
      console.log(`Table ${TABLE_NAME} already exists`);
    } else {
      throw error;
    }
  }
}

// Allow running directly
if (require.main === module) {
  createTable().catch(console.error);
}
