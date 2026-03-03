import { beforeAll, afterAll } from 'vitest';

beforeAll(async () => {
  process.env.AWS_REGION = 'ap-southeast-1';
  process.env.DYNAMODB_ENDPOINT = 'http://localhost:4566';
  process.env.DYNAMODB_TABLE_NAME = 'TowCommand-test';
  process.env.REDIS_HOST = 'localhost';
  process.env.REDIS_PORT = '6379';
  process.env.EVENT_BUS_NAME = 'towcommand-test';
  process.env.STAGE = 'test';
});

afterAll(async () => {
  // Cleanup test resources
});
