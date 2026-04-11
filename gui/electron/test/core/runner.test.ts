import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import { validate } from '../../src/core/runner.ts';
import { REPORT_SCHEMA_VERSION } from '../../src/core/report.ts';
import { PROJECT_ROOT, RULER_DIR } from '../test-helpers.ts';

describe('validate', () => {
  it('validates ruler.json successfully', () => {
    const report = validate(path.join(RULER_DIR, 'ruler.json'), {
      basePath: PROJECT_ROOT,
    });
    expect(report.schema_version).toBe(REPORT_SCHEMA_VERSION);
    expect(report.summary.failed).toBe(0);
    expect(report.issues).toEqual([]);
    expect(report.summary.files_checked).toBeGreaterThan(0);
    expect(report.summary.passed).toBe(report.summary.files_checked);
  });

  it('validates ruler_products.json successfully', () => {
    const report = validate(path.join(RULER_DIR, 'ruler_products.json'), {
      basePath: PROJECT_ROOT,
    });
    expect(report.summary.failed).toBe(0);
    expect(report.issues).toEqual([]);
  });

  it('validates ruler_orders.json successfully', () => {
    const report = validate(path.join(RULER_DIR, 'ruler_orders.json'), {
      basePath: PROJECT_ROOT,
    });
    expect(report.summary.failed).toBe(0);
    expect(report.issues).toEqual([]);
  });

  it('validates ruler_employees.json successfully', () => {
    const report = validate(path.join(RULER_DIR, 'ruler_employees.json'), {
      basePath: PROJECT_ROOT,
    });
    expect(report.summary.failed).toBe(0);
    expect(report.issues).toEqual([]);
  });

  it('returns error for non-existent config file', () => {
    const report = validate('/nonexistent/ruler.json');
    expect(report.issues.length).toBeGreaterThan(0);
    expect(report.issues[0]!.severity).toBe('error');
    expect(report.issues[0]!.message).toContain('read config file error');
  });

  it('includes duration_ms in summary', () => {
    const report = validate(path.join(RULER_DIR, 'ruler.json'), {
      basePath: PROJECT_ROOT,
    });
    expect(report.summary.duration_ms).toBeGreaterThanOrEqual(0);
  });
});
