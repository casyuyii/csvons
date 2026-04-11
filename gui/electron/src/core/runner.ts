import * as path from 'node:path';
import { readConfigFile } from './io/config-reader.ts';
import type { ExistsRule, UniqueRule, VTypeRule, ValidationError } from './types.ts';
import { existsTest } from './validators/exists.ts';
import { uniqueTest } from './validators/unique.ts';
import { vtypeTest } from './validators/vtype.ts';
import {
  type ValidationReport,
  type ValidationIssue,
  REPORT_SCHEMA_VERSION,
} from './report.ts';

export interface ValidateOptions {
  basePath?: string;
}

/**
 * Main validation entry point.
 * Reads a ruler.json config, validates all referenced CSV files, returns a report.
 */
export function validate(
  configPath: string,
  options?: ValidateOptions,
): ValidationReport {
  const startAt = performance.now();

  const config = readConfigFile(configPath);
  if (!config) {
    return {
      schema_version: REPORT_SCHEMA_VERSION,
      summary: { files_checked: 0, passed: 0, failed: 0, duration_ms: 0 },
      issues: [
        {
          message: `read config file error: file_name=${configPath}`,
          severity: 'error',
        },
      ],
    };
  }

  const { rules, metadata } = config;
  const basePath =
    options?.basePath ?? path.dirname(path.resolve(configPath));

  const allIssues: ValidationIssue[] = [];
  const stems = Object.keys(rules);
  let failedCount = 0;

  for (const stem of stems) {
    const rawRules = rules[stem]!;
    let stemHasErrors = false;

    // Process exists rules
    if (rawRules.exists) {
      const existsRules = rawRules.exists as ExistsRule[];
      const errors = existsTest(stem, existsRules, metadata, basePath);
      if (errors.length > 0) {
        stemHasErrors = true;
        allIssues.push(...errors.map(toIssue));
      }
    }

    // Process unique rules
    if (rawRules.unique) {
      const uniqueRuler = rawRules.unique as UniqueRule;
      const errors = uniqueTest(stem, uniqueRuler, metadata, basePath);
      if (errors.length > 0) {
        stemHasErrors = true;
        allIssues.push(...errors.map(toIssue));
      }
    }

    // Process vtype rules
    if (rawRules.vtype) {
      const vtypeRules = rawRules.vtype as VTypeRule[];
      const errors = vtypeTest(stem, vtypeRules, metadata, basePath);
      if (errors.length > 0) {
        stemHasErrors = true;
        allIssues.push(...errors.map(toIssue));
      }
    }

    if (stemHasErrors) failedCount++;
  }

  const durationMs = Math.round(performance.now() - startAt);

  return {
    schema_version: REPORT_SCHEMA_VERSION,
    summary: {
      files_checked: stems.length,
      passed: stems.length - failedCount,
      failed: failedCount,
      duration_ms: durationMs,
    },
    issues: allIssues,
  };
}

function toIssue(err: ValidationError): ValidationIssue {
  const issue: ValidationIssue = {
    message: err.message,
    severity: err.severity,
  };
  if (err.file) issue.file = err.file;
  if (err.rule) issue.rule = err.rule;
  if (err.field) issue.field = err.field;
  if (err.row !== undefined) issue.row = err.row;
  if (err.value) issue.value = err.value;
  return issue;
}
