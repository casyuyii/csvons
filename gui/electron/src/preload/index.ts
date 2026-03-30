import { contextBridge, ipcRenderer } from 'electron';

export interface CsvonsAPI {
  validate(configPath: string): Promise<any>;
  readCsvPreview(
    filePath: string,
    maxRows?: number,
  ): Promise<{ headers: string[]; rows: string[][]; totalRows: number }>;
  loadState(): Promise<any>;
  saveState(state: any): Promise<void>;
  selectFile(options?: {
    filters?: { name: string; extensions: string[] }[];
  }): Promise<string | null>;
  selectDirectory(): Promise<string | null>;
  exportReport(
    report: any,
    filePath: string,
    format: 'json' | 'markdown',
  ): Promise<string>;
}

const api: CsvonsAPI = {
  validate: (configPath) =>
    ipcRenderer.invoke('csvons:validate', configPath),
  readCsvPreview: (filePath, maxRows = 6) =>
    ipcRenderer.invoke('csvons:csv-preview', filePath, maxRows),
  loadState: () => ipcRenderer.invoke('csvons:state-load'),
  saveState: (state) => ipcRenderer.invoke('csvons:state-save', state),
  selectFile: (options) =>
    ipcRenderer.invoke('csvons:select-file', options),
  selectDirectory: () => ipcRenderer.invoke('csvons:select-directory'),
  exportReport: (report, filePath, format) =>
    ipcRenderer.invoke('csvons:export-report', report, filePath, format),
};

contextBridge.exposeInMainWorld('csvonsAPI', api);
