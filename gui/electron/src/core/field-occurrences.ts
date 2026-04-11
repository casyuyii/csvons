import type { Metadata } from './types.ts';
import type { FieldExpr } from './field-expr.ts';
import {
  PlainField,
  RepeatField,
  NestedField,
  ComplexField,
} from './field-expr.ts';

/**
 * FieldOccurrence captures a value with its 1-based CSV row number.
 */
export interface FieldOccurrence {
  row: number;
  value: string;
}

/**
 * FieldOccurrenceProvider can yield field occurrences with row tracking.
 */
export interface FieldOccurrenceProvider extends FieldExpr {
  fieldOccurrences(
    fields: string[],
    records: string[][],
  ): Iterable<FieldOccurrence> | null;
}

export function isFieldOccurrenceProvider(
  expr: FieldExpr,
): expr is FieldOccurrenceProvider {
  return 'fieldOccurrences' in expr;
}

// Extend PlainField with FieldOccurrences
PlainField.prototype.fieldOccurrences = function* (
  this: PlainField,
  fields: string[],
  records: string[][],
): Generator<FieldOccurrence> {
  const fieldIndex = fields.indexOf(this.fieldName);
  if (fieldIndex === -1) return;

  for (let i = this.metadata.data_index; i < records.length; i++) {
    const record = records[i]!;
    if (fieldIndex < record.length) {
      yield { row: i + 1, value: record[fieldIndex]! };
    }
  }
};

// Extend RepeatField with FieldOccurrences
RepeatField.prototype.fieldOccurrences = function* (
  this: RepeatField,
  fields: string[],
  records: string[][],
): Generator<FieldOccurrence> {
  const fieldIndex = fields.indexOf(this.fieldName);
  if (fieldIndex === -1) return;

  for (let i = this.metadata.data_index; i < records.length; i++) {
    const record = records[i]!;
    if (fieldIndex < record.length) {
      const lev1Vals = record[fieldIndex]!.split(this.metadata.lev1_separator);
      for (const val of lev1Vals) {
        yield { row: i + 1, value: val };
      }
    }
  }
};

// Extend NestedField with FieldOccurrences
NestedField.prototype.fieldOccurrences = function* (
  this: NestedField,
  fields: string[],
  records: string[][],
): Generator<FieldOccurrence> {
  const fieldIndex = fields.indexOf(this.fieldName);
  if (fieldIndex === -1) return;

  for (let i = this.metadata.data_index; i < records.length; i++) {
    const record = records[i]!;
    if (fieldIndex < record.length) {
      const lev1Vals = record[fieldIndex]!.split(this.metadata.lev1_separator);
      for (const lev1Val of lev1Vals) {
        const lev2Vals = lev1Val.split(this.metadata.lev2_separator);
        if (this.index < lev2Vals.length) {
          yield { row: i + 1, value: lev2Vals[this.index]! };
        }
      }
    }
  }
};

// Extend ComplexField with FieldOccurrences
ComplexField.prototype.fieldOccurrences = function* (
  this: ComplexField,
  fields: string[],
  records: string[][],
): Generator<FieldOccurrence> {
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
      yield { row: i + 1, value: cpxStr };
    }
  }
};

// Augment the class interfaces for TypeScript
declare module './field-expr.ts' {
  interface PlainField extends FieldOccurrenceProvider {}
  interface RepeatField extends FieldOccurrenceProvider {}
  interface NestedField extends FieldOccurrenceProvider {}
  interface ComplexField extends FieldOccurrenceProvider {}
}
