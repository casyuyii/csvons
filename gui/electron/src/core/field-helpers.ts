import type { FieldExpr } from './field-expr.ts';
import type {
  FieldOccurrence,
  FieldOccurrenceProvider,
} from './field-occurrences.ts';
import { isFieldOccurrenceProvider } from './field-occurrences.ts';
import {
  type ValidationContext,
  makeValidationError,
  type ValidationError,
} from './types.ts';

/**
 * Validates a field expression and returns its values iterable.
 * Returns a ValidationError if the expression is nil or cannot resolve values.
 */
export function requiredFieldValues(
  fieldExpr: FieldExpr | null,
  fieldName: string,
  fields: string[],
  records: string[][],
): Iterable<string> | ValidationError {
  if (!fieldExpr) {
    return makeValidationError({}, 2, `field expression [${fieldName}] is nil`);
  }
  const vals = fieldExpr.fieldValues(fields, records);
  if (!vals) {
    return makeValidationError(
      {},
      2,
      `field expression [${fieldName}] cannot resolve values`,
    );
  }
  return vals;
}

/**
 * Validates a field expression and returns its occurrences iterable.
 * Returns a ValidationError if the expression is nil or cannot resolve values.
 */
export function requiredFieldOccurrences(
  fieldExpr: FieldExpr | null,
  fieldName: string,
  fields: string[],
  records: string[][],
  ctx: ValidationContext,
): Iterable<FieldOccurrence> | ValidationError {
  const fieldCtx = { ...ctx, field: fieldName };

  if (!fieldExpr) {
    return makeValidationError(
      fieldCtx,
      2,
      `field expression [${fieldName}] is nil`,
    );
  }

  if (!isFieldOccurrenceProvider(fieldExpr)) {
    return makeValidationError(
      fieldCtx,
      2,
      `field expression [${fieldName}] cannot resolve values`,
    );
  }

  const occurrences = (fieldExpr as FieldOccurrenceProvider).fieldOccurrences(
    fields,
    records,
  );
  if (!occurrences) {
    return makeValidationError(
      fieldCtx,
      2,
      `field expression [${fieldName}] cannot resolve values`,
    );
  }
  return occurrences;
}

/**
 * Type guard to check if a result is a ValidationError.
 */
export function isValidationError(
  result: unknown,
): result is ValidationError {
  return result instanceof Error && 'code' in result;
}
