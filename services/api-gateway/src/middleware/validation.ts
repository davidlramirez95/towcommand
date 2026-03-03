import { ZodSchema } from 'zod';
import { ValidationError } from '@towcommand/core';

export async function validateRequest<T>(
  data: unknown,
  schema: ZodSchema,
): Promise<T> {
  const result = await schema.safeParseAsync(data);

  if (!result.success) {
    throw new ValidationError(result.error);
  }

  return result.data as T;
}
