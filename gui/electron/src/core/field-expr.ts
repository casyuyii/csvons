import type { Metadata } from './types.ts';

/**
 * FieldExpr defines how to extract values from CSV records.
 */
export interface FieldExpr {
  fieldValues(fields: string[], records: string[][]): Iterable<string> | null;
  typeString(): string;
  init(metadata: Metadata, expr: string): void;
}

/**
 * PlainField: direct column name reference (e.g., "Username")
 */
export class PlainField implements FieldExpr {
  metadata!: Metadata;
  fieldName!: string;

  init(metadata: Metadata, expr: string): void {
    this.metadata = metadata;
    this.fieldName = expr;
  }

  *fieldValues(fields: string[], records: string[][]): Generator<string> {
    const fieldIndex = fields.indexOf(this.fieldName);
    if (fieldIndex === -1) return;

    for (let i = this.metadata.data_index; i < records.length; i++) {
      const record = records[i]!;
      if (fieldIndex < record.length) {
        yield record[fieldIndex]!;
      }
    }
  }

  typeString(): string {
    return 'plain';
  }
}

/**
 * RepeatField: array expansion (e.g., "Tags[]")
 * Splits cell values by lev1_separator and yields each element.
 */
export class RepeatField implements FieldExpr {
  metadata!: Metadata;
  fieldName!: string;

  init(metadata: Metadata, expr: string): void {
    this.metadata = metadata;
    this.fieldName = expr.slice(0, -2); // Remove trailing "[]"
  }

  *fieldValues(fields: string[], records: string[][]): Generator<string> {
    const fieldIndex = fields.indexOf(this.fieldName);
    if (fieldIndex === -1) return;

    for (let i = this.metadata.data_index; i < records.length; i++) {
      const record = records[i]!;
      if (fieldIndex < record.length) {
        const lev1Vals = record[fieldIndex]!.split(this.metadata.lev1_separator);
        for (const val of lev1Vals) {
          yield val;
        }
      }
    }
  }

  typeString(): string {
    return 'repeat';
  }
}

/**
 * NestedField: two-level nested array access (e.g., "marks{1}")
 * Splits by lev1_separator, then by lev2_separator, extracts value at index.
 */
export class NestedField implements FieldExpr {
  metadata!: Metadata;
  fieldName!: string;
  index!: number;

  init(metadata: Metadata, expr: string): void {
    this.metadata = metadata;
    const matches = expr.match(/^([a-zA-Z0-9]+)\{(\d+)\}$/);
    if (matches && matches[1] && matches[2]) {
      this.fieldName = matches[1];
      this.index = parseInt(matches[2], 10);
    }
  }

  *fieldValues(fields: string[], records: string[][]): Generator<string> {
    const fieldIndex = fields.indexOf(this.fieldName);
    if (fieldIndex === -1) return;

    for (let i = this.metadata.data_index; i < records.length; i++) {
      const record = records[i]!;
      if (fieldIndex < record.length) {
        const lev1Vals = record[fieldIndex]!.split(this.metadata.lev1_separator);
        for (const lev1Val of lev1Vals) {
          const lev2Vals = lev1Val.split(this.metadata.lev2_separator);
          if (this.index < lev2Vals.length) {
            yield lev2Vals[this.index]!;
          }
        }
      }
    }
  }

  typeString(): string {
    return 'nested';
  }
}

/**
 * ComplexField: multi-field concatenation (e.g., "{data}{key}")
 * Concatenates values from multiple columns using field_connector.
 */
export class ComplexField implements FieldExpr {
  metadata!: Metadata;
  fieldNames!: string[];

  init(metadata: Metadata, expr: string): void {
    this.metadata = metadata;
    this.fieldNames = [...expr.matchAll(/([a-zA-Z0-9]+)/g)].map((m) => m[1]!);
  }

  *fieldValues(fields: string[], records: string[][]): Generator<string> {
    const fieldIndexes = this.fieldNames.map((name) => fields.indexOf(name));
    if (fieldIndexes.some((idx) => idx === -1)) return;

    for (let i = this.metadata.data_index; i < records.length; i++) {
      const record = records[i]!;
      let cpxStr = '';
      let valid = true;
      for (const fieldIndex of fieldIndexes) {
        if (fieldIndex < record.length) {
          cpxStr += record[fieldIndex]! + this.metadata.field_connector;
        } else {
          valid = false;
          break;
        }
      }
      if (valid) {
        yield cpxStr;
      }
    }
  }

  typeString(): string {
    return 'complex';
  }
}

/**
 * Pattern map for field expression factory.
 */
const fieldExprPatterns: Array<{ pattern: RegExp; create: () => FieldExpr }> = [
  { pattern: /^[a-zA-Z0-9]+$/, create: () => new PlainField() },
  { pattern: /^[a-zA-Z0-9]+\[\]$/, create: () => new RepeatField() },
  { pattern: /^[a-zA-Z0-9]+\{\d+\}$/, create: () => new NestedField() },
  { pattern: /^\{[a-zA-Z0-9]+\}+$/, create: () => new ComplexField() },
];

/**
 * Creates and initializes a FieldExpr from a raw expression string.
 * Returns null if no pattern matches.
 */
export function generateFieldExpr(
  metadata: Metadata,
  fieldExpr: string,
): FieldExpr | null {
  if (!metadata) {
    throw new Error('metadata is nil');
  }

  for (const { pattern, create } of fieldExprPatterns) {
    if (pattern.test(fieldExpr)) {
      const expr = create();
      expr.init(metadata, fieldExpr);
      return expr;
    }
  }

  return null;
}
