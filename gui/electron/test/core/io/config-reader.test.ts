import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { RULER_DIR } from '../../test-helpers.ts';

describe('readConfigFile', () => {
  it('reads ruler.json correctly', () => {
    const result = readConfigFile(path.join(RULER_DIR, 'ruler.json'));
    expect(result).not.toBeNull();
    expect(result!.metadata.csv_file_folder).toBe('testdata');
    expect(result!.metadata.name_index).toBe(0);
    expect(result!.metadata.data_index).toBe(1);
    expect(result!.metadata.extension).toBe('.csv');
    expect(result!.metadata.lev1_separator).toBe(';');
    expect(result!.metadata.lev2_separator).toBe(':');
    expect(result!.rules).toHaveProperty('username');
  });

  it('reads ruler_products.json correctly', () => {
    const result = readConfigFile(path.join(RULER_DIR, 'ruler_products.json'));
    expect(result).not.toBeNull();
    expect(result!.rules).toHaveProperty('products');
  });

  it('returns null for non-existent file', () => {
    const result = readConfigFile('/nonexistent/ruler.json');
    expect(result).toBeNull();
  });
});
