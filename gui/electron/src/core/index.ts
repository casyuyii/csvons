// Types
export type {
  Metadata,
  ExistsRule,
  UniqueRule,
  VTypeRule,
  ConstraintsConfig,
  ValidationContext,
} from './types.ts';
export { ValidationError, csvFileName, makeValidationError } from './types.ts';

// Field expressions
export type { FieldExpr } from './field-expr.ts';
export {
  PlainField,
  RepeatField,
  NestedField,
  ComplexField,
  generateFieldExpr,
} from './field-expr.ts';

// Field occurrences
export type {
  FieldOccurrence,
  FieldOccurrenceProvider,
} from './field-occurrences.ts';
export { isFieldOccurrenceProvider } from './field-occurrences.ts';

// Field helpers
export {
  requiredFieldValues,
  requiredFieldOccurrences,
  isValidationError,
} from './field-helpers.ts';

// I/O
export { readCsvFile } from './io/csv-reader.ts';
export { readConfigFile } from './io/config-reader.ts';

// Validators
export { existsTest } from './validators/exists.ts';
export { uniqueTest } from './validators/unique.ts';
export { vtypeTest } from './validators/vtype.ts';

// Runner
export { validate } from './runner.ts';
export type { ValidateOptions } from './runner.ts';

// Report
export type {
  ValidationReport,
  ValidationSummary,
  ValidationIssue,
} from './report.ts';
export { REPORT_SCHEMA_VERSION } from './report.ts';
