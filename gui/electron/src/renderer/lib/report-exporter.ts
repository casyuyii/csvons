import type { ValidationIssue } from './issue-filters.ts';

export interface ValidationReport {
  schema_version: string;
  summary: {
    files_checked: number;
    passed: number;
    failed: number;
    duration_ms: number;
  };
  issues: ValidationIssue[];
}

export function toPrettyJson(report: ValidationReport): string {
  return JSON.stringify(report, null, 2) + '\n';
}

export function toMarkdown(report: ValidationReport): string {
  const lines: string[] = [];
  lines.push('# csvons Validation Report');
  lines.push('');
  lines.push('## Summary');
  lines.push('');
  lines.push(`- **Schema Version**: ${report.schema_version}`);
  lines.push(`- **Files Checked**: ${report.summary.files_checked}`);
  lines.push(`- **Passed**: ${report.summary.passed}`);
  lines.push(`- **Failed**: ${report.summary.failed}`);
  lines.push(`- **Duration**: ${report.summary.duration_ms}ms`);
  lines.push('');

  if (report.issues.length === 0) {
    lines.push('No issues found.');
  } else {
    lines.push('## Issues');
    lines.push('');
    lines.push(
      '| Severity | File | Rule | Field | Row | Value | Message |',
    );
    lines.push(
      '|----------|------|------|-------|-----|-------|---------|',
    );
    for (const issue of report.issues) {
      const md = (s?: string) =>
        (s ?? '').replace(/\|/g, '\\|').replace(/\n/g, ' ');
      lines.push(
        `| ${md(issue.severity)} | ${md(issue.file)} | ${md(issue.rule)} | ${md(issue.field)} | ${issue.row ?? ''} | ${md(issue.value)} | ${md(issue.message)} |`,
      );
    }
  }

  lines.push('');
  return lines.join('\n');
}
