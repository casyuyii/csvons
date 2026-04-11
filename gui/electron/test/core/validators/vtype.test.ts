import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import { vtypeTest } from '../../../src/core/validators/vtype.ts';
import { readConfigFile } from '../../../src/core/io/config-reader.ts';
import { PROJECT_ROOT, RULER_DIR } from '../../test-helpers.ts';
import type { VTypeRule } from '../../../src/core/types.ts';

function loadRulesAndMetadata(rulerFile: string) {
  const config = readConfigFile(path.join(RULER_DIR, rulerFile));
  if (!config) throw new Error(`Failed to load ${rulerFile}`);
  return config;
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

  it('detects out of range value (negative test)', () => {
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
});
