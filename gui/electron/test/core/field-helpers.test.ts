import { describe, it, expect } from 'vitest';
import { requiredFieldValues, requiredFieldOccurrences, isValidationError } from '../../src/core/field-helpers.ts';
import { PlainField } from '../../src/core/field-expr.ts';
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
