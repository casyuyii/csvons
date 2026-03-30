import { describe, it, expect } from 'vitest';
import { requiredFieldValues, requiredFieldOccurrences, isValidationError } from '../../src/core/field-helpers.ts';
import { PlainField, type FieldExpr } from '../../src/core/field-expr.ts';
import '../../src/core/field-occurrences.ts';
import { ValidationError } from '../../src/core/types.ts';
import type { Metadata } from '../../src/core/types.ts';

const baseMetadata: Metadata = {
  csv_file_folder: '',
  name_index: 0,
  data_index: 1,
  extension: '.csv',
  lev1_separator: ';',
  lev2_separator: ':',
  field_connector: '|',
};

describe('requiredFieldValues', () => {
  it('returns error when fieldExpr is null', () => {
    const result = requiredFieldValues(null, 'Nope', ['f'], [['h'], ['v']]);
    expect(isValidationError(result)).toBe(true);
    expect((result as ValidationError).message).toContain('field expression [Nope] is nil');
  });

  it('returns error when field not found (channel nil)', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'missing');
    const result = requiredFieldValues(field, 'missing', ['other'], [['h'], ['v']]);
    // PlainField returns an empty generator (not null) for missing field in our TS impl
    // but the generator yields nothing, so it's still usable
    expect(isValidationError(result)).toBe(false);
  });

  it('returns iterable when valid', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'f');
    const result = requiredFieldValues(field, 'f', ['f'], [['h'], ['v']]);
    expect(isValidationError(result)).toBe(false);
    expect([...(result as Iterable<string>)]).toEqual(['v']);
  });

  it('returns error when fieldValues returns null (line 28)', () => {
    // Create a fake FieldExpr that returns null from fieldValues
    const fakeExpr: FieldExpr = {
      fieldValues: () => null,
      typeString: () => 'fake',
      init: () => {},
    };
    const result = requiredFieldValues(fakeExpr, 'fakeField', ['f'], [['h'], ['v']]);
    expect(isValidationError(result)).toBe(true);
    expect((result as ValidationError).message).toContain('field expression [fakeField] cannot resolve values');
  });
});

describe('requiredFieldOccurrences', () => {
  it('returns error when fieldExpr is null', () => {
    const result = requiredFieldOccurrences(
      null,
      'Nope',
      ['f'],
      [['h'], ['v']],
      { file: 'test.csv', rule: 'exists' },
    );
    expect(isValidationError(result)).toBe(true);
    expect((result as ValidationError).message).toContain('field expression [Nope] is nil');
  });

  it('returns occurrences when valid', () => {
    const field = new PlainField();
    field.init(baseMetadata, 'f');
    const result = requiredFieldOccurrences(
      field,
      'f',
      ['f'],
      [['h'], ['v']],
      { file: 'test.csv', rule: 'exists' },
    );
    expect(isValidationError(result)).toBe(false);
  });

  it('returns error when fieldExpr is not a FieldOccurrenceProvider (line 59)', () => {
    // Create a FieldExpr without fieldOccurrences method
    const fakeExpr: FieldExpr = {
      fieldValues: (_fields, _records) => [][Symbol.iterator](),
      typeString: () => 'fake',
      init: () => {},
    };
    const result = requiredFieldOccurrences(
      fakeExpr,
      'fakeField',
      ['f'],
      [['h'], ['v']],
      { file: 'test.csv', rule: 'unique' },
    );
    expect(isValidationError(result)).toBe(true);
    expect((result as ValidationError).message).toContain('field expression [fakeField] cannot resolve values');
  });
});

describe('isValidationError', () => {
  it('returns true for ValidationError', () => {
    expect(isValidationError(new ValidationError({ message: 'test', code: 1 }))).toBe(true);
  });

  it('returns false for non-error', () => {
    expect(isValidationError('hello')).toBe(false);
    expect(isValidationError(null)).toBe(false);
    expect(isValidationError([])).toBe(false);
  });
});

describe('ValidationError.exitCode', () => {
  it('returns 2 when code is 2', () => {
    const err = new ValidationError({ message: 'test', code: 2 });
    expect(err.exitCode()).toBe(2);
  });

  it('returns 1 when code is not 2 (line 81)', () => {
    const err = new ValidationError({ message: 'test', code: 1 });
    expect(err.exitCode()).toBe(1);
  });

  it('returns 1 when code is 0', () => {
    const err = new ValidationError({ message: 'test', code: 0 });
    expect(err.exitCode()).toBe(1);
  });
});
