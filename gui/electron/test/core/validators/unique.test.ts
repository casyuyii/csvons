import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import * as fs from 'node:fs';
import * as os from 'node:os';
import { uniqueTest } from '../../../src/core/validators/unique.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { UniqueRule, Metadata } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
}

/** Create a temp dir with a two-column CSV and return the dir path + metadata */
function makeTempCsv(
  header: string,
  rows: string[],
): { tempDir: string; metadata: Metadata } {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'csvons-unique-'));
  const csvDir = path.join(tempDir, 'testdata');
  fs.mkdirSync(csvDir, { recursive: true });
  // PapaParse requires 3+ columns to auto-detect the comma delimiter
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

describe('uniqueTest', () => {
  it('passes for ruler.json (username unique checks)', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler.json');
    const uniqueRuler = (rules['username'] as any).unique as UniqueRule;
    const errors = uniqueTest('username', uniqueRuler, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_products.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_products.json');
    const uniqueRuler = (rules['products'] as any).unique as UniqueRule;
    const errors = uniqueTest('products', uniqueRuler, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_orders.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_orders.json');
    const uniqueRuler = (rules['orders'] as any).unique as UniqueRule;
    const errors = uniqueTest('orders', uniqueRuler, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('passes for ruler_employees.json', () => {
    const { rules, metadata } = loadRulesAndMetadata('ruler_employees.json');
    const uniqueRuler = (rules['employees'] as any).unique as UniqueRule;
    const errors = uniqueTest('employees', uniqueRuler, metadata, PROJECT_ROOT);
    expect(errors).toEqual([]);
  });

  it('detects duplicate values (negative test)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const ruler: UniqueRule = { fields: ['Username'] };
    const errors = uniqueTest('username-d1', ruler, metadata, PROJECT_ROOT);
    expect(errors.length).toBeGreaterThan(0);
    expect(errors[0]!.rule).toBe('unique');
    expect(errors[0]!.code).toBe(1);
  });

  it('returns error for empty fields', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const errors = uniqueTest('username', { fields: [] }, metadata, PROJECT_ROOT);
    expect(errors.length).toBe(1);
    expect(errors[0]!.code).toBe(2);
  });

  it('returns error when name_index < 0 (line 33-42)', () => {
    // Use the pre-built CSV from makeTempCsv (two-column format)
    const { tempDir, metadata } = makeTempCsv('Value', ['a', 'b']);
    try {
      const badMeta: Metadata = { ...metadata, name_index: -1 };
      const errors = uniqueTest('data', { fields: ['Value'] }, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('name_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when data_index <= name_index (line 44-53)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['a']);
    try {
      const badMeta: Metadata = { ...metadata, name_index: 1, data_index: 1 };
      const errors = uniqueTest('data', { fields: ['Value'] }, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('data_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when src CSV does not exist (line 57-66)', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const errors = uniqueTest('nonexistent_file', { fields: ['Username'] }, metadata, PROJECT_ROOT);
    expect(errors.length).toBe(1);
    expect(errors[0]!.code).toBe(2);
    expect(errors[0]!.message).toContain('src_records length');
  });

  it('returns error when field expression is invalid (line 81-82)', () => {
    const { tempDir, metadata } = makeTempCsv('Value', ['a', 'b']);
    try {
      // Use an invalid field expression like "invalid-name" which generates null FieldExpr
      const errors = uniqueTest('data', { fields: ['invalid-name'] }, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('field expression');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });
});
