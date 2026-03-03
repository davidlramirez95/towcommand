import { AppError } from '@towcommand/core';

export type Role = 'customer' | 'provider' | 'fleet_manager' | 'ops_agent' | 'admin';

export function requireRole(...allowedRoles: Role[]) {
  return (userType: string): void => {
    if (!allowedRoles.includes(userType as Role)) {
      throw AppError.forbidden(`Role '${userType}' not authorized. Required: ${allowedRoles.join(', ')}`);
    }
  };
}

export function requireOwnerOrRole(resourceOwnerId: string, requesterId: string, ...allowedRoles: Role[]) {
  return (userType: string): void => {
    if (resourceOwnerId === requesterId) return;
    if (!allowedRoles.includes(userType as Role)) {
      throw AppError.forbidden('Not authorized to access this resource');
    }
  };
}
