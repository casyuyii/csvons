import { describe, it, expect } from 'vitest';
import * as path from 'node:path';
import * as fs from 'node:fs';
import * as os from 'node:os';
import { validate } from '../../src/core/runner.ts';
import { REPORT_SCHEMA_VERSION } from '../../src/core/report.ts';
import { PROJECT_ROOT, RULER_DIR, TESTDATA_DIR } from '../test-helpers.ts';

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

  it('propagates vtype errors into issues (covers lines 78-79 and toIssue)', () => {
    // Write a temp ruler that has a vtype rule that will fail
    const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'csvons-test-'));
    try {
      // Copy the username CSV so it's in the expected location
      const csvSrcDir = path.join(tempDir, 'testdata');
      fs.mkdirSync(csvSrcDir, { recursive: true });
      fs.copyFileSync(
        path.join(TESTDATA_DIR, 'username.csv'),
        path.join(csvSrcDir, 'username.csv'),
      );

      const rulerConfig = {
        username: {
          vtype: [{ field: 'Height', type: 'int' }],
        },
        csvons_metadata: {
          csv_file_folder: 'testdata',
          name_index: 0,
          data_index: 1,
          extension: '.csv',
          lev1_separator: ';',
          lev2_separator: ':',
          field_connector: '|',
        },
      };
      const rulerPath = path.join(tempDir, 'ruler_vtype_fail.json');
      fs.writeFileSync(rulerPath, JSON.stringify(rulerConfig));

      const report = validate(rulerPath, { basePath: tempDir });
      expect(report.summary.failed).toBeGreaterThan(0);
      // Issues should contain vtype errors with file, rule, field, row, value
      const issue = report.issues.find((i) => i.rule === 'vtype');
      expect(issue).toBeDefined();
      expect(issue!.file).toBeDefined();
      expect(issue!.rule).toBe('vtype');
      expect(issue!.field).toBeDefined();
      expect(issue!.row).toBeDefined();
      expect(issue!.value).toBeDefined();
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('propagates exists errors into issues (covers toIssue with all fields)', () => {
    const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'csvons-test-'));
    try {
      const csvSrcDir = path.join(tempDir, 'testdata');
      fs.mkdirSync(csvSrcDir, { recursive: true });
      // username-d2.csv has different usernames to cause exists failure
      fs.copyFileSync(
        path.join(TESTDATA_DIR, 'username.csv'),
        path.join(csvSrcDir, 'username.csv'),
      );
      // Create a destination CSV with only one username so others will fail
      // Use 3+ columns to avoid PapaParse UndetectableDelimiter error
      fs.writeFileSync(
        path.join(csvSrcDir, 'lookup.csv'),
        'Username,Extra1,Extra2\nbooker12,x,y\n',
      );

      const rulerConfig = {
        username: {
          exists: [
            {
              dst_file_stem: 'lookup',
              fields: [{ src: 'Username', dst: 'Username' }],
            },
          ],
        },
        csvons_metadata: {
          csv_file_folder: 'testdata',
          name_index: 0,
          data_index: 1,
          extension: '.csv',
          lev1_separator: ';',
          lev2_separator: ':',
          field_connector: '|',
        },
      };
      const rulerPath = path.join(tempDir, 'ruler_exists_fail.json');
      fs.writeFileSync(rulerPath, JSON.stringify(rulerConfig));

      const report = validate(rulerPath, { basePath: tempDir });
      expect(report.summary.failed).toBeGreaterThan(0);
      const issue = report.issues.find((i) => i.rule === 'exists');
      expect(issue).toBeDefined();
      expect(issue!.file).toBeDefined();
      expect(issue!.row).toBeDefined();
      expect(issue!.value).toBeDefined();
    } finally {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('uses basePath from config file directory when options not provided', () => {
    // Pass ruler path without basePath - it should derive basePath from dirname
    const report = validate(path.join(RULER_DIR, 'ruler.json'));
    // This will fail because testdata is relative to PROJECT_ROOT not RULER_DIR
    // but it should not throw - just produce a report (possibly with errors)
    expect(report).toBeDefined();
    expect(report.schema_version).toBe(REPORT_SCHEMA_VERSION);
  });
});
