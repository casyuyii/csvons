export const REPORT_SCHEMA_VERSION = 'csvons.validation_report.v1';

export interface ValidationSummary {
  files_checked: number;
  passed: number;
  failed: number;
  duration_ms: number;
}

export interface ValidationIssue {
  file?: string;
  rule?: string;
  field?: string;
  row?: number;
  value?: string;
  message: string;
  severity: string;
}

export interface ValidationReport {
  schema_version: string;
  summary: ValidationSummary;
  issues: ValidationIssue[];
}
