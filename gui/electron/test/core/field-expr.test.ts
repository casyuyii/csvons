import { describe, it, expect } from 'vitest';
import {
  PlainField,
  RepeatField,
  NestedField,
  ComplexField,
  generateFieldExpr,
} from '../../src/core/field-expr.ts';
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

describe('PlainField', () => {
  it('extracts values from a single column', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'field1');

    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['value1', 'value2'],
      ['value3', 'value4'],
    ];

    const results = [...field.fieldValues(fields, records)!];
    expect(results).toEqual(['value1', 'value3']);
  });

  it('returns "plain" for typeString', () => {
    expect(new PlainField().typeString()).toBe('plain');
  });

  it('init sets metadata and fieldName', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'testField');
    expect(field.metadata).toBe(baseMetadata);
    expect(field.fieldName).toBe('testField');
  });

  it('returns empty for field not found', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'field1');
    const results = field.fieldValues(['otherField'], [['h'], ['v']]);
    expect([...results!]).toEqual([]);
  });

  it('returns empty for empty records', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'field1');
    const results = [...field.fieldValues(['field1'], [])!];
    expect(results).toEqual([]);
  });
});

describe('RepeatField', () => {
  it('splits cell values by lev1_separator', () => {
    const field = new RepeatField();
    field.init(baseMetadata, 'field1[]');

    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['a,b,c', 'value2'],
      ['x,y', 'value4'],
    ];

    const results = [...field.fieldValues(fields, records)!];
    expect(results).toEqual(['a', 'b', 'c', 'x', 'y']);
  });

  it('returns "repeat" for typeString', () => {
    expect(new RepeatField().typeString()).toBe('repeat');
  });

  it('init strips [] suffix', () => {
    const field = new RepeatField();
    field.init(baseMetadata, 'testField[]');
    expect(field.fieldName).toBe('testField');
  });
});

describe('NestedField', () => {
  it('extracts values at nested index', () => {
    const field = new NestedField();
    field.init(baseMetadata, 'field1{1}');

    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['a:b:c,d:e:f', 'value2'],
      ['x:y:z', 'value4'],
    ];

    const results = [...field.fieldValues(fields, records)!];
    expect(results).toEqual(['b', 'e', 'y']);
  });

  it('returns "nested" for typeString', () => {
    expect(new NestedField().typeString()).toBe('nested');
  });

  it('init parses fieldName and index', () => {
    const field = new NestedField();
    field.init(baseMetadata, 'testField{2}');
    expect(field.fieldName).toBe('testField');
    expect(field.index).toBe(2);
  });
});

describe('ComplexField', () => {
  it('concatenates values from multiple columns', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{field1}{field2}');

    const fields = ['field1', 'field2', 'field3'];
    const records = [
      ['header1', 'header2', 'header3'],
      ['value1', 'value2', 'value3'],
      ['value4', 'value5', 'value6'],
    ];

    const results = [...field.fieldValues(fields, records)!];
    expect(results).toEqual(['value1-value2-', 'value4-value5-']);
  });

  it('returns "complex" for typeString', () => {
    expect(new ComplexField().typeString()).toBe('complex');
  });

  it('init parses field names from braces', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{field1}{field2}');
    expect(field.fieldNames).toEqual(['field1', 'field2']);
  });

  it('returns empty when any field index is out of bounds for a record', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{field1}{field2}');

    // field1 is at index 0, field2 is at index 1
    // row with only 1 element: fieldIndex=1 >= record.length=1 → valid=false
    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['onlyone'],           // record.length=1, fieldIndex for field2=1 → out of bounds
      ['val1', 'val2'],     // valid row
    ];

    const results = [...field.fieldValues(fields, records)!];
    // only the valid row should be yielded
    expect(results).toEqual(['val1-val2-']);
  });

  it('returns empty when a required field name is not in fields list', () => {
    const field = new ComplexField();
    field.init(baseMetadata, '{field1}{missing}');

    const fields = ['field1', 'field2'];
    const records = [
      ['header1', 'header2'],
      ['value1', 'value2'],
    ];

    // missing field causes early return (fieldIndexes.some idx===-1)
    const results = [...field.fieldValues(fields, records)!];
    expect(results).toEqual([]);
  });
});

describe('generateFieldExpr', () => {
  const tests = [
    { name: 'Plain field', expr: 'field1', type: 'plain', nil: false },
    { name: 'Plain with numbers', expr: 'field123', type: 'plain', nil: false },
    { name: 'Repeat field', expr: 'field1[]', type: 'repeat', nil: false },
    { name: 'Repeat with numbers', expr: 'field123[]', type: 'repeat', nil: false },
    { name: 'Nested field', expr: 'field1{0}', type: 'nested', nil: false },
    { name: 'Nested with numbers', expr: 'field123{2}', type: 'nested', nil: false },
    { name: 'Complex single', expr: '{field1}', type: 'complex', nil: false },
    { name: 'Complex multiple', expr: '{field1}{field2}', type: '', nil: true },
    { name: 'Invalid with dash', expr: 'field-1', type: '', nil: true },
    { name: 'Empty string', expr: '', type: '', nil: true },
    { name: 'Special chars', expr: 'field@1', type: '', nil: true },
    { name: 'Mixed case', expr: 'Field1', type: 'plain', nil: false },
    { name: 'Mixed case repeat', expr: 'Field1[]', type: 'repeat', nil: false },
    { name: 'Mixed case nested', expr: 'Field1{1}', type: 'nested', nil: false },
    { name: 'Mixed case complex', expr: '{Field1}', type: 'complex', nil: false },
  ];

  for (const tt of tests) {
    it(tt.name, () => {
      const result = generateFieldExpr(baseMetadata, tt.expr);
      if (tt.nil) {
        expect(result).toBeNull();
      } else {
        expect(result).not.toBeNull();
        expect(result!.typeString()).toBe(tt.type);
      }
    });
  }

  it('throws on null metadata', () => {
    expect(() => generateFieldExpr(null as any, 'field1')).toThrow(
      'metadata is nil',
    );
  });
});
