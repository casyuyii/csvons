import { describe, test, expect } from 'vitest';
import {
  filterAndSortIssues,
  type ValidationIssue,
} from '../../src/renderer/lib/issue-filters.ts';

const issues: ValidationIssue[] = [
  { severity: 'error', file: 'a.csv', rule: 'exists', field: 'id', row: 2, value: '10', message: 'Missing ref' },
  { severity: 'warning', file: 'b.csv', rule: 'unique', field: 'name', row: 5, value: 'dup', message: 'Duplicate name' },
  { severity: 'error', file: 'a.csv', rule: 'vtype', field: 'age', row: 3, value: 'abc', message: 'Invalid int' },
  { severity: 'info', file: 'c.csv', rule: 'exists', field: 'code', row: 1, value: '', message: 'Info note' },
  { severity: 'critical', file: 'a.csv', rule: 'exists', field: 'id', row: 10, value: '', message: 'Critical failure' },
];

describe('filterAndSortIssues', () => {
  test('returns all issues when no filters', () => {
    const result = filterAndSortIssues({ issues });
    expect(result).toHaveLength(5);
  });

  test('filters by severity', () => {
    const result = filterAndSortIssues({ issues, severityFilter: 'error' });
    expect(result).toHaveLength(2);
    expect(result.every((i) => i.severity === 'error')).toBe(true);
  });

  test('filters by file', () => {
    const result = filterAndSortIssues({ issues, fileFilter: 'a.csv' });
    expect(result).toHaveLength(3);
  });

  test('filters by rule', () => {
    const result = filterAndSortIssues({ issues, ruleFilter: 'exists' });
    expect(result).toHaveLength(3);
  });

  test('filters by search query', () => {
    const result = filterAndSortIssues({ issues, query: 'duplicate' });
    expect(result).toHaveLength(1);
    expect(result[0].message).toBe('Duplicate name');
  });

  test('query matches across multiple fields', () => {
    const result = filterAndSortIssues({ issues, query: 'abc' });
    expect(result).toHaveLength(1);
    expect(result[0].value).toBe('abc');
  });

  test('combines multiple filters', () => {
    const result = filterAndSortIssues({
      issues,
      severityFilter: 'error',
      fileFilter: 'a.csv',
    });
    expect(result).toHaveLength(2);
  });

  test('sorts by severity rank (default)', () => {
    const result = filterAndSortIssues({ issues });
    expect(result[0].severity).toBe('critical');
    expect(result[result.length - 1].severity).toBe('info');
  });

  test('sorts by severity descending', () => {
    const result = filterAndSortIssues({ issues, ascending: false });
    expect(result[0].severity).toBe('info');
    expect(result[result.length - 1].severity).toBe('critical');
  });

  test('sorts by row', () => {
    const result = filterAndSortIssues({ issues, sortField: 'row' });
    const rows = result.map((i) => i.row);
    for (let i = 1; i < rows.length; i++) {
      expect(rows[i]!).toBeGreaterThanOrEqual(rows[i - 1]!);
    }
  });

  test('sorts by file', () => {
    const result = filterAndSortIssues({ issues, sortField: 'file' });
    expect(result[0].file).toBe('a.csv');
  });

  test('sorts by message', () => {
    const result = filterAndSortIssues({ issues, sortField: 'message' });
    expect(result[0].message).toBe('Critical failure');
  });

  test('returns empty array for no matches', () => {
    const result = filterAndSortIssues({ issues, query: 'zzzznotfound' });
    expect(result).toHaveLength(0);
  });

  test('handles empty issues array', () => {
    const result = filterAndSortIssues({ issues: [] });
    expect(result).toHaveLength(0);
  });

  test('severity filter is case-insensitive', () => {
    const result = filterAndSortIssues({ issues, severityFilter: 'ERROR' });
    expect(result).toHaveLength(2);
  });
});
