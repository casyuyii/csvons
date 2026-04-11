# Contract: Electron IPC API

## Preload API (window.csvonsAPI)

Exposed via `contextBridge.exposeInMainWorld('csvonsAPI', ...)`.

### validate(configPath: string): Promise<ValidationReport>
Runs validation via @csvons/core in main process.

### readCsvPreview(filePath: string, maxRows?: number): Promise<{headers: string[], rows: string[][]}>
Reads first N rows of a CSV file for preview display.

### loadState(): Promise<LocalState>
Reads persisted application state (recent paths).

### saveState(state: Partial<LocalState>): Promise<void>
Writes partial state update to persistence file.

### selectFile(options?: {filters?: {name: string, extensions: string[]}[]}): Promise<string | null>
Opens native file dialog. Returns selected path or null.

### selectDirectory(): Promise<string | null>
Opens native directory dialog. Returns selected path or null.

### exportReport(report: ValidationReport, filePath: string, format: 'json' | 'markdown'): Promise<string>
Exports validation report to file. Returns written file path.

## LocalState

```typescript
interface LocalState {
  recentRulerPaths: string[];     // Last 8 ruler.json paths
  recentWorkspacePaths: string[]; // Last 8 workspace directories
  recentExportPaths: string[];    // Last 8 export file paths
}
```
