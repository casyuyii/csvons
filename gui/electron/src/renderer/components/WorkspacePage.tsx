import { useState, useEffect } from 'react';
import { CsvPreviewCard } from './CsvPreviewCard.tsx';

declare global {
  interface Window {
    csvonsAPI: any;
  }
}

interface CsvFileEntry {
  name: string;
  path: string;
}

export function WorkspacePage() {
  const [workspacePath, setWorkspacePath] = useState('');
  const [loading, setLoading] = useState(false);
  const [csvFiles, setCsvFiles] = useState<CsvFileEntry[]>([]);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [preview, setPreview] = useState<{
    headers: string[];
    rows: string[][];
    totalRows: number;
  } | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [recentPaths, setRecentPaths] = useState<string[]>([]);

  useEffect(() => {
    loadState();
  }, []);

  async function loadState() {
    try {
      const state = await window.csvonsAPI.loadState();
      setRecentPaths(state.recentWorkspacePaths ?? []);
    } catch {}
  }

  async function handleSelectDirectory() {
    const path = await window.csvonsAPI.selectDirectory();
    if (path) setWorkspacePath(path);
  }

  async function handleScan() {
    if (!workspacePath) {
      setError('Please select a workspace directory');
      return;
    }

    setLoading(true);
    setError(null);
    setCsvFiles([]);
    setSelectedFile(null);
    setPreview(null);

    try {
      // Save recent path
      const newRecent = [
        workspacePath,
        ...recentPaths.filter((p) => p !== workspacePath),
      ].slice(0, 8);
      setRecentPaths(newRecent);
      await window.csvonsAPI.saveState({ recentWorkspacePaths: newRecent });

      // We use the Electron API to list directory contents
      // Since we can't directly list files from renderer, we'll scan for CSVs
      // using the csv-preview endpoint as a workaround, or we can add a dedicated IPC handler
      // For now, we'll prompt the user to select individual CSV files
      setError(
        'Workspace scanning requires selecting CSV files individually. Use Browse to select a CSV file for preview.',
      );
    } catch (err: any) {
      setError(err.message ?? 'Failed to scan workspace');
    } finally {
      setLoading(false);
    }
  }

  async function handleSelectCsvFile() {
    const filePath = await window.csvonsAPI.selectFile({
      filters: [{ name: 'CSV', extensions: ['csv'] }],
    });
    if (filePath) {
      setSelectedFile(filePath);
      loadPreview(filePath);
    }
  }

  async function loadPreview(filePath: string) {
    setPreviewLoading(true);
    setPreview(null);

    try {
      const result = await window.csvonsAPI.readCsvPreview(filePath, 6);
      setPreview(result);
    } catch (err: any) {
      setError(err.message ?? 'Failed to load preview');
    } finally {
      setPreviewLoading(false);
    }
  }

  return (
    <div style={{ maxWidth: 1100, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 16 }}>Workspace</h2>

      {/* Workspace path */}
      <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
        <input
          type="text"
          value={workspacePath}
          onChange={(e) => setWorkspacePath(e.target.value)}
          placeholder="Workspace directory path"
          style={{
            flex: 1,
            padding: '8px 12px',
            border: '1px solid #ccc',
            borderRadius: 4,
            fontSize: 14,
          }}
        />
        <button onClick={handleSelectDirectory} style={btnOutlined}>
          Browse...
        </button>
      </div>

      {/* Recent paths */}
      {recentPaths.length > 0 && (
        <div
          style={{
            display: 'flex',
            gap: 4,
            flexWrap: 'wrap',
            marginBottom: 12,
          }}
        >
          {recentPaths.slice(0, 4).map((p) => (
            <button
              key={p}
              onClick={() => setWorkspacePath(p)}
              style={chipStyle}
            >
              {p.split('/').pop()}
            </button>
          ))}
        </div>
      )}

      {/* Actions */}
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <button onClick={handleSelectCsvFile} style={btnFilled}>
          Open CSV File
        </button>
      </div>

      {error && (
        <div style={{ color: '#c62828', marginBottom: 12, fontSize: 14 }}>
          {error}
        </div>
      )}

      {/* CSV Preview */}
      {selectedFile && (
        <div>
          <h3 style={{ marginBottom: 8, fontSize: 16 }}>
            Preview: {selectedFile.split('/').pop()}
          </h3>
          {previewLoading ? (
            <div style={{ color: '#666' }}>Loading preview...</div>
          ) : preview ? (
            <CsvPreviewCard
              headers={preview.headers}
              rows={preview.rows}
              totalRows={preview.totalRows}
              filePath={selectedFile}
            />
          ) : (
            <div style={{ color: '#666' }}>No preview available</div>
          )}
        </div>
      )}
    </div>
  );
}

const btnFilled: React.CSSProperties = {
  padding: '10px 20px',
  background: '#3f51b5',
  color: '#fff',
  border: 'none',
  borderRadius: 4,
  cursor: 'pointer',
  fontSize: 14,
  fontWeight: 500,
};

const btnOutlined: React.CSSProperties = {
  padding: '8px 16px',
  background: 'transparent',
  color: '#3f51b5',
  border: '1px solid #3f51b5',
  borderRadius: 4,
  cursor: 'pointer',
  fontSize: 14,
};

const chipStyle: React.CSSProperties = {
  padding: '4px 10px',
  background: '#e8eaf6',
  border: 'none',
  borderRadius: 12,
  cursor: 'pointer',
  fontSize: 12,
  color: '#3f51b5',
};
