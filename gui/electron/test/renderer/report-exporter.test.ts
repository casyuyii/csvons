import { describe, test, expect } from 'vitest';
import {
  toPrettyJson,
  toMarkdown,
  type ValidationReport,
} from '../../src/renderer/lib/report-exporter.ts';

const report: ValidationReport = {
  schema_version: '1.0',
  summary: {
    files_checked: 3,
    passed: 2,
    failed: 1,
    duration_ms: 42,
  },
  issues: [
    {
      severity: 'error',
      file: 'data.csv',
      rule: 'exists',
      field: 'user_id',
      row: 5,
      value: '999',
      message: 'Reference not found',
    },
    {
      severity: 'warning',
      file: 'data.csv',
      rule: 'vtype',
      field: 'age',
      row: 3,
      value: 'abc',
      message: 'Invalid integer',
    },
  ],
};

const emptyReport: ValidationReport = {
  schema_version: '1.0',
  summary: { files_checked: 1, passed: 1, failed: 0, duration_ms: 10 },
  issues: [],
};

describe('toPrettyJson', () => {
  test('produces valid JSON', () => {
    const json = toPrettyJson(report);
    const parsed = JSON.parse(json);
    expect(parsed.schema_version).toBe('1.0');
    expect(parsed.issues).toHaveLength(2);
  });

  test('is indented with 2 spaces', () => {
    const json = toPrettyJson(report);
    expect(json).toContain('  "schema_version"');
  });

  test('ends with newline', () => {
    const json = toPrettyJson(report);
    expect(json.endsWith('\n')).toBe(true);
  });
});

describe('toMarkdown', () => {
  test('includes report header', () => {
    const md = toMarkdown(report);
    expect(md).toContain('# csvons Validation Report');
  });

  test('includes summary stats', () => {
    const md = toMarkdown(report);
    expect(md).toContain('**Files Checked**: 3');
    expect(md).toContain('**Passed**: 2');
    expect(md).toContain('**Failed**: 1');
    expect(md).toContain('**Duration**: 42ms');
  });

  test('includes issues table header', () => {
    const md = toMarkdown(report);
    expect(md).toContain('| Severity | File | Rule | Field | Row | Value | Message |');
  });

  test('includes issue rows', () => {
    const md = toMarkdown(report);
    expect(md).toContain('error');
    expect(md).toContain('Reference not found');
    expect(md).toContain('Invalid integer');
  });

  test('escapes pipe characters in values', () => {
    const reportWithPipe: ValidationReport = {
      ...report,
      issues: [
        { severity: 'error', message: 'has|pipe', file: 'f.csv', rule: 'r', field: 'f', row: 1, value: 'a|b' },
      ],
    };
    const md = toMarkdown(reportWithPipe);
    expect(md).toContain('has\\|pipe');
    expect(md).toContain('a\\|b');
  });

  test('handles empty issues', () => {
    const md = toMarkdown(emptyReport);
    expect(md).toContain('No issues found.');
    expect(md).not.toContain('## Issues');
  });
});
