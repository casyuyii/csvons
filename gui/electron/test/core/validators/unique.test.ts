import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import { uniqueTest } from '../../../src/core/validators/unique.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { UniqueRule } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
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
});
