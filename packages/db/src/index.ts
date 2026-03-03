// Re-export client
export { getDocClient, getTableName } from './client';

// Re-export table design
export { TABLE_CONFIG, KEY_PREFIXES, buildKey } from './table-design';

// Re-export entities
export * from './entities';

// Re-export repositories
export * from './repositories';

// Re-export migrations
export { createTable } from './migrations/v1-create-table';
