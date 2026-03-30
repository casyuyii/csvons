import { describe, it, expect } from 'vitest';
import { PlainField, RepeatField, NestedField, ComplexField } from '../../src/core/field-expr.ts';
import '../../src/core/field-occurrences.ts';
import type { FieldOccurrenceProvider } from '../../src/core/field-occurrences.ts';
import type { Metadata } from '../../src/core/types.ts';

const baseMetadata: Metadata = {
  csv_file_folder: '',
  name_index: 0,
  data_index: 1,
  extension: '.csv',
  lev1_separator: ',',
  lev2_separator: ':',
  field_connector: '-',
};

describe('PlainField.fieldOccurrences', () => {
  it('yields values with 1-based row numbers', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'field1');

    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['value1', 'value2'],
      ['value3', 'value4'],
    ];

    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(fields, records)!];
    expect(results).toEqual([
      { row: 2, value: 'value1' },
      { row: 3, value: 'value3' },
    ]);
  });

  it('returns empty for field not found', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'missing');
    const results = (field as unknown as FieldOccurrenceProvider).fieldOccurrences(['other'], [['h'], ['v']]);
    expect([...results!]).toEqual([]);
  });
});

describe('RepeatField.fieldOccurrences', () => {
  it('yields split values with row numbers', () => {
    const field = new RepeatField();
    field.init(baseMetadata, 'field1[]');

    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['field1'],
      [['header'], ['a,b'], ['c']],
    )!];

    expect(results).toEqual([
      { row: 2, value: 'a' },
      { row: 2, value: 'b' },
      { row: 3, value: 'c' },
    ]);
  });
});

describe('NestedField.fieldOccurrences', () => {
  it('yields nested values with row numbers', () => {
    const field = new NestedField();
    field.init(baseMetadata, 'field1{1}');

    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['field1'],
      [['header'], ['a:b,c:d']],
    )!];

    expect(results).toEqual([
      { row: 2, value: 'b' },
      { row: 2, value: 'd' },
    ]);
  });
});

describe('ComplexField.fieldOccurrences', () => {
  it('yields concatenated values with row numbers', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{f1}{f2}');

    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['f1', 'f2'],
      [['h1', 'h2'], ['a', 'b'], ['c', 'd']],
    )!];

    expect(results).toEqual([
      { row: 2, value: 'a-b-' },
      { row: 3, value: 'c-d-' },
    ]);
  });

  it('skips rows where a field index is out of bounds (valid=false branch)', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{f1}{f2}');

    // row with only 1 element: f2 is at index 1, record.length=1 → valid=false
    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['f1', 'f2'],
      [['h1', 'h2'], ['only_one'], ['x', 'y']],
    )!];

    // only the valid row (index 2) should be yielded
    expect(results).toEqual([{ row: 3, value: 'x-y-' }]);
  });

  it('returns empty when required field is missing from fields list', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{f1}{missing}');

    const results = [...(field as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['f1', 'f2'],
      [['h1', 'h2'], ['a', 'b']],
    )!];

    expect(results).toEqual([]);
  });

  it('returns empty for PlainField when field not found', () => {
    const plain = new PlainField();
    plain.init(baseMetadata, 'nothere');

    const results = [...(plain as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['other'],
      [['h'], ['v']],
    )!];
    expect(results).toEqual([]);
  });

  it('returns empty for RepeatField when field not found', () => {
    const repeat = new RepeatField();
    repeat.init(baseMetadata, 'nothere[]');

    const results = [...(repeat as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['other'],
      [['h'], ['a,b']],
    )!];
    expect(results).toEqual([]);
  });

  it('returns empty for NestedField when field not found', () => {
    const nested = new NestedField();
    nested.init(baseMetadata, 'nothere{0}');

    const results = [...(nested as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['other'],
      [['h'], ['a:b']],
    )!];
    expect(results).toEqual([]);
  });

  it('skips nested entry when index is out of bounds', () => {
    const nested = new NestedField();
    nested.init(baseMetadata, 'field1{5}');

    // index 5, but lev2 only has 2 elements
    const results = [...(nested as unknown as FieldOccurrenceProvider).fieldOccurrences(
      ['field1'],
      [['header'], ['a:b']],
    )!];
    expect(results).toEqual([]);
  });
});
