import { GetCommand, PutCommand, UpdateCommand, DeleteCommand, QueryCommand, type QueryCommandInput } from '@aws-sdk/lib-dynamodb';
import { getDocClient, getTableName } from '../client';

export abstract class BaseRepository {
  protected get client() {
    return getDocClient();
  }

  protected get tableName() {
    return getTableName();
  }

  protected async getItem<T>(pk: string, sk: string): Promise<T | null> {
    const result = await this.client.send(
      new GetCommand({ TableName: this.tableName, Key: { PK: pk, SK: sk } }),
    );
    return (result.Item as T) ?? null;
  }

  protected async putItem(item: Record<string, unknown>): Promise<void> {
    await this.client.send(
      new PutCommand({ TableName: this.tableName, Item: item }),
    );
  }

  protected async updateItem(
    pk: string, sk: string,
    updates: Record<string, unknown>,
  ): Promise<void> {
    const expressions: string[] = [];
    const names: Record<string, string> = {};
    const values: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(updates)) {
      const attrName = `#${key}`;
      const attrValue = `:${key}`;
      expressions.push(`${attrName} = ${attrValue}`);
      names[attrName] = key;
      values[attrValue] = value;
    }

    await this.client.send(
      new UpdateCommand({
        TableName: this.tableName,
        Key: { PK: pk, SK: sk },
        UpdateExpression: `SET ${expressions.join(', ')}`,
        ExpressionAttributeNames: names,
        ExpressionAttributeValues: values,
      }),
    );
  }

  protected async deleteItem(pk: string, sk: string): Promise<void> {
    await this.client.send(
      new DeleteCommand({ TableName: this.tableName, Key: { PK: pk, SK: sk } }),
    );
  }

  protected async query<T>(params: Omit<QueryCommandInput, 'TableName'>): Promise<T[]> {
    const result = await this.client.send(
      new QueryCommand({ TableName: this.tableName, ...params }),
    );
    return (result.Items as T[]) ?? [];
  }

  protected async queryWithPagination<T>(
    params: Omit<QueryCommandInput, 'TableName'>,
    limit = 25,
  ): Promise<{ items: T[]; lastKey?: Record<string, unknown> }> {
    const result = await this.client.send(
      new QueryCommand({ TableName: this.tableName, ...params, Limit: limit }),
    );
    return {
      items: (result.Items as T[]) ?? [],
      lastKey: result.LastEvaluatedKey as Record<string, unknown> | undefined,
    };
  }
}
