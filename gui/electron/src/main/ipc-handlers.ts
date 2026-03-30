import { ipcMain, dialog } from 'electron';
import * as fs from 'node:fs';
import * as path from 'node:path';
import Papa from 'papaparse';
import { validate } from '../core/runner.ts';
import type { ValidationReport } from '../core/report.ts';

const STATE_FILE = '.csvons_gui_state.json';

interface LocalState {
  recentRulerPaths: string[];
  recentWorkspacePaths: string[];
  recentExportPaths: string[];
}

function getStatePath(): string {
  const home =
    process.env.HOME || process.env.USERPROFILE || process.cwd();
  return path.join(home, STATE_FILE);
}

export function registerIpcHandlers(): void {
  // Validate
  ipcMain.handle(
    'csvons:validate',
    async (_event, configPath: string): Promise<ValidationReport> => {
      const basePath = path.dirname(path.resolve(configPath));
      return validate(configPath, { basePath });
    },
  );

  // CSV Preview
  ipcMain.handle(
    'csvons:csv-preview',
    async (
      _event,
      filePath: string,
      maxRows: number = 6,
    ): Promise<{ headers: string[]; rows: string[][]; totalRows: number }> => {
      const content = fs.readFileSync(filePath, 'utf-8');
      const result = Papa.parse<string[]>(content, {
        header: false,
        skipEmptyLines: false,
      });
      const data = result.data;
      const headers = data[0] ?? [];
      const rows = data.slice(1, 1 + maxRows);
      return { headers, rows, totalRows: Math.max(0, data.length - 1) };
    },
  );

  // Load state
  ipcMain.handle('csvons:state-load', async (): Promise<LocalState> => {
    try {
      const data = fs.readFileSync(getStatePath(), 'utf-8');
      return JSON.parse(data) as LocalState;
    } catch {
      return {
        recentRulerPaths: [],
        recentWorkspacePaths: [],
        recentExportPaths: [],
      };
    }
  });

  // Save state
  ipcMain.handle(
    'csvons:state-save',
    async (_event, state: Partial<LocalState>): Promise<void> => {
      let current: LocalState;
      try {
        const data = fs.readFileSync(getStatePath(), 'utf-8');
        current = JSON.parse(data) as LocalState;
      } catch {
        current = {
          recentRulerPaths: [],
          recentWorkspacePaths: [],
          recentExportPaths: [],
        };
      }
      const merged = { ...current, ...state };
      // Keep only last 8 entries
      for (const key of [
        'recentRulerPaths',
        'recentWorkspacePaths',
        'recentExportPaths',
      ] as const) {
        if (merged[key]) {
          merged[key] = merged[key].slice(0, 8);
        }
      }
      fs.writeFileSync(getStatePath(), JSON.stringify(merged, null, 2));
    },
  );

  // Select file
  ipcMain.handle(
    'csvons:select-file',
    async (
      _event,
      options?: { filters?: { name: string; extensions: string[] }[] },
    ): Promise<string | null> => {
      const result = await dialog.showOpenDialog({
        properties: ['openFile'],
        filters: options?.filters,
      });
      return result.canceled ? null : (result.filePaths[0] ?? null);
    },
  );

  // Select directory
  ipcMain.handle(
    'csvons:select-directory',
    async (): Promise<string | null> => {
      const result = await dialog.showOpenDialog({
        properties: ['openDirectory'],
      });
      return result.canceled ? null : (result.filePaths[0] ?? null);
    },
  );

  // Export report
  ipcMain.handle(
    'csvons:export-report',
    async (
      _event,
      report: ValidationReport,
      filePath: string,
      format: 'json' | 'markdown',
    ): Promise<string> => {
      let content: string;
      let ext: string;

      if (format === 'json') {
        content = JSON.stringify(report, null, 2) + '\n';
        ext = '.json';
      } else {
        content = reportToMarkdown(report);
        ext = '.md';
      }

      // Normalize extension
      if (!filePath.toLowerCase().endsWith(ext)) {
        filePath += ext;
      }

      // Ensure parent directory exists
      const dir = path.dirname(filePath);
      fs.mkdirSync(dir, { recursive: true });

      fs.writeFileSync(filePath, content, 'utf-8');
      return filePath;
    },
  );
}

function reportToMarkdown(report: ValidationReport): string {
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
