import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import * as fs from 'node:fs';
import * as os from 'node:os';
import { vtypeTest } from '../../../src/core/validators/vtype.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { VTypeRule, Metadata } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
}

/** Create a temp dir with a multi-column CSV and return the dir path + metadata */
function makeTempCsv(
  header: string,
  rows: string[],
): { tempDir: string; metadata: Metadata } {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'csvons-vtype-'));
  const csvDir = path.join(tempDir, 'testdata');
  fs.mkdirSync(csvDir, { recursive: true });
  // PapaParse requires 3+ columns to auto-detect the comma delimiter;
  // use a three-column CSV to avoid UndetectableDelimiter errors
  const fullHeader = `${header},_dummy1,_dummy2`;
  const fullRows = rows.map((r) => `${r},x,y`);
  fs.writeFileSync(path.join(csvDir, 'data.csv'), [fullHeader, ...fullRows].join('\n'));
  const metadata: Metadata = {
    csv_file_folder: 'testdata',
    name_index: 0,
    data_index: 1,
    extension: '.csv',
    lev1_separator: ';',
    lev2_separator: ':',
    field_connector: '|',
  };
  return { tempDir, metadata };
}

describe('vtypeTest', () => {
  it('passes for ruler.json (username vtype checks)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler.json');
    const vtypeRules = (rules['username'] as any).vtype as VTypeRule[];
    const errors = vtypeTest('username', vtypeRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_products.json (float64, int, bool)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_products.json');
    const vtypeRules = (rules['products'] as any).vtype as VTypeRule[];
    const errors = vtypeTest('products', vtypeRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_orders.json (nested int with range)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_orders.json');
    const vtypeRules = (rules['orders'] as any).vtype as VTypeRule[];
    const errors = vtypeTest('orders', vtypeRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_employees.json (int with range, bool)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_employees.json');
    const vtypeRules = (rules['employees'] as any).vtype as VTypeRule[];
    const errors = vtypeTest('employees', vtypeRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('detects invalid int value (negative test)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const badRules: VTypeRule[] = [{ field: 'Height', type: 'int' }];
    const errors = vtypeTest('username', badRules, metadata, PROJECT_ROOT);
    expect(errors.length).toBeGreaterThan(0);
    expect(errors[0]!.rule).toBe('vtype');
    expect(errors[0]!.code).toBe(1);
  });

  it('detects invalid bool value (negative test)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const badRules: VTypeRule[] = [{ field: 'Username', type: 'bool' }];
    const errors = vtypeTest('username', badRules, metadata, PROJECT_ROOT);
    expect(errors.length).toBeGreaterThan(0);
    expect(errors[0]!.rule).toBe('vtype');
    expect(errors[0]!.message).toContain('is not a bool');
  });

  it('detects out of range int value (negative test)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const badRules: VTypeRule[] = [
      { field: 'Age', type: 'int', range: { min: 1, max: 10 } },
    ];
    const errors = vtypeTest('username', badRules, metadata, PROJECT_ROOT);
    expect(errors.length).toBeGreaterThan(0);
    expect(errors[0]!.message).toContain('is not in the range');
  });

  it('returns error for empty rules', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const errors = vtypeTest('username', [], metadata, PROJECT_ROOT);
    expect(errors.length).toBe(1);
    expect(errors[0]!.code).toBe(2);
  });

  it('accepts all valid bool values', () => {
    const { metadata } = loadRulesAndMetadata('ruler_products.json');
    const boolRules: VTypeRule[] = [{ field: 'Available', type: 'bool' }];
    const errors = vtypeTest('products', boolRules, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('detects invalid float64 value (lines 151-168)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['notanumber']);
    try {
      const rules: VTypeRule[] = [{ field: 'Value', type: 'float64' }];
      const errors = vtypeTest('data', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.message).toContain('is not a float64');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('detects out of range float64 value (lines 169-186)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['3.14', '999.9']);
    try {
      const rules: VTypeRule[] = [
        { field: 'Value', type: 'float64', range: { min: 0, max: 5.0 } },
      ];
      const errors = vtypeTest('data', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.message).toContain('is not in the range');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('accepts valid float64 within range', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['3.14', '2.71']);
    try {
      const rules: VTypeRule[] = [
        { field: 'Value', type: 'float64', range: { min: 0, max: 10.0 } },
      ];
      const errors = vtypeTest('data', rules, metadata, tempDir);
      expect(errors).toEqual([]);
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('records default/unknown type as code-2 error (lines 211-224)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['hello']);
    try {
      // Cast to bypass TypeScript type check - simulate runtime unknown type
      const rules = [{ field: 'Value', type: 'unknown_type' }] as unknown as VTypeRule[];
      const errors = vtypeTest('data', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('is not a valid type');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error for invalid name_index (< 0)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['42']);
    try {
      const badMeta: Metadata = { ...metadata, name_index: -1 };
      const rules: VTypeRule[] = [{ field: 'Value', type: 'int' }];
      const errors = vtypeTest('data', rules, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('name_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when data_index <= name_index', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['42']);
    try {
      const badMeta: Metadata = { ...metadata, name_index: 1, data_index: 1 };
      const rules: VTypeRule[] = [{ field: 'Value', type: 'int' }];
      const errors = vtypeTest('data', rules, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('data_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });
});
