import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import { existsTest } from '../../../src/core/validators/exists.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { ExistsRule, Metadata } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
}

describe('existsTest', () => {
  it('passes for ruler.json (username exists checks)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler.json');
    const existsRules = (rules['username'] as any).exists as ExistsRule[];
    const errors = existsTest('username', existsRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_products.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_products.json');
    const existsRules = (rules['products'] as any).exists as ExistsRule[];
    const errors = existsTest('products', existsRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_orders.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_orders.json');
    const existsRules = (rules['orders'] as any).exists as ExistsRule[];
    const errors = existsTest('orders', existsRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_employees.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_employees.json');
    const existsRules = (rules['employees'] as any).exists as ExistsRule[];
    const errors = existsTest('employees', existsRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('detects missing value in destination (negative test)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const badRules: ExistsRule[] = [
      {
        dst_file_stem: 'nonexistent',
        fields: [{ src: 'Username', dst: 'Username' }],
      },
    ];
    const errors = existsTest('username', badRules, metadata, PROJECT_ROOT);
    expect(errors.length).toBeGreaterThan(0);
    expect(errors[0]!.code).toBe(2);
  });

  it('returns error for empty rules', () => {
    const metadata: Metadata = {
      csv_file_folder: 'testdata',
      name_index: 0,
      data_index: 1,
      extension: '.csv',
      lev1_separator: ';',
      lev2_separator: ':',
      field_connector: '|',
    };
    const errors = existsTest('username', [], metadata, PROJECT_ROOT);
    expect(errors.length).toBe(1);
    expect(errors[0]!.code).toBe(2);
  });
});
