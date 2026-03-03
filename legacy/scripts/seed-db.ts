import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, PutCommand } from '@aws-sdk/lib-dynamodb';

const client = DynamoDBDocumentClient.from(
  new DynamoDBClient({
    region: process.env.AWS_REGION ?? 'ap-southeast-1',
    ...(process.env.DYNAMODB_ENDPOINT ? { endpoint: process.env.DYNAMODB_ENDPOINT } : {}),
  }),
);

const TABLE = process.env.DYNAMODB_TABLE_NAME ?? 'TowCommand-dev';

const seedData = [
  // Test customer
  {
    PK: 'USER#user-001', SK: 'PROFILE',
    GSI1PK: 'EMAIL#juan@example.com', GSI1SK: 'USER',
    GSI5PK: 'PHONE#+639171234567', GSI5SK: 'PROFILE',
    entityType: 'User',
    userId: 'user-001', name: 'Juan dela Cruz', email: 'juan@example.com',
    phone: '+639171234567', userType: 'customer', trustTier: 'basic',
    language: 'fil', status: 'active',
    createdAt: new Date().toISOString(), updatedAt: new Date().toISOString(),
  },
  // Test vehicle
  {
    PK: 'USER#user-001', SK: 'VEH#veh-001',
    entityType: 'UserVehicle',
    vehicleId: 'veh-001', userId: 'user-001',
    make: 'Toyota', model: 'Vios', year: 2020,
    plateNumber: 'ABC-1234', weightClass: 'light', color: 'white', isDefault: true,
    createdAt: new Date().toISOString(),
  },
  // Test provider
  {
    PK: 'PROV#prov-001', SK: 'PROFILE',
    GSI3PK: 'TIER#verified#NCR', GSI3SK: '00450',
    entityType: 'Provider',
    providerId: 'prov-001', name: 'Pedro Santos Towing',
    phone: '+639181234567', email: 'pedro@towing.ph',
    status: 'active', trustTier: 'verified', truckType: 'flatbed',
    maxWeightCapacityKg: 4500, plateNumber: 'XYZ-5678',
    rating: 4.5, totalJobsCompleted: 127, acceptanceRate: 0.92,
    isOnline: true, currentLat: 14.5547, currentLng: 121.0244,
    serviceAreas: ['NCR'],
    createdAt: new Date().toISOString(), updatedAt: new Date().toISOString(),
  },
  // Test provider 2
  {
    PK: 'PROV#prov-002', SK: 'PROFILE',
    GSI3PK: 'TIER#basic#NCR', GSI3SK: '00380',
    entityType: 'Provider',
    providerId: 'prov-002', name: 'Maria Garcia Rescue',
    phone: '+639191234567', email: 'maria@rescue.ph',
    status: 'active', trustTier: 'basic', truckType: 'wheel_lift',
    maxWeightCapacityKg: 3000, plateNumber: 'DEF-9012',
    rating: 3.8, totalJobsCompleted: 45, acceptanceRate: 0.85,
    isOnline: true, currentLat: 14.5896, currentLng: 121.0606,
    serviceAreas: ['NCR'],
    createdAt: new Date().toISOString(), updatedAt: new Date().toISOString(),
  },
  // Service area
  {
    PK: 'REGION#NCR', SK: 'AREA',
    entityType: 'ServiceArea',
    code: 'NCR', name: 'National Capital Region',
    defaultRadiusKm: 10, surgeMultiplierCap: 1.5, isActive: true,
  },
];

async function seed() {
  console.log('🌱 Seeding TowCommand database...');

  for (const item of seedData) {
    await client.send(new PutCommand({ TableName: TABLE, Item: item }));
    console.log(`  ✅ ${item.entityType}: ${item.PK} | ${item.SK}`);
  }

  console.log(`\n✅ Seeded ${seedData.length} items to ${TABLE}`);
}

seed().catch(console.error);
