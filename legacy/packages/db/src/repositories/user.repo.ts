import { BaseRepository } from './base.repo';
import { userKeys, vehicleKeys, toUserItem, toVehicleItem } from '../entities/user';
import type { User, UserVehicle } from '@towcommand/core';

export class UserRepository extends BaseRepository {
  async getById(userId: string): Promise<User | null> {
    const { PK, SK } = userKeys(userId);
    return this.getItem<User>(PK, SK);
  }

  async create(user: User): Promise<void> {
    await this.putItem(toUserItem(user));
  }

  async update(userId: string, updates: Partial<User>): Promise<void> {
    const { PK, SK } = userKeys(userId);
    await this.updateItem(PK, SK, { ...updates, updatedAt: new Date().toISOString() });
  }

  async getVehicles(userId: string): Promise<UserVehicle[]> {
    return this.query<UserVehicle>({
      KeyConditionExpression: 'PK = :pk AND begins_with(SK, :sk)',
      ExpressionAttributeValues: { ':pk': `USER#${userId}`, ':sk': 'VEH#' },
    });
  }

  async addVehicle(vehicle: UserVehicle): Promise<void> {
    await this.putItem(toVehicleItem(vehicle));
  }

  async getByPhone(phone: string): Promise<User | null> {
    const items = await this.query<User>({
      IndexName: 'GSI5-PhoneIndex',
      KeyConditionExpression: 'GSI5PK = :pk AND GSI5SK = :sk',
      ExpressionAttributeValues: { ':pk': `PHONE#${phone}`, ':sk': 'PROFILE' },
    });
    return items[0] ?? null;
  }
}
