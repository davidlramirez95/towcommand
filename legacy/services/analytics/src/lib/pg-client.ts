// TODO: Uncomment when PostgreSQL/Aurora analytics is ready and budget allows
// This module provides PostgreSQL connectivity for analytics queries.
// Currently disabled as part of serverless-first approach to avoid AWS charges.

// import { Pool, PoolConfig, QueryResult } from 'pg';
//
// const poolConfig: PoolConfig = {
//   host: process.env.PG_HOST,
//   port: parseInt(process.env.PG_PORT || '5432'),
//   database: process.env.PG_DATABASE || 'towcommand_analytics',
//   user: process.env.PG_USER,
//   password: process.env.PG_PASSWORD,
//   max: parseInt(process.env.PG_MAX_CONNECTIONS || '10'),
//   idleTimeoutMillis: 30000,
//   connectionTimeoutMillis: 5000,
//   ssl: process.env.PG_SSL === 'true' ? { rejectUnauthorized: false } : false,
// };
//
// let pool: Pool | null = null;
//
// export function getPool(): Pool {
//   if (!pool) {
//     pool = new Pool(poolConfig);
//     pool.on('error', (err) => {
//       console.error('Unexpected pool error:', err);
//     });
//   }
//   return pool;
// }
//
// export async function closePool(): Promise<void> {
//   if (pool) {
//     await pool.end();
//     pool = null;
//   }
// }
//
// export async function query<T = unknown>(text: string, values?: unknown[]): Promise<T[]> {
//   const client = getPool();
//   const result: QueryResult<T> = await client.query(text, values);
//   return result.rows;
// }

// Active serverless stubs (DynamoDB/Athena approach)
export interface DatabaseConfig {
  // Placeholder for future serverless database implementation
}

export function getPool(): never {
  throw new Error('PostgreSQL not enabled. Analytics will use DynamoDB or Athena. Uncomment pg-client code when Aurora is provisioned.');
}

export async function closePool(): Promise<void> {
  // No-op: serverless databases do not require connection pools
}

export async function query<T = unknown>(text: string, values?: unknown[]): Promise<T[]> {
  throw new Error('PostgreSQL not enabled. Use DynamoDB queries or Athena instead. Uncomment pg-client code when Aurora is provisioned.');
}
