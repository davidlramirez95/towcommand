import { ZodError } from 'zod';
import { AppError, ErrorCode } from './app-error';

export class ValidationError extends AppError {
  public readonly fieldErrors: Record<string, string[]>;

  constructor(zodError: ZodError) {
    const fieldErrors: Record<string, string[]> = {};
    for (const issue of zodError.issues) {
      const path = issue.path.join('.');
      if (!fieldErrors[path]) {
        fieldErrors[path] = [];
      }
      fieldErrors[path].push(issue.message);
    }

    super(
      ErrorCode.VALIDATION_ERROR,
      'Validation failed',
      400,
      true,
      { fieldErrors },
    );
    this.fieldErrors = fieldErrors;
    Object.setPrototypeOf(this, ValidationError.prototype);
  }
}
