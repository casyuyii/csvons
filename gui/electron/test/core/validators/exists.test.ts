import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import * as fs from 'node:fs';
import * as os from 'node:os';
import { existsTest } from '../../../src/core/validators/exists.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { ExistsRule, Metadata } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
}

/** Create a temp dir with CSVs and return the dir path + base metadata */
function makeTempDir(): { tempDir: string; metadata: Metadata } {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'csvons-exists-'));
  const csvDir = path.join(tempDir, 'testdata');
  fs.mkdirSync(csvDir, { recursive: true });
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

  it('returns error when name_index < 0', () => {
    const { tempDir, metadata } = makeTempDir();
    try {
      fs.writeFileSync(path.join(tempDir, 'testdata', 'src.csv'), 'Name,Extra1,Extra2\nalice,1,2\n');
      const badMeta: Metadata = { ...metadata, name_index: -1 };
      const rules: ExistsRule[] = [{ dst_file_stem: 'src', fields: [{ src: 'Name', dst: 'Name' }] }];
      const errors = existsTest('src', rules, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('name_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when data_index <= name_index', () => {
    const { tempDir, metadata } = makeTempDir();
    try {
      fs.writeFileSync(path.join(tempDir, 'testdata', 'src.csv'), 'Name,Extra1,Extra2\nalice,1,2\n');
      const badMeta: Metadata = { ...metadata, name_index: 1, data_index: 1 };
      const rules: ExistsRule[] = [{ dst_file_stem: 'src', fields: [{ src: 'Name', dst: 'Name' }] }];
      const errors = existsTest('src', rules, badMeta, tempDir);
      expect(errors.length).toBe(1);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('data_index');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when src CSV does not exist', () => {
    const { metadata } = loadRulesAndMetadata('ruler.json');
    const rules: ExistsRule[] = [{ dst_file_stem: 'username', fields: [{ src: 'Username', dst: 'Username' }] }];
    const errors = existsTest('no_such_file', rules, metadata, PROJECT_ROOT);
    expect(errors.length).toBe(1);
    expect(errors[0]!.code).toBe(2);
    expect(errors[0]!.message).toContain('src_records length');
  });

  it('returns error when src field expression is invalid (line 101)', () => {
    const { tempDir, metadata } = makeTempDir();
    try {
      // PapaParse needs 3+ columns to auto-detect delimiter
      fs.writeFileSync(path.join(tempDir, 'testdata', 'src.csv'), 'Name,Extra1,Extra2\nalice,1,a\n');
      fs.writeFileSync(path.join(tempDir, 'testdata', 'dst.csv'), 'Name,Extra1,Extra2\nalice,1,a\n');
      // "invalid-expr" does not match any pattern → null FieldExpr
      const rules: ExistsRule[] = [{ dst_file_stem: 'dst', fields: [{ src: 'invalid-expr', dst: 'Name' }] }];
      const errors = existsTest('src', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('field expression');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('returns error when dst field expression is invalid (lines 114-115)', () => {
    const { tempDir, metadata } = makeTempDir();
    try {
      fs.writeFileSync(path.join(tempDir, 'testdata', 'src.csv'), 'Name,Extra1,Extra2\nalice,1,a\n');
      fs.writeFileSync(path.join(tempDir, 'testdata', 'dst.csv'), 'Name,Extra1,Extra2\nalice,1,a\n');
      // src is valid, but dst is invalid expression
      const rules: ExistsRule[] = [{ dst_file_stem: 'dst', fields: [{ src: 'Name', dst: 'invalid-expr' }] }];
      const errors = existsTest('src', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.code).toBe(2);
      expect(errors[0]!.message).toContain('field expression');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('detects src value not found in dst (line 132)', () => {
    const { tempDir, metadata } = makeTempDir();
    try {
      fs.writeFileSync(path.join(tempDir, 'testdata', 'src.csv'), 'Name,Extra1,Extra2\nalice,1,a\nbob,2,b\n');
      // dst only has alice
      fs.writeFileSync(path.join(tempDir, 'testdata', 'dst.csv'), 'Name,Extra1,Extra2\nalice,1,a\n');
      const rules: ExistsRule[] = [{ dst_file_stem: 'dst', fields: [{ src: 'Name', dst: 'Name' }] }];
      const errors = existsTest('src', rules, metadata, tempDir);
      expect(errors.length).toBeGreaterThan(0);
      expect(errors[0]!.code).toBe(1);
      expect(errors[0]!.message).toContain('not found in dst_records');
      expect(errors[0]!.value).toBe('bob');
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });
});
