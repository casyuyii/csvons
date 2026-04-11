/**
 * Metadata describes the structure and location of CSV files being validated.
 */
export interface Metadata {
  csv_file_folder: string;
  name_index: number;
  data_index: number;
  extension: string;
  lev1_separator: string;
  lev2_separator: string;
  field_connector: string;
}

/**
 * ExistsRule defines a cross-file existence constraint.
 */
export interface ExistsRule {
  dst_file_stem: string;
  fields: Array<{ src: string; dst: string }>;
}

/**
 * UniqueRule defines a column uniqueness constraint.
 */
export interface UniqueRule {
  fields: string[];
}

/**
 * VTypeRule defines a value type and optional range constraint.
 */
export interface VTypeRule {
  field: string;
  type: 'int' | 'float64' | 'bool';
  range?: { min: number; max: number };
}

/**
 * ConstraintsConfig represents per-file constraint rules from ruler.json.
 */
export interface ConstraintsConfig {
  exists?: ExistsRule[];
  unique?: UniqueRule;
  vtype?: VTypeRule[];
}

/**
 * ValidationError carries structured context about a validation or runtime failure.
 */
export class ValidationError extends Error {
  file: string;
  rule: string;
  field: string;
  row: number | undefined;
  value: string;
  severity: string;
  code: number;

  constructor(opts: {
    file?: string;
    rule?: string;
    field?: string;
    row?: number;
    value?: string;
    message: string;
    severity?: string;
    code: number;
  }) {
    super(opts.message);
    this.name = 'ValidationError';
    this.file = opts.file ?? '';
    this.rule = opts.rule ?? '';
    this.field = opts.field ?? '';
    this.row = opts.row;
    this.value = opts.value ?? '';
    this.severity = opts.severity ?? 'error';
    this.code = opts.code;
  }

  exitCode(): number {
    return this.code === 2 ? 2 : 1;
  }
}

/**
 * ValidationContext carries optional metadata for building errors.
 */
export interface ValidationContext {
  file?: string;
  rule?: string;
  field?: string;
  row?: number;
  value?: string;
  severity?: string;
}

export function makeValidationError(
  ctx: ValidationContext,
  code: number,
  message: string,
): ValidationError {
  return new ValidationError({
    file: ctx.file,
    rule: ctx.rule,
    field: ctx.field,
    row: ctx.row,
    value: ctx.value,
    message,
    severity: ctx.severity ?? 'error',
    code,
  });
}

export function csvFileName(stem: string, metadata: Metadata | null): string {
  if (!metadata) return stem;
  return stem + metadata.extension;
}
