import { beforeAll, afterAll } from 'vitest';

let apiUrl: string;

beforeAll(async () => {
  apiUrl = process.env.API_URL ?? 'http://localhost:3000';
  process.env.API_URL = apiUrl;
});

afterAll(async () => {
  // Cleanup
});

export { apiUrl };
