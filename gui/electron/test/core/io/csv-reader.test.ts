import { describe, it, expect } from 'vitest';
import { readCsvFile } from '../../../src/core/io/csv-reader.ts';
import { PROJECT_ROOT } from '../../test-helpers.ts';
import type { Metadata } from '../../../src/core/types.ts';

const metadata: Metadata = {
  csv_file_folder: 'testdata',
  name_index: 0,
  data_index: 1,
  extension: '.csv',
  lev1_separator: ';',
  lev2_separator: ':',
  field_connector: '|',
};

describe('readCsvFile', () => {
  it('reads username.csv correctly', () => {
    const records = readCsvFile('username', metadata, PROJECT_ROOT);
    expect(records).not.toBeNull();
    expect(records!.length).toBeGreaterThan(1);
    expect(records![0]).toContain('Username');
  });

  it('reads products.csv correctly', () => {
    const records = readCsvFile('products', metadata, PROJECT_ROOT);
    expect(records).not.toBeNull();
    expect(records![0]).toContain('ProductID');
  });

  it('returns null for non-existent file', () => {
    const records = readCsvFile('nonexistent', metadata, PROJECT_ROOT);
    expect(records).toBeNull();
  });

  it('reads employees.csv correctly', () => {
    const records = readCsvFile('employees', metadata, PROJECT_ROOT);
    expect(records).not.toBeNull();
    expect(records![0]).toContain('EmployeeID');
  });
});
