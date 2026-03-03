import { KEY_PREFIXES, buildKey } from '../table-design';
import type { User, UserVehicle } from '@towcommand/core';

export function userKeys(userId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.USER, userId),
    SK: KEY_PREFIXES.PROFILE,
  };
}

export function userGSI1Keys(email: string) {
  return {
    GSI1PK: buildKey(KEY_PREFIXES.EMAIL, email),
    GSI1SK: 'USER',
  };
}

export function userPhoneGSI5Keys(phone: string) {
  return {
    GSI5PK: buildKey(KEY_PREFIXES.PHONE, phone),
    GSI5SK: KEY_PREFIXES.PROFILE,
  };
}

export function vehicleKeys(userId: string, vehicleId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.USER, userId),
    SK: buildKey(KEY_PREFIXES.VEHICLE, vehicleId),
  };
}

export function toUserItem(user: User) {
  return {
    ...userKeys(user.userId),
    ...userGSI1Keys(user.email),
    ...userPhoneGSI5Keys(user.phone),
    entityType: 'User',
    ...user,
  };
}

export function toVehicleItem(vehicle: UserVehicle) {
  return {
    ...vehicleKeys(vehicle.userId, vehicle.vehicleId),
    entityType: 'UserVehicle',
    ...vehicle,
  };
}
